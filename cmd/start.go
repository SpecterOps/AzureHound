// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"

	"github.com/certmichelin/azurehound/v2/client/rest"
	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/panicrecovery"
	"github.com/certmichelin/azurehound/v2/pipeline"
)

const (
	BHEAuthSignature string = "bhesignature"
)

var ErrExceededRetryLimit = errors.New("exceeded max retry limit for ingest batch, proceeding with next batch...")

func init() {
	configs := append(config.AzureConfig, config.BloodHoundEnterpriseConfig...)
	configs = append(configs, config.CollectionConfig...)
	config.Init(startCmd, configs)
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:               "start",
	Short:             "Start Azure data collection service for BloodHound Enterprise",
	Run:               startCmdImpl,
	PersistentPreRunE: persistentPreRunE,
	SilenceUsage:      true,
}

func startCmdImpl(cmd *cobra.Command, args []string) {
	start(cmd.Context())
}

func start(ctx context.Context) {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	sigChan := make(chan os.Signal)
	go func() {
		stacktrace := make([]byte, 8192)
		for range sigChan {
			length := runtime.Stack(stacktrace, true)
			fmt.Println(string(stacktrace[:length]))
		}
	}()
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	if azClient := connectAndCreateClient(); azClient == nil {
		exit(fmt.Errorf("azClient is unexpectedly nil"))
	} else if bheInstance, err := url.Parse(config.BHEUrl.Value().(string)); err != nil {
		exit(fmt.Errorf("unable to parse BHE url: %w", err))
	} else if bheClient, err := newSigningHttpClient(BHEAuthSignature, config.BHETokenId.Value().(string), config.BHEToken.Value().(string), config.Proxy.Value().(string)); err != nil {
		exit(fmt.Errorf("failed to create new signing HTTP client: %w", err))
	} else if updatedClient, err := updateClient(ctx, *bheInstance, bheClient); err != nil {
		exit(fmt.Errorf("failed to update client: %w", err))
	} else if err := endOrphanedJob(ctx, *bheInstance, bheClient, updatedClient); err != nil {
		exit(fmt.Errorf("failed to end orphaned job: %w", err))
	} else {
		log.Info("connected successfully! waiting for jobs...")
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		var (
			jobQueued    sync.Mutex
			currentJobID atomic.Int64
		)

		for {
			select {
			case <-ticker.C:
				if jobID := currentJobID.Load(); jobID != 0 {
					log.V(1).Info("collection in progress...", "jobId", jobID)
					if err := checkin(ctx, *bheInstance, bheClient); err != nil {
						log.Error(err, "bloodhound enterprise service checkin failed")
					}
				} else if jobQueued.TryLock() {
					go func() {
						defer panicrecovery.PanicRecovery()
						defer jobQueued.Unlock()
						defer bheClient.CloseIdleConnections()
						defer azClient.CloseIdleConnections()

						ctx, stop := context.WithCancel(ctx)
						panicrecovery.HandleBubbledPanic(ctx, stop, log)

						log.V(2).Info("checking for available collection jobs")
						if jobs, err := getAvailableJobs(ctx, *bheInstance, bheClient); err != nil {
							log.Error(err, "unable to fetch available jobs for azurehound")
						} else {
							// Get only the jobs that have reached their execution time
							executableJobs := []models.ClientJob{}
							now := time.Now()
							for _, job := range jobs {
								if job.Status == models.JobStatusReady && job.ExecutionTime.Before(now) || job.ExecutionTime.Equal(now) {
									executableJobs = append(executableJobs, job)
								}
							}

							// Sort jobs in ascending order by execution time
							sort.Slice(executableJobs, func(i, j int) bool {
								return executableJobs[i].ExecutionTime.Before(executableJobs[j].ExecutionTime)
							})

							if len(executableJobs) == 0 {
								log.V(2).Info("there are no jobs for azurehound to complete at this time")
							} else {
								defer currentJobID.Store(0)
								queuedJobID := executableJobs[0].ID
								currentJobID.Store(int64(queuedJobID))
								// Notify BHE instance of job start
								if err := startJob(ctx, *bheInstance, bheClient, queuedJobID); err != nil {
									log.Error(err, "failed to start job, will retry on next heartbeat")
									return
								}

								start := time.Now()

								// Batch data out for ingestion
								stream := listAll(ctx, azClient)
								batches := pipeline.Batch(ctx.Done(), stream, config.ColBatchSize.Value().(int), 10*time.Second)
								hasIngestErr := ingest(ctx, *bheInstance, bheClient, batches)

								// Notify BHE instance of job end
								duration := time.Since(start)

								message := "Collection completed successfully"
								if hasIngestErr {
									message = "Collection completed with errors during ingest"
								}
								if err := endJob(ctx, *bheInstance, bheClient, models.JobStatusComplete, message); err != nil {
									log.Error(err, "failed to end job")
								} else {
									log.Info(message, "id", queuedJobID, "duration", duration.String())
								}
							}
						}
					}()
				}
			case <-ctx.Done():
				return
			}
		}
	}
}

func ingest(ctx context.Context, bheUrl url.URL, bheClient *http.Client, in <-chan []interface{}) bool {
	endpoint := bheUrl.ResolveReference(&url.URL{Path: "/api/v2/ingest"})

	var (
		hasErrors           = false
		maxRetries          = 3
		unrecoverableErrMsg = fmt.Sprintf("ending current ingest job due to unrecoverable error while requesting %v", endpoint)
	)

	for data := range pipeline.OrDone(ctx.Done(), in) {
		var (
			body bytes.Buffer
			gw   = gzip.NewWriter(&body)
		)

		ingestData := models.IngestRequest{
			Meta: models.Meta{
				Type: "azure",
			},
			Data: data,
		}

		err := json.NewEncoder(gw).Encode(ingestData)
		if err != nil {
			log.Error(err, unrecoverableErrMsg)
		}
		gw.Close()

		if req, err := http.NewRequestWithContext(ctx, "POST", endpoint.String(), &body); err != nil {
			log.Error(err, unrecoverableErrMsg)
			return true
		} else {
			req.Header.Set("User-Agent", constants.UserAgent())
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			for retry := 0; retry < maxRetries; retry++ {
				// No retries on regular err cases, only on HTTP 504 Gateway Timeout and HTTP 503 Service Unavailable
				if response, err := bheClient.Do(req); err != nil {
					if rest.IsClosedConnectionErr(err) {
						// try again on force closed connection
						log.Error(err, fmt.Sprintf("remote host force closed connection while requesting %s; attempt %d/%d; trying again", req.URL, retry+1, maxRetries))
						rest.ExponentialBackoff(retry)

						if retry == maxRetries-1 {
							log.Error(ErrExceededRetryLimit, "")
							hasErrors = true
						}

						continue
					}
					log.Error(err, unrecoverableErrMsg)
					return true
				} else if response.StatusCode == http.StatusGatewayTimeout || response.StatusCode == http.StatusServiceUnavailable || response.StatusCode == http.StatusBadGateway {
					serverError := fmt.Errorf("received server error %d while requesting %v; attempt %d/%d; trying again", response.StatusCode, endpoint, retry+1, maxRetries)
					log.Error(serverError, "")

					rest.ExponentialBackoff(retry)

					if retry == maxRetries-1 {
						log.Error(ErrExceededRetryLimit, "")
						hasErrors = true
					}
					if err := response.Body.Close(); err != nil {
						log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
					continue
				} else if response.StatusCode != http.StatusAccepted {
					if bodyBytes, err := io.ReadAll(response.Body); err != nil {
						log.Error(fmt.Errorf("received unexpected response code from %v: %s; failure reading response body", endpoint, response.Status), unrecoverableErrMsg)
					} else {
						log.Error(fmt.Errorf("received unexpected response code from %v: %s %s", req.URL, response.Status, bodyBytes), unrecoverableErrMsg)
					}
					if err := response.Body.Close(); err != nil {
						log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
					return true
				} else {
					if err := response.Body.Close(); err != nil {
						log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
				}
			}
		}
	}
	return hasErrors
}

// TODO: create/use a proper bloodhound client
func do(bheClient *http.Client, req *http.Request) (*http.Response, error) {
	var (
		res        *http.Response
		maxRetries = 3
	)

	// copy the bytes in case we need to retry the request
	if body, err := rest.CopyBody(req); err != nil {
		return nil, err
	} else {
		for retry := 0; retry < maxRetries; retry++ {
			// Reusing http.Request requires rewinding the request body
			// back to a working state
			if body != nil && retry > 0 {
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}

			if res, err = bheClient.Do(req); err != nil {
				if rest.IsClosedConnectionErr(err) {
					// try again on force closed connections
					log.Error(err, fmt.Sprintf("remote host force closed connection while requesting %s; attempt %d/%d; trying again", req.URL, retry+1, maxRetries))
					rest.ExponentialBackoff(retry)
					continue
				}
				// normal client error, dont attempt again
				return nil, err
			} else if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
				if res.StatusCode >= http.StatusInternalServerError {
					// Internal server error, backoff and try again.
					serverError := fmt.Errorf("received server error %d while requesting %v", res.StatusCode, req.URL)
					log.Error(serverError, fmt.Sprintf("attempt %d/%d; trying again", retry+1, maxRetries))

					rest.ExponentialBackoff(retry)
					continue
				}
				// bad request we do not need to retry
				var body json.RawMessage
				defer res.Body.Close()
				if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
					return nil, fmt.Errorf("received unexpected response code from %v: %s; failure reading response body", req.URL, res.Status)
				} else {
					return nil, fmt.Errorf("received unexpected response code from %v: %s %s", req.URL, res.Status, body)
				}
			} else {
				return res, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to complete request to url=%s; attempts=%d;", req.URL, maxRetries)
}

type basicResponse[T any] struct {
	Data T `json:"data"`
}

func getAvailableJobs(ctx context.Context, bheUrl url.URL, bheClient *http.Client) ([]models.ClientJob, error) {
	var (
		endpoint = bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/available"})
		response basicResponse[[]models.ClientJob]
	)

	if req, err := rest.NewRequest(ctx, "GET", endpoint, nil, nil, nil); err != nil {
		return nil, err
	} else if res, err := do(bheClient, req); err != nil {
		return nil, err
	} else {
		defer res.Body.Close()
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return nil, err
		} else {
			return response.Data, nil
		}
	}
}

func checkin(ctx context.Context, bheUrl url.URL, bheClient *http.Client) error {
	endpoint := bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/current"})

	if req, err := rest.NewRequest(ctx, "GET", endpoint, nil, nil, nil); err != nil {
		return err
	} else if res, err := do(bheClient, req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

func startJob(ctx context.Context, bheUrl url.URL, bheClient *http.Client, jobId int) error {
	log.Info("beginning collection job", "id", jobId)
	var (
		endpoint = bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/start"})
		body     = map[string]int{
			"id": jobId,
		}
	)

	if req, err := rest.NewRequest(ctx, "POST", endpoint, body, nil, nil); err != nil {
		return err
	} else if res, err := do(bheClient, req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

func endJob(ctx context.Context, bheUrl url.URL, bheClient *http.Client, status models.JobStatus, message string) error {
	endpoint := bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/end"})

	body := models.CompleteJobRequest{
		Status:  status.String(),
		Message: message,
	}

	if req, err := rest.NewRequest(ctx, "POST", endpoint, body, nil, nil); err != nil {
		return err
	} else if res, err := do(bheClient, req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

func updateClient(ctx context.Context, bheUrl url.URL, bheClient *http.Client) (*models.UpdateClientResponse, error) {
	var (
		endpoint = bheUrl.ResolveReference(&url.URL{Path: "/api/v2/clients/update"})
		response = basicResponse[models.UpdateClientResponse]{}
	)
	if addr, err := dial(bheUrl.String()); err != nil {
		return nil, err
	} else {
		// hostname is nice to have but we don't really need it
		hostname, _ := os.Hostname()

		body := models.UpdateClientRequest{
			Address:  addr,
			Hostname: hostname,
			Version:  constants.Version,
		}

		log.V(2).Info("updating client info", "info", body)

		if req, err := rest.NewRequest(ctx, "PUT", endpoint, body, nil, nil); err != nil {
			return nil, err
		} else if res, err := do(bheClient, req); err != nil {
			return nil, err
		} else {
			defer res.Body.Close()
			if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
				return nil, err
			} else {
				return &response.Data, nil
			}
		}
	}
}

func endOrphanedJob(ctx context.Context, bheUrl url.URL, bheClient *http.Client, updatedClient *models.UpdateClientResponse) error {
	if updatedClient.CurrentJob.Status == models.JobStatusRunning {
		log.Info("the service started with an orphaned job in progress, sending job completion notice...", "jobId", updatedClient.CurrentJobID)
		return endJob(ctx, bheUrl, bheClient, models.JobStatusFailed, "This job has been orphaned. Re-run collection for complete data.")
	} else {
		return nil
	}
}

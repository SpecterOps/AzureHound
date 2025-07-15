package bloodhound

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

	"github.com/bloodhoundad/azurehound/v2/client/rest"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/go-logr/logr"
)

const (
	BHEAuthSignature string = "bhesignature"
)

var ErrExceededRetryLimit = errors.New("exceeded max retry limit for ingest batch, proceeding with next batch...")

// BloodHoundClient represents the methods for interacting with an instance of BloodHound
type BloodHoundClient interface {
	SendRequest(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
	Ingest(ctx context.Context, in <-chan []interface{}) bool
	GetAvailableJobs(ctx context.Context) ([]models.ClientJob, error)
	Checkin(ctx context.Context) error
	StartJob(ctx context.Context, jobId int) error
	EndJob(ctx context.Context, status models.JobStatus, message string) error
	UpdateClient(ctx context.Context) (*models.UpdateClientResponse, error)
	EndOrphanedJob(ctx context.Context, updatedClient *models.UpdateClientResponse) error
}

// BHEClient implements the BloodHoundClient interface to communicate with a BloodHound Enterprise instance
type BHEClient struct {
	httpClient *http.Client
	bheUrl     url.URL
	log        logr.Logger
}

// NewBHEClient creates a new BloodHoundClient using the values from the application's config
// TODO: the values from the global config should be provided via arguments in the constructor
func NewBHEClient(bheUrl url.URL) (BloodHoundClient, error) {
	client, err := rest.NewHTTPClient(config.Proxy.Value().(string))
	if err != nil {
		return nil, err
	}

	client.Transport = signingTransport{
		base:      client.Transport,
		tokenId:   config.BHETokenId.Value().(string),
		token:     config.BHEToken.Value().(string),
		signature: BHEAuthSignature,
	}

	return &BHEClient{
		httpClient: client,
		bheUrl:     bheUrl,
	}, nil
}

// SendRequest sends a given request to the BHE instance. In the event of an error, 3 retries will be attempted
func (s *BHEClient) SendRequest(req *http.Request) (*http.Response, error) {
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

			if res, err = s.httpClient.Do(req); err != nil {
				if rest.IsClosedConnectionErr(err) {
					// try again on force closed connections
					s.log.Error(err, fmt.Sprintf("remote host force closed connection while requesting %s; attempt %d/%d; trying again", req.URL, retry+1, maxRetries))
					rest.ExponentialBackoff(retry)
					continue
				}
				// normal client error, dont attempt again
				return nil, err
			} else if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
				if res.StatusCode >= http.StatusInternalServerError {
					// Internal server error, backoff and try again.
					serverError := fmt.Errorf("received server error %d while requesting %v", res.StatusCode, req.URL)
					s.log.Error(serverError, fmt.Sprintf("attempt %d/%d; trying again", retry+1, maxRetries))

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

func (s *BHEClient) Ingest(ctx context.Context, in <-chan []interface{}) bool {
	endpoint := s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/ingest"})

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
			s.log.Error(err, unrecoverableErrMsg)
		}
		gw.Close()

		if req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), &body); err != nil {
			s.log.Error(err, unrecoverableErrMsg)
			return true
		} else {
			req.Header.Set("User-Agent", constants.UserAgent())
			req.Header.Set("Accept", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			for retry := 0; retry < maxRetries; retry++ {
				// No retries on regular err cases, only on HTTP 504 Gateway Timeout and HTTP 503 Service Unavailable
				if response, err := s.httpClient.Do(req); err != nil {
					if rest.IsClosedConnectionErr(err) {
						// try again on force closed connection
						s.log.Error(err, fmt.Sprintf("remote host force closed connection while requesting %s; attempt %d/%d; trying again", req.URL, retry+1, maxRetries))
						rest.ExponentialBackoff(retry)

						if retry == maxRetries-1 {
							s.log.Error(ErrExceededRetryLimit, "")
							hasErrors = true
						}

						continue
					}
					s.log.Error(err, unrecoverableErrMsg)
					return true
				} else if response.StatusCode == http.StatusGatewayTimeout || response.StatusCode == http.StatusServiceUnavailable || response.StatusCode == http.StatusBadGateway {
					serverError := fmt.Errorf("received server error %d while requesting %v; attempt %d/%d; trying again", response.StatusCode, endpoint, retry+1, maxRetries)
					s.log.Error(serverError, "")

					rest.ExponentialBackoff(retry)

					if retry == maxRetries-1 {
						s.log.Error(ErrExceededRetryLimit, "")
						hasErrors = true
					}
					if err := response.Body.Close(); err != nil {
						s.log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
					continue
				} else if response.StatusCode != http.StatusAccepted {
					if bodyBytes, err := io.ReadAll(response.Body); err != nil {
						s.log.Error(fmt.Errorf("received unexpected response code from %v: %s; failure reading response body", endpoint, response.Status), unrecoverableErrMsg)
					} else {
						s.log.Error(fmt.Errorf("received unexpected response code from %v: %s %s", req.URL, response.Status, bodyBytes), unrecoverableErrMsg)
					}
					if err := response.Body.Close(); err != nil {
						s.log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
					return true
				} else {
					if err := response.Body.Close(); err != nil {
						s.log.Error(fmt.Errorf("failed to close ingest body: %w", err), unrecoverableErrMsg)
					}
				}
			}
		}
	}
	return hasErrors
}

// GetAvailableJobs sends a request to BHE to get the list of available jobs
func (s *BHEClient) GetAvailableJobs(ctx context.Context) ([]models.ClientJob, error) {
	var (
		endpoint = s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/available"})
		response bloodhoundResponse[[]models.ClientJob]
	)

	if req, err := rest.NewRequest(ctx, "GET", endpoint, nil, nil, nil); err != nil {
		return nil, err
	} else if res, err := s.SendRequest(req); err != nil {
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

// Checkin sends a request to BHE indicating that the client is running
func (s *BHEClient) Checkin(ctx context.Context) error {
	endpoint := s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/current"})

	if req, err := rest.NewRequest(ctx, "GET", endpoint, nil, nil, nil); err != nil {
		return err
	} else if res, err := s.SendRequest(req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

// StartJob sends a request to BHE instructing it to start a job
func (s *BHEClient) StartJob(ctx context.Context, jobId int) error {
	s.log.Info("beginning collection job", "id", jobId)
	var (
		endpoint = s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/start"})
		body     = map[string]int{
			"id": jobId,
		}
	)

	if req, err := rest.NewRequest(ctx, "POST", endpoint, body, nil, nil); err != nil {
		return err
	} else if res, err := s.SendRequest(req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

// EndJob sends a request to BHE instructing it to end a job
func (s *BHEClient) EndJob(ctx context.Context, status models.JobStatus, message string) error {
	endpoint := s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/jobs/end"})

	body := models.CompleteJobRequest{
		Status:  status.String(),
		Message: message,
	}

	if req, err := rest.NewRequest(ctx, "POST", endpoint, body, nil, nil); err != nil {
		return err
	} else if res, err := s.SendRequest(req); err != nil {
		return err
	} else {
		res.Body.Close()
		return nil
	}
}

// UpdateClient sends a request to BHE and updates the AzureHound client info
func (s *BHEClient) UpdateClient(ctx context.Context) (*models.UpdateClientResponse, error) {
	var (
		endpoint = s.bheUrl.ResolveReference(&url.URL{Path: "/api/v2/clients/update"})
		response = bloodhoundResponse[models.UpdateClientResponse]{}
	)
	if addr, err := rest.Dial(s.log, s.bheUrl.String()); err != nil {
		return nil, err
	} else {
		// hostname is nice to have, but we don't really need it
		hostname, _ := os.Hostname()

		body := models.UpdateClientRequest{
			Address:  addr,
			Hostname: hostname,
			Version:  constants.Version,
		}

		s.log.V(2).Info("updating client info", "info", body)

		if req, err := rest.NewRequest(ctx, "PUT", endpoint, body, nil, nil); err != nil {
			return nil, err
		} else if res, err := s.SendRequest(req); err != nil {
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

// EndOrphanedJob if a job is running, sends a request to BHE to end the current job with a failed status
func (s *BHEClient) EndOrphanedJob(ctx context.Context, updatedClient *models.UpdateClientResponse) error {
	if updatedClient.CurrentJob.Status == models.JobStatusRunning {
		s.log.Info("the service started with an orphaned job in progress, sending job completion notice...", "jobId", updatedClient.CurrentJobID)
		return s.EndJob(ctx, models.JobStatusFailed, "This job has been orphaned. Re-run collection for complete data.")
	} else {
		return nil
	}
}

// CloseIdleConnections closes all idle connections on the internal http.Client
func (s *BHEClient) CloseIdleConnections() {
	s.httpClient.CloseIdleConnections()
}

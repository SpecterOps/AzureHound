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
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/bloodhound"
	"github.com/spf13/cobra"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
)

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

	if azClient := connectAndCreateClient(); azClient == nil {
		exit(fmt.Errorf("azClient is unexpectedly nil"))
	} else if bheInstance, err := url.Parse(config.BHEUrl.Value().(string)); err != nil {
		exit(fmt.Errorf("unable to parse BHE url: %w", err))
	} else if bheClient, err := bloodhound.NewBHEClient(*bheInstance, config.BHETokenId.Value().(string), config.BHEToken.Value().(string), config.Proxy.Value().(string), config.BHEMaxReqPerConn.Value().(int), 3, log); err != nil {
		exit(fmt.Errorf("failed to create new signing HTTP client: %w", err))
	} else if updatedClient, err := bheClient.UpdateClient(ctx); err != nil {
		exit(fmt.Errorf("failed to update client: %w", err))
	} else if err := bheClient.EndOrphanedJob(ctx, updatedClient); err != nil {
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
					if err := bheClient.Checkin(ctx); err != nil {
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
						if jobs, err := bheClient.GetAvailableJobs(ctx); err != nil {
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
								if err := bheClient.StartJob(ctx, queuedJobID); err != nil {
									log.Error(err, "failed to start job, will retry on next heartbeat")
									return
								}

								start := time.Now()

								// Batch data out for ingestion
								stream := listAll(ctx, azClient)
								batches := pipeline.Batch(ctx.Done(), stream, config.ColBatchSize.Value().(int), 10*time.Second)
								hasIngestErr := bheClient.Ingest(ctx, batches)

								// Notify BHE instance of job end
								duration := time.Since(start)

								message := "Collection completed successfully"
								if hasIngestErr {
									message = "Collection completed with errors during ingest"
								}
								if err := bheClient.EndJob(ctx, models.JobStatusComplete, message); err != nil {
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

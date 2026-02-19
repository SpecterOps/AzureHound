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
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listAppFICCmd)
}

var listAppFICCmd = &cobra.Command{
	Use:          "appfics",
	Long:         "Lists Entra ID Application Federated Identity Credentials",
	Run:          listAppFICsCmdImpl,
	SilenceUsage: true,
}

func listAppFICsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure app federated identity credentials...")
	start := time.Now()
	stream := listAppFICs(ctx, azClient, listApps(ctx, azClient))
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listAppFICs(ctx context.Context, client client.AzureClient, apps <-chan azureWrapper[models.App]) <-chan azureWrapper[models.AppFICs] {
	var (
		out     = make(chan azureWrapper[models.AppFICs])
		streams = pipeline.Demux(ctx.Done(), apps, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
		params  = query.GraphParams{}
	)

	wg.Add(len(streams))
	for i := range streams {
		stream := streams[i]
		go func() {
			defer panicrecovery.PanicRecovery()
			defer wg.Done()
			for app := range stream {
				var (
					data = models.AppFICs{
						AppId: app.Data.AppId,
					}
					count = 0
				)
				for item := range client.ListAzureADAppFICs(ctx, app.Data.Id, params) {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing federated identity credentials for this app", "appId", app.Data.AppId)
					} else {
						appFIC := models.AppFIC{
							FIC:   item.Ok,
							AppId: app.Data.AppId,
						}
						log.V(2).Info("found app FIC", "appFic", appFIC)
						count++
						data.FICs = append(data.FICs, appFIC)
					}
				}

				if data.FICs != nil {
					if ok := pipeline.Send(ctx.Done(), out, NewAzureWrapper(
						enums.KindAZFederatedIdentityCredential,
						data,
					)); !ok {
						return
					}
				}
				log.V(1).Info("finished listing app fics", "appId", app.Data.AppId, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all app fics")
	}()

	return out
}

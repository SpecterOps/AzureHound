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
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listServicePrincipalsCmd)
}

var listServicePrincipalsCmd = &cobra.Command{
	Use:          "service-principals",
	Long:         "Lists Azure Active Directory Service Principals",
	Run:          listServicePrincipalsCmdImpl,
	SilenceUsage: true,
}

func listServicePrincipalsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory service principals...")
	start := time.Now()
	stream := listServicePrincipals(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listServicePrincipals(ctx context.Context, client client.AzureClient) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)
		count := 0
		for item := range client.ListAzureADServicePrincipals(ctx, query.GraphParams{}) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing service principals")
				return
			} else {
				log.V(2).Info("found service principal", "name", item.Ok.DisplayName)
				count++
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZServicePrincipal,
					Data: models.ServicePrincipal{
						ServicePrincipal: item.Ok,
						TenantId:         client.TenantInfo().TenantId,
						TenantName:       client.TenantInfo().DisplayName,
					},
				}); !ok {
					return
				}
			}
		}
		log.Info("finished listing all service principals", "count", count)
	}()

	return out
}

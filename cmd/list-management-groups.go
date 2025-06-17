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
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listManagementGroupsCmd)
}

var listManagementGroupsCmd = &cobra.Command{
	Use:          "management-groups",
	Long:         "Lists Azure Active Directory Management Groups",
	Run:          listManagementGroupsCmdImpl,
	SilenceUsage: true,
}

func listManagementGroupsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory management groups...")
	start := time.Now()

	stream := listManagementGroups(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listManagementGroups(ctx context.Context, client client.AzureClient) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)
		count := 0
		for item := range client.ListAzureManagementGroups(ctx, "") {
			if item.Error != nil {
				log.Info("warning: unable to process azure management groups; either the organization has no management groups or azurehound does not have the reader role on the root management group.")
				return
			} else if len(config.AzMgmtGroupId.Value().([]string)) == 0 || contains(config.AzMgmtGroupId.Value().([]string), item.Ok.Name) {
				log.V(2).Info("found management group", "managementGroup", item)
				count++
				mgmtGroup := models.ManagementGroup{
					ManagementGroup: item.Ok,
					TenantId:        item.Ok.Properties.TenantId,
				}

				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZManagementGroup,
					Data: mgmtGroup,
				}); !ok {
					return
				}
			}
		}
		log.Info("finished listing all management groups", "count", count)
	}()

	return out
}

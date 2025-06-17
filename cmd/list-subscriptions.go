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
	"os"
	"os/signal"
	"time"

	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listSubscriptionsCmd)
}

var listSubscriptionsCmd = &cobra.Command{
	Use:          "subscriptions",
	Long:         "Lists Azure Active Directory Subscriptions",
	Run:          listSubscriptionsCmdImpl,
	SilenceUsage: true,
}

func listSubscriptionsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory subscriptions...")
	start := time.Now()
	stream := listSubscriptions(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listSubscriptions(ctx context.Context, client client.AzureClient) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)
		var (
			count                = 0
			selectedSubIds       = config.AzSubId.Value().([]string)
			selectedMgmtGroupIds = config.AzMgmtGroupId.Value().([]string)
			filterOnSubs         = len(selectedSubIds) != 0 || len(selectedMgmtGroupIds) != 0
		)

		if len(selectedMgmtGroupIds) != 0 {
			descendantChannel := listManagementGroupDescendants(ctx, client, listManagementGroups(ctx, client))
			for i := range descendantChannel {
				if item, ok := i.(AzureWrapper).Data.(azure.DescendantInfo); !ok {
					log.Error(fmt.Errorf("failed type assertion"), "unable to continue evaluating management group descendants", "result", i)
					return
				} else if item.Type == "Microsoft.Management/managementGroups/subscriptions" {
					selectedSubIds = append(selectedSubIds, item.Name)
				}
			}
		}
		uniqueSubIds := unique(selectedSubIds)

		for item := range client.ListAzureSubscriptions(ctx) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing subscriptions")
				return
			} else if !filterOnSubs || contains(uniqueSubIds, item.Ok.SubscriptionId) {
				log.V(2).Info("found subscription", "subscription", item)
				count++
				// the embedded struct's values override top-level properties so TenantId
				// needs to be explicitly set.
				data := models.Subscription{
					Subscription: item.Ok,
				}
				data.TenantId = item.Ok.TenantId
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZSubscription,
					Data: data,
				}); !ok {
					return
				}
			}
		}
		log.Info("finished listing all subscriptions", "count", count)
	}()

	return out
}

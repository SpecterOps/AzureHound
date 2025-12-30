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
	"sync"
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
	listRootCmd.AddCommand(listStorageContainersCmd)
}

var listStorageContainersCmd = &cobra.Command{
	Use:          "storage-containers",
	Long:         "Lists Azure Storage Containers",
	Run:          listStorageContainersCmdImpl,
	SilenceUsage: true,
}

func listStorageContainersCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure storage containers...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	storageAccounts := listStorageAccounts(ctx, azClient, subscriptions)
	stream := listStorageContainers(ctx, azClient, storageAccounts)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listStorageContainers(ctx context.Context, client client.AzureClient, storageAccounts <-chan interface{}) <-chan interface{} {
	var (
		out = make(chan interface{})
		ids = make(chan interface{})
		// The original size of the demuxxer cascaded into error messages for a lot of collection steps.
		// Decreasing the demuxxer size only here is sufficient to prevent the cascade
		// The error message with higher values for size is
		// "The request was throttled."
		// See issue #7: https://github.com/bloodhoundad/azurehound/issues/7
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)
		for result := range pipeline.OrDone(ctx.Done(), storageAccounts) {
			if storageAccount, ok := result.(AzureWrapper).Data.(models.StorageAccount); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating storage containers", "result", result)
				return
			} else {
				if ok := pipeline.SendAny(ctx.Done(), ids, storageAccount); !ok {
					return
				}
			}
		}
	}()

	wg.Add(len(streams))
	for i := range streams {
		stream := streams[i]
		go func() {
			defer panicrecovery.PanicRecovery()
			defer wg.Done()
			for stAccount := range stream {
				count := 0
				for item := range client.ListAzureStorageContainers(ctx, stAccount.(models.StorageAccount).SubscriptionId, stAccount.(models.StorageAccount).ResourceGroupName, stAccount.(models.StorageAccount).Name, "", "deleted", "") {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing storage containers for this subscription", "subscriptionId", stAccount.(models.StorageAccount).SubscriptionId, "storageAccountName", stAccount.(models.StorageAccount).Name)
					} else {
						storageContainer := models.StorageContainer{
							StorageContainer:  item.Ok,
							StorageAccountId:  stAccount.(models.StorageAccount).StorageAccount.Id,
							SubscriptionId:    "/subscriptions/" + stAccount.(models.StorageAccount).SubscriptionId,
							ResourceGroupId:   item.Ok.ResourceGroupId(),
							ResourceGroupName: item.Ok.ResourceGroupName(),
							TenantId:          client.TenantInfo().TenantId,
						}
						log.V(2).Info("found storage container", "id", storageContainer.Id, "name", storageContainer.Name, "tenant", storageContainer.TenantId)
						count++
						if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
							Kind: enums.KindAZStorageContainer,
							Data: storageContainer,
						}); !ok {
							return
						}
					}
					log.V(1).Info("finished listing storage containers", "subscriptionId", stAccount.(models.StorageAccount).SubscriptionId, "count", count)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all storage containers")
	}()

	return out
}

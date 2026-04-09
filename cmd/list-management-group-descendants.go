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
	listRootCmd.AddCommand(listManagementGroupDescendantsCmd)
}

var listManagementGroupDescendantsCmd = &cobra.Command{
	Use:          "management-group-descendants",
	Long:         "Lists Azure Management Group Descendants",
	Run:          listManagementGroupDescendantsCmdImpl,
	SilenceUsage: true,
}

func listManagementGroupDescendantsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure management group descendants...")
	start := time.Now()
	stream := listManagementGroupDescendants(ctx, azClient, listManagementGroups(ctx, azClient))
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listManagementGroupDescendants(ctx context.Context, client client.AzureClient, managementGroups <-chan interface{}) <-chan interface{} {
	var (
		out     = make(chan interface{})
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), managementGroups) {
			if managementGroup, ok := result.(AzureWrapper).Data.(models.ManagementGroup); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating management group descendants", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, managementGroup.Name); !ok {
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
			for id := range stream {
				count := 0
				for item := range client.ListAzureManagementGroupDescendants(ctx, id, 3000) {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing descendants for this management group", "managementGroupId", id)
					} else {
						log.V(2).Info("found management group descendant", "name", item.Ok.Name)
						count++
						if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
							Kind: enums.KindAZManagementGroupDescendant,
							Data: item.Ok,
						}); !ok {
							return
						}
					}
				}
				log.V(1).Info("finished listing management group descendants", "managementGroupId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all management group descendants")
	}()

	return out
}

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
	listRootCmd.AddCommand(listSubscriptionRoleAssignmentsCmd)
}

var listSubscriptionRoleAssignmentsCmd = &cobra.Command{
	Use:          "subscription-role-assignments",
	Long:         "Lists Subscription Role Assignments",
	Run:          listSubscriptionRoleAssignmentsCmdImpl,
	SilenceUsage: true,
}

func listSubscriptionRoleAssignmentsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure subscription role assignments...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	stream := listSubscriptionRoleAssignments(ctx, azClient, subscriptions)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listSubscriptionRoleAssignments(ctx context.Context, client client.AzureClient, subscriptions <-chan interface{}) <-chan interface{} {
	var (
		out     = make(chan interface{})
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), subscriptions) {
			if subscription, ok := result.(AzureWrapper).Data.(models.Subscription); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating subscription role assignments", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, subscription.Id); !ok {
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
				var (
					subscriptionRoleAssignments = models.SubscriptionRoleAssignments{
						SubscriptionId: id,
					}
					count = 0
				)
				for item := range client.ListRoleAssignmentsForResource(ctx, id, "atScope()", "") {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing role assignments for this subscription", "subscriptionId", id)
					} else {
						subscriptionRoleAssignment := models.SubscriptionRoleAssignment{
							SubscriptionId: id,
							RoleAssignment: item.Ok,
						}
						log.V(2).Info("found subscription role assignment", "id", subscriptionRoleAssignment.RoleAssignment.Id, "name", subscriptionRoleAssignment.RoleAssignment.Name)
						count++
						subscriptionRoleAssignments.RoleAssignments = append(subscriptionRoleAssignments.RoleAssignments, subscriptionRoleAssignment)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZSubscriptionRoleAssignment,
					Data: subscriptionRoleAssignments,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing subscription role assignments", "subscriptionId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all subscription role assignments")
	}()

	return out
}

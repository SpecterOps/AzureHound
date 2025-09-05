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
	"path"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listSubscriptionOwnersCmd)
}

var listSubscriptionOwnersCmd = &cobra.Command{
	Use:          "subscription-owners",
	Long:         "Lists Azure Subscription Owners",
	Run:          listSubscriptionOwnersCmdImpl,
	SilenceUsage: true,
}

func listSubscriptionOwnersCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure subscription owners...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	roleAssignments := listSubscriptionRoleAssignments(ctx, azClient, subscriptions)
	stream := listSubscriptionOwners(ctx, azClient, roleAssignments)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listSubscriptionOwners(ctx context.Context, _ client.AzureClient, roleAssignments <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		for result := range pipeline.OrDone(ctx.Done(), roleAssignments) {
			if roleAssignments, ok := result.(AzureWrapper).Data.(models.SubscriptionRoleAssignments); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating subscription owners", "result", result)
				return
			} else {
				var (
					subscriptionOwners = models.SubscriptionOwners{
						SubscriptionId: roleAssignments.SubscriptionId,
					}
					count = 0
				)
				for _, item := range roleAssignments.RoleAssignments {
					roleDefinitionId := path.Base(item.RoleAssignment.Properties.RoleDefinitionId)

					if roleDefinitionId == constants.OwnerRoleID {
						subscriptionOwner := models.SubscriptionOwner{
							Owner:          item.RoleAssignment,
							SubscriptionId: item.SubscriptionId,
						}
						log.V(2).Info("found subscription owner", "subscriptionOwner", subscriptionOwner)
						count++
						subscriptionOwners.Owners = append(subscriptionOwners.Owners, subscriptionOwner)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZSubscriptionOwner,
					Data: subscriptionOwners,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing subscription owners", "subscriptionId", roleAssignments.SubscriptionId, "count", count)
			}
		}
	}()

	return out
}

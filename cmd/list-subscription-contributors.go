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
	listRootCmd.AddCommand(listSubscriptionContributorsCmd)
}

var listSubscriptionContributorsCmd = &cobra.Command{
	Use:          "subscription-contributors",
	Long:         "Lists Azure Subscription Contributors",
	Run:          listSubscriptionContributorsCmdImpl,
	SilenceUsage: true,
}

func listSubscriptionContributorsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure subscription contributors...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	roleAssignments := listSubscriptionRoleAssignments(ctx, azClient, subscriptions)
	stream := listSubscriptionContributors(ctx, azClient, roleAssignments)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listSubscriptionContributors(ctx context.Context, client client.AzureClient, roleAssignments <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		for result := range pipeline.OrDone(ctx.Done(), roleAssignments) {
			if roleAssignments, ok := result.(AzureWrapper).Data.(models.SubscriptionRoleAssignments); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating subscription contributors", "result", result)
				return
			} else {
				var (
					subscriptionContributors = models.SubscriptionContributors{
						SubscriptionId: roleAssignments.SubscriptionId,
					}
					count = 0
				)
				for _, item := range roleAssignments.RoleAssignments {
					roleDefinitionId := path.Base(item.RoleAssignment.Properties.RoleDefinitionId)

					if roleDefinitionId == constants.ContributorRoleID {
						subscriptionContributor := models.SubscriptionContributor{
							Contributor:    item.RoleAssignment,
							SubscriptionId: item.SubscriptionId,
						}
						log.V(2).Info("found subscription contributor", "name", subscriptionContributor.Contributor.Name)
						count++
						subscriptionContributors.Contributors = append(subscriptionContributors.Contributors, subscriptionContributor)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZSubscriptionContributor,
					Data: subscriptionContributors,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing subscription contributors", "subscriptionId", roleAssignments.SubscriptionId, "count", count)
			}
		}
	}()

	return out
}


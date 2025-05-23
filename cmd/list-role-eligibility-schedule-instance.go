// Copyright (C) 2025 Specter Ops, Inc.
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
	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

func init() {
	listRootCmd.AddCommand(listUnifiedRoleEligibilityScheduleInstanceCmd)
}

var listUnifiedRoleEligibilityScheduleInstanceCmd = &cobra.Command{
	Use:          "unified-role-eligibility-schedule-instances",
	Long:         "Lists Unified Role Eligibility Schedule Instances",
	SilenceUsage: true,
	Run:          listUnifiedRoleEligibilityScheduleInstancesCmdImpl,
}

func listUnifiedRoleEligibilityScheduleInstancesCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	azClient := connectAndCreateClient()
	log.V(1).Info("collecting azure unified role eligibility schedule instances")
	start := time.Now()
	stream := listRoleEligibilityScheduleInstances(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.V(1).Info("collection completed", "duration", duration.String())
}

func listRoleEligibilityScheduleInstances(ctx context.Context, client client.AzureClient) <-chan interface{} {
	var (
		out = make(chan interface{})
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)
		count := 0

		for item := range client.ListAzureUnifiedRoleEligibilityScheduleInstances(ctx, query.GraphParams{}) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing unified role eligibility instance schedules")
				return
			} else {
				log.V(2).Info("found unified role eligibility instance schedule", "unifiedRoleEligibilitySchedule", item)
				count++
				result := item.Ok
				if ok := pipeline.SendAny(ctx.Done(), out, azureWrapper[models.RoleEligibilityScheduleInstance]{
					Kind: enums.KindAZRoleEligibilityScheduleInstance,
					Data: models.RoleEligibilityScheduleInstance{
						Id:               result.Id,
						RoleDefinitionId: result.RoleDefinitionId,
						PrincipalId:      result.PrincipalId,
						DirectoryScopeId: result.DirectoryScopeId,
						StartDateTime:    result.StartDateTime,
						TenantId:         client.TenantInfo().TenantId,
					},
				}); !ok {
					return
				}
			}
		}
		log.V(1).Info("finished listing unified role eligibility schedule instances", "count", count)
	}()

	return out
}

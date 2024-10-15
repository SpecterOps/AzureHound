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

	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/enums"
	"github.com/certmichelin/azurehound/v2/internal"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/panicrecovery"
	"github.com/certmichelin/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listManagementGroupOwnersCmd)
}

var listManagementGroupOwnersCmd = &cobra.Command{
	Use:          "management-group-owners",
	Long:         "Lists Azure Management Group Owners",
	Run:          listManagementGroupOwnersCmdImpl,
	SilenceUsage: true,
}

func listManagementGroupOwnersCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure management group owners...")
	start := time.Now()
	managementGroups := listManagementGroups(ctx, azClient)
	roleAssignments := listManagementGroupRoleAssignments(ctx, azClient, managementGroups)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	stream := listManagementGroupOwners(ctx, roleAssignments)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listManagementGroupOwners(
	ctx context.Context,
	roleAssignments <-chan azureWrapper[models.ManagementGroupRoleAssignments],
) <-chan any {
	return pipeline.Map(ctx.Done(), roleAssignments, func(ra azureWrapper[models.ManagementGroupRoleAssignments]) any {
		filteredAssignments := internal.Filter(ra.Data.RoleAssignments, mgmtGroupRoleAssignmentFilter(constants.OwnerRoleID))
		owners := internal.Map(filteredAssignments, func(ra models.ManagementGroupRoleAssignment) models.ManagementGroupOwner {
			return models.ManagementGroupOwner{
				Owner:             ra.RoleAssignment,
				ManagementGroupId: ra.ManagementGroupId,
			}
		})
		return NewAzureWrapper(enums.KindAZManagementGroupOwner, models.ManagementGroupOwners{
			ManagementGroupId: ra.Data.ManagementGroupId,
			Owners:            owners,
		})
	})
}

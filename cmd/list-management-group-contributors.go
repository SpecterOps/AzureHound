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

	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/internal"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listManagementGroupContributorsCmd)
}

var listManagementGroupContributorsCmd = &cobra.Command{
	Use:          "management-group-contributors",
	Long:         "Lists Azure Management Group Contributors",
	Run:          listManagementGroupContributorsCmdImpl,
	SilenceUsage: true,
}

func listManagementGroupContributorsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure management group contributors...")
	start := time.Now()
	managementGroups := listManagementGroups(ctx, azClient)
	roleAssignments := listManagementGroupRoleAssignments(ctx, azClient, managementGroups)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	stream := listManagementGroupContributors(ctx, roleAssignments)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listManagementGroupContributors(
	ctx context.Context,
	roleAssignments <-chan azureWrapper[models.ManagementGroupRoleAssignments],
) <-chan any {
	return pipeline.Map(ctx.Done(), roleAssignments, func(ra azureWrapper[models.ManagementGroupRoleAssignments]) any {
		filteredAssignments := internal.Filter(ra.Data.RoleAssignments, mgmtGroupRoleAssignmentFilter(constants.ContributorRoleID))
		contributors := internal.Map(filteredAssignments, func(ra models.ManagementGroupRoleAssignment) models.ManagementGroupContributor {
			return models.ManagementGroupContributor{
				Contributor:       ra.RoleAssignment,
				ManagementGroupId: ra.ManagementGroupId,
			}
		})
		return NewAzureWrapper(enums.KindAZManagementGroupContributor, models.ManagementGroupContributors{
			ManagementGroupId: ra.Data.ManagementGroupId,
			Contributors:      contributors,
		})
	})
}


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
	listRootCmd.AddCommand(listVirtualMachineContributorsCmd)
}

var listVirtualMachineContributorsCmd = &cobra.Command{
	Use:          "virtual-machine-contributors",
	Long:         "Lists Azure Virtual Machine Contributors",
	Run:          listVirtualMachineContributorsCmdImpl,
	SilenceUsage: true,
}

func listVirtualMachineContributorsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure virtual machine contributors...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	vms := listVirtualMachines(ctx, azClient, subscriptions)
	vmRoleAssignments := listVirtualMachineRoleAssignments(ctx, azClient, vms)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	stream := listVirtualMachineContributors(ctx, vmRoleAssignments)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listVirtualMachineContributors(
	ctx context.Context,
	roleAssignments <-chan azureWrapper[models.VirtualMachineRoleAssignments],
) <-chan any {
	return pipeline.Map(ctx.Done(), roleAssignments, func(ra azureWrapper[models.VirtualMachineRoleAssignments]) any {
		filteredAssignments := internal.Filter(ra.Data.RoleAssignments, vmRoleAssignmentFilter(constants.ContributorRoleID))
		contributors := internal.Map(filteredAssignments, func(ra models.VirtualMachineRoleAssignment) models.VirtualMachineContributor {
			return models.VirtualMachineContributor{
				VirtualMachineId: ra.VirtualMachineId,
				Contributor:      ra.RoleAssignment,
			}
		})
		return NewAzureWrapper(enums.KindAZVMContributor, models.VirtualMachineContributors{
			VirtualMachineId: ra.Data.VirtualMachineId,
			Contributors:     contributors,
		})
	})
}

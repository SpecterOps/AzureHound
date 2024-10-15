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
	listRootCmd.AddCommand(listVirtualMachineAvereContributorsCmd)
}

var listVirtualMachineAvereContributorsCmd = &cobra.Command{
	Use:          "virtual-machine-avere-contributors",
	Long:         "Lists Azure Virtual Machine Avere Contributors",
	Run:          listVirtualMachineAvereContributorsCmdImpl,
	SilenceUsage: true,
}

func listVirtualMachineAvereContributorsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure virtual machine averecontributors...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	vms := listVirtualMachines(ctx, azClient, subscriptions)
	vmRoleAssignments := listVirtualMachineRoleAssignments(ctx, azClient, vms)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	stream := listVirtualMachineAvereContributors(ctx, vmRoleAssignments)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listVirtualMachineAvereContributors(
	ctx context.Context,
	roleAssignments <-chan azureWrapper[models.VirtualMachineRoleAssignments],
) <-chan any {
	return pipeline.Map(ctx.Done(), roleAssignments, func(ra azureWrapper[models.VirtualMachineRoleAssignments]) any {
		filteredAssignments := internal.Filter(ra.Data.RoleAssignments, vmRoleAssignmentFilter(constants.AvereContributorRoleID))
		avereContributors := internal.Map(filteredAssignments, func(ra models.VirtualMachineRoleAssignment) models.VirtualMachineAvereContributor {
			return models.VirtualMachineAvereContributor{
				VirtualMachineId: ra.VirtualMachineId,
				AvereContributor: ra.RoleAssignment,
			}
		})
		return NewAzureWrapper(enums.KindAZVMAvereContributor, models.VirtualMachineAvereContributors{
			VirtualMachineId:  ra.Data.VirtualMachineId,
			AvereContributors: avereContributors,
		})
	})
}

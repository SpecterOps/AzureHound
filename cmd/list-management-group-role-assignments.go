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

	"github.com/certmichelin/azurehound/v2/client"
	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/enums"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/panicrecovery"
	"github.com/certmichelin/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listManagementGroupRoleAssignmentsCmd)
}

var listManagementGroupRoleAssignmentsCmd = &cobra.Command{
	Use:          "management-group-role-assignments",
	Long:         "Lists Management Group Role Assignments",
	Run:          listManagementGroupRoleAssignmentsCmdImpl,
	SilenceUsage: true,
}

func listManagementGroupRoleAssignmentsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure management group role assignments...")
	start := time.Now()
	managementGroups := listManagementGroups(ctx, azClient)
	stream := listManagementGroupRoleAssignments(ctx, azClient, managementGroups)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listManagementGroupRoleAssignments(ctx context.Context, client client.AzureClient, managementGroups <-chan interface{}) <-chan azureWrapper[models.ManagementGroupRoleAssignments] {
	var (
		out     = make(chan azureWrapper[models.ManagementGroupRoleAssignments])
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), managementGroups) {
			if managementGroup, ok := result.(AzureWrapper).Data.(models.ManagementGroup); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating management group role assignments", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, managementGroup.Id); !ok {
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
					managementGroupRoleAssignments = models.ManagementGroupRoleAssignments{
						ManagementGroupId: id,
					}
					count = 0
				)
				for item := range client.ListRoleAssignmentsForResource(ctx, id, "atScope()", "") {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing role assignments for this managementGroup", "managementGroupId", id)
					} else {
						managementGroupRoleAssignment := models.ManagementGroupRoleAssignment{
							ManagementGroupId: id,
							RoleAssignment:    item.Ok,
						}
						log.V(2).Info("found managementGroup role assignment", "managementGroupRoleAssignment", managementGroupRoleAssignment)
						count++
						managementGroupRoleAssignments.RoleAssignments = append(managementGroupRoleAssignments.RoleAssignments, managementGroupRoleAssignment)
					}
				}
				if ok := pipeline.Send(ctx.Done(), out, NewAzureWrapper(
					enums.KindAZManagementGroupRoleAssignment,
					managementGroupRoleAssignments,
				)); !ok {
					return
				}
				log.V(1).Info("finished listing managementGroup role assignments", "managementGroupId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all management group role assignments")
	}()

	return out
}

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
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listRoleAssignmentsCmd)
}

var listRoleAssignmentsCmd = &cobra.Command{
	Use:          "role-assignments",
	Long:         "Lists Azure Active Directory Role Assignments",
	Run:          listRoleAssignmentsCmdImpl,
	SilenceUsage: true,
}

func listRoleAssignmentsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory role assignments...")
	start := time.Now()
	roles := listRoles(ctx, azClient)
	stream := listRoleAssignments(ctx, azClient, roles)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listRoleAssignments(ctx context.Context, client client.AzureClient, roles <-chan interface{}) <-chan interface{} {
	var (
		out     = make(chan interface{})
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), roles) {
			if role, ok := result.(AzureWrapper).Data.(models.Role); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating role assignments", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, role.ID); !ok {
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
					roleAssignments = models.RoleAssignments{
						RoleDefinitionId: id,
						TenantId:         client.TenantInfo().TenantId,
					}
					count  = 0
					filter = fmt.Sprintf("roleDefinitionId eq '%s'", id)
				)
				// We expand directoryScope in order to obtain the appId from app specific scoped role assignments
				for item := range client.ListAzureADRoleAssignments(ctx, query.GraphParams{Filter: filter, Expand: "directoryScope"}) {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing role assignments for this role", "roleDefinitionId", id)
					} else {
						log.V(2).Info("found role assignment", "roleAssignments", item)
						count++
						// To ensure proper linking to AZApp nodes we want to supply the AppId instead when role assignments are app specific scoped
						if item.Ok.DirectoryScopeId != "/" {
							item.Ok.DirectoryScopeId = fmt.Sprintf("/%s", item.Ok.DirectoryScope.AppId)
						}
						roleAssignments.RoleAssignments = append(roleAssignments.RoleAssignments, item.Ok)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZRoleAssignment,
					Data: roleAssignments,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing role assignments", "roleDefinitionId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all role assignments")
	}()

	return out
}

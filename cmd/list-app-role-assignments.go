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
	listRootCmd.AddCommand(listAppRoleAssignmentsCmd)
}

var listAppRoleAssignmentsCmd = &cobra.Command{
	Use:          "app-role-assignments",
	Long:         "Lists Azure Active Directory App Role Assignments",
	Run:          listAppRoleAssignmentsCmdImpl,
	SilenceUsage: true,
}

func listAppRoleAssignmentsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory app role assignments...")
	start := time.Now()
	servicePrincipals := listServicePrincipals(ctx, azClient)
	stream := listAppRoleAssignments(ctx, azClient, servicePrincipals)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listAppRoleAssignments(ctx context.Context, client client.AzureClient, servicePrincipals <-chan interface{}) <-chan interface{} {
	var (
		out         = make(chan interface{})
		filteredSPs = make(chan models.ServicePrincipal)
		streams     = pipeline.Demux(ctx.Done(), filteredSPs, config.ColStreamCount.Value().(int))
		wg          sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(filteredSPs)

		for result := range pipeline.OrDone(ctx.Done(), servicePrincipals) {
			if servicePrincipal, ok := result.(AzureWrapper).Data.(models.ServicePrincipal); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating app role assignments", "result", result)
				return
			} else {
				if len(servicePrincipal.AppRoles) != 0 {
					if ok := pipeline.Send(ctx.Done(), filteredSPs, servicePrincipal); !ok {
						return
					}
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
			for servicePrincipal := range stream {
				var (
					count = 0
				)
				for item := range client.ListAzureADAppRoleAssignments(ctx, servicePrincipal.ID, query.GraphParams{}) {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing app role assignments for this service principal", "servicePrincipalId", servicePrincipal)
					} else {
						log.V(2).Info("found app role assignment", "roleAssignments", item)
						count++
						if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
							Kind: enums.KindAZAppRoleAssignment,
							Data: models.AppRoleAssignment{
								AppRoleAssignment: item.Ok,
								AppId:             servicePrincipal.AppId,
								TenantId:          client.TenantInfo().TenantId,
							},
						}); !ok {
							return
						}
					}
				}
				log.V(1).Info("finished listing app role assignments", "appId", servicePrincipal.AppId, "servicePrincipalId", servicePrincipal.ID, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all app role assignments")
	}()

	return out
}

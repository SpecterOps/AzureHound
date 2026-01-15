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
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listKeyVaultRoleAssignmentsCmd)
}

var listKeyVaultRoleAssignmentsCmd = &cobra.Command{
	Use:          "key-vault-role-assignments",
	Long:         "Lists Key Vault Role Assignments",
	Run:          listKeyVaultRoleAssignmentsCmdImpl,
	SilenceUsage: true,
}

func listKeyVaultRoleAssignmentsCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure key vault role assignments...")
	start := time.Now()
	subscriptions := listSubscriptions(ctx, azClient)
	stream := listKeyVaultRoleAssignments(ctx, azClient, listKeyVaults(ctx, azClient, subscriptions))
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listKeyVaultRoleAssignments(ctx context.Context, client client.AzureClient, keyVaults <-chan interface{}) <-chan azureWrapper[models.KeyVaultRoleAssignments] {
	var (
		out     = make(chan azureWrapper[models.KeyVaultRoleAssignments])
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), keyVaults) {
			if keyVault, ok := result.(AzureWrapper).Data.(models.KeyVault); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating key vault role assignments", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, keyVault.Id); !ok {
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
					keyVaultRoleAssignments = models.KeyVaultRoleAssignments{
						KeyVaultId: id,
					}
					count = 0
				)
				for item := range client.ListRoleAssignmentsForResource(ctx, id, "", "") {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing role assignments for this key vault", "keyVaultId", id)
					} else {
						keyVaultRoleAssignment := models.KeyVaultRoleAssignment{
							KeyVaultId:     id,
							RoleAssignment: item.Ok,
						}
						log.V(2).Info("found key vault role assignment", "name", keyVaultRoleAssignment.RoleAssignment.Name)
						count++
						keyVaultRoleAssignments.RoleAssignments = append(keyVaultRoleAssignments.RoleAssignments, keyVaultRoleAssignment)
					}
				}
				if ok := pipeline.Send(ctx.Done(), out, NewAzureWrapper(enums.KindAZKeyVaultRoleAssignment, keyVaultRoleAssignments)); !ok {
					return
				}
				log.V(1).Info("finished listing key vault role assignments", "keyVaultId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all key vault role assignments")
	}()

	return out
}

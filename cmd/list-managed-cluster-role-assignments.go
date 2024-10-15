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
	listRootCmd.AddCommand(listManagedClusterRoleAssignment)
}

var listManagedClusterRoleAssignment = &cobra.Command{
	Use:          "managed-cluster-role-assignments",
	Long:         "Lists AKS Managed Cluster Role Assignments",
	Run:          listManagedClusterRoleAssignmentImpl,
	SilenceUsage: true,
}

func listManagedClusterRoleAssignmentImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	if err := testConnections(); err != nil {
		exit(err)
	} else if azClient, err := newAzureClient(); err != nil {
		exit(err)
	} else {
		log.Info("collecting azure managed cluster role assignments...")
		start := time.Now()
		subscriptions := listSubscriptions(ctx, azClient)
		stream := listManagedClusterRoleAssignments(ctx, azClient, listManagedClusters(ctx, azClient, subscriptions))
		panicrecovery.HandleBubbledPanic(ctx, stop, log)
		outputStream(ctx, stream)
		duration := time.Since(start)
		log.Info("collection completed", "duration", duration.String())
	}
}

func listManagedClusterRoleAssignments(ctx context.Context, client client.AzureClient, managedClusters <-chan interface{}) <-chan interface{} {
	var (
		out     = make(chan interface{})
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), managedClusters) {
			if managedCluster, ok := result.(AzureWrapper).Data.(models.ManagedCluster); !ok {
				log.Error(fmt.Errorf("failed type assertion"), "unable to continue enumerating managed cluster role assignments", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, managedCluster.Id); !ok {
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
					managedClusterRoleAssignments = models.AzureRoleAssignments{
						ObjectId: id,
					}
					count = 0
				)
				for item := range client.ListRoleAssignmentsForResource(ctx, id, "", "") {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing role assignments for this managed cluster", "managedClusterId", id)
					} else {
						roleDefinitionId := path.Base(item.Ok.Properties.RoleDefinitionId)

						managedClusterRoleAssignment := models.AzureRoleAssignment{
							Assignee:         item.Ok,
							ObjectId:         id,
							RoleDefinitionId: roleDefinitionId,
						}
						log.V(2).Info("found managed cluster role assignment", "managedClusterRoleAssignment", managedClusterRoleAssignment)
						count++
						managedClusterRoleAssignments.RoleAssignments = append(managedClusterRoleAssignments.RoleAssignments, managedClusterRoleAssignment)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZManagedClusterRoleAssignment,
					Data: managedClusterRoleAssignments,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing managed cluster role assignments", "managedClusterId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing all managed cluster role assignments")
	}()

	return out
}

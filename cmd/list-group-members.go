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
	listRootCmd.AddCommand(listGroupMembersCmd)
	listGroupMembersCmd.Flags().StringSliceVar(&listGroupMembersSelect, "select", []string{"id,displayName,createdDateTime"}, `Select properties to include. Use "" for Azure default properties. Azurehound default is "id,displayName,createdDateTime" if flag is not supplied.`)
}

var listGroupMembersCmd = &cobra.Command{
	Use:          "group-members",
	Long:         "Lists Azure AD Group Members",
	Run:          listGroupMembersCmdImpl,
	SilenceUsage: true,
}

var listGroupMembersSelect []string

func listGroupMembersCmdImpl(cmd *cobra.Command, _ []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure group members...")
	start := time.Now()
	stream := listGroupMembers(ctx, azClient, listGroups(ctx, azClient))
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listGroupMembers(ctx context.Context, client client.AzureClient, groups <-chan interface{}) <-chan interface{} {
	var (
		out     = make(chan interface{})
		ids     = make(chan string)
		streams = pipeline.Demux(ctx.Done(), ids, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
		params  = query.GraphParams{
			Select: unique(listGroupMembersSelect),
			Filter: "",
			Count:  false,
			Search: "",
			Top:    0,
			Expand: "",
		}
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(ids)

		for result := range pipeline.OrDone(ctx.Done(), groups) {
			if group, ok := result.(AzureWrapper).Data.(models.Group); !ok {
				log.Error(fmt.Errorf("failed group type assertion"), "unable to continue enumerating group members", "result", result)
				return
			} else {
				if ok := pipeline.Send(ctx.Done(), ids, group.ID); !ok {
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
					data = models.GroupMembers{
						GroupId: id,
					}
					count = 0
				)
				for item := range client.ListAzureADGroupMembers(ctx, id, params) {
					if item.Error != nil {
						log.Error(item.Error, "unable to continue processing members for this group", "groupId", id)
					} else {
						groupMember := models.GroupMember{
							Member:  item.Ok,
							GroupId: id,
						}
						log.V(2).Info("found group member", "groupMember", groupMember)
						count++
						data.Members = append(data.Members, groupMember)
					}
				}
				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
					Kind: enums.KindAZGroupMember,
					Data: data,
				}); !ok {
					return
				}
				log.V(1).Info("finished listing group memberships", "groupId", id, "count", count)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
		log.Info("finished listing members for all groups")
	}()

	return out
}

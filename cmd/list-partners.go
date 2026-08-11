// Copyright (C) 2026 Specter Ops, Inc.
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
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

var externalLinkSuffix = regexp.MustCompile(`\s@\([^)]*\)`)

func init() {
	listRootCmd.AddCommand(listPartnersCmd)
}

var listPartnersCmd = &cobra.Command{
	Use:          "partners",
	Long:         "Lists Azure Active Directory Delegated Partners",
	Run:          listPartnersCmdImpl,
	SilenceUsage: true,
}

func listPartnersCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure active directory delegated partners...")
	start := time.Now()
	stream := listPartners(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listPartners(ctx context.Context, client client.AzureClient) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)
		count := 0
		partnerTenants := make(map[string]azure.Tenant, 10)

		for partner := range client.ListAzureADPartners(ctx, query.GraphParams{}) {
			if partner.Error != nil {
				log.Error(partner.Error, "unable to continue processing partners")
				return
			}

			log.V(2).Info("found partner", "companyName", partner.Ok.CompanyName, "partnerTenantId", partner.Ok.PartnerTenantId)
			count++

			// Begin by fetching the partner tenant information
			externalTenant, err := client.GetAzureADTenantInfoById(ctx, partner.Ok.PartnerTenantId)
			if err != nil {
				log.Error(err, "failed to retrieve tenant information for external partner", "companyName", partner.Ok.CompanyName, "partnerTenantId", partner.Ok.PartnerTenantId)
			}

			externalTenant.Id = fmt.Sprintf("/tenants/%s", externalTenant.TenantId)
			externalTenant.TenantType = partner.Ok.CompanyType

			if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{
				Kind: enums.KindAZTenant,
				Data: models.Tenant{
					Tenant:   externalTenant,
					External: true,
				},
			}); !ok {
				return
			}

			partnerTenants[externalTenant.TenantId] = externalTenant
		}
		log.Info("finished listing all delegated partners", "count", count)

		count = 0

		// This part is a bit hacky but i'll try to explain what's going on:
		//
		// For partners, associated pricipal data is stored in their tenant.
		// This means that you unfortunately can't just directly query our own
		// list of groups/users/service principals and get back the principal
		// information directly. For some reason Microsoft decided to lock this
		// info behind calls that let you query information via `$expand` queries.
		//
		// While I'd love to filter based on `principalOrganizationId`, this field
		// seems to be some dynamic magic field on the backend and therefor can't be
		// filtered on. Morover, if you just use `$expand` on `principal` and try to
		// list all role assignments, you still won't get the information you're looking
		// for. So far the only way I'm able to reliably filter external tenant's
		// information is by passing a `roleDefinitionId` filter on `roleAssignments`
		// which results in the `principal` field and the `principalOrganizationId` field
		// being present.
		//
		// If you find a more efficient way of getting this info I'd love to see an
		// improved version :)
		observedPrincipalIds := make(map[string]bool)

		for role := range client.ListAzureADRoles(ctx, query.GraphParams{}) {
			if role.Error != nil {
				log.Error(role.Error, "unable to continue processing partner roles")
				break
			}

			for item := range client.ListAzureADRoleAssignments(ctx, query.GraphParams{
				Filter: fmt.Sprintf("roleDefinitionId eq '%s'", role.Ok.Id),
				Expand: "principal",
			}) {
				if item.Error != nil {
					log.Error(item.Error, "unable to continue processing partner role assignments")
					break
				}

				tenant, exists := partnerTenants[item.Ok.PrincipalOrganizationId]
				if !exists {
					continue
				}

				var header struct {
					Type        string `json:"@odata.type"`
					Id          string `json:"Id"`
					DisplayName string `json:"DisplayName,omitempty"`
				}

				if err := json.Unmarshal(item.Ok.Principal, &header); err != nil {
					log.Error(err, "unable to determine principal type")
					continue
				}

				if _, ok := observedPrincipalIds[header.Id]; ok {
					continue
				}

				var (
					kind enums.Kind
					data any
				)

				switch header.Type {
				case "#microsoft.graph.user":
					var user azure.User
					if err := json.Unmarshal(item.Ok.Principal, &user); err != nil {
						log.Error(err, "unable to unmarshal user principal")
						continue
					}
					user.DisplayName = externalLinkSuffix.ReplaceAllString(user.DisplayName, "")
					kind = enums.KindAZUser
					data = models.User{User: user, TenantId: tenant.TenantId, TenantName: tenant.DisplayName}
					log.V(2).Info("found partner user information", "id", item.Ok.Id)

				case "#microsoft.graph.group":
					var group azure.Group
					if err := json.Unmarshal(item.Ok.Principal, &group); err != nil {
						log.Error(err, "unable to unmarshal group principal")
						continue
					}
					group.DisplayName = externalLinkSuffix.ReplaceAllString(group.DisplayName, "")
					kind = enums.KindAZGroup
					data = models.Group{Group: group, TenantId: tenant.TenantId, TenantName: tenant.DisplayName}
					log.V(2).Info("found partner group information", "id", item.Ok.Id)

				case "#microsoft.graph.servicePrincipal":
					var sp azure.ServicePrincipal
					if err := json.Unmarshal(item.Ok.Principal, &sp); err != nil {
						log.Error(err, "unable to unmarshal service principal")
						continue
					}
					sp.DisplayName = externalLinkSuffix.ReplaceAllString(sp.DisplayName, "")
					kind = enums.KindAZServicePrincipal
					data = models.ServicePrincipal{ServicePrincipal: sp, TenantId: tenant.TenantId, TenantName: tenant.DisplayName}
					log.V(2).Info("found partner service principal information", "id", item.Ok.Id)

				default:
					log.V(2).Info("skipping unknown principal type", "type", header.Type)
					continue
				}

				observedPrincipalIds[header.Id] = true

				if ok := pipeline.SendAny(ctx.Done(), out, AzureWrapper{Kind: kind, Data: data}); !ok {
					break
				}

				count++
			}
		}

		log.Info("finished listing all delegated partner principals", "count", count)
	}()

	return out
}

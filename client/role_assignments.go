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

package client

import (
	"context"
	"fmt"

	"github.com/certmichelin/azurehound/v2/client/query"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/models/azure"
)

// ListAzureADRoleAssignments https://learn.microsoft.com/en-us/graph/api/rbacapplication-list-roleassignments?view=graph-rest-beta
func (s *azureClient) ListAzureADRoleAssignments(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleAssignment] {
	var (
		out  = make(chan AzureResult[azure.UnifiedRoleAssignment])
		path = fmt.Sprintf("/%s/roleManagement/directory/roleAssignments", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 999
	}

	go getAzureObjectList[azure.UnifiedRoleAssignment](s.msgraph, ctx, path, params, out)
	return out
}

// ListRoleAssignmentsForResource https://learn.microsoft.com/en-us/rest/api/authorization/role-assignments/list-for-resource?view=rest-authorization-2015-07-01
func (s *azureClient) ListRoleAssignmentsForResource(ctx context.Context, resourceId string, filter, tenantId string) <-chan AzureResult[azure.RoleAssignment] {
	var (
		out    = make(chan AzureResult[azure.RoleAssignment])
		path   = fmt.Sprintf("%s/providers/Microsoft.Authorization/roleAssignments", resourceId)
		params = query.RMParams{ApiVersion: "2015-07-01", Filter: filter, TenantId: tenantId}
	)

	go getAzureObjectList[azure.RoleAssignment](s.resourceManager, ctx, path, params, out)

	return out
}

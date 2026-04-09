// Copyright (C) 2025 Specter Ops, Inc.
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

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// AzureRoleManagementClient defines the methods to interface with the Azure role based access control (RBAC) API
type AzureRoleManagementClient interface {
	ListAzureUnifiedRoleEligibilityScheduleInstances(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance]
	ListRoleAssignmentPolicies(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleManagementPolicyAssignment]
}

// ListAzureUnifiedRoleEligibilityScheduleInstances https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityscheduleinstance?view=graph-rest-1.0
func (s *azureClient) ListAzureUnifiedRoleEligibilityScheduleInstances(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance] {
	var (
		out  = make(chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance])
		path = fmt.Sprintf("/%s/roleManagement/directory/roleEligibilityScheduleInstances", constants.GraphApiVersion)
	)

	go getAzureObjectList[azure.UnifiedRoleEligibilityScheduleInstance](s.msgraph, ctx, path, params, out)

	return out
}

// ListRoleAssignmentPolicies makes a GET request to  https://graph.microsoft.com/v1.0/policies/roleManagementPolicyAssignments
// This endpoint requires the RoleManagement.Read.All permission
// https://learn.microsoft.com/en-us/graph/permissions-reference#rolemanagementreadall
// Endpoint documentation: https://learn.microsoft.com/en-us/graph/api/policyroot-list-rolemanagementpolicyassignments?view=graph-rest-1.0&tabs=http
func (s *azureClient) ListRoleAssignmentPolicies(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleManagementPolicyAssignment] {
	var (
		out  = make(chan AzureResult[azure.UnifiedRoleManagementPolicyAssignment])
		path = fmt.Sprintf("/%s/policies/roleManagementPolicyAssignments", constants.GraphApiVersion)
	)

	go getAzureObjectList[azure.UnifiedRoleManagementPolicyAssignment](s.msgraph, ctx, path, params, out)

	return out
}

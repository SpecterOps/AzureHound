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
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// ListAzureRoleEligibilityScheduleInstances https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleeligibilityscheduleinstance?view=graph-rest-1.0
func (s *azureClient) ListAzureRoleEligibilityScheduleInstances(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance] {
	var (
		out  = make(chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance])
		path = fmt.Sprintf("/subscriptions/%s/resourcegroups", subscriptionId)
	)

	if params.ApiVersion == "" {
		params.ApiVersion = "2021-04-01"
	}

	go getAzureObjectList[azure.UnifiedRoleEligibilityScheduleInstance](s.resourceManager, ctx, path, params, out)

	return out
}

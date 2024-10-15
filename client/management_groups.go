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
	"github.com/certmichelin/azurehound/v2/models/azure"
)

// ListAzureManagementGroups https://learn.microsoft.com/en-us/rest/api/managementgroups/management-groups/list?view=rest-managementgroups-2020-05-01
func (s *azureClient) ListAzureManagementGroups(ctx context.Context, skipToken string) <-chan AzureResult[azure.ManagementGroup] {
	var (
		out    = make(chan AzureResult[azure.ManagementGroup])
		path   = "/providers/Microsoft.Management/managementGroups"
		params = query.RMParams{ApiVersion: "2020-05-01", SkipToken: skipToken}
	)

	go getAzureObjectList[azure.ManagementGroup](s.resourceManager, ctx, path, params, out)

	return out
}

// ListAzureManagementGroupDescendants https://learn.microsoft.com/en-us/rest/api/managementgroups/management-groups/get-descendants?view=rest-managementgroups-2020-05-01
func (s *azureClient) ListAzureManagementGroupDescendants(ctx context.Context, groupId string, top int32) <-chan AzureResult[azure.DescendantInfo] {
	var (
		out    = make(chan AzureResult[azure.DescendantInfo])
		path   = fmt.Sprintf("/providers/Microsoft.Management/managementGroups/%s/descendants", groupId)
		params = query.RMParams{ApiVersion: "2020-05-01", Top: top}
	)

	go getAzureObjectList[azure.DescendantInfo](s.resourceManager, ctx, path, params, out)

	return out
}

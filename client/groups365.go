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
	"encoding/json"
	"fmt"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// ListAzureGroups Microsoft 365 https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-beta
func (s *azureClient) ListAzureADGroups365(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Group365] {
	var (
		out  = make(chan AzureResult[azure.Group365])
		path = fmt.Sprintf("/%s/groups", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 99
	}

	go getAzureObjectList[azure.Group365](s.msgraph, ctx, path, params, out)

	return out
}

// ListAzureADGroupOwners Microsoft 365 https://learn.microsoft.com/en-us/graph/api/group-list-owners?view=graph-rest-beta
func (s *azureClient) ListAzureADGroup365Owners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
	var (
		out  = make(chan AzureResult[json.RawMessage])
		path = fmt.Sprintf("/%s/groups/%s/owners", constants.GraphApiBetaVersion, objectId)
	)

	if params.Top == 0 {
		params.Top = 99
	}

	go getAzureObjectList[json.RawMessage](s.msgraph, ctx, path, params, out)

	return out
}

// ListAzureADGroupMembers Microsoft 365 https://learn.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-beta
func (s *azureClient) ListAzureADGroup365Members(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
	var (
		out  = make(chan AzureResult[json.RawMessage])
		path = fmt.Sprintf("/%s/groups/%s/members", constants.GraphApiBetaVersion, objectId)
	)

	go getAzureObjectList[json.RawMessage](s.msgraph, ctx, path, params, out)

	return out
}

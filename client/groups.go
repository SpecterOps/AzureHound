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

	"github.com/certmichelin/azurehound/v2/client/query"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/models/azure"
)

// ListAzureADGroups https://learn.microsoft.com/en-us/graph/api/group-list?view=graph-rest-beta
func (s *azureClient) ListAzureADGroups(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Group] {
	var (
		out  = make(chan AzureResult[azure.Group])
		path = fmt.Sprintf("/%s/groups", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 99
	}

	go getAzureObjectList[azure.Group](s.msgraph, ctx, path, params, out)

	return out
}

// ListAzureADGroupOwners https://learn.microsoft.com/en-us/graph/api/group-list-owners?view=graph-rest-beta
func (s *azureClient) ListAzureADGroupOwners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
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

// ListAzureADGroupMembers https://learn.microsoft.com/en-us/graph/api/group-list-members?view=graph-rest-beta
func (s *azureClient) ListAzureADGroupMembers(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
	var (
		out  = make(chan AzureResult[json.RawMessage])
		path = fmt.Sprintf("/%s/groups/%s/members", constants.GraphApiBetaVersion, objectId)
	)

	go getAzureObjectList[json.RawMessage](s.msgraph, ctx, path, params, out)

	return out
}

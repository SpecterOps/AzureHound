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

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// List Azure OAuth2 Permission Grant https://learn.microsoft.com/en-us/graph/api/resources/oauth2permissiongrant?view=graph-rest-1.0
func (s *azureClient) ListAzureOauth2PermissionGrants(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.OAuth2PermissionGrant] {
	var (
		out  = make(chan AzureResult[azure.OAuth2PermissionGrant])
		path = fmt.Sprintf("/%s/oauth2PermissionGrants", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 99
	}

	go getAzureObjectList[azure.OAuth2PermissionGrant](s.msgraph, ctx, path, params, out)

	return out
}

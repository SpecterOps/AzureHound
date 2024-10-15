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

// GetAzureADAppRoleAssignments https://learn.microsoft.com/en-us/graph/api/serviceprincipal-list-approleassignedto?view=graph-rest-1.0
func (s *azureClient) ListAzureADAppRoleAssignments(ctx context.Context, servicePrincipalId string, params query.GraphParams) <-chan AzureResult[azure.AppRoleAssignment] {
	var (
		out  = make(chan AzureResult[azure.AppRoleAssignment])
		path = fmt.Sprintf("/%s/servicePrincipals/%s/appRoleAssignedTo", constants.GraphApiVersion, servicePrincipalId)
	)

	if params.Top == 0 {
		params.Top = 999
	}

	go getAzureObjectList[azure.AppRoleAssignment](s.msgraph, ctx, path, params, out)

	return out
}

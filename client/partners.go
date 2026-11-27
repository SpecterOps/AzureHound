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

package client

import (
	"context"
	"fmt"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// ListAzureADPartners
// Attempts to list partners using the (undocumented) `/directory/partners` API that can be
// seen being called when visiting partner relationships tab in Entra ID <https://portal.azure.com/#view/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/~/PartnerRelationships>
func (s *azureClient) ListAzureADPartners(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Partner] {
	var (
		out  = make(chan AzureResult[azure.Partner])
		path = fmt.Sprintf("/%s/directory/partners", constants.GraphApiVersion)
	)

	go getAzureObjectList[azure.Partner](s.msgraph, ctx, path, params, out)

	return out
}

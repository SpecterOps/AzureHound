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

// ListAzureKeyVaults https://learn.microsoft.com/en-us/rest/api/keyvault/keyvault/vaults/list-by-subscription?view=rest-keyvault-keyvault-2019-09-01
func (s *azureClient) ListAzureKeyVaults(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.KeyVault] {
	var (
		out  = make(chan AzureResult[azure.KeyVault])
		path = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.KeyVault/vaults", subscriptionId)
	)

	if params.ApiVersion == "" {
		params.ApiVersion = "2019-09-01"
	}

	go getAzureObjectList[azure.KeyVault](s.resourceManager, ctx, path, params, out)

	return out
}

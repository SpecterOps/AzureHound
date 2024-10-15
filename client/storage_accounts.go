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

// ListAzureStorageAccounts https://learn.microsoft.com/en-us/rest/api/storagerp/storage-accounts/list?view=rest-storagerp-2022-05-01
func (s *azureClient) ListAzureStorageAccounts(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.StorageAccount] {
	var (
		out    = make(chan AzureResult[azure.StorageAccount])
		path   = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Storage/storageAccounts", subscriptionId)
		params = query.RMParams{ApiVersion: "2022-05-01"}
	)

	go getAzureObjectList[azure.StorageAccount](s.resourceManager, ctx, path, params, out)

	return out
}

// ==
// Storage containers
// ==

// ListAzureStorageContainers https://learn.microsoft.com/en-us/rest/api/storagerp/blob-containers/list?view=rest-storagerp-2022-05-01
func (s *azureClient) ListAzureStorageContainers(ctx context.Context, subscriptionId string, resourceGroupName string, saName string, filter string, includeDeleted string, maxPageSize string) <-chan AzureResult[azure.StorageContainer] {
	var (
		out    = make(chan AzureResult[azure.StorageContainer])
		path   = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Storage/storageAccounts/%s/blobServices/default/containers", subscriptionId, resourceGroupName, saName)
		params = query.RMParams{ApiVersion: "2022-05-01", Filter: filter, IncludeDeleted: includeDeleted, MaxPageSize: maxPageSize}
	)

	go getAzureObjectList[azure.StorageContainer](s.resourceManager, ctx, path, params, out)

	return out
}

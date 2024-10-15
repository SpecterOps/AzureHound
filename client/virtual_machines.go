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

// ListAzureVirtualMachines https://learn.microsoft.com/en-us/rest/api/compute/virtual-machines/list-all?view=rest-compute-2021-07-01
func (s *azureClient) ListAzureVirtualMachines(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.VirtualMachine] {
	var (
		out  = make(chan AzureResult[azure.VirtualMachine])
		path = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Compute/virtualMachines", subscriptionId)
	)

	if params.ApiVersion == "" {
		params.ApiVersion = "2021-07-01"
	}

	go getAzureObjectList[azure.VirtualMachine](s.resourceManager, ctx, path, params, out)

	return out
}

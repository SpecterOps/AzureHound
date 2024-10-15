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

// ListAzureLogicApps https://learn.microsoft.com/en-us/rest/api/logic/workflows/list-by-subscription?view=rest-logic-2016-06-01
func (s *azureClient) ListAzureLogicApps(ctx context.Context, subscriptionId string, filter string, top int32) <-chan AzureResult[azure.LogicApp] {
	var (
		out    = make(chan AzureResult[azure.LogicApp])
		path   = fmt.Sprintf("/subscriptions/%s/providers/Microsoft.Logic/workflows", subscriptionId)
		params = query.RMParams{ApiVersion: "2016-06-01", Filter: filter, Top: top}
	)

	go getAzureObjectList[azure.LogicApp](s.resourceManager, ctx, path, params, out)

	return out
}

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
	"github.com/certmichelin/azurehound/v2/client/rest"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/models/azure"
)

func (s *azureClient) GetAzureADOrganization(ctx context.Context, selectCols []string) (*azure.Organization, error) {
	var (
		path     = fmt.Sprintf("/%s/organization", constants.GraphApiVersion)
		response azure.OrganizationList
	)
	if res, err := s.msgraph.Get(ctx, path, query.GraphParams{Select: selectCols}, nil); err != nil {
		return nil, err
	} else if err := rest.Decode(res.Body, &response); err != nil {
		return nil, err
	} else {
		return &response.Value[0], nil
	}
}

func (s *azureClient) GetAzureADTenants(ctx context.Context, includeAllTenantCategories bool) (azure.TenantList, error) {
	var (
		path     = "/tenants"
		params   = query.RMParams{ApiVersion: "2020-01-01", IncludeAllTenantCategories: includeAllTenantCategories}
		headers  map[string]string
		response azure.TenantList
	)

	if res, err := s.resourceManager.Get(ctx, path, params, headers); err != nil {
		return response, err
	} else if err := rest.Decode(res.Body, &response); err != nil {
		return response, err
	} else {
		return response, nil
	}
}

// ListAzureADTenants https://learn.microsoft.com/en-us/rest/api/subscription/tenants/list?view=rest-subscription-2020-01-01
func (s *azureClient) ListAzureADTenants(ctx context.Context, includeAllTenantCategories bool) <-chan AzureResult[azure.Tenant] {
	var (
		out    = make(chan AzureResult[azure.Tenant])
		path   = "/tenants"
		params = query.RMParams{ApiVersion: "2020-01-01", IncludeAllTenantCategories: includeAllTenantCategories}
	)

	go getAzureObjectList[azure.Tenant](s.resourceManager, ctx, path, params, out)

	return out
}

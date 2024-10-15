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

package azure

import "github.com/certmichelin/azurehound/v2/enums"

type Tenant struct {
	Country               string               `json:"country,omitempty"`               // Country/region name of the address for the tenant.
	CountryCode           string               `json:"countryCode,omitempty"`           // Country/region abbreviation for the tenant.
	DefaultDomain         string               `json:"defaultDomain,omitempty"`         // The default domain for the tenant.
	DisplayName           string               `json:"displayName,omitempty"`           // The display name of the tenant.
	Domains               []string             `json:"domains,omitempty"`               // The list of domains for the tenant
	Id                    string               `json:"id,omitempty"`                    // The fully qualified ID of the tenant. E.g. "/tenants/00000000-0000-0000-0000-000000000000"
	TenantBrandingLogoUrl string               `json:"tenantBrandingLogoUrl,omitempty"` // The tenant's branding logo URL. Only available for 'Home' TenantCategory
	TenantCategory        enums.TenantCategory `json:"tenantCategory,omitempty"`        // The category of the tenant.
	TenantId              string               `json:"tenantId,omitempty"`              // Then tenant ID. E.g. "00000000-0000-0000-0000-000000000000"
	TenantType            string               `json:"tenantType,omitempty"`            // The tenant type. Only available for 'Home' TenantCategory
}

type TenantList struct {
	NextLink string   `json:"nextLink,omitempty"` // The URL to use for getting the next set of values.
	Value    []Tenant `json:"value"`              // A list of tenants.
}

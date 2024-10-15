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

import (
	"github.com/certmichelin/azurehound/v2/enums"
)

type Subscription struct {
	Entity

	// The authorization source of the request. Valid values are one or more combinations of Legacy, RoleBased,
	// Bypassed, Direct and Management. For example, 'Legacy, RoleBased'.
	AuthorizationSource string `json:"authorizationSource,omitempty"`

	// The subscription display name.
	DisplayName string `json:"displayName,omitempty"`

	// A list of tenants managing the subscription.
	ManagedByTenants []ManagedByTenant `json:"managedByTenants,omitempty"`

	// The subscription state.
	State enums.SubscriptionState `json:"state,omitempty"`

	// The subscription ID.
	SubscriptionId string `json:"subscriptionId,omitempty"`

	// The subscription policies.
	SubscriptionPolicies SubscriptionPolicies `json:"subscriptionPolicies,omitempty"`

	// The tags attached to the subscription.
	Tags map[string]string `json:"tags,omitempty"`

	// The subscription tenant ID.
	TenantId string `json:"tenantId,omitempty"`
}

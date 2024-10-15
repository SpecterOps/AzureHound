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

// Deprecated, use ManagedIdentity
type VirtualMachineIdentity ManagedIdentity

// Managed identity.
type ManagedIdentity struct {
	// The principal id of the managed identity. The property will only be provided for a system assigned
	// identity.
	PrincipalId string `json:"principalId,omitempty"`

	// The tenant id associated with the managed identity. This property will only be provided for a system assigned
	// identity.
	TenantId string `json:"tenantId,omitempty"`

	// The type of identity used.
	Type enums.Identity `json:"type,omitempty"`

	// The list of user identities associated with the Managed identity. The user identity dictionary key references will be
	// ARM resource ids in the form:
	// '/subscriptions/{subscriptionId}/resourceGroups/{groupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'
	UserAssignedIdentities map[string]UserAssignedIdentity `json:"userAssignedIdentities,omitempty"`
}

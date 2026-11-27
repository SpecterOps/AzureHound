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

import "encoding/json"

type UnifiedRoleAssignment struct {
	Entity

	// Identifier of the role definition the assignment is for.
	// Read only.
	// Supports $filter (eq, in).
	RoleDefinitionId string `json:"roleDefinitionId,omitempty"`

	// Identifier of the principal to which the assignment is granted.
	// Supports $filter (eq, in).
	PrincipalId string `json:"principalId,omitempty"`

	// Identifier of the organization (tenant) the principal belongs to
	// Read only and does not support $filter.
	PrincipalOrganizationId string `json:"principalOrganizationId,omitempty"`

	// Identifier of the directory object representing the scope of the assignment.
	// Either this property or appScopeId is required.
	// The scope of an assignment determines the set of resources for which the principal has been granted access.
	// Directory scopes are shared scopes stored in the directory that are understood by multiple applications.
	//
	// Use / for tenant-wide scope.
	// Use appScopeId to limit the scope to an application only.
	//
	// Supports $filter (eq, in).
	DirectoryScopeId string `json:"directoryScopeId,omitempty"`

	// Identifier of the resource representing the scope of the assignment.
	ResourceScope string `json:"resourceScope,omitempty"`

	// Identifier of the app-specific scope when the assignment scope is app-specific.
	// Either this property or directoryScopeId is required.
	// App scopes are scopes that are defined and understood by this application only.
	//
	// Use / for tenant-wide app scopes.
	// Use directoryScopeId to limit the scope to particular directory objects, for example, administrative units.
	//
	// Supports $filter (eq, in).
	AppScopeId string `json:"appScopeId,omitempty"`

	// Referencing the assigned principal.
	// Read-only.
	// Supports $expand.
	Principal json.RawMessage `json:"principal,omitempty"`

	// The roleDefinition the assignment is for.
	// Supports $expand. roleDefinition.Id will be auto expanded.
	RoleDefinition UnifiedRoleDefinition `json:"roleDefinition,omitempty"`

	// The directory object that is the scope of the assignment.
	// Read-only.
	// Supports $expand.
	DirectoryScope Application `json:"directoryScope,omitempty"`

	// Read-only property with details of the app specific scope when the assignment scope is app specific.
	// Containment entity.
	// Supports $expand.
	AppScope AppScope `json:"appScope,omitempty"`
}

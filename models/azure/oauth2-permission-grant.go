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

// OAuth2PermissionGrant represents an OAuth2 Permission Grant in Azure Active Directory.
// It contains information about the client, consent type, principal, resource, and scope.
// For more detail see https://learn.microsoft.com/en-us/graph/api/resources/oauth2permissiongrant?view=graph-rest-1.0
type OAuth2PermissionGrant struct {
	DirectoryObject

	// Client ID of the application that requested this permission grant.
	ClientId string `json:"clientId,omitempty"`

	// Type of the Consent. Possible values are: "AllPrincipals", "Principal", "Application".
	ConsentType string `json:"consentType,omitempty"`

	// Id of the OAuth2PermissionGrant.
	Id string `json:"id,omitempty"`

	// PrincipalId of the user or service principal that the permission grant is for. (null if consentType is "AllPrincipals")
	PrincipalId string `json:"principalId,omitempty"`

	// ResourceId of the resource that the permission grant is for.
	ResourceId string `json:"resourceId,omitempty"`

	// Scope of the permission grant.
	Scope string `json:"scope,omitempty"`
}

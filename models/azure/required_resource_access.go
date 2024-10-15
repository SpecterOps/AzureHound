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
	"github.com/gofrs/uuid"
)

// Object used to specify an OAuth 2.0 permission scope or an app role that an application requires.
type ResourceAccess struct {
	// The unique identifier for one of the OAuth2PermissionScopes or AppRole instances that the resource application
	// exposes.
	Id uuid.UUID `json:"id,omitempty"`

	// Specifies whether the {@link Id} property references an OAuth2PermissionScope or AppRole.
	Type enums.AccessType `json:"type,omitempty"`
}

// Specifies the set of OAuth 2.0 permission scopes and app roles under the specified resource that an application
// requires access to. The application may request the specified OAuth 2.0 permission scopes or app roles through the
// requiredResourceAccess property.
type RequiredResourceAccess struct {
	// The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource.
	ResourceAccess []ResourceAccess `json:"resourceAccess,omitempty"`

	// The unique identifier for the resource that the application requires access to. This should be equal to the
	// {@link AppId} declared on the target resource application.
	ResourceAppId string `json:"resourceAppId,omitempty"`
}

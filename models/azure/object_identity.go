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

// Represents an identity used to sign in to a user account. An identity can be provided by Microsoft (a.k.a. local
// account), by organizations, or by 3rd party identity providers such as Facebook or Google that are tied to a user
// account.
// Note: For `$filter` both {@link Issuer} and {@link IssuerAssignedId} must be supplied.
// For more detail, see https://docs.microsoft.com/en-us/graph/api/resources/objectidentity?view=graph-rest-1.0
type ObjectIdentity struct {
	// Specifies the user sign-in types in your directory.
	// Federated represents a unique identifier for a user from an issuer, that can be in any format chosen by the
	// issuer.
	// Setting or updating a UserPrincipalName identity will update the value of the userPrincipalName property on the
	// user object. The validations performed on the userPrincipalName property on the user object, for example,
	// verified domains and acceptable characters, will be performed when setting or updating a UserPrincipalName
	// identity.
	// Additional validation is enforced on issuerAssignedId when the sign-in type is set to Email or UserName.
	// This property can also be set to any custom string; use string(SignInType) or enums.signintype(someValue) to
	// convert appropriately.
	SignInType enums.SigninType `json:"signInType,omitempty"`

	// Specifies the issuer of the identity.
	// **Notes:**
	// * For local accounts where {@link SignInType} is not `federated`, the value is the local B2C tenant default domain
	//   name.
	// * For external users from other Azure AD organizations, this will be the domain of the federated organization.
	//
	// Supports `$filter` w/ 512 character limit.
	Issuer string `json:"issuer,omitempty"`

	// Specifies the unique identifier assigned to the user by the issuer. The combination of issuer and
	// issuerAssignedId must be unique within the organization.
	// For more detail, see https://docs.microsoft.com/en-us/graph/api/resources/objectidentity?view=graph-rest-1.0
	//
	// Supports `$filter` w/ 100 character limit
	IssuerAssignedId string `json:"issuerAssignedId,omitempty"`
}

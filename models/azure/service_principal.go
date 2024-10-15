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

// Represents an instance of an application in a directory.
// For more detail see https://docs.microsoft.com/en-us/graph/api/resources/serviceprincipal?view=graph-rest-1.0
type ServicePrincipal struct {
	DirectoryObject

	// true if the service principal account is enabled; otherwise, false.
	// Supports $filter (eq, ne, NOT, in).
	AccountEnabled bool `json:"accountEnabled,omitempty"`

	// Defines custom behavior that a consuming service can use to call an app in specific contexts.
	// For example, applications that can render file streams may set the addIns property for its "FileHandler"
	// functionality. This will let services like Microsoft 365 call the application in the context of a document the
	// user is working on.
	AddIns []AddIn `json:"addIns,omitempty"`

	// Used to retrieve service principals by subscription, identify resource group and full resource ids for managed
	// identities.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	AlternativeNames []string `json:"alternativeNames,omitempty"`

	// The description exposed by the associated application.
	AppDescription string `json:"appDescription,omitempty"`

	// The display name exposed by the associated application.
	AppDisplayName string `json:"appDisplayName,omitempty"`

	// The unique identifier for the associated application (its appId property).
	AppId string `json:"appId,omitempty"`

	// Unique identifier of the applicationTemplate that the servicePrincipal was created from.
	// Read-only.
	// Supports $filter (eq, ne, NOT, startsWith).
	ApplicationTemplateId string `json:"applicationTemplateId,omitempty"`

	// Contains the tenant id where the application is registered.
	// This is applicable only to service principals backed by applications.
	// Supports $filter (eq, ne, NOT, ge, le).
	AppOwnerOrganizationId string `json:"appOwnerOrganizationId,omitempty"`

	// Specifies whether users or other service principals need to be granted an app role assignment for this service
	// principal before users can sign in or apps can get tokens.
	// The default value is false.
	// Not nullable.
	// Supports $filter (eq, ne, NOT).
	AppRoleAssignmentRequired bool `json:"appRoleAssignmentRequired,omitempty"`

	// The roles exposed by the application which this service principal represents.
	// For more information see the appRoles property definition on the application entity.
	// Not nullable.
	AppRoles []AppRole `json:"appRoles,omitempty"`

	// The date and time the service principal was deleted.
	// Read-only.
	DeletedDateTime string `json:"deletedDateTime,omitempty"`

	// Free text field to provide an internal end-user facing description of the service principal.
	// End-user portals such MyApps will display the application description in this field.
	// The maximum allowed size is 1024 characters.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith) and $search.
	Description string `json:"description,omitempty"`

	// Specifies whether Microsoft has disabled the registered application.
	// Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement
	// (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services
	// Agreement).
	// Supports $filter (eq, ne, NOT).
	DisabledByMicrosoftStatus string `json:"disabledByMicrosoftStatus,omitempty"`

	// The display name for the service principal.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith), $search, and $orderBy.
	DisplayName string `json:"displayName,omitempty"`

	// Home page or landing page of the application.
	Homepage string `json:"homepage,omitempty"`

	// Basic profile information of the acquired application such as app's marketing, support, terms of service and
	// privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent
	// experience.
	// Supports $filter (eq, ne, NOT, ge, le).
	Info InformationalUrl `json:"info,omitempty"`

	// The collection of key credentials associated with the service principal.
	// Not nullable.
	// Supports $filter (eq, NOT, ge, le).
	KeyCredentials []KeyCredential `json:"keyCredentials,omitempty"`

	// Specifies the URL where the service provider redirects the user to Azure AD to authenticate.
	// Azure AD uses the URL to launch the application from Microsoft 365 or the Azure AD My Apps. When blank, Azure AD
	// performs IDP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the
	// application from Microsoft 365, the Azure AD My Apps, or the Azure AD SSO URL.
	LoginUrl string `json:"loginUrl,omitempty"`

	// Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect
	// front-channel, back-channel or SAML logout protocols.
	LogoutUrl string `json:"logoutUrl,omitempty"`

	// Free text field to capture information about the service principal, typically used for operational purposes.
	// Maximum allowed size is 1024 characters.
	Notes string `json:"notes,omitempty"`

	// Specifies the list of email addresses where Azure AD sends a notification when the active certificate is near the
	// expiration date. This is only for the certificates used to sign the SAML token issued for Azure AD Gallery
	// applications.
	NotificationEmailAddresses []string `json:"notificationEmailAddresses,omitempty"`

	// The delegated permissions exposed by the application.
	OAuth2PermissionScopes []PermissionScope `json:"oauth2PermissionScopes,omitempty"`

	// The collection of password credentials associated with the application.
	// Not nullable.
	PasswordCredentials []PasswordCredential `json:"passwordCredentials,omitempty"`

	// Specifies the single sign-on mode configured for this application.
	// Azure AD uses the preferred single sign-on mode to launch the application from Microsoft 365 or the Azure AD My Apps.
	// The supported values are password, saml, notSupported, and oidc.
	PreferredSingleSignOnMode string `json:"preferredSingleSignOnMode,omitempty"`

	// The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that
	// OAuth 2.0 authorization codes and access tokens are sent to for the associated application.
	// Not nullable.
	ReplyUrls []string `json:"replyUrls,omitempty"`

	// The collection for settings related to saml single sign-on.
	SamlSingleSignOnSettings SamlSingleSignOnSettings `json:"samlSingleSignOnSettings,omitempty"`

	// Contains the list of identifiersUris, copied over from the associated application.
	// Additional values can be added to hybrid applications.
	// These values can be used to identify the permissions exposed by this app within Azure AD.
	// For example, Client apps can specify a resource URI which is based on the values of this property to acquire an
	// access token, which is the URI returned in the “aud” claim.
	// The any operator is required for filter expressions on multi-valued properties.
	// Not nullable.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	ServicePrincipalNames []string `json:"servicePrincipalNames,omitempty"`

	// Identifies whether the service principal represents an application, a managed identity, or a legacy application.
	// This is set by Azure AD internally.
	ServicePrincipalType enums.ServicePrincipalType `json:"servicePrincipalType,omitempty"`

	// Specifies the Microsoft accounts that are supported for the current application.
	// Read-only.
	SignInAudience enums.SigninAudience `json:"signInAudience,omitempty"`

	// Custom strings that can be used to categorize and identify the service principal. Not nullable.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	Tags []string `json:"tags,omitempty"`

	// Specifies the keyId of a public key from the keyCredentials collection.
	// When configured, Azure AD issues tokens for this application encrypted using the key specified by this property.
	// The application code that receives the encrypted token must use the matching private key to decrypt the token
	// before it can be used for the signed-in user.
	TokenEncryptionKeyId string `json:"tokenEncryptionKeyId,omitempty"`

	// Specifies the verified publisher of the application which this service principal represents.
	VerifiedPublisher VerifiedPublisher `json:"verifiedPublisher,omitempty"`
}

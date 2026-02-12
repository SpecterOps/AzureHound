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
	"github.com/bloodhoundad/azurehound/v2/enums"
)

// Represents an Azure Active Directory user account.
type User struct {
	DirectoryObject

	// A freeform text entry field for the user to describe themselves.
	//
	// Returned only on `$select`
	AboutMe string `json:"aboutMe,omitempty"`

	// `true` if the account is enabled; otherwise `false`. This property is required when a user is created.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	AccountEnabled bool `json:"accountEnabled,omitempty"`

	// Sets the age group of the user.
	//
	// Allowed values: `null`, `minor`, `notAdult` and `adult`
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	AgeGroup enums.AgeGroup `json:"ageGroup,omitempty"`

	// The licenses that are assigned to the user, including inherited (group-based) licenses.
	//
	// Not nullable
	// Returned only on `$select`
	// Supports `$filter` (eq, NOT)
	AssignedLicenses []AssignedLicense `json:"assignedLicenses,omitempty"`

	// The plans that are assigned to the user.
	//
	// Read-Only
	// Not Nullable
	// Returned only on `$select`
	// Supports `$filter` (eq, NOT)
	AssignedPlans []AssignedPlan `json:"assignedPlans,omitempty"`

	// The birthday of the user using ISO 8601 format.
	//
	// Returned only on `$select`
	Birthday string `json:"birthday,omitempty"`

	// The telephone numbers for the user.
	// Note: Although this is a string collection, only one number can be set for this property.
	//
	// Read-only for users synced from on-premises directory
	// Supports `$filter` (eq, NOT)
	BusinessPhones []string `json:"businessPhones,omitempty"`

	// The city in which the user is located.
	//
	// Max length is 128 characters
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	City string `json:"city,omitempty"`

	// The company name which the user is associated. Useful for describing a company for an external user.
	//
	// Max length is 64 characters
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	CompanyName string `json:"companyName,omitempty"`

	// Sets whether conset has been obtained for minors.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	ConsentProvidedForMinor enums.ConsentForMinor `json:"consentProvidedForMinor,omitempty"`

	// The country/region in which the user is located.
	//
	// Max length is 128 characters
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	Country string `json:"country,omitempty"`

	// The created date of the user object.
	//
	// Read-only
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in)
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// Indicates the method through which the user account was created.
	//
	// Read-only
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	CreationType enums.CreationType `json:"creationType,omitempty"`

	// The date and time the user was deleted.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in)
	DeletedDateTime string `json:"deletedDateTime,omitempty"`

	// The name for the department in which the user works.
	//
	// Max length is 64 characters
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in)
	Department string `json:"department,omitempty"`

	// The name displayed in the address book for the user. This is usually the combination of the user's first name,
	// middle initial and last name. Required on creation and cannot be cleared during updates.
	//
	// Max length is 256 characters
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith), `$orderBy` and `$search`
	DisplayName string `json:"displayName,omitempty"`

	// The data and time the user was hired or will start work in case of a future hire.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in)
	EmployeeHireDate string `json:"employeeHireDate,omitempty"`

	// The employee identifier assigned to the user bu the organization.
	//
	// Returned only on `$select`
	// Supports `filter` (eq,ne,NOT,ge,le,in,startsWith)
	EmployeeId string `json:"employeeId,omitempty"`

	// Represents organization data (e.g. division and costCenter) associated with a user.
	//
	// Returned only in `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in)
	EmployeeOrgData EmployeeOrgData `json:"employeeOrgData,omitempty"`

	// Captures enterprise worker type.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	EmployeeType string `json:"employeeType,omitempty"`

	// For an external user invited to the tenant, this represents the invited user's invitation status.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	ExternalUserState enums.ExternalUserState `json:"externalUserState,omitempty"`

	// Shows the timestamp for the latest change to the {@link ExternalUserState} property.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,in)
	ExternalUserStateChangeDateTime string `json:"externalUserStateChangeDateTime,omitempty"`

	// The fax number of the user.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	FaxNumber string `json:"faxNumber,omitempty"`

	// The given name (first name) of the user.
	//
	// Max length is 64 characters
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	GivenName string `json:"givenName,omitempty"`

	// The hire date of the user using ISO 8601 format.
	// Note: This property is specific to SharePoint Online. Use {@link EmployeeHireDate} to set or update.
	//
	// Returned only on `$select`
	HireDate string `json:"hireDate,omitempty"`

	// Represents the identities that can be used to sign in to this user account.
	//
	// Returned only on `$select`
	// Supports `$filter` (eq) only where the SignInType is not `userPrincipalName`
	Identities []ObjectIdentity `json:"identities,omitempty"`

	// The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for this user.
	//
	// Read-only
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,startsWith)
	ImAddresses []string `json:"imAddresses,omitempty"`

	// A list for the user to describe their interests.
	//
	// Returned only on `$select`
	Interests []string `json:"interests,omitempty"`

	// The user's job title.
	//
	// Max length is 128 characters
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	JobTitle string `json:"jobTitle,omitempty"`

	// The time when this Azure AD user last changed their password or when their password was created using ISO 8601
	// format in UTC time.
	//
	// Returned only on `$select`
	LastPasswordChangeDateTime string `json:"lastPasswordChangeDateTime,omitempty"`

	// The last time the user successfully signed in, in ISO 8601 format (UTC time).
	//
	// This value is sourced from the Microsoft Graph `signInActivity` dictionary.
	// To populate it, the client must include `signInActivity` in `$select`.
	LastSuccessfulSignInDateTime string `json:"lastSuccessfulSignInDateTime,omitempty"`

	// Used by enterprise applications to determine the legal age group of the user.
	//
	// Returned only on `$select`
	LegalAgeGroupClassification enums.LegalAgeGroup `json:"legalAgeGroupClassification,omitempty"`

	// State of license assignments for this user.
	//
	// Read-only
	// Returned only on `$select`
	LicenseAssignmentStates []LicenseAssignmentState `json:"licenseAssignmentStates,omitempty"`

	// The SMTP address for the user.
	//
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith,endsWith)
	Mail string `json:"mail,omitempty"`

	// Settings for the primary mailbox of the signed-in user.
	//
	// Returned only on `$select`
	MailboxSettings MailboxSettings `json:"mailboxSettings,omitempty"`

	// The mail alias for the user.
	//
	// Max length is 64 characters
	// Returned only on `$select`
	// Supports `$filter` (eq,ne,NOT,ge,le,in,startsWith)
	MailNickname string `json:"mailNickname,omitempty"`

	// The primary cellular telephone number for the user. Read-only for users synced from on-premises directory.
	// Maximum length is 64 characters.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	MobilePhone string `json:"mobilePhone,omitempty"`

	// The URL for the user's personal site.
	// Returned only on $select.
	MySite string `json:"mySite,omitempty"`

	// The office location in the user's place of business.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	OfficeLocation string `json:"officeLocation,omitempty"`

	// Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers
	// who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect.
	// Read-only.
	// Returned only on $select.
	OnPremisesDistinguishedName string `json:"onPremisesDistinguishedName,omitempty"`

	// Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory.
	// The property is only populated for customers who are synchronizing their on-premises directory to Azure Active
	// Directory via Azure AD Connect.
	// Read-only.
	// Returned only on $select.
	OnPremisesDomainName string `json:"onPremisesDomainName,omitempty"`

	// Contains extensionAttributes 1-15 for the user
	// Note that the individual extension attributes are neither selectable nor filterable.
	// For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is
	// read-only.
	// For a cloud-only user (where onPremisesSyncEnabled is false), these properties may be set during creation or
	// update. These extension attributes are also known as Exchange custom attributes 1-15.
	// Returned only on $select. Supports $filter (eq, NOT, ge, le, in).
	OnPremisesExtensionAttributes OnPremisesExtensionAttributes `json:"onPremisesExtensionAttributes,omitempty"`

	// This property is used to associate an on-premises Active Directory user account to their Azure AD user object.
	// This property must be specified when creating a new user account in the Graph if you are using a federated domain
	// for the user's userPrincipalName (UPN) property.
	// NOTE: The $ and _ characters cannot be used when specifying this property.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in)
	OnPremisesImmutableId string `json:"onPremisesImmutableId,omitempty"`

	// Indicates the last time at which the object was synced with the on-premises directory;
	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	// Read-only.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	OnPremisesLastSyncDateTime string `json:"onPremisesLastSyncDateTime,omitempty"`

	// Errors when using Microsoft synchronization product during provisioning.
	// Returned only on $select. Supports $filter (eq, NOT, ge, le).
	OnPremisesProvisioningErrors []OnPremisesProvisioningError `json:"onPremisesProvisioningErrors,omitempty"`

	// Contains the on-premises samAccountName synchronized from the on-premises directory.
	// The property is only populated for customers who are synchronizing their on-premises directory to Azure Active
	// Directory via Azure AD Connect.
	// Read-only.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	OnPremisesSamAccountName string `json:"onPremisesSamAccountName,omitempty"`

	// Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the
	// cloud.
	// Read-only.
	// Returned only on $select.
	OnPremisesSecurityIdentifier string `json:"onPremisesSecurityIdentifier,omitempty"`

	// true if this object is synced from an on-premises directory; false if this object was originally synced from an
	// on-premises directory but is no longer synced; null if this object has never been synced from an on-premises
	// directory (default).
	// Read-only.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, in).
	OnPremisesSyncEnabled bool `json:"onPremisesSyncEnabled,omitempty"`

	// Contains the on-premises userPrincipalName synchronized from the on-premises directory.
	// The property is only populated for customers who are synchronizing their on-premises directory to Azure Active
	// Directory via Azure AD Connect.
	// Read-only.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	OnPremisesUserPrincipalName string `json:"onPremisesUserPrincipalName,omitempty"`

	// A list of additional email addresses for the user; for example: ["bob@contoso.com", "Robert@fabrikam.com"].
	// NOTE: This property cannot contain accent characters.
	// Returned only on $select.
	// Supports $filter (eq, NOT, ge, le, in, startsWith).
	OtherMails []string `json:"otherMails,omitempty"`

	// Specifies password policies for the user. This value is an enumeration with one possible value being
	// DisableStrongPassword, which allows weaker passwords than the default policy to be specified.
	// DisablePasswordExpiration can also be specified.
	// The two may be specified together; for example: DisablePasswordExpiration, DisableStrongPassword.
	// Returned only on $select.
	// Supports $filter (ne, NOT).
	PasswordPolicies string `json:"passwordPolicies,omitempty"`

	// Specifies the password profile for the user.
	// The profile contains the user’s password. This property is required when a user is created. The password in the
	// profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong
	// password is required.
	//
	// NOTE: For Azure B2C tenants, the forceChangePasswordNextSignIn property should be set to false and instead use
	// custom policies and user flows to force password reset at first logon. See Force password reset at first logon.
	//
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, in).
	PasswordProfile PasswordProfile `json:"passwordProfile,omitempty"`

	// A list for the user to enumerate their past projects.
	// Returned only on $select.
	PastProjects []string `json:"pastProjects,omitempty"`

	// The postal code for the user's postal address. The postal code is specific to the user's country/region. In the
	// United States of America, this attribute contains the ZIP code.
	// Maximum length is 40 characters.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	PostalCode string `json:"postalCode,omitempty"`

	// The preferred data location for the user.
	PreferredDataLocation string `json:"preferredDataLocation,omitempty"`

	// The preferred language for the user. Should follow ISO 639-1 Code; for example en-US.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith)
	PreferredName string `json:"preferredName,omitempty"`

	// The plans that are provisioned for the user.
	// Read-only.
	// Not nullable.
	// Returned only on $select.
	// Supports $filter (eq, NOT, ge, le).
	ProvisionedPlans []ProvisionedPlan `json:"provisionedPlans,omitempty"`

	// For example: ["SMTP: bob@contoso.com", "smtp: bob@sales.contoso.com"].
	// For Azure AD B2C accounts, this property has a limit of ten unique addresses.
	// Read-only,
	// Not nullable.
	// Returned only on $select.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	ProxyAddresses []string `json:"proxyAddresses,omitempty"`

	// Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications
	// will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access
	// APIs such as Microsoft Graph). If this happens, the application will need to acquire a new refresh token by
	// making a request to the authorize endpoint.
	// Returned only on $select.
	// Read-only.
	RefreshTokensValidFromDateTime string `json:"refreshTokensValidFromDateTime,omitempty"`

	// A list for the user to enumerate their responsibilities.
	// Returned only on $select
	Responsibilities []string `json:"responsibilities,omitempty"`

	// A list for the user to enumerate the schools they have attended.
	// Returned only on $select.
	Schools []string `json:"schools,omitempty"`

	// true if the Outlook global address list should contain this user, otherwise false. If not set, this will be
	// treated as true. For users invited through the invitation manager, this property will be set to false.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, in).
	ShowInAddressList bool `json:"showInAddressList,omitempty"`

	// A list for the user to enumerate their skills.
	// Returned only on $select.
	Skills []string `json:"skills,omitempty"`

	// Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications
	// will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access
	// APIs such as Microsoft Graph). If this happens, the application will need to acquire a new refresh token by
	// making a request to the authorize endpoint.
	// Read-only.
	// Returned only on $select.
	SignInSessionsValidFromDateTime string `json:"signInSessionsValidFromDateTime,omitempty"`

	// The signin activity for the User
	SignInActivity SignInActivity `json:"signInActivity,omitempty"`

	// The state or province in the user's address.
	// Maximum length is 128 characters.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	State string `json:"state,omitempty"`

	// The street address of the user's place of business.
	// Maximum length is 1024 characters.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	StreetAddress string `json:"streetAddress,omitempty"`

	// The user's surname (family name or last name). Maximum length is 64 characters.
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	Surname string `json:"surname,omitempty"`

	// A two letter country code (ISO standard 3166). Required for users that will be assigned licenses due to legal
	// requirement to check for availability of services in countries. Examples include: US, JP, and GB.
	// Not nullable.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith).
	UsageLocation string `json:"usageLocation,omitempty"`

	// The user principal name (UPN) of the user.
	// The UPN is an Internet-style login name for the user based on the Internet standard RFC 822. By convention, this
	// should map to the user's email name. The general format is alias@domain, where domain must be present in the
	// tenant's collection of verified domains. This property is required when a user is created. The verified domains
	// for the tenant can be accessed from the verifiedDomains property of organization.
	//
	// NOTE: This property cannot contain accent characters.
	//
	// Returned by default.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith, endsWith) and $orderBy.
	UserPrincipalName string `json:"userPrincipalName,omitempty"`

	// A string value that can be used to classify user types in your directory, such as Member and Guest.
	// Returned only on $select.
	// Supports $filter (eq, ne, NOT, in).
	UserType string `json:"userType,omitempty"`
}

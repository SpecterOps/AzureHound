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
	"fmt"
)

// Represents the Azure Active Directory tenant that the user or application is signed in to
type Organization struct {
	DirectoryObject

	// The collection of service plans associated with the tenant
	AssignedPlans []AssignedPlan `json:"assignedPlans,omitempty"`

	// Telephone number for the organization. Although this is a string collection, only one number can be set for this
	// property.
	BusinessPhones []string `json:"businessPhones,omitempty"`

	// City name of the address for the organization.
	City string `json:"city,omitempty"`

	// Country or region name of the address for the organization.
	Country string `json:"country,omitempty"`

	// Country or region abbreviation for the organization in ISO 3166-2 format.
	CountryLetterCode string `json:"countryLetterCode,omitempty"`

	// Timestamp of when the organization was created. The value cannot be modified and is automatically populated when
	// the organization is created. The Timestamp type represents date and time information using ISO 8601 format and is
	// always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	//
	// Read-only.
	CreatedDateTime string `json:"createdDateTime,omitempty"`

	// Represents date and time of when the Azure AD tenant was deleted using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	//
	// Read-only.
	DeletedDateTime string `json:"deletedDateTime,omitempty"`

	// The display name for the tenant.
	DisplayName string `json:"displayName,omitempty"`

	// `true` if organization is Multi-Geo enabled; false if organization is not Multi-Geo enabled;
	//
	// null (default).
	// Read-only.
	IsMultipleDataLocationsForServicesEnabled bool `json:"isMultipleDataLocationsForServicesEnabled,omitempty"`

	MarketingNotificationEmails []string `json:"marketingNotificationEmails,omitempty"`

	// The time and date at which the tenant was last synced with the on-premises directory.
	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.
	//
	// Read-only.
	OnPremisesLastSyncDateTime string `json:"onPremisesLastSyncDateTime,omitempty"`

	// `true` if this object is synced from an on-premises directory;
	// `false` if this object was originally synced from an on-premises directory but is no longer synced.
	// `null` if this object has never been synced from an on-premises directory (default).
	OnPremisesSyncEnabled *bool `json:"onPremisesSyncEnabled,omitempty"`

	// Postal code of the address for the organization.
	PostalCode string `json:"postalCode,omitempty"`

	// The preferred language for the organization.
	// Should follow ISO 639-1 Code; for example, `en`.
	PreferredLanguage string `json:"preferredLanguage,omitempty"`

	// The privacy profile of an organization.
	PrivacyProfile PrivacyProfile `json:"privacyProfile,omitempty"`

	ProvisionedPlans []ProvisionedPlan `json:"provisionedPlans,omitempty"`

	SecurityComplianceNotificationMails []string `json:"securityComplianceNotificationMails,omitempty"`

	SecurityComplianceNotificationPhones []string `json:"securityComplianceNotificationPhones,omitempty"`

	// State name of the address for the organization
	State string `json:"state,omitempty"`

	// Street name of the address for the organization.
	Street string `json:"streetAddress,omitempty"`

	TechnicalNotificationMails []string `json:"technicalNotificationMails,omitempty"`

	// The tenant type. Only available for 'Home' TenantCategory
	TenantType string `json:"tenantType,omitempty"`

	// The collection of domains associated with this tenant.
	VerifiedDomains []VerifiedDomain `json:"verifiedDomains,omitempty"`
}

func (s Organization) ToTenant() Tenant {
	var (
		defaultDomain string
		domains       []string
	)

	for _, domain := range s.VerifiedDomains {
		if domain.IsDefault {
			defaultDomain = domain.Name
		}
		domains = append(domains, domain.Name)
	}

	return Tenant{
		Country:       s.Country,
		CountryCode:   s.CountryLetterCode,
		DefaultDomain: defaultDomain,
		DisplayName:   s.DisplayName,
		Domains:       domains,
		Id:            fmt.Sprintf("/tenants/%s", s.ID),
		TenantType:    s.TenantType,
		TenantId:      s.ID,
	}
}

type OrganizationList struct {
	Count    int            `json:"@odata.count,omitempty"`    // The total count of all results
	NextLink string         `json:"@odata.nextLink,omitempty"` // The URL to use for getting the next set of values.
	Value    []Organization `json:"value"`                     // A list of organizations.
}

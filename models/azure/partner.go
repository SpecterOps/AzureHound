// Copyright (C) 2026 Specter Ops, Inc.
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

// Specifies a delegated third-party partner object inside of Entra ID.
// A visual overview can be seen at <https://admin.cloud.microsoft/?#/partners>
// More information <https://learn.microsoft.com/en-us/entra/identity/users/directory-delegated-administration-primer>
type Partner struct {
	// Tenant ID of the external partner
	PartnerTenantId string `json:"partnerTenantId,omitempty"`

	// What kind of external partner?
	// Scraped a bunch of possible values from: <https://hosting.portal.azure.net/iam/Content/Dynamic/iHKmV4w1gTZm.js>
	//
	// Observed variants:
	// - microsoftSupport
	// - breadthPartner
	// - breadthPartnerDelegatedAdmin
	// - syndicatePartner
	// - resellerPartnerDelegatedAdmin (reseller)
	// - valueAddedResellerPartnerDelegatedAdmin (indiret reseller)
	//
	CompanyType string `json:"companyType,omitempty"`

	// The name of the external partner
	CompanyName string `json:"companyName,omitempty"`

	// Link to the partner sales portal
	CommerceUrl string `json:"commerceUrl,omitempty"`

	// Link to the partner help portal
	HelpUrl string `json:"helpUrl,omitempty"`

	// Link to the partner support portal
	// Unsure how this is different from the `HelpUrl`
	SupportUrl string `json:"supportUrl,omitempty"`

	// List of telephone numbers used for support
	SupportTelephones []string `json:"supportTelephones,omitempty"`

	// List of e-mail addresses used for support
	SupportEmails []string `json:"supportEmails,omitempty"`

	// What type of contract does the current tenant have with the partner?
	// Scraped values from: <https://res.cdn.office.net/admincenter/admin-main/inline/inline.en.c957b871d96842c5.chunk.js>
	//
	// Observed variants:
	// - resellerPartnerContract
	// - breadthPartnerContract
	ContractType string `json:"contractType,omitempty"`

	// List of Role ID's
	// Unsure what type this is since the observed versions had no data
	// Assuming string since that's the most logical based on the naming convention
	RoleIDs []string `json:"roleIds,omitempty"`
}

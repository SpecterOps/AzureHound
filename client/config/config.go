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

package config

import (
	"strings"

	"github.com/bloodhoundad/azurehound/v2/constants"
)

type Config struct {
	ApplicationId           string   // The Application Id that the  Azure app registration portal assigned when the app was registered.
	Authority               string   // The Azure ActiveDirectory Authority URL
	ClientSecret            string   // The Application Secret that was generated for the app in the app registration portal.
	ClientCert              string   // The certificate uploaded to the app registration portal."
	ClientKey               string   // The key for a certificate uploaded to the app registration portal."
	ClientKeyPass           string   // The passphrase to use in conjuction with the associated key of a certificate uploaded to the app registration portal."
	Graph                   string   // The Microsoft Graph URL
	JWT                     string   // The JSON web token that will be used to authenticate requests sent to Azure APIs
	Management              string   // The Azure ResourceManager URL
	MgmtGroupId             []string // The Management Group Id to use as a filter
	ManagedIdentity         bool     // If true then the client will use a managed identity to authenticate to Azure APIs
	ManagedIdentitySDK      bool     // If true then the client will use a managed identity sdk to authenticate to Azure APIs
	ManagedIdentityClientId string   // Client ID of user-assigned managed idenity used to authenticate via managed identity SDK
	Password                string   // The password associated with the user principal name associated with the Azure portal.
	ProxyUrl                string   // The forward proxy url
	RefreshToken            string   // The refresh token that will be used to authenticate requests sent to Azure APIs
	Region                  string   // The region of the Azure Cloud deployment.
	SubscriptionId          []string // The Subscription Id(s) to use as a filter
	Tenant                  string   // The directory tenant that you want to request permission from. This can be in GUID or friendly name format
	Username                string   // The user principal name associated with the Azure portal.
}

func AuthorityUrl(region string, defaultUrl string) string {
	switch region {
	case constants.China:
		return constants.AzureChina().ActiveDirectoryAuthority
	case constants.Cloud:
		return constants.AzureCloud().ActiveDirectoryAuthority
	case constants.Germany:
		return constants.AzureGermany().ActiveDirectoryAuthority
	case constants.USGovL4:
		return constants.AzureUSGovernment().ActiveDirectoryAuthority
	case constants.USGovL5:
		return constants.AzureUSGovernmentL5().ActiveDirectoryAuthority
	default:
		return defaultUrl
	}
}

func (s Config) AuthorityUrl() string {
	return AuthorityUrl(s.Region, s.Authority)
}

func GraphUrl(region string, defaultUrl string) string {
	switch region {
	case constants.China:
		return constants.AzureChina().MicrosoftGraphUrl
	case constants.Cloud:
		return constants.AzureCloud().MicrosoftGraphUrl
	case constants.Germany:
		return constants.AzureGermany().MicrosoftGraphUrl
	case constants.USGovL4:
		return constants.AzureUSGovernment().MicrosoftGraphUrl
	case constants.USGovL5:
		return constants.AzureUSGovernmentL5().MicrosoftGraphUrl
	default:
		return defaultUrl
	}
}

func (s Config) GraphUrl() string {
	return strings.TrimSuffix(GraphUrl(s.Region, s.Graph), "/")
}

func ResourceManagerUrl(region string, defaultUrl string) string {
	switch region {
	case constants.China:
		return constants.AzureChina().ResourceManagerUrl
	case constants.Cloud:
		return constants.AzureCloud().ResourceManagerUrl
	case constants.Germany:
		return constants.AzureGermany().ResourceManagerUrl
	case constants.USGovL4:
		return constants.AzureUSGovernment().ResourceManagerUrl
	case constants.USGovL5:
		return constants.AzureUSGovernmentL5().ResourceManagerUrl
	default:
		return defaultUrl
	}
}

func (s Config) ResourceManagerUrl() string {
	return strings.TrimSuffix(ResourceManagerUrl(s.Region, s.Graph), "/")
}

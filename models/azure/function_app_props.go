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

type FunctionAppProperties struct {
	AvailabilityState           enums.SiteAvailabilityState `json:"availabilityState,omitempty"`
	ClientAffinityEnabled       bool                        `json:"clientAffinityEnabled,omitempty"`
	ClientCertEnabled           bool                        `json:"clientCertEnabled,omitempty"`
	ClientCertExclusionPaths    string                      `json:"clientCertExclusionPaths,omitempty"`
	ClientCertMode              enums.ClientCertMode        `json:"clientCertMode,omitempty"`
	CloningInfo                 CloningInfo                 `json:"cloningInfo,omitempty"`
	ContainerSize               int                         `json:"containerSize,omitempty"`
	CustomDomainVerificationId  string                      `json:"customDomainVerificationId,omitempty"`
	DailyMemoryTimeQuota        int                         `json:"dailyMemoryTimeQuota,omitempty"`
	DefaultHostName             string                      `json:"defaultHostName,omitempty"`
	Enabled                     bool                        `json:"enabled,omitempty"`
	EnabledHostnames            []string                    `json:"enabledHostnames,omitempty"`
	HostingEnvironmentProfile   HostingEnvironmentProfile   `json:"hostingEnvironmentProfile,omitempty"`
	Hostnames                   []string                    `json:"hostNames,omitempty"`
	HostNamesDisabled           bool                        `json:"hostNamesDisabled,omitempty"`
	HostNameSslStates           []HostNameSslState          `json:"hostNameSslStates,omitempty"`
	HttpsOnly                   bool                        `json:"httpsOnly,omitempty"`
	HyperV                      bool                        `json:"hyperV,omitempty"`
	InProgressOperationId       string                      `json:"inProgressOperationId,omitempty"`
	IsDefaultContainer          bool                        `json:"isDefaultContainer,omitempty"`
	IsXenon                     bool                        `json:"isXenon,omitempty"`
	KeyVaultReferenceIdentity   string                      `json:"keyVaultReferenceIdentity,omitempty"`
	LastModifiedTimeUTC         string                      `json:"lastModifiedTimeUtc,omitempty"`
	MaxNumberOfWorkers          int                         `json:"maxNumberOfWorkers,omitempty"`
	OutboundIpAddresses         string                      `json:"outboundIpAddresses,omitempty"`
	PossibleOutboundIpAddresses string                      `json:"possibleOutboundIpAddresses,omitempty"`
	PublicNetworkAccess         string                      `json:"publicNetworkAccess,omitempty"`
	RedundancyMode              enums.RedundancyMode        `json:"redundancyMode,omitempty"`
	RepositorySiteName          string                      `json:"repositorySiteName,omitempty"`
	Reserved                    bool                        `json:"reserved,omitempty"`
	ResourceGroup               string                      `json:"resourceGroup,omitempty"`
	ScmSiteAlsoStopped          bool                        `json:"scmSiteAlsoStopped,omitempty"`
	ServerFarmId                string                      `json:"serverFarmId,omitempty"`
	SiteConfig                  SiteConfig                  `json:"siteConfig,omitempty"`
	SlotSwapStatus              SlotSwapStatus              `json:",omitempty"`
	State                       string                      `json:"state,omitempty"`
	StorageAccountRequired      bool                        `json:"storageAccountRequired,omitempty"`
	SuspendedTill               string                      `json:"suspendedTill,omitempty"`
	TargetSwapSlot              string                      `json:"targetSwapSlot,omitempty"`
	TrafficManagerHostNames     []string                    `json:"trafficManagerHostNames,omitempty"`
	UsageState                  enums.UsageState            `json:"usageState,omitempty"`
	VirtualNetworkSubnetId      string                      `json:"virtualNetworkSubnetId,omitempty"`
	VnetContentShareEnabled     bool                        `json:"vnetContentShareEnabled,omitempty"`
	VnetImagePullEnabled        bool                        `json:"vnetImagePullEnabled,omitempty"`
	VnetRouteAllEnabled         bool                        `json:"vnetRouteAllEnabled,omitempty"`

	// Following elements have been found in testing within the returned object, but not present in the official documentation
	AdminEnabled                bool   `json:"adminEnabled,omitempty"`
	ComputeMode                 string `json:"computeMode,omitempty"`
	ContainerAllocationSubnet   string `json:"containerAllocationSubnet,omitempty"`
	ContentAvailabilityState    string `json:"contentAvailabilityState,omitempty"`
	FtpsHostName                string `json:"ftpsHostName,omitempty"`
	FtpUsername                 string `json:"ftpUsername,omitempty"`
	InboundIPAddress            string `json:"inboundIpAddress,omitempty"`
	Kind                        string `json:"kind,omitempty"`
	Name                        string `json:"name,omitempty"`
	PossibleInboundIpAddresses  string `json:"possibleInboundIpAddresses,omitempty"`
	PrivateEndpointConnections  string `json:"privateEndpointConnections,omitempty"`
	RuntimeAvailabilityState    string `json:"runtimeAvailabilityState,omitempty"`
	SelfLink                    string `json:"selfLink,omitempty"`
	StorageRecoveryDefaultState string `json:"storageRecoveryDefaultState,omitempty"`
}

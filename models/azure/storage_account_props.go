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

type StorageAccountProperties struct {
	AccessTier                            enums.StorageAccountAccessTier        `json:"accessTier,omitempty"`
	AllowBlobPublicAccess                 bool                                  `json:"allowBlobPublicAccess,omitempty"`
	AllowCrossTenantReplication           bool                                  `json:"allowCrossTenantReplication,omitempty"`
	AllowSharedKeyAccess                  bool                                  `json:"allowSharedKeyAccess,omitempty"`
	AllowedCopyScope                      enums.AllowedCopyScope                `json:"allowedCopyScope,omitempty"`
	AzureFilesIdentityBasedAuthentication AzureFilesIdentityBasedAuthentication `json:"azureFilesIdentityBasedAuthentication,omitempty"`
	BlobRestoreStatus                     BlobRestoreStatus                     `json:"blobRestoreStatus,omitempty"`
	CreationTime                          string                                `json:"creationTime,omitempty"`
	CustomDomain                          StorageAccountCustomDomain            `json:"customDomain,omitempty"`
	DefaultToOAuthAuthentication          bool                                  `json:"defaultToOAuthAuthentication,omitempty"`
	DnsEndpointType                       enums.DnsEndpointType                 `json:"dnsEndpointType,omitempty"`
	Encryption                            StorageAccountEncryptionProperties    `json:"encryption,omitempty"`
	FailoverInProgress                    bool                                  `json:"failoverInProgress,omitempty"`
	GeoReplicationStats                   GeoReplicationStats                   `json:"geoReplicationStats,omitempty"`
	ImmutableStorageWithVersioning        ImmutableStorageAccount               `json:"immutableStorageWithVersioning,omitempty"`
	IsHnsEnabled                          bool                                  `json:"isHnsEnabled,omitempty"`
	IsLocalUserEnabled                    bool                                  `json:"isLocalUserEnabled,omitempty"`
	IsNfsV3Enabled                        bool                                  `json:"isNfsV3Enabled,omitempty"`
	IsSftpEnabled                         bool                                  `json:"isSftpEnabled,omitempty"`
	KeyCreationTime                       StorageAccountKeyCreationTime         `json:"keyCreationTime,omitempty"`
	KeyPolicy                             StorageAccountKeyPolicy               `json:"keyPolicy,omitempty"`
	LargeFileSharesState                  enums.GenericEnabledDisabled          `json:"largeFileSharesState,omitempty"`
	LastGeoFailoverTime                   string                                `json:"lastGeoFailoverTime,omitempty"`
	MinimumTlsVersion                     enums.MinimumTlsVersion               `json:"minimumTlsVersion,omitempty"`
	NetworkAcls                           NetworkRuleSet                        `json:"networkAcls,omitempty"`
	PrimaryEndpoints                      Endpoints                             `json:"primaryEndpoints,omitempty"`
	PrimaryLocation                       string                                `json:"primaryLocation,omitempty"`
	PrivateEndpointConnections            []PrivateEndpointConnection           `json:"privateEndpointConnections"`
	ProvisioningState                     enums.ProvisioningState               `json:"provisioningState,omitempty"`
	PublicNetworkAccess                   enums.GenericEnabledDisabled          `json:"availabilitySet,omitempty"`
	RoutingPreference                     RoutingPreference                     `json:"routingPreference,omitempty"`
	SasPolicy                             SasPolicy                             `json:"sasPolicy,omitempty"`
	SecondaryEndpoints                    Endpoints                             `json:"secondaryEndpoints,omitempty"`
	SecondaryLocation                     string                                `json:"secondaryLocation,omitempty"`
	StatusOfPrimary                       enums.AccountStatus                   `json:"statusOfPrimary,omitempty"`
	StatusOfSecondary                     enums.AccountStatus                   `json:"statusOfSecondary,omitempty"`
	StorageAccountSkuConversionStatus     StorageAccountSkuConversionStatus     `json:"storageAccountSkuConversionStatus,omitempty"`

	SupportsHttpsTrafficOnly bool `json:"supportsHttpsTrafficOnly,omitempty"`
}

type StorageAccountCustomDomain struct {
	Name             string `json:"name,omitempty"`
	UseSubDomainName bool   `json:"useSubDomainName,omitempty"`
}

type StorageAccountKeyCreationTime struct {
	Key1 string `json:"key1,omitempty"`
	Key2 string `json:"key2,omitempty"`
}

type StorageAccountKeyPolicy struct {
	KeyExpirationPeriodInDays int `json:"keyExpirationPeriodInDays,omitempty"`
}

type StorageAccountLargeFileSharesState struct {
}

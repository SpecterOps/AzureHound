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

type AutomationAccountEncryptionProperties struct {
	Identity           ManagedIdentity               `json:"identity,omitempty"`
	KeySource          enums.EncryptionKeySourceType `json:"keySource,omitempty"`
	KeyVaultProperties KeyVaultProperties            `json:"keyVaultProperties,omitempty"`
}

type StorageAccountEncryptionProperties struct {
	Identity                        StorageAccountEncryptionIdentity `json:"identity,omitempty"`
	KeySource                       enums.EncryptionKeySourceType    `json:"keySource,omitempty"`
	Keyvaultproperties              KeyVaultProperties               `json:"keyvaultproperties,omitempty"`
	RequireInfrastructureEncryption bool                             `json:"requireInfrastructureEncryption,omitempty"`
	Services                        EncryptionServices               `json:"services,omitempty"`
}

type KeyVaultProperties struct {
	CurrentVersionedKeyExpirationTimestamp string `json:"currentVersionedKeyExpirationTimestamp,omitempty"`
	CurrentVersionedKeyIdentifier          string `json:"currentVersionedKeyIdentifier,omitempty"`
	KeyName                                string `json:"keyName,omitempty"`
	KeyVersion                             string `json:"keyVersion,omitempty"`
	KeyvaultUri                            string `json:"keyvaultUri,omitempty"`
	LastKeyRotationTimestamp               string `json:"lastKeyRotationTimestamp,omitempty"`
}

type StorageAccountEncryptionIdentity struct {
	FederatedIdentityClientId string `json:"federatedIdentityClientId,omitempty"`
	UserAssignedIdentity      string `json:"userAssignedIdentity,omitempty"`
}

type EncryptionServices struct {
	Blob  EncryptionService `json:"blob,omitempty"`
	File  EncryptionService `json:"file,omitempty"`
	Queue EncryptionService `json:"queue,omitempty"`
	Table EncryptionService `json:"table,omitempty"`
}

type EncryptionService struct {
	Enabled         bool                    `json:"enabled,omitempty"`
	KeyType         enums.EncryptionKeyType `json:"keyType,omitempty"`
	LastEnabledTime string                  `json:"lastEnabledTime,omitempty"`
}

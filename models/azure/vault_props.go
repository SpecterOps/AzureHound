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

// Properties of the vault
type VaultProperties struct {
	// An array of 0 to 1024 identities that have access to the key vault. All identities in the array must use the same
	// tenant ID as the key vault's tenant ID. When createMode is set to recover, access policies are not required.
	// Otherwise, access policies are required.
	AccessPolicies []AccessPolicyEntry `json:"accessPolicies,omitempty"`

	// The vault's create mode to indicate whether the vault need to be recovered or not.
	CreateMode enums.CreateMode `json:"createMode,omitempty"`

	// Property specifying whether protection against purge is enabled for this vault.
	// Setting this property to true activates protection against purge for this vault and its content - only the
	// Key Vault service may initiate a hard, irrecoverable deletion.
	// The setting is effective only if soft delete is also enabled.
	// Enabling this functionality is irreversible - that is, the property does not accept false as its value.
	EnablePurgeProtection bool `json:"enablePurgeProtection,omitempty"`

	// Property that controls how data actions are authorized.
	// When true, the key vault will use Role Based Access Control (RBAC) for authorization of data actions, and the
	// access policies specified in vault properties will be ignored. When false, the key vault will use the access
	// policies specified in vault properties, and any policy stored on Azure Resource Manager will be ignored. If null
	// or not specified, the vault is created with the default value of false.
	// Note that management actions are always authorized with RBAC.
	EnableRbacAuthorization bool `json:"enableRbacAuthorization,omitempty"`

	// Property to specify whether the 'soft delete' functionality is enabled for this key vault.
	// If it's not set to any value(true or false) when creating new key vault, it will be set to true by default.
	// Once set to true, it cannot be reverted to false.
	EnableSoftDelete bool `json:"enableSoftDelete,omitempty"`

	// Property to specify whether Azure Virtual Machines are permitted to retrieve certificates stored as secrets from
	// the key vault.
	EnabledForDeployment bool `json:"enabledForDeployment,omitempty"`

	// Property to specify whether Azure Disk Encryption is permitted to retrieve secrets from the vault and unwrap keys.
	EnabledForDiskEncryption bool `json:"enabledForDiskEncryption,omitempty"`

	// Property to specify whether Azure Resource Manager is permitted to retrieve secrets from the key vault.
	EnabledForTemplateDeployment bool `json:"enabledForTemplateDeployment,omitempty"`

	// The resource ID of HSM Pool.
	HsmPoolResourceId string `json:"hsmPoolResourceId,omitempty"`

	// Rules governing the accessibility of the key vault from specific network locations.
	NetworkAcls NetworkRuleSet `json:"networkAcls,omitempty"`

	// List of private endpoint connections associated with the key vault.
	PrivateEndpointConnections []PrivateEndpointConnectionItem `json:"privateEndpointConnections,omitempty"`

	// Provisioning state of the vault.
	ProvisioningState enums.VaultProvisioningState `json:"provisioningState,omitempty"`

	// SKU details
	Sku Sku `json:"sku,omitempty"`

	// softDelete data retention days. It accepts >=7 and <=90.
	SoftDeleteRetentionInDays int `json:"softDeleteRetentionInDays,omitempty"`

	// The Azure Active Directory tenant ID that should be used for authenticating requests to the key vault.
	TenantId string `json:"tenantId,omitempty"`

	// The URI of the vault for performing operations on keys and secrets. This property is readonly.
	VaultUri string `json:"vaultUri,omitempty"`
}

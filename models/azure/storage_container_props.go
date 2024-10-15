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

type StorageContainerProperties struct {
	DefaultEncryptionScope         string                         `json:"defaultEncryptionScope,omitempty"`
	Deleted                        bool                           `json:"deleted,omitempty"`
	DeletedTime                    string                         `json:"deletedTime,omitempty"`
	DenyEncryptionScopeOverride    bool                           `json:"denyEncryptionScopeOverride,omitempty"`
	EnableNfsV3AllSquash           bool                           `json:"enableNfsV3AllSquash,omitempty"`
	EnableNfsV3RootSquash          bool                           `json:"enableNfsV3RootSquash,omitempty"`
	HasImmutabilityPolicy          bool                           `json:"hasImmutabilityPolicy,omitempty"`
	HasLegalHold                   bool                           `json:"hasLegalHold,omitempty"`
	ImmutabilityPolicy             ImmutabilityPolicy             `json:"immutabilityPolicy,omitempty"`
	ImmutableStorageWithVersioning ImmutableStorageWithVersioning `json:"immutableStorageWithVersioning,omitempty"`
	LastModifiedTime               string                         `json:"lastModifiedTime,omitempty"`
	LeaseDuration                  enums.LeaseDuration            `json:"leaseDuration,omitempty"`
	LeaseState                     enums.LeaseState               `json:"leaseState,omitempty"`
	LeaseStatus                    enums.LeaseStatus              `json:"leaseStatus,omitempty"`
	LegalHold                      LegalHoldProperties            `json:"legalHold,omitempty"`
	Metadata                       interface{}                    `json:"metadata,omitempty"`
	PublicAccess                   enums.PublicAccess             `json:"publicAccess,omitempty"`
	RemainingRetentionDays         int                            `json:"remainingRetentionDays,omitempty"`
	Version                        string                         `json:"version,omitempty"`
}

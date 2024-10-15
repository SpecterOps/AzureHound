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

// Represents a device registered in the organization. Devices are created in the cloud using the Device Registration
// Service or by Intune. They're used by conditional access policies for multi-factor authentication. These devices can
// range from desktop and laptop machines to phones and tablets.
// For more detail see https://docs.microsoft.com/en-us/graph/api/resources/device?view=graph-rest-1.0
type Device struct {
	DirectoryObject

	// true if the account is enabled; otherwise, false. Required. Default is true.
	// Supports $filter (eq, ne, NOT, in).
	// Only callers in Global Administrator and Cloud Device Administrator roles can set this property.
	AccountEnabled bool `json:"accountEnabled,omitempty"`

	// For internal use only. Not nullable. Supports $filter (eq, NOT, ge, le).
	AlternativeSecurityIds []AlternativeSecurityId `json:"alternativeSecurityIds,omitempty"`

	// The timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
	// Supports $filter (eq, ne, NOT, ge, le) and $orderBy.
	ApproximateLastSignInDateTime string `json:"approximateLastSignInDateTime,omitempty"`

	// The timestamp when the device is no longer deemed compliant.
	// The timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.
	ComplianceExpirationDateTime string `json:"complianceExpirationDateTime,omitempty"`

	// Unique identifier set by Azure Device Registration Service at the time of registration.
	// Supports $filter (eq, ne, NOT, startsWith).
	DeviceId string `json:"deviceId,omitempty"`

	// For internal use only. Set to null.
	DeviceMetadata string `json:"deviceMetadata,omitempty"`

	// For internal use only.
	DeviceVersion int32 `json:"deviceVersion,omitempty"`

	// The display name for the device.
	// Required.
	// Supports $filter (eq, ne, NOT, ge, le, in, startsWith), $search, and $orderBy.
	DisplayName string `json:"displayName,omitempty"`

	// Contains extension attributes 1-15 for the device.
	// The individual extension attributes are not selectable.
	// These properties are mastered in cloud and can be set during creation or update of a device object in Azure AD.
	// Supports $filter (eq, NOT, startsWith).
	ExtensionAttributes OnPremisesExtensionAttributes `json:"onPremisesExtensionAttributes,omitempty"`

	// true if the device complies with Mobile Device Management (MDM) policies; otherwise, false.
	// Read-only.
	// This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices.
	// Supports $filter (eq, ne, NOT).
	IsCompliant bool `json:"isCompliant,omitempty"`

	// true if the device is managed by a Mobile Device Management (MDM) app; otherwise, false.
	// This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices.
	// Supports $filter (eq, ne, NOT).
	IsManaged bool `json:"isManaged,omitempty"`

	// Manufacturer of the device.
	// Read-only.
	Manufacturer string `json:"manufacturer,omitempty"`

	// Application identifier used to register device into MDM.
	// Read-only.
	// Supports $filter (eq, ne, NOT, startsWith).
	MdmAppId string `json:"mdmAppId,omitempty"`

	// Model of the device.
	// Read-only.
	Model string `json:"model,omitempty"`

	// The last time at which the object was synced with the on-premises directory.
	// The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time.
	// For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z Read-only.
	// Supports $filter (eq, ne, NOT, ge, le, in).
	OnPremisesLastSyncDateTime string `json:"onPremisesLastSyncDateTime,omitempty"`

	// true if this object is synced from an on-premises directory; false if this object was originally synced from an
	// on-premises directory but is no longer synced; null if this object has never been synced from an on-premises
	// directory (default).
	// Read-only.
	// Supports $filter (eq, ne, NOT, in).
	OnPremisesSyncEnabled bool `json:"onPremisesSyncEnabled,omitempty"`

	// The type of operating system on the device.
	// Required.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith).
	OperatingSystem string `json:"operatingSystem,omitempty"`

	// The version of the operating system on the device.
	// Required.
	// Supports $filter (eq, ne, NOT, ge, le, startsWith).
	OperatingSystemVersion string `json:"operatingSystemVersion,omitempty"`

	// For internal use only.
	// Not nullable.
	// Supports $filter (eq, NOT, ge, le, startsWith).
	PhysicalIds []string `json:"physicalIds,omitempty"`

	// The profile type of the device.
	ProfileType enums.DeviceProfile `json:"profileType,omitempty"`

	// List of labels applied to the device by the system.
	SystemLabels []string `json:"systemLabels,omitempty"`

	// Type of trust for the joined device.
	// Read-only.
	TrustType enums.TrustType `json:"trustType,omitempty"`
}

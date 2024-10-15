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

type VirtualMachineProperties struct {
	AdditionalCapabilities  AdditionalCapabilities     `json:"additionalCapabilities,omitempty"`
	ApplicationProfile      ApplicationProfile         `json:"applicationProfile,omitempty"`
	AvailabilitySet         SubResource                `json:"availabilitySet,omitempty"`
	BillingProfile          BillingProfile             `json:"billingProfile,omitempty"`
	CapacityReservation     CapacityReservationProfile `json:"capacityReservation,omitempty"`
	DiagnosticsProfile      DiagnosticsProfile         `json:"diagnosticsProfile,omitempty"`
	EvictionPolicy          enums.VMEvictionPolicy     `json:"evictionPolicy,omitempty"`
	ExtensionsTimeBudget    string                     `json:"extensionsTimeBudget,omitempty"`
	HardwareProfile         HardwareProfile            `json:"hardwareProfile,omitempty"`
	Host                    SubResource                `json:"host,omitempty"`
	HostGroup               SubResource                `json:"hostGroup,omitempty"`
	InstanceView            VirtualMachineInstanceView `json:"instanceView,omitempty"`
	LicenseType             string                     `json:"licenseType,omitempty"`
	NetworkProfile          NetworkProfile             `json:"networkProfile,omitempty"`
	OSProfile               OSProfile                  `json:"osProfile,omitempty"`
	PlatformFaultDomain     int                        `json:"platformFaultDomain,omitempty"`
	Priority                enums.VMPriority           `json:"priority,omitempty"`
	ProvisioningState       string                     `json:"provisioningState,omitempty"`
	ProximityPlacementGroup SubResource                `json:"proximityPlacementGroup,omitempty"`
	ScheduledEventsProfile  ScheduledEventsProfile     `json:"scheduledEventsProfile,omitempty"`
	SecurityProfile         SecurityProfile            `json:"securityProfile,omitempty"`
	StorageProfile          StorageProfile             `json:"storageProfile,omitempty"`
	UserData                string                     `json:"userData,omitempty"`
	VirtualMachineScaleSet  SubResource                `json:"virtualMachineScaleSet,omitempty"`
	VMId                    string                     `json:"vmId,omitempty"`
}

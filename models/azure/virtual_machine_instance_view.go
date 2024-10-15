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

// The instance view of a virtual machine.
type VirtualMachineInstanceView struct {
	// Resource id of the dedicated host, on which the virtual machine is allocated through automatic placement, when
	// the virtual machine is associated with a dedicated host group that has automatic placement enabled.
	// Minimum api-version: 2020-06-01.
	AssignedHost string `json:"assignedHost,omitempty"`

	// Boot Diagnostics is a debugging feature which allows you to view Console Output and Screenshot to diagnose VM
	// status.
	// You can easily view the output of your console log.
	// Azure also enables you to see a screenshot of the VM from the hypervisor.
	BootDiagnotics BootDiagnoticsInstanceView `json:"bootDiagnotics,omitempty"`

	// The computer name assigned to the virtual machine.
	ComputerName string `json:"computerName,omitempty"`

	// The virtual machine disk information.
	Disks []DiskInstanceView `json:"disks,omitempty"`

	// The extensions information.
	Extensions []VirtualMachineExtensionInstanceView `json:"extensions,omitempty"`

	// Specifies the HyperVGeneration Type associated with a resource.
	HyperVGeneration enums.HyperVGeneration `json:"hyperVGeneration,omitempty"`

	// The Maintenance Operation status on the virtual machine.
	MaintenanceRedeployStatus MaintenanceRedeployStatus `json:"maintenanceRedeployStatus,omitempty"`

	// The Operating System running on the virtual machine.
	OSName string `json:"osName,omitempty"`

	// The version of Operating System running on the virtual machine.
	OSVersion string `json:"osVersion,omitempty"`

	// [Preview Feature] The status of the virtual machine patch operations.
	PatchStatus VirtualMachinePatchStatus `json:"patchStatus,omitempty"`

	// Specifies the fault domain of the virtual machine.
	PlatformFaultDomain int `json:"platformFaultDomain,omitempty"`

	// Specifies the update domain of the virtual machine.
	PlatformUpdateDomain int `json:"platformUpdateDomain,omitempty"`

	// The remote desktop certificate thumbprint.
	RDPThumbPrint string `json:"rdpThumbPrint,omitempty"`

	// The resource status information.
	Statuses []InstanceViewStatus `json:"statuses,omitempty"`

	// The VM Agent running on the virtual machine.
	VMAgent VirtualMachineAgentInstanceView `json:"vmAgent,omitempty"`

	// The health status for the VM.
	VMHealth VirtualMachineHealthStatus `json:"vmHealth,omitempty"`
}

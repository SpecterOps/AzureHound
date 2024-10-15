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

type VMNetworkInterfaceConfigurationProperties struct {
	// Specify what happens to the network interface when the VM is deleted.
	DeleteOption enums.VMDeleteOption `json:"deleteOption,omitempty"`

	// The dns settings to be applied on the network interfaces.
	DNSSettings VMNetworkInterfaceDNSSettings `json:"dnsSettings,omitempty"`

	// The DSCP resource to be applied to the network interfaces.
	DSCPConfiguration SubResource `json:"dscpConfiguration,omitempty"`

	// Specifies whether the network interface is accelerated networking-enabled.
	EnabledAcceleratedNetworking bool `json:"enabledAcceleratedNetworking,omitempty"`

	// Specifies whether the network is FPGA networking-enabled
	EnableFpga bool `json:"enableFpga,omitempty"`

	// Whether IP forwarding is enabled on this NIC.
	EnableIPForwarding bool `json:"enableIPForwarding,omitempty"`

	// Specifies the IP configurations of the network interface.
	IPConfigurations []VMNetworkInterfaceIPConfig `json:"ipConfigurations,omitempty"`

	// The network security group.
	NetworkSecurityGroup SubResource `json:"networkSecurityGroup,omitempty"`

	// Specifies the primary network interface in case the virtual machine has more than one.
	Primary bool `json:"primary,omitempty"`
}

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

// A set of rules governing the network accessibility of a vault.
type NetworkRuleSet struct {
	// Tells what traffic can bypass network rules. This can be 'AzureServices' or 'None'.
	// If not specified the default is 'AzureServices'.
	Bypass enums.BypassOption `json:"bypass,omitempty"`

	// The default action when no rule from ipRules and from virtualNetworkRules match.
	// This is only used after the bypass property has been evaluated.
	DefaultAction enums.NetworkAction `json:"defaultAction,omitempty"`

	// The list of IP address rules.
	IPRules []IPRule `json:"ipRules,omitempty"`

	// The list of virtual network rules.
	VirtualNetworkRules []VirtualNetworkRule `json:"virtualNetworkRules,omitempty"`
}

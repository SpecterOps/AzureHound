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
	"encoding/json"
	"strings"
)

// The properties of the parent management group.
type DescendantParentGroupInfo struct {
	// The fully qualified ID for the parent management group.
	//
	// For example:
	// - /providers/Microsoft.Management/managementGroups/0000000-0000-0000-0000-000000000000
	Id string `json:"id,omitempty"`
}

// DescendantInfoProperties describes the properties of the management group descendant.
type DescendantInfoProperties struct {
	// The friendly name of the management group.
	DisplayName string `json:"display_name,omitempty"`

	// The properties of the parent management group.
	Parent DescendantParentGroupInfo `json:"parent,omitempty"`
}

// DescendantInfo is a management group descendant.
type DescendantInfo struct {
	// The fully qualified ID for the descendant.
	//
	// For example:
	// - /providers/Microsoft.Management/managementGroups/0000000-0000-0000-0000-000000000000
	// - /subscriptions/0000000-0000-0000-0000-000000000000
	Id string `json:"id,omitempty"`

	// The name of the descendant.
	//
	// For example:
	// - 00000000-0000-0000-0000-000000000000
	Name string `json:"name,omitempty"`

	// The properties of the management group descendant.
	Properties DescendantInfoProperties `json:"properties,omitempty"`

	// The type of the resource.
	//
	// For example:
	// - Microsoft.Management/managementGroups
	// - /subscriptions
	Type string `json:"type,omitempty"`
}

func (s DescendantInfo) MarshalJSON() ([]byte, error) {
	type Alias DescendantInfo
	a := Alias(s)
	a.Id = strings.ToUpper(a.Id)
	a.Properties.Parent.Id = strings.ToUpper(a.Properties.Parent.Id)
	return json.Marshal(a)
}

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

package models

import (
	"encoding/json"
	"strings"

	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

type VirtualMachineRoleAssignment struct {
	RoleAssignment   azure.RoleAssignment `json:"roleAssignment"`
	VirtualMachineId string               `json:"virtualMachineId"`
}

// MarshalJSON uppercases the VirtualMachineId and the RoleAssignment endpoint
// identifiers so the raw (use_raw_object_id) ingest path matches the normalized
// node ObjectIDs. The input is not mutated.
func (s VirtualMachineRoleAssignment) MarshalJSON() ([]byte, error) {
	type Alias VirtualMachineRoleAssignment
	a := Alias(s)
	a.VirtualMachineId = strings.ToUpper(a.VirtualMachineId)
	a.RoleAssignment = UpperRoleAssignment(a.RoleAssignment)
	return json.Marshal(a)
}

type VirtualMachineRoleAssignments struct {
	RoleAssignments  []VirtualMachineRoleAssignment `json:"roleAssignments"`
	VirtualMachineId string                         `json:"virtualMachineId"`
}

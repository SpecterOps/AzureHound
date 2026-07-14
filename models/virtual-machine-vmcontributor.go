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

type VirtualMachineVMContributor struct {
	VMContributor    azure.RoleAssignment `json:"vmContributor"`
	VirtualMachineId string               `json:"virtualMachineId"`
}

func (s VirtualMachineVMContributor) MarshalJSON() ([]byte, error) {
	type Alias VirtualMachineVMContributor
	a := Alias(s)
	a.VirtualMachineId = strings.ToUpper(a.VirtualMachineId)
	a.VMContributor = UpperRoleAssignment(a.VMContributor)
	return json.Marshal(a)
}

type VirtualMachineVMContributors struct {
	VMContributors   []VirtualMachineVMContributor `json:"vmContributors"`
	VirtualMachineId string                        `json:"virtualMachineId"`
}

func (s VirtualMachineVMContributors) MarshalJSON() ([]byte, error) {
	type Alias VirtualMachineVMContributors
	a := Alias(s)
	a.VirtualMachineId = strings.ToUpper(a.VirtualMachineId)
	return json.Marshal(a)
}

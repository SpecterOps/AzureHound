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

type ResourceGroupOwner struct {
	Owner           azure.RoleAssignment `json:"owner"`
	ResourceGroupId string               `json:"resourceGroupId"`
}

func (s ResourceGroupOwner) MarshalJSON() ([]byte, error) {
	type Alias ResourceGroupOwner
	a := Alias(s)
	a.ResourceGroupId = strings.ToUpper(a.ResourceGroupId)
	a.Owner = UpperRoleAssignment(a.Owner)
	return json.Marshal(a)
}

type ResourceGroupOwners struct {
	Owners          []ResourceGroupOwner `json:"owners"`
	ResourceGroupId string               `json:"resourceGroupId"`
}

func (s ResourceGroupOwners) MarshalJSON() ([]byte, error) {
	type Alias ResourceGroupOwners
	a := Alias(s)
	a.ResourceGroupId = strings.ToUpper(a.ResourceGroupId)
	return json.Marshal(a)
}

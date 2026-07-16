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

type ResourceGroupRoleAssignment struct {
	RoleAssignment  azure.RoleAssignment `json:"roleAssignment"`
	ResourceGroupId string               `json:"resourceGroupId"`
}

// MarshalJSON uppercases the ResourceGroupId and the RoleAssignment endpoint
// identifiers so the raw (use_raw_object_id) ingest path matches the normalized
// node ObjectIDs. The input is not mutated.
func (s ResourceGroupRoleAssignment) MarshalJSON() ([]byte, error) {
	type Alias ResourceGroupRoleAssignment
	a := Alias(s)
	a.ResourceGroupId = strings.ToUpper(a.ResourceGroupId)
	a.RoleAssignment = UpperRoleAssignment(a.RoleAssignment)
	return json.Marshal(a)
}

type ResourceGroupRoleAssignments struct {
	RoleAssignments []ResourceGroupRoleAssignment `json:"roleAssignments"`
	ResourceGroupId string                        `json:"resourceGroupId"`
}

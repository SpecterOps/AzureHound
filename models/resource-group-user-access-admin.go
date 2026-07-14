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

type ResourceGroupUserAccessAdmin struct {
	UserAccessAdmin azure.RoleAssignment `json:"userAccessAdmin"`
	ResourceGroupId string               `json:"resourceGroupId"`
}

func (s ResourceGroupUserAccessAdmin) MarshalJSON() ([]byte, error) {
	type Alias ResourceGroupUserAccessAdmin
	a := Alias(s)
	a.ResourceGroupId = strings.ToUpper(a.ResourceGroupId)
	a.UserAccessAdmin = UpperRoleAssignment(a.UserAccessAdmin)
	return json.Marshal(a)
}

type ResourceGroupUserAccessAdmins struct {
	UserAccessAdmins []ResourceGroupUserAccessAdmin `json:"userAccessAdmins"`
	ResourceGroupId  string                         `json:"resourceGroupId"`
}

func (s ResourceGroupUserAccessAdmins) MarshalJSON() ([]byte, error) {
	type Alias ResourceGroupUserAccessAdmins
	a := Alias(s)
	a.ResourceGroupId = strings.ToUpper(a.ResourceGroupId)
	return json.Marshal(a)
}

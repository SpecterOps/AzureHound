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

type KeyVaultUserAccessAdmin struct {
	UserAccessAdmin azure.RoleAssignment `json:"userAccessAdmin"`
	KeyVaultId      string               `json:"keyVaultId"`
}

func (s KeyVaultUserAccessAdmin) MarshalJSON() ([]byte, error) {
	type Alias KeyVaultUserAccessAdmin
	a := Alias(s)
	a.KeyVaultId = strings.ToUpper(a.KeyVaultId)
	a.UserAccessAdmin = UpperRoleAssignment(a.UserAccessAdmin)
	return json.Marshal(a)
}

type KeyVaultUserAccessAdmins struct {
	UserAccessAdmins []KeyVaultUserAccessAdmin `json:"userAccessAdmins"`
	KeyVaultId       string                    `json:"keyVaultId"`
}

func (s KeyVaultUserAccessAdmins) MarshalJSON() ([]byte, error) {
	type Alias KeyVaultUserAccessAdmins
	a := Alias(s)
	a.KeyVaultId = strings.ToUpper(a.KeyVaultId)
	return json.Marshal(a)
}

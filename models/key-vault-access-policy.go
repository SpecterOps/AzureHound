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

type KeyVaultAccessPolicy struct {
	azure.AccessPolicyEntry
	KeyVaultId string `json:"keyVaultId"`
}

func (s KeyVaultAccessPolicy) MarshalJSON() ([]byte, error) {
	type Alias KeyVaultAccessPolicy
	a := Alias(s)
	a.ObjectId = strings.ToUpper(a.ObjectId)
	a.ApplicationId = strings.ToUpper(a.ApplicationId)
	a.TenantId = strings.ToUpper(a.TenantId)
	a.KeyVaultId = strings.ToUpper(a.KeyVaultId)
	return json.Marshal(a)
}

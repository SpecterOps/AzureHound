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

type AppRoleAssignment struct {
	azure.AppRoleAssignment
	AppId    string `json:"appId"`
	TenantId string `json:"tenantId"`
}

func (s AppRoleAssignment) MarshalJSON() ([]byte, error) {
	type Alias AppRoleAssignment
	a := Alias(s)
	a.ResourceId = strings.ToUpper(a.ResourceId)
	a.TenantId = strings.ToUpper(a.TenantId)

	// PrincipalId is a uuid.UUID and cannot hold an uppercased string, so emit
	// it through a map override alongside the aliased fields.
	raw, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	var output map[string]json.RawMessage
	if err := json.Unmarshal(raw, &output); err != nil {
		return nil, err
	}
	if _, ok := output["principalId"]; ok {
		pid := strings.ToUpper(s.PrincipalId.String())
		output["principalId"], _ = json.Marshal(pid)
	}
	return json.Marshal(output)
}

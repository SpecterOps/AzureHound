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

type RoleAssignments struct {
	RoleAssignments  []azure.UnifiedRoleAssignment `json:"roleAssignments"`
	RoleDefinitionId string                        `json:"roleDefinitionId"`
	TenantId         string                        `json:"tenantId"`
}

func (s RoleAssignments) MarshalJSON() ([]byte, error) {
	type Alias RoleAssignments
	a := Alias(s)
	a.RoleDefinitionId = strings.ToUpper(a.RoleDefinitionId)
	a.TenantId = strings.ToUpper(a.TenantId)

	if s.RoleAssignments != nil {
		assignments := make([]azure.UnifiedRoleAssignment, len(s.RoleAssignments))
		for i, assignment := range s.RoleAssignments {
			assignment.PrincipalId = strings.ToUpper(assignment.PrincipalId)
			assignment.DirectoryScopeId = strings.ToUpper(assignment.DirectoryScopeId)
			if len(assignment.Principal) > 0 {
				if principal, err := OmitEmptyUpper(assignment.Principal, "id"); err != nil {
					return nil, err
				} else {
					assignment.Principal = principal
				}
			}
			assignments[i] = assignment
		}
		a.RoleAssignments = assignments
	}
	return json.Marshal(a)
}

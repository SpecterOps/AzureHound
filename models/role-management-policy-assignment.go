// Copyright (C) 2025 Specter Ops, Inc.
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

type RoleManagementPolicyAssignment struct {
	azure.UnifiedRoleManagementPolicyAssignment

	Id                                                string   `json:"id,omitempty"`
	RoleDefinitionId                                  string   `json:"roleDefinitionId,omitempty"`
	EndUserAssignmentRequiresApproval                 bool     `json:"endUserAssignmentRequiresApproval,omitempty"`
	EndUserAssignmentRequiresCAPAuthenticationContext bool     `json:"endUserAssignmentRequiresCAPAuthenticationContext,omitempty"`
	EndUserAssignmentUserApprovers                    []string `json:"endUserAssignmentUserApprovers,omitempty"`
	EndUserAssignmentGroupApprovers                   []string `json:"endUserAssignmentGroupApprovers,omitempty"`
	EndUserAssignmentRequiresMFA                      bool     `json:"endUserAssignmentRequiresMFA,omitempty"`
	EndUserAssignmentRequiresJustification            bool     `json:"endUserAssignmentRequiresJustification,omitempty"`
	EndUserAssignmentRequiresTicketInformation        bool     `json:"endUserAssignmentRequiresTicketInformation,omitempty"`
	TenantId                                          string   `json:"tenantId,omitempty"`
}

func (s RoleManagementPolicyAssignment) MarshalJSON() ([]byte, error) {
	type Alias RoleManagementPolicyAssignment
	a := Alias(s)
	a.Id = strings.ToUpper(a.Id)
	a.RoleDefinitionId = strings.ToUpper(a.RoleDefinitionId)
	a.TenantId = strings.ToUpper(a.TenantId)
	a.EndUserAssignmentUserApprovers = upperStrings(a.EndUserAssignmentUserApprovers)
	a.EndUserAssignmentGroupApprovers = upperStrings(a.EndUserAssignmentGroupApprovers)
	return json.Marshal(a)
}

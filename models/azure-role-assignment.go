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

type AzureRoleAssignment struct {
	Assignee         azure.RoleAssignment `json:"assignee"`
	ObjectId         string               `json:"objectId"`
	RoleDefinitionId string               `json:"roleDefinitionId"`
}

// MarshalJSON uppercases the ObjectId and the Assignee endpoint identifiers so
// the raw (use_raw_object_id) ingest path matches the normalized node ObjectIDs.
// RoleDefinitionId is left untouched because ingest matches it against lowercase
// role-definition constants. The input is not mutated.
func (s AzureRoleAssignment) MarshalJSON() ([]byte, error) {
	type Alias AzureRoleAssignment
	a := Alias(s)
	a.ObjectId = strings.ToUpper(a.ObjectId)
	a.Assignee = UpperRoleAssignment(a.Assignee)
	return json.Marshal(a)
}

type AzureRoleAssignments struct {
	RoleAssignments []AzureRoleAssignment `json:"assignees"`
	ObjectId        string                `json:"objectId"`
}

// MarshalJSON uppercases the top-level ObjectId (the resource id BloodHound
// ingest reads as the RBAC edge target endpoint for the resource-scoped
// role-assignment convertors) so the raw (use_raw_object_id) ingest path matches
// the normalized resource node ObjectIDs. The nested assignees uppercase their
// own identifiers via AzureRoleAssignment.MarshalJSON. The input is not mutated.
func (s AzureRoleAssignments) MarshalJSON() ([]byte, error) {
	type Alias AzureRoleAssignments
	a := Alias(s)
	a.ObjectId = strings.ToUpper(a.ObjectId)
	return json.Marshal(a)
}

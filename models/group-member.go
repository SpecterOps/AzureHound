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
)

type GroupMember struct {
	Member  json.RawMessage `json:"member"`
	GroupId string          `json:"groupId"`
}

// MarshalJSON uppercases GroupId and the embedded member.id for raw
// (use_raw_object_id) ingest. An empty or nil member is emitted as null to
// avoid unmarshaling it. Non-mutating.
func (s GroupMember) MarshalJSON() ([]byte, error) {
	type Alias GroupMember
	a := Alias(s)
	a.GroupId = strings.ToUpper(a.GroupId)

	if len(a.Member) > 0 {
		member, err := OmitEmptyUpper(a.Member, "id")
		if err != nil {
			return nil, err
		}
		a.Member = member
	} else {
		a.Member = nil
	}
	return json.Marshal(a)
}

type GroupMembers struct {
	Members []GroupMember `json:"members"`
	GroupId string        `json:"groupId"`
}

func (s GroupMembers) MarshalJSON() ([]byte, error) {
	type Alias GroupMembers
	a := Alias(s)
	a.GroupId = strings.ToUpper(a.GroupId)
	return json.Marshal(a)
}

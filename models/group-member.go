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

func (s GroupMember) MarshalJSON() ([]byte, error) {
	output := make(map[string]any)
	output["groupId"] = strings.ToUpper(s.GroupId)

	if member, err := OmitEmptyUpper(s.Member, "id"); err != nil {
		return nil, err
	} else {
		output["member"] = member
		return json.Marshal(output)
	}
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

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

type AppMember struct {
	json.RawMessage
	AppId string `json:"appId"`
}

func (s AppMember) MarshalJSON() ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal(s.RawMessage, &data); err != nil {
		return nil, err
	} else {
		StripEmptyEntries(data)
		if id, ok := data["id"].(string); ok {
			data["id"] = strings.ToUpper(id)
		}
		data["appId"] = strings.ToUpper(s.AppId)
		return json.Marshal(data)
	}
}

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
)

type AppFIC struct {
	FIC   json.RawMessage `json:"fic"`
	AppId string          `json:"appId"`
}

func (s *AppFIC) MarshalJSON() ([]byte, error) {
	output := make(map[string]any)
	output["appId"] = s.AppId

	if fic, err := OmitEmpty(s.FIC); err != nil {
		return nil, err
	} else {
		output["fic"] = fic
		return json.Marshal(output)
	}
}

type AppFICs struct {
	FICs  []AppFIC `json:"fics"`
	AppId string   `json:"appId"`
}

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

type AppFIC struct {
	FIC   json.RawMessage `json:"fic"`
	AppId string          `json:"appId"`
}

// MarshalJSON uppercases AppId and the embedded fic.id for raw
// (use_raw_object_id) ingest; display-only fic fields are untouched. An empty
// or nil fic is emitted as null to avoid unmarshaling it. Non-mutating.
func (s *AppFIC) MarshalJSON() ([]byte, error) {
	type Alias AppFIC
	a := Alias(*s)
	a.AppId = strings.ToUpper(a.AppId)

	if len(a.FIC) > 0 {
		fic, err := OmitEmptyUpper(a.FIC, "id")
		if err != nil {
			return nil, err
		}
		a.FIC = fic
	} else {
		a.FIC = nil
	}
	return json.Marshal(a)
}

type AppFICs struct {
	FICs       []AppFIC `json:"fics"`
	AppId      string   `json:"appId"`
	TenantId   string   `json:"tenantId"`
	TenantName string   `json:"tenantName"`
}

type FICData struct {
	Audiences   []string `json:"audiences"`
	ID          string   `json:"id"`
	Issuer      string   `json:"issuer"`
	Name        string   `json:"name"`
	Subject     string   `json:"subject"`
	Description string   `json:"description"`
}

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

// MarshalJSON uppercases the AppId and the embedded fic.id so the raw
// (use_raw_object_id) ingest path matches the normalized node ObjectIDs.
// BloodHound ingest reads fic.id as the AZFederatedIdentityCredential node
// ObjectID and AZAuthenticatesTo source, and AppId as the AZApp endpoint. The
// remaining fic fields (issuer, subject, name, audiences) are display-only and
// are left untouched. The input is not mutated.
func (s *AppFIC) MarshalJSON() ([]byte, error) {
	output := make(map[string]any)
	output["appId"] = strings.ToUpper(s.AppId)

	if s.FIC == nil {
		return nil, nil
	}

	if fic, err := OmitEmptyUpper(s.FIC, "id"); err != nil {
		return nil, err
	} else {
		output["fic"] = fic
		return json.Marshal(output)
	}
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

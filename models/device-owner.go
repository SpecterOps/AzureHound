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

type DeviceOwner struct {
	Owner    json.RawMessage `json:"owner"`
	DeviceId string          `json:"deviceId"`
}

// MarshalJSON uppercases the DeviceId and the embedded owner.id so the raw
// (use_raw_object_id) ingest path matches the normalized node ObjectIDs. When
// owner is empty or nil it is emitted as null rather than passed to
// OmitEmptyUpper, which would fail to unmarshal an empty raw message. The input
// is not mutated.
func (s DeviceOwner) MarshalJSON() ([]byte, error) {
	type Alias DeviceOwner
	a := Alias(s)
	a.DeviceId = strings.ToUpper(a.DeviceId)

	if len(a.Owner) > 0 {
		owner, err := OmitEmptyUpper(a.Owner, "id")
		if err != nil {
			return nil, err
		}
		a.Owner = owner
	} else {
		a.Owner = nil
	}
	return json.Marshal(a)
}

type DeviceOwners struct {
	Owners   []DeviceOwner `json:"owners"`
	DeviceId string        `json:"deviceId"`
}

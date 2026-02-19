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

package azure

import (
	"encoding/json"
	"testing"
)

func TestUserUnmarshal_PopulatesLastSuccessfulSignInDateTimeFromSignInActivity(t *testing.T) {
	payload := []byte(`{
		"id":"3fb2a5fc-3a42-4c11-8200-85302657dc1a",
		"displayName":"test-user",
		"signInActivity":{
			"lastSignInDateTime":"2025-01-27T22:20:22Z",
			"lastSignInRequestId":"af4c2c83-9463-434d-a8e5-fbce099b2600",
			"lastNonInteractiveSignInDateTime":null,
			"lastNonInteractiveSignInRequestId":null,
			"lastSuccessfulSignInDateTime":"2025-01-27T22:20:22Z",
			"lastSuccessfulSignInRequestId":"af4c2c83-9463-434d-a8e5-fbce099b2600"
		}
	}`)

	var u User
	if err := json.Unmarshal(payload, &u); err != nil {
		t.Fatalf("unexpected unmarshal error: %v", err)
	}

	if u.SignInActivity.LastSuccessfulSignInDateTime != "2025-01-27T22:20:22Z" {
		t.Fatalf("expected LastSuccessfulSignInDateTime to be populated, got %q", u.SignInActivity.LastSuccessfulSignInDateTime)
	}
}

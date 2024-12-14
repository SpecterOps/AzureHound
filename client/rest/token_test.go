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

package rest

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTokenUnmarshall(t *testing.T) {
	
	// example response from oath2 auth code grant request
	expiresInIsInt := `{
		"access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIsIng1dCI6Ik5HVEZ2ZEstZnl0aEV1Q...",
		"token_type": "Bearer",
		"expires_in": 3599,
		"scope": "https%3A%2F%2Fgraph.microsoft.com%2Fmail.read",
		"refresh_token": "AwABAAAAvPM1KaPlrEqdFSBzjqfTGAMxZGUTdM0t4B4...",
		"id_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJub25lIn0.eyJhdWQiOiIyZDRkMTFhMi1mODE0LTQ2YTctOD..."
	}`
	
	// example response from managed identity auth request
	expiresInIsString := `{
  		"access_token": "eyJ0eXAi...",
  		"refresh_token": "",
 		 "expires_in": "3599",
 		 "expires_on": "1506484173",
  		"not_before": "1506480273",
  		"resource": "https://management.azure.com/",
  		"token_type": "Bearer"
	}`
	var t1 = Token{}
	var reader = strings.NewReader(expiresInIsInt)
	json.NewDecoder(reader).Decode(&t1)
	
	if t1.expiresIn != 3599 {
		t.Errorf("expected 3599 got %d", t1.expiresIn)
	}
	
	var t2 = Token{}
	reader = strings.NewReader(expiresInIsString)
	json.NewDecoder(reader).Decode(&t2)
	
	if t2.expiresIn != 3599 {
		t.Errorf("expected 3599 got %d", t2.expiresIn)		
	}
	
}



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
	"fmt"
	"time"
)

type Token struct {
	accessToken  string
	expiresIn    int
	extExpiresIn int
	expires      time.Time
}

func (s Token) IsExpired() bool {
	return time.Now().After(s.expires.Add(-10 * time.Second))
}

func (s Token) String() string {
	return fmt.Sprintf("Bearer %s", s.accessToken)
}

func (s *Token) UnmarshalJSON(data []byte) error {
	var res struct {
		AccessToken  string `json:"access_token"`   			// The token to use in calls to Microsoft Graph API
		ExpiresIn    IntOrStringInt   `json:"expires_in"`     	// How long the access token is valid in seconds (sometime a string of digits, sometimes an int)
		ExtExpiresIn IntOrStringInt    `json:"ext_expires_in"` 	// How long the access token is valid in seconds (sometime a string of digits, sometimes an int)
		TokenType    string `json:"token_type"`     			// Indicates the token type value. The only type currently supported by Azure AD is `bearer`
	}

	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}
	s.accessToken = res.AccessToken
	s.expiresIn = int(res.ExpiresIn)
	s.extExpiresIn = int(res.ExtExpiresIn)
	s.expires = time.Now().Add(time.Duration(int(res.ExpiresIn)) * time.Second)
		
	return nil
}

// Copyright (C) 2026 Specter Ops, Inc.
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
	"context"
	"net/url"
	"testing"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/constants"
)

func TestUserAgent_DefaultsWhenUnset(t *testing.T) {
	config.UserAgent.Set("")
	t.Cleanup(func() { config.UserAgent.Set("") })

	if got, want := UserAgent(), constants.UserAgent(); got != want {
		t.Fatalf("UserAgent() = %q, want default %q", got, want)
	}
}

func TestUserAgent_HonorsConfigValue(t *testing.T) {
	const custom = "my-custom-agent/1.2.3"
	config.UserAgent.Set(custom)
	t.Cleanup(func() { config.UserAgent.Set("") })

	if got := UserAgent(); got != custom {
		t.Fatalf("UserAgent() = %q, want %q", got, custom)
	}
}

func TestNewRequest_AppliesCustomUserAgent(t *testing.T) {
	const custom = "my-custom-agent/1.2.3"
	config.UserAgent.Set(custom)
	t.Cleanup(func() { config.UserAgent.Set("") })

	endpoint, _ := url.Parse("http://example.com/")
	req, err := NewRequest(context.Background(), "GET", endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if got := req.Header.Get("User-Agent"); got != custom {
		t.Fatalf("User-Agent header = %q, want %q", got, custom)
	}
}

func TestNewRequest_DefaultUserAgentWhenConfigEmpty(t *testing.T) {
	config.UserAgent.Set("")
	t.Cleanup(func() { config.UserAgent.Set("") })

	endpoint, _ := url.Parse("http://example.com/")
	req, err := NewRequest(context.Background(), "GET", endpoint, nil, nil, nil)
	if err != nil {
		t.Fatalf("NewRequest error: %v", err)
	}
	if got, want := req.Header.Get("User-Agent"), constants.UserAgent(); got != want {
		t.Fatalf("User-Agent header = %q, want default %q", got, want)
	}
}

// Copyright (C) 2024 Specter Ops, Inc.
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
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/bloodhoundad/azurehound/v2/client/config"
)

func TestClosedConnection(t *testing.T) {
	var testServer *httptest.Server
	attempt := 0
	var mockHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		attempt++
		testServer.CloseClientConnections()
	}

	testServer = httptest.NewServer(mockHandler)
	defer testServer.Close()

	defaultConfig := config.Config{
		Username:  "azurehound",
		Password:  "we_collect",
		Authority: testServer.URL,
	}

	if client, err := NewRestClient(testServer.URL, defaultConfig); err != nil {
		t.Fatalf("error initializing rest client %v", err)
	} else {
		requestCompleted := false

		// make request in separate goroutine so its not blocking after we validated the retry
		go func() {
			req := httptest.NewRequest("GET", testServer.URL, nil)
			client.AddAuthenticationToRequest(req) // Authenticate() because it uses the internal client.send method.
			// CloseClientConnections should block the request from completing, however if it completes then the test fails.
			requestCompleted = true
		}()

		// block until attempt is > 2 or request succeeds
		for attempt <= 2 {
			if attempt > 1 || requestCompleted {
				break
			}
		}

		if requestCompleted {
			t.Fatalf("expected an attempted retry but the request completed")
		}
	}
}

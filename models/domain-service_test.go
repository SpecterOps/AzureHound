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
	"testing"
)

func TestDomainServiceMarshalJSON(t *testing.T) {
	var domainService DomainService
	raw := []byte(`{
		"id":"/subscriptions/sub/resourceGroups/rg/providers/Microsoft.AAD/domainServices/example.com",
		"name":"example.com",
		"properties":{
			"tenantId":"aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee",
			"domainName":"example.com",
			"syncApplicationId":"11111111-2222-3333-4444-555555555555",
			"ldapsSettings":{
				"ldaps":"Enabled",
				"externalAccess":"Disabled",
				"publicCertificate":"certificate-data",
				"certificateNotAfter":"2030-01-01T00:00:00Z",
				"certificateThumbprint":"thumbprint"
			}
		}
	}`)
	if err := json.Unmarshal(raw, &domainService); err != nil {
		t.Fatalf("failed to unmarshal domain service: %v", err)
	}

	domainService.SubscriptionID = "/subscriptions/sub"
	domainService.ResourceGroupID = "/subscriptions/sub/resourceGroups/rg"
	encoded, err := json.Marshal(domainService)
	if err != nil {
		t.Fatalf("failed to marshal domain service: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(encoded, &result); err != nil {
		t.Fatalf("failed to inspect marshaled domain service: %v", err)
	}
	properties := result["properties"].(map[string]any)
	ldapsSettings := properties["ldapsSettings"].(map[string]any)

	if result["id"] != "/SUBSCRIPTIONS/SUB/RESOURCEGROUPS/RG/PROVIDERS/MICROSOFT.AAD/DOMAINSERVICES/EXAMPLE.COM" {
		t.Errorf("expected uppercased resource id, got %v", result["id"])
	}
	if properties["tenantId"] != "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE" {
		t.Errorf("expected uppercased tenant id, got %v", properties["tenantId"])
	}
	if properties["syncApplicationId"] != "11111111-2222-3333-4444-555555555555" {
		t.Errorf("expected sync application id to be retained, got %v", properties["syncApplicationId"])
	}
	for _, excluded := range []string{"publicCertificate", "certificateNotAfter", "certificateThumbprint"} {
		if _, found := ldapsSettings[excluded]; found {
			t.Errorf("did not expect %s in marshaled LDAPS settings", excluded)
		}
	}
}

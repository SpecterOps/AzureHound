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

package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bloodhoundad/azurehound/v2/constants"
)

func TestParseProductVersion(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{
			name:    "valid version",
			args:    []string{"cmd", "v1.2.3"},
			want:    "v1.2.3",
			wantErr: false,
		},
		{
			name:    "valid version with build metadata",
			args:    []string{"cmd", "v0.0.0-rolling+5f8807a4107f0b80debaf79b2d245bfa7078a54b"},
			want:    "v0.0.0-rolling+5f8807a4107f0b80debaf79b2d245bfa7078a54b",
			wantErr: false,
		},
		{
			name:    "no version provided",
			args:    []string{"cmd"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty version string",
			args:    []string{"cmd", ""},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original os.Args and restore after test
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			os.Args = tt.args

			got, err := parseProductVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("parseProductVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseProductVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildWinresConfig(t *testing.T) {
	productVersion := "v2.9.0-test"
	config := buildWinresConfig(productVersion)

	// Verify icon resource
	rtGroupIcon, ok := config["RT_GROUP_ICON"].(map[string]interface{})
	if !ok {
		t.Fatal("RT_GROUP_ICON not found or wrong type")
	}

	app, ok := rtGroupIcon["APP"].(map[string]interface{})
	if !ok {
		t.Fatal("APP icon not found or wrong type")
	}

	iconPath, ok := app["0000"].(string)
	if !ok {
		t.Fatal("icon path not found or wrong type")
	}

	if iconPath != iconFile {
		t.Errorf("icon path = %v, want %v", iconPath, iconFile)
	}

	// Verify top-level structure
	rtVersion, ok := config["RT_VERSION"].(map[string]interface{})
	if !ok {
		t.Fatal("RT_VERSION not found or wrong type")
	}

	// Navigate to the version info
	level1, ok := rtVersion["#1"].(map[string]interface{})
	if !ok {
		t.Fatal("#1 not found or wrong type")
	}

	level2, ok := level1["0000"].(map[string]interface{})
	if !ok {
		t.Fatal("0000 not found or wrong type")
	}

	// Test fixed section
	fixed, ok := level2["fixed"].(map[string]interface{})
	if !ok {
		t.Fatal("fixed section not found")
	}

	if fixed["file_version"] != fileVersion {
		t.Errorf("file_version = %v, want %v", fixed["file_version"], fileVersion)
	}

	if fixed["product_version"] != fileVersion {
		t.Errorf("product_version = %v, want %v", fixed["product_version"], fileVersion)
	}

	// Test info section
	info, ok := level2["info"].(map[string]interface{})
	if !ok {
		t.Fatal("info section not found")
	}

	langInfo, ok := info[langCodeUSEn].(map[string]string)
	if !ok {
		t.Fatal("language info not found")
	}

	// Verify all required fields
	requiredFields := map[string]string{
		"FileDescription":  constants.Description,
		"ProductName":      constants.DisplayName,
		"CompanyName":      constants.AuthorRef,
		"ProductVersion":   productVersion,
		"FileVersion":      fileVersion,
		"OriginalFilename": "azurehound.exe",
	}

	for field, expected := range requiredFields {
		if got := langInfo[field]; got != expected {
			t.Errorf("%s = %v, want %v", field, got, expected)
		}
	}

	// Verify copyright contains current year and company
	copyright := langInfo["LegalCopyright"]
	currentYear := time.Now().Year()
	expectedCopyright := fmt.Sprintf("Copyright (C) %d %s", currentYear, constants.Company)
	if copyright != expectedCopyright {
		t.Errorf("LegalCopyright = %v, want %v", copyright, expectedCopyright)
	}
}

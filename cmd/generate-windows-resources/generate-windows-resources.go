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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/bloodhoundad/azurehound/v2/constants"
)

const (
	winresDir      = "winres"
	winresJSONFile = "winres.json"
	iconFile       = "favicon.ico"
	langCodeUSEn   = "0409"    // US English
	fileVersion    = "0.0.0.0" // Windows PE file version; we will update 'productVersion' field instead of this one
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	productVersion, err := parseProductVersion()
	if err != nil {
		return err
	}

	config := buildWinresConfig(productVersion)

	if err := writeWinresConfig(config); err != nil {
		return fmt.Errorf("failed to write winres config: %w", err)
	}

	if err := runWinres(); err != nil {
		return fmt.Errorf("failed to generate windows resources: %w", err)
	}

	fmt.Printf("✓ Windows resources generated successfully!\n")
	fmt.Printf("  Product Version: %s\n", productVersion)
	return nil
}

func parseProductVersion() (string, error) {
	if len(os.Args) < 2 {
		return "", fmt.Errorf("usage: %s <product-version>", filepath.Base(os.Args[0]))
	}
	version := os.Args[1]
	if version == "" {
		return "", fmt.Errorf("product version cannot be empty")
	}
	return version, nil
}

func buildWinresConfig(productVersion string) map[string]interface{} {
	return map[string]interface{}{
		// Icon resource
		"RT_GROUP_ICON": map[string]interface{}{
			"APP": map[string]interface{}{
				"0000": iconFile,
			},
		},
		// Version information
		"RT_VERSION": map[string]interface{}{
			"#1": map[string]interface{}{
				"0000": map[string]interface{}{
					"fixed": map[string]interface{}{
						"file_version":    fileVersion,
						"product_version": fileVersion,
					},
					"info": map[string]interface{}{
						langCodeUSEn: map[string]string{
							"FileDescription":  constants.Description,
							"ProductName":      constants.DisplayName,
							"CompanyName":      constants.AuthorRef,
							"LegalCopyright":   fmt.Sprintf("Copyright (C) %d %s", time.Now().Year(), constants.Company),
							"ProductVersion":   productVersion,
							"FileVersion":      fileVersion,
							"OriginalFilename": "azurehound.exe",
						},
					},
				},
			},
		},
	}
}

func writeWinresConfig(config map[string]interface{}) error {
	if err := os.MkdirAll(winresDir, 0755); err != nil {
		return fmt.Errorf("failed to create winres directory: %w", err)
	}

	configPath := filepath.Join(winresDir, winresJSONFile)
	f, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", configPath, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(config); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func runWinres() error {
	cmd := exec.Command("go", "tool", "go-winres", "make")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-winres command failed: %w", err)
	}

	return nil
}

// Copyright (C) 2026 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

package cmd

import "testing"

func TestValidateMinInt(t *testing.T) {
	tests := []struct {
		name    string
		minimum int
		input   string
		wantErr bool
	}{
		{name: "minimum", minimum: 1, input: "1"},
		{name: "above minimum", minimum: 1, input: "100"},
		{name: "zero allowed", minimum: 0, input: "0"},
		{name: "below minimum", minimum: 1, input: "0", wantErr: true},
		{name: "negative", minimum: 0, input: "-1", wantErr: true},
		{name: "not an integer", minimum: 0, input: "one", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateMinInt(test.minimum)(test.input)
			if test.wantErr && err == nil {
				t.Fatal("expected validation error")
			}
			if !test.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

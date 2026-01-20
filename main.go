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

//go:build !windows
// +build !windows

//go:debug x509negativeserial=1

package main

import (
	"fmt"
	"os"

	"github.com/bloodhoundad/azurehound/v2/cmd"
	"github.com/bloodhoundad/azurehound/v2/constants"
)

func main() {
	fmt.Fprintf(os.Stderr, "%s %s\n%s\n\n", constants.DisplayName, constants.Version, constants.AuthorRef)
	cmd.Execute()
}

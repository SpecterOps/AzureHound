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

package main

import (
	"fmt"

	"golang.org/x/sys/windows/svc"

	"github.com/certmichelin/azurehound/v2/cmd"
	"github.com/certmichelin/azurehound/v2/constants"
)

func main() {
	fmt.Printf("%s %s\n%s\n\n", constants.DisplayName, constants.Version, constants.AuthorRef)

	if isWinSvc, err := svc.IsWindowsService(); err != nil {
		panic(err)
	} else if isWinSvc {
		if err := cmd.StartWindowsService(); err != nil {
			panic(err)
		}
	} else {
		cmd.Execute()
	}
}

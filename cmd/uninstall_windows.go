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

package cmd

import (
	"fmt"

	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/spf13/cobra"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:               "uninstall",
	Short:             "Removes AzureHound as a system service",
	Run:               uninstallCmdImpl,
	PersistentPreRunE: persistentPreRunE,
	SilenceUsage:      true,
}

func uninstallCmdImpl(cmd *cobra.Command, args []string) {
	if err := uninstallService(constants.Name); err != nil {
		exit(fmt.Errorf("failed to uninstall service: %w", err))
	}
}

func uninstallService(name string) error {
	if wsm, err := mgr.Connect(); err != nil {
		return err
	} else {
		defer wsm.Disconnect()

		if service, err := wsm.OpenService(name); err != nil {
			return err
		} else {
			defer service.Close()

			if err := service.Delete(); err != nil {
				return err
			} else if err := eventlog.Remove(name); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

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
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/spf13/cobra"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func init() {
	rootCmd.AddCommand(installCmd)
}

var installCmd = &cobra.Command{
	Use:               "install",
	Short:             "Installs AzureHound as a system service for BloodHound Enterprise",
	Run:               installCmdImpl,
	PersistentPreRunE: persistentPreRunE,
	SilenceUsage:      true,
}

func installCmdImpl(cmd *cobra.Command, args []string) {
	var (
		config = mgr.Config{
			DisplayName:      constants.DisplayName,
			Description:      constants.Description,
			StartType:        mgr.StartAutomatic,
			DelayedAutoStart: true,
		}
		recoveryActions = []mgr.RecoveryAction{
			{Type: mgr.ServiceRestart, Delay: 5 * time.Second},
			{Type: mgr.ServiceRestart, Delay: 30 * time.Second},
			{Type: mgr.ServiceRestart, Delay: 60 * time.Second},
		}
	)

	if err := configureService(); err != nil {
		exit(fmt.Errorf("failed to configure service: %w", err))
	} else if err := installService(constants.DisplayName, config, recoveryActions); err != nil {
		exit(fmt.Errorf("failed to install service: %w", err))
	}
}

func configureService() error {
	var (
		configDir  = config.SystemConfigDirs()[0]
		sysConfig  = filepath.Join(configDir, "config.json")
		userConfig = config.ConfigFile.Value().(string)
	)

	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return err
	}

	// Confirm use of existing service config
	if shouldUseConfig(sysConfig) {
		return nil
	}

	// Confirm use of existing user config
	if shouldUseConfig(userConfig) {
		return copyFile(userConfig, sysConfig)
	}

	config.ConfigFile.Set(sysConfig)
	return configure()
}

func shouldUseConfig(config string) bool {
	if _, err := os.Stat(config); err != nil {
		return false
	} else {
		fmt.Fprintf(os.Stderr, "Detected configuration at %s.\n", config)
		return confirm("Use these settings to configure the service", true)
	}
}

func copyFile(src, dest string) error {
	if srcFile, err := os.Open(src); err != nil {
		return err
	} else if destFile, err := os.Create(dest); err != nil {
		return err
	} else {
		defer srcFile.Close()
		defer destFile.Close()
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return err
		}
	}
	return nil
}

func installService(name string, config mgr.Config, recoveryActions []mgr.RecoveryAction, args ...string) error {
	if exe, err := getExePath(); err != nil {
		return err
	} else if wsm, err := mgr.Connect(); err != nil {
		return err
	} else {
		defer wsm.Disconnect()

		if err := createService(wsm, name, exe, config, recoveryActions, args...); err != nil {
			return err
		} else {
			return nil
		}
	}
}

func createService(wsm *mgr.Mgr, name string, exe string, config mgr.Config, recoveryActions []mgr.RecoveryAction, args ...string) error {
	if service, err := wsm.OpenService(name); err == nil {
		service.Close()
		return fmt.Errorf("service %s already exists", name)
	} else if service, err := wsm.CreateService(name, exe, config, args...); err != nil {
		return err
	} else {
		defer service.Close()

		if err := eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info); err != nil {
			service.Delete()
			return fmt.Errorf("failed to add %s to event log: %w", name, err)
		}

		if recoveryActions != nil {
			if err := service.SetRecoveryActions(recoveryActions, 60); err != nil {
				service.Delete()
				return fmt.Errorf("failed to set recovery actions: %w", err)
			}
		}

		return nil
	}
}

func getExePath() (string, error) {
	exe := os.Args[0]
	if exePath, err := filepath.Abs(exe); err != nil {
		return "", err
	} else if path, info, err := stat(exePath); err != nil {
		return "", err
	} else if info.IsDir() {
		return "", fmt.Errorf("%s is a directory", path)
	} else {
		return path, nil
	}
}

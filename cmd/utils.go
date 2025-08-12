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
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"runtime/pprof"

	"github.com/bloodhoundad/azurehound/v2/client/rest"
	"github.com/spf13/cobra"
	"golang.org/x/net/proxy"

	"github.com/bloodhoundad/azurehound/v2/client"
	client_config "github.com/bloodhoundad/azurehound/v2/client/config"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/logger"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/bloodhoundad/azurehound/v2/sinks"
)

func init() {
	proxy.RegisterDialerType("http", rest.NewProxyDialer)
	proxy.RegisterDialerType("https", rest.NewProxyDialer)
}

func exit(err error) {
	log.Error(err, "encountered unrecoverable error")
	log.GetSink()
	os.Exit(1)
}

func persistentPreRunE(cmd *cobra.Command, args []string) error {
	// need to set config flag value explicitly
	if cmd != nil {
		if configFlag := cmd.Flag(config.ConfigFile.Name).Value.String(); configFlag != "" {
			config.ConfigFile.Set(configFlag)
		}
	}

	config.LoadValues(cmd, config.Options())
	config.SetAzureDefaults()

	if logr, err := logger.GetLogger(); err != nil {
		return err
	} else {
		log = *logr
		config.CheckCollectionConfigSanity(log)

		if config.ConfigFileUsed() != "" {
			log.V(1).Info(fmt.Sprintf("Config File: %v", config.ConfigFileUsed()))
		}

		if config.LogFile.Value() != "" {
			log.V(1).Info(fmt.Sprintf("Log File: %v", config.LogFile.Value()))
		}

		return nil
	}
}

func gracefulShutdown(stop context.CancelFunc) {
	stop()
	fmt.Fprintln(os.Stderr, "\nshutting down gracefully, press ctrl+c again to force")
	if profile := pprof.Lookup(config.Pprof.Value().(string)); profile != nil {
		profile.WriteTo(os.Stderr, 1)
	}
}

func testConnections() error {
	if _, err := rest.Dial(log, config.AzAuthUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzAuthUrl.Value(), err)
	} else if _, err := rest.Dial(log, config.AzGraphUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzGraphUrl.Value(), err)
	} else if _, err := rest.Dial(log, config.AzMgmtUrl.Value().(string)); err != nil {
		return fmt.Errorf("unable to connect to %s: %w", config.AzMgmtUrl.Value(), err)
	} else {
		return nil
	}
}

func newAzureClient() (client.AzureClient, error) {
	var (
		certFile   = config.AzCert.Value()
		keyFile    = config.AzKey.Value()
		clientCert string
		clientKey  string
	)

	if file, ok := certFile.(string); ok && file != "" {
		if content, err := os.ReadFile(certFile.(string)); err != nil {
			return nil, fmt.Errorf("unable to read provided certificate: %w", err)
		} else {
			clientCert = string(content)
		}
	}

	if file, ok := keyFile.(string); ok && file != "" {
		if content, err := os.ReadFile(keyFile.(string)); err != nil {
			return nil, fmt.Errorf("unable to read provided key file: %w", err)
		} else {
			clientKey = string(content)
		}
	}

	config := client_config.Config{
		ApplicationId:   config.AzAppId.Value().(string),
		Authority:       config.AzAuthUrl.Value().(string),
		ClientSecret:    config.AzSecret.Value().(string),
		ClientCert:      clientCert,
		ClientKey:       clientKey,
		ClientKeyPass:   config.AzKeyPass.Value().(string),
		Graph:           config.AzGraphUrl.Value().(string),
		JWT:             config.JWT.Value().(string),
		Management:      config.AzMgmtUrl.Value().(string),
		MgmtGroupId:     config.AzMgmtGroupId.Value().([]string),
		Password:        config.AzPassword.Value().(string),
		ProxyUrl:        config.Proxy.Value().(string),
		RefreshToken:    config.RefreshToken.Value().(string),
		Region:          config.AzRegion.Value().(string),
		SubscriptionId:  config.AzSubId.Value().([]string),
		Tenant:          config.AzTenant.Value().(string),
		Username:        config.AzUsername.Value().(string),
		ManagedIdentity: config.AzUseManagedIdentity.Value().(bool),
	}
	return client.NewClient(config)
}

func contains[T comparable](collection []T, value T) bool {
	for _, item := range collection {
		if item == value {
			return true
		}
	}
	return false
}

func unique(collection []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, item := range collection {
		if _, found := keys[item]; !found {
			keys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func stat(path string) (string, fs.FileInfo, error) {
	if info, err := os.Stat(path); err == nil {
		return path, info, nil
	} else {
		p := path + ".exe"
		info, err := os.Stat(p)
		return p, info, err
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

func setupLogger() {
	if logger, err := logger.GetLogger(); err != nil {
		panic(err)
	} else {
		log = *logger
	}
}

// deprecated: use azureWrapper instead
type AzureWrapper struct {
	Kind enums.Kind  `json:"kind"`
	Data interface{} `json:"data"`
}

type azureWrapper[T any] struct {
	Kind enums.Kind `json:"kind"`
	Data T          `json:"data"`
}

func NewAzureWrapper[T any](kind enums.Kind, data T) azureWrapper[T] {
	return azureWrapper[T]{
		Kind: kind,
		Data: data,
	}
}

func outputStream[T any](ctx context.Context, stream <-chan T) {
	formatted := pipeline.FormatJson(ctx.Done(), stream)
	if path := config.OutputFile.Value().(string); path != "" {
		if err := sinks.WriteToFile(ctx, path, formatted); err != nil {
			exit(fmt.Errorf("failed to write stream to file: %w", err))
		}
	} else {
		sinks.WriteToConsole(ctx, formatted)
	}
}

func kvRoleAssignmentFilter(roleId string) func(models.KeyVaultRoleAssignment) bool {
	return func(ra models.KeyVaultRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func vmRoleAssignmentFilter(roleId string) func(models.VirtualMachineRoleAssignment) bool {
	return func(ra models.VirtualMachineRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func rgRoleAssignmentFilter(roleId string) func(models.ResourceGroupRoleAssignment) bool {
	return func(ra models.ResourceGroupRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func mgmtGroupRoleAssignmentFilter(roleId string) func(models.ManagementGroupRoleAssignment) bool {
	return func(ra models.ManagementGroupRoleAssignment) bool {
		return path.Base(ra.RoleAssignment.Properties.RoleDefinitionId) == roleId
	}
}

func connectAndCreateClient() client.AzureClient {
	log.V(1).Info("testing connections")
	if err := testConnections(); err != nil {
		exit(fmt.Errorf("failed to test connections: %w", err))
	} else if azClient, err := newAzureClient(); err != nil {
		exit(fmt.Errorf("failed to create new Azure client: %w", err))
	} else {
		return azClient
	}

	panic("unexpectedly failed to create azClient without error")
}

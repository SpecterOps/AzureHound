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

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	config "github.com/bloodhoundad/azurehound/v2/config/internal"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/enums"
)

type Config = config.Config

var (
	homeDir, _ = os.UserHomeDir()

	// DefaultConfigFile is the path to the default configuration file.
	//
	// - $HOME/.config/azurehound/config.json (Unix/Darwin)
	// - %USERPROFILE%\.config\azurehound\config.json (Windows)
	DefaultConfigFile = filepath.Join(homeDir, ".config", "azurehound", "config.json")
)

func SystemConfigDirs() []string {
	prefixes := func() []string {
		switch runtime.GOOS {
		case "darwin":
			return []string{"/Library/Application Support"}
		case "linux":
			if xdgDirs := os.Getenv("XDG_CONFIG_DIRS"); xdgDirs != "" {
				return strings.Split(xdgDirs, ":")
			} else {
				return []string{"/etc/xdg"}
			}
		case "windows":
			return []string{os.Getenv("PROGRAMDATA")}
		default:
			panic("unsupported operating system")
		}
	}()

	configDirs := []string{}
	for _, dir := range prefixes {
		path := filepath.Join(dir, "azurehound")
		configDirs = append(configDirs, path)
	}
	return configDirs
}

const EnvPrefix string = "AZUREHOUND"

var AzRegions = []string{
	constants.China,
	constants.Cloud,
	constants.Germany,
	constants.USGovL4,
	constants.USGovL5,
}

var (
	// Global Configurations
	ConfigFile = Config{
		Name:       "config",
		Shorthand:  "c",
		Usage:      fmt.Sprintf("AzureHound configuration file (default: %s)", DefaultConfigFile),
		Persistent: true,
		Default:    DefaultConfigFile,
	}
	VerbosityLevel = Config{
		Name:       "verbosity",
		Shorthand:  "v",
		Usage:      fmt.Sprintf("AzureHound verbosity level (defaults to %d) [Min: %d, Max: %d]", 0, -1, 2),
		Persistent: true,
		Default:    0,
	}
	JsonLogs = Config{
		Name:       "json",
		Shorthand:  "",
		Usage:      "Output logs as json",
		Persistent: true,
		Default:    false,
	}
	JWT = Config{
		Name:       "jwt",
		Shorthand:  "j",
		Usage:      "Use an acquired JWT to authenticate into Azure",
		Persistent: true,
		Default:    "",
	}
	LogFile = Config{
		Name:       "log-file",
		Shorthand:  "",
		Usage:      "Output logs to this file",
		Persistent: true,
		Default:    "",
	}
	Proxy = Config{
		Name:       "proxy",
		Shorthand:  "",
		Usage:      "Sets the proxy URL for the AzureHound service",
		Persistent: true,
		Default:    "",
	}
	RefreshToken = Config{
		Name:       "refresh-token",
		Shorthand:  "r",
		Usage:      "Use an acquired refresh token to authenticate into Azure",
		Persistent: true,
		Default:    "",
	}
	Pprof = Config{
		Name:       "pprof",
		Usage:      "During graceful shutdown, prints the pprof profile with the provided name to stderr",
		Persistent: true,
		Default:    "",
	}

	// Azure Configurations
	AzAppId = Config{
		Name:       "app",
		Shorthand:  "a",
		Usage:      "The Application Id that the Azure app registration portal assigned when the app was registered.",
		Persistent: true,
		Default:    "",
	}
	AzSecret = Config{
		Name:       "secret",
		Shorthand:  "s",
		Usage:      "The Application Secret that was generated for the app in the app registration portal.",
		Persistent: true,
		Default:    "",
	}
	AzCert = Config{
		Name:       "cert",
		Shorthand:  "",
		Usage:      "The path to the certificate uploaded to the app registration portal.",
		Persistent: true,
		Default:    "",
	}
	AzKey = Config{
		Name:       "key",
		Shorthand:  "k",
		Usage:      "The path to the key file for a certificate uploaded to the app registration portal.",
		Persistent: true,
		Default:    "",
	}
	AzKeyPass = Config{
		Name:       "keypass",
		Shorthand:  "",
		Usage:      "The passphrase to use in conjuction with --key ${key file}.",
		Persistent: true,
		Default:    "",
	}
	AzRegion = Config{
		Name:       "region",
		Shorthand:  "",
		Usage:      fmt.Sprintf("The region of the Azure Cloud deployment (defaults to '%s') [%s]", constants.Cloud, strings.Join(AzRegions, ", ")),
		Persistent: true,
		Default:    constants.Cloud,
	}
	AzTenant = Config{
		Name:       "tenant",
		Shorthand:  "t",
		Usage:      "The directory tenant that you want to request permission from. This can be in GUID or friendly name format.",
		Required:   true,
		Persistent: true,
		Default:    "",
	}
	AzAuthUrl = Config{
		Name:       "auth",
		Shorthand:  "",
		Usage:      "The Azure ActiveDirectory Authority URL.",
		Persistent: true,
		Default:    "",
	}
	AzGraphUrl = Config{
		Name:       "graph",
		Shorthand:  "",
		Usage:      "The Microsoft Graph URL.",
		Persistent: true,
		Default:    "",
	}
	AzMgmtUrl = Config{
		Name:       "mgmt",
		Shorthand:  "",
		Usage:      "The URL of the Azure Resource Manager.",
		Persistent: true,
		Default:    "",
	}
	AzUsername = Config{
		Name:       "username",
		Shorthand:  "u",
		Usage:      "The user principal name for the Azure Portal",
		Persistent: true,
		Default:    "",
	}
	AzPassword = Config{
		Name:       "password",
		Shorthand:  "p",
		Usage:      "The user's password for the Azure Portal",
		Persistent: true,
		Default:    "",
	}
	AzSubId = Config{
		Name:       "subscriptionId",
		Shorthand:  "b",
		Usage:      "The subscription ID to use as a filter.",
		Persistent: true,
		Default:    []string{},
	}
	AzMgmtGroupId = Config{
		Name:       "mgmtGroupId",
		Shorthand:  "m",
		Usage:      "The management group ID to use as a filter.",
		Persistent: true,
		Default:    []string{},
	}

	AzUseManagedIdentity = Config{
		Name:       "managed-identity",
		Shorthand:  "",
		Usage:      "If true then authentication is done via managed identity (default false).",
		Persistent: true,
		Default:    bool(false),
	}

	// BHE Configurations
	BHEUrl = Config{
		Name:       "instance",
		Shorthand:  "i",
		Usage:      "The BloodHound Enterprise instance URL.",
		Persistent: true,
		Required:   true,
		Default:    "",
	}

	BHEToken = Config{
		Name:       "token",
		Shorthand:  "",
		Usage:      "The BloodHound Enterprise token.",
		Persistent: true,
		Required:   true,
		Default:    "",
	}

	BHETokenId = Config{
		Name:       "tokenId",
		Shorthand:  "",
		Usage:      "The BloodHound Enterprise token ID.",
		Persistent: true,
		Required:   true,
		Default:    "",
	}

	BHEMaxReqPerConn = Config{
		Name:      "maxReqsPerConn",
		Shorthand: "",
		Usage:     "The number of requests a single HTTP connection can make, when this limit is reached, a new HTTP connection will be established with the server",
		Required:  false,
		Default:   10_000,
		MinValue:  1,
		MaxValue:  10_000,
	}

	ColBatchSize = Config{
		Name:       "batchSize",
		Shorthand:  "",
		Usage:      "The number of resources to send in a single batch sent to the server.",
		Persistent: true,
		Required:   false,
		Default:    100,
		MinValue:   1,
		MaxValue:   256,
	}

	ColMaxConnsPerHost = Config{
		Name:       "maxConnsPerHost",
		Shorthand:  "",
		Usage:      "The maximum number of connections made during collection.",
		Persistent: true,
		Required:   false,
		Default:    20,
		MinValue:   1,
		MaxValue:   200,
	}

	ColMaxIdleConnsPerHost = Config{
		Name:       "maxIdleConnsPerHost",
		Shorthand:  "",
		Usage:      "The maximum number of idle connections allowed during collection.",
		Persistent: true,
		Required:   false,
		Default:    20,
		MinValue:   1,
		MaxValue:   200,
	}

	ColStreamCount = Config{
		Name:       "streamCount",
		Shorthand:  "",
		Usage:      "The number of threads to use when collecting various resources.",
		Persistent: true,
		Required:   false,
		Default:    25,
		MinValue:   1,
		MaxValue:   50,
	}

	// Command specific configurations
	KeyVaultAccessTypes = Config{
		Name:       "access-types",
		Shorthand:  "",
		Usage:      fmt.Sprintf("Filter key vault policies by one or more access type. [%s]\n\tNote: may be used multiple times or values may be provided as comma-separated list\n", strings.Join(enums.KeyVaultAccessPolicies(), ", ")),
		Persistent: true,
		Default:    []enums.KeyVaultAccessType{},
	}

	OutputFile = Config{
		Name:       "output",
		Shorthand:  "o",
		Usage:      "The path to the file in which to output data",
		Persistent: true,
		Default:    "",
	}

	GlobalConfig = []Config{
		ConfigFile,
		VerbosityLevel,
		JsonLogs,
		JWT,
		LogFile,
		Proxy,
		RefreshToken,
		Pprof,
	}

	AzureConfig = []Config{
		AzAppId,
		AzSecret,
		AzCert,
		AzKey,
		AzKeyPass,
		AzRegion,
		AzTenant,
		AzAuthUrl,
		AzGraphUrl,
		AzMgmtUrl,
		AzUsername,
		AzPassword,
		AzSubId,
		AzMgmtGroupId,
		AzUseManagedIdentity,
	}

	BloodHoundEnterpriseConfig = []Config{
		BHEUrl,
		BHETokenId,
		BHEToken,
		BHEMaxReqPerConn,
	}

	CollectionConfig = []Config{
		ColBatchSize,
		ColMaxConnsPerHost,
		ColMaxIdleConnsPerHost,
		ColStreamCount,
	}
)

func ConfigFileUsed() string {
	return config.ConfigFileUsed()
}

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
	"net/url"
	"os"

	client "github.com/bloodhoundad/azurehound/v2/client/config"
	config "github.com/bloodhoundad/azurehound/v2/config/internal"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/go-logr/logr"
)

var Init = config.Init
var LoadValues = config.LoadValues

func SetAzureDefaults() {
	if AzAuthUrl.Value() == "" {
		region := AzRegion.Value().(string)
		url := client.AuthorityUrl(region, constants.AzureCloud().ActiveDirectoryAuthority)
		AzAuthUrl.Set(url)
	}

	if AzGraphUrl.Value() == "" {
		region := AzRegion.Value().(string)
		url := client.GraphUrl(region, constants.AzureCloud().MicrosoftGraphUrl)
		AzGraphUrl.Set(url)
	}

	if AzMgmtUrl.Value() == "" {
		region := AzRegion.Value().(string)
		url := client.ResourceManagerUrl(region, constants.AzureCloud().ResourceManagerUrl)
		AzMgmtUrl.Set(url)
	}
}

func CheckCollectionConfigSanity(log logr.Logger) {
	useSaneIntValues(ColBatchSize, log)
	useSaneIntValues(ColMaxConnsPerHost, log)
	useSaneIntValues(ColMaxIdleConnsPerHost, log)
	useSaneIntValues(ColStreamCount, log)
}

// ValidateLoggingConfig checks file logging settings before the logger is
// created. Logging limits are ignored when file logging is disabled.
func ValidateLoggingConfig() error {
	if logFile, ok := LogFile.Value().(string); !ok || logFile == "" {
		return nil
	} else if fileInfo, err := os.Stat(logFile); err == nil && fileInfo.IsDir() {
		return fmt.Errorf("%s must reference a file, not a directory: %q", LogFile.Name, logFile)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("could not inspect %s %q: %w", LogFile.Name, logFile, err)
	}

	if value := LogMaxSize.Value().(int); value < LogMaxSize.MinValue {
		return fmt.Errorf("%s must be at least %d", LogMaxSize.Name, LogMaxSize.MinValue)
	}
	if value := LogMaxAge.Value().(int); value < LogMaxAge.MinValue {
		return fmt.Errorf("%s must be at least %d", LogMaxAge.Name, LogMaxAge.MinValue)
	}
	if value := LogMaxBackups.Value().(int); value < LogMaxBackups.MinValue {
		return fmt.Errorf("%s must be at least %d", LogMaxBackups.Name, LogMaxBackups.MinValue)
	}

	return nil
}

func useSaneIntValues(c config.Config, log logr.Logger) {
	val := c.Value().(int)
	if val < c.MinValue {
		log.V(1).Info(fmt.Sprintf("Provided value %d for config option %s is less than minimum value %d. Using default value %d.", val, c.Name, c.MinValue, c.Default))
		c.Set(c.Default)
	} else if val > c.MaxValue {
		log.V(1).Info(fmt.Sprintf("Provided value %d for config option %s is greater than maximum value %d. Using default value %d.", val, c.Name, c.MaxValue, c.Default))
		c.Set(c.Default)
	}
}

func ValidateURL(input string) error {
	if parsedURL, err := url.Parse(input); err != nil {
		return err
	} else if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return fmt.Errorf("invalid URL")
	} else {
		return nil
	}
}

func Options() config.Options {
	return config.Options{
		ConfigFile:  ConfigFile.Value().(string),
		ConfigName:  "config",
		ConfigPaths: SystemConfigDirs(),
		EnvPrefix:   EnvPrefix,
	}
}

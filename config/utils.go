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

	client "github.com/certmichelin/azurehound/v2/client/config"
	config "github.com/certmichelin/azurehound/v2/config/internal"
	"github.com/certmichelin/azurehound/v2/constants"
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

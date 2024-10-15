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

package config_test

import (
	"testing"

	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/logger"
)

func TestCheckCollectionConfigSanity(t *testing.T) {
	config.JsonLogs.Set(true)

	if logr, err := logger.GetLogger(); err != nil {
		t.Errorf("Error creating logger: %v", err)
	} else {
		log := *logr
		config.CheckCollectionConfigSanity(log)

		if config.ColBatchSize.Value().(int) != config.ColBatchSize.Default {
			t.Errorf("ColBatchSize did not have the default value of %d. Actual: %d", config.ColBatchSize.Default, config.ColBatchSize.Value())
		}

		if config.ColMaxConnsPerHost.Value().(int) != config.ColMaxConnsPerHost.Default {
			t.Errorf("ColMaxConnsPerHost did not have the default value of %d. Actual: %d", config.ColMaxConnsPerHost.Default, config.ColMaxConnsPerHost.Value())
		}

		if config.ColMaxIdleConnsPerHost.Value().(int) != config.ColMaxIdleConnsPerHost.Default {
			t.Errorf("ColMaxIdleConnsPerHost did not have the default value of %d. Actual: %d", config.ColMaxIdleConnsPerHost.Default, config.ColMaxIdleConnsPerHost.Value())
		}

		if config.ColStreamCount.Value().(int) != config.ColStreamCount.Default {
			t.Errorf("ColStreamCount did not have the default value of %d. Actual: %d", config.ColStreamCount.Default, config.ColStreamCount.Value())
		}
	}
}

func TestCheckCollectionConfigSanityOutOfBounds(t *testing.T) {
	config.JsonLogs.Set(true)

	if logr, err := logger.GetLogger(); err != nil {
		t.Errorf("Error creating logger: %v", err)
	} else {
		log := *logr

		config.ColBatchSize.Set(9999)
		config.ColMaxConnsPerHost.Set(-9999)

		config.CheckCollectionConfigSanity(log)

		if config.ColBatchSize.Value().(int) != config.ColBatchSize.Default {
			t.Errorf("ColBatchSize should have reverted to the default value of %d. Actual: %d", config.ColBatchSize.Default, config.ColBatchSize.Value())
		}

		if config.ColMaxConnsPerHost.Value().(int) != config.ColMaxConnsPerHost.Default {
			t.Errorf("ColMaxConnsPerHost should have reverted to the default value of %d. Actual: %d", config.ColMaxConnsPerHost.Default, config.ColMaxConnsPerHost.Value())
		}
	}
}

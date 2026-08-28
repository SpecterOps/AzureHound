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
	"strings"
	"testing"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/logger"
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

func TestValidateLoggingConfig(t *testing.T) {
	oldLogFile := config.LogFile.Value()
	oldMaxSize := config.LogMaxSize.Value()
	oldMaxAge := config.LogMaxAge.Value()
	oldMaxBackups := config.LogMaxBackups.Value()
	defer func() {
		config.LogFile.Set(oldLogFile)
		config.LogMaxSize.Set(oldMaxSize)
		config.LogMaxAge.Set(oldMaxAge)
		config.LogMaxBackups.Set(oldMaxBackups)
	}()

	config.LogFile.Set("azurehound.log")
	config.LogMaxSize.Set(100)
	config.LogMaxAge.Set(14)
	config.LogMaxBackups.Set(20)
	if err := config.ValidateLoggingConfig(); err != nil {
		t.Fatalf("valid logging configuration returned an error: %v", err)
	}

	tests := []struct {
		name   string
		config config.Config
		value  int
	}{
		{name: "zero max size", config: config.LogMaxSize, value: 0},
		{name: "negative max age", config: config.LogMaxAge, value: -1},
		{name: "negative max backups", config: config.LogMaxBackups, value: -1},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			original := test.config.Value()
			test.config.Set(test.value)
			defer test.config.Set(original)
			if err := config.ValidateLoggingConfig(); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}

	config.LogMaxAge.Set(0)
	config.LogMaxBackups.Set(0)
	if err := config.ValidateLoggingConfig(); err != nil {
		t.Fatalf("disabled pruning returned an error: %v", err)
	}
}

func TestValidateLoggingConfigIgnoredWithoutLogFile(t *testing.T) {
	oldLogFile := config.LogFile.Value()
	oldMaxSize := config.LogMaxSize.Value()
	defer func() {
		config.LogFile.Set(oldLogFile)
		config.LogMaxSize.Set(oldMaxSize)
	}()

	config.LogFile.Set("")
	config.LogMaxSize.Set(0)
	if err := config.ValidateLoggingConfig(); err != nil {
		t.Fatalf("file logging limits should be ignored without a log file: %v", err)
	}
}

func TestValidateLoggingConfigRejectsDirectoryLogPath(t *testing.T) {
	oldLogFile := config.LogFile.Value()
	oldMaxSize := config.LogMaxSize.Value()
	oldMaxAge := config.LogMaxAge.Value()
	oldMaxBackups := config.LogMaxBackups.Value()
	defer func() {
		config.LogFile.Set(oldLogFile)
		config.LogMaxSize.Set(oldMaxSize)
		config.LogMaxAge.Set(oldMaxAge)
		config.LogMaxBackups.Set(oldMaxBackups)
	}()

	logPath := t.TempDir()
	config.LogFile.Set(logPath)
	config.LogMaxSize.Set(config.DefaultLogMaxSize)
	config.LogMaxAge.Set(config.DefaultLogMaxAge)
	config.LogMaxBackups.Set(config.DefaultLogMaxBackups)

	err := config.ValidateLoggingConfig()
	if err == nil {
		t.Fatal("expected a directory-valued log path to be rejected")
	}
	if !strings.Contains(err.Error(), "not a directory") {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

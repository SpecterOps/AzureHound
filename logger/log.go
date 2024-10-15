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

package logger

import (
	"io"
	"os"

	"github.com/certmichelin/azurehound/v2/config"
	logger "github.com/certmichelin/azurehound/v2/logger/internal"
	"github.com/go-logr/logr"
)

func setupLogger() (*logr.Logger, error) {
	options := logger.Options{
		Level:      config.VerbosityLevel.Value().(int),
		Structured: config.JsonLogs.Value().(bool),
		Colors:     true,
		Writers:    []io.Writer{os.Stderr},
	}

	// emit logs to file if configured
	if fileLogWriter := getFileLogLevelWriter(); fileLogWriter != nil {
		options.Writers = append(options.Writers, fileLogWriter)
	}

	logr := logger.NewLogger(options)
	return &logr, nil
}

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

package logger

import (
	"io"

	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/go-logr/logr"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	log           *logr.Logger
	fileLogWriter io.Writer
)

func getFileLogLevelWriter() (io.Writer, error) {
	if fileLogWriter != nil {
		return fileLogWriter, nil
	} else if err := config.ValidateLoggingConfig(); err != nil {
		return nil, err
	} else if logfile, ok := config.LogFile.Value().(string); !ok || logfile == "" {
		return nil, nil
	} else {
		fileLogWriter = &lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    config.LogMaxSize.Value().(int),
			MaxAge:     config.LogMaxAge.Value().(int),
			MaxBackups: config.LogMaxBackups.Value().(int),
			Compress:   config.LogCompress.Value().(bool),
		}
		return fileLogWriter, nil
	}
}

func GetLogger() (*logr.Logger, error) {
	if log != nil {
		return log, nil
	}

	if logr, err := setupLogger(); err != nil {
		return nil, err
	} else {
		log = logr
		return log, nil
	}
}

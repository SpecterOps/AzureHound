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
	"os"

	"github.com/certmichelin/azurehound/v2/config"
	"github.com/go-logr/logr"
)

var (
	log           *logr.Logger
	fileLogWriter io.Writer
)

func getFileLogLevelWriter() io.Writer {
	if fileLogWriter != nil {
		return fileLogWriter
	} else if logfile, ok := config.LogFile.Value().(string); !ok || logfile == "" {
		return nil
	} else if file, err := os.OpenFile(logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666); err != nil {
		return nil
	} else {
		return file
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

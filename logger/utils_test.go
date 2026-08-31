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
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/SpecterOps/AzureHound/v2/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestGetFileLogLevelWriterUsesLumberjack(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "azurehound.log")
	restore := setFileLoggingConfig(logPath)
	defer restore()

	config.LogMaxSize.Set(25)
	config.LogMaxAge.Set(7)
	config.LogMaxBackups.Set(5)
	config.LogCompress.Set(false)

	fileWriter, err := getFileLogLevelWriter()
	if err != nil {
		t.Fatal(err)
	}
	writer, ok := fileWriter.(*lumberjack.Logger)
	if !ok {
		t.Fatalf("expected *lumberjack.Logger, got %T", fileLogWriter)
	}
	if writer.Filename != logPath {
		t.Errorf("Filename = %q, want %q", writer.Filename, logPath)
	}
	if writer.MaxSize != 25 {
		t.Errorf("MaxSize = %d, want 25", writer.MaxSize)
	}
	if writer.MaxAge != 7 {
		t.Errorf("MaxAge = %d, want 7", writer.MaxAge)
	}
	if writer.MaxBackups != 5 {
		t.Errorf("MaxBackups = %d, want 5", writer.MaxBackups)
	}
	if writer.Compress {
		t.Error("Compress = true, want false")
	}
}

func TestGetLoggerRejectsDirectoryLogPath(t *testing.T) {
	logPath := t.TempDir()
	restore := setFileLoggingConfig(logPath)
	defer restore()

	if _, err := GetLogger(); err == nil {
		t.Fatal("expected logger setup to reject a directory-valued log path")
	}

	fileInfo, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("stat log directory: %v", err)
	}
	if !fileInfo.IsDir() {
		t.Fatalf("log path was changed from a directory: mode = %v", fileInfo.Mode())
	}
	if fileLogWriter != nil {
		t.Fatalf("file log writer was created for a directory: %T", fileLogWriter)
	}
}

func TestGetLoggerRejectsInaccessibleLogFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file modes do not reliably prevent writes on Windows")
	}

	logPath := filepath.Join(t.TempDir(), "azurehound.log")
	if err := os.WriteFile(logPath, nil, 0400); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chmod(logPath, 0600); err != nil {
			t.Errorf("restore log file permissions: %v", err)
		}
	}()

	if file, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY, 0); err == nil {
		_ = file.Close()
		t.Skip("test requires write permissions to be enforced")
	}

	restore := setFileLoggingConfig(logPath)
	defer restore()

	if _, err := GetLogger(); err == nil {
		t.Fatal("expected logger setup to reject an inaccessible log file")
	}

	fileInfo, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("stat log file: %v", err)
	}
	if !fileInfo.Mode().IsRegular() {
		t.Fatalf("log path was changed from a regular file: mode = %v", fileInfo.Mode())
	}
	if fileLogWriter != nil {
		t.Fatalf("file log writer was created for an inaccessible file: %T", fileLogWriter)
	}
	if log != nil {
		t.Fatal("logger was cached for an inaccessible log file")
	}
}

func TestGetLoggerWritesToRotatingFile(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "logs", "azurehound.log")
	restore := setFileLoggingConfig(logPath)
	defer restore()

	log, err := GetLogger()
	if err != nil {
		t.Fatal(err)
	}
	log.Info("rotating file log", "answer", 42)
	resetLoggerForTest()

	contents, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "rotating file log") {
		t.Fatalf("log file does not contain emitted message: %s", contents)
	}
}

func setFileLoggingConfig(logPath string) func() {
	oldLogFile := config.LogFile.Value()
	oldMaxSize := config.LogMaxSize.Value()
	oldMaxAge := config.LogMaxAge.Value()
	oldMaxBackups := config.LogMaxBackups.Value()
	oldCompress := config.LogCompress.Value()
	oldJSON := config.JsonLogs.Value()
	oldVerbosity := config.VerbosityLevel.Value()

	resetLoggerForTest()
	config.LogFile.Set(logPath)
	config.LogMaxSize.Set(config.DefaultLogMaxSize)
	config.LogMaxAge.Set(config.DefaultLogMaxAge)
	config.LogMaxBackups.Set(config.DefaultLogMaxBackups)
	config.LogCompress.Set(true)
	config.JsonLogs.Set(true)
	config.VerbosityLevel.Set(0)

	return func() {
		resetLoggerForTest()
		config.LogFile.Set(oldLogFile)
		config.LogMaxSize.Set(oldMaxSize)
		config.LogMaxAge.Set(oldMaxAge)
		config.LogMaxBackups.Set(oldMaxBackups)
		config.LogCompress.Set(oldCompress)
		config.JsonLogs.Set(oldJSON)
		config.VerbosityLevel.Set(oldVerbosity)
	}
}

func resetLoggerForTest() {
	if closer, ok := fileLogWriter.(io.Closer); ok {
		_ = closer.Close()
	}
	fileLogWriter = nil
	log = nil
}

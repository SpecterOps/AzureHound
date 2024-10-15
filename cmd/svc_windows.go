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

	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/logger"
	"github.com/judwhite/go-svc"
)

func StartWindowsService() error {
	if err := svc.Run(&azurehoundSvc{}); err != nil {
		return err
	} else {
		return nil
	}
}

type azurehoundSvc struct {
	cancel context.CancelFunc
}

func (s *azurehoundSvc) Init(env svc.Environment) error {
	config.LoadValues(nil, config.Options())
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

func (s *azurehoundSvc) Start() error {
	if err := s.Stop(); err != nil {
		return err
	} else {
		log.Info("starting azurehound service...", "config", config.ConfigFileUsed())
		ctx, stop := context.WithCancel(context.Background())
		s.cancel = stop
		go start(ctx)
		return nil
	}
}

func (s *azurehoundSvc) Stop() error {
	if s.cancel != nil {
		log.Info("stopping azurehound service...")
		s.cancel()
		s = nil
	}
	return nil
}

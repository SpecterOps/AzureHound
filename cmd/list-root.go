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
	"os"
	"os/signal"
	"time"

	"github.com/certmichelin/azurehound/v2/client"
	"github.com/certmichelin/azurehound/v2/config"
	"github.com/certmichelin/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	config.Init(listRootCmd, append(config.AzureConfig, config.OutputFile))
	rootCmd.AddCommand(listRootCmd)
}

var listRootCmd = &cobra.Command{
	Use:               "list",
	Short:             "Lists Azure Objects",
	Run:               listCmdImpl,
	PersistentPreRunE: persistentPreRunE,
	SilenceUsage:      true,
}

func listCmdImpl(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		exit(fmt.Errorf("unsupported subcommand: %v", args))
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting azure objects...")
	start := time.Now()
	stream := listAll(ctx, azClient)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listAll(ctx context.Context, client client.AzureClient) <-chan interface{} {
	var (
		azureAD = listAllAD(ctx, client)
		azureRM = listAllRM(ctx, client)
	)
	return pipeline.Mux(ctx.Done(), azureAD, azureRM)
}

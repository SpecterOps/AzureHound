// File: cmd/list-intune-devices.go
// Copyright (C) 2022 SpecterOps
// Command implementation for listing Intune managed devices

package cmd

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listIntuneDevicesCmd)
}

var listIntuneDevicesCmd = &cobra.Command{
	Use:          "intune-devices",
	Long:         "Lists Intune Managed Devices",
	Run:          listIntuneDevicesCmdImpl,
	SilenceUsage: true,
}

func listIntuneDevicesCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting intune managed devices...")
	start := time.Now()
	stream := listIntuneDevices(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listIntuneDevices(ctx context.Context, client client.AzureClient) <-chan interface{} {
	var (
		out    = make(chan interface{})
		params = query.GraphParams{
			Filter: "operatingSystem eq 'Windows'", // Focus on Windows devices for BloodHound
		}
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		count := 0
		for item := range client.ListIntuneManagedDevices(ctx, params) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing intune devices")
			} else {
				log.V(2).Info("found intune device", "device", item.Ok)
				count++
				select {
				case out <- NewAzureWrapper(enums.KindAZIntuneDevice, item.Ok):
				case <-ctx.Done():
					return
				}
			}
		}
		log.V(1).Info("finished listing intune devices", "count", count)
	}()

	return out
}
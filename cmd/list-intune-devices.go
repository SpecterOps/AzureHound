// File: cmd/list-intune-devices.go
// Copyright (C) 2022 SpecterOps
// Command implementation for listing Intune managed devices

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
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

	// Add proper error handling for client creation
	azClient, err := connectAndCreateClientWithError()
	if err != nil {
		log.Error(err, "failed to create Azure client")
		os.Exit(1)
	}

	if devices, err := listIntuneDevices(ctx, azClient); err != nil {
		exit(err)
	} else {
		// Simple output - just print device count for now
		fmt.Printf("Found %d Intune devices\n", len(devices))

		// Print basic device info
		for _, device := range devices {
			fmt.Printf("Device: %s (%s) - %s\n",
				device.DeviceName,
				device.OperatingSystem,
				device.ComplianceState)
		}
	}
}

func listIntuneDevices(ctx context.Context, azClient client.AzureClient) ([]azure.IntuneDevice, error) {
	var (
		out     = make([]azure.IntuneDevice, 0)
		devices = azClient.ListIntuneDevices(ctx, query.GraphParams{})
	)

	for result := range devices {
		if result.Error != nil {
			return nil, result.Error
		} else {
			out = append(out, result.Ok)
		}
	}

	return out, nil
}

// Helper function to create client with proper error handling
// Helper function to create client with proper error handling
func connectAndCreateClientWithError() (azClient client.AzureClient, err error) {
	// Use named return parameters to allow deferred function to modify them
	defer func() {
		if r := recover(); r != nil {
			// Convert panic to error instead of re-panicking
			err = fmt.Errorf("failed to create client: %v", r)
			azClient = nil
		}
	}()

	azClient = connectAndCreateClient()
	return azClient, nil
}

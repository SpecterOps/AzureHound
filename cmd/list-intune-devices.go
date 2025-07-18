// File: cmd/list-intune-devices.go
// Copyright (C) 2022 SpecterOps
// Command implementation for listing Intune managed devices

package cmd

import (
	"context"
	"fmt"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/models/intune"
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
	ctx, stop := context.WithCancel(cmd.Context())
	defer stop()

	azClient := connectAndCreateClient()

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

func listIntuneDevices(ctx context.Context, azClient client.AzureClient) ([]intune.ManagedDevice, error) {
	var (
		out     = make([]intune.ManagedDevice, 0)
		devices = azClient.ListIntuneManagedDevices(ctx, query.GraphParams{})
		count   = 0
	)

	for result := range devices {
		if result.Error != nil {
			return nil, result.Error
		} else {
			count++
			out = append(out, result.Ok)
		}
	}

	return out, nil
}

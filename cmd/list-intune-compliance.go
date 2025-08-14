// File: cmd/list-intune-compliance.go
// Command for listing Intune device compliance information

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/config"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models/intune"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func createBasicComplianceState(device intune.ManagedDevice, suffix string) intune.ComplianceState {
	return intune.ComplianceState{
		Id:         device.Id + suffix,
		DeviceId:   device.Id,
		DeviceName: device.DeviceName,
		State:      device.ComplianceState,
		Version:    1,
	}
}

func init() {
	listRootCmd.AddCommand(listIntuneComplianceCmd)

	listIntuneComplianceCmd.Flags().StringP("state", "s", "", "Filter by compliance state: compliant, noncompliant, conflict, error, unknown")
	listIntuneComplianceCmd.Flags().BoolP("details", "d", false, "Include detailed compliance settings")
}

var listIntuneComplianceCmd = &cobra.Command{
	Use:   "intune-compliance",
	Short: "List Intune device compliance information",
	Long: `List compliance information for Intune managed devices.

Examples:
  # List all device compliance
  azurehound list intune-compliance --jwt $JWT

  # List only non-compliant devices
  azurehound list intune-compliance --state noncompliant --jwt $JWT

  # Include detailed compliance settings
  azurehound list intune-compliance --details --jwt $JWT`,
	Run:          listIntuneComplianceCmdImpl,
	SilenceUsage: true,
}

func listIntuneComplianceCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	// Get flag values locally to avoid global variable issues
	complianceState, _ := cmd.Flags().GetString("state")
	includeDetails, _ := cmd.Flags().GetBool("details")

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting intune device compliance...")
	start := time.Now()
	stream := listIntuneCompliance(ctx, azClient, complianceState, includeDetails)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listIntuneCompliance(ctx context.Context, client client.AzureClient, complianceState string, includeDetails bool) <-chan interface{} {
	var (
		out = make(chan interface{})
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		// First get all managed devices
		devices := getComplianceTargetDevices(ctx, client, complianceState)

		// Then collect compliance data for each device
		collectDeviceCompliance(ctx, client, devices, out, includeDetails)
	}()

	return out
}

func getComplianceTargetDevices(ctx context.Context, client client.AzureClient, complianceState string) <-chan intune.ManagedDevice {
	var (
		out    = make(chan intune.ManagedDevice)
		params = query.GraphParams{
			Filter: "operatingSystem eq 'Windows'",
		}
	)

	// Apply compliance state filter if specified
	if complianceState != "" {
		if params.Filter != "" {
			params.Filter += " and "
		}
		params.Filter += fmt.Sprintf("complianceState eq '%s'", complianceState)
	}

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		count := 0
		for item := range client.ListIntuneManagedDevices(ctx, params) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing devices")
			} else {
				log.V(2).Info("found device for compliance check", "device", item.Ok.DeviceName)
				count++
				select {
				case out <- item.Ok:
				case <-ctx.Done():
					return
				}
			}
		}
		log.V(1).Info("finished collecting target devices", "count", count)
	}()

	return out
}

func collectDeviceCompliance(ctx context.Context, client client.AzureClient, devices <-chan intune.ManagedDevice, out chan<- interface{}, includeDetails bool) {
	var (
		streams = pipeline.Demux(ctx.Done(), devices, config.ColStreamCount.Value().(int))
		wg      sync.WaitGroup
	)

	wg.Add(len(streams))
	for i := range streams {
		stream := streams[i]
		go func() {
			defer panicrecovery.PanicRecovery()
			defer wg.Done()

			for device := range stream {
				if includeDetails {
					collectDetailedCompliance(ctx, client, device, out)
				} else {
					basicCompliance := createBasicComplianceState(device, "-basic")
					select {
					case out <- NewAzureWrapper(enums.KindAZIntuneCompliance, basicCompliance):
					case <-ctx.Done():
						return
					}
				}
			}
		}()
	}
	wg.Wait()
}

func collectDetailedCompliance(ctx context.Context, client client.AzureClient, device intune.ManagedDevice, out chan<- interface{}) {
	log.V(2).Info("collecting detailed compliance", "device", device.DeviceName)

	params := query.GraphParams{}
	count := 0

	for complianceResult := range client.GetIntuneDeviceCompliance(ctx, device.Id, params) {
		if complianceResult.Error != nil {
			log.Error(complianceResult.Error, "failed to get detailed compliance", "device", device.DeviceName)

			// Fall back to basic compliance info using helper
			basicCompliance := createBasicComplianceState(device, "-fallback")
			select {
			case out <- NewAzureWrapper(enums.KindAZIntuneCompliance, basicCompliance):
			case <-ctx.Done():
				return
			}
			return // Changed from continue to return as suggested
		}

		log.V(2).Info("found detailed compliance state",
			"device", device.DeviceName,
			"state", complianceResult.Ok.State,
			"settingsCount", len(complianceResult.Ok.SettingStates))

		count++
		select {
		case out <- NewAzureWrapper(enums.KindAZIntuneCompliance, complianceResult.Ok):
		case <-ctx.Done():
			return
		}
	}

	if count > 0 {
		log.V(1).Info("finished detailed compliance collection", "device", device.DeviceName, "policies", count)
	}
}

// File: cmd/list-intune-compliance.go
// Command for listing Intune device compliance information with configurable OS filter

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
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

var (
	complianceState string
	includeDetails  bool
	operatingSystem string // New flag for OS filter
)

func init() {
	listRootCmd.AddCommand(listIntuneComplianceCmd)

	listIntuneComplianceCmd.Flags().StringVar(&complianceState, "state", "", "Filter by compliance state: compliant, noncompliant, conflict, error, unknown")
	listIntuneComplianceCmd.Flags().BoolVar(&includeDetails, "details", false, "Include detailed compliance settings")
	listIntuneComplianceCmd.Flags().StringVar(&operatingSystem, "os", "Windows", "Filter by operating system (e.g., Windows, Android, iOS, macOS). Use 'all' for no OS filtering")
}

var listIntuneComplianceCmd = &cobra.Command{
	Use:   "intune-compliance",
	Short: "List Intune device compliance information",
	Long: `List compliance information for Intune managed devices with configurable OS filtering.

Examples:
  # List all Windows device compliance (default)
  azurehound list intune-compliance --jwt $JWT

  # List compliance for all operating systems
  azurehound list intune-compliance --os all --jwt $JWT

  # List only Android devices that are non-compliant
  azurehound list intune-compliance --os Android --state noncompliant --jwt $JWT

  # Include detailed compliance settings for iOS devices
  azurehound list intune-compliance --os iOS --details --jwt $JWT`,
	Run:          listIntuneComplianceCmdImpl,
	SilenceUsage: true,
}

func listIntuneComplianceCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting intune device compliance...")
	start := time.Now()
	stream := listIntuneCompliance(ctx, azClient)
	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listIntuneCompliance(ctx context.Context, client client.AzureClient) <-chan interface{} {
	var (
		out = make(chan interface{})
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		// First get all managed devices with configurable OS filter
		devices := getComplianceTargetDevices(ctx, client, operatingSystem, complianceState)

		// Then collect compliance data for each device
		collectDeviceCompliance(ctx, client, devices, out)
	}()

	return out
}

// getComplianceTargetDevices retrieves devices based on configurable OS and compliance state filters
// Parameters:
//   - ctx: Context for cancellation
//   - client: AzureClient instance
//   - osFilter: Operating system filter ("Windows", "Android", "iOS", "macOS", or "all" for no filtering)
//   - complianceFilter: Compliance state filter (optional)
//
// Returns a channel of ManagedDevice objects matching the specified filters
func getComplianceTargetDevices(ctx context.Context, client client.AzureClient, osFilter, complianceFilter string) <-chan intune.ManagedDevice {
	var (
		out     = make(chan intune.ManagedDevice)
		params  = query.GraphParams{}
		filters []string
	)

	// Apply OS filtering if not "all"
	if osFilter != "" && strings.ToLower(osFilter) != "all" {
		filters = append(filters, fmt.Sprintf("operatingSystem eq '%s'", osFilter))
		log.V(1).Info("applying OS filter", "operatingSystem", osFilter)
	} else {
		log.V(1).Info("no OS filtering applied - collecting all operating systems")
	}

	// Apply compliance state filter if specified
	if complianceFilter != "" {
		filters = append(filters, fmt.Sprintf("complianceState eq '%s'", complianceFilter))
		log.V(1).Info("applying compliance filter", "complianceState", complianceFilter)
	}

	// Combine filters with AND operator
	if len(filters) > 0 {
		params.Filter = strings.Join(filters, " and ")
		log.V(1).Info("final filter applied", "filter", params.Filter)
	}

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		count := 0
		skipped := 0
		for item := range client.ListIntuneManagedDevices(ctx, params) {
			if item.Error != nil {
				log.Error(item.Error, "unable to continue processing devices")
			} else {
				// Additional client-side filtering for edge cases where server-side filtering might not work perfectly
				if shouldIncludeDevice(item.Ok, osFilter, complianceFilter) {
					log.V(2).Info("found device for compliance check",
						"device", item.Ok.DeviceName,
						"os", item.Ok.OperatingSystem,
						"compliance", item.Ok.ComplianceState)
					count++
					select {
					case out <- item.Ok:
					case <-ctx.Done():
						return
					}
				} else {
					skipped++
					log.V(2).Info("skipping device due to filter mismatch",
						"device", item.Ok.DeviceName,
						"os", item.Ok.OperatingSystem,
						"compliance", item.Ok.ComplianceState)
				}
			}
		}
		log.V(1).Info("finished collecting target devices", "included", count, "skipped", skipped)
	}()

	return out
}

// shouldIncludeDevice performs additional client-side validation of filters
func shouldIncludeDevice(device intune.ManagedDevice, osFilter, complianceFilter string) bool {
	// Check OS filter
	if osFilter != "" && strings.ToLower(osFilter) != "all" {
		if !strings.EqualFold(device.OperatingSystem, osFilter) {
			return false
		}
	}

	// Check compliance filter
	if complianceFilter != "" {
		if !strings.EqualFold(device.ComplianceState, complianceFilter) {
			return false
		}
	}

	return true
}

func collectDeviceCompliance(ctx context.Context, client client.AzureClient, devices <-chan intune.ManagedDevice, out chan<- interface{}) {
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
			continue
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

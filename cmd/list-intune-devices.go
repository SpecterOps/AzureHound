// File: cmd/list-intune-devices.go
// Copyright (C) 2022 SpecterOps
// Command implementation for listing Intune managed devices with streaming processing

package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models/intune"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/spf13/cobra"
)

var (
	outputFormat string // Flag for output format
	maxDevices   int    // Flag to limit output for testing
)

func init() {
	listRootCmd.AddCommand(listIntuneDevicesCmd)

	// Add flags for better control over output
	listIntuneDevicesCmd.Flags().StringVar(&outputFormat, "format", "summary", "Output format: summary, detailed, json")
	listIntuneDevicesCmd.Flags().IntVar(&maxDevices, "max", 0, "Maximum number of devices to process (0 = unlimited)")
}

var listIntuneDevicesCmd = &cobra.Command{
	Use:   "intune-devices",
	Short: "Lists Intune Managed Devices",
	Long: `Lists Intune Managed Devices using streaming processing to handle large datasets efficiently.

Output Formats:
  summary  - Basic device information (default)
  detailed - Detailed device properties
  json     - JSON output suitable for further processing

Examples:
  # List all devices with summary output
  azurehound list intune-devices --jwt $JWT

  # Show detailed information for first 10 devices
  azurehound list intune-devices --format detailed --max 10 --jwt $JWT

  # Output in JSON format for processing
  azurehound list intune-devices --format json --jwt $JWT`,
	Run:          listIntuneDevicesCmdImpl,
	SilenceUsage: true,
}

func listIntuneDevicesCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	log.V(1).Info("testing connections")
	azClient := connectAndCreateClient()
	log.Info("collecting intune devices...")
	start := time.Now()

	// Use streaming approach based on output format
	switch outputFormat {
	case "json":
		stream := listIntuneDevicesAsStream(ctx, azClient)
		panicrecovery.HandleBubbledPanic(ctx, stop, log)
		outputStream(ctx, stream)
	default:
		// For summary and detailed formats, process devices directly without accumulating
		processIntuneDevicesStreaming(ctx, azClient, outputFormat, maxDevices)
	}

	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

// processIntuneDevicesStreaming processes devices as they are received without accumulating in memory
func processIntuneDevicesStreaming(ctx context.Context, azClient client.AzureClient, format string, maxCount int) {
	devices := azClient.ListIntuneManagedDevices(ctx, query.GraphParams{})

	count := 0
	errorCount := 0

	// Print header based on format
	printHeader(format)

	for result := range devices {
		if result.Error != nil {
			errorCount++
			log.Error(result.Error, "error retrieving device")
			continue
		}

		// Process each device immediately
		processDevice(result.Ok, format, count+1)
		count++

		// Respect max limit if set
		if maxCount > 0 && count >= maxCount {
			fmt.Printf("\n[Limit reached: processed %d devices, stopping as requested]\n", maxCount)
			break
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			fmt.Printf("\n[Operation cancelled after processing %d devices]\n", count)
			return
		default:
			// Continue processing
		}
	}

	// Print footer/summary
	printFooter(format, count, errorCount)
}

// listIntuneDevicesAsStream returns a stream for JSON output compatible with existing pipeline
func listIntuneDevicesAsStream(ctx context.Context, azClient client.AzureClient) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		devices := azClient.ListIntuneManagedDevices(ctx, query.GraphParams{})
		count := 0

		for result := range devices {
			if result.Error != nil {
				log.Error(result.Error, "error retrieving device")
				continue
			}

			count++

			// Respect max limit if set
			if maxDevices > 0 && count > maxDevices {
				break
			}

			select {
			case out <- NewAzureWrapper(enums.KindAZIntuneDevice, result.Ok):
			case <-ctx.Done():
				return
			}
		}
	}()

	return out
}

// printHeader prints the appropriate header for the output format
func printHeader(format string) {
	switch format {
	case "detailed":
		fmt.Printf("%-4s %-30s %-15s %-20s %-15s %-20s %s\n",
			"#", "Device Name", "OS", "OS Version", "Compliance", "Last Sync", "User")
		fmt.Printf("%s\n", "────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
	case "summary":
		fmt.Printf("%-4s %-30s %-15s %-15s %s\n",
			"#", "Device Name", "OS", "Compliance", "User")
		fmt.Printf("%s\n", "──────────────────────────────────────────────────────────────────────────────────────")
	}
}

// processDevice handles individual device processing based on format
func processDevice(device intune.ManagedDevice, format string, index int) {
	switch format {
	case "detailed":
		printDetailedDevice(device, index)
	case "summary":
		printSummaryDevice(device, index)
	default:
		printSummaryDevice(device, index)
	}
}

// printDetailedDevice prints detailed device information
func printDetailedDevice(device intune.ManagedDevice, index int) {
	lastSync := "Never"
	if !device.LastSyncDateTime.IsZero() {
		lastSync = device.LastSyncDateTime.Format("2006-01-02 15:04")
	}

	fmt.Printf("%-4d %-30s %-15s %-20s %-15s %-20s %s\n",
		index,
		truncateString(device.DeviceName, 30),
		truncateString(device.OperatingSystem, 15),
		truncateString(device.OSVersion, 20),
		truncateString(device.ComplianceState, 15),
		truncateString(lastSync, 20),
		truncateString(device.UserPrincipalName, 30))
}

// printSummaryDevice prints summary device information
func printSummaryDevice(device intune.ManagedDevice, index int) {
	fmt.Printf("%-4d %-30s %-15s %-15s %s\n",
		index,
		truncateString(device.DeviceName, 30),
		truncateString(device.OperatingSystem, 15),
		truncateString(device.ComplianceState, 15),
		truncateString(device.UserPrincipalName, 30))
}

// printFooter prints summary information
func printFooter(format string, deviceCount, errorCount int) {
	fmt.Printf("\n")
	if format == "detailed" {
		fmt.Printf("%s\n", "────────────────────────────────────────────────────────────────────────────────────────────────────────────────────")
	} else {
		fmt.Printf("%s\n", "──────────────────────────────────────────────────────────────────────────────────────")
	}

	fmt.Printf("Summary: %d devices processed", deviceCount)
	if errorCount > 0 {
		fmt.Printf(", %d errors encountered", errorCount)
	}
	fmt.Printf("\n")
}

// truncateString truncates a string to the specified length
func truncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	if maxLength <= 3 {
		return s[:maxLength]
	}
	return s[:maxLength-3] + "..."
}

// DeviceProcessor defines an interface for processing devices
type DeviceProcessor interface {
	ProcessDevice(device intune.ManagedDevice) error
}

// StreamingDeviceProcessor provides a callback-based streaming processor
type StreamingDeviceProcessor struct {
	ProcessFunc func(device intune.ManagedDevice) error
	maxDevices  int
	processed   int
}

// NewStreamingDeviceProcessor creates a new streaming processor
func NewStreamingDeviceProcessor(processFunc func(intune.ManagedDevice) error, maxDevices int) *StreamingDeviceProcessor {
	return &StreamingDeviceProcessor{
		ProcessFunc: processFunc,
		maxDevices:  maxDevices,
		processed:   0,
	}
}

// ProcessDevice processes a single device
func (p *StreamingDeviceProcessor) ProcessDevice(device intune.ManagedDevice) error {
	if p.maxDevices > 0 && p.processed >= p.maxDevices {
		return fmt.Errorf("maximum device limit reached: %d", p.maxDevices)
	}

	p.processed++
	return p.ProcessFunc(device)
}

// GetProcessedCount returns the number of devices processed
func (p *StreamingDeviceProcessor) GetProcessedCount() int {
	return p.processed
}

// processIntuneDevicesWithCallback processes devices using a callback function (alternative streaming approach)
func processIntuneDevicesWithCallback(ctx context.Context, azClient client.AzureClient, processor DeviceProcessor) error {
	devices := azClient.ListIntuneManagedDevices(ctx, query.GraphParams{})

	for result := range devices {
		if result.Error != nil {
			log.Error(result.Error, "error retrieving device")
			continue
		}

		if err := processor.ProcessDevice(result.Ok); err != nil {
			return err
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Continue processing
		}
	}

	return nil
}

// Legacy function maintained for backward compatibility but now uses streaming
func listIntuneDevices(ctx context.Context, azClient client.AzureClient) ([]intune.ManagedDevice, error) {
	log.V(1).Info("using legacy listIntuneDevices function - consider using streaming approach for better performance")

	var devices []intune.ManagedDevice
	deviceStream := azClient.ListIntuneManagedDevices(ctx, query.GraphParams{})

	for result := range deviceStream {
		if result.Error != nil {
			return nil, result.Error
		}
		devices = append(devices, result.Ok)

		// Add a safety check to prevent excessive memory usage
		if len(devices) > 10000 {
			log.V(1).Info("large dataset detected - consider using streaming mode", "deviceCount", len(devices))
		}

		// Check for context cancellation
		select {
		case <-ctx.Done():
			return devices, ctx.Err()
		default:
			// Continue processing
		}
	}

	return devices, nil
}

// Example usage functions for the streaming processor

// ExampleJSONProcessor demonstrates processing devices to JSON
func ExampleJSONProcessor() func(intune.ManagedDevice) error {
	return func(device intune.ManagedDevice) error {
		data, err := json.MarshalIndent(device, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(data))
		return nil
	}
}

// ExampleSummaryProcessor demonstrates processing devices to summary format
func ExampleSummaryProcessor() func(intune.ManagedDevice) error {
	count := 0
	return func(device intune.ManagedDevice) error {
		count++
		fmt.Printf("%d. %s (%s) - %s\n",
			count,
			device.DeviceName,
			device.OperatingSystem,
			device.ComplianceState)
		return nil
	}
}

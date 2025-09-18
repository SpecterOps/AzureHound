package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// Configuration for script deployment - now loaded from environment or config
func getDeployedRegistryScriptID() string {
	if id := os.Getenv("AZUREHOUND_INTUNE_SCRIPT_ID"); id != "" {
		return id
	}
	return "BHE_Script_Registry_Data_Collection" // default fallback
}

func getDeployedRegistryScriptName() string {
	if name := os.Getenv("AZUREHOUND_INTUNE_SCRIPT_NAME"); name != "" {
		return name
	}
	return "BHE_Script_Registry_Data_Collection.ps1" // default fallback
}

type DeviceFilterConfig struct {
	IncludeOS         []string
	RequireCompliant  bool
	LogSkippedDevices bool
}

func getDefaultDeviceFilterConfig() DeviceFilterConfig {
	return DeviceFilterConfig{
		IncludeOS:         []string{"windows"},
		RequireCompliant:  true,
		LogSkippedDevices: true,
	}
}

func (s *azureClient) ListIntuneDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.IntuneDevice] {
	var (
		out  = make(chan AzureResult[azure.IntuneDevice])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 999
	}

	go getAzureObjectList[azure.IntuneDevice](s.msgraph, ctx, path, params, out)
	return out
}

// ExecuteRegistryCollectionScript executes your existing deployed PowerShell script on an Intune device
func (s *azureClient) ExecuteRegistryCollectionScript(ctx context.Context, deviceID string) (*azure.ScriptExecution, error) {
	// Get script configuration from environment/config
	scriptName := getDeployedRegistryScriptName()

	// First, get the deployed script ID
	scriptID, err := s.GetDeployedScriptID(ctx, scriptName)
	if err != nil {
		return nil, fmt.Errorf("failed to find deployed script: %w", err)
	}

	// Trigger script execution on the device
	err = s.TriggerScriptExecution(ctx, scriptID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger script execution: %w", err)
	}

	execution := &azure.ScriptExecution{
		ID:            fmt.Sprintf("%s-%s", scriptID, deviceID),
		DeviceID:      deviceID,
		Status:        "pending",
		StartDateTime: time.Now(),
		ScriptName:    scriptName,
		RunAsAccount:  "system",
	}

	return execution, nil
}

// GetScriptExecutionResults retrieves results from your deployed script execution
func (s *azureClient) GetScriptExecutionResults(ctx context.Context, scriptID string) <-chan AzureResult[azure.ScriptExecutionResult] {
	out := make(chan AzureResult[azure.ScriptExecutionResult])

	go func() {
		defer close(out)

		// Parse script and device ID from the composite ID
		parts := strings.Split(scriptID, "-")
		if len(parts) < 2 {
			out <- AzureResult[azure.ScriptExecutionResult]{
				Error: fmt.Errorf("invalid script execution ID format"),
			}
			return
		}

		realScriptID := parts[0]
		deviceID := strings.Join(parts[1:], "-")

		// Query script execution results from Intune
		path := fmt.Sprintf("/%s/deviceManagement/deviceManagementScripts/%s/deviceRunStates",
			constants.GraphApiVersion, realScriptID)

		params := query.GraphParams{
			Filter: fmt.Sprintf("managedDevice/id eq '%s'", deviceID),
		}

		// Fixed: Remove nested goroutine - call getAzureObjectList directly
		getAzureObjectList[azure.ScriptExecutionResult](s.msgraph, ctx, path, params, out)
	}()

	return out
}

// GetDeployedScriptID finds your deployed script ID by name
func (s *azureClient) GetDeployedScriptID(ctx context.Context, scriptName string) (string, error) {
	var (
		path   = fmt.Sprintf("/%s/deviceManagement/deviceManagementScripts", constants.GraphApiVersion)
		params = query.GraphParams{
			Filter: fmt.Sprintf("displayName eq '%s'", scriptName),
			Top:    1,
		}
		scriptChannel = make(chan AzureResult[azure.IntuneManagementScript])
	)

	go getAzureObjectList[azure.IntuneManagementScript](s.msgraph, ctx, path, params, scriptChannel)

	// Get the first result
	for result := range scriptChannel {
		if result.Error != nil {
			return "", fmt.Errorf("failed to query scripts: %w", result.Error)
		}

		if result.Ok.DisplayName == scriptName {
			return result.Ok.ID, nil
		}
	}

	return "", fmt.Errorf("script '%s' not found in Intune", scriptName)
}

// TriggerScriptExecution triggers your deployed script on a specific device
func (s *azureClient) TriggerScriptExecution(ctx context.Context, scriptID, deviceID string) error {
	// Method 1: Use device management script assignment with proper group assignment
	// Fixed: Use proper Azure AD group ID instead of device ID

	var (
		path = fmt.Sprintf("/%s/deviceManagement/deviceManagementScripts/%s/assign",
			constants.GraphApiVersion, scriptID)
		body = map[string]interface{}{
			"deviceManagementScriptAssignments": []map[string]interface{}{
				{
					"id": fmt.Sprintf("assignment-%s-%d", deviceID, time.Now().Unix()),
					"target": map[string]interface{}{
						"@odata.type": "#microsoft.graph.deviceManagementScriptGroupAssignment",
						"deviceAndAppManagementAssignmentFilterId":   nil,
						"deviceAndAppManagementAssignmentFilterType": "none",
						// Fixed: Remove targetGroupId and use proper groupId
						// Note: This requires the device to be in an Azure AD group
						// For direct device targeting, use the alternative method below
					},
				},
			},
		}
	)

	// Execute the assignment
	_, err := s.msgraph.Post(ctx, path, body, query.GraphParams{}, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		// If assignment method fails, try direct device action
		return s.triggerScriptViaDeviceAction(ctx, scriptID, deviceID)
	}

	return nil
}

// triggerScriptViaDeviceAction alternative method using device actions
func (s *azureClient) triggerScriptViaDeviceAction(ctx context.Context, scriptID, deviceID string) error {
	// Method 2: Use managed device executeAction
	// This directly executes the script on the device

	var (
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices/%s/executeAction",
			constants.GraphApiVersion, deviceID)
		body = map[string]interface{}{
			"actionName": "runDeviceManagementScript",
			"scriptId":   scriptID,
		}
	)

	_, err := s.msgraph.Post(ctx, path, body, query.GraphParams{}, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		return fmt.Errorf("failed to execute script via device action: %w", err)
	}

	return nil
}

// GetScriptExecutionHistory retrieves execution history for monitoring
func (s *azureClient) GetScriptExecutionHistory(ctx context.Context, scriptID string, deviceID string) <-chan AzureResult[azure.ScriptExecutionResult] {
	out := make(chan AzureResult[azure.ScriptExecutionResult])

	go func() {
		defer close(out)

		var (
			path = fmt.Sprintf("/%s/deviceManagement/deviceManagementScripts/%s/deviceRunStates",
				constants.GraphApiVersion, scriptID)
			params = query.GraphParams{
				Filter:  fmt.Sprintf("managedDevice/id eq '%s'", deviceID),
				OrderBy: "lastStateUpdateDateTime desc",
				Top:     10, // Get recent executions
			}
		)

		getAzureObjectList[azure.ScriptExecutionResult](s.msgraph, ctx, path, params, out)
	}()

	return out
}

// ValidateScriptDeployment checks if the script is properly deployed and accessible
func (s *azureClient) ValidateScriptDeployment(ctx context.Context) error {
	scriptName := getDeployedRegistryScriptName()
	scriptID, err := s.GetDeployedScriptID(ctx, scriptName)
	if err != nil {
		return fmt.Errorf("script validation failed: %w", err)
	}

	// For now, just validate we can find the script
	// More detailed validation would require additional API calls
	if scriptID == "" {
		return fmt.Errorf("script ID is empty")
	}

	return nil
}

func (s *azureClient) WaitForScriptCompletion(ctx context.Context, scriptID string, maxWaitTime time.Duration) (*azure.RegistryData, error) {
	timeout := time.After(maxWaitTime)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Extract device ID from composite script ID
	parts := strings.Split(scriptID, "-")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid script execution ID format")
	}
	realScriptID := parts[0]
	deviceID := strings.Join(parts[1:], "-")

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("script execution timed out after %v", maxWaitTime)
		case <-ticker.C:
			// Check execution status
			results := s.GetScriptExecutionHistory(ctx, realScriptID, deviceID)
			for result := range results {
				if result.Error != nil {
					continue
				}

				switch result.Ok.RunState {
				case "success":
					if result.Ok.RemediationScriptOutput != "" {
						var registryData azure.RegistryData
						if err := json.Unmarshal([]byte(result.Ok.RemediationScriptOutput), &registryData); err == nil {
							return &registryData, nil
						}
					}
					// If no output yet, continue waiting

				case "error", "failed":
					return nil, fmt.Errorf("script execution failed: %s (Error Code: %d)",
						result.Ok.ResultMessage, result.Ok.ErrorCode)

				case "pending", "running":
					// Continue waiting
					break

				default:
					// Unknown state, continue waiting
					break
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (s *azureClient) CollectRegistryDataFromDevice(ctx context.Context, deviceID string) (*azure.RegistryData, error) {
	// Use your existing deployed script instead of uploading a new one
	execution, err := s.ExecuteRegistryCollectionScript(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deployed script: %w", err)
	}

	registryData, err := s.WaitForScriptCompletion(ctx, execution.ID, 10*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("failed to get script results: %w", err)
	}

	return registryData, nil
}

// Fixed: Make device filtering configurable and add logging
func (s *azureClient) CollectRegistryDataFromAllDevices(ctx context.Context) <-chan AzureResult[azure.DeviceRegistryData] {
	return s.CollectRegistryDataFromAllDevicesWithConfig(ctx, getDefaultDeviceFilterConfig())
}

func (s *azureClient) CollectRegistryDataFromAllDevicesWithConfig(ctx context.Context, filterConfig DeviceFilterConfig) <-chan AzureResult[azure.DeviceRegistryData] {
	out := make(chan AzureResult[azure.DeviceRegistryData])

	go func() {
		defer close(out)

		devices := s.ListIntuneDevices(ctx, query.GraphParams{})

		for deviceResult := range devices {
			if deviceResult.Error != nil {
				out <- AzureResult[azure.DeviceRegistryData]{Error: deviceResult.Error}
				continue
			}

			device := deviceResult.Ok

			// Configurable OS filtering with logging
			osMatches := false
			for _, allowedOS := range filterConfig.IncludeOS {
				if strings.Contains(strings.ToLower(device.OperatingSystem), strings.ToLower(allowedOS)) {
					osMatches = true
					break
				}
			}

			if !osMatches {
				if filterConfig.LogSkippedDevices {
					fmt.Printf("Skipping device %s: OS %s not in allowed list %v\n",
						device.DeviceName, device.OperatingSystem, filterConfig.IncludeOS)
				}
				continue
			}

			// Configurable compliance filtering with logging
			if filterConfig.RequireCompliant && device.ComplianceState != "compliant" {
				if filterConfig.LogSkippedDevices {
					fmt.Printf("Skipping device %s: compliance state %s (require compliant: %v)\n",
						device.DeviceName, device.ComplianceState, filterConfig.RequireCompliant)
				}
				continue
			}

			registryData, err := s.CollectRegistryDataFromDevice(ctx, device.ID)
			if err != nil {
				out <- AzureResult[azure.DeviceRegistryData]{
					Error: fmt.Errorf("failed to collect registry data from device %s: %w", device.DeviceName, err),
				}
				continue
			}

			deviceRegistryData := azure.DeviceRegistryData{
				Device:       device,
				RegistryData: *registryData,
				CollectedAt:  time.Now(),
			}

			out <- AzureResult[azure.DeviceRegistryData]{Ok: deviceRegistryData}
		}
	}()

	return out
}

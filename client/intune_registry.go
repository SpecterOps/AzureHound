package client

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/sirupsen/logrus"
)

// Configuration for script deployment - now loaded from environment/config
var (
	registryScriptID   string
	registryScriptName string
	maxConcurrentJobs  int
	pollInterval       time.Duration
	maxWaitTime        time.Duration
)

// Initialize configuration from environment variables or defaults
func init() {
	registryScriptID = getEnvWithDefault("AZUREHOUND_REGISTRY_SCRIPT_ID", "BHE_Script_Registry_Data_Collection")
	registryScriptName = getEnvWithDefault("AZUREHOUND_REGISTRY_SCRIPT_NAME", "BHE_Script_Registry_Data_Collection.ps1")

	// Concurrency control - default 5 concurrent jobs
	if val := os.Getenv("AZUREHOUND_MAX_CONCURRENT_REGISTRY_JOBS"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			maxConcurrentJobs = int(parsed)
		}
	}
	if maxConcurrentJobs <= 0 {
		maxConcurrentJobs = 5
	}

	// Polling configuration
	if val := os.Getenv("AZUREHOUND_REGISTRY_POLL_INTERVAL"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			pollInterval = parsed
		}
	}
	if pollInterval <= 0 {
		pollInterval = 30 * time.Second
	}

	if val := os.Getenv("AZUREHOUND_REGISTRY_MAX_WAIT"); val != "" {
		if parsed, err := time.ParseDuration(val); err == nil {
			maxWaitTime = parsed
		}
	}
	if maxWaitTime <= 0 {
		maxWaitTime = 10 * time.Minute
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
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

// ExecuteRegistryCollectionScript executes the configured PowerShell script on an Intune device
func (s *azureClient) ExecuteRegistryCollectionScript(ctx context.Context, deviceID string) (*azure.ScriptExecution, error) {
	// First, get the deployed script ID
	scriptID, err := s.GetDeployedScriptID(ctx, registryScriptName)
	if err != nil {
		return nil, fmt.Errorf("failed to find deployed script: %w", err)
	}

	// Trigger script execution on the device
	err = s.TriggerScriptExecution(ctx, scriptID, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to trigger script execution: %w", err)
	}

	execution := &azure.ScriptExecution{
		ID:            createCompositeScriptID(scriptID, deviceID),
		DeviceID:      deviceID,
		Status:        "pending",
		StartDateTime: time.Now(),
		ScriptName:    registryScriptName,
		RunAsAccount:  "system",
	}

	return execution, nil
}

// createCompositeScriptID creates a safe composite ID with separator
func createCompositeScriptID(scriptID, deviceID string) string {
	return fmt.Sprintf("%s|%s", scriptID, deviceID)
}

// parseCompositeScriptID safely parses composite script ID
func parseCompositeScriptID(compositeID string) (scriptID, deviceID string, err error) {
	parts := strings.SplitN(compositeID, "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid composite script ID format: expected 'scriptID|deviceID', got '%s'", compositeID)
	}

	scriptID = strings.TrimSpace(parts[0])
	deviceID = strings.TrimSpace(parts[1])

	if scriptID == "" || deviceID == "" {
		return "", "", fmt.Errorf("invalid composite script ID: scriptID and deviceID cannot be empty")
	}

	return scriptID, deviceID, nil
}

// GetScriptExecutionResults retrieves results from script execution with improved error handling
func (s *azureClient) GetScriptExecutionResults(ctx context.Context, compositeScriptID string) <-chan AzureResult[azure.ScriptExecutionResult] {
	out := make(chan AzureResult[azure.ScriptExecutionResult])

	go func() {
		defer close(out)

		// Parse composite ID safely
		realScriptID, deviceID, err := parseCompositeScriptID(compositeScriptID)
		if err != nil {
			out <- AzureResult[azure.ScriptExecutionResult]{
				Error: fmt.Errorf("failed to parse script execution ID: %w", err),
			}
			return
		}

		// Query script execution results from Intune
		path := fmt.Sprintf("/%s/deviceManagement/deviceManagementScripts/%s/deviceRunStates",
			constants.GraphApiVersion, realScriptID)

		params := query.GraphParams{
			Filter: fmt.Sprintf("managedDevice/id eq '%s'", deviceID),
		}

		// Use the existing getAzureObjectList function
		go getAzureObjectList[azure.ScriptExecutionResult](s.msgraph, ctx, path, params, out)
	}()

	return out
}

// GetDeployedScriptID finds script ID by name
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

// TriggerScriptExecution triggers script on a specific device with improved error handling
func (s *azureClient) TriggerScriptExecution(ctx context.Context, scriptID, deviceID string) error {
	// Method 1: Use device management script assignment with Azure AD group
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
						// Note: Microsoft Graph requires group assignment, not individual device assignment
						"groupId": nil, // This should be set to an Azure AD group containing the device
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
		// Log the original assignment error for debugging
		logrus.WithError(err).Error("Script assignment via group failed, trying device action fallback")

		// If assignment method fails, try direct device action
		return s.triggerScriptViaDeviceAction(ctx, scriptID, deviceID)
	}

	return nil
}

// triggerScriptViaDeviceAction alternative method using device actions with error logging
func (s *azureClient) triggerScriptViaDeviceAction(ctx context.Context, scriptID, deviceID string) error {
	// Method 2: Use managed device executeAction
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
		logrus.WithError(err).Error("Failed to execute script via device action")
		return fmt.Errorf("failed to execute script via device action: %w", err)
	}

	return nil
}

// ValidateScriptDeployment checks if the script is properly deployed and accessible
func (s *azureClient) ValidateScriptDeployment(ctx context.Context) error {
	scriptID, err := s.GetDeployedScriptID(ctx, registryScriptName)
	if err != nil {
		return fmt.Errorf("script validation failed: %w", err)
	}

	if scriptID == "" {
		return fmt.Errorf("script ID is empty")
	}

	return nil
}

// WaitForScriptCompletion waits for script completion with configurable polling and exponential backoff
func (s *azureClient) WaitForScriptCompletion(ctx context.Context, compositeScriptID string, maxWaitTime time.Duration) (*azure.RegistryData, error) {
	timeout := time.After(maxWaitTime)
	currentInterval := pollInterval
	backoffMultiplier := 1.5
	maxInterval := 5 * time.Minute

	// Parse composite ID
	realScriptID, deviceID, err := parseCompositeScriptID(compositeScriptID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse script execution ID: %w", err)
	}

	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("script execution timed out after %v", maxWaitTime)
		case <-time.After(currentInterval):
			// Check execution status
			results := s.GetScriptExecutionHistory(ctx, realScriptID, deviceID)
			for result := range results {
				if result.Error != nil {
					logrus.WithError(result.Error).Warn("Error checking script execution status")
					continue
				}

				switch result.Ok.RunState {
				case "success":
					if result.Ok.RemediationScriptOutput != "" {
						var registryData azure.RegistryData
						if err := json.Unmarshal([]byte(result.Ok.RemediationScriptOutput), &registryData); err != nil {
							logrus.WithError(err).Error("Failed to unmarshal script output as JSON")
							return nil, fmt.Errorf("failed to parse script output: %w", err)
						}
						// Reset interval on success for future calls
						currentInterval = pollInterval
						return &registryData, nil
					}
					// If no output yet, continue waiting

				case "error", "failed":
					return nil, fmt.Errorf("script execution failed: %s (Error Code: %d)",
						result.Ok.ResultMessage, result.Ok.ErrorCode)

				case "pending", "running":
					// Continue waiting - implement exponential backoff
					currentInterval = time.Duration(float64(currentInterval) * backoffMultiplier)
					if currentInterval > maxInterval {
						currentInterval = maxInterval
					}

				default:
					// Unknown state, continue waiting with backoff
					currentInterval = time.Duration(float64(currentInterval) * backoffMultiplier)
					if currentInterval > maxInterval {
						currentInterval = maxInterval
					}
				}
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (s *azureClient) CollectRegistryDataFromDevice(ctx context.Context, deviceID string) (*azure.RegistryData, error) {
	// Use configured script instead of uploading a new one
	execution, err := s.ExecuteRegistryCollectionScript(ctx, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute deployed script: %w", err)
	}

	registryData, err := s.WaitForScriptCompletion(ctx, execution.ID, maxWaitTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get script results: %w", err)
	}

	return registryData, nil
}

// CollectRegistryDataFromAllDevices with concurrency control and improved device filtering
func (s *azureClient) CollectRegistryDataFromAllDevices(ctx context.Context) <-chan AzureResult[azure.DeviceRegistryData] {
	out := make(chan AzureResult[azure.DeviceRegistryData])

	go func() {
		defer close(out)

		devices := s.ListIntuneDevices(ctx, query.GraphParams{})

		// Create a semaphore for concurrency control
		semaphore := make(chan struct{}, maxConcurrentJobs)
		var wg sync.WaitGroup

		for deviceResult := range devices {
			if deviceResult.Error != nil {
				out <- AzureResult[azure.DeviceRegistryData]{Error: deviceResult.Error}
				continue
			}

			device := deviceResult.Ok

			// Only collect from Windows devices (removed compliance check)
			if !strings.Contains(strings.ToLower(device.OperatingSystem), "windows") {
				continue
			}

			// Acquire semaphore
			wg.Add(1)
			go func(dev azure.IntuneDevice) {
				defer wg.Done()
				semaphore <- struct{}{}        // Acquire
				defer func() { <-semaphore }() // Release

				registryData, err := s.CollectRegistryDataFromDevice(ctx, dev.ID)
				if err != nil {
					out <- AzureResult[azure.DeviceRegistryData]{
						Error: fmt.Errorf("failed to collect registry data from device %s: %w", dev.DeviceName, err),
					}
					return
				}

				deviceRegistryData := azure.DeviceRegistryData{
					Device:       dev,
					RegistryData: *registryData,
					CollectedAt:  time.Now(),
				}

				out <- AzureResult[azure.DeviceRegistryData]{Ok: deviceRegistryData}
			}(device)
		}

		// Wait for all goroutines to complete
		wg.Wait()
	}()

	return out
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

		go getAzureObjectList[azure.ScriptExecutionResult](s.msgraph, ctx, path, params, out)
	}()

	return out
}

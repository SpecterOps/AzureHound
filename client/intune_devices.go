// File: client/intune_devices.go
// Copyright (C) 2022 SpecterOps
// Implementation of Intune device management API calls

package client

import (
	"context"
	"fmt"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/intune"
)

func setDefaultParams(params query.GraphParams) query.GraphParams {
	if params.Top == 0 {
		params.Top = 999
	}
	return params
}

// ListIntuneManagedDevices retrieves all managed devices from Intune
// GET /deviceManagement/managedDevices
func (s *azureClient) ListIntuneManagedDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[intune.ManagedDevice] {
	var (
		out  = make(chan AzureResult[intune.ManagedDevice])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices", constants.GraphApiVersion)
	)

	params = setDefaultParams(params)

	go getAzureObjectList[intune.ManagedDevice](s.msgraph, ctx, path, params, out)
	return out
}

// GetIntuneDeviceCompliance retrieves compliance information for a specific device
// GET /deviceManagement/managedDevices/{id}/deviceCompliancePolicyStates
func (s *azureClient) GetIntuneDeviceCompliance(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ComplianceState] {
	if deviceId == "" {
		out := make(chan AzureResult[intune.ComplianceState])
		go func() {
			defer close(out)
			out <- AzureResult[intune.ComplianceState]{Error: fmt.Errorf("deviceId cannot be empty")}
		}()
		return out
	}

	var (
		out  = make(chan AzureResult[intune.ComplianceState])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices/%s/deviceCompliancePolicyStates", constants.GraphApiVersion, deviceId)
	)

	params = setDefaultParams(params)

	go getAzureObjectList[intune.ComplianceState](s.msgraph, ctx, path, params, out)
	return out
}

// GetIntuneDeviceConfiguration retrieves configuration information for a specific device
// GET /deviceManagement/managedDevices/{id}/deviceConfigurationStates
func (s *azureClient) GetIntuneDeviceConfiguration(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ConfigurationState] {
	if deviceId == "" {
		out := make(chan AzureResult[intune.ConfigurationState])
		go func() {
			defer close(out)
			out <- AzureResult[intune.ConfigurationState]{Error: fmt.Errorf("deviceId cannot be empty")}
		}()
		return out
	}

	var (
		out  = make(chan AzureResult[intune.ConfigurationState])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices/%s/deviceConfigurationStates", constants.GraphApiVersion, deviceId)
	)

	params = setDefaultParams(params)

	go getAzureObjectList[intune.ConfigurationState](s.msgraph, ctx, path, params, out)
	return out
}

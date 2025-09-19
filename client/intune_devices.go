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

// Maximum page size for Microsoft Graph API requests
const MaxGraphAPIPageSize = 999

func setDefaultParams(params *query.GraphParams) {
	if params.Top == 0 {
		params.Top = MaxGraphAPIPageSize
	}
}

// ListIntuneManagedDevices retrieves all managed devices from Intune
// GET /deviceManagement/managedDevices
func (s *azureClient) ListIntuneManagedDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[intune.ManagedDevice] {
	var (
		out  = make(chan AzureResult[intune.ManagedDevice])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices", constants.GraphApiVersion)
	)

	setDefaultParams(&params)

	go getAzureObjectList[intune.ManagedDevice](s.msgraph, ctx, path, params, out)
	return out
}

// GetIntuneDeviceCompliance retrieves compliance information for a specific device
// GET /deviceManagement/managedDevices/{id}/deviceCompliancePolicyStates
func (s *azureClient) GetIntuneDeviceCompliance(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ComplianceState] {
	var (
		out  = make(chan AzureResult[intune.ComplianceState])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices/%s/deviceCompliancePolicyStates", constants.GraphApiVersion, deviceId)
	)

	setDefaultParams(&params)

	go getAzureObjectList[intune.ComplianceState](s.msgraph, ctx, path, params, out)
	return out
}

// GetIntuneDeviceConfiguration retrieves configuration information for a specific device
// GET /deviceManagement/managedDevices/{id}/deviceConfigurationStates
func (s *azureClient) GetIntuneDeviceConfiguration(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ConfigurationState] {
	var (
		out  = make(chan AzureResult[intune.ConfigurationState])
		path = fmt.Sprintf("/%s/deviceManagement/managedDevices/%s/deviceConfigurationStates", constants.GraphApiVersion, deviceId)
	)

	setDefaultParams(&params)

	go getAzureObjectList[intune.ConfigurationState](s.msgraph, ctx, path, params, out)
	return out
}

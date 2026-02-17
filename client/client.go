// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package client

//go:generate go run go.uber.org/mock/mockgen -destination=./mocks/client.go -package=mocks . AzureClient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/config"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/client/rest"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/bloodhoundad/azurehound/v2/models/intune"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
)

func NewClient(config config.Config) (AzureClient, error) {
	if msgraph, err := rest.NewRestClient(config.GraphUrl(), config); err != nil {
		return nil, err
	} else if resourceManager, err := rest.NewRestClient(config.ResourceManagerUrl(), config); err != nil {
		return nil, err
	} else {
		if config.JWT != "" {
			if aud, err := rest.ParseAud(config.JWT); err != nil {
				return nil, err
			} else if aud == config.GraphUrl() {
				return initClientViaGraph(msgraph, resourceManager)
			} else if aud == config.ResourceManagerUrl() {
				if body, err := rest.ParseBody(config.JWT); err != nil {
					return nil, err
				} else {
					return initClientViaRM(msgraph, resourceManager, body["tid"])
				}
			} else {
				return nil, fmt.Errorf("error: invalid token audience")
			}
		} else {
			return initClientViaGraph(msgraph, resourceManager)
		}
	}
}

func initClientViaRM(msgraph, resourceManager rest.RestClient, tid interface{}) (AzureClient, error) {
	client := &azureClient{
		msgraph:         msgraph,
		resourceManager: resourceManager,
	}
	if result, err := client.GetAzureADTenants(context.Background(), true); err != nil {
		return nil, err
	} else {
		for _, tenant := range result.Value {
			if tenant.TenantId == tid.(string) {
				client.tenant = tenant
				break
			}
		}
		return client, nil
	}
}

func initClientViaGraph(msgraph, resourceManager rest.RestClient) (AzureClient, error) {
	client := &azureClient{
		msgraph:         msgraph,
		resourceManager: resourceManager,
	}
	if org, err := client.GetAzureADOrganization(context.Background(), nil); err != nil {
		return nil, err
	} else {
		client.tenant = org.ToTenant()
		return client, nil
	}
}

type AzureResult[T any] struct {
	Error error
	Ok    T
}

func getAzureObjectList[T any](client rest.RestClient, ctx context.Context, path string, params query.Params, out chan AzureResult[T]) {
	defer panicrecovery.PanicRecovery()
	defer close(out)

	var (
		errResult AzureResult[T]
		nextLink  string
	)

	for {
		var (
			list struct {
				CountGraph    int    `json:"@odata.count,omitempty"`    // The total count of all graph results
				NextLinkGraph string `json:"@odata.nextLink,omitempty"` // The URL to use for getting the next set of graph values.
				ContextGraph  string `json:"@odata.context,omitempty"`
				NextLinkRM    string `json:"nextLink,omitempty"` // The URL to use for getting the next set of rm values.
				Value         []T    `json:"value"`              // A list of azure values
			}
			res *http.Response
			err error
		)

		if nextLink != "" {
			if nextUrl, err := url.Parse(nextLink); err != nil {
				errResult.Error = err
				_ = pipeline.Send(ctx.Done(), out, errResult)
				return
			} else {
				paramsMap := make(map[string]string)
				if params != nil {
					paramsMap = params.AsMap()
				}
				if req, err := rest.NewRequest(ctx, "GET", nextUrl, nil, paramsMap, nil); err != nil {
					errResult.Error = err
					_ = pipeline.Send(ctx.Done(), out, errResult)
					return
				} else if res, err = client.Send(req); err != nil {
					errResult.Error = err
					_ = pipeline.Send(ctx.Done(), out, errResult)
					return
				}
			}
		} else {
			if res, err = client.Get(ctx, path, params, nil); err != nil {
				errResult.Error = err
				_ = pipeline.Send(ctx.Done(), out, errResult)
				return
			}
		}

		if err := rest.Decode(res.Body, &list); err != nil {
			errResult.Error = err
			_ = pipeline.Send(ctx.Done(), out, errResult)
			return
		} else {
			for _, u := range list.Value {
				if ok := pipeline.Send(ctx.Done(), out, AzureResult[T]{Ok: u}); !ok {
					return
				}
			}
		}

		if list.NextLinkRM == "" && list.NextLinkGraph == "" {
			break
		} else if list.NextLinkGraph != "" {
			nextLink = list.NextLinkGraph
		} else if list.NextLinkRM != "" {
			nextLink = list.NextLinkRM
		}
	}
}

type azureClient struct {
	msgraph         rest.RestClient
	resourceManager rest.RestClient
	tenant          azure.Tenant
}

// Core AzureGraphClient interface - unchanged to preserve compatibility
type AzureGraphClient interface {
	GetAzureADOrganization(ctx context.Context, selectCols []string) (*azure.Organization, error)

	ListAzureADGroups(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Group]
	ListAzureADGroupMembers(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage]
	ListAzureADGroupOwners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage]

	ListAzureADUsers(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.User]

	ListAzureADApps(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Application]
	ListAzureADAppOwners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage]

	ListAzureADServicePrincipals(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.ServicePrincipal]
	ListAzureADServicePrincipalOwners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage]

	ListAzureADRoles(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Role]
	ListAzureADRoleAssignments(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleAssignment]

	ListAzureDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.Device]
	ListAzureDeviceRegisteredOwners(ctx context.Context, objectId string, params query.GraphParams) <-chan AzureResult[json.RawMessage]

	ListAzureADAppRoleAssignments(ctx context.Context, servicePrincipalId string, params query.GraphParams) <-chan AzureResult[azure.AppRoleAssignment]
	ListAzureUnifiedRoleEligibilityScheduleInstances(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleEligibilityScheduleInstance]

	ListAzureADTenants(ctx context.Context, includeAllTenantCategories bool) <-chan AzureResult[azure.Tenant]
	GetAzureADTenants(ctx context.Context, includeAllTenantCategories bool) (azure.TenantList, error)

	ListRoleAssignmentPolicies(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.UnifiedRoleManagementPolicyAssignment]
}

type AzureResourceManagerClient interface {
	GetAzureADTenants(ctx context.Context, includeAllTenantCategories bool) (azure.TenantList, error)

	ListRoleAssignmentsForResource(ctx context.Context, resourceId string, filter, tenantId string) <-chan AzureResult[azure.RoleAssignment]
	ListAzureADTenants(ctx context.Context, includeAllTenantCategories bool) <-chan AzureResult[azure.Tenant]
	ListAzureContainerRegistries(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.ContainerRegistry]
	ListAzureWebApps(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.WebApp]
	ListAzureManagedClusters(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.ManagedCluster]
	ListAzureVMScaleSets(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.VMScaleSet]
	ListAzureKeyVaults(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.KeyVault]
	ListAzureManagementGroups(ctx context.Context, skipToken string) <-chan AzureResult[azure.ManagementGroup]
	ListAzureManagementGroupDescendants(ctx context.Context, groupId string, top int32) <-chan AzureResult[azure.DescendantInfo]
	ListAzureResourceGroups(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.ResourceGroup]
	ListAzureSubscriptions(ctx context.Context) <-chan AzureResult[azure.Subscription]
	ListAzureVirtualMachines(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.VirtualMachine]
	ListAzureStorageAccounts(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.StorageAccount]
	ListAzureStorageContainers(ctx context.Context, subscriptionId string, resourceGroupName string, saName string, filter string, includeDeleted string, maxPageSize string) <-chan AzureResult[azure.StorageContainer]
	ListAzureAutomationAccounts(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.AutomationAccount]
	ListAzureLogicApps(ctx context.Context, subscriptionId string, filter string, top int32) <-chan AzureResult[azure.LogicApp]
	ListAzureFunctionApps(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.FunctionApp]
}

// New interface for Intune-specific Graph operations
type IntuneGraphClient interface {
	ValidateScriptDeployment(ctx context.Context) error

	ListIntuneDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.IntuneDevice]
	ExecuteRegistryCollectionScript(ctx context.Context, deviceID string) (*azure.ScriptExecution, error)
	GetScriptExecutionResults(ctx context.Context, scriptID string) <-chan AzureResult[azure.ScriptExecutionResult]
	WaitForScriptCompletion(ctx context.Context, scriptID string, maxWaitTime time.Duration) (*azure.RegistryData, error)
	CollectRegistryDataFromDevice(ctx context.Context, deviceID string) (*azure.RegistryData, error)
	CollectRegistryDataFromAllDevices(ctx context.Context) <-chan AzureResult[azure.DeviceRegistryData]
	GetDeployedScriptID(ctx context.Context, scriptName string) (string, error)
	TriggerScriptExecution(ctx context.Context, scriptID, deviceID string) error

	// Add Intune methods
	ListIntuneManagedDevices(ctx context.Context, params query.GraphParams) <-chan AzureResult[intune.ManagedDevice]
	GetIntuneDeviceCompliance(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ComplianceState]
	GetIntuneDeviceConfiguration(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[intune.ConfigurationState]
}

// Composed interface that includes both core and Intune functionality
type ExtendedAzureGraphClient interface {
	AzureGraphClient
	IntuneGraphClient
}

// Core AzureClient interface - unchanged to preserve compatibility
type AzureClient interface {
	ExtendedAzureGraphClient

	ListAzureSubscriptions(ctx context.Context) <-chan AzureResult[azure.Subscription]
	ListAzureResourceGroups(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.ResourceGroup]
	ListAzureVirtualMachines(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.VirtualMachine]
	ListAzureVMScaleSets(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.VMScaleSet]
	ListAzureKeyVaults(ctx context.Context, subscriptionId string, params query.RMParams) <-chan AzureResult[azure.KeyVault]
	ListAzureStorageAccounts(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.StorageAccount]
	ListAzureStorageContainers(ctx context.Context, subscriptionId, resourceGroupName, saName, filter, includeDeleted, maxPageSize string) <-chan AzureResult[azure.StorageContainer]
	ListAzureContainerRegistries(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.ContainerRegistry]
	ListAzureWebApps(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.WebApp]
	ListAzureFunctionApps(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.FunctionApp]
	ListAzureLogicApps(ctx context.Context, subscriptionId, filter string, top int32) <-chan AzureResult[azure.LogicApp]
	ListAzureAutomationAccounts(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.AutomationAccount]
	ListAzureManagedClusters(ctx context.Context, subscriptionId string) <-chan AzureResult[azure.ManagedCluster]

	ListAzureManagementGroups(ctx context.Context, skipToken string) <-chan AzureResult[azure.ManagementGroup]
	ListAzureManagementGroupDescendants(ctx context.Context, groupId string, top int32) <-chan AzureResult[azure.DescendantInfo]

	ListRoleAssignmentsForResource(ctx context.Context, resourceId, filter, tenantId string) <-chan AzureResult[azure.RoleAssignment]

	TenantInfo() azure.Tenant
	CloseIdleConnections()
}

func (s azureClient) TenantInfo() azure.Tenant {
	return s.tenant
}

func (s azureClient) CloseIdleConnections() {
	s.msgraph.CloseIdleConnections()
	s.resourceManager.CloseIdleConnections()
}

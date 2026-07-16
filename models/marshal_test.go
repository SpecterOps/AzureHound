package models_test

import (
	"encoding/json"
	"testing"

	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
)

func marshalToMap(t *testing.T, v any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	return out
}

func TestUserMarshalJSONUppercasesIdentifiers(t *testing.T) {
	user := models.User{TenantId: "tenant-abc"}
	user.Id = "user-def"
	user.OnPremisesSecurityIdentifier = "s-1-5-21-abc-def"

	out := marshalToMap(t, user)

	require.Equal(t, "USER-DEF", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	// On-prem SID is the AZ->AD hybrid join key and must be uppercased.
	require.Equal(t, "S-1-5-21-ABC-DEF", out["onPremisesSecurityIdentifier"])
	// Source is unchanged.
	require.Equal(t, "user-def", user.Id)
	require.Equal(t, "tenant-abc", user.TenantId)
	require.Equal(t, "s-1-5-21-abc-def", user.OnPremisesSecurityIdentifier)
}

func TestGroupMarshalJSONUppercasesOnPremSID(t *testing.T) {
	group := models.Group{TenantId: "tenant-abc"}
	group.Id = "group-def"
	group.OnPremisesSecurityIdentifier = "s-1-5-21-ghi-jkl"

	out := marshalToMap(t, group)

	require.Equal(t, "GROUP-DEF", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "S-1-5-21-GHI-JKL", out["onPremisesSecurityIdentifier"])
	// Source is unchanged.
	require.Equal(t, "group-def", group.Id)
	require.Equal(t, "s-1-5-21-ghi-jkl", group.OnPremisesSecurityIdentifier)
}

func TestAppMarshalJSONUppercasesIdentifiers(t *testing.T) {
	app := models.App{TenantId: "tenant-abc"}
	app.Id = "app-def"
	app.AppId = "appid-ghi"

	out := marshalToMap(t, app)

	require.Equal(t, "APP-DEF", out["id"])
	require.Equal(t, "APPID-GHI", out["appId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	// Source is unchanged.
	require.Equal(t, "app-def", app.Id)
	require.Equal(t, "appid-ghi", app.AppId)
	require.Equal(t, "tenant-abc", app.TenantId)
}

func TestVirtualMachineMarshalJSONUppercasesIdentityNonMutating(t *testing.T) {
	vm := models.VirtualMachine{
		SubscriptionId:  "sub-1",
		ResourceGroupId: "/subscriptions/sub-1/resourcegroups/rg-1",
		TenantId:        "tenant-1",
	}
	vm.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/vm-1"
	vm.Identity.PrincipalId = "principal-sys"
	vm.Identity.UserAssignedIdentities = map[string]azure.UserAssignedIdentity{
		"/uai/one": {ClientId: "client-1", PrincipalId: "principal-uai"},
	}

	out := marshalToMap(t, vm)

	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1/PROVIDERS/VM-1", out["id"])
	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1", out["resourceGroupId"])
	require.Equal(t, "TENANT-1", out["tenantId"])

	identity := out["identity"].(map[string]any)
	require.Equal(t, "PRINCIPAL-SYS", identity["principalId"])
	uais := identity["userAssignedIdentities"].(map[string]any)
	uai := uais["/uai/one"].(map[string]any)
	require.Equal(t, "PRINCIPAL-UAI", uai["principalId"])
	require.Equal(t, "client-1", uai["clientId"])

	// Source identity map is unchanged.
	require.Equal(t, "principal-sys", vm.Identity.PrincipalId)
	require.Equal(t, "principal-uai", vm.Identity.UserAssignedIdentities["/uai/one"].PrincipalId)
}

func TestGroupOwnerMarshalJSONUppercasesOwnerId(t *testing.T) {
	owner := models.GroupOwner{
		GroupId: "group-1",
		Owner:   json.RawMessage(`{"id":"owner-1","@odata.type":"#microsoft.graph.user"}`),
	}

	out := marshalToMap(t, owner)

	require.Equal(t, "GROUP-1", out["groupId"])
	ownerBlob := out["owner"].(map[string]any)
	require.Equal(t, "OWNER-1", ownerBlob["id"])
}

func TestServicePrincipalOwnersMarshalJSONUppercasesServicePrincipalId(t *testing.T) {
	owners := models.ServicePrincipalOwners{
		ServicePrincipalId: "sp-1",
		Owners: []models.ServicePrincipalOwner{
			{
				ServicePrincipalId: "sp-1",
				Owner:              json.RawMessage(`{"id":"owner-1","@odata.type":"#microsoft.graph.user"}`),
			},
		},
	}

	out := marshalToMap(t, owners)

	// Top-level id is the edge endpoint and must match the AZServicePrincipal node ObjectID.
	require.Equal(t, "SP-1", out["servicePrincipalId"])
	entry := out["owners"].([]any)[0].(map[string]any)
	require.Equal(t, "SP-1", entry["servicePrincipalId"])
	ownerBlob := entry["owner"].(map[string]any)
	require.Equal(t, "OWNER-1", ownerBlob["id"])
}

func TestGroupOwnersMarshalJSONUppercasesGroupId(t *testing.T) {
	owners := models.GroupOwners{
		GroupId: "group-1",
		Owners: []models.GroupOwner{
			{
				GroupId: "group-1",
				Owner:   json.RawMessage(`{"id":"owner-1","@odata.type":"#microsoft.graph.user"}`),
			},
		},
	}

	out := marshalToMap(t, owners)

	require.Equal(t, "GROUP-1", out["groupId"])
	entry := out["owners"].([]any)[0].(map[string]any)
	require.Equal(t, "GROUP-1", entry["groupId"])
	ownerBlob := entry["owner"].(map[string]any)
	require.Equal(t, "OWNER-1", ownerBlob["id"])
}

func TestAppOwnersMarshalJSONUppercasesAppId(t *testing.T) {
	owners := models.AppOwners{
		AppId: "app-1",
		Owners: []models.AppOwner{
			{
				AppId: "app-1",
				Owner: json.RawMessage(`{"id":"owner-1","@odata.type":"#microsoft.graph.user"}`),
			},
		},
	}

	out := marshalToMap(t, owners)

	require.Equal(t, "APP-1", out["appId"])
	entry := out["owners"].([]any)[0].(map[string]any)
	require.Equal(t, "APP-1", entry["appId"])
	ownerBlob := entry["owner"].(map[string]any)
	require.Equal(t, "OWNER-1", ownerBlob["id"])
}

func TestKeyVaultOwnersMarshalJSONUppercasesScopeAndPrincipal(t *testing.T) {
	owners := models.KeyVaultOwners{
		KeyVaultId: "/subscriptions/s/kv-1",
		Owners: []models.KeyVaultOwner{
			{
				KeyVaultId: "/subscriptions/s/kv-1",
				Owner: azure.RoleAssignment{
					Properties: azure.RoleAssignmentPropertiesWithScope{
						PrincipalId: "principal-1",
						Scope:       "/subscriptions/s/kv-1",
					},
				},
			},
		},
	}

	out := marshalToMap(t, owners)

	require.Equal(t, "/SUBSCRIPTIONS/S/KV-1", out["keyVaultId"])
	entry := out["owners"].([]any)[0].(map[string]any)
	require.Equal(t, "/SUBSCRIPTIONS/S/KV-1", entry["keyVaultId"])
	ownerAssignment := entry["owner"].(map[string]any)
	props := ownerAssignment["properties"].(map[string]any)
	// Scope and target id are uppercased consistently so ingest == still holds.
	require.Equal(t, "/SUBSCRIPTIONS/S/KV-1", props["scope"])
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
}

func TestAppRoleAssignmentMarshalJSONUppercasesUUIDFields(t *testing.T) {
	principal := uuid.FromStringOrNil("aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee")
	appRole := uuid.FromStringOrNil("11111111-2222-3333-4444-555555555555")

	assignment := models.AppRoleAssignment{
		AppId:    "app-1",
		TenantId: "tenant-1",
	}
	assignment.PrincipalId = principal
	assignment.ResourceId = "resource-1"
	assignment.AppRoleId = appRole

	out := marshalToMap(t, assignment)

	require.Equal(t, "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE", out["principalId"])
	require.Equal(t, "RESOURCE-1", out["resourceId"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	// AppRoleId is used for lowercase matching in ingest and must remain untouched.
	require.Equal(t, "11111111-2222-3333-4444-555555555555", out["appRoleId"])
}

func TestAzureRoleAssignmentsMarshalJSONUppercasesEndpointsInSlice(t *testing.T) {
	assignments := models.AzureRoleAssignments{
		ObjectId: "/subscriptions/s/rg/cr-1",
		RoleAssignments: []models.AzureRoleAssignment{
			{
				ObjectId:         "/subscriptions/s/rg/cr-1",
				RoleDefinitionId: "b24988ac-6180-42a0-ab88-20f7382dd24c",
				Assignee: azure.RoleAssignment{
					Name: "af7c7710-3443-40e9-be55-f7d8eefba417",
					Properties: azure.RoleAssignmentPropertiesWithScope{
						PrincipalId:      "principal-1",
						Scope:            "/subscriptions/s/rg/cr-1",
						RoleDefinitionId: "b24988ac-6180-42a0-ab88-20f7382dd24c",
					},
				},
			},
		},
	}

	out := marshalToMap(t, assignments)

	entry := out["assignees"].([]any)[0].(map[string]any)
	require.Equal(t, "/SUBSCRIPTIONS/S/RG/CR-1", entry["objectId"])
	assignee := entry["assignee"].(map[string]any)
	props := assignee["properties"].(map[string]any)
	// Edge endpoints are uppercased so raw-path ingest matches node ObjectIDs.
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, "/SUBSCRIPTIONS/S/RG/CR-1", props["scope"])
	// RoleDefinitionId is matched against lowercase constants and must be preserved.
	require.Equal(t, "b24988ac-6180-42a0-ab88-20f7382dd24c", entry["roleDefinitionId"])
	require.Equal(t, "b24988ac-6180-42a0-ab88-20f7382dd24c", props["roleDefinitionId"])
	// Source is unchanged.
	require.Equal(t, "/subscriptions/s/rg/cr-1", assignments.RoleAssignments[0].ObjectId)
	require.Equal(t, "principal-1", assignments.RoleAssignments[0].Assignee.Properties.PrincipalId)
}

func TestSubscriptionRoleAssignmentMarshalJSONUppercasesEndpoints(t *testing.T) {
	ra := models.SubscriptionRoleAssignment{
		SubscriptionId: "sub-1",
		RoleAssignment: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId: "principal-1",
				Scope:       "/subscriptions/sub-1",
			},
		},
	}

	out := marshalToMap(t, ra)

	require.Equal(t, "SUB-1", out["subscriptionId"])
	props := out["roleAssignment"].(map[string]any)["properties"].(map[string]any)
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, "/SUBSCRIPTIONS/SUB-1", props["scope"])
	require.Equal(t, "sub-1", ra.SubscriptionId)
}

func TestResourceGroupRoleAssignmentMarshalJSONUppercasesEndpoints(t *testing.T) {
	ra := models.ResourceGroupRoleAssignment{
		ResourceGroupId: "/subscriptions/s/resourcegroups/rg-1",
		RoleAssignment: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId: "principal-1",
				Scope:       "/subscriptions/s/resourcegroups/rg-1",
			},
		},
	}

	out := marshalToMap(t, ra)

	require.Equal(t, "/SUBSCRIPTIONS/S/RESOURCEGROUPS/RG-1", out["resourceGroupId"])
	props := out["roleAssignment"].(map[string]any)["properties"].(map[string]any)
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, "/SUBSCRIPTIONS/S/RESOURCEGROUPS/RG-1", props["scope"])
}

func TestManagementGroupRoleAssignmentMarshalJSONUppercasesEndpoints(t *testing.T) {
	ra := models.ManagementGroupRoleAssignment{
		ManagementGroupId: "/providers/managementgroups/mg-1",
		RoleAssignment: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId: "principal-1",
				Scope:       "/providers/managementgroups/mg-1",
			},
		},
	}

	out := marshalToMap(t, ra)

	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-1", out["managementGroupId"])
	props := out["roleAssignment"].(map[string]any)["properties"].(map[string]any)
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-1", props["scope"])
}

func TestVirtualMachineRoleAssignmentMarshalJSONUppercasesEndpoints(t *testing.T) {
	ra := models.VirtualMachineRoleAssignment{
		VirtualMachineId: "/subscriptions/s/resourcegroups/rg/providers/vm-1",
		RoleAssignment: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId: "principal-1",
				Scope:       "/subscriptions/s/resourcegroups/rg/providers/vm-1",
			},
		},
	}

	out := marshalToMap(t, ra)

	require.Equal(t, "/SUBSCRIPTIONS/S/RESOURCEGROUPS/RG/PROVIDERS/VM-1", out["virtualMachineId"])
	props := out["roleAssignment"].(map[string]any)["properties"].(map[string]any)
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
}

func TestKeyVaultRoleAssignmentMarshalJSONUppercasesEndpoints(t *testing.T) {
	ra := models.KeyVaultRoleAssignment{
		KeyVaultId: "/subscriptions/s/kv-1",
		RoleAssignment: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId: "principal-1",
				Scope:       "/subscriptions/s/kv-1",
			},
		},
	}

	out := marshalToMap(t, ra)

	// KeyVaultId serializes under the legacy "virtualMachineId" json tag.
	require.Equal(t, "/SUBSCRIPTIONS/S/KV-1", out["virtualMachineId"])
	props := out["roleAssignment"].(map[string]any)["properties"].(map[string]any)
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, "/SUBSCRIPTIONS/S/KV-1", props["scope"])
}

func TestDescendantInfoMarshalJSONUppercasesIds(t *testing.T) {
	descendant := azure.DescendantInfo{
		Id: "/providers/managementgroups/mg-child",
	}
	descendant.Properties.Parent.Id = "/providers/managementgroups/mg-parent"

	out := marshalToMap(t, descendant)

	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-CHILD", out["id"])
	props := out["properties"].(map[string]any)
	parent := props["parent"].(map[string]any)
	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-PARENT", parent["id"])
}

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
	user := models.User{TenantId: "tenant-abc", TenantName: "contoso.onmicrosoft.com"}
	user.Id = "user-def"
	user.OnPremisesSecurityIdentifier = "s-1-5-21-abc-def"
	user.DisplayName = "Alice Example"

	out := marshalToMap(t, user)

	require.Equal(t, "USER-DEF", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO.ONMICROSOFT.COM", out["tenantName"])
	// On-prem SID is the AZ->AD hybrid join key and must be uppercased.
	require.Equal(t, "S-1-5-21-ABC-DEF", out["onPremisesSecurityIdentifier"])
	require.Equal(t, "ALICE EXAMPLE", out["displayName"])
	// Source is unchanged.
	require.Equal(t, "user-def", user.Id)
	require.Equal(t, "tenant-abc", user.TenantId)
	require.Equal(t, "s-1-5-21-abc-def", user.OnPremisesSecurityIdentifier)
	require.Equal(t, "Alice Example", user.DisplayName)
}

func TestGroupMarshalJSONUppercasesOnPremSID(t *testing.T) {
	group := models.Group{TenantId: "tenant-abc"}
	group.Id = "group-def"
	group.OnPremisesSecurityIdentifier = "s-1-5-21-ghi-jkl"
	group.DisplayName = "Engineering"

	out := marshalToMap(t, group)

	require.Equal(t, "GROUP-DEF", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "S-1-5-21-GHI-JKL", out["onPremisesSecurityIdentifier"])
	require.Equal(t, "ENGINEERING", out["displayName"])
	// Source is unchanged.
	require.Equal(t, "group-def", group.Id)
	require.Equal(t, "s-1-5-21-ghi-jkl", group.OnPremisesSecurityIdentifier)
	require.Equal(t, "Engineering", group.DisplayName)
}

func TestAppMarshalJSONUppercasesIdentifiers(t *testing.T) {
	app := models.App{TenantId: "tenant-abc", TenantName: "contoso.onmicrosoft.com"}
	app.Id = "app-def"
	app.AppId = "appid-ghi"
	app.DisplayName = "My App"

	out := marshalToMap(t, app)

	require.Equal(t, "APP-DEF", out["id"])
	require.Equal(t, "APPID-GHI", out["appId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO.ONMICROSOFT.COM", out["tenantName"])
	require.Equal(t, "MY APP", out["displayName"])
	// Source is unchanged.
	require.Equal(t, "app-def", app.Id)
	require.Equal(t, "appid-ghi", app.AppId)
	require.Equal(t, "tenant-abc", app.TenantId)
	require.Equal(t, "My App", app.DisplayName)
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

func TestGroupMembersMarshalJSONUppercasesGroupId(t *testing.T) {
	members := models.GroupMembers{
		GroupId: "group-1",
		Members: []models.GroupMember{
			{
				GroupId: "group-1",
				Member:  json.RawMessage(`{"id":"member-1","@odata.type":"#microsoft.graph.user"}`),
			},
		},
	}

	out := marshalToMap(t, members)

	// Top-level groupId is the edge endpoint and must match the AZGroup node ObjectID.
	require.Equal(t, "GROUP-1", out["groupId"])
	entry := out["members"].([]any)[0].(map[string]any)
	require.Equal(t, "GROUP-1", entry["groupId"])
	memberBlob := entry["member"].(map[string]any)
	require.Equal(t, "MEMBER-1", memberBlob["id"])
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

func TestAppFICMarshalJSONUppercasesAppIdAndFicId(t *testing.T) {
	fics := models.AppFICs{
		AppId:      "app-1",
		TenantId:   "tenant-1",
		TenantName: "SpecterOps Development",
		FICs: []models.AppFIC{
			{
				AppId: "app-1",
				FIC:   json.RawMessage(`{"id":"fic-1","issuer":"https://token.example/","subject":"repo:example:ref"}`),
			},
		},
	}

	out := marshalToMap(t, fics)

	// Top-level wrapper identifiers are uppercased, including tenantName.
	require.Equal(t, "APP-1", out["appId"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	require.Equal(t, "SPECTEROPS DEVELOPMENT", out["tenantName"])
	entry := out["fics"].([]any)[0].(map[string]any)
	// appId is the AZApp edge endpoint; fic.id is the FIC node ObjectID / source.
	require.Equal(t, "APP-1", entry["appId"])
	ficBlob := entry["fic"].(map[string]any)
	require.Equal(t, "FIC-1", ficBlob["id"])
	// Display-only fields are left untouched.
	require.Equal(t, "https://token.example/", ficBlob["issuer"])
	require.Equal(t, "repo:example:ref", ficBlob["subject"])
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
	require.Equal(t, "app-1", out["appId"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	// AppRoleId is used for lowercase matching in ingest and must remain untouched.
	require.Equal(t, "11111111-2222-3333-4444-555555555555", out["appRoleId"])
}

func TestAzureRoleAssignmentsMarshalJSONUppercasesEndpointsInSlice(t *testing.T) {
	// Use a mixed-case ARM resource id (as Azure returns it) to catch the raw
	// original casing, not just pure-lowercase.
	const mixedID = "/subscriptions/s/resourceGroups/BHE_RG/providers/Microsoft.ContainerRegistry/registries/specterDev"
	const upperID = "/SUBSCRIPTIONS/S/RESOURCEGROUPS/BHE_RG/PROVIDERS/MICROSOFT.CONTAINERREGISTRY/REGISTRIES/SPECTERDEV"

	assignments := models.AzureRoleAssignments{
		ObjectId: mixedID,
		RoleAssignments: []models.AzureRoleAssignment{
			{
				ObjectId:         mixedID,
				RoleDefinitionId: "b24988ac-6180-42a0-ab88-20f7382dd24c",
				Assignee: azure.RoleAssignment{
					Name: "af7c7710-3443-40e9-be55-f7d8eefba417",
					Properties: azure.RoleAssignmentPropertiesWithScope{
						PrincipalId:      "principal-1",
						Scope:            mixedID,
						RoleDefinitionId: "b24988ac-6180-42a0-ab88-20f7382dd24c",
					},
				},
			},
		},
	}

	out := marshalToMap(t, assignments)

	// Top-level objectId is the resource id BHE reads as the RBAC edge target
	// (data.ObjectId) for the resource-scoped role-assignment convertors; it must
	// be uppercased so the raw-path ingest does not create a mixed-case stub node.
	require.Equal(t, upperID, out["objectId"])

	entry := out["assignees"].([]any)[0].(map[string]any)
	require.Equal(t, upperID, entry["objectId"])
	assignee := entry["assignee"].(map[string]any)
	props := assignee["properties"].(map[string]any)
	// Edge endpoints are uppercased so raw-path ingest matches node ObjectIDs.
	require.Equal(t, "PRINCIPAL-1", props["principalId"])
	require.Equal(t, upperID, props["scope"])
	// RoleDefinitionId is matched against lowercase constants and must be preserved.
	require.Equal(t, "b24988ac-6180-42a0-ab88-20f7382dd24c", entry["roleDefinitionId"])
	require.Equal(t, "b24988ac-6180-42a0-ab88-20f7382dd24c", props["roleDefinitionId"])
	// Source is unchanged.
	require.Equal(t, mixedID, assignments.ObjectId)
	require.Equal(t, mixedID, assignments.RoleAssignments[0].ObjectId)
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
	descendant.Properties.DisplayName = "Child MG"

	out := marshalToMap(t, descendant)

	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-CHILD", out["id"])
	props := out["properties"].(map[string]any)
	parent := props["parent"].(map[string]any)
	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-PARENT", parent["id"])
	require.Equal(t, "CHILD MG", props["display_name"])
}

func TestServicePrincipalMarshalJSONUppercasesIdentifiers(t *testing.T) {
	sp := models.ServicePrincipal{TenantId: "tenant-abc", TenantName: "contoso.onmicrosoft.com"}
	sp.Id = "sp-def"
	sp.AppId = "appid-ghi"
	sp.AppOwnerOrganizationId = "owner-org-1"
	sp.DisplayName = "My SP"

	out := marshalToMap(t, sp)

	require.Equal(t, "SP-DEF", out["id"])
	require.Equal(t, "APPID-GHI", out["appId"])
	require.Equal(t, "OWNER-ORG-1", out["appOwnerOrganizationId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO.ONMICROSOFT.COM", out["tenantName"])
	require.Equal(t, "MY SP", out["displayName"])
	// Source is unchanged.
	require.Equal(t, "sp-def", sp.Id)
	require.Equal(t, "My SP", sp.DisplayName)
}

func TestSubscriptionMarshalJSONUppercasesDisplayName(t *testing.T) {
	sub := models.Subscription{TenantId: "tenant-abc"}
	sub.Id = "/subscriptions/sub-1"
	sub.SubscriptionId = "sub-guid-1"
	sub.DisplayName = "Prod Subscription"

	out := marshalToMap(t, sub)

	require.Equal(t, "/SUBSCRIPTIONS/SUB-1", out["id"])
	require.Equal(t, "SUB-GUID-1", out["subscriptionId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "PROD SUBSCRIPTION", out["displayName"])
	require.Equal(t, "Prod Subscription", sub.DisplayName)
}

func TestTenantMarshalJSONUppercasesDisplayName(t *testing.T) {
	tenant := models.Tenant{}
	tenant.Id = "/tenants/tenant-abc"
	tenant.TenantId = "tenant-abc"
	tenant.DisplayName = "Contoso"

	out := marshalToMap(t, tenant)

	require.Equal(t, "/TENANTS/TENANT-ABC", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO", out["displayName"])
	require.Equal(t, "Contoso", tenant.DisplayName)
}

func TestRoleMarshalJSONUppercasesDisplayName(t *testing.T) {
	role := models.Role{TenantId: "tenant-abc", TenantName: "contoso.onmicrosoft.com"}
	role.Id = "role-def"
	role.DisplayName = "Global Administrator"

	out := marshalToMap(t, role)

	require.Equal(t, "ROLE-DEF", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO.ONMICROSOFT.COM", out["tenantName"])
	require.Equal(t, "GLOBAL ADMINISTRATOR", out["displayName"])
	require.Equal(t, "Global Administrator", role.DisplayName)
}

func TestDeviceMarshalJSONUppercasesDisplayName(t *testing.T) {
	device := models.Device{TenantId: "tenant-abc", TenantName: "contoso.onmicrosoft.com"}
	device.Id = "device-def"
	device.DeviceId = "device-guid-1"
	device.DisplayName = "Laptop-01"

	out := marshalToMap(t, device)

	require.Equal(t, "DEVICE-DEF", out["id"])
	require.Equal(t, "DEVICE-GUID-1", out["deviceId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	require.Equal(t, "CONTOSO.ONMICROSOFT.COM", out["tenantName"])
	require.Equal(t, "LAPTOP-01", out["displayName"])
	require.Equal(t, "Laptop-01", device.DisplayName)
}

func TestManagementGroupMarshalJSONUppercasesDisplayName(t *testing.T) {
	mg := models.ManagementGroup{TenantId: "tenant-abc"}
	mg.Id = "/providers/managementgroups/mg-1"
	mg.Properties.DisplayName = "Root MG"
	mg.Properties.TenantId = "tenant-abc"

	out := marshalToMap(t, mg)

	require.Equal(t, "/PROVIDERS/MANAGEMENTGROUPS/MG-1", out["id"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	props := out["properties"].(map[string]any)
	require.Equal(t, "ROOT MG", props["displayName"])
	require.Equal(t, "TENANT-ABC", props["tenantId"])
	require.Equal(t, "Root MG", mg.Properties.DisplayName)
}

func TestKeyVaultMarshalJSONUppercasesPropertiesTenantId(t *testing.T) {
	kv := models.KeyVault{
		SubscriptionId: "sub-1",
		ResourceGroup:  "rg-1",
		TenantId:       "tenant-abc",
	}
	kv.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/kv-1"
	kv.Properties.TenantId = "tenant-abc"

	out := marshalToMap(t, kv)

	require.Equal(t, "TENANT-ABC", out["tenantId"])
	props := out["properties"].(map[string]any)
	require.Equal(t, "TENANT-ABC", props["tenantId"])
	require.Equal(t, "tenant-abc", kv.Properties.TenantId)
}

func TestUpperManagedIdentityUppercasesTenantId(t *testing.T) {
	vm := models.VirtualMachine{
		SubscriptionId:  "sub-1",
		ResourceGroupId: "/subscriptions/sub-1/resourcegroups/rg-1",
		TenantId:        "tenant-1",
	}
	vm.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/vm-1"
	vm.Identity.PrincipalId = "principal-sys"
	vm.Identity.TenantId = "tenant-1"

	out := marshalToMap(t, vm)

	identity := out["identity"].(map[string]any)
	require.Equal(t, "PRINCIPAL-SYS", identity["principalId"])
	require.Equal(t, "TENANT-1", identity["tenantId"])
	// Source is unchanged.
	require.Equal(t, "tenant-1", vm.Identity.TenantId)
}

func TestManagedClusterMarshalJSONUppercasesNodeResourceGroupAndIdentity(t *testing.T) {
	mc := models.ManagedCluster{
		SubscriptionId:  "sub-1",
		ResourceGroupId: "/subscriptions/sub-1/resourcegroups/rg-1",
		TenantId:        "tenant-1",
	}
	mc.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/mc-1"
	mc.Properties.NodeResourceGroup = "mc_rg-1_aks"
	mc.Identity.PrincipalId = "principal-sys"

	out := marshalToMap(t, mc)

	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1/PROVIDERS/MC-1", out["id"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	props := out["properties"].(map[string]any)
	require.Equal(t, "MC_RG-1_AKS", props["nodeResourceGroup"])
	identity := out["identity"].(map[string]any)
	require.Equal(t, "PRINCIPAL-SYS", identity["principalId"])
	// Source is unchanged.
	require.Equal(t, "mc_rg-1_aks", mc.Properties.NodeResourceGroup)
	require.Equal(t, "principal-sys", mc.Identity.PrincipalId)
}

func TestStorageAccountMarshalJSONUppercasesIdentifiersAndIdentity(t *testing.T) {
	sa := models.StorageAccount{
		SubscriptionId:    "sub-1",
		ResourceGroupId:   "/subscriptions/sub-1/resourcegroups/rg-1",
		ResourceGroupName: "rg-1",
		TenantId:          "tenant-1",
	}
	sa.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/sa-1"
	sa.Identity.PrincipalId = "principal-sys"

	out := marshalToMap(t, sa)

	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1/PROVIDERS/SA-1", out["id"])
	require.Equal(t, "SUB-1", out["subscriptionId"])
	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1", out["resourceGroupId"])
	require.Equal(t, "RG-1", out["resourceGroupName"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	identity := out["identity"].(map[string]any)
	require.Equal(t, "PRINCIPAL-SYS", identity["principalId"])
	require.Equal(t, "principal-sys", sa.Identity.PrincipalId)
}

func TestStorageContainerMarshalJSONUppercasesIdentifiers(t *testing.T) {
	sc := models.StorageContainer{
		SubscriptionId:    "sub-1",
		ResourceGroupId:   "/subscriptions/sub-1/resourcegroups/rg-1",
		ResourceGroupName: "rg-1",
		StorageAccountId:  "/subscriptions/sub-1/resourcegroups/rg-1/providers/sa-1",
		TenantId:          "tenant-1",
	}
	sc.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/sa-1/blobservices/default/containers/c-1"

	out := marshalToMap(t, sc)

	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1/PROVIDERS/SA-1/BLOBSERVICES/DEFAULT/CONTAINERS/C-1", out["id"])
	require.Equal(t, "SUB-1", out["subscriptionId"])
	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1", out["resourceGroupId"])
	require.Equal(t, "RG-1", out["resourceGroupName"])
	require.Equal(t, "/SUBSCRIPTIONS/SUB-1/RESOURCEGROUPS/RG-1/PROVIDERS/SA-1", out["storageAccountId"])
	require.Equal(t, "TENANT-1", out["tenantId"])
	// Source is unchanged.
	require.Equal(t, "sub-1", sc.SubscriptionId)
}

func TestRoleAssignmentsMarshalJSONUppercasesInnerEndpoints(t *testing.T) {
	ra := models.RoleAssignments{
		RoleDefinitionId: "role-def-1",
		TenantId:         "tenant-abc",
		RoleAssignments: []azure.UnifiedRoleAssignment{
			{
				RoleDefinitionId: "role-def-1",
				PrincipalId:      "principal-1",
				DirectoryScopeId: "/administrativeUnits/au-1",
			},
		},
	}
	// The assignment id is a case-sensitive base64url Graph identifier.
	ra.RoleAssignments[0].Id = "lAPPGGoDgkmB_Xyz-123"

	out := marshalToMap(t, ra)

	require.Equal(t, "ROLE-DEF-1", out["roleDefinitionId"])
	require.Equal(t, "TENANT-ABC", out["tenantId"])
	entry := out["roleAssignments"].([]any)[0].(map[string]any)
	// The per-assignment id is base64url and must be preserved as-is.
	require.Equal(t, "lAPPGGoDgkmB_Xyz-123", entry["id"])
	require.Equal(t, "ROLE-DEF-1", entry["roleDefinitionId"])
	require.Equal(t, "PRINCIPAL-1", entry["principalId"])
	require.Equal(t, "/ADMINISTRATIVEUNITS/AU-1", entry["directoryScopeId"])
	// Source is unchanged.
	require.Equal(t, "role-def-1", ra.RoleAssignments[0].RoleDefinitionId)
	require.Equal(t, "principal-1", ra.RoleAssignments[0].PrincipalId)
}

func TestKeyVaultMarshalJSONUppercasesAccessPolicies(t *testing.T) {
	kv := models.KeyVault{
		SubscriptionId: "sub-1",
		ResourceGroup:  "rg-1",
		TenantId:       "tenant-abc",
	}
	kv.Id = "/subscriptions/sub-1/resourcegroups/rg-1/providers/kv-1"
	kv.Properties.AccessPolicies = []azure.AccessPolicyEntry{
		{
			ObjectId:      "object-1",
			ApplicationId: "app-1",
			TenantId:      "tenant-abc",
		},
	}

	out := marshalToMap(t, kv)

	props := out["properties"].(map[string]any)
	policy := props["accessPolicies"].([]any)[0].(map[string]any)
	require.Equal(t, "OBJECT-1", policy["objectId"])
	require.Equal(t, "APP-1", policy["applicationId"])
	require.Equal(t, "TENANT-ABC", policy["tenantId"])
	// Source is unchanged.
	require.Equal(t, "object-1", kv.Properties.AccessPolicies[0].ObjectId)
	require.Equal(t, "app-1", kv.Properties.AccessPolicies[0].ApplicationId)
}

func TestRoleManagementPolicyAssignmentMarshalJSONUppercasesApproverIds(t *testing.T) {
	rule := `{
		"@odata.type": "#microsoft.graph.unifiedRoleManagementPolicyApprovalRule",
		"setting": {
			"approvalStages": [
				{
					"primaryApprovers": [
						{"@odata.type": "#microsoft.graph.groupMembers", "groupId": "group-1"},
						{"@odata.type": "#microsoft.graph.singleUser", "userId": "user-1"}
					]
				}
			]
		}
	}`
	rmpa := models.RoleManagementPolicyAssignment{
		Id:                              "policy-assignment-1",
		RoleDefinitionId:                "role-def-1",
		TenantId:                        "tenant-abc",
		EndUserAssignmentGroupApprovers: []string{"group-1"},
	}
	rmpa.Policy.Rules = []json.RawMessage{json.RawMessage(rule)}

	out := marshalToMap(t, rmpa)

	require.Equal(t, "POLICY-ASSIGNMENT-1", out["id"])
	require.Equal(t, "ROLE-DEF-1", out["roleDefinitionId"])
	require.Equal(t, "GROUP-1", out["endUserAssignmentGroupApprovers"].([]any)[0])

	// Approver ids nested in the raw policy rules are uppercased.
	rules := out["policy"].(map[string]any)["rules"].([]any)
	setting := rules[0].(map[string]any)["setting"].(map[string]any)
	stage := setting["approvalStages"].([]any)[0].(map[string]any)
	approvers := stage["primaryApprovers"].([]any)
	require.Equal(t, "GROUP-1", approvers[0].(map[string]any)["groupId"])
	require.Equal(t, "USER-1", approvers[1].(map[string]any)["userId"])
	// Source is unchanged.
	require.Contains(t, string(rmpa.Policy.Rules[0]), "group-1")
}

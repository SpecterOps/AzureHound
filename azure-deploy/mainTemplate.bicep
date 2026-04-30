// mainTemplate.bicep
targetScope = 'subscription'

@description('Name of the resource group to deploy into')
param resourceGroupName string

@description('Location for all resources')
param location string

@description('Name of the container user managed identity')
param containerUMIName string

@description('Resource group for the container UMI (defaults to main resource group if not specified)')
param containerUMIResourceGroupName string = resourceGroupName

@description('Name of the deployment user managed identity')
param deploymentUMIName string

@description('Resource group containing the deployment user managed identity')
param deploymentUMIResourceGroupName string

@description('Azure tenant ID to analyze')
param azureTenantId string

@description('Bloodhound instance domain')
param bloodhoundInstanceDomain string

@description('Bloodhound token ID')
@secure()
param bloodhoundTokenId string

@description('Bloodhound token')
@secure()
param bloodhoundToken string

var containerName = '${resourceGroupName}-container-group'
var imageName = 'azurehounddeploy.azurecr.io/azurehound:latest'

// Create main resource group if it doesn't exist
resource mainRG 'Microsoft.Resources/resourceGroups@2023-07-01' = {
  name: resourceGroupName
  location: location
}

// Create container UMI resource group if different from main and doesn't exist
resource containerUMIRG 'Microsoft.Resources/resourceGroups@2023-07-01' = if (containerUMIResourceGroupName != resourceGroupName) {
  name: containerUMIResourceGroupName
  location: location
}

// Create deployment UMI resource group if different from main
resource deploymentUMIRG 'Microsoft.Resources/resourceGroups@2023-07-01' = if (deploymentUMIResourceGroupName != resourceGroupName) {
  name: deploymentUMIResourceGroupName
  location: location
}


// Deploy or reference the container UMI
module containerIdentity './modules/containerIdentity.bicep' = {
  name: 'containerIdentity-deployment'
  scope: resourceGroup(containerUMIResourceGroupName)
  params: {
    location: location
    containerUMIName: containerUMIName
  }
  dependsOn: [
    mainRG
    containerUMIRG
  ]
}

// Call module to create deployment UMI in correct RG
module deploymentIdentity 'modules/deploymentIdentity.bicep' = {
  name: 'deploymentIdentity-deployment'
  scope: resourceGroup(deploymentUMIResourceGroupName)
  params: {
    location: location
    deploymentUMIName: deploymentUMIName
  }
  dependsOn: [
    mainRG
    deploymentUMIRG
  ]
}

// Assign Reader role at subscription scope
resource readerRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(subscription().subscriptionId, containerUMIName, 'Reader')
  scope: subscription()
  properties: {
    principalId: containerIdentity.outputs.principalId
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'acdd72a7-3385-48ef-bd42-f606fba81ae7') // Reader role
    principalType: 'ServicePrincipal'
  }
  dependsOn: [
    containerIdentity
  ]
}

// Deploy script to configure Graph API permissions
module graphPermissions './modules/graphPermissions.bicep' = {
  name: 'graph-permissions-deployment'
  scope: resourceGroup(resourceGroupName)
  params: {
    location: location
    deploymentUMIName: deploymentUMIName
    deploymentUMIResourceGroupName: deploymentUMIResourceGroupName
    containerUMIPrincipalId: containerIdentity.outputs.principalId
  }
  dependsOn: [
    mainRG
    containerIdentity
  ]
}

// Deploy the container instance
module containerInstance './modules/containerInstance.bicep' = {
  name: 'container-instance-deployment'
  scope: resourceGroup(resourceGroupName)
  params: {
    location: location
    containerGroupName: containerName
    containerUMIResourceId: containerIdentity.outputs.resourceId
    imageName: imageName
    bloodhoundInstanceDomain: bloodhoundInstanceDomain
    azureTenantId: azureTenantId
    bloodhoundTokenId: bloodhoundTokenId
    bloodhoundToken: bloodhoundToken
  }
  dependsOn: [
    mainRG
    graphPermissions
  ]
}

output containerUMIResourceId string = containerIdentity.outputs.resourceId
output containerUMIPrincipalId string = containerIdentity.outputs.principalId
output permissionSetupRequired bool = graphPermissions.outputs.needsManualSetup
output permissionStatus string = graphPermissions.outputs.statusMessage
output assignedPermissions array = graphPermissions.outputs.assignedPermissions
output deploymentUMIPrincipalId string = deploymentIdentity.outputs.principalId
output deploymentUMIId string = deploymentIdentity.outputs.resourceId
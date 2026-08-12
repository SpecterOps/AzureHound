// modules/containerIdentity.bicep
param location string
param containerUMIName string

// Reference existing identity if it exists, create if it doesn't
resource containerUMI 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: containerUMIName
  location: location
}

// Check if Reader role is already assigned
resource existingReaderRole 'Microsoft.Authorization/roleAssignments@2022-04-01' existing = {
  scope: subscription()
  name: guid(subscription().id, containerUMI.id, 'Reader')
}

// Assign Reader role at subscription scope if not already assigned
resource readerRoleAssignment 'Microsoft.Authorization/roleAssignments@2022-04-01' = {
  name: guid(subscription().id, containerUMI.id, 'Reader')
  properties: {
    principalId: containerUMI.properties.principalId
    roleDefinitionId: subscriptionResourceId('Microsoft.Authorization/roleDefinitions', 'acdd72a7-3385-48ef-bd42-f606fba81ae7') // Reader role
    principalType: 'ServicePrincipal'
  }
}

output resourceId string = containerUMI.id
output principalId string = containerUMI.properties.principalId

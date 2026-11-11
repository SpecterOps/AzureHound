// modules/deploymentIdentity.bicep
param location string
param deploymentUMIName string

resource deploymentUMI 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' = {
  name: deploymentUMIName
  location: location
}

output resourceId string = deploymentUMI.id
output principalId string = deploymentUMI.properties.principalId

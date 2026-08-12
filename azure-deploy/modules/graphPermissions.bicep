// modules/graphPermissions.bicep
param location string
param containerUMIPrincipalId string
param deploymentUMIResourceGroupName string
param deploymentUMIName string

// Reference the deployment UMI to ensure it exists before using it
resource deploymentUMI 'Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31' existing = {
  name: deploymentUMIName
  scope: resourceGroup(deploymentUMIResourceGroupName)
}

resource graphPermissionScript 'Microsoft.Resources/deploymentScripts@2023-08-01' = {
  name: 'graph-permissions-script'
  location: location
  kind: 'AzurePowerShell'
  identity: {
    type: 'UserAssigned'
    userAssignedIdentities: {
      '${deploymentUMI.id}': {}
    }
  }
  properties: {
    azPowerShellVersion: '9.7'
    retentionInterval: 'P1D'
    timeout: 'PT30M'
    cleanupPreference: 'Always'
    scriptContent: '''
      $ErrorActionPreference = "Continue"
      
      # Initialize arrays for tracking
      $warningsList = @()
      $successList = @()
      $needManualSetup = $false
      
      try {
          $token = (Get-AzAccessToken -ResourceUrl "https://graph.microsoft.com/").Token
          $graphAppId = "00000003-0000-0000-c000-000000000000"
          
          $graphSp = Get-AzADServicePrincipal -ApplicationId $graphAppId
          if (-not $graphSp) {
              $needManualSetup = $true
          }
          
          if ($graphSp) {
              $headers = @{
                  'Authorization' = "Bearer $token"
                  'Content-Type' = 'application/json'
              }

              try {
                  $apiUrl = "https://graph.microsoft.com/v1.0/servicePrincipals/$($env:ContainerUMIPrincipalId)/appRoleAssignments"
                  $existingAssignments = Invoke-RestMethod -Uri $apiUrl -Headers $headers -Method Get
                  
                  $directoryReadAllId = "7ab1d382-f21e-4acd-a863-ba3e13f7da61"
                  $existingAssignment = $existingAssignments.value | Where-Object { 
                      $_.appRoleId -eq $directoryReadAllId -and 
                      $_.resourceId -eq $graphSp.Id
                  }

                  if (-not $existingAssignment) {
                      try {
                          $body = @{
                              principalId = $env:ContainerUMIPrincipalId
                              resourceId = $graphSp.Id
                              appRoleId = $directoryReadAllId
                          } | ConvertTo-Json

                          $result = Invoke-RestMethod -Uri $apiUrl -Headers $headers -Method Post -Body $body
                          $successList += "Directory.Read.All"
                      }
                      catch {
                          $needManualSetup = $true
                      }
                  }
                  else {
                      $successList += "Directory.Read.All"
                  }
              }
              catch {
                  $needManualSetup = $true
              }
          }
      }
      catch {
          $needManualSetup = $true
      }

      # Create a structured permission status message
      $statusMessage = if ($needManualSetup) {
          @"
MANUAL PERMISSION SETUP REQUIRED
------------------------------
The container's managed identity requires the following Microsoft Graph permission:
- Directory.Read.All

Please run the provided setup script to configure these permissions:
./setup-container-permissions.ps1 -PrincipalId $($env:ContainerUMIPrincipalId)
"@
      } else {
          "All required permissions have been configured successfully."
      }
      
      # Output the results
      $DeploymentScriptOutputs = @{
          needsManualSetup = $needManualSetup
          statusMessage = $statusMessage
          assignedPermissions = $successList
      }
    '''
    environmentVariables: [
      {
        name: 'ContainerUMIPrincipalId'
        value: containerUMIPrincipalId
      }
    ]
  }
}

// Output these values so they can be captured by the main template
output needsManualSetup bool = graphPermissionScript.properties.outputs.needsManualSetup
output statusMessage string = graphPermissionScript.properties.outputs.statusMessage
output assignedPermissions array = graphPermissionScript.properties.outputs.assignedPermissions

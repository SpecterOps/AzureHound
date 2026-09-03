param(
    [Parameter(Mandatory=$true)]
    [string]$ResourceGroupName,
    [Parameter(Mandatory=$true)]
    [string]$Location,
    [Parameter(Mandatory=$true)]
    [string]$ContainerUMIName,
    [Parameter(Mandatory=$true)]
    [string]$AzureTenantId,
    [Parameter(Mandatory=$true)]
    [string]$BloodhoundInstanceDomain,
    [Parameter(Mandatory=$true)]
    [string]$BloodhoundTokenId,
    [Parameter(Mandatory=$true)]
    [string]$BloodhoundToken,
    [Parameter(Mandatory=$true)]
    [string]$RegistryPassword
)

# Ensure we're connected to Azure
if (-not (Get-AzContext)) {
    Write-Error "Not connected to Azure. Please run Connect-AzAccount first."
    exit 1
}

function New-ManagedIdentity {
    param(
        [Parameter(Mandatory=$true)]
        [PSObject]$ResourceGroup,
        [Parameter(Mandatory=$true)]
        [string]$IdentityName
    )

    $identity = New-AzUserAssignedIdentity -ResourceGroupName $ResourceGroup.ResourceGroupName -Name $IdentityName -Location $ResourceGroup.Location
    if (-not $identity) {
        throw "Failed to create managed identity"
    }

    # Wait for identity to propagate to AAD
    Write-Host "Waiting for identity $IdentityName to propagate to Azure AD..."
    $timeout = (Get-Date).AddMinutes(2)
    $servicePrincipal = $null
    
    while ((Get-Date) -lt $timeout) {
        $servicePrincipal = Get-AzADServicePrincipal -ObjectId $identity.PrincipalId -ErrorAction SilentlyContinue
        if ($servicePrincipal) {
            Write-Host "Identity propagation confirmed"
            break
        }
        Write-Host "Waiting for identity propagation..."
        Start-Sleep -Seconds 10
    }

    if (-not $servicePrincipal) {
        throw "Timeout waiting for identity to propagate to Azure AD"
    }

    return $identity
}


function Add-RoleAssignment {
    param(
        [Parameter(Mandatory=$true)]
        [string]$PrincipalId,
        [Parameter(Mandatory=$true)]
        [string]$RoleDefinitionName,
        [Parameter(Mandatory=$true)]
        [string]$Scope
    )

    New-AzRoleAssignment -ObjectId $PrincipalId `
                        -RoleDefinitionName $RoleDefinitionName `
                        -Scope $Scope
}

function Test-RoleAssignment {
    param(
        [string]$PrincipalId,
        [string]$RoleDefinitionName,
        [string]$Scope
    )
    
    $assignment = Get-AzRoleAssignment `
        -ObjectId $PrincipalId `
        -RoleDefinitionName $RoleDefinitionName `
        -Scope $Scope `
        -ErrorAction SilentlyContinue

    return $null -ne $assignment
}

#
#Add-RoleAssignment `
# -PrincipalId $containerUMI.PrincipalId `
# -RoleDefinitionName "Reader" `
# -Scope $subscriptionScope
#

function Add-GraphApiPermissionWithPropogationTest {
    param(
        [Parameter(Mandatory=$true)]
        [string]$PrincipalId,
        [Parameter(Mandatory=$true)]
        [string]$PermissionName,
        [int]$TimeoutInMinutes = 2
    )

    $token = (Get-AzAccessToken -ResourceUrl "https://graph.microsoft.com/").Token
    $graphAppId = "00000003-0000-0000-c000-000000000000"
    $graphSp = Get-AzADServicePrincipal -ApplicationId $graphAppId

    # Define permission IDs
    $permissionIds = @{
        "Directory.Read.All" = "7ab1d382-f21e-4acd-a863-ba3e13f7da61"
        "User.Read" = "e1fe6dd8-ba31-4d61-89e7-88639da4683d"
    }

    if (-not $permissionIds.ContainsKey($PermissionName)) {
        throw "Unknown permission: $PermissionName. Supported permissions are: $($permissionIds.Keys -join ', ')"
    }

    $headers = @{
        'Authorization' = "Bearer $token"
        'Content-Type' = 'application/json'
    }

    # Function to check if permission exists
    function Test-PermissionAssignment {
        $existingAssignments = Invoke-RestMethod -Uri "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId/appRoleAssignments" -Headers $headers -Method Get
        return $existingAssignments.value | Where-Object { $_.appRoleId -eq $permissionIds[$PermissionName] }
    }

    # Check if permission is already assigned
    $existingAssignment = Test-PermissionAssignment
    if ($existingAssignment) {
        Write-Host "Permission $PermissionName is already assigned"
        return
    }

    $body = @{
        principalId = $PrincipalId
        resourceId = $graphSp.Id
        appRoleId = $permissionIds[$PermissionName]
    } | ConvertTo-Json

    try {
        $apiUrl = "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId/appRoleAssignments"
        $result = Invoke-RestMethod -Uri $apiUrl -Headers $headers -Method Post -Body $body
        Write-Host "Permission assignment initiated for $PermissionName"

        # Wait for permission to propagate
        Write-Host "Waiting for permission $PermissionName to propagate..."
        $timeout = (Get-Date).AddMinutes($TimeoutInMinutes)
        $permissionConfirmed = $false
        
        while ((Get-Date) -lt $timeout) {
            if (Test-PermissionAssignment) {
                $permissionConfirmed = $true
                Write-Host "Permission $PermissionName successfully propagated"
                break
            }
            Write-Host "Waiting for permission to propagate..."
            Start-Sleep -Seconds 10
        }

        if (-not $permissionConfirmed) {
            throw "Timeout waiting for permission $PermissionName to propagate"
        }
    }
    catch {
        $errorMessage = $_.ErrorDetails.Message | ConvertFrom-Json -ErrorAction SilentlyContinue
        Write-Warning "Failed to assign Graph API permission: $PermissionName"
        Write-Warning "Error: $($errorMessage.error.message)"
        throw
    }
}


# Create User Assigned Managed Identity if it doesn't exist
function New-AzureHoundContainerGroup {
    param(
        [Parameter(Mandatory=$true)]
        [PSObject]$ResourceGroup,
        [Parameter(Mandatory=$true)]
        [string]$ContainerGroupName,
        [Parameter(Mandatory=$true)]
        [string]$ContainerInstanceName,
        [Parameter(Mandatory=$true)]
        [PSObject]$ContainerRegistry,
        [Parameter(Mandatory=$true)]
        [string]$ContainerImage,
        [Parameter(Mandatory=$true)]
        [string[]]$ContainerEntrypoint,
        [Parameter(Mandatory=$true)]
        [string]$ConfigJson,
        [Parameter(Mandatory=$true)]
        [PSObject]$ContainerUMI
    )

    # Convert config.json to base64
    $configJsonBase64 = [Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes($configJson))

    # Create container instance
    $containerInstance = New-AzContainerInstanceObject `
        -Name $ContainerGroupName `
        -Image $ContainerImage `
        -VolumeMount @(
            @{
                Name = "config-volume"
                MountPath = "/home/nonroot/.config/azurehound"
            }
        ) `
        -Command $ContainerEntrypoint


    # Create container group
    $containerGroup = New-AzContainerGroup `
        -ResourceGroupName $ResourceGroup.ResourceGroupName `
        -Name $ContainerGroupName `
        -Location $ResourceGroup.Location `
        -IdentityType "UserAssigned" `
        -IdentityUserAssignedIdentity @{ 
            $ContainerUMI.Id = @{} 
        } `
        -Volume @{
                Name = "config-volume"
                Secret = @{
                  "config.json" = $configJsonBase64
            }
        } `
        -ImageRegistryCredential @(
            @{
                Server = $ContainerRegistry.LoginServer
                Username = "ditkinreg"
                Password = $RegistryPassword
            }
        ) `
        -Container $containerInstance `
        -OsType Linux `
        -RestartPolicy Never `
        
    
    return $containerGroup
}

try {
    # Create Resource Group if it doesn't exist assign existing or new resource group to variable rg
    if (-not ($resourceGroup = Get-AzResourceGroup -Name $ResourceGroupName -ErrorAction SilentlyContinue)) {
        Write-Host "Creating resource group: $ResourceGroupName"
        $resourceGroup = New-AzResourceGroup -Name $ResourceGroupName -Location $Location
    } else {
        Write-Host "Resource group $ResourceGroupName already exists"
    }

    # Identity name derived from the resource group name
    if (-not $ContainerUMIName) {
        $containerUMIName = "$ResourceGroupName-container-umi"
        Write-Host "Container UMI name not provided, using default: $containerUMIName"
    } else {
        Write-Host "Using provided container UMI name: $ContainerUMIName"
        $containerUMIName = $ContainerUMIName
    }
    $containerName = "$ResourceGroupName-container-group"

    # Get the container registry
    $containerRegistryName = "ditkinreg"
    $containerRegistryResourceGroup = "ditkin-test-registry"

    # TODO: This should be parameterized
    $imageName = "ditkin-test-image:latest"

    # Authenticate with ContainerRegistry
    Connect-AzContainerRegistry -Name $containerRegistryName 

    # Authenticate with ContainerRegistry
    $acr = Get-AzContainerRegistry -ResourceGroupName $containerRegistryResourceGroup -Name $containerRegistryName

    if (-not ($containerUMI = Get-AzUserAssignedIdentity -ResourceGroupName $resourceGroup.ResourceGroupName -Name $containerUMIName -ErrorAction SilentlyContinue)){
        Write-Host "Creating managed identity: $containerUMIName"
        $containerUMI = New-ManagedIdentity -ResourceGroup $resourceGroup -IdentityName $containerUMIName
    } else {
        Write-Host "Managed identity $containerUMIName already exists"
    }

    # json config for azurehound
    $config = @{
        app = "appValue"
        auth = ""
        batchsize = 100
        config = "/home/nonroot/.config/azurehound/config.json"
        instance = "https://${BloodhoundInstanceDomain}/"
        json = $false
        'managed-identity' = $true
        maxconnsperhost = 20
        maxidleconnsperhost = 20
        region = "cloud"
        streamcount = 25
        tenant = "${AzureTenantId}"
        token = "${BloodhoundToken}"
        tokenid = "${BloodhoundTokenId}"
        verbosity = 0
    }

    $configJson = $config | ConvertTo-Json

    if (-not ($containerGroup = Get-AzContainerGroup -ResourceGroupName $ResourceGroupName -Name $containerName -ErrorAction SilentlyContinue)) {
        $containerGroup = New-AzureHoundContainerGroup `
            -ResourceGroup $resourceGroup `
            -ContainerGroupName $containerName `
            -ContainerInstanceName "azurehound" `
            -ContainerRegistry $acr `
            -ContainerImage "$($acr.LoginServer)/$($imageName)" `
            -ContainerEntrypoint @("sleep", "infinity") `
            -ConfigJson $configJson `
            -ContainerUMI $containerUMI        
    }

    # Add permissions to the container UMI
    # TODO: I think the scope needs to be based on the tenant we are analyzing
    $subscriptionScope = "/subscriptions/$((Get-AzContext).Subscription.Id)"
    Add-RoleAssignment `
        -PrincipalId $containerUMI.PrincipalId `
        -RoleDefinitionName "Reader" `
        -Scope $subscriptionScope

    # Then add Graph API permission as before
    Add-GraphApiPermissionWithPropogationTest -PrincipalId $containerUMI.PrincipalId -PermissionName "Directory.Read.All"
        
    # Output the identity details needed for deployment
    Write-Host "`nSetup complete! Use these values in your deployment:" -ForegroundColor Green
    Write-Host "Identity Resource ID: $($containerUMI.Id)"
    Write-Host "Principal ID: $($containerUMI.PrincipalId)"
    
} catch {
    Write-Error "Error: $($_.Exception.Message)"
    Write-Error "Stack Trace: $($_.ScriptStackTrace)"
}

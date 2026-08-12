param(
    [Parameter(Mandatory=$true)]
    [string]$ResourceGroupName,
    [Parameter(Mandatory=$true)]
    [string]$Location,
    [Parameter(Mandatory=$true)]
    [string]$IdentityName
)

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

    try {
	    New-AzRoleAssignment -ObjectId $PrincipalId `
                        -RoleDefinitionName $RoleDefinitionName `
                        -Scope $Scope
    } catch {
        Write-Host "Failed to add role assignment skipping"
    }
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

try {
    # Create the managed identity
    if (-not ($resourceGroup = Get-AzResourceGroup -Name $ResourceGroupName -ErrorAction SilentlyContinue)) {
        Write-Host "Creating resource group: $ResourceGroupName"
        $resourceGroup = New-AzResourceGroup -Name $ResourceGroupName -Location $Location
    }

    if (-not ($containerUMI = Get-AzUserAssignedIdentity -ResourceGroupName $resourceGroup.ResourceGroupName -Name $IdentityName -ErrorAction SilentlyContinue)){
        $containerUMI = New-ManagedIdentity -ResourceGroup $resourceGroup -IdentityName $IdentityName
    }

    # Add permissions to the container UMI
    $subscriptionScope = "/subscriptions/$((Get-AzContext).Subscription.Id)"
    Add-RoleAssignment `
        -PrincipalId $containerUMI.PrincipalId `
        -RoleDefinitionName "Reader" `
        -Scope $subscriptionScope

    # Then add Graph API permission
    Add-GraphApiPermissionWithPropogationTest -PrincipalId $containerUMI.PrincipalId -PermissionName "Directory.Read.All"
        
    # Output the identity details needed for deployment
    Write-Host "`nSetup complete! Use these values in your deployment:" -ForegroundColor Green
    Write-Host "Identity Resource ID: $($containerUMI.Id)"
    Write-Host "Principal ID: $($containerUMI.PrincipalId)"
} catch {
    Write-Error "Error: $($_.Exception.Message)"
}



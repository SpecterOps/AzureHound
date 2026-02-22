param(
    [Parameter(Mandatory=$true)]
    [string]$ResourceGroupName,
    [Parameter(Mandatory=$true)]
    [string]$Location,
    [Parameter(Mandatory=$true)]
    [string]$IdentityName
)

# Ensure we're connected to Azure
if (-not (Get-AzContext)) {
    Write-Error "Not connected to Azure. Please run Connect-AzAccount first."
    exit 1
}

function New-DeploymentManagedIdentity {
    param(
        [Parameter(Mandatory=$true)]
        [string]$ResourceGroupName,
        [Parameter(Mandatory=$true)]
        [string]$Location,
        [Parameter(Mandatory=$true)]
        [string]$IdentityName
    )

    # Create Resource Group if it doesn't exist
    if (-not (Get-AzResourceGroup -Name $ResourceGroupName -ErrorAction SilentlyContinue)) {
        Write-Host "Creating resource group: $ResourceGroupName"
        New-AzResourceGroup -Name $ResourceGroupName -Location $Location
    }

    # Create the managed identity
    Write-Host "Creating managed identity: $IdentityName"
    $identity = New-AzUserAssignedIdentity -ResourceGroupName $ResourceGroupName -Name $IdentityName -Location $Location

    # Wait for identity to propagate to AAD
    Write-Host "Waiting for identity to propagate to Azure AD..."
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

function Add-GraphPermissions {
    param(
        [Parameter(Mandatory=$true)]
        [string]$PrincipalId,
        [Parameter(Mandatory=$true)]
        [hashtable]$RequiredPermissions
    )

    $token = (Get-AzAccessToken -ResourceUrl "https://graph.microsoft.com/").Token
    $graphAppId = "00000003-0000-0000-c000-000000000000"
    $graphSp = Get-AzADServicePrincipal -ApplicationId $graphAppId

    $assignedPermissions = @{}

    foreach($permission in $RequiredPermissions.Keys) {
        $assignedPermissions[$permission] = $false
        $permissionId = $RequiredPermissions[$permission]

        $headers = @{
            'Authorization' = "Bearer $token"
            'Content-Type' = 'application/json'
        }

        # Check if permission exists
        $existingAssignments = Invoke-RestMethod -Uri "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId/appRoleAssignments" -Headers $headers -Method Get
        $existingAssignment = $existingAssignments.value | Where-Object { $_.appRoleId -eq $permissionId }

        if ($existingAssignment) {
            Write-Host "$permission permission is already assigned"
            $assignedPermissions[$permission] = $true
            continue
        }

        $body = @{
            principalId = $PrincipalId
            resourceId = $graphSp.Id
            appRoleId = $permissionId
        } | ConvertTo-Json

        try {
            $apiUrl = "https://graph.microsoft.com/v1.0/servicePrincipals/$PrincipalId/appRoleAssignments"
            $result = Invoke-RestMethod -Uri $apiUrl -Headers $headers -Method Post -Body $body
            Write-Host "Successfully assigned $permission permission"
            $assignedPermissions[$permission] = $true
        }
        catch {
            Write-Error "Failed to assign $permission permission: $_"
            continue
        }
    }
    return $assignedPermissions
}

function Add-SubscriptionRoleAssignments {
    param(
        [Parameter(Mandatory=$true)]
        [string]$PrincipalId,
        [Parameter(Mandatory=$true)]
        [hashtable]$RequiredRoles,
        [Parameter(Mandatory=$false)]
        [string]$SubscriptionId = (Get-AzContext).Subscription.Id
    )

    $assignedRoles = @{}

    $subscriptionScope = "/subscriptions/$SubscriptionId"

    foreach ($role in $RequiredRoles.Keys) {
        $assignedRoles[$role] = $false
        $roleDefinitionId = $RequiredRoles[$role]

        # Check if role assignment exists
        $existingAssignment = Get-AzRoleAssignment `
            -ObjectId $PrincipalId `
            -RoleDefinitionId $roleDefinitionId `
            -Scope $subscriptionScope `
            -ErrorAction SilentlyContinue

        if ($existingAssignment) {
            Write-Host "$role is already assigned"
            $assignedRoles[$role] = $true
            continue
        }

        try {
            Write-Host "Assigning $role..."
            $assigned = $false
            $role_assignment = New-AzRoleAssignment `
                -ObjectId $PrincipalId `
                -RoleDefinitionId $roleDefinitionId `
                -Scope $subscriptionScope

            Write-Host "Did role assignment $($info)"
            if (-not $role_assignment) {
                Write-Host "Failed to assign role $role to $PrincipalId skipping this role assignment"
                continue
            }

            # Wait for role assignment to propagate
            $timeout = (Get-Date).AddMinutes(2)
            
            while ((Get-Date) -lt $timeout -and -not $assigned) {
                Write-Host "Waiting for $role assignment to propagate..."
                Start-Sleep -Seconds 10
                
                $assigned = Get-AzRoleAssignment `
                    -ObjectId $PrincipalId `
                    -RoleDefinitionId $roleDefinitionId `
                    -Scope $subscriptionScope `
                    -ErrorAction SilentlyContinue
            }

            if ($assigned) {
                Write-Host "$role assignment confirmed"
                $assignedRoles[$role] = $true
            } 
        }
        catch {
            Write-Error "Failed to assign $role. skipping this role assignment"
        }
    }
    return $assignedRoles
}

try {
    # Create the deployment managed identity
    $identity = New-DeploymentManagedIdentity -ResourceGroupName $ResourceGroupName -Location $Location -IdentityName $IdentityName

    $graphPermissionsRequired = @{
        "Application.ReadWrite.All" = "1bfefb4e-e0b5-418b-a88f-73c46d2cc8e9"
    }
    $graphPermissions = Add-GraphPermissions -PrincipalId $identity.PrincipalId -RequiredPermissions $graphPermissionsRequired

    $roleDefinitionsRequired = @{
        'Managed Identity Contributor' = 'f1a07417-d97a-45cb-824c-7a7467783830'
        'User Access Administrator' = '18d7d88d-d35e-4fb5-a5c3-7773c20a72d9'
    }
    $roles = Add-SubscriptionRoleAssignments -PrincipalId $identity.PrincipalId -RequiredRoles $roleDefinitionsRequired

    # Output the identity details needed for deployment
    Write-Host "`nDeployment identity created! Use these values in your ARM template:" -ForegroundColor Green
    Write-Host "Identity Resource ID: $($identity.Id)"
    Write-Host "Principal ID: $($identity.PrincipalId)"
    Write-Host "Client ID: $($identity.ClientId)"
    
    Write-Host "`nIMPORTANT: Ensure that all permissions have been assigned!" -ForegroundColor Yellow
    Write-Host "Please have an administrator assign the following permissions to the managed identity:"
    Write-Host "1. Azure RBAC Roles (at subscription scope):"
    foreach ($role in $roleDefinitionsRequired.Keys) {
        if (-not $roles[$role]) {
            Write-Host "Have administrator assign $role $($roleDefinitionsRequired[$role])"
        }
    }
    Write-Host "2. Graph API Permissions:"
    foreach ($graphRole in $graphPermissionsRequired.Keys) {
        if (-not $graphPermissions[$graphRole]) {
            Write-Host "Have administrator assign $graphRole $($graphPermissionsRequired[$graphRole])"
        }
    }    
} catch {
    Write-Error "Error: $($_.Exception.Message)"
    Write-Error "Stack Trace: $($_.ScriptStackTrace)"
    throw
}
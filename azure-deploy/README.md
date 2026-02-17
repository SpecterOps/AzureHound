# AzureHoundDeploy

## Overview

AzureHound now supports Managed Identity authentication.  This allows AzureHound to be run in an Azure Container Instance.  The Container Instance must be associated with a Managed Identity that has the required RBAC Roles and Graph Permissions that AzureHound requires.

This repository containes the Azure Resource Manager (ARM) Template along with supporting scripts that
allow a user to conveniently deploy and configure AzureHound to an Azure Container Instance.  Specifically this ARM template provides the following functionallity.

1) Deploy an AzureHound Instance that runs in an Azure Container Instance
2) Creates or uses an existing Container Instance
3) Creates or uses an existing Managed Identity for the Container that provides Azure permissions
4) Provides a wizard that configures the AzureHound Instance

## Process


## Prerequisites

In order for this ARM Template to create the container's User Managed Identity, the ARM Template requires an existing User Managed Identity with the following permissions:

   - Application.ReadWrite.All
   - Managed Identity Contributor
   - User Access Administrator

This repository contains a `create-deployment-umi-with-perms.ps1` script that can be used to create deployment's user managed identity.  Alternatively you can create the User Managed Identity in the Azure portal.

## AzureHound Required Permissions

AzureHound requires the following Azure permissions

 - Directory.Read.All
 - Reader 
 
The ARM Template will create the Container Instance along with a User Managed Identity that provides this permissions to the AzureHound Instance.

## Supporting Scripts

The ARM Template is designed to create the Container Instance along with a User Managed Identity that will provide AzureHound with all the permissions it needs to run.  However, it

- Create/Fix Managed Identity For The ARM Template
   `create-deployment-umi-with-perms.ps1`
- Create/Fix Managed Identity For the Container
   `create-container-umi.ps1`
- Full end to end script.
   `single-script-full-deployment.ps1`

## Notes About Approach
ManagedIdentities can be assigned permissions just like App Registration (Enterprise Applications), however you are assigning the permissions to 
the managed identity's application object id.  After creation of a Managed Identity it takes some amount of time before the application id is associated with the managed identity.  Therefore we add retry logic.

## Permissions DeploymentScript requires
The `managed-identity-permissions.sh` script will require 
the following permissions to be assigned to a managed identity.  

# Issues to document

## `single-script-full-deployment`

- it requires tenant-id, but this can be retrieved with `(Get-AzTenant).Id` after logging in with `Connect-AzAccount`
- does it require an existing container registry.
   - Maybe use `$registry = New-AzContainerRegistry -ResourceGroupName "myResourceGroup" -Name "mycontainerregistry" -EnableAdminUser -Sku Basic -Location EastUS`
   - but permissions required may be a problem.

   **maybe create a separate script**

## The docker image must be loaded into the container registry before hand.
To do this this is the following procedure:

   1) Have docker installed and running
   2) An image registry in Azure must exist or be created
   3) Have the authentication information for the existing azure image registry
   4) At the command line authenticate with the azure registry
         - 
   5) pull the latest AzureHound image
      `docker pull ghcr.io/bloodhoundad/azurehound:<tag>`  Note: Use the tag 'latest' unless you want a specific version of azurehound.
   6) tag the image to associate it with an the azure registry
      `docker tag ghcr.io/owner/repository:tag <AzureRegistryName>.azurecr.io/<AzureRepositoryName>:tag Note: Best to use latest unless you need to maintain multiple versions in Azure's registry.
   7) docker push <AzureRegistryName>.azurecr.io/<AzureRepositoryName>:tag





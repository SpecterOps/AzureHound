# Azure app registration details
$clientId = ""
$clientSecret = "" 
$tenantId = ""

# Get access token
$tokenBody = @{
    grant_type = "client_credentials"
    client_id = $clientId
    client_secret = $clientSecret
    scope = "https://graph.microsoft.com/.default"
}

Write-Host "Getting access token..." -ForegroundColor Yellow
$tokenResponse = Invoke-RestMethod -Uri "https://login.microsoftonline.com/$tenantId/oauth2/v2.0/token" -Method POST -Body $tokenBody
$headers = @{ Authorization = "Bearer $($tokenResponse.access_token)" }

Write-Host $tokenResponse.access_token -ForegroundColor Yellow
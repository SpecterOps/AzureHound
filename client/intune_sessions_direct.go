// client/intune_sessions_direct.go - Correct implementation using AzureHound patterns
package client

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// Helper function to set default parameters (copied from intune_devices.go)
func setDefaultSessionParams(params *query.GraphParams) {
	if params.Top == 0 {
		params.Top = 100 // Smaller default for sign-in logs
	}
}

// CollectSessionDataDirectly collects session data from Microsoft Graph Sign-In Logs
func (s *azureClient) CollectSessionDataDirectly(ctx context.Context) <-chan AzureResult[azure.DeviceSessionData] {
	out := make(chan AzureResult[azure.DeviceSessionData])

	go func() {
		defer close(out)

		fmt.Printf("🔍 Starting session data collection from Graph API...\n")

		// Get sign-in logs using the correct AzureHound pattern
		params := query.GraphParams{
			Filter: fmt.Sprintf("createdDateTime ge %s", time.Now().AddDate(0, 0, -7).Format(time.RFC3339)),
			Select: []string{
				"id", "createdDateTime", "userDisplayName", "userPrincipalName",
				"userId", "appDisplayName", "clientAppUsed", "ipAddress",
				"isInteractive", "status", "deviceDetail", "riskState", "riskLevelAggregated",
			},
		}
		setDefaultSessionParams(&params)

		signInLogsChan := s.listSignInLogs(ctx, params)

		var signInLogs []SignInEvent
		errorCount := 0

		// Collect all sign-in logs
		for result := range signInLogsChan {
			if result.Error != nil {
				fmt.Printf("❌ Error collecting sign-in log: %v\n", result.Error)
				errorCount++
				if errorCount > 5 { // Stop after too many errors
					out <- AzureResult[azure.DeviceSessionData]{
						Error: fmt.Errorf("too many errors collecting sign-in logs: %w", result.Error),
					}
					return
				}
				continue
			}
			signInLogs = append(signInLogs, result.Ok)
		}

		fmt.Printf("📊 Retrieved %d sign-in events (%d errors)\n", len(signInLogs), errorCount)

		if len(signInLogs) == 0 {
			fmt.Printf("⚠️  No sign-in logs found. This could mean:\n")
			fmt.Printf("   • No users signed in recently (last 7 days)\n")
			fmt.Printf("   • Missing AuditLog.Read.All permission\n")
			fmt.Printf("   • Azure AD Premium license required for audit logs\n")
			fmt.Printf("   • Sign-in logs not available in this tenant\n")

			out <- AzureResult[azure.DeviceSessionData]{
				Error: fmt.Errorf("no sign-in logs found - check permissions and recent user activity"),
			}
			return
		}

		// Process the logs into device sessions
		deviceSessions := s.processSignInLogsSimple(signInLogs)
		fmt.Printf("🔄 Created %d device session records\n", len(deviceSessions))

		// Send results
		for _, sessionData := range deviceSessions {
			out <- AzureResult[azure.DeviceSessionData]{Ok: sessionData}
		}
	}()

	return out
}

// GetUserSignInActivity retrieves sign-in activity for a specific user
func (s *azureClient) GetUserSignInActivity(ctx context.Context, userPrincipalName string, days int) ([]SignInEvent, error) {
	fmt.Printf("🔍 Getting sign-in activity for user: %s\n", userPrincipalName)

	params := query.GraphParams{
		Filter: fmt.Sprintf("userPrincipalName eq '%s' and createdDateTime ge %s",
			userPrincipalName, time.Now().AddDate(0, 0, -days).Format(time.RFC3339)),
		Top: 50,
	}

	signInLogsChan := s.listSignInLogs(ctx, params)

	var signInLogs []SignInEvent
	for result := range signInLogsChan {
		if result.Error != nil {
			return nil, result.Error
		}
		signInLogs = append(signInLogs, result.Ok)
	}

	return signInLogs, nil
}

// GetDeviceSignInActivity retrieves sign-in activity for a specific device
func (s *azureClient) GetDeviceSignInActivity(ctx context.Context, deviceId string, days int) ([]SignInEvent, error) {
	fmt.Printf("🔍 Getting sign-in activity for device: %s\n", deviceId)

	params := query.GraphParams{
		Filter: fmt.Sprintf("deviceDetail/deviceId eq '%s' and createdDateTime ge %s",
			deviceId, time.Now().AddDate(0, 0, -days).Format(time.RFC3339)),
		Top: 50,
	}

	signInLogsChan := s.listSignInLogs(ctx, params)

	var signInLogs []SignInEvent
	for result := range signInLogsChan {
		if result.Error != nil {
			return nil, result.Error
		}
		signInLogs = append(signInLogs, result.Ok)
	}

	return signInLogs, nil
}

// listSignInLogs follows the exact AzureHound pattern from intune_devices.go
func (s *azureClient) listSignInLogs(ctx context.Context, params query.GraphParams) <-chan AzureResult[SignInEvent] {
	var (
		out  = make(chan AzureResult[SignInEvent])
		path = fmt.Sprintf("/%s/auditLogs/signIns", constants.GraphApiVersion)
	)

	setDefaultSessionParams(&params)

	// Use the exact same pattern as AzureHound - call getAzureObjectList
	go getAzureObjectList[SignInEvent](s.msgraph, ctx, path, params, out)
	return out
}

// processSignInLogsSimple converts sign-in logs to device session data
func (s *azureClient) processSignInLogsSimple(signInLogs []SignInEvent) []azure.DeviceSessionData {
	fmt.Printf("🔄 Processing %d sign-in logs into device sessions\n", len(signInLogs))

	// Group sign-ins by device
	deviceGroups := make(map[string][]SignInEvent)

	for _, signIn := range signInLogs {
		deviceKey := signIn.DeviceDetail.DeviceId
		if deviceKey == "" {
			deviceKey = signIn.DeviceDetail.DisplayName
		}
		if deviceKey == "" {
			deviceKey = fmt.Sprintf("Unknown_%s", signIn.IPAddress)
		}

		deviceGroups[deviceKey] = append(deviceGroups[deviceKey], signIn)
	}

	fmt.Printf("📊 Grouped sign-ins into %d devices\n", len(deviceGroups))

	var results []azure.DeviceSessionData

	for deviceKey, sessions := range deviceGroups {
		fmt.Printf("🔄 Processing device: %s (%d sessions)\n", deviceKey, len(sessions))
		// FIXED: Call the correct function name
		sessionData := s.createDeviceSessionData(deviceKey, sessions)
		results = append(results, sessionData)
	}

	return results
}

// createDeviceSessionData creates session data for a device (main function)
func (s *azureClient) createDeviceSessionData(deviceKey string, signIns []SignInEvent) azure.DeviceSessionData {
	now := time.Now()

	// Create basic device info
	var deviceInfo azure.IntuneDevice
	if len(signIns) > 0 {
		first := signIns[0]
		deviceInfo = azure.IntuneDevice{
			ID:                first.DeviceDetail.DeviceId,
			DeviceName:        first.DeviceDetail.DisplayName,
			OperatingSystem:   first.DeviceDetail.OperatingSystem,
			UserPrincipalName: first.UserPrincipalName,
			UserDisplayName:   first.UserDisplayName,
			LastSyncDateTime:  first.CreatedDateTime,
			ComplianceState:   getComplianceString(first.DeviceDetail.IsCompliant),
			AzureADDeviceID:   first.DeviceDetail.DeviceId,
		}

		if deviceInfo.DeviceName == "" {
			deviceInfo.DeviceName = deviceKey
		}
	}

	// Process sessions
	var activeSessions []azure.ActiveSession
	var loggedOnUsers []azure.LoggedOnUser
	userMap := make(map[string]bool)

	adminCount := 0
	serviceCount := 0
	suspiciousActivities := []azure.SuspiciousActivity{}

	for i, signIn := range signIns {
		// Only process successful sign-ins
		if signIn.Status.ErrorCode == 0 {
			// FIXED: Use the correct function names with fallback approach
			isAdmin := s.isAdminUserEnhanced(context.Background(), signIn.UserPrincipalName)
			isService := s.isServiceUserEnhanced(context.Background(), signIn.UserPrincipalName)

			if isAdmin {
				adminCount++
			}

			session := azure.ActiveSession{
				SessionID:     i + 1,
				UserName:      getUsernameFromUPN(signIn.UserPrincipalName),
				DomainName:    getDomainFromUPN(signIn.UserPrincipalName),
				SessionType:   getSessionType(signIn.ClientAppUsed),
				SessionState:  "Active",
				LogonTime:     signIn.CreatedDateTime,
				IdleTime:      getIdleTime(signIn.CreatedDateTime),
				ClientName:    signIn.DeviceDetail.DisplayName,
				ClientAddress: signIn.IPAddress,
				IsElevated:    isAdmin,
			}
			activeSessions = append(activeSessions, session)

			// Add unique users
			if !userMap[signIn.UserPrincipalName] {
				userMap[signIn.UserPrincipalName] = true

				if isService {
					serviceCount++
				}

				user := azure.LoggedOnUser{
					UserName:         getUsernameFromUPN(signIn.UserPrincipalName),
					DomainName:       getDomainFromUPN(signIn.UserPrincipalName),
					SID:              signIn.UserId,
					LogonType:        "Interactive",
					AuthPackage:      "AzureAD",
					LogonTime:        signIn.CreatedDateTime,
					LogonServer:      "login.microsoftonline.com",
					HasCachedCreds:   true,
					IsServiceAccount: isService,
					TokenPrivileges:  getTokenPrivileges(isAdmin),
				}
				loggedOnUsers = append(loggedOnUsers, user)
			}
		}

		// Check for suspicious activities
		if signIn.RiskState == "atRisk" || signIn.RiskLevelAggregated == "high" {
			activity := azure.SuspiciousActivity{
				ActivityType: "High_Risk_Sign_In",
				Description:  fmt.Sprintf("High risk sign-in for %s from %s", signIn.UserDisplayName, signIn.IPAddress),
				RiskLevel:    "High",
				Evidence:     []string{fmt.Sprintf("Risk: %s", signIn.RiskState)},
				DetectedAt:   signIn.CreatedDateTime,
				UserName:     signIn.UserPrincipalName,
				SessionID:    0,
			}
			suspiciousActivities = append(suspiciousActivities, activity)
		}
	}

	// Create session data
	sessionData := azure.SessionData{
		DeviceInfo: azure.DeviceInfo{
			ComputerName:  deviceInfo.DeviceName,
			Domain:        "AZUREAD",
			User:          "SYSTEM",
			Timestamp:     now.Format(time.RFC3339),
			ScriptVersion: "azurehound-graph-1.0",
		},
		ActiveSessions: activeSessions,
		LoggedOnUsers:  loggedOnUsers,
		SecurityIndicators: azure.SessionSecurityInfo{
			AdminSessionsActive:     adminCount > 0,
			RemoteSessionsActive:    false,
			ServiceAccountSessions:  serviceCount > 0,
			CredentialTheftRisk:     getRiskLevel(adminCount),
			PrivilegeEscalationRisk: getRiskLevel(adminCount),
			SuspiciousActivities:    suspiciousActivities,
		},
		Summary: azure.SessionDataSummary{
			TotalActiveSessions: len(activeSessions),
			UniqueUsers:         len(loggedOnUsers),
			AdminSessions:       adminCount,
			RemoteSessions:      0,
			ServiceSessions:     serviceCount,
			CredentialExposure:  len(loggedOnUsers),
		},
	}

	fmt.Printf("✅ Created session data for %s: %d sessions, %d users, %d admin sessions, %d service accounts\n",
		deviceInfo.DeviceName, len(activeSessions), len(loggedOnUsers), adminCount, serviceCount)

	return azure.DeviceSessionData{
		Device:      deviceInfo,
		SessionData: sessionData,
		CollectedAt: now,
	}
}

// Utility helper functions
func getComplianceString(isCompliant bool) string {
	if isCompliant {
		return "compliant"
	}
	return "noncompliant"
}

func getUsernameFromUPN(upn string) string {
	if upn == "" {
		return "Unknown"
	}
	parts := strings.Split(upn, "@")
	return parts[0]
}

func getDomainFromUPN(upn string) string {
	if upn == "" {
		return "AZUREAD"
	}
	parts := strings.Split(upn, "@")
	if len(parts) > 1 {
		return strings.ToUpper(parts[1])
	}
	return "AZUREAD"
}

func getSessionType(clientApp string) string {
	lower := strings.ToLower(clientApp)
	if strings.Contains(lower, "mobile") {
		return "Mobile"
	}
	if strings.Contains(lower, "browser") {
		return "Browser"
	}
	if strings.Contains(lower, "desktop") {
		return "Desktop"
	}
	return "Interactive"
}

func getIdleTime(logonTime time.Time) string {
	duration := time.Since(logonTime)
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

// Helper function for GUID detection
func isGUIDLike(s string) bool {
	guidPattern := `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`
	matched, _ := regexp.MatchString(guidPattern, s)
	return matched
}

func getTokenPrivileges(isAdmin bool) []string {
	if isAdmin {
		return []string{"SeDebugPrivilege", "SeImpersonatePrivilege"}
	}
	return []string{}
}

func getRiskLevel(adminCount int) string {
	if adminCount > 2 {
		return "High"
	}
	if adminCount > 0 {
		return "Medium"
	}
	return "Low"
}
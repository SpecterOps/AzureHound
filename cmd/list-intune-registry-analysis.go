// cmd/list-intune-registry-analysis.go
package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listIntuneRegistryAnalysisCmd)
}

var listIntuneRegistryAnalysisCmd = &cobra.Command{
	Use:          "intune-registry-analysis",
	Long:         "Performs security analysis on collected registry data and formats for BloodHound",
	Run:          listIntuneRegistryAnalysisCmdImpl,
	SilenceUsage: true,
}

func listIntuneRegistryAnalysisCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := context.WithCancel(cmd.Context())
	defer stop()

	azClient := connectAndCreateClient()

	// Add flag to choose between real and simulated analysis
	useRealAnalysis, _ := cmd.Flags().GetBool("real-analysis")

	var analysisResults []azure.DeviceSecurityAnalysis
	var err error

	if useRealAnalysis {
		fmt.Printf("🔍 Performing real registry security analysis...\n")
		// Validate script deployment first
		if err := azClient.ValidateScriptDeployment(ctx); err != nil {
			fmt.Printf("⚠️ Script validation failed: %v\n", err)
			fmt.Printf("ℹ️ Falling back to device compliance analysis\n")
			analysisResults, err = performDeviceAnalysisWithoutScripts(ctx, azClient)
		} else {
			analysisResults, err = performRealRegistrySecurityAnalysis(ctx, azClient)
		}
	} else {
		fmt.Printf("ℹ️ Performing device compliance analysis (no registry scripts)\n")
		analysisResults, err = performDeviceAnalysisWithoutScripts(ctx, azClient)
	}

	if err != nil {
		exit(err)
	}

	displayAnalysisResults(analysisResults)
}

// Add flag to command initialization
func init() {
	listRootCmd.AddCommand(listIntuneRegistryAnalysisCmd)
	listIntuneRegistryAnalysisCmd.Flags().Bool("real-analysis", false, "Perform real registry analysis using deployed scripts")
}

// cmd/list-intune-registry-analysis.go - Add this function

func displayAnalysisResults(results []azure.DeviceSecurityAnalysis) {
	fmt.Printf("\n=== INTUNE DEVICE SECURITY ANALYSIS RESULTS ===\n\n")

	if len(results) == 0 {
		fmt.Printf("❌ No devices were analyzed\n")
		return
	}

	// Calculate summary statistics
	summary := calculateSummaryStats(results)
	displaySummary(summary, len(results))

	// Display detailed results for each device
	fmt.Printf("📱 DEVICE DETAILS:\n")
	fmt.Printf("═══════════════════════════════════════════════════════════════\n\n")

	for i, result := range results {
		displayDeviceResult(i+1, result)
	}

	// Display recommendations
	displayRecommendations(results)
}

func displaySummary(summary map[string]interface{}, totalDevices int) {
	fmt.Printf("📊 ANALYSIS SUMMARY\n")
	fmt.Printf("─────────────────────────────────────────────────────────────\n")

	// Compliance summary
	if complianceSummary, ok := summary["compliance_summary"].(map[string]interface{}); ok {
		fmt.Printf("🎯 Compliance Overview:\n")
		fmt.Printf("   • Total Devices: %d\n", totalDevices)
		fmt.Printf("   • Compliant: %v\n", complianceSummary["compliant"])
		fmt.Printf("   • Partially Compliant: %v\n", complianceSummary["partially_compliant"])
		fmt.Printf("   • Non-Compliant: %v\n", complianceSummary["non_compliant"])
		fmt.Printf("   • Compliance Rate: %v\n", complianceSummary["compliance_rate"])
		fmt.Printf("\n")
	}

	// Risk summary
	if riskSummary, ok := summary["risk_summary"].(map[string]interface{}); ok {
		fmt.Printf("⚠️  Risk Assessment:\n")
		fmt.Printf("   • Average Risk Score: %v/100\n", riskSummary["average_risk_score"])
		fmt.Printf("   • Total Security Findings: %v\n", riskSummary["total_findings"])

		if findingsBySeverity, ok := riskSummary["findings_by_severity"].(map[string]int); ok {
			fmt.Printf("   • Critical: %d | High: %d | Medium: %d | Low: %d\n",
				findingsBySeverity["CRITICAL"],
				findingsBySeverity["HIGH"],
				findingsBySeverity["MEDIUM"],
				findingsBySeverity["LOW"])
		}
		fmt.Printf("\n")
	}

	// Device breakdown
	if deviceBreakdown, ok := summary["device_breakdown"].(map[string]interface{}); ok {
		fmt.Printf("🔍 Risk Distribution:\n")
		fmt.Printf("   • High Risk (70-100): %v devices\n", deviceBreakdown["high_risk_devices"])
		fmt.Printf("   • Medium Risk (30-69): %v devices\n", deviceBreakdown["medium_risk_devices"])
		fmt.Printf("   • Low Risk (0-29): %v devices\n", deviceBreakdown["low_risk_devices"])
		fmt.Printf("\n")
	}
}

func displayDeviceResult(index int, result azure.DeviceSecurityAnalysis) {
	// Device header with risk level emoji
	riskEmoji := getRiskEmoji(result.RiskScore)
	statusEmoji := getComplianceEmoji(result.ComplianceStatus)

	fmt.Printf("%s %s Device #%d: %s\n",
		riskEmoji, statusEmoji, index, result.Device.DeviceName)

	// Basic device info
	fmt.Printf("   🆔 Device ID: %s\n", result.Device.ID)
	fmt.Printf("   💻 OS: %s %s\n", result.Device.OperatingSystem, result.Device.OSVersion)
	fmt.Printf("   👤 User: %s\n", getDisplayValue(result.Device.UserPrincipalName))
	fmt.Printf("   📊 Risk Score: %d/100 (%s)\n", result.RiskScore, getRiskLevel(result.RiskScore))
	fmt.Printf("   ✅ Compliance: %s\n", result.ComplianceStatus)
	fmt.Printf("   🕒 Last Analysis: %s\n", result.AnalysisTimestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("   🔄 Last Sync: %s\n", result.Device.LastSyncDateTime.Format("2006-01-02 15:04:05"))

	// Security findings
	if len(result.SecurityFindings) > 0 {
		fmt.Printf("   🚨 Security Findings (%d):\n", len(result.SecurityFindings))

		// Group findings by severity
		findingsBySeverity := groupFindingsBySeverity(result.SecurityFindings)

		for _, severity := range []string{"CRITICAL", "HIGH", "MEDIUM", "LOW", "INFO"} {
			if findings, exists := findingsBySeverity[severity]; exists && len(findings) > 0 {
				for _, finding := range findings {
					emoji := getSecurityEmoji(finding.Severity)
					fmt.Printf("      %s %s\n", emoji, finding.Title)
					fmt.Printf("         📝 %s\n", finding.Description)

					// Show evidence for high/critical findings
					if finding.Severity == "HIGH" || finding.Severity == "CRITICAL" {
						if len(finding.Evidence) > 0 {
							fmt.Printf("         🔍 Evidence: %s\n", finding.Evidence[0])
						}
						if len(finding.Recommendations) > 0 {
							fmt.Printf("         💡 Recommendation: %s\n", finding.Recommendations[0])
						}
					}
				}
			}
		}
	} else {
		fmt.Printf("   ✅ No security findings detected\n")
	}

	// Escalation vectors
	if len(result.EscalationVectors) > 0 {
		fmt.Printf("   ⚡ Privilege Escalation Vectors (%d):\n", len(result.EscalationVectors))
		for _, vector := range result.EscalationVectors {
			fmt.Printf("      🎯 %s: %s → %s\n", vector.Type, vector.Source, vector.Target)
			fmt.Printf("         Method: %s (Complexity: %s)\n", vector.Method, vector.Complexity)
		}
	}

	fmt.Printf("\n")
}

func displayRecommendations(results []azure.DeviceSecurityAnalysis) {
	criticalCount := 0
	highCount := 0
	nonCompliantCount := 0

	for _, result := range results {
		if result.ComplianceStatus == "NON_COMPLIANT" {
			nonCompliantCount++
		}

		for _, finding := range result.SecurityFindings {
			switch finding.Severity {
			case "CRITICAL":
				criticalCount++
			case "HIGH":
				highCount++
			}
		}
	}

	if criticalCount > 0 || highCount > 0 || nonCompliantCount > 0 {
		fmt.Printf("🎯 IMMEDIATE ACTIONS REQUIRED\n")
		fmt.Printf("─────────────────────────────────────────────────────────────\n")

		if criticalCount > 0 {
			fmt.Printf("🔥 CRITICAL: %d critical security issues need immediate attention\n", criticalCount)
		}

		if highCount > 0 {
			fmt.Printf("🚨 HIGH: %d high-severity issues should be addressed soon\n", highCount)
		}

		if nonCompliantCount > 0 {
			fmt.Printf("📋 COMPLIANCE: %d devices are non-compliant with policies\n", nonCompliantCount)
		}

		fmt.Printf("\n💡 Recommended Actions:\n")
		fmt.Printf("   1. Address all CRITICAL and HIGH severity findings immediately\n")
		fmt.Printf("   2. Review and remediate non-compliant devices\n")
		fmt.Printf("   3. Update device compliance policies if needed\n")
		fmt.Printf("   4. Schedule regular security assessments\n")
		fmt.Printf("   5. Consider additional endpoint protection measures\n\n")
	} else {
		fmt.Printf("✅ GOOD NEWS!\n")
		fmt.Printf("─────────────────────────────────────────────────────────────\n")
		fmt.Printf("No critical security issues were found in the analyzed devices.\n")
		fmt.Printf("Continue regular monitoring to maintain security posture.\n\n")
	}
}

// Helper functions for display formatting

func getRiskEmoji(riskScore int) string {
	switch {
	case riskScore >= 70:
		return "🔴" // High risk
	case riskScore >= 30:
		return "🟡" // Medium risk
	default:
		return "🟢" // Low risk
	}
}

func getComplianceEmoji(status string) string {
	switch status {
	case "COMPLIANT":
		return "✅"
	case "PARTIALLY_COMPLIANT":
		return "⚠️"
	case "NON_COMPLIANT":
		return "❌"
	default:
		return "❓"
	}
}

func getSecurityEmoji(severity string) string {
	switch severity {
	case "CRITICAL":
		return "🔥"
	case "HIGH":
		return "🚨"
	case "MEDIUM":
		return "⚠️"
	case "LOW":
		return "ℹ️"
	case "INFO":
		return "📋"
	default:
		return "❓"
	}
}

func getRiskLevel(riskScore int) string {
	switch {
	case riskScore >= 70:
		return "HIGH RISK"
	case riskScore >= 30:
		return "MEDIUM RISK"
	default:
		return "LOW RISK"
	}
}

func getDisplayValue(value string) string {
	if value == "" {
		return "Not specified"
	}
	return value
}

func groupFindingsBySeverity(findings []azure.SecurityFinding) map[string][]azure.SecurityFinding {
	grouped := make(map[string][]azure.SecurityFinding)

	for _, finding := range findings {
		grouped[finding.Severity] = append(grouped[finding.Severity], finding)
	}

	return grouped
}

// calculateSummaryStats function (referenced in the display)
func calculateSummaryStats(results []azure.DeviceSecurityAnalysis) map[string]interface{} {
	if len(results) == 0 {
		return map[string]interface{}{}
	}

	compliantCount := 0
	partiallyCompliantCount := 0
	nonCompliantCount := 0
	totalRiskScore := 0
	totalFindings := 0
	severityCounts := map[string]int{
		"CRITICAL": 0,
		"HIGH":     0,
		"MEDIUM":   0,
		"LOW":      0,
		"INFO":     0,
	}

	for _, result := range results {
		switch result.ComplianceStatus {
		case "COMPLIANT":
			compliantCount++
		case "PARTIALLY_COMPLIANT":
			partiallyCompliantCount++
		case "NON_COMPLIANT":
			nonCompliantCount++
		}

		totalRiskScore += result.RiskScore
		totalFindings += len(result.SecurityFindings)

		for _, finding := range result.SecurityFindings {
			severityCounts[finding.Severity]++
		}
	}

	avgRiskScore := float64(totalRiskScore) / float64(len(results))
	complianceRate := float64(compliantCount) / float64(len(results)) * 100

	return map[string]interface{}{
		"compliance_summary": map[string]interface{}{
			"compliant":           compliantCount,
			"partially_compliant": partiallyCompliantCount,
			"non_compliant":       nonCompliantCount,
			"compliance_rate":     fmt.Sprintf("%.1f%%", complianceRate),
		},
		"risk_summary": map[string]interface{}{
			"average_risk_score":   fmt.Sprintf("%.1f", avgRiskScore),
			"total_findings":       totalFindings,
			"findings_by_severity": severityCounts,
		},
		"device_breakdown": map[string]interface{}{
			"high_risk_devices":   countDevicesByRiskLevel(results, 70, 100),
			"medium_risk_devices": countDevicesByRiskLevel(results, 30, 69),
			"low_risk_devices":    countDevicesByRiskLevel(results, 0, 29),
		},
	}
}

func countDevicesByRiskLevel(results []azure.DeviceSecurityAnalysis, minRisk, maxRisk int) int {
	count := 0
	for _, result := range results {
		if result.RiskScore >= minRisk && result.RiskScore <= maxRisk {
			count++
		}
	}
	return count
}

func performDeviceAnalysisWithoutScripts(ctx context.Context, azClient client.AzureClient) ([]azure.DeviceSecurityAnalysis, error) {
	fmt.Printf("Starting device analysis without script execution...")

	var results []azure.DeviceSecurityAnalysis

	// Just analyze devices based on Intune compliance data
	devices := azClient.ListIntuneDevices(ctx, query.GraphParams{})

	for deviceResult := range devices {
		if deviceResult.Error != nil {
			fmt.Printf("Error getting device: %v", deviceResult.Error)
			continue
		}

		device := deviceResult.Ok

		// Skip non-Windows devices
		if !strings.Contains(strings.ToLower(device.OperatingSystem), "windows") {
			continue
		}

		// Create analysis based on device compliance state
		analysis := analyzeDeviceComplianceOnly(device)
		results = append(results, analysis)
	}

	fmt.Printf("Analyzed %d devices based on compliance data", len(results))
	return results, nil
}

func analyzeDeviceComplianceOnly(device azure.IntuneDevice) azure.DeviceSecurityAnalysis {
	analysis := azure.DeviceSecurityAnalysis{
		Device:            device,
		AnalysisTimestamp: time.Now(),
		SecurityFindings:  []azure.SecurityFinding{},
		EscalationVectors: []azure.EscalationVector{},
		RiskScore:         0,
		ComplianceStatus:  "COMPLIANT",
	}

	// Analyze based on device properties
	if device.ComplianceState != "compliant" {
		finding := azure.SecurityFinding{
			ID:              "DEVICE_NON_COMPLIANT",
			Title:           "Device Non-Compliant",
			Severity:        "MEDIUM",
			Category:        "Compliance",
			Description:     "Device does not meet compliance requirements",
			Evidence:        []string{fmt.Sprintf("State: %s", device.ComplianceState)},
			Recommendations: []string{"Review compliance policies"},
			MITREAttack:     []string{"T1562"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore = 30
		analysis.ComplianceStatus = "NON_COMPLIANT"
	}

	// Check for old sync dates
	if time.Since(device.LastSyncDateTime) > 7*24*time.Hour {
		finding := azure.SecurityFinding{
			ID:              "DEVICE_STALE_SYNC",
			Title:           "Device Not Recently Synced",
			Severity:        "LOW",
			Category:        "Management",
			Description:     "Device hasn't synced with Intune recently",
			Evidence:        []string{fmt.Sprintf("Last sync: %s", device.LastSyncDateTime.Format("2006-01-02"))},
			Recommendations: []string{"Check device connectivity", "Force sync"},
			MITREAttack:     []string{},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 10
	}

	return analysis
}

// performRealRegistrySecurityAnalysis performs actual registry data collection and analysis
func performRealRegistrySecurityAnalysis(ctx context.Context, azClient client.AzureClient) ([]azure.DeviceSecurityAnalysis, error) {
	var (
		out          = make([]azure.DeviceSecurityAnalysis, 0)
		successCount = 0
		errorCount   = 0
	)

	fmt.Printf("Starting real registry security analysis...")

	// Use the real registry collection function from your client
	deviceRegistryData := azClient.CollectRegistryDataFromAllDevices(ctx)

	for registryResult := range deviceRegistryData {
		if registryResult.Error != nil {
			fmt.Printf("Error collecting registry data: %v", registryResult.Error)
			errorCount++
			continue
		}

		// Perform real security analysis on the collected registry data
		analysis := performBasicDeviceSecurityAnalysis(registryResult.Ok)

		// Enhance the analysis with additional checks
		enhanceSecurityAnalysis(&analysis, registryResult.Ok)

		out = append(out, analysis)
		successCount++

		fmt.Printf("Analyzed device %s: %d findings, risk score %d",
			analysis.Device.DeviceName,
			len(analysis.SecurityFindings),
			analysis.RiskScore)
	}

	fmt.Printf("Registry analysis completed: %d successful, %d errors", successCount, errorCount)

	if successCount == 0 && errorCount > 0 {
		return nil, fmt.Errorf("failed to analyze any devices successfully (%d errors)", errorCount)
	}

	return out, nil
}

// performBasicDeviceSecurityAnalysis - your existing real analysis function
func performBasicDeviceSecurityAnalysis(deviceData azure.DeviceRegistryData) azure.DeviceSecurityAnalysis {
	analysis := azure.DeviceSecurityAnalysis{
		Device:            deviceData.Device,
		AnalysisTimestamp: deviceData.CollectedAt,
		SecurityFindings:  []azure.SecurityFinding{},
		EscalationVectors: []azure.EscalationVector{},
		RiskScore:         0,
		ComplianceStatus:  "COMPLIANT",
	}

	// UAC Disabled Check
	if deviceData.RegistryData.SecurityIndicators.UACDisabled {
		finding := azure.SecurityFinding{
			ID:              "UAC_DISABLED",
			Title:           "User Account Control Disabled",
			Severity:        "HIGH",
			Category:        "Privilege Escalation",
			Description:     "UAC is disabled, allowing privilege escalation attacks",
			Evidence:        []string{"UAC is disabled in registry"},
			Recommendations: []string{"Enable UAC through Group Policy or registry"},
			MITREAttack:     []string{"T1548.002"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 25

		// Add escalation vector for UAC bypass
		vector := azure.EscalationVector{
			VectorID:      "UAC_BYPASS_001",
			Type:          "Privilege Escalation",
			Source:        "Standard User",
			Target:        "Administrator",
			Method:        "UAC Disabled",
			RequiredPrivs: []string{"User"},
			Complexity:    "Low",
			Impact:        "High",
			Conditions:    []string{"UAC disabled"},
		}
		analysis.EscalationVectors = append(analysis.EscalationVectors, vector)
	}

	// Auto Admin Logon Check
	if deviceData.RegistryData.SecurityIndicators.AutoAdminLogon {
		finding := azure.SecurityFinding{
			ID:              "AUTO_ADMIN_LOGON",
			Title:           "Automatic Administrator Logon Enabled",
			Severity:        "CRITICAL",
			Category:        "Credential Exposure",
			Description:     "Automatic administrator logon exposes admin credentials",
			Evidence:        []string{"AutoAdminLogon is enabled in registry"},
			Recommendations: []string{"Disable automatic administrator logon", "Use secure credential storage"},
			MITREAttack:     []string{"T1552.002"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 40

		// Add escalation vector for credential access
		vector := azure.EscalationVector{
			VectorID:      "CRED_ACCESS_001",
			Type:          "Credential Access",
			Source:        "Local Access",
			Target:        "Administrator Credentials",
			Method:        "Registry Credential Storage",
			RequiredPrivs: []string{"Local Access"},
			Complexity:    "Low",
			Impact:        "Critical",
			Conditions:    []string{"AutoAdminLogon enabled"},
		}
		analysis.EscalationVectors = append(analysis.EscalationVectors, vector)
	}

	// Weak Service Permissions Check
	if deviceData.RegistryData.SecurityIndicators.WeakServicePermissions {
		finding := azure.SecurityFinding{
			ID:              "WEAK_SERVICE_PERMS",
			Title:           "Weak Service Permissions Detected",
			Severity:        "MEDIUM",
			Category:        "Privilege Escalation",
			Description:     "Services with weak permissions can be exploited for privilege escalation",
			Evidence:        []string{"Weak service permissions found in registry"},
			Recommendations: []string{"Review and restrict service permissions", "Apply principle of least privilege"},
			MITREAttack:     []string{"T1543.003"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 15
	}

	// Suspicious Startup Items Check
	if len(deviceData.RegistryData.SecurityIndicators.SuspiciousStartupItems) > 0 {
		evidence := make([]string, 0, len(deviceData.RegistryData.SecurityIndicators.SuspiciousStartupItems))
		for _, item := range deviceData.RegistryData.SecurityIndicators.SuspiciousStartupItems {
			evidence = append(evidence, fmt.Sprintf("%s: %s", item.Name, item.Value))
		}

		finding := azure.SecurityFinding{
			ID:              "SUSPICIOUS_STARTUP",
			Title:           "Suspicious Startup Items Detected",
			Severity:        "MEDIUM",
			Category:        "Persistence",
			Description:     fmt.Sprintf("Found %d suspicious startup items", len(deviceData.RegistryData.SecurityIndicators.SuspiciousStartupItems)),
			Evidence:        evidence,
			Recommendations: []string{"Review startup items", "Remove unauthorized persistence mechanisms"},
			MITREAttack:     []string{"T1547.001"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 10
	}

	// Set compliance status based on risk score
	if analysis.RiskScore >= 50 {
		analysis.ComplianceStatus = "NON_COMPLIANT"
	} else if analysis.RiskScore >= 25 {
		analysis.ComplianceStatus = "PARTIALLY_COMPLIANT"
	}

	return analysis
}

// enhanceSecurityAnalysis adds additional security checks and analysis
func enhanceSecurityAnalysis(analysis *azure.DeviceSecurityAnalysis, deviceData azure.DeviceRegistryData) {
	// Check device compliance state from Intune
	if deviceData.Device.ComplianceState != "compliant" {
		finding := azure.SecurityFinding{
			ID:              "DEVICE_NON_COMPLIANT",
			Title:           "Device Non-Compliant with Intune Policies",
			Severity:        "MEDIUM",
			Category:        "Compliance",
			Description:     "Device does not meet Intune compliance requirements",
			Evidence:        []string{fmt.Sprintf("Compliance state: %s", deviceData.Device.ComplianceState)},
			Recommendations: []string{"Review device compliance policies", "Update device configuration"},
			MITREAttack:     []string{"T1562"},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
		analysis.RiskScore += 15
	}

	// Check for old OS versions (basic heuristic)
	if deviceData.Device.OSVersion != "" && len(deviceData.Device.OSVersion) > 0 {
		// Add informational finding about OS version
		finding := azure.SecurityFinding{
			ID:              "OS_VERSION_INFO",
			Title:           "Operating System Information",
			Severity:        "INFO",
			Category:        "Information",
			Description:     "Device OS version recorded for security posture assessment",
			Evidence:        []string{fmt.Sprintf("OS: %s, Version: %s", deviceData.Device.OperatingSystem, deviceData.Device.OSVersion)},
			Recommendations: []string{"Ensure OS is up to date with latest security patches"},
			MITREAttack:     []string{},
		}
		analysis.SecurityFindings = append(analysis.SecurityFindings, finding)
	}

	// Update compliance status if it was degraded
	if analysis.RiskScore >= 50 {
		analysis.ComplianceStatus = "NON_COMPLIANT"
	} else if analysis.RiskScore >= 25 {
		analysis.ComplianceStatus = "PARTIALLY_COMPLIANT"
	}
}

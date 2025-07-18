// cmd/list-group-membership.go
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listGroupMembershipCmd)
}

var listGroupMembershipCmd = &cobra.Command{
	Use:          "group-membership",
	Long:         "Collects Azure AD group membership and user role assignment data (focused on BloodHound essentials)",
	Run:          listGroupMembershipCmdImpl,
	SilenceUsage: true,
}

func listGroupMembershipCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	azClient := connectAndCreateClient()

	fmt.Printf("🔄 Collecting focused Azure AD data for BloodHound...\n\n")
	startTime := time.Now()

	// Collect only essential data
	result, err := collectAllGraphData(ctx, azClient)
	if err != nil {
		exit(err)
	}

	duration := time.Since(startTime)
	result.CollectionTime = duration

	// Display focused results
	displayGraphDataResults(result)

	// Export focused BloodHound data
	err = exportGraphDataToBloodHound(result)
	if err != nil {
		fmt.Printf("⚠️  Warning: Failed to export BloodHound data: %v\n", err)
	}
}

func collectAllGraphData(ctx context.Context, azClient client.AzureClient) (*azure.GraphDataCollectionResult, error) {
	result := &azure.GraphDataCollectionResult{
		GroupMemberships:    []azure.GroupMembershipData{},
		UserRoleAssignments: []azure.UserRoleData{},
		SignInActivity:      []azure.SignIn{},           // Keep struct but don't populate
		DeviceAccess:        []azure.DeviceAccessData{}, // Keep struct but don't populate
		Errors:              []string{},
	}

	// Collect Group Memberships (focused on relevant groups)
	fmt.Printf("🏢 Collecting Azure AD groups and memberships...\n")
	groupResults := azClient.CollectGroupMembershipData(ctx)
	for groupResult := range groupResults {
		if groupResult.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Group error: %v", groupResult.Error))
		} else {
			// Only include privileged groups or groups with members
			if isPrivilegedGroup(groupResult.Ok.Group.DisplayName) || len(groupResult.Ok.Members) > 0 {
				result.GroupMemberships = append(result.GroupMemberships, groupResult.Ok)
				result.TotalGroups++
			}
		}
	}
	fmt.Printf("   ✓ Collected %d relevant groups\n", result.TotalGroups)

	// Collect User Role Assignments (focused on users with roles)
	fmt.Printf("👤 Collecting user role assignments...\n")
	userResults := azClient.CollectUserRoleAssignments(ctx)
	for userResult := range userResults {
		if userResult.Error != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("User error: %v", userResult.Error))
		} else {
			// Only include users with role assignments
			if len(userResult.Ok.RoleAssignments) > 0 {
				result.UserRoleAssignments = append(result.UserRoleAssignments, userResult.Ok)
				result.TotalUsers++
			}
		}
	}
	fmt.Printf("   ✓ Collected role assignments for %d users\n", result.TotalUsers)

	// Skip device and sign-in collection for focused approach
	fmt.Printf("ℹ️  Skipping device and sign-in data (focused collection)\n")

	return result, nil
}

func displayGraphDataResults(result *azure.GraphDataCollectionResult) {
	fmt.Printf("\n=== AZURE AD DATA COLLECTION RESULTS ===\n")
	fmt.Printf("⏱️  Collection Time: %v\n", result.CollectionTime)
	fmt.Printf("📊 Data Summary:\n")
	fmt.Printf("   • Relevant Groups: %d\n", result.TotalGroups)
	fmt.Printf("   • Users with Roles: %d\n", result.TotalUsers)
	fmt.Printf("   • Errors: %d\n", len(result.Errors))
	fmt.Printf("\n")

	// Group Analysis
	if len(result.GroupMemberships) > 0 {
		fmt.Printf("🏢 GROUP MEMBERSHIP ANALYSIS:\n")

		totalMembers := 0
		privilegedGroups := 0

		for _, group := range result.GroupMemberships {
			totalMembers += len(group.Members)

			// Check for privileged groups
			groupName := group.Group.DisplayName
			if isPrivilegedGroup(groupName) {
				privilegedGroups++
				fmt.Printf("   🔴 Privileged Group: %s (%d members)\n", groupName, len(group.Members))
			}
		}

		fmt.Printf("   • Total Group Memberships: %d\n", totalMembers)
		fmt.Printf("   • Privileged Groups Found: %d\n", privilegedGroups)
		fmt.Printf("\n")
	}

	// User Rights Analysis
	if len(result.UserRoleAssignments) > 0 {
		fmt.Printf("👤 USER RIGHTS ANALYSIS:\n")

		totalAssignments := 0
		privilegedUsers := 0

		for _, user := range result.UserRoleAssignments {
			totalAssignments += len(user.RoleAssignments)

			if hasPrivilegedRoles(user.RoleAssignments) {
				privilegedUsers++
				fmt.Printf("   🔴 Privileged User: %s (%d roles)\n",
					user.User.DisplayName, len(user.RoleAssignments))
			}
		}

		fmt.Printf("   • Total Role Assignments: %d\n", totalAssignments)
		fmt.Printf("   • Users with Privileged Roles: %d\n", privilegedUsers)
		fmt.Printf("\n")
	}

	// Error Summary
	if len(result.Errors) > 0 {
		fmt.Printf("⚠️  ERRORS ENCOUNTERED:\n")
		for i, err := range result.Errors {
			if i >= 10 { // Limit error display
				fmt.Printf("   ... and %d more errors\n", len(result.Errors)-10)
				break
			}
			fmt.Printf("   • %s\n", err)
		}
		fmt.Printf("\n")
	}
}

func exportGraphDataToBloodHound(result *azure.GraphDataCollectionResult) error {
	fmt.Printf("📥 Exporting data to BloodHound format...\n")

	bloodhoundData := convertToBloodHoundFormat(result)

	// Write to file
	filename := fmt.Sprintf("bloodhound_azuread_focused_%s.json", time.Now().Format("20060102_150405"))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(bloodhoundData); err != nil {
		return fmt.Errorf("failed to write BloodHound data: %w", err)
	}

	fmt.Printf("✓ Focused BloodHound data exported to: %s\n", filename)
	fmt.Printf("   • Group Memberships: %d\n", len(bloodhoundData.GroupMemberships))
	fmt.Printf("   • Role Assignments: %d\n", len(bloodhoundData.UserRoleAssignments))

	return nil
}

func convertToBloodHoundFormat(result *azure.GraphDataCollectionResult) azure.BloodHoundGraphData {
	bloodhoundData := azure.BloodHoundGraphData{
		Meta: azure.BloodHoundMeta{
			Type:        "azuread-focused",
			Count:       result.TotalGroups + result.TotalUsers,
			Version:     "1.0",
			Methods:     2, // Groups + Users only (focused)
			CollectedBy: "AzureHound-Focused",
			CollectedAt: time.Now(),
		},
		GroupMemberships:    []azure.BloodHoundGroupMembership{},
		UserRoleAssignments: []azure.BloodHoundUserRoleAssignment{},
		DeviceOwnerships:    []azure.BloodHoundDeviceOwnership{}, // Keep empty
		SignInActivity:      []azure.BloodHoundSignInActivity{},  // Keep empty
	}

	// Convert Groups and Memberships only
	for _, groupData := range result.GroupMemberships {
		// Convert memberships
		for _, memberRaw := range groupData.Members {
			var member map[string]interface{}
			json.Unmarshal(memberRaw, &member)

			if memberID, ok := member["id"].(string); ok {
				memberName := ""
				memberType := "User"

				if displayName, ok := member["displayName"].(string); ok {
					memberName = displayName
				}
				if odataType, ok := member["@odata.type"].(string); ok {
					memberType = extractTypeFromOData(odataType)
				}

				membership := azure.BloodHoundGroupMembership{
					GroupId:          groupData.Group.Id,
					GroupName:        groupData.Group.DisplayName,
					MemberId:         memberID,
					MemberName:       memberName,
					MemberType:       memberType,
					RelationshipType: "MemberOf",
				}
				bloodhoundData.GroupMemberships = append(bloodhoundData.GroupMemberships, membership)
			}
		}

		// Convert owners
		for _, ownerRaw := range groupData.Owners {
			var owner map[string]interface{}
			json.Unmarshal(ownerRaw, &owner)

			if ownerID, ok := owner["id"].(string); ok {
				ownerName := ""
				ownerType := "User"

				if displayName, ok := owner["displayName"].(string); ok {
					ownerName = displayName
				}
				if odataType, ok := owner["@odata.type"].(string); ok {
					ownerType = extractTypeFromOData(odataType)
				}

				ownership := azure.BloodHoundGroupMembership{
					GroupId:          groupData.Group.Id,
					GroupName:        groupData.Group.DisplayName,
					MemberId:         ownerID,
					MemberName:       ownerName,
					MemberType:       ownerType,
					RelationshipType: "OwnerOf",
				}
				bloodhoundData.GroupMemberships = append(bloodhoundData.GroupMemberships, ownership)
			}
		}
	}

	// Convert User Role Assignments only
	for _, userData := range result.UserRoleAssignments {
		for _, assignment := range userData.RoleAssignments {
			roleAssignment := azure.BloodHoundUserRoleAssignment{
				UserId:          userData.User.Id,
				UserName:        userData.User.DisplayName,
				RoleId:          assignment.AppRoleId.String(),
				RoleName:        assignment.PrincipalDisplayName,
				ResourceId:      assignment.ResourceId,
				ResourceName:    assignment.ResourceDisplayName,
				AssignmentType:  "DirectAssignment",
				CreatedDateTime: time.Now(),
			}
			bloodhoundData.UserRoleAssignments = append(bloodhoundData.UserRoleAssignments, roleAssignment)
		}
	}

	// Skip device and sign-in conversion (focused approach)
	// Set data wrapper to empty since we're only exporting relationships
	bloodhoundData.Data = azure.BloodHoundGraphDataWrapper{
		Users:   []azure.BloodHoundUser{},   // Empty - focus on relationships
		Groups:  []azure.BloodHoundGroup{},  // Empty - focus on relationships
		Devices: []azure.BloodHoundDevice{}, // Empty - not collected
	}

	return bloodhoundData
}

// Helper functions with improved privileged group detection
func isPrivilegedGroup(groupName string) bool {
	privilegedGroups := []string{
		"Global Administrator",
		"Privileged Role Administrator",
		"Security Administrator",
		"User Administrator",
		"Exchange Administrator",
		"SharePoint Administrator",
		"Application Administrator",
		"Cloud Application Administrator",
		"Authentication Administrator",
		"Privileged Authentication Administrator",
		"Domain Admins",
		"Enterprise Admins",
		"Schema Admins",
		"Administrators",
	}

	normalizedGroupName := strings.ToLower(strings.TrimSpace(groupName))

	for _, privileged := range privilegedGroups {
		normalizedPrivileged := strings.ToLower(privileged)

		// Use exact match or word boundary check to avoid false positives
		if normalizedGroupName == normalizedPrivileged {
			return true
		}

		// Check for word boundaries to avoid matching "Non-Administrator" to "Administrator"
		words := strings.Fields(normalizedGroupName)
		privilegedWords := strings.Fields(normalizedPrivileged)

		// If all privileged words are found as complete words in group name
		if containsAllWords(words, privilegedWords) {
			return true
		}
	}
	return false
}

// containsAllWords checks if all words in 'privileged' exist as complete words in 'groupWords'
func containsAllWords(groupWords, privilegedWords []string) bool {
	if len(privilegedWords) == 0 {
		return false
	}

	for _, privilegedWord := range privilegedWords {
		found := false
		for _, groupWord := range groupWords {
			if groupWord == privilegedWord {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func hasPrivilegedRoles(assignments []azure.AppRoleAssignment) bool {
	privilegedRoles := []string{
		"Global Administrator",
		"Privileged Role Administrator",
		"Security Administrator",
		"User Administrator",
		"Directory.ReadWrite.All",
		"RoleManagement.ReadWrite.Directory",
		"Application.ReadWrite.All",
	}

	for _, assignment := range assignments {
		// Use available fields from AppRoleAssignment
		assignmentName := assignment.PrincipalDisplayName
		resourceName := assignment.ResourceDisplayName

		for _, privileged := range privilegedRoles {
			if strings.Contains(strings.ToLower(assignmentName), strings.ToLower(privileged)) ||
				strings.Contains(strings.ToLower(resourceName), strings.ToLower(privileged)) {
				return true
			}
		}
	}
	return false
}

func extractTypeFromOData(odataType string) string {
	if strings.Contains(odataType, "user") {
		return "User"
	} else if strings.Contains(odataType, "group") {
		return "Group"
	} else if strings.Contains(odataType, "servicePrincipal") {
		return "ServicePrincipal"
	} else if strings.Contains(odataType, "application") {
		return "Application"
	}
	return "Unknown"
}

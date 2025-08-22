// Updated implementation for client/user_roles.go

package client

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// Package-level constants for admin role names
var adminRoleNames = map[string]bool{
	"global administrator":                    true,
	"privileged role administrator":           true,
	"security administrator":                  true,
	"user administrator":                      true,
	"helpdesk administrator":                  true,
	"exchange administrator":                  true,
	"sharepoint administrator":                true,
	"teams administrator":                     true,
	"intune administrator":                    true,
	"cloud application administrator":         true,
	"application administrator":               true,
	"authentication administrator":            true,
	"privileged authentication administrator": true,
	"directory readers":                       true,
	"directory writers":                       true,
}

// MemberOfObject represents objects returned by memberOf API
type MemberOfObject struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	ODataType   string `json:"@odata.type"`
}

// GetUserMemberOf gets all groups and roles a user belongs to
func (s *azureClient) GetUserMemberOf(ctx context.Context, userPrincipalName string) ([]azure.DirectoryObject, error) {
	// URL encode the userPrincipalName to handle special characters
	encodedUPN := url.QueryEscape(userPrincipalName)
	path := fmt.Sprintf("/v1.0/users/%s/memberOf", encodedUPN)

	params := query.GraphParams{
		Select: []string{"id", "displayName", "@odata.type"},
	}

	resultChan := make(chan AzureResult[MemberOfObject])
	go getAzureObjectList[MemberOfObject](s.msgraph, ctx, path, params, resultChan)

	var objects []azure.DirectoryObject
	for result := range resultChan {
		if result.Error != nil {
			return nil, result.Error
		}

		// Convert to our model
		obj := azure.DirectoryObject{
			ID:          result.Ok.ID,
			DisplayName: result.Ok.DisplayName,
			ObjectType:  result.Ok.ODataType,
		}
		objects = append(objects, obj)
	}

	return objects, nil
}

// Enhanced role detection using memberOf API
func (s *azureClient) isAdminUserEnhanced(ctx context.Context, userPrincipalName string) bool {
	memberOf, err := s.GetUserMemberOf(ctx, userPrincipalName)
	if err != nil {
		// Fall back to heuristic if API fails
		return isAdminUserHeuristic(userPrincipalName)
	}

	adminGroupPatterns := []string{
		"admin", "administrator", "domain admin", "enterprise admin",
		"schema admin", "global admin", "tenant admin", "it admin",
		"system admin", "security admin", "compliance admin",
	}

	for _, obj := range memberOf {
		displayName := strings.ToLower(obj.DisplayName)

		// Check if it's a directory role (admin role)
		if obj.ObjectType == "#microsoft.graph.directoryRole" {
			if adminRoleNames[displayName] {
				return true
			}
		}

		// Check if it's an admin group
		if obj.ObjectType == "#microsoft.graph.group" {
			for _, pattern := range adminGroupPatterns {
				if strings.Contains(displayName, pattern) {
					return true
				}
			}
		}
	}

	return false
}

func (s *azureClient) isServiceUserEnhanced(ctx context.Context, userPrincipalName string) bool {
	memberOf, err := s.GetUserMemberOf(ctx, userPrincipalName)
	if err != nil {
		// Fall back to heuristic if API fails
		return isServiceUserHeuristic(userPrincipalName)
	}

	// Service account group patterns
	serviceGroupPatterns := []string{
		"service", "svc", "application", "app", "automation",
		"system", "daemon", "bot", "sync", "backup", "monitoring",
		"service account", "application account", "system account",
	}

	for _, obj := range memberOf {
		if obj.ObjectType == "#microsoft.graph.group" {
			groupName := strings.ToLower(obj.DisplayName)
			for _, pattern := range serviceGroupPatterns {
				if strings.Contains(groupName, pattern) {
					return true
				}
			}
		}
	}

	// Also check the UPN itself for service patterns
	return isServiceUserHeuristic(userPrincipalName)
}

// Alternative: Get only directory roles (more specific for admin detection)
func (s *azureClient) GetUserDirectoryRoles(ctx context.Context, userPrincipalName string) ([]azure.DirectoryRole, error) {
	// URL encode the userPrincipalName to handle special characters
	encodedUPN := url.QueryEscape(userPrincipalName)
	path := fmt.Sprintf("/v1.0/users/%s/memberOf", encodedUPN)

	params := query.GraphParams{
		Filter: "@odata.type eq 'microsoft.graph.directoryRole'",
		Select: []string{"id", "displayName", "description"},
	}

	resultChan := make(chan AzureResult[azure.DirectoryRole])
	go getAzureObjectList[azure.DirectoryRole](s.msgraph, ctx, path, params, resultChan)

	var roles []azure.DirectoryRole
	for result := range resultChan {
		if result.Error != nil {
			return nil, result.Error
		}
		roles = append(roles, result.Ok)
	}

	return roles, nil
}

// Helper functions for role detection
func isAdminUserHeuristic(upn string) bool {
	if upn == "" {
		return false
	}

	lower := strings.ToLower(upn)
	adminPatterns := []string{
		"admin", "administrator", "root", "sysadmin", "systemadmin",
		"domain-admin", "domainadmin", "global-admin", "globaladmin",
		"tenant-admin", "it-admin", "adm-", "-adm",
	}

	for _, pattern := range adminPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	return false
}

func isServiceUserHeuristic(upn string) bool {
	if upn == "" {
		return false
	}

	lower := strings.ToLower(upn)

	// More comprehensive service account patterns
	servicePatterns := []string{
		"service", "svc", "srv", "system", "daemon", "app-",
		"application", "azure-", "microsoft", "msonline", "sync_",
		"exchange", "sharepoint", "teams", "bot", "automation",
		"backup", "monitoring",
	}

	for _, pattern := range servicePatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	// Check for machine account suffix or GUID-like names
	return strings.HasSuffix(lower, "$") || isGUIDLike(upn)
}

// Helper function for GUID detection
func isGUIDLike(s string) bool {
	// Basic GUID pattern check
	if len(s) != 36 {
		return false
	}

	// Check for GUID format: 8-4-4-4-12
	parts := strings.Split(s, "-")
	if len(parts) != 5 {
		return false
	}

	expectedLengths := []int{8, 4, 4, 4, 12}
	for i, part := range parts {
		if len(part) != expectedLengths[i] {
			return false
		}
		// Check if all characters are hexadecimal
		for _, char := range part {
			if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f') || (char >= 'A' && char <= 'F')) {
				return false
			}
		}
	}

	return true
}

// client/intune_groups_direct.go
package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

// ListUserAppRoleAssignments - Get app role assignments for a user (User Rights)
func (s *azureClient) ListUserAppRoleAssignments(ctx context.Context, userID string, params query.GraphParams) <-chan AzureResult[azure.AppRoleAssignment] {
	var (
		out  = make(chan AzureResult[azure.AppRoleAssignment])
		path = fmt.Sprintf("/%s/users/%s/appRoleAssignments", constants.GraphApiVersion, userID)
	)

	if params.Top == 0 {
		params.Top = 999
	}

	go getAzureObjectList[azure.AppRoleAssignment](s.msgraph, ctx, path, params, out)
	return out
}

// ListSignIns - Get sign-in activity (for active sessions context)
func (s *azureClient) ListSignIns(ctx context.Context, params query.GraphParams) <-chan AzureResult[azure.SignIn] {
	var (
		out  = make(chan AzureResult[azure.SignIn])
		path = fmt.Sprintf("/%s/auditLogs/signIns", constants.GraphApiVersion)
	)

	if params.Top == 0 {
		params.Top = 100 // Sign-ins can be large datasets
	}

	go getAzureObjectList[azure.SignIn](s.msgraph, ctx, path, params, out)
	return out
}

// GetDeviceRegisteredUsers - Get users registered to a device
func (s *azureClient) GetDeviceRegisteredUsers(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
	var (
		out  = make(chan AzureResult[json.RawMessage])
		path = fmt.Sprintf("/%s/devices/%s/registeredUsers", constants.GraphApiVersion, deviceId)
	)

	go getAzureObjectList[json.RawMessage](s.msgraph, ctx, path, params, out)
	return out
}

// GetDeviceRegisteredOwners - Get owners of a device
func (s *azureClient) GetDeviceRegisteredOwners(ctx context.Context, deviceId string, params query.GraphParams) <-chan AzureResult[json.RawMessage] {
	var (
		out  = make(chan AzureResult[json.RawMessage])
		path = fmt.Sprintf("/%s/devices/%s/registeredOwners", constants.GraphApiVersion, deviceId)
	)

	go getAzureObjectList[json.RawMessage](s.msgraph, ctx, path, params, out)
	return out
}

// CollectGroupMembershipData - Collect all group membership data using existing methods
func (s *azureClient) CollectGroupMembershipData(ctx context.Context) <-chan AzureResult[azure.GroupMembershipData] {
	out := make(chan AzureResult[azure.GroupMembershipData])

	go func() {
		defer close(out)

		// Use existing ListAzureADGroups method
		groups := s.ListAzureADGroups(ctx, query.GraphParams{})

		for groupResult := range groups {
			if groupResult.Error != nil {
				out <- AzureResult[azure.GroupMembershipData]{Error: groupResult.Error}
				continue
			}

			group := groupResult.Ok

			// Use existing ListAzureADGroupMembers method - fix field name
			members := s.ListAzureADGroupMembers(ctx, group.Id, query.GraphParams{})

			var membersList []json.RawMessage
			for memberResult := range members {
				if memberResult.Error != nil {
					continue // Skip individual member errors
				}
				membersList = append(membersList, memberResult.Ok)
			}

			// Use existing ListAzureADGroupOwners method - fix field name
			owners := s.ListAzureADGroupOwners(ctx, group.Id, query.GraphParams{})

			var ownersList []json.RawMessage
			for ownerResult := range owners {
				if ownerResult.Error != nil {
					continue // Skip individual owner errors
				}
				ownersList = append(ownersList, ownerResult.Ok)
			}

			groupData := azure.GroupMembershipData{
				Group:   group,
				Members: membersList,
				Owners:  ownersList,
			}

			out <- AzureResult[azure.GroupMembershipData]{Ok: groupData}
		}
	}()

	return out
}

// CollectUserRoleAssignments - Collect user rights assignments from Graph API
func (s *azureClient) CollectUserRoleAssignments(ctx context.Context) <-chan AzureResult[azure.UserRoleData] {
	out := make(chan AzureResult[azure.UserRoleData])

	go func() {
		defer close(out)

		// Use existing ListAzureADUsers method
		users := s.ListAzureADUsers(ctx, query.GraphParams{})

		for userResult := range users {
			if userResult.Error != nil {
				out <- AzureResult[azure.UserRoleData]{Error: userResult.Error}
				continue
			}

			user := userResult.Ok

			// Get app role assignments for this user - fix field name
			roleAssignments := s.ListUserAppRoleAssignments(ctx, user.Id, query.GraphParams{})

			var assignments []azure.AppRoleAssignment
			for assignmentResult := range roleAssignments {
				if assignmentResult.Error != nil {
					continue // Skip individual assignment errors
				}
				assignments = append(assignments, assignmentResult.Ok)
			}

			userData := azure.UserRoleData{
				User:            user,
				RoleAssignments: assignments,
			}

			out <- AzureResult[azure.UserRoleData]{Ok: userData}
		}
	}()

	return out
}

// CollectDeviceAccessData - Collect device access and ownership data
func (s *azureClient) CollectDeviceAccessData(ctx context.Context) <-chan AzureResult[azure.DeviceAccessData] {
	out := make(chan AzureResult[azure.DeviceAccessData])

	go func() {
		defer close(out)

		// Get all Intune devices
		devices := s.ListIntuneDevices(ctx, query.GraphParams{})

		for deviceResult := range devices {
			if deviceResult.Error != nil {
				out <- AzureResult[azure.DeviceAccessData]{Error: deviceResult.Error}
				continue
			}

			device := deviceResult.Ok

			// Try to find corresponding Azure AD device using existing method
			azureDevices := s.ListAzureDevices(ctx, query.GraphParams{
				Filter: fmt.Sprintf("deviceId eq '%s'", device.AzureADDeviceID),
			})

			var azureDevice *azure.Device
			for azureDeviceResult := range azureDevices {
				if azureDeviceResult.Error == nil {
					deviceData := azureDeviceResult.Ok
					azureDevice = &deviceData
					break
				}
			}

			var registeredUsers []json.RawMessage
			var registeredOwners []json.RawMessage

			if azureDevice != nil {
				// Get registered users - fix field name
				users := s.GetDeviceRegisteredUsers(ctx, azureDevice.Id, query.GraphParams{})
				for userResult := range users {
					if userResult.Error == nil {
						registeredUsers = append(registeredUsers, userResult.Ok)
					}
				}

				// Get registered owners - fix field name
				owners := s.GetDeviceRegisteredOwners(ctx, azureDevice.Id, query.GraphParams{})
				for ownerResult := range owners {
					if ownerResult.Error == nil {
						registeredOwners = append(registeredOwners, ownerResult.Ok)
					}
				}
			}

			deviceAccessData := azure.DeviceAccessData{
				IntuneDevice:     device,
				AzureDevice:      azureDevice,
				RegisteredUsers:  registeredUsers,
				RegisteredOwners: registeredOwners,
			}

			out <- AzureResult[azure.DeviceAccessData]{Ok: deviceAccessData}
		}
	}()

	return out
}

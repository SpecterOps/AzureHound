// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"context"
	"testing"

	"github.com/bloodhoundad/azurehound/v2/client/mocks"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListResourceGroupUserAccessAdmins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockRoleAssignmentsChannel := make(chan azureWrapper[models.ResourceGroupRoleAssignments])
	mockTenant := azure.Tenant{}
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	channel := listResourceGroupUserAccessAdmins(ctx, mockRoleAssignmentsChannel)

	go func() {
		defer close(mockRoleAssignmentsChannel)

		mockRoleAssignmentsChannel <- NewAzureWrapper(
			enums.KindAZResourceGroupRoleAssignment,
			models.ResourceGroupRoleAssignments{
				ResourceGroupId: "foo",
				RoleAssignments: []models.ResourceGroupRoleAssignment{
					{
						RoleAssignment: azure.RoleAssignment{
							Name: "uaa-assignment",
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.UserAccessAdminRoleID,
							},
						},
					},
					{
						RoleAssignment: azure.RoleAssignment{
							Name: "owner-assignment",
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.OwnerRoleID,
							},
						},
					},
				},
			},
		)
	}()

	result, ok := <-channel
	if !ok {
		t.Fatalf("failed to receive from channel")
	}

	wrapper, ok := result.(azureWrapper[models.ResourceGroupUserAccessAdmins])
	if !ok {
		t.Fatalf("unexpected type in channel: %T", result)
	}

	if wrapper.Data.ResourceGroupId != "foo" {
		t.Errorf("expected ResourceGroupId 'foo', got '%s'", wrapper.Data.ResourceGroupId)
	}

	if len(wrapper.Data.UserAccessAdmins) != 1 {
		t.Fatalf("expected 1 user access admin, got %d", len(wrapper.Data.UserAccessAdmins))
	}

	if wrapper.Data.UserAccessAdmins[0].UserAccessAdmin.Name != "uaa-assignment" {
		t.Errorf("expected user access admin name 'uaa-assignment', got '%s'", wrapper.Data.UserAccessAdmins[0].UserAccessAdmin.Name)
	}

	if _, ok := <-channel; ok {
		t.Error("should not have received from channel")
	}
}

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
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListSubscriptionContributors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockRoleAssignmentsChannel := make(chan interface{})
	mockTenant := azure.Tenant{}
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	channel := listSubscriptionContributors(ctx, mockClient, mockRoleAssignmentsChannel)

	go func() {
		defer close(mockRoleAssignmentsChannel)

		mockRoleAssignmentsChannel <- AzureWrapper{
			Data: models.SubscriptionRoleAssignments{
				SubscriptionId: "foo",
				RoleAssignments: []models.SubscriptionRoleAssignment{
					{
						RoleAssignment: azure.RoleAssignment{
							Name: constants.ContributorRoleID,
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.ContributorRoleID,
							},
						},
					},
				},
			},
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.SubscriptionContributors); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.SubscriptionContributors{})
	} else {
		if data.SubscriptionId != "foo" {
			t.Errorf("got SubscriptionId %q, want %q", data.SubscriptionId, "foo")
		}
		if len(data.Contributors) != 1 {
			t.Fatalf("got %v contributors, want 1", len(data.Contributors))
		}
		if data.Contributors[0].Contributor.Name != constants.ContributorRoleID {
			t.Errorf("got Contributor.Name %q, want %q", data.Contributors[0].Contributor.Name, constants.ContributorRoleID)
		}
		if data.Contributors[0].Contributor.Properties.RoleDefinitionId != constants.ContributorRoleID {
			t.Errorf("got RoleDefinitionId %q, want %q", data.Contributors[0].Contributor.Properties.RoleDefinitionId, constants.ContributorRoleID)
		}
	}

	if _, ok := <-channel; ok {
		t.Error("should not have received from channel")
	}
}

func TestListSubscriptionContributors_Filters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockRoleAssignmentsChannel := make(chan interface{})
	mockTenant := azure.Tenant{}
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	channel := listSubscriptionContributors(ctx, mockClient, mockRoleAssignmentsChannel)

	go func() {
		defer close(mockRoleAssignmentsChannel)

		// Send role assignments with Owner and UserAccessAdmin roles — neither should pass the Contributor filter
		mockRoleAssignmentsChannel <- AzureWrapper{
			Data: models.SubscriptionRoleAssignments{
				SubscriptionId: "foo",
				RoleAssignments: []models.SubscriptionRoleAssignment{
					{
						RoleAssignment: azure.RoleAssignment{
							Name: constants.OwnerRoleID,
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.OwnerRoleID,
							},
						},
					},
					{
						RoleAssignment: azure.RoleAssignment{
							Name: constants.UserAccessAdminRoleID,
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.UserAccessAdminRoleID,
							},
						},
					},
				},
			},
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.SubscriptionContributors); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.SubscriptionContributors{})
	} else if len(data.Contributors) != 0 {
		t.Errorf("got %v contributors, want 0", len(data.Contributors))
	}

	if _, ok := <-channel; ok {
		t.Error("should not have received from channel")
	}
}


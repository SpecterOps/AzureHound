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

func TestListSubscriptionOwners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockRoleAssignmentsChannel := make(chan interface{})
	mockTenant := azure.Tenant{}
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	channel := listSubscriptionOwners(ctx, mockClient, mockRoleAssignmentsChannel)

	go func() {
		defer close(mockRoleAssignmentsChannel)

		mockRoleAssignmentsChannel <- AzureWrapper{
			Data: models.SubscriptionRoleAssignments{
				SubscriptionId: "foo",
				RoleAssignments: []models.SubscriptionRoleAssignment{
					{
						RoleAssignment: azure.RoleAssignment{
							Name: "matching-assignment",
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.OwnerRoleID,
							},
						},
					},
					{
						RoleAssignment: azure.RoleAssignment{
							Name: "non-matching-assignment",
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.ContributorRoleID,
							},
						},
					},
				},
			},
		}
	}()

	result, ok := <-channel
	if !ok {
		t.Fatalf("failed to receive from channel")
	}

	wrapper, ok := result.(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	}

	data, ok := wrapper.Data.(models.SubscriptionOwners)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", wrapper.Data, models.SubscriptionOwners{})
	}

	if data.SubscriptionId != "foo" {
		t.Errorf("expected SubscriptionId 'foo', got '%s'", data.SubscriptionId)
	}

	if len(data.Owners) != 1 {
		t.Fatalf("expected 1 owner, got %d", len(data.Owners))
	}

	if data.Owners[0].Owner.Name != "matching-assignment" {
		t.Errorf("expected owner name 'matching-assignment', got '%s'", data.Owners[0].Owner.Name)
	}

	if _, ok := <-channel; ok {
		t.Error("should not have received from channel")
	}
}

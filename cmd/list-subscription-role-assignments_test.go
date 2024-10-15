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
	"fmt"
	"testing"

	"github.com/certmichelin/azurehound/v2/client"
	"github.com/certmichelin/azurehound/v2/client/mocks"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListSubscriptionRoleAssignments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockSubscriptionsChannel := make(chan interface{})
	mockSubscriptionRoleAssignmentChannel := make(chan client.AzureResult[azure.RoleAssignment])
	mockSubscriptionRoleAssignmentChannel2 := make(chan client.AzureResult[azure.RoleAssignment])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListRoleAssignmentsForResource(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockSubscriptionRoleAssignmentChannel).Times(1)
	mockClient.EXPECT().ListRoleAssignmentsForResource(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockSubscriptionRoleAssignmentChannel2).Times(1)
	channel := listSubscriptionRoleAssignments(ctx, mockClient, mockSubscriptionsChannel)

	go func() {
		defer close(mockSubscriptionsChannel)
		mockSubscriptionsChannel <- AzureWrapper{
			Data: models.Subscription{},
		}
		mockSubscriptionsChannel <- AzureWrapper{
			Data: models.Subscription{},
		}
	}()
	go func() {
		defer close(mockSubscriptionRoleAssignmentChannel)
		mockSubscriptionRoleAssignmentChannel <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.ContributorRoleID,
				},
			},
		}
		mockSubscriptionRoleAssignmentChannel <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.OwnerRoleID,
				},
			},
		}
	}()
	go func() {
		defer close(mockSubscriptionRoleAssignmentChannel2)
		mockSubscriptionRoleAssignmentChannel2 <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.OwnerRoleID,
				},
			},
		}
		mockSubscriptionRoleAssignmentChannel2 <- client.AzureResult[azure.RoleAssignment]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.SubscriptionRoleAssignments); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.SubscriptionRoleAssignments{})
	} else if len(data.RoleAssignments) != 2 {
		t.Errorf("got %v, want %v", len(data.RoleAssignments), 2)
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.SubscriptionRoleAssignments); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.SubscriptionRoleAssignments{})
	} else if len(data.RoleAssignments) != 1 {
		t.Errorf("got %v, want %v", len(data.RoleAssignments), 2)
	}
}

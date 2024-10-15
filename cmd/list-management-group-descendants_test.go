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
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListManagementGroupDescendants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockManagementGroupsChannel := make(chan interface{})
	mockManagementGroupDescendantChannel := make(chan client.AzureResult[azure.DescendantInfo])
	mockManagementGroupDescendantChannel2 := make(chan client.AzureResult[azure.DescendantInfo])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureManagementGroupDescendants(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockManagementGroupDescendantChannel).Times(1)
	mockClient.EXPECT().ListAzureManagementGroupDescendants(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockManagementGroupDescendantChannel2).Times(1)
	channel := listManagementGroupDescendants(ctx, mockClient, mockManagementGroupsChannel)

	go func() {
		defer close(mockManagementGroupsChannel)
		mockManagementGroupsChannel <- AzureWrapper{
			Data: models.ManagementGroup{},
		}
		mockManagementGroupsChannel <- AzureWrapper{
			Data: models.ManagementGroup{},
		}
	}()
	go func() {
		defer close(mockManagementGroupDescendantChannel)
		mockManagementGroupDescendantChannel <- client.AzureResult[azure.DescendantInfo]{}
		mockManagementGroupDescendantChannel <- client.AzureResult[azure.DescendantInfo]{}
	}()
	go func() {
		defer close(mockManagementGroupDescendantChannel2)
		mockManagementGroupDescendantChannel2 <- client.AzureResult[azure.DescendantInfo]{}
		mockManagementGroupDescendantChannel2 <- client.AzureResult[azure.DescendantInfo]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(azure.DescendantInfo); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, azure.DescendantInfo{})
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(azure.DescendantInfo); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, azure.DescendantInfo{})
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(azure.DescendantInfo); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, azure.DescendantInfo{})
	}

	if _, ok := <-channel; ok {
		t.Error("expected channel to close from an error result but it did not")
	}
}

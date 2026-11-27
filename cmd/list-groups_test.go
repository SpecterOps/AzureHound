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

	"github.com/SpecterOps/AzureHound/v2/client"
	"github.com/SpecterOps/AzureHound/v2/client/mocks"
	"github.com/SpecterOps/AzureHound/v2/models"
	"github.com/SpecterOps/AzureHound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.Group])
	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroups(gomock.Any(), gomock.Any()).Return(mockChannel)

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.Group]{
			Ok: azure.Group{},
		}
		mockChannel <- client.AzureResult[azure.Group]{
			Error: mockError,
		}
		mockChannel <- client.AzureResult[azure.Group]{
			Ok: azure.Group{},
		}
	}()

	channel := listGroups(ctx, mockClient)
	result := <-channel
	if _, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	}

	if _, ok := <-channel; ok {
		t.Error("expected channel to close from an error result but it did not")
	}
}

func TestListGroupsCloudManaged(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.Group])
	mockTenant := azure.Tenant{}
	cloudManaged := true

	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroups(gomock.Any(), gomock.Any()).Return(mockChannel)
	// The group carries an on-premises SID, so collection resolves its onPremisesSyncBehavior exactly once.
	mockClient.EXPECT().
		GetAzureADGroupOnPremisesSyncBehavior(gomock.Any(), "group-with-sid").
		Return(&azure.OnPremisesSyncBehavior{IsCloudManaged: &cloudManaged}, nil).
		Times(1)

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.Group]{
			Ok: azure.Group{
				DirectoryObject:              azure.DirectoryObject{Id: "group-with-sid"},
				OnPremisesSecurityIdentifier: "S-1-5-21-111-222-333-1369",
			},
		}
	}()

	channel := listGroups(ctx, mockClient)
	wrapper, ok := (<-channel).(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got non-AzureWrapper result")
	}
	group, ok := wrapper.Data.(models.Group)
	if !ok {
		t.Fatalf("failed type assertion on Data: got %T, want models.Group", wrapper.Data)
	}
	if group.IsCloudManaged == nil || !*group.IsCloudManaged {
		t.Errorf("expected IsCloudManaged=true, got %v", group.IsCloudManaged)
	}
}

// A group without an on-premises SID can never be a writeback target, so collection must NOT spend a
// per-group onPremisesSyncBehavior call on it. The strict mock controller fails the test if the call happens
// (no EXPECT registered), and the emitted group must leave IsCloudManaged unresolved (nil).
func TestListGroupsCloudManagedSkippedWithoutSID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.Group])
	mockTenant := azure.Tenant{}

	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroups(gomock.Any(), gomock.Any()).Return(mockChannel)
	// No GetAzureADGroupOnPremisesSyncBehavior EXPECT: any such call is an unexpected-call failure.

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.Group]{
			Ok: azure.Group{DirectoryObject: azure.DirectoryObject{Id: "cloud-only-group"}},
		}
	}()

	channel := listGroups(ctx, mockClient)
	wrapper, ok := (<-channel).(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got non-AzureWrapper result")
	}
	group, ok := wrapper.Data.(models.Group)
	if !ok {
		t.Fatalf("failed type assertion on Data: got %T, want models.Group", wrapper.Data)
	}
	if group.IsCloudManaged != nil {
		t.Errorf("expected IsCloudManaged=nil for a group without an on-premises SID, got %v", *group.IsCloudManaged)
	}
}

// If resolving onPremisesSyncBehavior fails (e.g. the beta endpoint 404s or the token lacks the permission),
// collection must tolerate it: the group is still emitted, just with IsCloudManaged unresolved (nil).
func TestListGroupsCloudManagedErrorTolerated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.Group])
	mockTenant := azure.Tenant{}

	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroups(gomock.Any(), gomock.Any()).Return(mockChannel)
	mockClient.EXPECT().
		GetAzureADGroupOnPremisesSyncBehavior(gomock.Any(), "group-with-sid").
		Return(nil, fmt.Errorf("insufficient privileges")).
		Times(1)

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.Group]{
			Ok: azure.Group{
				DirectoryObject:              azure.DirectoryObject{Id: "group-with-sid"},
				OnPremisesSecurityIdentifier: "S-1-5-21-111-222-333-1369",
			},
		}
	}()

	channel := listGroups(ctx, mockClient)
	wrapper, ok := (<-channel).(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got non-AzureWrapper result")
	}
	group, ok := wrapper.Data.(models.Group)
	if !ok {
		t.Fatalf("failed type assertion on Data: got %T, want models.Group", wrapper.Data)
	}
	if group.IsCloudManaged != nil {
		t.Errorf("expected IsCloudManaged=nil when resolution errors, got %v", *group.IsCloudManaged)
	}
}

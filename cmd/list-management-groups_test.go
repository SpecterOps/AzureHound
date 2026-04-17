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

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/mocks"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListManagementGroups(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	homeTenantId := "home-tenant-id"
	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.ManagementGroup])
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(azure.Tenant{TenantId: homeTenantId}).AnyTimes()
	mockClient.EXPECT().ListAzureManagementGroups(gomock.Any(), gomock.Any()).Return(mockChannel)

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Ok: azure.ManagementGroup{
				Properties: azure.ManagementGroupProperties{TenantId: homeTenantId},
			},
		}
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Error: mockError,
		}
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Ok: azure.ManagementGroup{
				Properties: azure.ManagementGroupProperties{TenantId: homeTenantId},
			},
		}
	}()

	channel := listManagementGroups(ctx, mockClient)
	result := <-channel
	if _, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	}

	if _, ok := <-channel; ok {
		t.Error("expected channel to close from an error result but it did not")
	}
}

func TestListManagementGroups_FiltersForeignTenants(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	const (
		homeTenantId    = "home-tenant-aaaa-bbbb-cccc"
		foreignTenantId = "foreign-tenant-dddd-eeee-ffff"
	)

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.ManagementGroup])
	mockClient.EXPECT().TenantInfo().Return(azure.Tenant{TenantId: homeTenantId}).AnyTimes()
	mockClient.EXPECT().ListAzureManagementGroups(gomock.Any(), gomock.Any()).Return(mockChannel)

	go func() {
		defer close(mockChannel)
		// Management group belonging to the home tenant — should be collected
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Ok: azure.ManagementGroup{
				Entity: azure.Entity{Id: "/providers/Microsoft.Management/managementGroups/HomeMG"},
				Name:   "HomeMG",
				Properties: azure.ManagementGroupProperties{
					TenantId:    homeTenantId,
					DisplayName: "Home Management Group",
				},
			},
		}
		// Management group belonging to a foreign tenant — should be filtered out
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Ok: azure.ManagementGroup{
				Entity: azure.Entity{Id: "/providers/Microsoft.Management/managementGroups/ForeignMG"},
				Name:   "ForeignMG",
				Properties: azure.ManagementGroupProperties{
					TenantId:    foreignTenantId,
					DisplayName: "Foreign Management Group",
				},
			},
		}
		// Another home tenant management group — should be collected
		mockChannel <- client.AzureResult[azure.ManagementGroup]{
			Ok: azure.ManagementGroup{
				Entity: azure.Entity{Id: "/providers/Microsoft.Management/managementGroups/HomeMG2"},
				Name:   "HomeMG2",
				Properties: azure.ManagementGroupProperties{
					TenantId:    homeTenantId,
					DisplayName: "Home Management Group 2",
				},
			},
		}
	}()

	channel := listManagementGroups(ctx, mockClient)

	var results []models.ManagementGroup
	for item := range channel {
		wrapper := item.(AzureWrapper)
		mg := wrapper.Data.(models.ManagementGroup)
		results = append(results, mg)
	}

	if len(results) != 2 {
		t.Fatalf("expected 2 management groups (home tenant only), got %d", len(results))
	}

	for _, mg := range results {
		if mg.TenantId != homeTenantId {
			t.Errorf("expected all management groups to have tenantId %q, got %q (name: %s)",
				homeTenantId, mg.TenantId, mg.Name)
		}
	}

	if results[0].Name != "HomeMG" {
		t.Errorf("expected first result to be HomeMG, got %s", results[0].Name)
	}
	if results[1].Name != "HomeMG2" {
		t.Errorf("expected second result to be HomeMG2, got %s", results[1].Name)
	}
}

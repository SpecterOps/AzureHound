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
	"sort"
	"testing"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/mocks"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListKeyVaultRoleAssignments(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockKeyVaultsChannel := make(chan interface{})
	mockKeyVaultRoleAssignmentChannel := make(chan client.AzureResult[azure.RoleAssignment])
	mockKeyVaultRoleAssignmentChannel2 := make(chan client.AzureResult[azure.RoleAssignment])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListRoleAssignmentsForResource(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockKeyVaultRoleAssignmentChannel).Times(1)
	mockClient.EXPECT().ListRoleAssignmentsForResource(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(mockKeyVaultRoleAssignmentChannel2).Times(1)
	channel := listKeyVaultRoleAssignments(ctx, mockClient, mockKeyVaultsChannel)

	go func() {
		defer close(mockKeyVaultsChannel)
		mockKeyVaultsChannel <- AzureWrapper{
			Data: models.KeyVault{},
		}
		mockKeyVaultsChannel <- AzureWrapper{
			Data: models.KeyVault{},
		}
	}()
	go func() {
		defer close(mockKeyVaultRoleAssignmentChannel)
		mockKeyVaultRoleAssignmentChannel <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.KeyVaultContributorRoleID,
				},
			},
		}
		mockKeyVaultRoleAssignmentChannel <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.ContributorRoleID,
				},
			},
		}
	}()
	go func() {
		defer close(mockKeyVaultRoleAssignmentChannel2)
		mockKeyVaultRoleAssignmentChannel2 <- client.AzureResult[azure.RoleAssignment]{
			Ok: azure.RoleAssignment{
				Properties: azure.RoleAssignmentPropertiesWithScope{
					RoleDefinitionId: constants.KeyVaultAdministratorRoleID,
				},
			},
		}
		mockKeyVaultRoleAssignmentChannel2 <- client.AzureResult[azure.RoleAssignment]{
			Error: mockError,
		}
	}()

	// Collect both results before asserting. listKeyVaultRoleAssignments fans out across
	// goroutines via pipeline.Demux, so output order is non-deterministic. Asserting
	// positionally produces a flaky test; instead we assert on the set of counts.
	var roleCounts []int
	for i := 0; i < 2; i++ {
		result, ok := <-channel
		if !ok {
			t.Fatalf("failed to receive result %d from channel", i+1)
		}
		roleCounts = append(roleCounts, len(result.Data.RoleAssignments))
	}

	sort.Ints(roleCounts)
	if roleCounts[0] != 1 || roleCounts[1] != 2 {
		t.Errorf("expected role assignment counts [1 2] (in any order), got %v", roleCounts)
	}
}

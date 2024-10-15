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

	"github.com/certmichelin/azurehound/v2/client/mocks"
	"github.com/certmichelin/azurehound/v2/constants"
	"github.com/certmichelin/azurehound/v2/enums"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListKeyVaultKVContributors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockRoleAssignmentsChannel := make(chan azureWrapper[models.KeyVaultRoleAssignments])
	mockTenant := azure.Tenant{}
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	channel := listKeyVaultKVContributors(ctx, mockRoleAssignmentsChannel)

	go func() {
		defer close(mockRoleAssignmentsChannel)

		mockRoleAssignmentsChannel <- NewAzureWrapper(
			enums.KindAZKeyVaultRoleAssignment,
			models.KeyVaultRoleAssignments{
				KeyVaultId: "foo",
				RoleAssignments: []models.KeyVaultRoleAssignment{
					{
						RoleAssignment: azure.RoleAssignment{
							Name: constants.KeyVaultContributorRoleID,
							Properties: azure.RoleAssignmentPropertiesWithScope{
								RoleDefinitionId: constants.KeyVaultContributorRoleID,
							},
						},
					},
				},
			},
		)
	}()

	if _, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	}

	if _, ok := <-channel; ok {
		t.Error("should not have recieved from channel")
	}
}

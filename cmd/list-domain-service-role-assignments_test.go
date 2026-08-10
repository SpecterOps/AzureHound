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

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/mocks"
	"github.com/bloodhoundad/azurehound/v2/constants"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func TestListDomainServiceRoleAssignments(t *testing.T) {
	var (
		ctx                   = context.Background()
		controller            = gomock.NewController(t)
		mockClient            = mocks.NewMockAzureClient(controller)
		domainServices        = make(chan interface{})
		roleAssignmentResults = make(chan client.AzureResult[azure.RoleAssignment])
		resourceID            = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.AAD/domainServices/example.com"
	)

	mockClient.EXPECT().ListRoleAssignmentsForResource(gomock.Any(), resourceID, "atScope()", "").Return(roleAssignmentResults)
	results := listDomainServiceRoleAssignments(ctx, mockClient, domainServices)

	go func() {
		defer close(domainServices)
		domainServices <- AzureWrapper{Data: models.DomainService{DomainService: azure.DomainService{Entity: azure.Entity{Id: resourceID}}}}
	}()
	go func() {
		defer close(roleAssignmentResults)
		roleAssignmentResults <- client.AzureResult[azure.RoleAssignment]{Ok: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId:      "inherited-principal",
				RoleDefinitionId: "/providers/Microsoft.Authorization/roleDefinitions/" + constants.OwnerRoleID,
				Scope:            "/subscriptions/sub/resourceGroups/rg",
			},
		}}
		roleAssignmentResults <- client.AzureResult[azure.RoleAssignment]{Ok: azure.RoleAssignment{
			Properties: azure.RoleAssignmentPropertiesWithScope{
				PrincipalId:      "principal",
				RoleDefinitionId: "/providers/Microsoft.Authorization/roleDefinitions/" + constants.DomainServicesContributorRoleID,
				Scope:            resourceID,
			},
		}}
	}()

	result, ok := <-results
	if !ok {
		t.Fatal("failed to receive role assignments")
	}
	wrapper, ok := result.(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	}
	if wrapper.Kind != enums.KindAZEntraDSRoleAssignment {
		t.Errorf("unexpected kind: %s", wrapper.Kind)
	}
	roleAssignments, ok := wrapper.Data.(models.AzureRoleAssignments)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", wrapper.Data, models.AzureRoleAssignments{})
	}
	if len(roleAssignments.RoleAssignments) != 1 {
		t.Fatalf("expected one role assignment, got %d", len(roleAssignments.RoleAssignments))
	}
	if roleAssignments.RoleAssignments[0].RoleDefinitionId != constants.DomainServicesContributorRoleID {
		t.Errorf("unexpected role definition id: %s", roleAssignments.RoleAssignments[0].RoleDefinitionId)
	}
}

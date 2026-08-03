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
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func TestListDomainServices(t *testing.T) {
	var (
		ctx                  = context.Background()
		controller           = gomock.NewController(t)
		mockClient           = mocks.NewMockAzureClient(controller)
		subscriptions        = make(chan interface{})
		domainServiceResults = make(chan client.AzureResult[azure.DomainService])
		resourceID           = "/subscriptions/sub/resourceGroups/rg/providers/Microsoft.AAD/domainServices/example.com"
	)

	mockClient.EXPECT().ListAzureDomainServices(gomock.Any(), "sub").Return(domainServiceResults)
	results := listDomainServices(ctx, mockClient, subscriptions)

	go func() {
		defer close(subscriptions)
		subscriptions <- AzureWrapper{Data: models.Subscription{Subscription: azure.Subscription{SubscriptionId: "sub"}}}
	}()
	go func() {
		defer close(domainServiceResults)
		domainServiceResults <- client.AzureResult[azure.DomainService]{Ok: azure.DomainService{
			Entity: azure.Entity{Id: resourceID},
			Name:   "example.com",
		}}
	}()

	result, ok := <-results
	if !ok {
		t.Fatal("failed to receive domain service")
	}
	wrapper, ok := result.(AzureWrapper)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	}
	domainService, ok := wrapper.Data.(models.DomainService)
	if !ok {
		t.Fatalf("failed type assertion: got %T, want %T", wrapper.Data, models.DomainService{})
	}
	if domainService.ResourceGroupID != "/subscriptions/sub/resourceGroups/rg" {
		t.Errorf("unexpected resource group id: %s", domainService.ResourceGroupID)
	}
	if _, ok := <-results; ok {
		t.Error("expected domain service stream to close")
	}
}

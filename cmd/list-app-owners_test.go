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
	"encoding/json"
	"fmt"
	"testing"

	"github.com/certmichelin/azurehound/v2/client"
	"github.com/certmichelin/azurehound/v2/client/mocks"
	"github.com/certmichelin/azurehound/v2/enums"
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListAppOwners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockAppsChannel := make(chan azureWrapper[models.App])
	mockAppOwnerChannel := make(chan client.AzureResult[json.RawMessage])
	mockAppOwnerChannel2 := make(chan client.AzureResult[json.RawMessage])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADAppOwners(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAppOwnerChannel).Times(1)
	mockClient.EXPECT().ListAzureADAppOwners(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockAppOwnerChannel2).Times(1)
	channel := listAppOwners(ctx, mockClient, mockAppsChannel)

	go func() {
		defer close(mockAppsChannel)
		mockAppsChannel <- NewAzureWrapper(enums.KindAZApp, models.App{})
		mockAppsChannel <- NewAzureWrapper(enums.KindAZApp, models.App{})
	}()
	go func() {
		defer close(mockAppOwnerChannel)
		mockAppOwnerChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockAppOwnerChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
	}()
	go func() {
		defer close(mockAppOwnerChannel2)
		mockAppOwnerChannel2 <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockAppOwnerChannel2 <- client.AzureResult[json.RawMessage]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if len(result.Data.Owners) != 2 {
		t.Errorf("got %v, want %v", len(result.Data.Owners), 2)
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if len(result.Data.Owners) != 1 {
		t.Errorf("got %v, want %v", len(result.Data.Owners), 2)
	}
}

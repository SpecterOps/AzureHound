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

func TestListKeyVaults(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockSubscriptionsChannel := make(chan interface{})
	mockKeyVaultChannel := make(chan client.AzureResult[azure.KeyVault])
	mockKeyVaultChannel2 := make(chan client.AzureResult[azure.KeyVault])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureKeyVaults(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockKeyVaultChannel).Times(1)
	mockClient.EXPECT().ListAzureKeyVaults(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockKeyVaultChannel2).Times(1)
	channel := listKeyVaults(ctx, mockClient, mockSubscriptionsChannel)

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
		defer close(mockKeyVaultChannel)
		mockKeyVaultChannel <- client.AzureResult[azure.KeyVault]{
			Ok: azure.KeyVault{},
		}
		mockKeyVaultChannel <- client.AzureResult[azure.KeyVault]{
			Ok: azure.KeyVault{},
		}
	}()
	go func() {
		defer close(mockKeyVaultChannel2)
		mockKeyVaultChannel2 <- client.AzureResult[azure.KeyVault]{
			Ok: azure.KeyVault{},
		}
		mockKeyVaultChannel2 <- client.AzureResult[azure.KeyVault]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(models.KeyVault); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.KeyVault{})
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(models.KeyVault); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.KeyVault{})
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if _, ok := wrapper.Data.(models.KeyVault); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.KeyVault{})
	}

	if _, ok := <-channel; ok {
		t.Error("should not have recieved from channel")
	}
}

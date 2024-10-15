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
	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListDeviceOwners(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockDevicesChannel := make(chan interface{})
	mockDeviceOwnerChannel := make(chan client.AzureResult[json.RawMessage])
	mockDeviceOwnerChannel2 := make(chan client.AzureResult[json.RawMessage])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureDeviceRegisteredOwners(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDeviceOwnerChannel).Times(1)
	mockClient.EXPECT().ListAzureDeviceRegisteredOwners(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockDeviceOwnerChannel2).Times(1)
	channel := listDeviceOwners(ctx, mockClient, mockDevicesChannel)

	go func() {
		defer close(mockDevicesChannel)
		mockDevicesChannel <- AzureWrapper{
			Data: models.Device{},
		}
		mockDevicesChannel <- AzureWrapper{
			Data: models.Device{},
		}
	}()
	go func() {
		defer close(mockDeviceOwnerChannel)
		mockDeviceOwnerChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockDeviceOwnerChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
	}()
	go func() {
		defer close(mockDeviceOwnerChannel2)
		mockDeviceOwnerChannel2 <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockDeviceOwnerChannel2 <- client.AzureResult[json.RawMessage]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.DeviceOwners); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.DeviceOwners{})
	} else if len(data.Owners) != 2 {
		t.Errorf("got %v, want %v", len(data.Owners), 2)
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.DeviceOwners); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.DeviceOwners{})
	} else if len(data.Owners) != 1 {
		t.Errorf("got %v, want %v", len(data.Owners), 2)
	}
}

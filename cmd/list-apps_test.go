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
	"github.com/certmichelin/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListApps(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)
	mockChannel := make(chan client.AzureResult[azure.Application])
	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADApps(gomock.Any(), gomock.Any()).Return(mockChannel)

	go func() {
		defer close(mockChannel)
		mockChannel <- client.AzureResult[azure.Application]{
			Ok: azure.Application{},
		}
		mockChannel <- client.AzureResult[azure.Application]{
			Error: mockError,
		}
		mockChannel <- client.AzureResult[azure.Application]{
			Ok: azure.Application{},
		}
	}()

	channel := listApps(ctx, mockClient)
	<-channel
	if _, ok := <-channel; ok {
		t.Error("expected channel to close from an error result but it did not")
	}
}

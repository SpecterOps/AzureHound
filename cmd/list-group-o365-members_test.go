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

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/mocks"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"go.uber.org/mock/gomock"
)

func init() {
	setupLogger()
}

func TestListGroup365Members(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockGroups365Channel := make(chan interface{})
	mockGroup365MemberChannel := make(chan client.AzureResult[json.RawMessage])
	mockGroup365MemberChannel2 := make(chan client.AzureResult[json.RawMessage])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroup365Members(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGroup365MemberChannel).Times(1)
	mockClient.EXPECT().ListAzureADGroup365Members(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGroup365MemberChannel2).Times(1)
	channel := listGroup365Members(ctx, mockClient, mockGroups365Channel)

	go func() {
		defer close(mockGroups365Channel)
		mockGroups365Channel <- AzureWrapper{
			Data: models.Group365{},
		}
		mockGroups365Channel <- AzureWrapper{
			Data: models.Group365{},
		}
	}()
	go func() {
		defer close(mockGroup365MemberChannel)
		mockGroup365MemberChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockGroup365MemberChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
	}()
	go func() {
		defer close(mockGroup365MemberChannel2)
		mockGroup365MemberChannel2 <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockGroup365MemberChannel2 <- client.AzureResult[json.RawMessage]{
			Error: mockError,
		}
	}()

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.Group365Members); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.Group365Members{})
	} else if len(data.Members) != 2 {
		t.Errorf("got %v, want %v", len(data.Members), 2)
	}

	if result, ok := <-channel; !ok {
		t.Fatalf("failed to receive from channel")
	} else if wrapper, ok := result.(AzureWrapper); !ok {
		t.Errorf("failed type assertion: got %T, want %T", result, AzureWrapper{})
	} else if data, ok := wrapper.Data.(models.Group365Members); !ok {
		t.Errorf("failed type assertion: got %T, want %T", wrapper.Data, models.Group365Members{})
	} else if len(data.Members) != 1 {
		t.Errorf("got %v, want %v", len(data.Members), 1)
	}
}

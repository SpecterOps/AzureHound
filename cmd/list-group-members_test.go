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
	"sort"
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

func TestListGroupMembers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.Background()

	mockClient := mocks.NewMockAzureClient(ctrl)

	mockGroupsChannel := make(chan interface{})
	mockGroupMemberChannel := make(chan client.AzureResult[json.RawMessage])
	mockGroupMemberChannel2 := make(chan client.AzureResult[json.RawMessage])

	mockTenant := azure.Tenant{}
	mockError := fmt.Errorf("I'm an error")
	mockClient.EXPECT().TenantInfo().Return(mockTenant).AnyTimes()
	mockClient.EXPECT().ListAzureADGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGroupMemberChannel).Times(1)
	mockClient.EXPECT().ListAzureADGroupMembers(gomock.Any(), gomock.Any(), gomock.Any()).Return(mockGroupMemberChannel2).Times(1)
	channel := listGroupMembers(ctx, mockClient, mockGroupsChannel)

	go func() {
		defer close(mockGroupsChannel)
		mockGroupsChannel <- AzureWrapper{
			Data: models.Group{},
		}
		mockGroupsChannel <- AzureWrapper{
			Data: models.Group{},
		}
	}()
	go func() {
		defer close(mockGroupMemberChannel)
		mockGroupMemberChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockGroupMemberChannel <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
	}()
	go func() {
		defer close(mockGroupMemberChannel2)
		mockGroupMemberChannel2 <- client.AzureResult[json.RawMessage]{
			Ok: json.RawMessage{},
		}
		mockGroupMemberChannel2 <- client.AzureResult[json.RawMessage]{
			Error: mockError,
		}
	}()

	var memberCounts []int
	for i := 0; i < 2; i++ {
		result, ok := <-channel
		if !ok {
			t.Fatalf("failed to receive result %d from channel", i+1)
		}
		wrapper, ok := result.(AzureWrapper)
		if !ok {
			t.Fatalf("result %d: failed type assertion: got %T, want %T", i+1, result, AzureWrapper{})
		}
		data, ok := wrapper.Data.(models.GroupMembers)
		if !ok {
			t.Fatalf("result %d: failed type assertion: got %T, want %T", i+1, wrapper.Data, models.GroupMembers{})
		}
		memberCounts = append(memberCounts, len(data.Members))
	}

	sort.Ints(memberCounts)
	if memberCounts[0] != 1 || memberCounts[1] != 2 {
		t.Errorf("expected member counts [1 2] (in any order), got %v", memberCounts)
	}
}

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

package azure

import "github.com/certmichelin/azurehound/v2/enums"

// Subscription Policies
type SubscriptionPolicies struct {
	// The subscription location placement ID. The ID indicates which regions are visible for a subscription.
	// For example, a subscription with a location placement Id of Public_2014-09-01 has access to Azure public regions.
	LocationPlacementId string `json:"locationPlacementId,omitempty"`

	// The subscription quota ID.
	QuotaId string `json:"quotaId,omitempty"`

	// The subscription spending limit.
	SpendingLimit enums.SpendingLimit `json:"spendingLimit,omitempty"`
}

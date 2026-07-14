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

package models

import (
	"encoding/json"
	"strings"

	"github.com/bloodhoundad/azurehound/v2/models/azure"
)

type SubscriptionContributor struct {
	Contributor    azure.RoleAssignment `json:"contributor"`
	SubscriptionId string               `json:"subscriptionId"`
}

func (s SubscriptionContributor) MarshalJSON() ([]byte, error) {
	type Alias SubscriptionContributor
	a := Alias(s)
	a.SubscriptionId = strings.ToUpper(a.SubscriptionId)
	a.Contributor = UpperRoleAssignment(a.Contributor)
	return json.Marshal(a)
}

type SubscriptionContributors struct {
	Contributors   []SubscriptionContributor `json:"contributors"`
	SubscriptionId string                    `json:"subscriptionId"`
}

func (s SubscriptionContributors) MarshalJSON() ([]byte, error) {
	type Alias SubscriptionContributors
	a := Alias(s)
	a.SubscriptionId = strings.ToUpper(a.SubscriptionId)
	return json.Marshal(a)
}


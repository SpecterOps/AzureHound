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

type SubscriptionUserAccessAdmin struct {
	UserAccessAdmin azure.RoleAssignment `json:"userAccessAdmin"`
	SubscriptionId  string               `json:"subscriptionId"`
}

func (s SubscriptionUserAccessAdmin) MarshalJSON() ([]byte, error) {
	type Alias SubscriptionUserAccessAdmin
	a := Alias(s)
	a.SubscriptionId = strings.ToUpper(a.SubscriptionId)
	a.UserAccessAdmin = UpperRoleAssignment(a.UserAccessAdmin)
	return json.Marshal(a)
}

type SubscriptionUserAccessAdmins struct {
	UserAccessAdmins []SubscriptionUserAccessAdmin `json:"userAccessAdmins"`
	SubscriptionId   string                        `json:"subscriptionId"`
}

func (s SubscriptionUserAccessAdmins) MarshalJSON() ([]byte, error) {
	type Alias SubscriptionUserAccessAdmins
	a := Alias(s)
	a.SubscriptionId = strings.ToUpper(a.SubscriptionId)
	return json.Marshal(a)
}

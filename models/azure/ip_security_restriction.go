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

type IpSecurityRestriction struct {
	Action               string            `json:"action,omitempty"`
	Description          string            `json:"description,omitempty"`
	Headers              interface{}       `json:"headers,omitempty"`
	IpAddress            string            `json:"ipAddress,omitempty"`
	Name                 string            `json:"name,omitempty"`
	Priority             int               `json:"priority,omitempty"`
	SubnetMask           string            `json:"subnetMask,omitempty"`
	SubnetTrafficTag     int               `json:"subnetTrafficTag,omitempty"`
	Tag                  enums.IpFilterTag `json:"tag,omitempty"`
	VnetSubnetResourceId string            `json:"vnetSubnetResourceId,omitempty"`
	VnetTrafficTag       int               `json:"vnetTrafficTag,omitempty"`
}

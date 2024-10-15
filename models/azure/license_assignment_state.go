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

import (
	"github.com/certmichelin/azurehound/v2/enums"
)

// Provides details about license assignments to a user.
type LicenseAssignmentState struct {
	// The id of the group that assigns this license. If direct-assigned this field will be null.
	//
	// Read-only
	AssignedByGroup string `json:"assignedByGroup,omitempty"`

	// The service plans that are disabled in this assignment.
	//
	// Read-only
	DisabledPlans string `json:"disabledPlans,omitempty"`

	// License assignment failure error.
	Error enums.LicenseError `json:"error,omitempty"`

	// The unique identifier for the SKU
	//
	// Read-only
	SkuId string `json:"skuId,omitempty"`

	// Indicates the current state of this assignment.
	//
	// Read-only
	State enums.LicenseState `json:"state,omitempty"`
}

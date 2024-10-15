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

// Maintenance operations status.
type MaintenanceRedeployStatus struct {
	// True if customer is allowed to perform maintenance.
	IsCustomerInitiatedMaintenanceAllowed bool `json:"isCustomerInitiatedMaintenanceAllowed,omitempty"`

	// Message returned for the last maintenance operation.
	LastOperationMessage string `json:"lastOperationMessage,omitempty"`

	// The last maintenance operation result code.
	LastOperationResultCode enums.MaintenanceOperationCode `json:"lastOperationResultCode,omitempty"`

	// End time for the maintenance window.
	MaintenanceWindowEndTime string `json:"maintenanceWindowEndTime,omitempty"`

	// Start time for the maintenance window.
	MaintenanceWindowStartTime string `json:"maintenanceWindowStartTime,omitempty"`

	// End time for the pre maintenance window.
	PreMaintenanceWindowEndTime string `json:"preMaintenanceWindowEndTime,omitempty"`

	// Start time for the pre maintenance window.
	PreMaintenanceWindowStartTime string `json:"preMaintenanceWindowStartTime,omitempty"`
}

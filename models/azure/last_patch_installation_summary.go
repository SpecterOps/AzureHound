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

// Describes the properties of the last installed patch summary.
type LastPatchInstallationSummary struct {
	// The errors that were encountered during execution of the operation. The details array contains the list of them.
	Error ODataError `json:"error,omitempty"`

	// The number of all available patches but excluded explicitly by a customer-specified exclusion list match.
	ExcludedPatchCount int `json:"excludedPatchCount,omitempty"`

	// The count of patches that failed installation.
	FailedPatchCount int `json:"failedPatchCount,omitempty"`

	// The activity ID of the operation that produced this result. It is used to correlate across CRP and extension logs.
	InstallationActivityId string `json:"installationActivityId,omitempty"`

	// The count of patches that successfully installed.
	InstalledPatchCount int `json:"installedPatchCount,omitempty"`

	// The UTC timestamp when the operation began.
	LastModifiedTime string `json:"lastModifiedTime,omitempty"`

	// Describes whether the operation ran out of time before it completed all its intended actions.
	MaintenanceWindowExceeded bool `json:"maintenanceWindowExceeded,omitempty"`

	// The number of all available patches but not going to be installed because it didn't match a classification or
	// inclusion list entry.
	NotSelectedPatchCount int `json:"notSelectedPatchCount,omitempty"`

	// The number of all available patches expected to be installed over the course of the patch installation operation.
	PendingPatchCount int `json:"pendingPatchCount,omitempty"`

	// The UTC timestamp when the operation began.
	StartTime string `json:"startTime,omitempty"`

	// The overall success or failure status of the operation. It remains "InProgress" until the operation completes.
	// At that point it will become "Unknown", "Failed", "Succeeded", or "CompletedWithWarnings."
	Status enums.PatchStatus `json:"status,omitempty"`
}

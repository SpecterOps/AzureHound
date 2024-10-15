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

// Describes the properties of a virtual machine instance view for available patch summary.
type AvailablePatchSummary struct {
	// The activity ID of the operation that produced this result. It is used to correlate across CRP and extension logs.
	AssessmentActivityId string `json:"assessmentActivityId,omitempty"`

	// The number of critical or security patches that have been detected as available and not yet installed.
	CriticalAndSecurityPatchCount int `json:"criticalAndSecurityPatchCount,omitempty"`

	// The errors that were encountered during execution of the operation. The details array contains the list of them.
	Error ODataError `json:"error,omitempty"`

	// The UTC timestamp when the operation began.
	LastModifiedTime string `json:"lastModifiedTime,omitempty"`

	// The number of all available patches excluding critical and security.
	OtherPatchCount int `json:"otherPatchCount,omitempty"`

	// The overall reboot status of the VM. It will be true when partially installed patches require a reboot to
	// complete installation but the reboot has not yet occurred.
	RebootPending bool `json:"rebootPending,omitempty"`

	// The UTC timestamp when the operation began.
	StartTime string `json:"startTime,omitempty"`

	// The overall success or failure status of the operation. It remains "InProgress" until the operation completes.
	// At that point it will become "Unknown", "Failed", "Succeeded", or "CompletedWithWarnings."
	Status enums.PatchStatus `json:"status,omitempty"`
}

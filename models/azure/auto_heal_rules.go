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

type AutoHealRules struct {
	Actions  AutoHealActions  `json:"actions,omitempty"`
	Triggers AutoHealTriggers `json:"triggers,omitempty"`
}

type AutoHealActions struct {
	ActionType              enums.AutoHealActionType `json:"actionType,omitempty"`
	CustomAction            AutoHealCustomAction     `json:"customAction,omitempty"`
	MinProcessExecutionTime string                   `json:"minProcessExecutionTime,omitempty"`
}

type AutoHealCustomAction struct {
	Exe        string `json:"exe,omitempty"`
	Parameters string `json:"parameters,omitempty"`
}

type AutoHealTriggers struct {
	PrivateBytesInKB     int                            `json:"privateBytesInKB,omitempty"`
	Requests             RequestsBasedTrigger           `json:"requests,omitempty"`
	SlowRequests         SlowRequestsBasedTrigger       `json:"slowRequests,omitempty"`
	SlowRequestsWithPath []SlowRequestsBasedTrigger     `json:"slowRequestsWithPath,omitempty"`
	StatusCodes          []StatusCodesBasedTrigger      `json:"statusCodes,omitempty"`
	StatusCodesRange     []StatusCodesRangeBasedTrigger `json:"statusCodesRange,omitempty"`
}

type RequestsBasedTrigger struct {
	Count        int    `json:"count,omitempty"`
	TimeInterval string `json:"timeinterval,omitempty"`
}

type SlowRequestsBasedTrigger struct {
	Count        int    `json:"count,omitempty"`
	Path         string `json:"path,omitempty"`
	TimeInterval string `json:"timeInterval,omitempty"`
	TimeTaken    string `json:"timeTaken,omitempty"`
}

type StatusCodesBasedTrigger struct {
	Count        int    `json:"count,omitempty"`
	Path         string `json:"path,omitempty"`
	Status       int    `json:"status,omitempty"`
	SubStatus    int    `json:"subStatus,omitempty"`
	TimeInterval string `json:"timeInterval,omitempty"`
	Win32Status  int    `json:"win32Status,omitempty"`
}

type StatusCodesRangeBasedTrigger struct {
	Count        int    `json:"count,omitempty"`
	Path         string `json:"path,omitempty"`
	StatusCodes  string `json:"statusCodes,omitempty"`
	TimeInterval string `json:"timeInterval,omitempty"`
}

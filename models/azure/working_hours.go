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

type WorkingHours struct {
	// The days of the week on which the user works.
	DaysOfWeek []enums.DayOfWeek `json:"daysOfWeek,omitempty"`

	// The time of the day that the user starts working.
	StartTime string `json:"startTime,omitempty"`

	// The time of the day that the user stops working.
	EndTime string `json:"endTime,omitempty"`

	// The time zone to which the working hours apply.
	TimeZone TimeZoneBase `json:"timeZoneBase,omitempty"`
}

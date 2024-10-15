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

type MailboxSettings struct {
	// Folder ID of an archive folder for the user.
	ArchiveFolder string `json:"archiveFolder,omitempty"`

	// Configuration settings to automatically notify the sender of an incoming email with a message from the signed-in
	// user.
	AutomaticRepliesSetting AutomaticRepliesSetting `json:"automaticRepliesSetting,omitempty"`

	// The date format for the user's mailbox.
	DateFormat string `json:"dateFormat,omitempty"`

	// If the user has a calendar delegate, this specifies whether the delegate, mailbox owner, or both receive meeting
	// messages and meeting responses.
	DelegateMeetingMessageDeliveryOptions enums.MessageDeliveryOptions `json:"delegateMeetingMessageDeliveryOptions,omitempty"`

	// The locale information for the user, including the preferred language and country/region.
	Language LocaleInfo `json:"language,omitempty"`

	// The time format for the user's mailbox.
	TimeFormat string `json:"timeFormat,omitempty"`

	// The default time zone for the user's mailbox.
	TimeZone string `json:"timeZone,omitempty"`

	// The days of the week and hours in a specific time zone that the user works.
	WorkingHours WorkingHours `json:"workingHours,omitempty"`
}

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

// Specifies parental control settings for an application. These settings control the consent experience.
// For more detail see https://docs.microsoft.com/en-us/graph/api/resources/parentalcontrolsettings?view=graph-rest-1.0
type ParentalControlSettings struct {
	// Specifies the ISO 3166 country codes for which access to the application will be blocked for minors.
	CountriesBlockedForMinors []string `json:"countriesBlockedForMinors,omitempty"`
	// Specifies the legal age group rule that applies to users of the app.
	LegalAgeGroupRule enums.LegalAgeGroupRule `json:"legalAgeGroupRule,omitempty"`
}

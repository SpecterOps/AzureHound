// Copyright (C) 2025 Specter Ops, Inc.
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

type UnifiedRoleEligibilityScheduleInstance struct {
	Entity

	PrincipalId string `json:"principalId,omitempty"`

	RoleDefinitionId string `json:"roleDefinitionId,omitempty"`

	DirectoryScopeId string `json:"directoryScopeId,omitempty"`

	AppScopeId string `json:"appScopeId,omitempty"`

	StartDateTime string `json:"startDateTime,omitempty"`

	EndDateTime string `json:"endDateTime,omitempty"`

	MemberType string `json:"memberType,omitempty"`

	RoleEligibilityScheduleId string `json:"roleEligibilityScheduleId,omitempty"`

	AppScope AppScope `json:"appScope,omitempty"`

	RoleDefinition UnifiedRoleDefinition `json:"roleDefinition,omitempty"`

	Principal User `json:"principal,omitempty"`
}

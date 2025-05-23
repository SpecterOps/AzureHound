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

import "encoding/json"

type RoleManagementPolicyAssignment struct {
	Entity
	PolicyId         string `json:"policyId,omitempty"`
	ScopeId          string `json:"scopeId,omitempty"`
	RoleDefinitionId string `json:"roleDefinitionId,omitempty"`

	EndUserAssignmentRequiresApproval                 bool     `json:"endUserAssignmentRequiresApproval,omitempty"`
	EndUserAssignmentRequiresCAPAuthenticationContext bool     `json:"endUserAssignmentRequiresCAPAuthenticationContext,omitempty"`
	EndUserAssignmentUserApprovers                    []string `json:"endUserAssignmentUserApprovers,omitempty"`
	EndUserAssignmentGroupApprovers                   []string `json:"endUserAssignmentGroupApprovers,omitempty"`
	EndUserAssignmentRequiresMFA                      bool     `json:"endUserAssignmentRequiresMFA,omitempty"`
	EndUserAssignmentRequiresJustification            bool     `json:"endUserAssignmentRequiresJustification,omitempty"`
	EndUserAssignmentRequiresTicketInformation        bool     `json:"endUserAssignmentRequiresTicketInformation,omitempty"`
	TenantId                                          string   `json:"tenantId,omitempty"`
}

type RoleManagementPolicy struct {
	Id                    string          `json:"id,omitempty"`
	DisplayName           string          `json:"displayName,omitempty"`
	Description           string          `json:"description,omitempty"`
	IsOrganizationDefault bool            `json:"isOrganizationDefault,omitempty"`
	ScopeId               string          `json:"scopeId,omitempty"`
	ScopeType             string          `json:"scopeType,omitempty"`
	LastModifiedDateTime  string          `json:"lastModifiedDateTime,omitempty"`
	LastModifiedBy        string          `json:"lastModifiedBy,omitempty"`
	Rules                 json.RawMessage `json:"rules,omitempty"`
}

type UnifiedRoleManagementPolicyExpirationRule struct {
	Id                   string                         `json:"id,omitempty"`
	IsExpirationRequired bool                           `json:"isExpirationRequired,omitempty"`
	MaximumDuration      string                         `json:"maximumDuration,omitempty"`
	Target               RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

type UnifiedRoleManagementPolicyEnablementRule struct {
	Id           string                         `json:"id,omitempty"`
	EnabledRules []string                       `json:"enabledRules,omitempty"`
	Target       RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

type UnifiedRoleManagementPolicyNotificationRule struct {
	Id                         string                         `json:"id,omitempty"`
	NotificationType           string                         `json:"notificationType,omitempty"`
	RecipientType              string                         `json:"recipientType,omitempty"`
	NotificationLevel          string                         `json:"notificationLevel,omitempty"`
	IsDefaultRecipientsEnabled bool                           `json:"isDefaultRecipientsEnabled,omitempty"`
	NotificationRecipients     []string                       `json:"notificationRecipients,omitempty"`
	Target                     RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

type UnifiedRoleManagementPolicyAuthenticationContextRule struct {
	Id         string                         `json:"id,omitempty"`
	IsEnabled  bool                           `json:"isEnabled,omitempty"`
	ClaimValue string                         `json:"claimValue,omitempty"`
	Target     RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

type RoleManagementPolicyRuleTarget struct {
	Caller              string   `json:"caller,omitempty"`
	Operations          []string `json:"operations,omitempty"`
	Level               string   `json:"level,omitempty"`
	InheritableSettings []string `json:"inheritableSettings,omitempty"`
	EnforcedSettings    []string `json:"enforcedSettings,omitempty"`
}

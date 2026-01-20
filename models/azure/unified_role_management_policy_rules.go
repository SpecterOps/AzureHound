package azure

// UnifiedRoleManagementPolicyApprovalRule represents the unifiedRoleManagementPolicyApprovalRule resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyapprovalrule?view=graph-rest-1.0
type UnifiedRoleManagementPolicyApprovalRule struct {
	Entity

	IsExpirationRequired bool                           `json:"isExpirationRequired,omitempty"`
	MaximumDuration      string                         `json:"maximumDuration,omitempty"`
	Target               RoleManagementPolicyRuleTarget `json:"target,omitempty"`
	Setting              ApprovalSettings               `json:"setting,omitempty"`
}

// UnifiedRoleManagementPolicyExpirationRule represents the unifiedRoleManagementPolicyExpirationRule resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyexpirationrule?view=graph-rest-1.0
type UnifiedRoleManagementPolicyExpirationRule struct {
	Entity

	IsExpirationRequired bool                           `json:"isExpirationRequired,omitempty"`
	MaximumDuration      string                         `json:"maximumDuration,omitempty"`
	Target               RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

// UnifiedRoleManagementPolicyEnablementRule represents the unifiedRoleManagementPolicyEnablementRule resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyenablementrule?view=graph-rest-1.0
type UnifiedRoleManagementPolicyEnablementRule struct {
	Entity

	EnabledRules []string                       `json:"enabledRules,omitempty"`
	Target       RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

// UnifiedRoleManagementPolicyNotificationRule represents the unifiedRoleManagementPolicyNotificationRule resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicynotificationrule?view=graph-rest-1.0
type UnifiedRoleManagementPolicyNotificationRule struct {
	Entity

	NotificationType           string                         `json:"notificationType,omitempty"`
	RecipientType              string                         `json:"recipientType,omitempty"`
	NotificationLevel          string                         `json:"notificationLevel,omitempty"`
	IsDefaultRecipientsEnabled bool                           `json:"isDefaultRecipientsEnabled,omitempty"`
	NotificationRecipients     []string                       `json:"notificationRecipients,omitempty"`
	Target                     RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

// UnifiedRoleManagementPolicyAuthenticationContextRule represents the unifiedRoleManagementPolicyAuthenticationContextRule resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyauthenticationcontextrule?view=graph-rest-1.0
type UnifiedRoleManagementPolicyAuthenticationContextRule struct {
	Entity

	IsEnabled  bool                           `json:"isEnabled,omitempty"`
	ClaimValue string                         `json:"claimValue,omitempty"`
	Target     RoleManagementPolicyRuleTarget `json:"target,omitempty"`
}

// RoleManagementPolicyRuleTarget represents the unifiedRoleManagementPolicyRuleTarget resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyruletarget?view=graph-rest-1.0
type RoleManagementPolicyRuleTarget struct {
	Caller              string   `json:"caller,omitempty"`
	Operations          []string `json:"operations,omitempty"`
	Level               string   `json:"level,omitempty"`
	InheritableSettings []string `json:"inheritableSettings,omitempty"`
	EnforcedSettings    []string `json:"enforcedSettings,omitempty"`
}

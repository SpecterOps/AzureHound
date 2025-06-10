package enums

type RoleManagementPolicyRuleType string

const (
	PolicyRuleApproval              = "#microsoft.graph.unifiedRoleManagementPolicyApprovalRule"
	PolicyRuleExpiration            = "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule"
	PolicyRuleEnablement            = "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule"
	PolicyRuleNotification          = "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule"
	PolicyRuleAuthenticationContext = "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule"
)

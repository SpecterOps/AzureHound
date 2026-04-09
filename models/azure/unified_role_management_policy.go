package azure

import "encoding/json"

// UnifiedRoleManagementPolicy represents the unifiedRoleManagementPolicy resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicy?view=graph-rest-1.0
type UnifiedRoleManagementPolicy struct {
	Entity

	DisplayName           string   `json:"displayName,omitempty"`
	Description           string   `json:"description,omitempty"`
	IsOrganizationDefault bool     `json:"isOrganizationDefault,omitempty"`
	ScopeId               string   `json:"scopeId,omitempty"`
	ScopeType             string   `json:"scopeType,omitempty"`
	LastModifiedDateTime  string   `json:"lastModifiedDateTime,omitempty"`
	LastModifiedBy        Identity `json:"lastModifiedBy,omitempty"`

	// Rules represents an abstract type that may be one of multiple resource types which will be determined at runtime
	// https://learn.microsoft.com/en-us/graph/api/resources/unifiedrolemanagementpolicyrule?view=graph-rest-1.0
	Rules []json.RawMessage `json:"rules,omitempty"`
}

package azure

import "encoding/json"

// Identity defines the model for the Azure Identity resource type
// https://learn.microsoft.com/en-us/graph/api/resources/identity?view=graph-rest-1.0
type Identity struct {
	Entity

	DisplayName string          `json:"displayName,omitempty"`
	TenantId    string          `json:"tenantId,omitempty"`
	Thumbnails  json.RawMessage `json:"thumbnails,omitempty"`
}

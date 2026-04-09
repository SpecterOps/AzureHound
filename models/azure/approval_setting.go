package azure

// ApprovalSettings represents the approvalSettings resource type
// https://learn.microsoft.com/en-us/graph/api/resources/approvalsettings?view=graph-rest-1.0
type ApprovalSettings struct {
	Type                             string                  `json:"@odata.type,omitempty"`
	ApprovalMode                     string                  `json:"approvalMode,omitempty"`
	ApprovalStages                   []UnifiedApprovalStages `json:"approvalStages,omitempty"`
	IsApprovalRequired               bool                    `json:"isApprovalRequired,omitempty"`
	IsApprovalRequiredForExtension   bool                    `json:"isApprovalRequiredForExtension,omitempty"`
	IsRequestorJustificationRequired bool                    `json:"isRequestorJustificationRequired,omitempty"`
}

package azure

// UnifiedApprovalStages represents the unifiedApprovalStage resource type
// https://learn.microsoft.com/en-us/graph/api/resources/unifiedapprovalstage?view=graph-rest-1.0
type UnifiedApprovalStages struct {
	Type                            string                `json:"@odata.type,omitempty"`
	ApprovalStageTimeOutInDays      int32                 `json:"approvalStageTimeOutInDays,omitempty"`
	IsApproverJustificationRequired bool                  `json:"isApproverJustificationRequired,omitempty"`
	EscalationTimeInMinutes         int32                 `json:"escalationTimeInMinutes,omitempty"`
	PrimaryApprovers                []PrimaryApprovers    `json:"primaryApprovers,omitempty"`
	IsEscalationEnabled             bool                  `json:"isEscalationEnabled,omitempty"`
	EscalationApprovers             []EscalationApprovers `json:"escalationApprovers,omitempty"`
}

// PrimaryApprovers is a subjectSet collection
type PrimaryApprovers struct {
	Type    string `json:"@odata.type,omitempty"`
	UserId  string `json:"userId,omitempty"`
	GroupId string `json:"groupId,omitempty"`
}

// EscalationApprovers is a subjectSet collection
type EscalationApprovers struct {
	Type string `json:"@odata.type,omitempty"`
}

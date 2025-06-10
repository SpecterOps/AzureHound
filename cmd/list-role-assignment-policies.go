package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"slices"
	"time"

	"github.com/bloodhoundad/azurehound/v2/client"
	"github.com/bloodhoundad/azurehound/v2/client/query"
	"github.com/bloodhoundad/azurehound/v2/enums"
	"github.com/bloodhoundad/azurehound/v2/models"
	"github.com/bloodhoundad/azurehound/v2/models/azure"
	"github.com/bloodhoundad/azurehound/v2/panicrecovery"
	"github.com/bloodhoundad/azurehound/v2/pipeline"
	"github.com/spf13/cobra"
)

func init() {
	listRootCmd.AddCommand(listRoleAssignmentPoliciesCmd)
}

var listRoleAssignmentPoliciesCmd = &cobra.Command{
	Use:          "unified-role-assignment-policies",
	Short:        "Lists Unified Role Assignment Policies",
	Run:          listUnifiedRoleAssignmentPoliciesCmdImpl,
	SilenceUsage: true,
}

func listUnifiedRoleAssignmentPoliciesCmdImpl(cmd *cobra.Command, args []string) {
	ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, os.Kill)
	defer gracefulShutdown(stop)

	var (
		azClient = connectAndCreateClient()
		start    = time.Now()
		stream   = listRoleAssignmentPolicies(ctx, azClient)
	)

	panicrecovery.HandleBubbledPanic(ctx, stop, log)
	outputStream(ctx, stream)
	duration := time.Since(start)
	log.Info("collection completed", "duration", duration.String())
}

func listRoleAssignmentPolicies(ctx context.Context, azClient client.AzureClient) <-chan any {
	var (
		out = make(chan any)
	)

	go func() {
		defer panicrecovery.PanicRecovery()
		defer close(out)

		count := 0
		log.Info("collecting azure unified role assignment policies...")
		for item := range azClient.ListRoleAssignmentPolicies(ctx, query.GraphParams{
			Filter: "scopeId eq '/' and scopeType eq 'Directory'",
			Expand: "policy($expand=rules)",
		}) {
			if item.Error != nil {
				log.Error(item.Error, item.Error.Error())
				return
			} else {
				formattedItem, err := formatRoleManagementPolicyAssignment(item.Ok)
				if err != nil {
					log.Error(err, err.Error())
					continue
				}

				formattedItem.TenantId = azClient.TenantInfo().TenantId

				log.V(2).Info("found unified role assignment policy", "unifiedRoleAssignmentPolicy", formattedItem)
				count++

				if ok := pipeline.SendAny(ctx.Done(), out, azureWrapper[models.RoleManagementPolicyAssignment]{
					Data: formattedItem,
					Kind: enums.KindAZRoleManagementPolicyAssignment,
				}); !ok {
					return
				}
			}
		}

		log.V(1).Info("finished listing unified role assignment policies", "count", count)
	}()

	return out
}

type tempRuleType struct {
	Type enums.RoleManagementPolicyRuleType `json:"@odata.type"`
}

// formatRoleManagementPolicyAssignment takes a reference to a UnifiedRoleManagementPolicyAssignment and unmarshalls the model's Policy.Rules into their respective types
func formatRoleManagementPolicyAssignment(assignment azure.UnifiedRoleManagementPolicyAssignment) (models.RoleManagementPolicyAssignment, error) {
	rmPolicyAssignment := models.RoleManagementPolicyAssignment{
		UnifiedRoleManagementPolicyAssignment: assignment,

		Id:               assignment.Id,
		RoleDefinitionId: assignment.RoleDefinitionId,
	}

	rules := assignment.Policy.Rules
	for _, rule := range rules {
		var ruleType tempRuleType
		if err := json.Unmarshal(rule, &ruleType); err != nil {
			return rmPolicyAssignment, err
		}

		switch ruleType.Type {
		case enums.PolicyRuleApproval:
			if err := unmarshallPolicyRuleApproval(rule, &rmPolicyAssignment); err != nil {
				return rmPolicyAssignment, err
			}
		case enums.PolicyRuleEnablement:
			if err := unmarshallPolicyRuleEnablement(rule, &rmPolicyAssignment); err != nil {
				return rmPolicyAssignment, err
			}
		case enums.PolicyRuleAuthenticationContext:
			if err := unmarshallPolicyRuleAuthenticationContext(rule, &rmPolicyAssignment); err != nil {
				return rmPolicyAssignment, err
			}
		default:
			continue
		}
	}

	return rmPolicyAssignment, nil
}

// unmarshallPolicyRuleApproval unmarshalls the provided data into a UnifiedRoleManagementPolicyApprovalRule, extracts the relevant fields, and applies them to the rmPolicyAssignment
// Note: The provided rmPolicyAssignment will be modified when using this function
func unmarshallPolicyRuleApproval(data json.RawMessage, rmPolicyAssignment *models.RoleManagementPolicyAssignment) error {
	var rule azure.UnifiedRoleManagementPolicyApprovalRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return fmt.Errorf("error unmarshalling PolicyRuleApproval: %w", err)
	}

	var (
		userApprovers  []string
		groupApprovers []string
	)

	for _, approvalStage := range rule.Setting.ApprovalStages {
		for _, approver := range approvalStage.PrimaryApprovers {
			switch approver.Type {
			case enums.ApprovalStageSingleUser:
				userApprovers = append(userApprovers, approver.UserId)
			case enums.ApprovalStageGroupMembers:
				groupApprovers = append(groupApprovers, approver.GroupId)
			}
		}
	}

	rmPolicyAssignment.EndUserAssignmentUserApprovers = userApprovers
	rmPolicyAssignment.EndUserAssignmentGroupApprovers = groupApprovers
	rmPolicyAssignment.EndUserAssignmentRequiresApproval = rule.Setting.IsApprovalRequired

	return nil
}

// unmarshallPolicyRuleEnablement unmarshalls the provided data into a UnifiedRoleManagementPolicyEnablementRule, extracts the relevant fields, and applies them to the rmPolicyAssignment
// Note: The provided rmPolicyAssignment will be modified when using this function
func unmarshallPolicyRuleEnablement(data json.RawMessage, rmPolicyAssignment *models.RoleManagementPolicyAssignment) error {
	var rule azure.UnifiedRoleManagementPolicyEnablementRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return fmt.Errorf("error unmarshalling PolicyRuleEnablement: %w", err)
	}

	rmPolicyAssignment.EndUserAssignmentRequiresMFA = slices.Contains(rule.EnabledRules, "MultiFactorAuthentication")
	rmPolicyAssignment.EndUserAssignmentRequiresJustification = slices.Contains(rule.EnabledRules, "Justification")
	rmPolicyAssignment.EndUserAssignmentRequiresTicketInformation = slices.Contains(rule.EnabledRules, "Ticketing")

	return nil
}

// unmarshallPolicyRuleAuthenticationContext unmarshalls the provided data into a UnifiedRoleManagementPolicyAuthenticationContextRule, extracts the relevant fields, and applies them to the rmPolicyAssignment
// Note: The provided rmPolicyAssignment will be modified when using this function
func unmarshallPolicyRuleAuthenticationContext(data json.RawMessage, rmPolicyAssignment *models.RoleManagementPolicyAssignment) error {
	var rule azure.UnifiedRoleManagementPolicyAuthenticationContextRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return fmt.Errorf("error unmarshalling PolicyRuleAuthenticationContext: %w", err)
	}

	rmPolicyAssignment.EndUserAssignmentRequiresCAPAuthenticationContext = rule.IsEnabled

	return nil
}

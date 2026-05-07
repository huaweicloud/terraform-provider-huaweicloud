package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	approvalBatchActionResolve = "resolve"
	approvalBatchActionReject  = "reject"
	approvalBatchActionRecall  = "recall"
)

var architectureApprovalsBatchActionNonUpdatableParams = []string{
	"workspace_id",
	"approval_ids",
	"message",
	"action",
}

// @API DataArtsStudio PUT /v2/{project_id}/design/approvals
// @API DataArtsStudio PUT /v2/{project_id}/design/approvals/action
func ResourceArchitectureApprovalsBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureApprovalsBatchActionCreate,
		ReadContext:   resourceArchitectureApprovalsBatchActionRead,
		UpdateContext: resourceArchitectureApprovalsBatchActionUpdate,
		DeleteContext: resourceArchitectureApprovalsBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(architectureApprovalsBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the approval is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the approval belongs.`,
			},
			"approval_ids": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The IDs of the approvals.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action of the approval status to be approved.`,
			},

			// Optional parameters.
			"message": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The approval message is required for the approval action.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildArchitectureProcessApprovalBatchActionBodyParams(approvalIds, message string) map[string]interface{} {
	return map[string]interface{}{
		"ids": strings.Split(approvalIds, ","),
		"msg": message,
	}
}

func actionArchitectureProcessApprovals(client *golangsdk.ServiceClient, workspaceId, approvalIds, action, message string) error {
	httpUrl := "v2/{project_id}/design/approvals/action"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = fmt.Sprintf(path+"?action-id=%s", action)

	approvalsActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildArchitectureProcessApprovalBatchActionBodyParams(approvalIds, message)),
	}

	_, err := client.Request("PUT", path, &approvalsActionOpt)

	return err
}

func actionArchitectureWithdrawApprovals(client *golangsdk.ServiceClient, workspaceId, approvalIds string) error {
	httpUrl := "v2/{project_id}/design/approvals"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = fmt.Sprintf(path+"?ids=%s", approvalIds)

	approvalsActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(workspaceId),
		OkCodes:          []int{200, 204},
	}

	_, err := client.Request("PUT", path, &approvalsActionOpt)

	return err
}

func resourceArchitectureApprovalsBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		approvalIds = d.Get("approval_ids").(string)
		action      = d.Get("action").(string)
		message     = d.Get("message").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	switch action {
	case approvalBatchActionResolve, approvalBatchActionReject:
		err = actionArchitectureProcessApprovals(client, workspaceId, approvalIds, action, message)
		if err != nil {
			return diag.Errorf("error action processing architecture approvals (%s): %s", approvalIds, err)
		}

	case approvalBatchActionRecall:
		err = actionArchitectureWithdrawApprovals(client, workspaceId, approvalIds)
		if err != nil {
			return diag.Errorf("error action withdrawing architecture approvals (%s): %s", approvalIds, err)
		}

	default:
		return diag.Errorf("invalid action type: %s, expected 'resolve', 'reject', or 'recall'", action)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceArchitectureApprovalsBatchActionRead(ctx, d, meta)
}

func resourceArchitectureApprovalsBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for reject/resolve/recall the API.
	// There is no API for the provider to query the history of this API action.
	return nil
}

func resourceArchitectureApprovalsBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceArchitectureApprovalsBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	msg := `This resource is only a one-time action resource for operating approvals. Deleting this resource will not
clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  msg,
		},
	}
}

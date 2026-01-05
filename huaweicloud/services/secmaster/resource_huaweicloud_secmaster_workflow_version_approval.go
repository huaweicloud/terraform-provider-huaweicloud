package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var workflowVersionApprovalNonUpdatableParams = []string{"workspace_id", "version_id", "content", "result"}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}/approval
func ResourceWorkflowVersionApproval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowVersionApprovalCreate,
		UpdateContext: resourceWorkflowVersionApprovalUpdate,
		ReadContext:   resourceWorkflowVersionApprovalRead,
		DeleteContext: resourceWorkflowVersionApprovalDelete,

		CustomizeDiff: config.FlexibleForceNew(workflowVersionApprovalNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceWorkflowVersionApprovalCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/versions/{version_id}/approval"
		workspaceId   = d.Get("workspace_id").(string)
		versionId     = d.Get("version_id").(string)
	)

	createClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)
	createPath = strings.ReplaceAll(createPath, "{version_id}", versionId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: buildWorkflowVersionApprovalBodyParams(d),
	}

	_, err = createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error approvaling workflow version: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return nil
}

func buildWorkflowVersionApprovalBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"content": d.Get("content"),
		"result":  d.Get("result"),
	}

	return bodyParams
}

func resourceWorkflowVersionApprovalRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowVersionApprovalUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowVersionApprovalDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource to approval workflow version. Deleting this resource will
		not change the status of the currently resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

package secmaster

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var workflowVersionValidationNonUpdatableParams = []string{
	"workspace_id",
	"aopworkflow_id",
	"mode",
	"taskconfig",
	"taskflow",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/workflows/validation
func ResourceWorkflowVersionValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWorkflowVersionValidationCreate,
		UpdateContext: resourceWorkflowVersionValidationUpdate,
		ReadContext:   resourceWorkflowVersionValidationRead,
		DeleteContext: resourceWorkflowVersionValidationDelete,

		CustomizeDiff: config.FlexibleForceNew(workflowVersionValidationNonUpdatableParams),

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
			"aopworkflow_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"taskconfig": {
				Type:     schema.TypeString,
				Required: true,
			},
			"taskflow": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"result": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWorkflowVersionValidationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/workflows/validation"
		workspaceId   = d.Get("workspace_id").(string)
	)

	createClient, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := createClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
		JSONBody: buildWorkflowVersionValidationBodyParams(d),
	}

	createResp, err := createClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error validating workflow version: %s", err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	mErr := multierror.Append(nil,
		d.Set("result", utils.PathSearch("data.result", respBody, nil)),
		d.Set("detail", utils.PathSearch("data.detail", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildWorkflowVersionValidationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aopworkflow_id": d.Get("aopworkflow_id"),
		"mode":           d.Get("mode"),
		"taskconfig":     d.Get("taskconfig"),
		"taskflow":       d.Get("taskflow"),
	}

	return bodyParams
}

func resourceWorkflowVersionValidationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowVersionValidationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceWorkflowVersionValidationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource to validate workflow version. Deleting this resource will
		not change the status of the currently resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

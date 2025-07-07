package workspace

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var desktopPoolActionNonUpdateParams = []string{"pool_id", "op_type", "type"}

// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/action
// @API Workspace GET /v2/{project_id}/workspace-jobs/{job_id}
func ResourceDesktopPoolAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopPoolActionCreate,
		ReadContext:   resourceDesktopPoolActionRead,
		UpdateContext: resourceDesktopPoolActionUpdate,
		DeleteContext: resourceDesktopPoolActionDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopPoolActionNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktop pool is located.`,
			},
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop pool.`,
			},
			"op_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of desktop pool operation.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Whether the operation is mandatory.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of operation dispatch task.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDesktopPoolActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"op_type": d.Get("op_type"),
		"type":    d.Get("type"),
	}
}

func resourceDesktopPoolActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktop-pools/{pool_id}/action"
		poolID  = d.Get("pool_id").(string)
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace App client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{pool_id}", poolID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDesktopPoolActionBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error dispatching desktop pool operation: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable find job ID from API response")
	}
	// Backup job ID proves that the request was successful
	d.SetId(jobId)

	status, err := waitForJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) to completed: %s", jobId, err)
	}
	d.Set("status", status)

	return resourceDesktopPoolActionRead(ctx, d, meta)
}

func resourceDesktopPoolActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopPoolActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopPoolActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to operate desktop pool. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

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

var desktopPoolNotificationNonUpdateParams = []string{"pool_id", "notifications"}

// @API Workspace POST /v2/{project_id}/desktop-pools/{pool_id}/notifications
// @API Workspace GET /v2/{project_id}/workspace-jobs/{job_id}
func ResourceDesktopPoolNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopPoolNotificationCreate,
		ReadContext:   resourceDesktopPoolNotificationRead,
		UpdateContext: resourceDesktopPoolNotificationUpdate,
		DeleteContext: resourceDesktopPoolNotificationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopPoolNotificationNonUpdateParams),

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
			"notifications": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The message want to dispatch.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of notification dispatch task.`,
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

func buildDesktopPoolNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"notifications": d.Get("notifications"),
	}
}

func resourceDesktopPoolNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktop-pools/{pool_id}/notifications"
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
		JSONBody: buildDesktopPoolNotificationBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating desktop pool notification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to job ID from API response")
	}
	// Backup job ID proves that the request was successful
	d.SetId(jobId)

	status, err := waitForJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) to completed: %s", jobId, err)
	}
	d.Set("status", status)

	return resourceDesktopPoolNotificationRead(ctx, d, meta)
}

func resourceDesktopPoolNotificationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopPoolNotificationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopPoolNotificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to dispatch desktop pool message. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

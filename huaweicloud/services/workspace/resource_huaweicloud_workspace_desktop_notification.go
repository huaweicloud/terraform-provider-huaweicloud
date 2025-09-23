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

var desktopNotificationNonUpdateParams = []string{"desktop_ids", "notifications"}

// @API Workspace POST /v2/{project_id}/desktops/notifications
// @API Workspace GET /v2/{project_id}/workspace-jobs/{job_id}
func ResourceDesktopNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopNotificationCreate,
		ReadContext:   resourceDesktopNotificationRead,
		UpdateContext: resourceDesktopNotificationUpdate,
		DeleteContext: resourceDesktopNotificationDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(desktopNotificationNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktop is located.`,
			},
			"desktop_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of desktop IDs.`,
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

func buildDesktopNotificationBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"desktop_ids":   utils.ExpandToStringList(d.Get("desktop_ids").([]interface{})),
		"notifications": d.Get("notifications"),
	}
}

func resourceDesktopNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/desktops/notifications"
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildDesktopNotificationBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating desktop notification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to get job ID from API response")
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

func resourceDesktopNotificationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopNotificationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopNotificationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to dispatch desktop message. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

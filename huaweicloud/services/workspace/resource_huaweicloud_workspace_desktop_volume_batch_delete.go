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

var desktopVolumeBatchDeleteNonUpdatableParams = []string{"desktop_id", "volume_ids"}

// @API Workspace POST /v2/{project_id}/desktops/{desktop_id}/volumes/batch-delete
// @API Workspace GET /v2/{project_id}/workspace-jobs/{job_id}
func ResourceDesktopVolumeBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDesktopVolumeBatchDeleteCreate,
		ReadContext:   resourceDesktopVolumeBatchDeleteRead,
		UpdateContext: resourceDesktopVolumeBatchDeleteUpdate,
		DeleteContext: resourceDesktopVolumeBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(desktopVolumeBatchDeleteNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the desktop is located.`,
			},

			// Required parameters.
			"desktop_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop to which the data volumes to be deleted belongs.`,
			},
			"volume_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of desktop data volume IDs to be deleted.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDesktopVolumeBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"volume_ids": d.Get("volume_ids"),
	}
}

func resourceDesktopVolumeBatchDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/desktops/{desktop_id}/volumes/batch-delete"
		desktopId = d.Get("desktop_id").(string)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{desktop_id}", desktopId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildDesktopVolumeBatchDeleteBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error deleting desktop data volumes: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find job ID from API response")
	}
	d.SetId(jobId)

	status, err := waitForJobCompleted(ctx, client, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the job (%s) to complete: %s", jobId, err)
	}
	if status == "FAIL" {
		return diag.Errorf("the job (%s) of delete desktop data volume has failed.", jobId)
	}

	return resourceDesktopVolumeBatchDeleteRead(ctx, d, meta)
}

func resourceDesktopVolumeBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopVolumeBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDesktopVolumeBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch deleting desktop data volumes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

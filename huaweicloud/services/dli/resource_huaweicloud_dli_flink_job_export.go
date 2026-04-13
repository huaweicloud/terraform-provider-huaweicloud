package dli

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

var flinkJobExportNonUpdatableParams = []string{
	"obs_path",
	"job_ids",
}

// @API DLI POST /v1.0/{project_id}/streaming/jobs/export
func ResourceFlinkJobExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFlinkJobExportCreate,
		ReadContext:   resourceFlinkJobExportRead,
		UpdateContext: resourceFlinkJobExportUpdate,
		DeleteContext: resourceFlinkJobExportDelete,

		CustomizeDiff: config.FlexibleForceNew(flinkJobExportNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where to export flink jobs.`,
			},

			// Required parameters.
			"obs_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS save path for the exported job file.`,
			},
			"job_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The set of job IDs to be exported.`,
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
					}),
			},
		},
	}
}

func buildFlinkJobExportBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"obs_dir":      d.Get("obs_path"),
		"job_selected": d.Get("job_ids"),

		// This request parameter means 'whether to export flink job', but the name of the API is 'export flink job'.
		// So we set the default value to true.
		"is_selected": true,
	}

	return bodyParams
}

func resourceFlinkJobExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI client: %s", err)
	}

	httpUrl := "v1.0/{project_id}/streaming/jobs/export"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildFlinkJobExportBodyParams(d),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error exporting DLI Flink jobs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// is_success: the type of the API is string
	if utils.PathSearch("is_success", respBody, "false").(string) == "false" {
		return diag.Errorf("unable to export the jobs: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceFlinkJobExportRead(ctx, d, meta)
}

func resourceFlinkJobExportRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkJobExportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceFlinkJobExportDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for exporting Flink jobs. Deleting this resource will not
clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

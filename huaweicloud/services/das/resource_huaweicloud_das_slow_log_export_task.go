package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	slowLogExportTaskNonUpdatableParams = []string{
		"instance_id",
		"bucket_name",
		"start_time",
		"end_time",
		"file_path",
		"export_type",
		"sort_field",
		"sort_asc",
		"user_name",
		"client_ip_address",
		"killed",
		"execute_time_min",
		"execute_time_max",
		"min_avg_execute_time",
		"max_avg_execute_time",
		"rows_max_examined",
		"rows_min_examined",
		"fuzzy_sql",
		"operation",
		"time_zone",
	}
)

// @API DAS POST /v3/{project_id}/instances/{instance_id}/slow-log/create-slow-log-export-task
func ResourceSlowLogExportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSlowLogExportTaskCreate,
		ReadContext:   resourceSlowLogExportTaskRead,
		UpdateContext: resourceSlowLogExportTaskUpdate,
		DeleteContext: resourceSlowLogExportTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(slowLogExportTaskNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the slow log export task is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance ID.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS bucket name.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time, in RFC3339 format.`,
			},

			// Optional parameters.
			"file_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OBS file directory.`,
			},
			"export_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The export type.`,
			},
			"sort_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort field.`,
			},
			"sort_asc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to sort in ascending order.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The user name.`,
			},
			"client_ip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The client IP address.`,
			},
			"killed": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The execution status.`,
			},
			"execute_time_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The minimum execution time, in milliseconds.`,
			},
			"execute_time_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum execution time, in milliseconds.`,
			},
			"min_avg_execute_time": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The minimum average execution time.`,
			},
			"max_avg_execute_time": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: `The maximum average execution time.`,
			},
			"rows_max_examined": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum number of scanned rows.`,
			},
			"rows_min_examined": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The minimum number of scanned rows.`,
			},
			"fuzzy_sql": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The fuzzy SQL.`,
			},
			"operation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The operation type.`,
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time zone.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildSlowLogExportTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		// required
		"bucket_name": d.Get("bucket_name"),
		"start_time":  utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string)),
		"end_time":    utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string)),

		// optional
		"file_path":            utils.ValueIgnoreEmpty(d.Get("file_path")),
		"export_type":          utils.ValueIgnoreEmpty(d.Get("export_type")),
		"sort_field":           utils.ValueIgnoreEmpty(d.Get("sort_field")),
		"sort_asc":             utils.ValueIgnoreEmpty(d.Get("sort_asc")),
		"killed":               utils.ValueIgnoreEmpty(d.Get("killed")),
		"execute_time_min":     utils.ValueIgnoreEmpty(d.Get("execute_time_min")),
		"execute_time_max":     utils.ValueIgnoreEmpty(d.Get("execute_time_max")),
		"min_avg_execute_time": utils.ValueIgnoreEmpty(d.Get("min_avg_execute_time")),
		"max_avg_execute_time": utils.ValueIgnoreEmpty(d.Get("max_avg_execute_time")),
		"rows_max_examined":    utils.ValueIgnoreEmpty(d.Get("rows_max_examined")),
		"rows_min_examined":    utils.ValueIgnoreEmpty(d.Get("rows_min_examined")),
		"fuzzy_sql":            utils.ValueIgnoreEmpty(d.Get("fuzzy_sql")),
		"operation":            utils.ValueIgnoreEmpty(d.Get("operation")),
		"time_zone":            utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"user":                 utils.ValueIgnoreEmpty(d.Get("user_name")),
		"client":               utils.ValueIgnoreEmpty(d.Get("client_ip_address")),
	}

	return result
}

func resourceSlowLogExportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/slow-log/create-slow-log-export-task"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
		JSONBody: utils.RemoveNil(buildSlowLogExportTaskBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS slow log export task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	taskId := utils.PathSearch("task_id", respBody, nil)
	if taskId == nil {
		return diag.Errorf("unable to find the task ID from the API response")
	}
	d.SetId(fmt.Sprintf("%v", taskId))

	return resourceSlowLogExportTaskRead(ctx, d, meta)
}

func resourceSlowLogExportTaskRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSlowLogExportTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSlowLogExportTaskDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for exporting slow log tasks. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information
from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

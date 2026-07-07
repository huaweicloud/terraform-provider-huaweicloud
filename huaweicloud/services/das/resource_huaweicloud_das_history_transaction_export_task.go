package das

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	historyTransactionExportTaskStatusQuotaFull  float64 = -1
	historyTransactionExportTaskStatusWaiting    float64 = 0
	historyTransactionExportTaskStatusRunning    float64 = 1
	historyTransactionExportTaskStatusFailed     float64 = 2
	historyTransactionExportTaskStatusSuccess    float64 = 3
	historyTransactionExportTaskStatusTimeout    float64 = 4
	historyTransactionExportTaskStatusObsDeleted float64 = 5
)

var (
	historyTransactionExportTaskNonUpdatableParams = []string{
		"instance_id",
		"bucket_name",
		"start_time",
		"end_time",
		"file_path",
		"time_zone",
		"order_field",
		"order_by",
		"last_sec_min",
		"last_sec_max",
	}

	historyTransactionExportTaskNotFoundCodes = []string{
		"DAS.2002", // The exported task does not exist during deletion.
	}
)

// @API DAS POST /v3/{project_id}/transaction/{instance_id}/create-export-task
// @API DAS GET /v3/{project_id}/transaction/{instance_id}/get-export-task-info
// @API DAS POST /v3/{project_id}/transaction/{instance_id}/delete-export-task
func ResourceHistoryTransactionExportTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHistoryTransactionExportTaskCreate,
		ReadContext:   resourceHistoryTransactionExportTaskRead,
		UpdateContext: resourceHistoryTransactionExportTaskUpdate,
		DeleteContext: resourceHistoryTransactionExportTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(historyTransactionExportTaskNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceHistoryTransactionExportTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the history transaction export task is located.`,
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
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time zone.`,
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort field.`,
			},
			"order_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sort order.`,
			},
			"last_sec_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The minimum duration.`,
			},
			"last_sec_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The maximum duration.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The task status.`,
			},
			"export_line_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of exported lines.`,
			},
			"download_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The download URL of the exported file.`,
			},
			"created_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task creation time, in RFC3339 format.`,
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

func buildHistoryTransactionExportTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"bucket_name": d.Get("bucket_name"),
		"start_at":    utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string)),
		"end_at":      utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string)),
	}

	if v, ok := d.GetOk("file_path"); ok {
		result["file_path"] = v
	}
	if v, ok := d.GetOk("time_zone"); ok {
		result["time_zone"] = v
	}
	if v, ok := d.GetOk("order_field"); ok {
		result["order"] = v
	}
	if v, ok := d.GetOk("order_by"); ok {
		result["order_by"] = v
	}
	if v, ok := d.GetOk("last_sec_min"); ok {
		result["last_sec_min"] = v
	}
	if v, ok := d.GetOk("last_sec_max"); ok {
		result["last_sec_max"] = v
	}

	return result
}

func GetHistoryTransactionExportTask(client *golangsdk.ServiceClient, instanceId, taskId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/transaction/{instance_id}/get-export-task-info?task_id={task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func refreshHistoryTransactionExportTaskStatusFunc(client *golangsdk.ServiceClient, instanceId, taskId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetHistoryTransactionExportTask(client, instanceId, taskId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("task_status", respBody, float64(-1)).(float64)
		statusStr := fmt.Sprintf("%v", status)

		switch status {
		case historyTransactionExportTaskStatusQuotaFull, historyTransactionExportTaskStatusFailed,
			historyTransactionExportTaskStatusTimeout, historyTransactionExportTaskStatusObsDeleted:
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", statusStr)
		case historyTransactionExportTaskStatusSuccess:
			return respBody, "COMPLETED", nil
		default:
			return respBody, "PENDING", nil
		}
	}
}

func waitForHistoryTransactionExportTaskComplete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	instanceId := d.Get("instance_id").(string)
	taskId := d.Id()

	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshHistoryTransactionExportTaskStatusFunc(client, instanceId, taskId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the DAS history transaction export task (%s) to complete: %s", taskId, err)
	}
	return nil
}

func resourceHistoryTransactionExportTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/transaction/{instance_id}/create-export-task"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildHistoryTransactionExportTaskBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS history transaction export task: %s", err)
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

	if err = waitForHistoryTransactionExportTaskComplete(ctx, client, d); err != nil {
		return diag.Errorf("error waiting for the DAS history transaction export task to complete: %s", err)
	}

	return resourceHistoryTransactionExportTaskRead(ctx, d, meta)
}

func resourceHistoryTransactionExportTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		taskId     = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := GetHistoryTransactionExportTask(client, instanceId, taskId)
	if err != nil {
		return common.CheckDeletedDiag(
			d, err, fmt.Sprintf("error retrieving DAS history transaction export task (%s)", taskId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", respBody, nil)),
		d.Set("status", utils.PathSearch("task_status", respBody, nil)),
		d.Set("export_line_num", utils.PathSearch("export_line_num", respBody, nil)),
		d.Set("download_url", utils.PathSearch("download_url", respBody, nil)),
		d.Set("created_time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_at", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceHistoryTransactionExportTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceHistoryTransactionExportTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		taskId     = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/transaction/{instance_id}/delete-export-task"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"task_id": d.Id(),
		},
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errorCodeStr",
			historyTransactionExportTaskNotFoundCodes...),
			fmt.Sprintf("error deleting DAS history transaction export task (%s)", taskId))
	}

	return nil
}

func resourceHistoryTransactionExportTaskImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	).ErrorOrNil()
}

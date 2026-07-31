package dsc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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

var (
	databaseWatermarkEmbedTaskNonUpdatableParams = []string{
		"source_db_info.*.ins_id",
		"source_db_info.*.ins_name",
		"source_db_info.*.db_type",
		"target_db_info.*.ins_id",
		"target_db_info.*.ins_name",
		"target_db_info.*.db_type",
	}

	databaseWatermarkEmbedTaskNotFoundErrCodes = []string{
		"dsc.10000009", // The DSC instance does not exist.
	}
)

// @API DSC POST /v1/{project_id}/data-watermark-embed-tasks
// @API DSC GET /v1/{project_id}/data-watermark-embed-tasks
// @API DSC PUT /v1/{project_id}/data-watermark-embed-tasks/{id}
// @API DSC DELETE /v1/{project_id}/data-watermark-embed-tasks/{id}
func ResourceDatabaseWatermarkEmbedTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseWatermarkEmbedTaskCreate,
		ReadContext:   resourceDatabaseWatermarkEmbedTaskRead,
		UpdateContext: resourceDatabaseWatermarkEmbedTaskUpdate,
		DeleteContext: resourceDatabaseWatermarkEmbedTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(databaseWatermarkEmbedTaskNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the database watermark embed task is located.`,
			},

			// Required parameters.
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database watermark embed task.`,
			},
			"water_mark": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The watermark content to be embedded into the database.`,
			},
			"watermark_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The watermark algorithm version.`,
			},
			"db_water_param": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        databaseWatermarkEmbedTaskDbWaterParamSchema(),
				Description: `The database watermark embedding configuration.`,
			},
			"source_db_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        databaseWatermarkEmbedTaskDbInfoSchema(),
				Description: `The source database information from which data is read.`,
			},
			"target_db_info": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        databaseWatermarkEmbedTaskDbInfoSchema(),
				Description: `The target database information to which watermarked data is written.`,
			},
			"error_code": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The watermark error correction level.`,
			},

			// Optional parameters.
			"selected_fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        databaseWatermarkEmbedTaskColumnInfoSchema(),
				Description: `The selected field list used for watermark embedding.`,
			},
			"watermark_describe": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the watermark.`,
			},
			"schedule_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable task scheduling.`,
			},
			"schedule_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schedule type of the task.`,
			},
			"start_now": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to start the watermark embed task immediately after creation.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The scheduled start time of the task, in RFC3339 format.`,
			},

			// Attributes.
			"task_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The running state of the watermark embed task.`,
			},
			"task_create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the task, in RFC3339 format.`,
			},
			"task_end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time of the task, in RFC3339 format.`,
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

func databaseWatermarkEmbedTaskDbWaterParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"embed_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The watermark embed mode.`,
			},
			"watermark_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The watermark key used to embed and extract the watermark.`,
			},
			"params": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        databaseWatermarkEmbedTaskEmbedParamSchema(),
				Description: `The fake-column embed parameter list.`,
			},
			"row_spacing": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The row spacing used by fake-row watermark.`,
			},
		},
	}
}

func databaseWatermarkEmbedTaskEmbedParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"new_column_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the new fake column to be created.`,
			},
			"new_column_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data type of the new fake column.`,
			},
			"fake_strategy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The strategy used to generate fake data.`,
			},
			"fake_param": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem:        databaseWatermarkEmbedTaskFakeParamSchema(),
				Description: `The configuration of fake data generation.`,
			},
		},
	}
}

func databaseWatermarkEmbedTaskFakeParamSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"address_accuracy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The accuracy of generated address data.`,
			},
			"date_begin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start date of the generated date range, in RFC3339 format.`,
			},
			"date_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The end date of the generated date range, in RFC3339 format.`,
			},
			"random_accuracy": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The precision of generated random numbers.`,
			},
			"random_begin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The lower bound of the generated random value range.`,
			},
			"random_distribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The distribution mode of generated random values.`,
			},
			"random_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The upper bound of the generated random value range.`,
			},
			"string_distribute": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The distribution mode of generated string values.`,
			},
		},
	}
}

func databaseWatermarkEmbedTaskColumnInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"column_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the selected column.`,
			},
			"column_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The data type of the selected column.`,
			},
		},
	}
}

func databaseWatermarkEmbedTaskDbInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the authorized database.`,
			},
			"db_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database.`,
			},
			"db_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},
			"ins_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the database instance.`,
			},
			"ins_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database instance.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the database table.`,
			},
			"schema_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schema name of the database.`,
			},
		},
	}
}

func buildDatabaseWatermarkEmbedTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := utils.RemoveNil(map[string]interface{}{
		// Required parameters.
		"task_name":         d.Get("task_name"),
		"water_mark":        d.Get("water_mark"),
		"watermark_version": d.Get("watermark_version"),
		"db_water_param":    buildDatabaseWatermarkEmbedTaskDbWaterParamBodyParams(d.Get("db_water_param").([]interface{})),
		"source_db_info":    buildDatabaseWatermarkEmbedTaskDbInfoBodyParams(d.Get("source_db_info").([]interface{})),
		"target_db_info":    buildDatabaseWatermarkEmbedTaskDbInfoBodyParams(d.Get("target_db_info").([]interface{})),
		"error_code":        d.Get("error_code"),
		// Optional parameters.
		"start_now":          utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "start_now"),
		"watermark_describe": utils.ValueIgnoreEmpty(d.Get("watermark_describe")),
		"schedule_switch":    utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "schedule_switch"),
		"schedule_type":      utils.ValueIgnoreEmpty(d.Get("schedule_type")),
		"start_time":         utils.ValueIgnoreEmpty(utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))),
		"selected_fields":    buildDatabaseWatermarkEmbedTaskColumnInfoBodyParams(d.Get("selected_fields").([]interface{})),
	})

	// For API, the `selected_fields` is required.
	if params["selected_fields"] == nil {
		params["selected_fields"] = make([]interface{}, 0)
	}

	return params
}

func buildDatabaseWatermarkEmbedTaskDbWaterParamBodyParams(dbWaterParams []interface{}) map[string]interface{} {
	if len(dbWaterParams) == 0 {
		return nil
	}

	dbWaterParam := dbWaterParams[0]
	return map[string]interface{}{
		"embed_mode":    utils.PathSearch("embed_mode", dbWaterParam, nil),
		"watermark_key": utils.ValueIgnoreEmpty(utils.PathSearch("watermark_key", dbWaterParam, nil)),
		"params": buildDatabaseWatermarkEmbedTaskEmbedParamBodyParams(utils.PathSearch("params", dbWaterParam,
			make([]interface{}, 0)).([]interface{})),
		"row_spacing": utils.ValueIgnoreEmpty(utils.PathSearch("row_spacing", dbWaterParam, nil)),
	}
}

func buildDatabaseWatermarkEmbedTaskEmbedParamBodyParams(params []interface{}) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(params))
	for _, v := range params {
		rst = append(rst, map[string]interface{}{
			"new_column_name": utils.PathSearch("new_column_name", v, nil),
			"new_column_type": utils.PathSearch("new_column_type", v, nil),
			"fake_strategy":   utils.PathSearch("fake_strategy", v, nil),
			"fake_param": buildDatabaseWatermarkEmbedTaskFakeParamBodyParams(utils.PathSearch("fake_param", v,
				make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func buildDatabaseWatermarkEmbedTaskFakeParamBodyParams(fakeParams []interface{}) map[string]interface{} {
	if len(fakeParams) == 0 || fakeParams[0] == nil {
		return nil
	}

	fakeParam := fakeParams[0]
	return map[string]interface{}{
		"address_accuracy": utils.ValueIgnoreEmpty(utils.PathSearch("address_accuracy", fakeParam, nil)),
		"date_begin": utils.ValueIgnoreEmpty(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("date_begin",
			fakeParam, "").(string))),
		"date_end": utils.ValueIgnoreEmpty(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("date_end",
			fakeParam, "").(string))),
		"random_accuracy":   utils.ValueIgnoreEmpty(utils.PathSearch("random_accuracy", fakeParam, nil)),
		"random_begin":      utils.ValueIgnoreEmpty(utils.PathSearch("random_begin", fakeParam, nil)),
		"random_distribute": utils.ValueIgnoreEmpty(utils.PathSearch("random_distribute", fakeParam, nil)),
		"random_end":        utils.ValueIgnoreEmpty(utils.PathSearch("random_end", fakeParam, nil)),
		"string_distribute": utils.ValueIgnoreEmpty(utils.PathSearch("string_distribute", fakeParam, nil)),
	}
}

func buildDatabaseWatermarkEmbedTaskDbInfoBodyParams(dbInfos []interface{}) map[string]interface{} {
	if len(dbInfos) == 0 {
		return nil
	}

	dbInfo := dbInfos[0]
	return map[string]interface{}{
		"db_id":       utils.PathSearch("db_id", dbInfo, nil),
		"db_name":     utils.PathSearch("db_name", dbInfo, nil),
		"db_type":     utils.PathSearch("db_type", dbInfo, nil),
		"ins_id":      utils.PathSearch("ins_id", dbInfo, nil),
		"ins_name":    utils.PathSearch("ins_name", dbInfo, nil),
		"table_name":  utils.PathSearch("table_name", dbInfo, nil),
		"schema_name": utils.ValueIgnoreEmpty(utils.PathSearch("schema_name", dbInfo, nil)),
	}
}

func buildDatabaseWatermarkEmbedTaskColumnInfoBodyParams(selectedFields []interface{}) []map[string]interface{} {
	if len(selectedFields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(selectedFields))
	for _, v := range selectedFields {
		rst = append(rst, map[string]interface{}{
			"column_name": utils.ValueIgnoreEmpty(utils.PathSearch("column_name", v, nil)),
			"column_type": utils.ValueIgnoreEmpty(utils.PathSearch("column_type", v, nil)),
		})
	}

	return rst
}

func createDatabaseWatermarkEmbedTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/data-watermark-embed-tasks"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDatabaseWatermarkEmbedTaskBodyParams(d),
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceDatabaseWatermarkEmbedTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createDatabaseWatermarkEmbedTask(client, d)
	if err != nil {
		return diag.Errorf("error creating database watermark embed task: %s", err)
	}

	tasks, err := listDatabaseWatermarkEmbedTasks(client)
	if err != nil {
		return diag.Errorf("error retrieving database watermark embed tasks: %s", err)
	}

	taskId, ok := utils.PathSearch(fmt.Sprintf("[?task_name=='%s']|[0].id", d.Get("task_name").(string)), tasks, float64(0)).(float64)
	if !ok || taskId == 0 {
		return diag.Errorf("unable to find the database watermark embed task ID from API response")
	}

	taskIdStr := strconv.Itoa(int(taskId))
	d.SetId(taskIdStr)

	if isImmediateStart, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "start_now").(bool); ok && isImmediateStart {
		if _, err := waitForDatabaseWatermarkEmbedTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), taskIdStr); err != nil {
			return diag.Errorf("error waiting for database watermark embed task to be completed: %s", err)
		}
	}

	return resourceDatabaseWatermarkEmbedTaskRead(ctx, d, meta)
}

func waitForDatabaseWatermarkEmbedTaskCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration,
	taskId string) (interface{}, error) {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDatabaseWatermarkEmbedTaskStatusFunc(client, taskId, []string{"FINISHED", "CLOSED"}),
		Timeout:      timeout,
		Delay:        15 * time.Second,
		PollInterval: 15 * time.Second,
	}

	serverResp, err := stateConf.WaitForStateContext(ctx)
	return serverResp, err
}

func refreshDatabaseWatermarkEmbedTaskStatusFunc(client *golangsdk.ServiceClient, taskId string, targets []string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		tasks, err := listDatabaseWatermarkEmbedTasks(client, fmt.Sprintf("&id=%v", taskId))
		if err != nil {
			return tasks, "ERROR", err
		}

		// The status enumeration values are as follows: `WAIT`, `READY`, `RUNNING`, `ERROR`, `FINISHED` and `CLOSED`.
		// RUNNING: The task is executing.
		// READY: The task is ready to execute.
		status := utils.PathSearch("[0].task_state", tasks, "").(string)
		if status == "ERROR" {
			return nil, "FAILED", fmt.Errorf("unexpected status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return tasks, "COMPLETED", nil
		}

		return tasks, "PENDING", nil
	}
}

func listDatabaseWatermarkEmbedTasks(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/data-watermark-embed-tasks"
		limit   = 100
		// `offset` means the page number, start from 1.
		offset  = 1
		results = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		tasks := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, tasks...)
		if len(tasks) < limit {
			break
		}

		offset++
	}

	return results, nil
}

// GetDatabaseWatermarkEmbedTaskById is a method used to query the database watermark embed task by ID.
func GetDatabaseWatermarkEmbedTaskById(client *golangsdk.ServiceClient, taskId string) (interface{}, error) {
	tasks, err := listDatabaseWatermarkEmbedTasks(client, fmt.Sprintf("&id=%v", taskId))
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/data-watermark-embed-tasks",
				RequestId: "NONE",
				Body:      fmt.Appendf(nil, "the database watermark embed task (%s) does not exist", taskId),
			},
		}
	}

	return tasks[0], nil
}

func resourceDatabaseWatermarkEmbedTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		taskId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	task, err := GetDatabaseWatermarkEmbedTaskById(client, taskId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", databaseWatermarkEmbedTaskNotFoundErrCodes...),
			fmt.Sprintf("error retrieving database watermark embed task (%s)", taskId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("task_name", utils.PathSearch("task_name", task, nil)),
		d.Set("water_mark", utils.PathSearch("water_mark", task, nil)),
		d.Set("watermark_version", utils.PathSearch("watermark_version", task, nil)),
		d.Set("db_water_param", flattenDatabaseWatermarkEmbedTaskDbWaterParam(utils.PathSearch("db_water_param", task, nil))),
		d.Set("source_db_info", flattenDatabaseWatermarkEmbedTaskDbInfo(utils.PathSearch("source_db_info", task, nil))),
		d.Set("target_db_info", flattenDatabaseWatermarkEmbedTaskDbInfo(utils.PathSearch("target_db_info", task, nil))),
		d.Set("error_code", utils.PathSearch("error_code", task, nil)),
		// Optional parameters.
		d.Set("start_now", utils.PathSearch("start_now", task, nil)),
		d.Set("watermark_describe", utils.PathSearch("watermark_describe", task, nil)),
		d.Set("schedule_switch", utils.PathSearch("schedule_switch", task, nil)),
		d.Set("schedule_type", utils.PathSearch("schedule_type", task, nil)),
		d.Set("start_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("start_time", task, float64(0)).(float64))/1000, false)),
		d.Set("selected_fields", flattenDatabaseWatermarkEmbedTaskColumnInfos(utils.PathSearch("selected_fields", task,
			make([]interface{}, 0)).([]interface{}))),
		// Attributes.
		d.Set("task_state", utils.PathSearch("task_state", task, nil)),
		d.Set("task_create_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("task_create_time",
			task, float64(0)).(float64))/1000, false)),
		d.Set("task_end_time", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("task_end_time",
			task, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatabaseWatermarkEmbedTaskDbWaterParam(dbWaterParam interface{}) []map[string]interface{} {
	if dbWaterParam == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"embed_mode":    utils.PathSearch("embed_mode", dbWaterParam, nil),
			"watermark_key": utils.PathSearch("watermark_key", dbWaterParam, nil),
			"params": flattenDatabaseWatermarkEmbedTaskEmbedParams(utils.PathSearch("params", dbWaterParam,
				make([]interface{}, 0)).([]interface{})),
			"row_spacing": utils.PathSearch("row_spacing", dbWaterParam, nil),
		},
	}
}

func flattenDatabaseWatermarkEmbedTaskEmbedParams(params []interface{}) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(params))
	for _, v := range params {
		rst = append(rst, map[string]interface{}{
			"new_column_name": utils.PathSearch("new_column_name", v, nil),
			"new_column_type": utils.PathSearch("new_column_type", v, nil),
			"fake_strategy":   utils.PathSearch("fake_strategy", v, nil),
			"fake_param":      flattenDatabaseWatermarkEmbedTaskFakeParam(utils.PathSearch("fake_param", v, nil)),
		})
	}

	return rst
}

func flattenDatabaseWatermarkEmbedTaskFakeParam(fakeParam interface{}) []map[string]interface{} {
	if fakeParam == nil {
		return nil
	}

	beginTime, err := strconv.ParseInt(utils.PathSearch("date_begin", fakeParam, "").(string), 10, 64)
	if err != nil {
		log.Printf("[ERROR] error parsing 'date_begin' field to Integer: %s", err)
	}

	endTime, err := strconv.ParseInt(utils.PathSearch("date_end", fakeParam, "").(string), 10, 64)
	if err != nil {
		log.Printf("[ERROR] error parsing 'date_end' field to Integer: %s", err)
	}

	return []map[string]interface{}{
		{
			"address_accuracy":  utils.PathSearch("address_accuracy", fakeParam, nil),
			"date_begin":        utils.FormatTimeStampRFC3339(beginTime/1000, false),
			"date_end":          utils.FormatTimeStampRFC3339(endTime/1000, false),
			"random_accuracy":   utils.PathSearch("random_accuracy", fakeParam, nil),
			"random_begin":      utils.PathSearch("random_begin", fakeParam, nil),
			"random_distribute": utils.PathSearch("random_distribute", fakeParam, nil),
			"random_end":        utils.PathSearch("random_end", fakeParam, nil),
			"string_distribute": utils.PathSearch("string_distribute", fakeParam, nil),
		},
	}
}

func flattenDatabaseWatermarkEmbedTaskDbInfo(dbInfo interface{}) []map[string]interface{} {
	if dbInfo == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"db_id":       utils.PathSearch("db_id", dbInfo, nil),
			"db_name":     utils.PathSearch("db_name", dbInfo, nil),
			"db_type":     utils.PathSearch("db_type", dbInfo, nil),
			"ins_id":      utils.PathSearch("ins_id", dbInfo, nil),
			"ins_name":    utils.PathSearch("ins_name", dbInfo, nil),
			"table_name":  utils.PathSearch("table_name", dbInfo, nil),
			"schema_name": utils.PathSearch("schema_name", dbInfo, nil),
		},
	}
}

func flattenDatabaseWatermarkEmbedTaskColumnInfos(fields []interface{}) []map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(fields))
	for _, v := range fields {
		rst = append(rst, map[string]interface{}{
			"column_name": utils.PathSearch("column_name", v, nil),
			"column_type": utils.PathSearch("column_type", v, nil),
		})
	}

	return rst
}

func updateDatabaseWatermarkEmbedTask(ctx context.Context, client *golangsdk.ServiceClient, taskId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/data-watermark-embed-tasks/{id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", taskId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDatabaseWatermarkEmbedTaskBodyParams(d)),
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		_, err := client.Request("PUT", updatePath, &updateOpt)
		isRetry, err := handleOperationDatabaseWatermarkEmbedTaskError(err)
		return nil, isRetry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 15 * time.Second,
		PollInterval: 15 * time.Second,
	})

	return err
}

func resourceDatabaseWatermarkEmbedTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		taskId = d.Id()
	)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChangesExcept("enable_force_new") {
		err = updateDatabaseWatermarkEmbedTask(ctx, client, taskId, d)
		if err != nil {
			return diag.Errorf("error updating database watermark embed task (%v): %s", taskId, err)
		}

		if isImmediateStart, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "start_now").(bool); ok && isImmediateStart {
			if _, err := waitForDatabaseWatermarkEmbedTaskCompleted(ctx, client, d.Timeout(schema.TimeoutUpdate), taskId); err != nil {
				return diag.Errorf("error waiting for database watermark embed task to be completed: %s", err)
			}
		}
	}

	return resourceDatabaseWatermarkEmbedTaskRead(ctx, d, meta)
}

func handleOperationDatabaseWatermarkEmbedTaskError(err error) (bool, error) {
	if err == nil {
		// The operation was executed successfully and does not need to be executed again.
		return false, nil
	}

	if errCode, ok := err.(golangsdk.ErrDefault400); ok {
		var apiError interface{}
		if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
			return false, jsonErr
		}

		// dsc.40000028: The task is running and cannot be dispatched.
		if utils.PathSearch("error_code", apiError, "").(string) == "dsc.40000028" {
			return true, err
		}
	}

	return false, err
}

func deleteDatabaseWatermarkEmbedTask(ctx context.Context, client *golangsdk.ServiceClient, taskId string, timeout time.Duration) error {
	httpUrl := "v1/{project_id}/data-watermark-embed-tasks/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", taskId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	retryFunc := func() (interface{}, bool, error) {
		res, err := client.Request("DELETE", deletePath, &deleteOpt)
		isRetry, err := handleOperationDatabaseWatermarkEmbedTaskError(err)
		return res, isRetry, err
	}
	_, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		Timeout:      timeout,
		DelayTimeout: 15 * time.Second,
		PollInterval: 15 * time.Second,
	})

	return err
}

func resourceDatabaseWatermarkEmbedTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		taskId = d.Id()
	)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteDatabaseWatermarkEmbedTask(ctx, client, taskId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", databaseWatermarkEmbedTaskNotFoundErrCodes...),
			fmt.Sprintf("error deleting database watermark embed task (%v)", taskId),
		)
	}

	return nil
}

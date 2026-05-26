package das

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	binlogParseTaskExportNonUpdatableParams = []string{
		"user_id",
		"task_id",
		"bucket_name",
		"filter_condition",
		"filter_condition.*.db_names",
		"filter_condition.*.tb_names",
		"filter_condition.*.file_names",
		"filter_condition.*.start_time",
		"filter_condition.*.end_time",
		"filter_condition.*.types",
		"filter_condition.*.parse_double_insert",
		"filter_condition.*.columns",
		"filter_condition.*.columns.*.name",
		"filter_condition.*.columns.*.value",
	}

	binlogParseTaskExportNotFoundCodes = []string{
		"DAS.2002", // The exported task does not exist during deletion.
	}
)

// @API DAS POST /v3/{project_id}/connections/{connection_id}/binlog-parse/export
// @API DAS GET /v3/{project_id}/connections/{connection_id}/binlog-parse/get-export-task-info
// @API DAS DELETE /v3/{project_id}/connections/{connection_id}/binlog-parse/delete-export-task
func ResourceBinlogParseTaskExport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBinlogParseTaskExportCreate,
		ReadContext:   resourceBinlogParseTaskExportRead,
		UpdateContext: resourceBinlogParseTaskExportUpdate,
		DeleteContext: resourceBinlogParseTaskExportDelete,

		CustomizeDiff: config.FlexibleForceNew(binlogParseTaskExportNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceBinlogParseTaskExportImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the binlog parse task export is located.`,
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user ID of the database connection.`,
			},
			"task_id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The binlog parse task ID.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS bucket name.`,
			},

			// The API defines 'filter_condition' as 'info'. And 'info' is optional, but actually it is required.
			"filter_condition": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of database names to filter.`,
						},
						"tb_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of table names to filter.`,
						},
						"file_names": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of file names to filter.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The start time of the export range, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The end time of the export range, in RFC3339 format.`,
						},
						"types": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of SQL types to filter.`,
						},
						"columns": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The column name.`,
									},
									"value": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `The column value.`,
									},
								},
							},
							Description: `The list of columns to filter.`,
						},
						"parse_double_insert": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Whether to export UPDATE statements as two INSERT statements.`,
						},
					},
				},
				Description: `The filter conditions for the export task.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The task status.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The instance ID.`,
			},
			"last_record_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last record time, in RFC3339 format.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task creation time, in RFC3339 format.`,
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

func buildFilterCondition(filterCondition map[string]interface{}) map[string]interface{} {
	result := map[string]interface{}{
		"db_names":            utils.ValueIgnoreEmpty(filterCondition["db_names"]),
		"tb_names":            utils.ValueIgnoreEmpty(filterCondition["tb_names"]),
		"file_names":          utils.ValueIgnoreEmpty(filterCondition["file_names"]),
		"types":               utils.ValueIgnoreEmpty(filterCondition["types"]),
		"parse_double_insert": utils.ValueIgnoreEmpty(filterCondition["parse_double_insert"]),
		"columns":             buildColumns(utils.ValueIgnoreEmpty(filterCondition["columns"])),
	}

	if v, ok := filterCondition["start_time"].(string); ok && v != "" {
		result["start_time"] = utils.ConvertTimeStrToNanoTimestamp(v) / 1000
	}
	if v, ok := filterCondition["end_time"].(string); ok && v != "" {
		result["end_time"] = utils.ConvertTimeStrToNanoTimestamp(v) / 1000
	}

	return result
}

func buildColumns(columnsRaw interface{}) []map[string]interface{} {
	if columnsRaw == nil {
		return nil
	}

	columns := columnsRaw.([]interface{})
	if len(columns) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(columns))
	for _, column := range columns {
		columnMap := column.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"column_name":  columnMap["name"],
			"column_value": columnMap["value"],
		})
	}

	return result
}

func buildBinlogParseTaskExportBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		"task_id":     d.Get("task_id"),
		"bucket_name": d.Get("bucket_name"),
		"info": buildFilterCondition(
			d.Get("filter_condition").([]interface{})[0].(map[string]interface{})),
	}

	return result
}

// GetBinlogParseTaskExport queries the binlog export task info by user ID and export task ID.
func GetBinlogParseTaskExport(client *golangsdk.ServiceClient, userId, exportedTaskId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/get-export-task-info?export_task_id={export_task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{connection_id}", userId)
	getPath = strings.ReplaceAll(getPath, "{export_task_id}", exportedTaskId)

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

func resourceBinlogParseTaskExportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		userId = d.Get("user_id").(string)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/export"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{connection_id}", userId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildBinlogParseTaskExportBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS binlog parse task export: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	// The response body only has an exported task ID without any json format.
	if respBody == nil {
		return diag.Errorf("unable to find the exported task ID from the API response")
	}
	d.SetId(fmt.Sprintf("%v", respBody))

	return resourceBinlogParseTaskExportRead(ctx, d, meta)
}

func resourceBinlogParseTaskExportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		userId         = d.Get("user_id").(string)
		exportedTaskId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := GetBinlogParseTaskExport(client, userId, exportedTaskId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving DAS binlog parse task export (%s)", exportedTaskId))
	}

	// The GET API will return a 'task_id', but the meaning of this 'task_id' is exported binlog task ID, not the
	// parsed binlog task ID.
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_id", userId),
		d.Set("instance_id", utils.PathSearch("instance_id", respBody, nil)),
		d.Set("status", utils.PathSearch("task_status", respBody, nil)),
		d.Set("export_line_num", utils.PathSearch("export_line_num", respBody, nil)),
		d.Set("download_url", utils.PathSearch("download_url", respBody, nil)),
		d.Set("last_record_time", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("last_record_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_at", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBinlogParseTaskExportUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBinlogParseTaskExportDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		userId         = d.Get("user_id").(string)
		exportedTaskId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/delete-export-task"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", userId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"export_task_id": exportedTaskId,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errorCodeStr",
			binlogParseTaskExportNotFoundCodes...), fmt.Sprintf("error deleting DAS binlog parse task export (%s)", exportedTaskId))
	}

	return nil
}

func resourceBinlogParseTaskExportImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, errors.New("invalid format specified for import ID, must be <user_id>/<bucket_name>/<id>")
	}

	d.SetId(parts[2])
	return []*schema.ResourceData{d}, multierror.Append(nil,
		d.Set("user_id", parts[0]),
		d.Set("bucket_name", parts[1]),
	).ErrorOrNil()
}

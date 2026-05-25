package das

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	binlogParseTaskNonUpdatableParams = []string{
		"user_id",
		"binlog_type",
		"file_name",
	}
	binlogParseTaskNotFoundCodes = []string{
		"DAS.5520", // The parsed binlog task does not exist during query or deletion.
	}
)

// @API DAS POST /v3/{project_id}/connections/{connection_id}/binlog-parse/create-task
// @API DAS GET /v3/{project_id}/connections/{connection_id}/binlog-parse/get-task-info
// @API DAS GET /v3/{project_id}/connections/{connection_id}/binlog-parse/delete-task
func ResourceBinlogParseTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBinlogParseTaskCreate,
		ReadContext:   resourceBinlogParseTaskRead,
		UpdateContext: resourceBinlogParseTaskUpdate,
		DeleteContext: resourceBinlogParseTaskDelete,

		CustomizeDiff: config.FlexibleForceNew(binlogParseTaskNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceBinlogParseTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the binlog parse task is located.`,
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user ID of the database connection.`,
			},
			"binlog_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The binlog type.`,
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The binlog file name.`,
			},

			// Optional parameters.
			"backup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The backup ID.`,
			},

			// Attributes.
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The task status.`,
			},
			"position": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The binlog file parse position.`,
			},
			"error_message": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The error message.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task creation time, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task modification time, in RFC3339 format.`,
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

func buildBinlogParseTaskCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"binlog_type": d.Get("binlog_type"),
		"file_name":   d.Get("file_name"),
		"backup_id":   utils.ValueIgnoreEmpty(d.Get("backup_id")),
	}
}

func GetBinlogParseTask(client *golangsdk.ServiceClient, userId, taskId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/get-task-info?task_id={task_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	// The `connection_id` is the DB user ID, you can use `data.huaweicloud_das_database_users` to get it.
	getPath = strings.ReplaceAll(getPath, "{connection_id}", userId)
	getPath = strings.ReplaceAll(getPath, "{task_id}", taskId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			// Use this header key-value pair to make the API response formatted as underscores
			// eg: `fileName` to `file_name`
			"X-Source-Service": "das",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func resourceBinlogParseTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	userId := d.Get("user_id").(string)
	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/create-task"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	// The `connection_id` is the DB user ID, you can use `data.huaweicloud_das_database_users` to get it.
	createPath = strings.ReplaceAll(createPath, "{connection_id}", userId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			// Use this header key-value pair to make the API response formatted as underscores
			// eg: `fileName` to `file_name`
			"X-Source-Service": "das",
		},
		JSONBody: utils.RemoveNil(buildBinlogParseTaskCreateBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS binlog parse task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("id", respBody, float64(0)).(float64)
	if taskId == 0 {
		return diag.Errorf("unable to find the task ID from the API response")
	}
	d.SetId(strconv.Itoa(int(taskId)))

	return resourceBinlogParseTaskRead(ctx, d, meta)
}

func resourceBinlogParseTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		userId = d.Get("user_id").(string)
		taskId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := GetBinlogParseTask(client, userId, taskId)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			binlogParseTaskNotFoundCodes...), fmt.Sprintf("error retrieving DAS binlog parse task (%s)", taskId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_id", utils.PathSearch("connection_id", respBody, nil)),
		d.Set("binlog_type", utils.PathSearch("binlog_type", respBody, nil)),
		d.Set("file_name", utils.PathSearch("file_name", respBody, nil)),
		d.Set("backup_id", utils.PathSearch("backup_id", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("position", utils.PathSearch("position", respBody, nil)),
		d.Set("error_message", utils.PathSearch("err_msg", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("gmt_create", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("gmt_modified", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceBinlogParseTaskUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceBinlogParseTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		userId = d.Get("user_id").(string)
		taskId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/connections/{connection_id}/binlog-parse/delete-task?task_id={task_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	// The `connection_id` is the DB user ID, you can use `data.huaweicloud_das_database_users` to get it.
	deletePath = strings.ReplaceAll(deletePath, "{connection_id}", userId)
	deletePath = strings.ReplaceAll(deletePath, "{task_id}", taskId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			// Use this header key-value pair to make the API response formatted as underscores
			// eg: `fileName` to `file_name`
			"X-Source-Service": "das",
		},
	}

	_, err = client.Request("GET", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			binlogParseTaskNotFoundCodes...), fmt.Sprintf("error deleting DAS binlog parse task (%s)", taskId))
	}

	return nil
}

func resourceBinlogParseTaskImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <user_id>/<id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("user_id", parts[0])
}

package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

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
	instanceLogNonUpdatableParams = []string{
		"instance_id",
		"log_type",
		"log_file_name",
		"log_group_name",
		"log_stream_name",
	}
)

// @API Kafka POST /v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
// @API Kafka GET /v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}
// @API LTS GET /v2/{project_id}/groups
// @API LTS GET /v2/{project_id}/groups/{log_group_id}/streams
// @API Kafka DELETE /v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}
func ResourceInstanceLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceLogCreate,
		ReadContext:   resourceInstanceLogRead,
		UpdateContext: resourceInstanceLogUpdate,
		DeleteContext: resourceInstanceLogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceLogImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(instanceLogNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region where the instance log is located.`,
			},
			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the log task belongs.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the log.`,
			},

			// Optional parameters.
			"log_file_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the log file.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the log group.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the log stream.`,
			},

			// Attributes.
			"log_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the log group.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the log stream.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the log.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the log, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the log, in RFC3339 format.`,
			},

			// Internal parameter(s).
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

func buildInstanceLogBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"log_task_list": []map[string]interface{}{
			{
				"log_file_name":   utils.ValueIgnoreEmpty(d.Get("log_file_name")),
				"log_group_name":  utils.ValueIgnoreEmpty(d.Get("log_group_name")),
				"log_stream_name": utils.ValueIgnoreEmpty(d.Get("log_stream_name")),
			},
		},
	}
}

func enableInstanceLog(client *golangsdk.ServiceClient, d *schema.ResourceData, instanceId, logType string) (interface{}, error) {
	httpUrl := "v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	httpUrl = strings.ReplaceAll(httpUrl, "{log_type}", logType)
	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
		JSONBody:         utils.RemoveNil(buildInstanceLogBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceInstanceLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		logType    = d.Get("log_type").(string)
	)

	// Lock the resource to prevent concurrent log tasks (error_code: DMS.00501008, error_msg: log task start in progress)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := enableInstanceLog(client, d, instanceId, logType)
	if err != nil {
		return diag.Errorf("error enabling log (%s) of the instance (%s): %s", logType, instanceId, err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the log (%s) of the instance (%s) to be enabled: %s", logType, instanceId, err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceId, logType))

	return resourceInstanceLogRead(ctx, d, meta)
}

func GetInstanceLog(client *golangsdk.ServiceClient, instanceId, logType string) (interface{}, error) {
	httpUrl := "v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{log_type}", logType)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	logResponse := utils.PathSearch("log_response_list[0]", respBody, make(map[string]interface{})).(map[string]interface{})
	if len(logResponse) == 0 || utils.PathSearch("status", logResponse, "").(string) == "CLOSE" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The log (%s) of the instance (%s) has been closed.", logType, instanceId)),
			},
		}
	}

	return logResponse, nil
}

func resourceInstanceLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		logType    = d.Get("log_type").(string)
	)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	ltsClient, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	respBody, err := GetInstanceLog(client, instanceId, logType)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting log (%s) of the instance (%s)", logType, instanceId))
	}

	logGroupId := utils.PathSearch("log_group_id", respBody, "").(string)
	logStreamId := utils.PathSearch("log_stream_id", respBody, "").(string)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", respBody, nil)),
		d.Set("log_type", utils.PathSearch("log_type", respBody, nil)),
		d.Set("log_file_name", utils.PathSearch("log_file_name", respBody, nil)),
		// Attributes.
		d.Set("log_group_id", logGroupId),
		d.Set("log_stream_id", logStreamId),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at",
			respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at",
			respBody, float64(0)).(float64))/1000, false)),
	)

	if logGroupId == "" || logStreamId == "" {
		return diag.FromErr(mErr.ErrorOrNil())
	}

	logGroupName, err := getInstanceLogGroupName(ltsClient, logGroupId)
	if err != nil {
		log.Printf("[ERROR] error getting log group name of the log (%s) of the instance (%s): %s", logType, instanceId, err)
	}

	logStreamName, err := getInstanceLogStreamName(ltsClient, logGroupId, logStreamId)
	if err != nil {
		log.Printf("[ERROR] error getting log stream name of the log (%s) of the instance (%s): %s", logType, instanceId, err)
	}

	mErr = multierror.Append(mErr,
		d.Set("log_group_name", logGroupName),
		d.Set("log_stream_name", logStreamName),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getInstanceLogGroupName(client *golangsdk.ServiceClient, logGroupId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	getPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch(fmt.Sprintf("log_groups[?log_group_id=='%s']|[0].log_group_name", logGroupId), respBody, nil), nil
}

func getInstanceLogStreamName(client *golangsdk.ServiceClient, logGroupId, logStreamId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/groups/{log_group_id}/streams"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{log_group_id}", logGroupId)
	getPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch(fmt.Sprintf("log_streams[?log_stream_id=='%s']|[0].log_stream_name", logStreamId), respBody, nil), nil
}

func resourceInstanceLogUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func disableInstanceLog(client *golangsdk.ServiceClient, instanceId, logType string) (interface{}, error) {
	httpUrl := "v2/{project_id}/kafka/instances/{instance_id}/logs/{log_type}"
	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	httpUrl = strings.ReplaceAll(httpUrl, "{log_type}", logType)
	deletePath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=utf-8"},
	}

	resp, err := client.Request("DELETE", deletePath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceInstanceLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		logType    = d.Get("log_type").(string)
	)

	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := disableInstanceLog(client, instanceId, logType)
	if err != nil {
		// DMS.00501007: The log is disabled.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DMS.00501007"),
			fmt.Sprintf("error disabling log task of the instance (%s)", instanceId),
		)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the log (%s) of the instance (%s) to be disabled: %s", logType, instanceId, err)
	}

	return nil
}

func resourceInstanceLogImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<log_type>', but got '%s'", importId)
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("log_type", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

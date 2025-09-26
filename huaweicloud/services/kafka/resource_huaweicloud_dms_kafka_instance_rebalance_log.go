package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceRebalanceLogNonUpdatableParams = []string{"instance_id"}

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log
// @API Kafka GET /v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log
// @API Kafka DELETE /v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
func ResourceInstanceRebalanceLog() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceRebalanceLogCreate,
		ReadContext:   resourceInstanceRebalanceLogRead,
		UpdateContext: resourceInstanceRebalanceLogUpdate,
		DeleteContext: resourceInstanceRebalanceLogDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(instanceRebalanceLogNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `The region where the Kafka instance rebalance log is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance to which the rebalance log belongs.`,
			},
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
			"dashboard_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the dashboard.`,
			},
			"log_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the log.`,
			},
			"log_file_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the log file.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the rebalance log.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the rebalance log, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the rebalance log, in RFC3339 format.`,
			},
		},
	}
}

func resourceInstanceRebalanceLogCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log"
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	createPath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error enabling rebalance log of the instance (%s): %s", instanceId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the rebalance log of the instance (%s) to be enabled: %s", instanceId, err)
	}

	d.SetId(instanceId)

	return resourceInstanceRebalanceLogRead(ctx, d, meta)
}

func GetInstanceRebalanceLog(client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	if utils.PathSearch("status", respBody, "") == "CLOSE" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("The rebalance log of the instance (%s) has been closed.", instanceId)),
			},
		}
	}

	return respBody, nil
}

func resourceInstanceRebalanceLogRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Id()
	)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := GetInstanceRebalanceLog(client, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting rebalance log of the instance (%s)", instanceId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("instance_id", respBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_group_id", respBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_stream_id", respBody, nil)),
		d.Set("dashboard_id", utils.PathSearch("dashboard_id", respBody, nil)),
		d.Set("log_type", utils.PathSearch("log_type", respBody, nil)),
		d.Set("log_file_name", utils.PathSearch("log_file_name", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at",
			respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("updated_at",
			respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceInstanceRebalanceLogUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceRebalanceLogDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v2/kafka/{project_id}/instances/{instance_id}/log/rebalance-log"
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	httpUrl = strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)
	httpUrl = strings.ReplaceAll(httpUrl, "{instance_id}", instanceId)
	deletePath := client.Endpoint + httpUrl

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf-8",
		},
	}

	resp, err := client.Request("DELETE", deletePath, &opt)
	// DMS.00501007: The log has been closed.
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "DMS.00501007"),
			fmt.Sprintf("error disabling rebalance log of the instance (%s)", instanceId),
		)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch("job_id", respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID from the API response")
	}

	err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for the rebalance log of the instance (%s) to be disabled: %s", instanceId, err)
	}

	return nil
}

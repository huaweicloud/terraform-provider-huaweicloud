package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/schedule-task/future-executions
func DataSourceAppScheduleTaskFutureExecutions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppScheduleTaskFutureExecutionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the schedule task is located.",
			},
			"scheduled_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of execution cycle.",
			},
			"scheduled_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The interval in days for the scheduled task is to be executed.",
			},
			"day_interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The days of the weeks when the scheduled task is to be executed.",
			},
			"week_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The days of the week for execution.",
			},
			"month_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The month when the scheduled task is to be executed",
			},
			"date_list": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The days of month when the scheduled task is to be executed.",
			},
			"scheduled_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The fixed date when the scheduled task is to be executed.",
			},
			"time_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The time zone of the schedule task.",
			},
			"expire_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The expiration time of the schedule task.",
			},
			"future_executions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of future execution times of the schedule task.",
			},
		},
	}
}

func buildAppScheduleTaskFutureExecutionsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"scheduled_type": d.Get("scheduled_type"),
		"scheduled_time": d.Get("scheduled_time"),
		// Optional parameters.
		"day_interval":   utils.ValueIgnoreEmpty(d.Get("day_interval")),
		"week_list":      utils.ValueIgnoreEmpty(d.Get("week_list")),
		"month_list":     utils.ValueIgnoreEmpty(d.Get("month_list")),
		"date_list":      utils.ValueIgnoreEmpty(d.Get("date_list")),
		"scheduled_date": utils.ValueIgnoreEmpty(d.Get("scheduled_date")),
		"time_zone":      utils.ValueIgnoreEmpty(d.Get("time_zone")),
		"expire_time":    utils.ValueIgnoreEmpty(d.Get("expire_time")),
	}
}

// Only the last 5 times can be queried.
func dataSourceAppScheduleTaskFutureExecutionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/schedule-task/future-executions"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAppScheduleTaskFutureExecutionsBodyParams(d)),
	}

	resp, err := client.Request("POST", getPath, &opt)
	if err != nil {
		return diag.Errorf("error retrieving future executions of the schedule task: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("future_executions", utils.PathSearch("future_executions", respBody, nil)),
		d.Set("time_zone", utils.PathSearch("time_zone", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

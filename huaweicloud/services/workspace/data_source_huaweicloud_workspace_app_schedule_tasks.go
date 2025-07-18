package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/schedule-task
func DataSourceAppScheduleTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppScheduleTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the schedule tasks are located.",
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the schedule task.",
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the schedule task.",
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the schedule task.",
						},
						"task_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the schedule task.",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the schedule task.",
						},
						"scheduled_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The execution cycle of the schedule task.",
						},
						"scheduled_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The scheduled time of the schedule task.",
						},
						"day_interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The interval in days of the schedule task.",
						},
						"week_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The days of the week of the schedule task.",
						},
						"month_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The month of the schedule task.",
						},
						"date_list": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The days of the month of the schedule task.",
						},
						"scheduled_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fixed time of the schedule task.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time zone of the schedule task.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expire time of the schedule task.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the schedule task.",
						},
						"is_enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the schedule task is enabled.",
						},
						"task_cron": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cron expression of the schedule task.",
						},
						"next_execution_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The next execution time of the schedule task.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The create time of the schedule task, in RFC3339 format.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lastest update time of the schedule task, in RFC3339 format.",
						},
						"last_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The last execution status of the schedule task.",
						},
					},
				},
				Description: "All schedule tasks that match the filter parameters.",
			},
		},
	}
}

func buildAppScheduleTasksQueryParams(d *schema.ResourceData) string {
	params := ""
	if v, ok := d.GetOk("task_type"); ok {
		params += fmt.Sprintf("&task_type=%s", v)
	}

	if v, ok := d.GetOk("task_name"); ok {
		params += fmt.Sprintf("&task_name=%s", v)
	}

	return params
}

func queryAppScheduleTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		listPath = "v1/{project_id}/schedule-task"
		offset   = 0
		// For API, limit default value is 10.
		limit   = 100
		results = make([]interface{}, 0)
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
	)

	listPath = client.Endpoint + listPath
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildAppScheduleTasksQueryParams(d)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		listResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		scheduleTasks := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		results = append(results, scheduleTasks...)
		if len(scheduleTasks) < limit {
			break
		}

		offset += len(scheduleTasks)
	}

	return results, nil
}

func dataSourceAppScheduleTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	items, err := queryAppScheduleTasks(client, d)
	if err != nil {
		return diag.Errorf("error getting Workspace APP schedule tasks: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tasks", flattenAppScheduleTasks(items)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppScheduleTasks(tasks []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(tasks))
	for i, v := range tasks {
		result[i] = map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"task_name":           utils.PathSearch("task_name", v, nil),
			"task_type":           utils.PathSearch("task_type", v, nil),
			"scheduled_type":      utils.PathSearch("scheduled_type", v, nil),
			"scheduled_time":      utils.PathSearch("scheduled_time", v, nil),
			"day_interval":        utils.PathSearch("day_interval", v, nil),
			"week_list":           utils.PathSearch("week_list", v, nil),
			"month_list":          utils.PathSearch("month_list", v, nil),
			"date_list":           utils.PathSearch("date_list", v, nil),
			"scheduled_date":      utils.PathSearch("scheduled_date", v, nil),
			"time_zone":           utils.PathSearch("time_zone", v, nil),
			"expire_time":         utils.PathSearch("expire_time", v, nil),
			"description":         utils.PathSearch("description", v, nil),
			"is_enable":           utils.PathSearch("is_enable", v, nil),
			"task_cron":           utils.PathSearch("task_cron", v, nil),
			"next_execution_time": utils.PathSearch("next_execution_time", v, nil),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("create_time",
				v, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("update_time",
				v, "").(string))/1000, false),
			"last_status": utils.PathSearch("last_status", v, nil),
		}
	}

	return result
}

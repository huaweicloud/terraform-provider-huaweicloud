package workspace

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/scheduled-tasks
func DataSourceScheduledTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduledTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the scheduled tasks are located.`,
			},

			// Optional parameters.
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the scheduled task to be queried.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the scheduled task to be queried.`,
			},
			"scheduled_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The execution cycle type of the scheduled task to be queried.`,
			},
			"last_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The last execution status of the scheduled task to be queried.`,
			},

			// Attributes.
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the scheduled task.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the scheduled task.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the scheduled task.`,
						},
						"scheduled_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution cycle type of the scheduled task.`,
						},
						"life_cycle_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The trigger scenario type of the scheduled task.`,
						},
						"last_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last execution status of the scheduled task.`,
						},
						"next_execution_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The next execution time of the scheduled task.`,
						},
						"enable": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the scheduled task is enabled.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the scheduled task.`,
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time zone of the scheduled task.`,
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The priority of the scheduled task.`,
						},
						"wait_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The wait time after the trigger scenario for the scheduled task.`,
						},
					},
				},
				Description: `The list of the scheduled tasks that matched filter parameters.`,
			},
		},
	}
}

func buildScheduledTasksQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("task_name"); ok {
		res = fmt.Sprintf("%s&task_name=%v", res, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		res = fmt.Sprintf("%s&task_type=%v", res, v)
	}
	if v, ok := d.GetOk("scheduled_type"); ok {
		res = fmt.Sprintf("%s&scheduled_type=%v", res, v)
	}
	if v, ok := d.GetOk("last_status"); ok {
		res = fmt.Sprintf("%s&last_status=%v", res, v)
	}

	return res
}

func listScheduledTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/scheduled-tasks?limit={limit}"
		limit   = 50
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildScheduledTasksQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		scheduledTasks := utils.PathSearch("scheduled_tasks", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, scheduledTasks...)
		if len(scheduledTasks) < limit {
			break
		}
		offset += len(scheduledTasks)
	}

	return result, nil
}

func flattenScheduledTasks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		result[i] = map[string]interface{}{
			"id":                  utils.PathSearch("id", item, nil),
			"name":                utils.PathSearch("task_name", item, nil),
			"type":                utils.PathSearch("task_type", item, nil),
			"scheduled_type":      utils.PathSearch("scheduled_type", item, nil),
			"life_cycle_type":     utils.PathSearch("life_cycle_type", item, nil),
			"last_status":         utils.PathSearch("last_status", item, nil),
			"next_execution_time": utils.PathSearch("next_execution_time", item, nil),
			"enable":              utils.PathSearch("enable", item, nil),
			"description":         utils.PathSearch("description", item, nil),
			"priority":            utils.PathSearch("priority", item, nil),
			"time_zone":           utils.PathSearch("time_zone", item, nil),
			"wait_time":           utils.PathSearch("wait_time", item, nil),
		}
	}

	return result
}

func dataSourceScheduledTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	resp, err := listScheduledTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace scheduled tasks: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenScheduledTasks(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

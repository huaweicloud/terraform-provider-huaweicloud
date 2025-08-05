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

// @API Workspace GET /v1/{project_id}/schedule-task/{task_id}/execute-history
func DataSourceAppScheduleTaskExecutions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppScheduleTaskExecutionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the schedule task is located.",
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the schedule task.",
			},
			"executions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the schedule task execution record.",
						},
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the schedule task.",
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
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the schedule task execution.",
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of subtasks.",
						},
						"success_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of successful subtasks.",
						},
						"failed_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of failed subtasks.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timezone of the schedule task.",
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The begin time of the schedule task execution, in UTC format.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time of the schedule task execution, in UTC format.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the schedule task, in UTC format.",
						},
					},
				},
				Description: "The list of schedule task executions.",
			},
		},
	}
}

func queryScheduleTaskExecutions(client *golangsdk.ServiceClient, taskId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/schedule-task/{task_id}/execute-history"
		result  = make([]interface{}, 0)
		limit   = 100
		offset  = 0
		opt     = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
		}
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{task_id}", taskId)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		executeHistories := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, executeHistories...)
		if len(executeHistories) < limit {
			break
		}
		offset += len(executeHistories)
	}

	return result, nil
}

func dataSourceAppScheduleTaskExecutionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		taskId = d.Get("task_id").(string)
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	executions, err := queryScheduleTaskExecutions(client, taskId)
	if err != nil {
		return diag.Errorf("error querying schedule task(%s) executions: %s", taskId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("executions", flattenScheduleTaskExecutions(executions)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScheduleTaskExecutions(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", item, nil),
			"task_id":        utils.PathSearch("task_id", item, nil),
			"task_type":      utils.PathSearch("task_type", item, nil),
			"scheduled_type": utils.PathSearch("scheduled_type", item, nil),
			"status":         utils.PathSearch("status", item, nil),
			"total_count":    utils.PathSearch("total_count", item, nil),
			"success_count":  utils.PathSearch("success_count", item, nil),
			"failed_count":   utils.PathSearch("failed_count", item, nil),
			"time_zone":      utils.PathSearch("time_zone", item, nil),
			"begin_time":     utils.PathSearch("begin_time", item, nil),
			"end_time":       utils.PathSearch("end_time", item, nil),
			"create_time":    utils.PathSearch("create_time", item, nil),
		})
	}
	return result
}

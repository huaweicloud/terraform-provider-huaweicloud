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

// @API Workspace GET /v1/{project_id}/schedule-task/{execute_history_id}/execute-detail
func DataSourceAppScheduleTaskExecuteDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppScheduleTaskExecuteDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the schedule task execution details are located.",
			},
			"execute_history_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the schedule task execution record.",
			},
			"execute_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the sub-task to be executed.",
						},
						"execute_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the schedule task execution record.",
						},
						"server_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the server being operated.",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the server being operated.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of schedule task execution.",
						},
						"task_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the schedule task.",
						},
						"time_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The timezone of the schedule task.",
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The start time of the sub-task, in UTC format.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The end time of the sub-task, in UTC format.",
						},
						"result_code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The error code when the task execution fails.",
						},
						"result_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The reason of task failure.",
						},
					},
				},
				Description: "The list of the sub-task execution details.",
			},
		},
	}
}

func queryScheduleTaskExecuteDetails(client *golangsdk.ServiceClient, executeHistoryId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/schedule-task/{execute_history_id}/execute-detail"
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
	listPath = strings.ReplaceAll(listPath, "{execute_history_id}", executeHistoryId)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	for {
		listPathWithOfset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOfset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		executeDetails := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, executeDetails...)
		if len(executeDetails) < limit {
			break
		}

		offset += len(executeDetails)
	}

	return result, nil
}

func dataSourceAppScheduleTaskExecuteDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		executeHistoryId = d.Get("execute_history_id").(string)
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	executeDetails, err := queryScheduleTaskExecuteDetails(client, executeHistoryId)
	if err != nil {
		return diag.Errorf("error querying sub-task execution details under specified task (%s): %s", executeHistoryId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("execute_details", flattenScheduleTaskExecuteDetails(executeDetails)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenScheduleTaskExecuteDetails(executeDetails []interface{}) []map[string]interface{} {
	if len(executeDetails) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(executeDetails))
	for _, item := range executeDetails {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("id", item, nil),
			"execute_id":     utils.PathSearch("execute_id", item, nil),
			"server_id":      utils.PathSearch("server_id", item, nil),
			"server_name":    utils.PathSearch("server_name", item, nil),
			"status":         utils.PathSearch("status", item, nil),
			"task_type":      utils.PathSearch("task_type", item, nil),
			"time_zone":      utils.PathSearch("time_zone", item, nil),
			"begin_time":     utils.PathSearch("begin_time", item, nil),
			"end_time":       utils.PathSearch("end_time", item, nil),
			"result_code":    utils.PathSearch("result_code", item, nil),
			"result_message": utils.PathSearch("result_message", item, nil),
		})
	}
	return result
}

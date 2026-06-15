package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-export-task-list
func DataSourceSlowLogExportedTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSlowLogExportedTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the slow log exported tasks are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance ID.`,
			},

			// Optional parameters.
			"export_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The export type.`,
			},

			// Attributes.
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of slow log exported tasks that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The task ID.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The instance ID.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The task status.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task creation time, in RFC3339 format.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task start time, in RFC3339 format.`,
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The task end time, in RFC3339 format.`,
						},
						"download_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The download URL.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceSlowLogExportedTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	tasks, err := listSlowLogExportedTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS slow log exported tasks: %s", err)
	}

	randomUUID, err := uuid.NewUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenSlowLogExportedTasks(tasks)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listSlowLogExportedTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-log/get-slow-log-export-task-list"
		perPage = 100
		curPage = 1
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":     "application/json",
			"X-Source-Service": "das",
		},
	}

	for {
		listPathWithParams := listPath + buildSlowLogExportedTasksQueryParams(d, curPage, perPage)
		requestResp, err := client.Request("GET", listPathWithParams, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		taskList := utils.PathSearch("task_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, taskList...)

		if len(taskList) < perPage {
			break
		}
		curPage++
	}

	return result, nil
}

func buildSlowLogExportedTasksQueryParams(d *schema.ResourceData, curPage, perPage int) string {
	res := fmt.Sprintf("&cur_page=%d&per_page=%d", curPage, perPage)

	if v, ok := d.GetOk("export_type"); ok {
		res = fmt.Sprintf("%s&export_type=%s", res, v.(string))
	}

	return "?" + res[1:]
}

func flattenSlowLogExportedTasks(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("task_id", item, nil),
			"status":       utils.PathSearch("task_status", item, nil),
			"instance_id":  utils.PathSearch("instance_id", item, nil),
			"download_url": utils.PathSearch("download_url", item, nil),
			"created_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_at", item, float64(0)).(float64))/1000, false),
			"start_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("start_time", item, float64(0)).(float64))/1000, false),
			"end_time": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("end_time", item, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

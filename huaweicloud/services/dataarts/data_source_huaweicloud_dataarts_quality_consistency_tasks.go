package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/quality/consistency-tasks
func DataSourceQualityConsistencyTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQualityConsistencyTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the consistency tasks are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace.`,
			},

			// Optional parameters.
			"category_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The category ID of the consistency task.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the consistency task.`,
			},
			"schedule_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schedule status of the consistency task.`,
			},
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The start time of the last run time query interval.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The end time of the last run time query interval.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the consistency task.`,
			},

			// Attributes.
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the consistency tasks that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the consistency task.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the consistency task.`,
						},
						"category_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The category ID of the consistency task.`,
						},
						"schedule_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The schedule status of the consistency task.`,
						},
						"schedule_period": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The schedule period of the consistency task.`,
						},
						"schedule_interval": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The schedule interval of the consistency task.`,
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the consistency task.`,
						},
						"last_run_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The last run time of the consistency task.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the consistency task.`,
						},
					},
				},
			},
		},
	}
}

func buildQualityMoreHeaders(workspaceId string) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type": "application/json",
	}

	if workspaceId != "" {
		moreHeaders["workspace"] = workspaceId
	}

	return moreHeaders
}

func flattenQualityConsistencyTasks(tasks []interface{}) []map[string]interface{} {
	if len(tasks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tasks))
	for _, task := range tasks {
		// The id of API response is float64, eg: 1497195981236936700
		taskId := utils.PathSearch("id", task, float64(0))
		if v, ok := taskId.(float64); ok {
			taskId = strconv.FormatInt(int64(v), 10)
		}
		r := map[string]interface{}{
			"id":                taskId,
			"name":              utils.PathSearch("name", task, nil),
			"category_id":       utils.PathSearch("category_id", task, nil),
			"schedule_status":   utils.PathSearch("schedule_status", task, nil),
			"schedule_period":   utils.PathSearch("schedule_period", task, nil),
			"schedule_interval": utils.PathSearch("schedule_interval", task, nil),
			"creator":           utils.PathSearch("creator", task, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", task,
				float64(0)).(float64))/1000, false),
			"last_run_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", task,
				float64(0)).(float64))/1000, false),
		}

		result = append(result, r)
	}

	return result
}

func buildQualityConsistencyTasksQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("category_id"); ok {
		res = fmt.Sprintf("%s&category_id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("schedule_status"); ok {
		res = fmt.Sprintf("%s&schedule_status=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}

	return res
}

func listQualityConsistencyTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/quality/consistency-tasks?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildQualityConsistencyTasksQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildQualityMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		tasks := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasks...)
		if len(tasks) < limit {
			break
		}

		offset += len(tasks)
	}

	return result, nil
}

func dataSourceQualityConsistencyTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	tasks, err := listQualityConsistencyTasks(client, d)
	if err != nil {
		return diag.Errorf("error querying consistency tasks: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenQualityConsistencyTasks(tasks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

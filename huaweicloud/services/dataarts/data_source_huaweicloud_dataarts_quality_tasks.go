package dataarts

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

// @API DataArtsStudio GET /v2/{project_id}/quality/quality-tasks
func DataSourceQualityTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQualityTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the quality tasks are located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The workspace ID to which the quality tasks belong.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the quality task.`,
			},
			"category_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The category ID to which the quality tasks belong.`,
			},
			"schedule_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The schedule status of the quality task.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The start time of the query interval for the most recent run time.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the quality task creator.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        qualityTaskSchema(),
				Computed:    true,
				Description: `All quality tasks that match the filter parameters.`,
			},
		},
	}
}

func qualityTaskSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: `The ID of the quality task.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the quality task.`,
			},
			"category_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The category ID to which the quality task belongs.`,
			},
			"schedule_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schedule status of the quality task.`,
			},
			"schedule_period": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schedule period of the quality task.`,
			},
			"schedule_interval": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The schedule interval of the quality task.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the quality task, in RFC3339 format.`,
			},
			"last_run_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last run time of the quality task, in RFC3339 format.`,
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the task creator.`,
			},
		},
	}
	return &sc
}

func buildQualityTasksQueryParams(d *schema.ResourceData) string {
	res := ""
	if apiName, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, apiName)
	}
	if categoryId, ok := d.GetOk("category_id"); ok {
		res = fmt.Sprintf("%s&category_id=%v", res, categoryId)
	}
	if scheduleStatus, ok := d.GetOk("schedule_status"); ok {
		res = fmt.Sprintf("%s&schedule_status=%v", res, scheduleStatus)
	}
	if startTime, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, utils.ConvertTimeStrToNanoTimestamp(startTime.(string)))
	}
	if creator, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, creator)
	}
	return res
}

func queryQualityTasks(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/quality/quality-tasks?limit=100"
		offset  = 0
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildQualityTasksQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"workspace":    d.Get("workspace_id").(string),
		},
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
		if len(tasks) < 1 {
			break
		}
		result = append(result, tasks...)
		offset += len(tasks)
	}

	return result, nil
}

func flattenQualityTasks(tasks []interface{}) []interface{} {
	result := make([]interface{}, 0)

	for _, task := range tasks {
		result = append(result, map[string]interface{}{
			"id":                utils.PathSearch("id", task, nil),
			"name":              utils.PathSearch("name", task, nil),
			"category_id":       fmt.Sprintf("%v", utils.PathSearch("category_id", task, float64(0)).(float64)),
			"schedule_status":   utils.PathSearch("schedule_status", task, nil),
			"schedule_period":   utils.PathSearch("schedule_period", task, nil),
			"schedule_interval": utils.PathSearch("schedule_interval", task, nil),
			"created_at":        utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", task, float64(0)).(float64))/1000, false),
			"last_run_time":     utils.FormatTimeStampRFC3339(int64(utils.PathSearch("last_run_time", task, float64(0)).(float64))/1000, false),
			"creator":           utils.PathSearch("creator", task, nil),
		})
	}

	return result
}

func dataSourceQualityTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	qualityTasks, err := queryQualityTasks(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenQualityTasks(qualityTasks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

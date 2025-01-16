package gaussdb

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

// @API GaussDB GET /v3/{project_id}/tasks
func DataSourceOpenGaussTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenGaussTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task status.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task name.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the task list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance ID.`,
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the instance name.`,
						},
						"instance_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the iInstance status.`,
						},
						"job_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task status.`,
						},
						"process": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task progress.`,
						},
						"fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task failure cause.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task creation time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"ended_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task end time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceOpenGaussTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/tasks"
		product = "opengauss"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	listBasePath := client.Endpoint + httpUrl
	listBasePath = strings.ReplaceAll(listBasePath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	limit := 100
	offset := 1
	res := make([]map[string]interface{}, 0)
	for {
		listPath := listBasePath + buildListOpenGaussTasksQueryParams(d, limit, offset)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB OpenGauss tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		tasks := flattenOpenGaussTasks(getRespBody)
		res = append(res, tasks...)
		if len(tasks) < limit {
			break
		}
		offset++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("tasks", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListOpenGaussTasksQueryParams(d *schema.ResourceData, limit, offset int) string {
	res := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		startTime := utils.ConvertTimeStrToNanoTimestamp(v.(string))
		res = fmt.Sprintf("%s&start_time=%v", res, startTime)
	}
	if v, ok := d.GetOk("end_time"); ok {
		endTime := utils.ConvertTimeStrToNanoTimestamp(v.(string))
		res = fmt.Sprintf("%s&end_time=%v", res, endTime)
	}
	return res
}

func flattenOpenGaussTasks(resp interface{}) []map[string]interface{} {
	tasksJson := utils.PathSearch("tasks", resp, make([]interface{}, 0))
	tasksArray := tasksJson.([]interface{})
	if len(tasksArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tasksArray))
	for _, task := range tasksArray {
		result = append(result, map[string]interface{}{
			"instance_id":     utils.PathSearch("instance_info.instance_id", task, nil),
			"instance_name":   utils.PathSearch("instance_info.instance_name", task, nil),
			"instance_status": utils.PathSearch("instance_info.instance_status", task, nil),
			"job_id":          utils.PathSearch("job_id", task, nil),
			"name":            utils.PathSearch("name", task, nil),
			"status":          utils.PathSearch("status", task, nil),
			"process":         utils.PathSearch("process", task, nil),
			"fail_reason":     utils.PathSearch("fail_reason", task, nil),
			"created_at":      utils.PathSearch("created_at", task, nil),
			"ended_at":        utils.PathSearch("ended_at", task, nil),
		})
	}
	return result
}

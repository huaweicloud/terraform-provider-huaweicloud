package taurusdb

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

// @API GaussDBforMySQL GET /v3/{project_id}/immediate-jobs
func DataSourceGaussDBMysqlInstantTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlInstantTasksRead,

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
				Description: `Specifies the task execution status.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task ID.`,
			},
			"job_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task name. Value options:`,
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
			"jobs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the task details.`,
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
						"job_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task name.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task execution status.`,
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
						"order_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the order ID.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the task creation time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"ended_time": {
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

func dataSourceGaussDBMysqlInstantTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/immediate-jobs"
		product = "gaussdb"
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

	listPath := ""
	limit := 50
	offset := 1
	res := make([]map[string]interface{}, 0)
	for {
		listPath = listBasePath + buildListInstantTasksQueryParams(d, limit, offset)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GaussDB MySQL instant tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		tasks := flattenGaussDBMysqlInstantTasks(getRespBody)
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
		d.Set("jobs", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListInstantTasksQueryParams(d *schema.ResourceData, limit, offset int) string {
	res := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("job_id"); ok {
		res = fmt.Sprintf("%s&job_id=%v", res, v)
	}
	if v, ok := d.GetOk("job_name"); ok {
		res = fmt.Sprintf("%s&job_name=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	return res
}

func flattenGaussDBMysqlInstantTasks(resp interface{}) []map[string]interface{} {
	jobsJson := utils.PathSearch("jobs", resp, make([]interface{}, 0))
	jobsArray := jobsJson.([]interface{})
	if len(jobsArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(jobsArray))
	for _, job := range jobsArray {
		result = append(result, map[string]interface{}{
			"instance_id":     utils.PathSearch("instance_id", job, nil),
			"instance_name":   utils.PathSearch("instance_name", job, nil),
			"instance_status": utils.PathSearch("instance_status", job, nil),
			"job_id":          utils.PathSearch("job_id", job, nil),
			"job_name":        utils.PathSearch("job_name", job, nil),
			"status":          utils.PathSearch("status", job, nil),
			"process":         utils.PathSearch("process", job, nil),
			"fail_reason":     utils.PathSearch("fail_reason", job, nil),
			"order_id":        utils.PathSearch("order_id", job, nil),
			"created_time":    utils.PathSearch("created_time", job, nil),
			"ended_time":      utils.PathSearch("ended_time", job, nil),
		})
	}
	return result
}

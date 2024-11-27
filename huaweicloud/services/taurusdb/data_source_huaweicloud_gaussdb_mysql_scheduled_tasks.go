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

// @API GaussDBforMySQL GET /v3/{project_id}/scheduled-jobs
func DataSourceGaussDBMysqlScheduledTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlScheduledTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datastore_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlScheduledTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/scheduled-jobs"
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

	var currentTotal int
	page := 1
	var scheduledTasks []interface{}

	for {
		listPath := listBasePath + buildGaussDBMysqlScheduledTasksQueryParams(d, page)
		listResp, err := client.Request("GET", listPath, &listOpt)
		if err != nil {
			return diag.FromErr(err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}
		tasks := flattenListGaussDBMysqlScheduledTasksResponseBody(listRespBody)
		scheduledTasks = append(scheduledTasks, tasks...)
		total := utils.PathSearch("total_count", listRespBody, float64(0)).(float64)
		currentTotal += len(tasks)
		if currentTotal >= int(total) {
			break
		}
		page++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("tasks", scheduledTasks),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGaussDBMysqlScheduledTasksQueryParams(d *schema.ResourceData, page int) string {
	res := fmt.Sprintf("?limit=100&offset=%v", page)
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%s", res, v)
	}
	if v, ok := d.GetOk("job_id"); ok {
		res = fmt.Sprintf("%s&job_id=%v", res, v)
	}
	if v, ok := d.GetOk("job_name"); ok {
		res = fmt.Sprintf("%s&job_name=%v", res, v)
	}
	return res
}

func flattenListGaussDBMysqlScheduledTasksResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("schedules", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"job_id":          utils.PathSearch("job_id", v, nil),
			"instance_id":     utils.PathSearch("instance_id", v, nil),
			"instance_name":   utils.PathSearch("instance_name", v, nil),
			"instance_status": utils.PathSearch("instance_status", v, nil),
			"project_id":      utils.PathSearch("project_id", v, nil),
			"job_name":        utils.PathSearch("job_name", v, nil),
			"create_time":     utils.PathSearch("create_time", v, nil),
			"start_time":      utils.PathSearch("start_time", v, nil),
			"end_time":        utils.PathSearch("end_time", v, nil),
			"job_status":      utils.PathSearch("job_status", v, nil),
			"datastore_type":  utils.PathSearch("datastore_type", v, nil),
		})
	}
	return rst
}

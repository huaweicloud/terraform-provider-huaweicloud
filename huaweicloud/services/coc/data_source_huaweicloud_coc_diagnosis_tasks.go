package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/diagnosis/tasks
func DataSourceCocDiagnosisTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocDiagnosisTasksRead,

		Schema: map[string]*schema.Schema{
			"task_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"work_order_id": {
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
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"instance_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCocDiagnosisTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/diagnosis/tasks"
		product = "coc"
		pageNo  = 1
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath += buildGetDiagnosisTasksParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var dataRes []map[string]interface{}
	for {
		getPath := basePath + fmt.Sprintf("&page_index=%d", pageNo)
		getResp, err := client.Request("GET", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving COC diagnosis tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		diagnosisTasks := flattenCocGetDiagnosisTasks(getRespBody)
		if len(diagnosisTasks) < 1 {
			break
		}
		dataRes = append(dataRes, diagnosisTasks...)
		pageNo++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("data", dataRes),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDiagnosisTasksParams(d *schema.ResourceData) string {
	res := "?page_size=100"
	if v, ok := d.GetOk("task_id"); ok {
		res = fmt.Sprintf("%s&task_id=%v", res, v)
	}
	if v, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&type=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("region"); ok {
		res = fmt.Sprintf("%s&region=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	return res
}

func flattenCocGetDiagnosisTasks(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	dataJson := utils.PathSearch("data.data", resp, make([]interface{}, 0))
	dataArray := dataJson.([]interface{})
	if len(dataArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(dataArray))
	for _, data := range dataArray {
		result = append(result, map[string]interface{}{
			"id":            utils.PathSearch("id", data, nil),
			"code":          utils.PathSearch("code", data, nil),
			"domain_id":     utils.PathSearch("domain_id", data, nil),
			"project_id":    utils.PathSearch("project_id", data, nil),
			"user_id":       utils.PathSearch("user_id", data, nil),
			"user_name":     utils.PathSearch("user_name", data, nil),
			"progress":      utils.PathSearch("progress", data, nil),
			"work_order_id": utils.PathSearch("work_order_id", data, nil),
			"instance_id":   utils.PathSearch("instance_id", data, nil),
			"instance_name": utils.PathSearch("instance_name", data, nil),
			"type":          utils.PathSearch("type", data, nil),
			"status":        utils.PathSearch("status", data, nil),
			"start_time":    utils.PathSearch("start_time", data, nil),
			"end_time":      utils.PathSearch("end_time", data, nil),
			"instance_num":  utils.PathSearch("instance_num", data, nil),
			"os_type":       utils.PathSearch("os_type", data, nil),
			"region":        utils.PathSearch("region", data, nil),
		})
	}
	return result
}

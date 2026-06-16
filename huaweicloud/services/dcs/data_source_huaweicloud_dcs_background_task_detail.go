package dcs

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}/progress
func DataSourceDcsBackgroundTaskDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsBackgroundTaskDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remain_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"step_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsBackgroundTaskStepDetailSchema(),
			},
		},
	}
}

func dcsBackgroundTaskStepDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"step_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"step_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"step_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_step_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsBackgroundTaskSubStepDetailSchema(),
			},
		},
	}
}

func dcsBackgroundTaskSubStepDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sub_step_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_step_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub_step_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"begin_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detail": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsBackgroundTaskDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/tasks/{task_id}/progress"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{task_id}", d.Get("task_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS background task detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening DCS background task detail response: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("progress", utils.PathSearch("progress", getRespBody, nil)),
		d.Set("remain_time", utils.PathSearch("remain_time", getRespBody, nil)),
		d.Set("step_details", flattenDcsBackgroundTaskStepDetails(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDcsBackgroundTaskStepDetails(resp interface{}) []interface{} {
	curJson := utils.PathSearch("step_details", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"step_id":          utils.PathSearch("step_id", v, nil),
			"step_name":        utils.PathSearch("step_name", v, nil),
			"step_status":      utils.PathSearch("step_status", v, nil),
			"begin_time":       utils.PathSearch("begin_time", v, nil),
			"end_time":         utils.PathSearch("end_time", v, nil),
			"error_code":       utils.PathSearch("error_code", v, nil),
			"sub_step_details": flattenDcsBackgroundTaskSubStepDetails(v),
		})
	}
	return res
}

func flattenDcsBackgroundTaskSubStepDetails(resp interface{}) []interface{} {
	curJson := utils.PathSearch("sub_step_details", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"sub_step_id":     utils.PathSearch("sub_step_id", v, nil),
			"sub_step_name":   utils.PathSearch("sub_step_name", v, nil),
			"sub_step_status": utils.PathSearch("sub_step_status", v, nil),
			"begin_time":      utils.PathSearch("begin_time", v, nil),
			"end_time":        utils.PathSearch("end_time", v, nil),
			"detail":          utils.PathSearch("detail", v, nil),
			"error_code":      utils.PathSearch("error_code", v, nil),
		})
	}
	return res
}

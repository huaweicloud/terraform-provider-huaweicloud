package drs

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

// @API DRS GET /v5/{project_id}/jobs
func DataSourceDrsJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsJobsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"instance_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobsSchema(),
			},
		},
	}
}

func jobsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"job_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_multi_az": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_new_framework": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"job_action": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobActionSchema(),
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     childrenSchema(),
			},
		},
	}
}

func jobActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"available_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"unavailable_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"current_action": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func childrenSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"net_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"charging_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"billing_tag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"job_direction": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"job_mode_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_multi_az": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"node_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_new_framework": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"job_action": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobActionSchema(),
			},
		},
	}
}

func buildJobsQueryParams(d *schema.ResourceData, offset int) string {
	rst := fmt.Sprintf("?job_type=%v", d.Get("job_type"))

	if v, ok := d.GetOk("name"); ok {
		rst += fmt.Sprintf("&name=%v", v)
	}

	if v, ok := d.GetOk("status"); ok {
		rst += fmt.Sprintf("&status=%v", v)
	}

	if v, ok := d.GetOk("engine_type"); ok {
		rst += fmt.Sprintf("&engine_type=%v", v)
	}

	if v, ok := d.GetOk("net_type"); ok {
		rst += fmt.Sprintf("&net_type=%v", v)
	}

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		rst += fmt.Sprintf("&enterprise_project_id=%v", v)
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if rawArray, ok := d.Get("instance_ids").([]interface{}); ok && len(rawArray) > 0 {
		for _, v := range rawArray {
			rst += fmt.Sprintf("&instance_ids=%v", v)
		}
	}

	if v, ok := d.GetOk("instance_ip"); ok {
		rst += fmt.Sprintf("&instance_ip=%v", v)
	}

	if offset > 0 {
		rst += fmt.Sprintf("&offset=%d", offset)
	}

	return rst
}

func dataSourceDrsJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/jobs"
		offset  = 0
		result  []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithParams := requestPath + buildJobsQueryParams(d, offset)
		resp, err := client.Request("GET", requestPathWithParams, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS jobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataList := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenJobs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobs(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", item, nil),
			"name":                  utils.PathSearch("name", item, nil),
			"status":                utils.PathSearch("status", item, nil),
			"description":           utils.PathSearch("description", item, nil),
			"create_time":           utils.PathSearch("create_time", item, nil),
			"engine_type":           utils.PathSearch("engine_type", item, nil),
			"net_type":              utils.PathSearch("net_type", item, nil),
			"charging_mode":         utils.PathSearch("charging_mode", item, nil),
			"billing_tag":           utils.PathSearch("billing_tag", item, nil),
			"job_direction":         utils.PathSearch("job_direction", item, nil),
			"job_type":              utils.PathSearch("job_type", item, nil),
			"task_type":             utils.PathSearch("task_type", item, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", item, nil),
			"job_mode":              utils.PathSearch("job_mode", item, nil),
			"job_mode_role":         utils.PathSearch("job_mode_role", item, nil),
			"is_multi_az":           utils.PathSearch("is_multi_az", item, nil),
			"node_role":             utils.PathSearch("node_role", item, nil),
			"node_new_framework":    utils.PathSearch("node_new_framework", item, nil),
			"job_action":            flattenJobV5Action(utils.PathSearch("job_action", item, nil)),
			"children":              flattenChildren(utils.PathSearch("children", item, nil)),
		})
	}
	return result
}

func flattenJobV5Action(jobAction interface{}) []interface{} {
	if jobAction == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"available_actions":   utils.PathSearch("available_actions", jobAction, nil),
			"unavailable_actions": utils.PathSearch("unavailable_actions", jobAction, nil),
			"current_action":      utils.PathSearch("current_action", jobAction, nil),
		},
	}
}

func flattenChildren(children interface{}) []interface{} {
	childrenRaw, ok := children.([]interface{})
	if !ok || len(childrenRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(childrenRaw))
	for _, item := range childrenRaw {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", item, nil),
			"name":                  utils.PathSearch("name", item, nil),
			"status":                utils.PathSearch("status", item, nil),
			"description":           utils.PathSearch("description", item, nil),
			"create_time":           utils.PathSearch("create_time", item, nil),
			"engine_type":           utils.PathSearch("engine_type", item, nil),
			"net_type":              utils.PathSearch("net_type", item, nil),
			"charging_mode":         utils.PathSearch("charging_mode", item, nil),
			"billing_tag":           utils.PathSearch("billing_tag", item, nil),
			"job_direction":         utils.PathSearch("job_direction", item, nil),
			"job_type":              utils.PathSearch("job_type", item, nil),
			"task_type":             utils.PathSearch("task_type", item, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", item, nil),
			"job_mode":              utils.PathSearch("job_mode", item, nil),
			"job_mode_role":         utils.PathSearch("job_mode_role", item, nil),
			"is_multi_az":           utils.PathSearch("is_multi_az", item, nil),
			"node_role":             utils.PathSearch("node_role", item, nil),
			"node_new_framework":    utils.PathSearch("node_new_framework", item, nil),
			"job_action":            flattenJobV5Action(utils.PathSearch("job_action", item, nil)),
		})
	}
	return result
}

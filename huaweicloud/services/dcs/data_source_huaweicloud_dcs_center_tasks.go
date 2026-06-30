package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/tasks
func DataSourceDcsCenterTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsCenterTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsCenterTaskSchema(),
			},
		},
	}
}

func dcsCenterTaskSchema() *schema.Resource {
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
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsCenterTaskDetailsSchema(),
			},
			"user_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"params": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_show": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dcsCenterTaskDetailsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"old_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_capacity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_public_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"public_ip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_ssl": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"old_cache_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_cache_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"old_resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_resource_spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"old_replica_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"new_replica_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"old_cache_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"new_cache_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"replica_az": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"old_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"new_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_only_adjust_charging": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rename_commands": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"updated_config_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsCenterTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/tasks"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListCenterTasksQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DCS center tasks: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("tasks", flattenListCenterTasksBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListCenterTasksQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListCenterTasksBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("tasks", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"details":     flattenListCenterTasksDetailsBody(v),
			"user_name":   utils.PathSearch("user_name", v, nil),
			"user_id":     utils.PathSearch("user_id", v, nil),
			"params":      utils.PathSearch("params", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
			"error_code":  utils.PathSearch("error_code", v, nil),
			"enable_show": utils.PathSearch("enable_show", v, nil),
			"job_id":      utils.PathSearch("job_id", v, nil),
		})
	}
	return res
}

func flattenListCenterTasksDetailsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("details", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"old_capacity":            utils.PathSearch("old_capacity", curJson, nil),
			"new_capacity":            utils.PathSearch("new_capacity", curJson, nil),
			"enable_public_ip":        utils.PathSearch("enable_public_ip", curJson, nil),
			"public_ip_id":            utils.PathSearch("public_ip_id", curJson, nil),
			"public_ip_address":       utils.PathSearch("public_ip_address", curJson, nil),
			"enable_ssl":              utils.PathSearch("enable_ssl", curJson, nil),
			"old_cache_mode":          utils.PathSearch("old_cache_mode", curJson, nil),
			"new_cache_mode":          utils.PathSearch("new_cache_mode", curJson, nil),
			"old_resource_spec_code":  utils.PathSearch("old_resource_spec_code", curJson, nil),
			"new_resource_spec_code":  utils.PathSearch("new_resource_spec_code", curJson, nil),
			"old_replica_num":         utils.PathSearch("old_replica_num", curJson, nil),
			"new_replica_num":         utils.PathSearch("new_replica_num", curJson, nil),
			"old_cache_type":          utils.PathSearch("old_cache_type", curJson, nil),
			"new_cache_type":          utils.PathSearch("new_cache_type", curJson, nil),
			"replica_ip":              utils.PathSearch("replica_ip", curJson, nil),
			"replica_az":              utils.PathSearch("replica_az", curJson, nil),
			"group_name":              utils.PathSearch("group_name", curJson, nil),
			"old_port":                utils.PathSearch("old_port", curJson, nil),
			"new_port":                utils.PathSearch("new_port", curJson, nil),
			"is_only_adjust_charging": utils.PathSearch("is_only_adjust_charging", curJson, nil),
			"account_name":            utils.PathSearch("account_name", curJson, nil),
			"source_ip":               utils.PathSearch("source_ip", curJson, nil),
			"target_ip":               utils.PathSearch("target_ip", curJson, nil),
			"node_name":               utils.PathSearch("node_name", curJson, nil),
			"rename_commands":         utils.PathSearch("rename_commands", curJson, nil),
			"updated_config_length":   utils.PathSearch("updated_config_length", curJson, nil),
		},
	}
	return res
}

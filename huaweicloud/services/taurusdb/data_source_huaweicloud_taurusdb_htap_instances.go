package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/htap
func DataSourceTaurusDBHtapInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapInstancesRead,

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
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     htapInstancesSchema(),
			},
			"max_htap_instance_num_of_taurus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func htapInstancesSchema() *schema.Resource {
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
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_state": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_fail_error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"fail_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wait_restart_for_params": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"create_at": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_frozen": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ha_mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pay_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alter_order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_vip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"readable_node_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"proxy_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data_vip_v6": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"available_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"az_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"current_actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"object_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"updated_at": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"volume_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dedicated_resource_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_net_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"ch_master_node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/htap"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving TaurusDB HTAP instances: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("instances", flattenHtapInstancesBody(getRespBody)),
		d.Set("max_htap_instance_num_of_taurus", utils.PathSearch("max_htap_instance_num_of_taurus", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHtapInstancesBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"engine_name":           utils.PathSearch("engine_name", v, nil),
			"engine_version":        utils.PathSearch("engine_version", v, nil),
			"project_id":            utils.PathSearch("project_id", v, nil),
			"instance_state":        flattenHtapInstanceState(v),
			"create_at":             utils.PathSearch("create_at", v, nil),
			"is_frozen":             utils.PathSearch("is_frozen", v, nil),
			"ha_mode":               utils.PathSearch("ha_mode", v, nil),
			"pay_model":             utils.PathSearch("pay_model", v, nil),
			"order_id":              utils.PathSearch("order_id", v, nil),
			"alter_order_id":        utils.PathSearch("alter_order_id", v, nil),
			"data_vip":              utils.PathSearch("data_vip", v, nil),
			"readable_node_infos":   flattenHtapReadableNodeInfos(v),
			"proxy_ips":             utils.PathSearch("proxy_ips", v, nil),
			"data_vip_v6":           utils.PathSearch("data_vip_v6", v, nil),
			"port":                  utils.PathSearch("port", v, nil),
			"available_zones":       flattenHtapAvailableZones(v),
			"current_actions":       flattenHtapCurrentActions(v),
			"volume_type":           utils.PathSearch("volume_type", v, nil),
			"server_type":           utils.PathSearch("server_type", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
			"dedicated_resource_id": utils.PathSearch("dedicated_resource_id", v, nil),
			"network":               flattenHtapNetwork(v),
			"ch_master_node_id":     utils.PathSearch("ch_master_node_id", v, nil),
			"node_num":              utils.PathSearch("node_num", v, nil),
		})
	}
	return res
}

func flattenHtapInstanceState(instance interface{}) []interface{} {
	curJson := utils.PathSearch("instance_state", instance, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"instance_status":         utils.PathSearch("instance_status", curJson, nil),
			"create_fail_error_code":  utils.PathSearch("create_fail_error_code", curJson, nil),
			"fail_message":            utils.PathSearch("fail_message", curJson, nil),
			"wait_restart_for_params": utils.PathSearch("wait_restart_for_params", curJson, nil),
		},
	}
}

func flattenHtapReadableNodeInfos(instance interface{}) []interface{} {
	curJson := utils.PathSearch("readable_node_infos", instance, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"data_ip":   utils.PathSearch("data_ip", v, nil),
			"node_id":   utils.PathSearch("node_id", v, nil),
			"node_name": utils.PathSearch("node_name", v, nil),
		})
	}
	return res
}

func flattenHtapAvailableZones(instance interface{}) []interface{} {
	curJson := utils.PathSearch("available_zones", instance, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"code":        utils.PathSearch("code", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"az_type":     utils.PathSearch("az_type", v, nil),
		})
	}
	return res
}

func flattenHtapCurrentActions(instance interface{}) []interface{} {
	curJson := utils.PathSearch("current_actions", instance, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"action":     utils.PathSearch("action", v, nil),
			"object_id":  utils.PathSearch("object_id", v, nil),
			"type":       utils.PathSearch("type", v, nil),
			"job_id":     utils.PathSearch("job_id", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"created_at": utils.PathSearch("created_at", v, nil),
			"updated_at": utils.PathSearch("updated_at", v, nil),
		})
	}
	return res
}

func flattenHtapNetwork(instance interface{}) []interface{} {
	curJson := utils.PathSearch("network", instance, nil)
	if curJson == nil {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			"vpc_id":            utils.PathSearch("vpc_id", curJson, nil),
			"sub_net_id":        utils.PathSearch("sub_net_id", curJson, nil),
			"security_group_id": utils.PathSearch("security_group_id", curJson, nil),
		},
	}
}

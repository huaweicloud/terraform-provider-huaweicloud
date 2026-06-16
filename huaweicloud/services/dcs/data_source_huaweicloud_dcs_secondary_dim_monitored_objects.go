package dcs

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

// @API DCS GET /v2/{project_id}/dims/monitored-objects/{instance_id}
func DataSourceDcsSecondaryDimMonitoredObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsSecondaryDimMonitoredObjectsRead,

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
			"dim_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"router": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"children": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dim_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dim_route": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcs_instance_id": {
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
					},
				},
			},
			"dcs_cluster_redis_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dcs_cluster_redis_node": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dcs_cluster_proxy_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dcs_cluster_proxy_node": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"dcs_cluster_proxy2_node": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dcs_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dcs_cluster_proxy2_node": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsSecondaryDimMonitoredObjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	dimName := d.Get("dim_name").(string)

	httpUrl := "v2/{project_id}/dims/monitored-objects/{instance_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)
	getPath += fmt.Sprintf("?dim_name=%s", dimName)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS secondary dimension monitored objects: %s", err)
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

	routerRaw := utils.PathSearch("router", getRespBody, nil)
	routerList := make([]string, 0)

	if routerRaw != nil {
		for _, v := range routerRaw.([]interface{}) {
			if str, ok := v.(string); ok {
				routerList = append(routerList, str)
			}
		}
	}

	mErr = multierror.Append(
		d.Set("router", routerList),
		d.Set("children", flattenGetSecondaryMonitoredObjectsChildren(getRespBody)),
		d.Set("instances", flattenGetSecondaryMonitoredObjectsInstances(getRespBody)),
		d.Set("dcs_cluster_redis_node", flattenGetSecondaryMonitoredObjectsRedisNodes(getRespBody)),
		d.Set("dcs_cluster_proxy_node", flattenGetSecondaryMonitoredObjectsProxyNodes(getRespBody)),
		d.Set("dcs_cluster_proxy2_node", flattenGetSecondaryMonitoredObjectsProxy2Nodes(getRespBody)),
		d.Set("total", utils.PathSearch("total", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetSecondaryMonitoredObjectsChildren(resp interface{}) []interface{} {
	curJson := utils.PathSearch("children", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dim_name":  utils.PathSearch("dim_name", v, nil),
			"dim_route": utils.PathSearch("dim_route", v, nil),
		})
	}
	return res
}

func flattenGetSecondaryMonitoredObjectsInstances(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instances", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dcs_instance_id": utils.PathSearch("dcs_instance_id", v, nil),
			"name":            utils.PathSearch("name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
		})
	}
	return res
}

func flattenGetSecondaryMonitoredObjectsRedisNodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dcs_cluster_redis_node", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dcs_instance_id":        utils.PathSearch("dcs_instance_id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"dcs_cluster_redis_node": utils.PathSearch("dcs_cluster_redis_node", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
		})
	}
	return res
}

func flattenGetSecondaryMonitoredObjectsProxyNodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dcs_cluster_proxy_node", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dcs_instance_id":        utils.PathSearch("dcs_instance_id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"dcs_cluster_proxy_node": utils.PathSearch("dcs_cluster_proxy_node", v, nil),
			"status":                 utils.PathSearch("status", v, nil),
		})
	}
	return res
}

func flattenGetSecondaryMonitoredObjectsProxy2Nodes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dcs_cluster_proxy2_node", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"dcs_instance_id":         utils.PathSearch("dcs_instance_id", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"dcs_cluster_proxy2_node": utils.PathSearch("dcs_cluster_proxy2_node", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
		})
	}
	return res
}

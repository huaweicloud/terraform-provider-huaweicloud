package dcs

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/nodes
func DataSourceDcsClusterInstanceTopology() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsClusterInstanceTopologyRead,

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
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"redis_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"cluster_proxy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"vpcendpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"elb": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"cluster_lvs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"cluster_admin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
			"master": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     topologyInfoResource(),
			},
		},
	}
}

func topologyInfoResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"qps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"input": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"output": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"db_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_idx": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keys": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"avg_ttl": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"relation_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"relation_port": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dims": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dim_k": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dim_v": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDcsClusterInstanceTopologyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v2/{project_id}/instances/{instance_id}/nodes"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)

	query := make([]string, 0)

	if d.HasChange("start") || d.Get("start") != nil {
		startVal := d.Get("start").(int)
		query = append(query, "start="+strconv.Itoa(startVal))
	}

	if d.HasChange("limit") || d.Get("limit") != nil {
		limitVal := d.Get("limit").(int)
		query = append(query, "limit="+strconv.Itoa(limitVal))
	}

	if len(query) > 0 {
		listPath += "?" + strings.Join(query, "&")
	}

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS cluster instance topology: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("redis_server", flattenTopologyInfoList(utils.PathSearch("redis_server", listRespBody, nil))),
		d.Set("cluster_proxy", flattenTopologyInfoList(utils.PathSearch("cluster_proxy", listRespBody, nil))),
		d.Set("vpcendpoint", flattenTopologyInfoList(utils.PathSearch("vpcendpoint", listRespBody, nil))),
		d.Set("elb", flattenTopologyInfoList(utils.PathSearch("elb", listRespBody, nil))),
		d.Set("cluster_lvs", flattenTopologyInfoList(utils.PathSearch("cluster_lvs", listRespBody, nil))),
		d.Set("cluster_admin", flattenTopologyInfoList(utils.PathSearch("cluster_admin", listRespBody, nil))),
		d.Set("master", flattenTopologyInfoList(utils.PathSearch("master", listRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTopologyInfoList(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	l := v.([]interface{})
	res := make([]interface{}, 0, len(l))
	for _, item := range l {
		res = append(res, flattenTopologyInfo(item))
	}
	return res
}

func flattenTopologyInfo(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	m := v.(map[string]interface{})
	res := map[string]interface{}{
		"node_id":       utils.PathSearch("node_id", m, nil),
		"node_name":     utils.PathSearch("node_name", m, nil),
		"ip":            utils.PathSearch("ip", m, nil),
		"port":          utils.PathSearch("port", m, nil),
		"node_type":     utils.PathSearch("node_type", m, nil),
		"max_memory":    utils.PathSearch("max_memory", m, nil),
		"used_memory":   utils.PathSearch("used_memory", m, nil),
		"qps":           utils.PathSearch("qps", m, nil),
		"db_num":        utils.PathSearch("db_num", m, nil),
		"relation_ip":   utils.PathSearch("relation_ip", m, nil),
		"relation_port": utils.PathSearch("relation_port", m, nil),
		"group_id":      utils.PathSearch("group_id", m, nil),
		"status":        utils.PathSearch("status", m, nil),
	}

	if bandwidth := utils.PathSearch("bandwidth", m, nil); bandwidth != nil {
		bandwidthMap := bandwidth.(map[string]interface{})
		res["bandwidth"] = []interface{}{
			map[string]interface{}{
				"input":  utils.PathSearch("input", bandwidthMap, nil),
				"output": utils.PathSearch("output", bandwidthMap, nil),
			},
		}
	}

	if dbs := utils.PathSearch("dbs", m, nil); dbs != nil {
		dbsMap := dbs.(map[string]interface{})
		res["dbs"] = []interface{}{
			map[string]interface{}{
				"db_idx":  utils.PathSearch("db_idx", dbsMap, nil),
				"keys":    utils.PathSearch("keys", dbsMap, nil),
				"expires": utils.PathSearch("expires", dbsMap, nil),
				"avg_ttl": utils.PathSearch("avg_ttl", dbsMap, nil),
			},
		}
	}

	if dims := utils.PathSearch("dims", m, nil); dims != nil {
		dimsMap := dims.(map[string]interface{})
		res["dims"] = []interface{}{
			map[string]interface{}{
				"dim_k": utils.PathSearch("dim_k", dimsMap, nil),
				"dim_v": utils.PathSearch("dim_v", dimsMap, nil),
			},
		}
	}

	return res
}

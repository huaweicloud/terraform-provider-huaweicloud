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

// @API DCS GET /v2/{project_id}/instances/{instance_id}/nodes
func DataSourceDcsInstanceTopology() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsInstanceTopologyRead,

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
			"redis_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"cluster_lvs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"cluster_admin": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"cluster_proxy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"master": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"vpc_endpoint": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
			"elb": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyInfoSchema(),
			},
		},
	}
}

func dcsInstanceTopologyInfoSchema() *schema.Resource {
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
				Elem:     dcsInstanceTopologyBandWidthSchema(),
			},
			"db_num": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsInstanceTopologyKeySpaceSchema(),
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
				Elem:     dcsInstanceTopologyDimsInfoSchema(),
			},
		},
	}
}

func dcsInstanceTopologyBandWidthSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dcsInstanceTopologyKeySpaceSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dcsInstanceTopologyDimsInfoSchema() *schema.Resource {
	return &schema.Resource{
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
	}
}

func dataSourceDcsInstanceTopologyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/nodes"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS instance topology: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("redis_server", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "redis_server")),
		d.Set("cluster_lvs", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "cluster_lvs")),
		d.Set("cluster_admin", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "cluster_admin")),
		d.Set("cluster_proxy", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "cluster_proxy")),
		d.Set("master", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "master")),
		d.Set("vpc_endpoint", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "vpcendpoint")),
		d.Set("elb", flattenGetDcsInstanceTopologyInfoBody(getRespBody, "elb")),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetDcsInstanceTopologyInfoBody(resp interface{}, param string) []interface{} {
	curJson := utils.PathSearch(param, resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"node_id":       utils.PathSearch("node_id", v, nil),
			"node_name":     utils.PathSearch("node_name", v, nil),
			"ip":            utils.PathSearch("ip", v, nil),
			"port":          utils.PathSearch("port", v, nil),
			"node_type":     utils.PathSearch("node_type", v, nil),
			"max_memory":    utils.PathSearch("max_memory", v, nil),
			"used_memory":   utils.PathSearch("used_memory", v, nil),
			"qps":           utils.PathSearch("qps", v, nil),
			"bandwidth":     flattenGetDcsInstanceTopologyBandWidthBody(v),
			"db_num":        utils.PathSearch("db_num", v, nil),
			"dbs":           flattenGetDcsInstanceTopologyKeySpaceBody(v),
			"relation_ip":   utils.PathSearch("relation_ip", v, nil),
			"relation_port": utils.PathSearch("relation_port", v, nil),
			"group_id":      utils.PathSearch("group_id", v, nil),
			"status":        utils.PathSearch("status", v, nil),
			"dims":          flattenGetDcsInstanceTopologyDimsInfoBody(v),
		})
	}
	return res
}

func flattenGetDcsInstanceTopologyBandWidthBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("bandwidth", resp, nil)
	if curJson == nil {
		return nil
	}
	res := []interface{}{
		map[string]interface{}{
			"input":  utils.PathSearch("input", curJson, nil),
			"output": utils.PathSearch("output", curJson, nil),
		},
	}
	return res
}

func flattenGetDcsInstanceTopologyKeySpaceBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dbs", resp, nil)
	if curJson == nil {
		return nil
	}
	res := []interface{}{
		map[string]interface{}{
			"db_idx":  utils.PathSearch("db_idx", curJson, nil),
			"keys":    utils.PathSearch("keys", curJson, nil),
			"expires": utils.PathSearch("expires", curJson, nil),
			"avg_ttl": utils.PathSearch("avg_ttl", curJson, nil),
		},
	}
	return res
}

func flattenGetDcsInstanceTopologyDimsInfoBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dims", resp, nil)
	if curJson == nil {
		return nil
	}
	res := []interface{}{
		map[string]interface{}{
			"dim_k": utils.PathSearch("dim_k", curJson, nil),
			"dim_v": utils.PathSearch("dim_v", curJson, nil),
		},
	}
	return res
}

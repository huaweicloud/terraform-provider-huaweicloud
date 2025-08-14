package hss

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

// @API HSS POST /v5/{project_id}/container/kubernetes/clusters/configs/batch-query
func DataSourceContainerKubernetesClustersConfigs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEConfigsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource. If omitted, the provider-level region will be used.",
			},
			"cluster_info_list": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies the cluster information list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the cluster ID.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the cluster name.",
						},
					},
				},
			},
			"cluster_id_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Specifies the cluster ID list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project that the server belongs to.",
			},
			"data_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of cluster configurations.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster ID.",
						},
						"protect_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of protected nodes.",
						},
						"protect_interrupt_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of nodes with interrupted protection.",
						},
						"protect_degradation_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of nodes with degraded protection.",
						},
						"unprotect_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of unprotected nodes.",
						},
						"node_total_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of nodes.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster name.",
						},
						"charging_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The charging mode.",
						},
						"prefer_packet_cycle": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to prefer packet cycle.",
						},
						"protect_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protection type.",
						},
						"protect_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protection status.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cluster type.",
						},
						"fail_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The failure reason.",
						},
					},
				},
			},
		},
	}
}

func buildCCEConfigsClusterInfoList(d *schema.ResourceData) []map[string]interface{} {
	clusterInfoList, ok := d.Get("cluster_info_list").([]interface{})
	if !ok || len(clusterInfoList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(clusterInfoList))
	for _, v := range clusterInfoList {
		clusterInfo, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		result = append(result, map[string]interface{}{
			"cluster_id":   clusterInfo["cluster_id"],
			"cluster_name": clusterInfo["cluster_name"],
		})
	}

	return result
}

func buildCCEConfigsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"cluster_info_list": buildCCEConfigsClusterInfoList(d),
	}

	if clusterIdList, ok := d.Get("cluster_id_list").([]interface{}); ok && len(clusterIdList) > 0 {
		rst["cluster_id_list"] = utils.ExpandToStringList(clusterIdList)
	}

	return rst
}

func dataSourceCCEConfigsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/container/kubernetes/clusters/configs/batch-query"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	if epsId != "" {
		requestPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCCEConfigsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS container kubernetes clusters configs: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data_list", flattenDataListAttribute(dataList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataListAttribute(dataList []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataList))
	for _, v := range dataList {
		result = append(result, map[string]interface{}{
			"cluster_id":                   utils.PathSearch("cluster_id", v, nil),
			"protect_node_num":             utils.PathSearch("protect_node_num", v, nil),
			"protect_interrupt_node_num":   utils.PathSearch("protect_interrupt_node_num", v, nil),
			"protect_degradation_node_num": utils.PathSearch("protect_degradation_node_num", v, nil),
			"unprotect_node_num":           utils.PathSearch("unprotect_node_num", v, nil),
			"node_total_num":               utils.PathSearch("node_total_num", v, nil),
			"cluster_name":                 utils.PathSearch("cluster_name", v, nil),
			"charging_mode":                utils.PathSearch("charging_mode", v, nil),
			"prefer_packet_cycle":          utils.PathSearch("prefer_packet_cycle", v, nil),
			"protect_type":                 utils.PathSearch("protect_type", v, nil),
			"protect_status":               utils.PathSearch("protect_status", v, nil),
			"cluster_type":                 utils.PathSearch("cluster_type", v, nil),
			"fail_reason":                  utils.PathSearch("fail_reason", v, nil),
		})
	}
	return result
}

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

// @API HSS GET /v5/{project_id}/kubernetes/clusters
func DataSourceContainerKubernetesClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"load_agent_info": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"scene": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"total_nodes_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"active_nodes_number": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creation_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agent_installed_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agent_install_failed_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agent_not_install_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agent_ds_install_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_ds_failed_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_operate_timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"last_scan_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sys_vul_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_vul_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"emg_vul_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"risk_assess_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"sec_comp_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesClustersQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}
	if d.Get("load_agent_info").(bool) {
		queryParams = fmt.Sprintf("%s&load_agent_info=%v", queryParams, true)
	}
	if v, ok := d.GetOk("scene"); ok {
		queryParams = fmt.Sprintf("%s&scene=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerKubernetesClustersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		epsId          = cfg.GetEnterpriseProjectID(d)
		product        = "hss"
		httpUrl        = "v5/{project_id}/kubernetes/clusters"
		offset         = 0
		result         = make([]interface{}, 0)
		totalNum       float64
		lastUpdateTime float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerKubernetesClustersQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		getResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS container kubernetes clusters: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		lastUpdateTime = utils.PathSearch("last_update_time", getRespBody, float64(0)).(float64)
		totalNum = utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		clusterInfoListResp := utils.PathSearch("cluster_info_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(clusterInfoListResp) == 0 {
			break
		}

		result = append(result, clusterInfoListResp...)
		offset += len(clusterInfoListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("last_update_time", lastUpdateTime),
		d.Set("total_num", totalNum),
		d.Set("cluster_info_list", flattenContainerKubernetesClustersDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesClustersDataList(clusterInfoListResp []interface{}) []interface{} {
	if len(clusterInfoListResp) == 0 {
		return nil
	}

	result := make([]interface{}, len(clusterInfoListResp))
	for i, clusterInfo := range clusterInfoListResp {
		result[i] = map[string]interface{}{
			"id":                       utils.PathSearch("id", clusterInfo, nil),
			"cluster_name":             utils.PathSearch("cluster_name", clusterInfo, nil),
			"cluster_id":               utils.PathSearch("cluster_id", clusterInfo, nil),
			"cluster_type":             utils.PathSearch("cluster_type", clusterInfo, nil),
			"status":                   utils.PathSearch("status", clusterInfo, nil),
			"version":                  utils.PathSearch("version", clusterInfo, nil),
			"total_nodes_number":       utils.PathSearch("total_nodes_number", clusterInfo, nil),
			"active_nodes_number":      utils.PathSearch("active_nodes_number", clusterInfo, nil),
			"creation_timestamp":       utils.PathSearch("creation_timestamp", clusterInfo, nil),
			"agent_installed_num":      utils.PathSearch("agent_installed_num", clusterInfo, nil),
			"agent_install_failed_num": utils.PathSearch("agent_install_failed_num", clusterInfo, nil),
			"agent_not_install_num":    utils.PathSearch("agent_not_install_num", clusterInfo, nil),
			"agent_ds_install_status":  utils.PathSearch("agent_ds_install_status", clusterInfo, nil),
			"agent_ds_failed_reason":   utils.PathSearch("agent_ds_failed_reason", clusterInfo, nil),
			"last_operate_timestamp":   utils.PathSearch("last_operate_timestamp", clusterInfo, nil),
			"last_scan_time":           utils.PathSearch("last_scan_time", clusterInfo, nil),
			"sys_vul_num":              utils.PathSearch("sys_vul_num", clusterInfo, nil),
			"app_vul_num":              utils.PathSearch("app_vul_num", clusterInfo, nil),
			"emg_vul_num":              utils.PathSearch("emg_vul_num", clusterInfo, nil),
			"risk_assess_num":          utils.PathSearch("risk_assess_num", clusterInfo, nil),
			"sec_comp_num":             utils.PathSearch("sec_comp_num", clusterInfo, nil),
		}
	}

	return result
}

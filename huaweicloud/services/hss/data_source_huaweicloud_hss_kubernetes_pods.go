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

// @API HSS GET /v5/{project_id}/kubernetes/pods
func DataSourceKubernetesPods() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesPodsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pod_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pod_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
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
					},
				},
			},
			"pod_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"pod_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_limit": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detect_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildKubernetesPodsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("pod_name"); ok {
		queryParams = fmt.Sprintf("%s&pod_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace_name"); ok {
		queryParams = fmt.Sprintf("%s&namespace_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceKubernetesPodsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		totalNum float64
		httpUrl  = "v5/{project_id}/kubernetes/pods"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildKubernetesPodsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS kubernetes pods: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// In the API documentation, the return of `data_ist` and `pod_info_ist` cannot determine the pagination logic,
	// so the resource temporarily does not use pagination parameters
	// and only sets the `limit` to the maximum value of `200`.
	totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	podInfoList := utils.PathSearch("pod_info_list", respBody, make([]interface{}, 0)).([]interface{})

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenKubernetesPodsDataList(dataList)),
		d.Set("pod_info_list", flattenKubernetesPodsPodInfoList(podInfoList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKubernetesPodsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"pod_name":       utils.PathSearch("pod_name", v, nil),
			"namespace_name": utils.PathSearch("namespace_name", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"node_name":      utils.PathSearch("node_name", v, nil),
			"cpu":            utils.PathSearch("cpu", v, nil),
			"memory":         utils.PathSearch("memory", v, nil),
			"cpu_limit":      utils.PathSearch("cpu_limit", v, nil),
			"memory_limit":   utils.PathSearch("memory_limit", v, nil),
			"node_ip":        utils.PathSearch("node_ip", v, nil),
			"pod_ip":         utils.PathSearch("pod_ip", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"region_id":      utils.PathSearch("region_id", v, nil),
			"id":             utils.PathSearch("id", v, nil),
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"cluster_type":   utils.PathSearch("cluster_type", v, nil),
		})
	}

	return rst
}

func flattenKubernetesPodsPodInfoList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"pod_name":       utils.PathSearch("pod_name", v, nil),
			"namespace_name": utils.PathSearch("namespace_name", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"cpu":            utils.PathSearch("cpu", v, nil),
			"memory":         utils.PathSearch("memory", v, nil),
			"cpu_limit":      utils.PathSearch("cpu_limit", v, nil),
			"memory_limit":   utils.PathSearch("memory_limit", v, nil),
			"pod_ip":         utils.PathSearch("pod_ip", v, nil),
			"protect_status": utils.PathSearch("protect_status", v, nil),
			"detect_result":  utils.PathSearch("detect_result", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
		})
	}

	return rst
}

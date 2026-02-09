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

// @API HSS GET /v5/{project_id}/container/kubernetes/clusters/daemonsets
func DataSourceContainerKubernetesClustersDaemonsets() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesClustersDaemonsetsRead,

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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"show_cluster_log_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_status": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"upgradeful_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"err_running_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"err_access_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"latest_version": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"agent_version": {
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
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ds_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"desired_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"current_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ready_num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"cluster_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"installed_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"combined_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_log_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"invoked_service": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"registry_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"registry_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"registry_addr": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"registry_username": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"namespace": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesClustersDaemonsetsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}
	if d.Get("show_cluster_log_status").(bool) {
		queryParams = fmt.Sprintf("%s&show_cluster_log_status=%v", queryParams, true)
	}
	if d.Get("access_status").(bool) {
		queryParams = fmt.Sprintf("%s&access_status=%v", queryParams, true)
	}

	return queryParams
}

func dataSourceContainerKubernetesClustersDaemonsetsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		product       = "hss"
		epsId         = cfg.GetEnterpriseProjectID(d)
		httpUrl       = "v5/{project_id}/container/kubernetes/clusters/daemonsets"
		result        = make([]interface{}, 0)
		offset        = 0
		totalNum      float64
		upgradefulNum float64
		errRunningNum float64
		errAccessNum  float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerKubernetesClustersDaemonsetsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container kubernetes clusters daemonsets: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		upgradefulNum = utils.PathSearch("upgradeful_num", respBody, float64(0)).(float64)
		errRunningNum = utils.PathSearch("err_running_num", respBody, float64(0)).(float64)
		errAccessNum = utils.PathSearch("err_access_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("upgradeful_num", upgradefulNum),
		d.Set("err_running_num", errRunningNum),
		d.Set("err_access_num", errAccessNum),
		d.Set("data_list", flattenDaemonsetsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDaemonsetsDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"latest_version":     utils.PathSearch("latest_version", v, nil),
			"agent_version":      utils.PathSearch("agent_version", v, nil),
			"cluster_name":       utils.PathSearch("cluster_name", v, nil),
			"cluster_id":         utils.PathSearch("cluster_id", v, nil),
			"namespace":          utils.PathSearch("namespace", v, nil),
			"cluster_type":       utils.PathSearch("cluster_type", v, nil),
			"node_num":           utils.PathSearch("node_num", v, nil),
			"ds_info":            flattenDataListDsInfo(utils.PathSearch("ds_info", v, nil)),
			"cluster_status":     utils.PathSearch("cluster_status", v, nil),
			"installed_status":   utils.PathSearch("installed_status", v, nil),
			"access_status":      utils.PathSearch("access_status", v, nil),
			"combined_status":    utils.PathSearch("combined_status", v, nil),
			"failed_message":     utils.PathSearch("failed_message", v, nil),
			"cluster_log_status": utils.PathSearch("cluster_log_status", v, nil),
			"invoked_service":    utils.PathSearch("invoked_service", v, nil),
			"registry_info": flattenDataListRegistryInfo(
				utils.PathSearch("registry_info", v, nil)),
		})
	}

	return result
}

func flattenDataListDsInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"desired_num": utils.PathSearch("desired_num", resp, nil),
			"current_num": utils.PathSearch("current_num", resp, nil),
			"ready_num":   utils.PathSearch("ready_num", resp, nil),
		},
	}
}

func flattenDataListRegistryInfo(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"registry_type":     utils.PathSearch("registry_type", resp, nil),
			"registry_addr":     utils.PathSearch("registry_addr", resp, nil),
			"registry_username": utils.PathSearch("registry_username", resp, nil),
			"namespace":         utils.PathSearch("namespace", resp, nil),
		},
	}
}

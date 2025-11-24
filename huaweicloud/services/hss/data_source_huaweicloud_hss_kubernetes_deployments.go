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

// @API HSS GET /v5/{project_id}/kubernetes/deployments
func DataSourceKubernetesDeployments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesDeploymentsRead,

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
			"deployment_name": {
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
			"last_update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resources_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
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
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pods_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"image_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"match_labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"val": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"create_time": {
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
					},
				},
			},
		},
	}
}

func buildKubernetesDeploymentsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("deployment_name"); ok {
		queryParams = fmt.Sprintf("%s&deployment_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace_name"); ok {
		queryParams = fmt.Sprintf("%s&namespace_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceKubernetesDeploymentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "hss"
		epsId          = cfg.GetEnterpriseProjectID(d)
		result         = make([]interface{}, 0)
		offset         = 0
		totalNum       float64
		lastUpdateTime float64
		typeResp       string
		httpUrl        = "v5/{project_id}/kubernetes/deployments"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildKubernetesDeploymentsQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS kubernetes deployments: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		lastUpdateTime = utils.PathSearch("last_update_time", respBody, float64(0)).(float64)
		typeResp = utils.PathSearch("type", respBody, "").(string)
		resourcesInfoResp := utils.PathSearch("resources_info_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(resourcesInfoResp) == 0 {
			break
		}

		result = append(result, resourcesInfoResp...)
		offset += len(resourcesInfoResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("last_update_time", lastUpdateTime),
		d.Set("type", typeResp),
		d.Set("resources_info_list", flattenKubernetesDeploymentsResourcesInfoList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenKubernetesDeploymentsResourcesInfoList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"name":           utils.PathSearch("name", v, nil),
			"namespace_name": utils.PathSearch("namespace_name", v, nil),
			"cluster_name":   utils.PathSearch("cluster_name", v, nil),
			"status":         utils.PathSearch("status", v, nil),
			"protect_status": utils.PathSearch("protect_status", v, nil),
			"pods_num":       utils.PathSearch("pods_num", v, nil),
			"image_name":     utils.PathSearch("image_name", v, nil),
			"match_labels": flattenDeploymentsMatchLabels(
				utils.PathSearch("match_labels", v, make([]interface{}, 0)).([]interface{})),
			"create_time":              utils.PathSearch("create_time", v, nil),
			"agent_installed_num":      utils.PathSearch("agent_installed_num", v, nil),
			"agent_install_failed_num": utils.PathSearch("agent_install_failed_num", v, nil),
			"agent_not_install_num":    utils.PathSearch("agent_not_install_num", v, nil),
		})
	}

	return rst
}

func flattenDeploymentsMatchLabels(labels []interface{}) []interface{} {
	if len(labels) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"key": utils.PathSearch("key", label, nil),
			"val": utils.PathSearch("val", label, nil),
		})
	}

	return result
}

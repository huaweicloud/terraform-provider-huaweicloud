package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/kubernetes/endpoint/detail
func DataSourceContainerKubernetesEndpointDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesEndpointDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"association_service": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"endpoint_pod_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pod_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"endpoint_port_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesEndpointDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?id=%s", d.Get("endpoint_id").(string))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerKubernetesEndpointDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/kubernetes/endpoint/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerKubernetesEndpointDetailQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS container kubernetes endpoint detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("endpoint_id").(string))

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("service_name", utils.PathSearch("service_name", getRespBody, nil)),
		d.Set("namespace", utils.PathSearch("namespace", getRespBody, nil)),
		d.Set("creation_timestamp", utils.PathSearch("creation_timestamp", getRespBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", getRespBody, nil)),
		d.Set("labels", utils.PathSearch("labels", getRespBody, nil)),
		d.Set("association_service", utils.PathSearch("association_service", getRespBody, nil)),
		d.Set("endpoint_pod_list", flattenContainerKubernetesEndpointDetailPodList(
			utils.PathSearch("endpoint_pod_list", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("endpoint_port_list", flattenContainerKubernetesEndpointDetailPortList(
			utils.PathSearch("endpoint_port_list", getRespBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesEndpointDetailPodList(endpointPodListResp []interface{}) []interface{} {
	if len(endpointPodListResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(endpointPodListResp))
	for _, v := range endpointPodListResp {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"endpoint_id": utils.PathSearch("endpoint_id", v, nil),
			"pod_ip":      utils.PathSearch("pod_ip", v, nil),
			"pod_name":    utils.PathSearch("pod_name", v, nil),
			"available":   utils.PathSearch("available", v, nil),
		})
	}

	return result
}

func flattenContainerKubernetesEndpointDetailPortList(endpointPortListResp []interface{}) []interface{} {
	if len(endpointPortListResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(endpointPortListResp))
	for _, v := range endpointPortListResp {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("id", v, nil),
			"endpoint_id":  utils.PathSearch("endpoint_id", v, nil),
			"name":         utils.PathSearch("name", v, nil),
			"protocol":     utils.PathSearch("protocol", v, nil),
			"port":         utils.PathSearch("port", v, nil),
			"app_protocol": utils.PathSearch("app_protocol", v, nil),
		})
	}

	return result
}

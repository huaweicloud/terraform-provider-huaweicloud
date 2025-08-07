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

// @API HSS GET /v5/{project_id}/kubernetes/endpoints
func DataSourceContainerKubernetesEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerKubernetesEndpointsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
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
			"endpoint_info_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
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
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"association_service": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerKubernetesEndpointsQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if v, ok := d.GetOk("name"); ok {
		queryParams = fmt.Sprintf("%s&name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_name"); ok {
		queryParams = fmt.Sprintf("%s&cluster_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerKubernetesEndpointsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		epsId          = cfg.GetEnterpriseProjectID(d)
		product        = "hss"
		httpUrl        = "v5/{project_id}/kubernetes/endpoints"
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
	requestPath += buildContainerKubernetesEndpointsQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		getResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS container kubernetes endpoints: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", getRespBody, float64(0)).(float64)
		lastUpdateTime = utils.PathSearch("last_update_time", getRespBody, float64(0)).(float64)
		endpointInfoListResp := utils.PathSearch("endpoint_info_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(endpointInfoListResp) == 0 {
			break
		}

		result = append(result, endpointInfoListResp...)
		offset += len(endpointInfoListResp)
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
		d.Set("endpoint_info_list", flattenContainerKubernetesEndpointsDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerKubernetesEndpointsDataList(endpointInfoListResp []interface{}) []interface{} {
	if len(endpointInfoListResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(endpointInfoListResp))
	for _, v := range endpointInfoListResp {
		result = append(result, map[string]interface{}{
			"id":                  utils.PathSearch("id", v, nil),
			"name":                utils.PathSearch("name", v, nil),
			"service_name":        utils.PathSearch("service_name", v, nil),
			"namespace":           utils.PathSearch("namespace", v, nil),
			"creation_timestamp":  utils.PathSearch("creation_timestamp", v, nil),
			"cluster_name":        utils.PathSearch("cluster_name", v, nil),
			"cluster_type":        utils.PathSearch("cluster_type", v, nil),
			"association_service": utils.PathSearch("association_service", v, nil),
		})
	}

	return result
}

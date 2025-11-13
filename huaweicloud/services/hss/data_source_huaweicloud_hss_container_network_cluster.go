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

// @API HSS GET /v5/{project_id}/container-network/cluster
func DataSourceContainerNetworkCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkClusterRead,

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
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:     schema.TypeString,
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"protection_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerNetworkClusterQueryParams(epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceContainerNetworkClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "hss"
		epsId    = cfg.GetEnterpriseProjectID(d)
		result   = make([]interface{}, 0)
		offset   = 0
		totalNum = 0
		httpUrl  = "v5/{project_id}/container-network/cluster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildContainerNetworkClusterQueryParams(epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS container network cluster: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataList) == 0 {
			break
		}

		result = append(result, dataList...)
		offset += len(dataList)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("data_list", flattenContainerNetworkClusterDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerNetworkClusterDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"cluster_id":        utils.PathSearch("cluster_id", v, nil),
			"cluster_name":      utils.PathSearch("cluster_name", v, nil),
			"cluster_type":      utils.PathSearch("cluster_type", v, nil),
			"version":           utils.PathSearch("version", v, nil),
			"mode":              utils.PathSearch("mode", v, nil),
			"namespace_num":     utils.PathSearch("namespace_num", v, nil),
			"policy_num":        utils.PathSearch("policy_num", v, nil),
			"protection_status": utils.PathSearch("protection_status", v, nil),
		})
	}

	return rst
}

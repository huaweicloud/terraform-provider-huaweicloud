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

// @API HSS GET /v5/{project_id}/container-network/{cluster_id}/policy
func DataSourceContainerNetworkPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNetworkPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keyword": {
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
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_content": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerNetworkPoliciesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("namespace"); ok {
		queryParams = fmt.Sprintf("%s&namespace=%v", queryParams, v)
	}
	if v, ok := d.GetOk("keyword"); ok {
		queryParams = fmt.Sprintf("%s&keyword=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceContainerNetworkPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		product        = "hss"
		epsId          = cfg.GetEnterpriseProjectID(d)
		result         = make([]interface{}, 0)
		offset         = 0
		totalNum       = 0
		lastUpdateTime = 0
		clusterID      = d.Get("cluster_id").(string)
		httpUrl        = "v5/{project_id}/container-network/{cluster_id}/policy"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)
	requestPath += buildContainerNetworkPoliciesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving container network policies: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = int(utils.PathSearch("total_num", respBody, float64(0)).(float64))
		lastUpdateTime = int(utils.PathSearch("last_update_time", respBody, float64(0)).(float64))
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
		d.Set("last_update_time", lastUpdateTime),
		d.Set("data_list", flattenContainerNetworkPoliciesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerNetworkPoliciesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"policy_id":      utils.PathSearch("policy_id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"namespace":      utils.PathSearch("namespace", v, nil),
			"policy_content": utils.PathSearch("policy_content", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
			"deploy_status":  utils.PathSearch("deploy_status", v, nil),
		})
	}

	return rst
}

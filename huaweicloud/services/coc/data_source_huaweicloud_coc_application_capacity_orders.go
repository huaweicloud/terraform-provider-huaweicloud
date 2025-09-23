package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC GET /v1/capacity/order
func DataSourceCocApplicationCapacityOrders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocApplicationCapacityOrdersRead,

		Schema: map[string]*schema.Schema{
			"cloud_service_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"component_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rank_list": {
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
									"value": {
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

func buildApplicationCapacityOrdersParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?provider=%v&type=%v", d.Get("cloud_service_name"), d.Get("type"))
	if v, ok := d.GetOk("application_id"); ok {
		res = fmt.Sprintf("%s&application_id=%v", res, v)
	}
	if v, ok := d.GetOk("component_id"); ok {
		res = fmt.Sprintf("%s&component_id=%v", res, v)
	}
	if v, ok := d.GetOk("group_id"); ok {
		res = fmt.Sprintf("%s&group_id=%v", res, v)
	}

	return res
}

func dataSourceCocApplicationCapacityOrdersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/capacity/order"
		product = "coc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPath := basePath + buildApplicationCapacityOrdersParams(d)
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving COC application capacity orders: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocGetApplicationCapacityOrders(
			utils.PathSearch("data", getRespBody, nil))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCocGetApplicationCapacityOrders(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"type": utils.PathSearch("type", raw, nil),
				"rank_list": flattenCocGetApplicationCapacityOrdersRankList(
					utils.PathSearch("rank_list", raw, nil)),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocGetApplicationCapacityOrdersRankList(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"id":    utils.PathSearch("id", raw, nil),
				"name":  utils.PathSearch("name", raw, nil),
				"value": utils.PathSearch("value", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

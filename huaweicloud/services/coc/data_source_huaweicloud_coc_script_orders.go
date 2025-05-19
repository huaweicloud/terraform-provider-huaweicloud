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

// @API COC GET /v1/job/script/orders
func DataSourceCocScriptOrders() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocScriptOrdersRead,

		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"order_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"execute_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_created": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_finished": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"execute_costs": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"creator": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"properties": dataDataProperties1(),
					},
				},
			},
		},
	}
}

func dataDataProperties1() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"region_ids": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func dataSourceCocScriptOrdersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/job/script/orders"
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

	marker := 0.0
	res := make([]map[string]interface{}, 0)
	for {
		getPath := basePath + buildGetScriptOrdersParams(d, marker)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving COC script orders: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		scriptOrders, nextMarker := flattenCocGetScriptOrders(getRespBody)
		if len(scriptOrders) < 1 {
			break
		}
		res = append(res, scriptOrders...)
		marker = nextMarker
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("data", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetScriptOrdersParams(d *schema.ResourceData, marker float64) string {
	res := "?limit=100"
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if v, ok := d.GetOk("creator"); ok {
		res = fmt.Sprintf("%s&creator=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if marker != 0 {
		res = fmt.Sprintf("%s&marker=%v", res, int(marker))
	}

	return res
}

func flattenCocGetScriptOrders(resp interface{}) ([]map[string]interface{}, float64) {
	scriptOrdersJson := utils.PathSearch("data.data", resp, make([]interface{}, 0))
	scriptOrdersArray := scriptOrdersJson.([]interface{})
	if len(scriptOrdersArray) == 0 {
		return nil, 0
	}

	result := make([]map[string]interface{}, 0, len(scriptOrdersArray))
	var marker float64
	for _, scriptOrder := range scriptOrdersArray {
		result = append(result, map[string]interface{}{
			"order_id":      utils.PathSearch("order_id", scriptOrder, nil),
			"order_name":    utils.PathSearch("order_name", scriptOrder, nil),
			"execute_uuid":  utils.PathSearch("execute_uuid", scriptOrder, nil),
			"gmt_created":   utils.PathSearch("gmt_created", scriptOrder, nil),
			"gmt_finished":  utils.PathSearch("gmt_finished", scriptOrder, nil),
			"execute_costs": utils.PathSearch("execute_costs", scriptOrder, nil),
			"creator":       utils.PathSearch("creator", scriptOrder, nil),
			"status":        utils.PathSearch("status", scriptOrder, nil),
			"properties":    flattenScriptOrderProperties(utils.PathSearch("properties", scriptOrder, nil)),
		})
		marker = utils.PathSearch("order_id", scriptOrder, float64(0)).(float64)
	}
	return result, marker
}

func flattenScriptOrderProperties(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"region_ids": utils.PathSearch("region_ids", param, nil),
		},
	}

	return rst
}

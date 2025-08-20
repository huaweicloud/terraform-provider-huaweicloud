package aad

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

// @API AAD GET /v2/aad/domains/waf-info/flow/bandwidth
func DataSourceBandwidthCurve() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBandwidthCurveRead,

		Schema: map[string]*schema.Schema{
			"value_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domains": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recent": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// In the API documentation, `overseas_type` is int type.
			// But it needs to meet the scenario of `0`, so it is defined here as a string type.
			"overseas_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"curve": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"out": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBandwidthCurveQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?value_type=%v", d.Get("value_type"))

	if v, ok := d.GetOk("domains"); ok {
		queryParams = fmt.Sprintf("%s&domains=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("recent"); ok {
		queryParams = fmt.Sprintf("%s&recent=%v", queryParams, v)
	}
	if v, ok := d.GetOk("overseas_type"); ok {
		queryParams = fmt.Sprintf("%s&overseas_type=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBandwidthCurveRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/domains/waf-info/flow/bandwidth"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildBandwidthCurveQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD bandwidth curve: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("curve",
			flattenBandwidthCurve(utils.PathSearch("curve", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBandwidthCurve(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"in":   utils.PathSearch("in", v, nil),
			"out":  utils.PathSearch("out", v, nil),
			"time": utils.PathSearch("time", v, nil),
		})
	}

	return rst
}

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

// @API AAD GET /v2/aad/domains/waf-info/flow/qps
func DataSourceQPSCurve() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceQPSCurveRead,

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
			"overseas_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"curve": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"attack": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"basic": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cc": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"custom_custom": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildQPSCurveQueryParams(d *schema.ResourceData) string {
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

func dataSourceQPSCurveRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/domains/waf-info/flow/qps"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildQPSCurveQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD QPS curve: %s", err)
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
		d.Set("curve", flattenQPSCurve(utils.PathSearch("curve", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenQPSCurve(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		result = append(result, map[string]interface{}{
			"time":          utils.PathSearch("time", v, nil),
			"total":         utils.PathSearch("total", v, nil),
			"attack":        utils.PathSearch("attack", v, nil),
			"basic":         utils.PathSearch("basic", v, nil),
			"cc":            utils.PathSearch("cc", v, nil),
			"custom_custom": utils.PathSearch("custom_custom", v, nil),
		})
	}

	return result
}

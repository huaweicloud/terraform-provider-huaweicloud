package cfw

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

// @API CFW GET /v1/{project_id}/cfw/logs/traffic-trend
func DataSourceTrafficLogTrend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrafficLogTrendRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agg_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vgw_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agg_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"in_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"out_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildTrafficLogTrendQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&log_type=%v&agg_type=%v",
		d.Get("fw_instance_id"), d.Get("log_type"), d.Get("agg_type"))

	if v, ok := d.GetOk("range"); ok {
		res = fmt.Sprintf("%s&range=%v", res, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if rawArray, ok := d.Get("vgw_id").([]interface{}); ok {
		for _, v := range rawArray {
			res = fmt.Sprintf("%s&vgw_id=%s", res, v.(string))
		}
	}
	if rawArray, ok := d.Get("ip").([]interface{}); ok {
		for _, v := range rawArray {
			res = fmt.Sprintf("%s&ip=%s", res, v.(string))
		}
	}

	return res
}

func dataSourceTrafficLogTrendRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/traffic-trend"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildTrafficLogTrendQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW traffic log trend: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenTrafficLogTrendData(
			utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenTrafficLogTrendData(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, data := range dataResp {
		result = append(result, map[string]interface{}{
			"agg_time": utils.PathSearch("agg_time", data, nil),
			"bps":      utils.PathSearch("bps", data, nil),
			"in_bps":   utils.PathSearch("in_bps", data, nil),
			"out_bps":  utils.PathSearch("out_bps", data, nil),
		})
	}

	return result
}

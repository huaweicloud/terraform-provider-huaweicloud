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

// @API CFW GET /v1/{project_id}/cfw/logs/flow-trend
func DataSourceFlowLogTrend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowLogTrendRead,

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
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
			"range": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"direction": {
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
			"asset_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vpc": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agg_time": {
							Type:     schema.TypeInt,
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

func buildFlowLogTrendQueryParams(d *schema.ResourceData) string {
	req := fmt.Sprintf("?fw_instance_id=%v&log_type=%v", d.Get("fw_instance_id"), d.Get("log_type"))

	if v, ok := d.GetOk("range"); ok {
		req = fmt.Sprintf("%s&range=%v", req, v)
	}
	if v, ok := d.GetOk("start_time"); ok {
		req = fmt.Sprintf("%s&start_time=%v", req, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		req = fmt.Sprintf("%s&end_time=%v", req, v)
	}
	if rawArray, ok := d.Get("vgw_id").([]interface{}); ok {
		for _, v := range rawArray {
			req = fmt.Sprintf("%s&vgw_id=%s", req, v.(string))
		}
	}
	if v, ok := d.GetOk("asset_type"); ok {
		req = fmt.Sprintf("%s&asset_type=%v", req, v)
	}
	if rawArray, ok := d.Get("ip").([]interface{}); ok {
		for _, v := range rawArray {
			req = fmt.Sprintf("%s&ip=%s", req, v.(string))
		}
	}
	if rawArray, ok := d.Get("vpc").([]interface{}); ok {
		for _, v := range rawArray {
			req = fmt.Sprintf("%s&vpc=%s", req, v.(string))
		}
	}

	return req
}

func dataSourceFlowLogTrendRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/flow-trend"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildFlowLogTrendQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW flow log trend: %s", err)
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
		d.Set("records", flattenFlowLogTrendRecords(
			utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlowLogTrendRecords(recordsRaw []interface{}) []interface{} {
	if len(recordsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(recordsRaw))
	for _, raw := range recordsRaw {
		result = append(result, map[string]interface{}{
			"agg_time": utils.PathSearch("agg_time", raw, nil),
			"in_bps":   utils.PathSearch("in_bps", raw, nil),
			"out_bps":  utils.PathSearch("out_bps", raw, nil),
		})
	}

	return result
}

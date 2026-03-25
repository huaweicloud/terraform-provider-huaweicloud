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

// @API CFW GET /v1/{project_id}/cfw/logs/flow-statistic
func DataSourceFlowLogStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowLogStatisticsRead,

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
			"item": {
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
			// The API documentation is of type int, but here it is changed to string to support scenarios set to `0`.
			"size": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"apps": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     flowLogStatisticsItemVOSchema(),
						},
						"associate_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"item": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agg_start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agg_end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     flowLogStatisticsItemVOSchema(),
						},
						"region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_byte": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"response_byte": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"sessions": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"src_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     flowLogStatisticsItemVOSchema(),
						},
						"dst_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     flowLogStatisticsItemVOSchema(),
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flowLogStatisticsItemVOSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildFlowLogStatisticsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&log_type=%v&item=%v",
		d.Get("fw_instance_id"), d.Get("log_type"), d.Get("item"))

	if v, ok := d.GetOk("range"); ok {
		res = fmt.Sprintf("%s&range=%v", res, v)
	}
	if v, ok := d.GetOk("direction"); ok {
		res = fmt.Sprintf("%s&direction=%v", res, v)
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
	if v, ok := d.GetOk("asset_type"); ok {
		res = fmt.Sprintf("%s&asset_type=%v", res, v)
	}
	if v, ok := d.GetOk("size"); ok {
		res = fmt.Sprintf("%s&size=%v", res, v)
	}

	return res
}

func dataSourceFlowLogStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/flow-statistic"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildFlowLogStatisticsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW flow log statistics: %s", err)
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
		d.Set("records", flattenFlowLogStatisticsRecords(
			utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlowLogStatisticsRecords(records []interface{}) []interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"apps": flattenFlowLogStatisticsItemVOs(
				utils.PathSearch("apps", record, make([]interface{}, 0)).([]interface{})),
			"associate_instance_type": utils.PathSearch("associate_instance_type", record, nil),
			"device_name":             utils.PathSearch("device_name", record, nil),
			"item":                    utils.PathSearch("item", record, nil),
			"last_time":               utils.PathSearch("last_time", record, nil),
			"agg_start_time":          utils.PathSearch("agg_start_time", record, nil),
			"agg_end_time":            utils.PathSearch("agg_end_time", record, nil),
			"ports": flattenFlowLogStatisticsItemVOs(
				utils.PathSearch("ports", record, make([]interface{}, 0)).([]interface{})),
			"region":        utils.PathSearch("region", record, nil),
			"request_byte":  utils.PathSearch("request_byte", record, nil),
			"response_byte": utils.PathSearch("response_byte", record, nil),
			"sessions":      utils.PathSearch("sessions", record, nil),
			"tags":          utils.PathSearch("tags", record, nil),
			"src_ip": flattenFlowLogStatisticsItemVOs(
				utils.PathSearch("src_ip", record, make([]interface{}, 0)).([]interface{})),
			"dst_ip": flattenFlowLogStatisticsItemVOs(
				utils.PathSearch("dst_ip", record, make([]interface{}, 0)).([]interface{})),
			"protocol": utils.PathSearch("protocol", record, nil),
		})
	}

	return result
}

func flattenFlowLogStatisticsItemVOs(rawResp []interface{}) []interface{} {
	if len(rawResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawResp))
	for _, raw := range rawResp {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", raw, nil),
			"name":  utils.PathSearch("name", raw, nil),
			"value": utils.PathSearch("value", raw, nil),
		})
	}

	return result
}

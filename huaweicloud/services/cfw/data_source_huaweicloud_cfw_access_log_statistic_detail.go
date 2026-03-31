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

// @API CFW GET /v1/{project_id}/cfw/logs/top-access-detail
func DataSourceAccessLogStatisticDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessLogStatisticDetailRead,

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
			"item": {
				Type:     schema.TypeString,
				Required: true,
			},
			"item_id": {
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
			"log_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: accessLogStatisticDetailDataSchema(),
				},
			},
		},
	}
}

func accessLogStatisticDetailDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"dst_ip_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dst_port_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"hit_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"protocol_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"recent_end_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"recent_start_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"record_total": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"records": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: accessLogInfoSchema(),
			},
		},
		"rule_hit_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"src_ip_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func accessLogInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"app": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"dst_region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_province_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_province_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_city_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_city_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hit_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"log_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rule_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rule_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"rule_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"src_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"src_region_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_region_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_province_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_province_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_city_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_city_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"vgw_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"qos_rule_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"qos_rule_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"qos_rule_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func buildAccessLogStatisticDetailQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&item=%v&item_id=%v",
		d.Get("fw_instance_id"), d.Get("item"), d.Get("item_id"))

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
	if v, ok := d.GetOk("log_type"); ok {
		res = fmt.Sprintf("%s&log_type=%v", res, v)
	}

	return res
}

func dataSourceAccessLogStatisticDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/top-access-detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAccessLogStatisticDetailQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW access log statistic detail: %s", err)
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
		d.Set("data", flattenAccessLogStatisticDetailData(
			utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessLogStatisticDetailData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"dst_ip_count":      utils.PathSearch("dst_ip_count", dataResp, nil),
		"dst_port_count":    utils.PathSearch("dst_port_count", dataResp, nil),
		"hit_count":         utils.PathSearch("hit_count", dataResp, nil),
		"protocol_count":    utils.PathSearch("protocol_count", dataResp, nil),
		"recent_end_time":   utils.PathSearch("recent_end_time", dataResp, nil),
		"recent_start_time": utils.PathSearch("recent_start_time", dataResp, nil),
		"record_total":      utils.PathSearch("record_total", dataResp, nil),
		"records": flattenAccessLogInfos(
			utils.PathSearch("records", dataResp, make([]interface{}, 0)).([]interface{})),
		"rule_hit_count": utils.PathSearch("rule_hit_count", dataResp, nil),
		"src_ip_count":   utils.PathSearch("src_ip_count", dataResp, nil),
	}}
}

func flattenAccessLogInfos(rawResp []interface{}) []interface{} {
	if len(rawResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawResp))
	for _, raw := range rawResp {
		result = append(result, map[string]interface{}{
			"action":            utils.PathSearch("action", raw, nil),
			"app":               utils.PathSearch("app", raw, nil),
			"url":               utils.PathSearch("url", raw, nil),
			"dst_host":          utils.PathSearch("dst_host", raw, nil),
			"dst_ip":            utils.PathSearch("dst_ip", raw, nil),
			"dst_port":          utils.PathSearch("dst_port", raw, nil),
			"dst_region_id":     utils.PathSearch("dst_region_id", raw, nil),
			"dst_region_name":   utils.PathSearch("dst_region_name", raw, nil),
			"dst_province_id":   utils.PathSearch("dst_province_id", raw, nil),
			"dst_province_name": utils.PathSearch("dst_province_name", raw, nil),
			"dst_city_id":       utils.PathSearch("dst_city_id", raw, nil),
			"dst_city_name":     utils.PathSearch("dst_city_name", raw, nil),
			"hit_time":          utils.PathSearch("hit_time", raw, nil),
			"log_id":            utils.PathSearch("log_id", raw, nil),
			"protocol":          utils.PathSearch("protocol", raw, nil),
			"rule_id":           utils.PathSearch("rule_id", raw, nil),
			"rule_name":         utils.PathSearch("rule_name", raw, nil),
			"rule_type":         utils.PathSearch("rule_type", raw, nil),
			"src_ip":            utils.PathSearch("src_ip", raw, nil),
			"src_port":          utils.PathSearch("src_port", raw, nil),
			"src_region_id":     utils.PathSearch("src_region_id", raw, nil),
			"src_region_name":   utils.PathSearch("src_region_name", raw, nil),
			"src_province_id":   utils.PathSearch("src_province_id", raw, nil),
			"src_province_name": utils.PathSearch("src_province_name", raw, nil),
			"src_city_id":       utils.PathSearch("src_city_id", raw, nil),
			"src_city_name":     utils.PathSearch("src_city_name", raw, nil),
			"vgw_id":            utils.PathSearch("vgw_id", raw, nil),
			"qos_rule_id":       utils.PathSearch("qos_rule_id", raw, nil),
			"qos_rule_name":     utils.PathSearch("qos_rule_name", raw, nil),
			"qos_rule_type":     utils.PathSearch("qos_rule_type", raw, nil),
		})
	}

	return result
}

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

// @API CFW GET /v1/{project_id}/cfw/logs/attack-detail
func DataSourceAttackLogStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAttackLogStatisticsRead,

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
			"action": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"item": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			// This field is a numeric type in the API, but we're defining it as a string type here.
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
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"attack_rule_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"attack_type_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dst_ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"dst_port_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"records": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackLogsStatisticsRecordsSchema(),
						},
						"src_ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataAttackLogsStatisticsRecordsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attack_rule": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attack_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"attack_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"direction": {
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
			"event_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"src_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"real_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tag": {
				Type:     schema.TypeInt,
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
		},
	}
}

func buildAttackLogStatisticsQueryParams(d *schema.ResourceData) string {
	rst := fmt.Sprintf(
		"?fw_instance_id=%s&log_type=%s&action=%d&item=%s&value=%s",
		d.Get("fw_instance_id").(string),
		d.Get("log_type").(string),
		d.Get("action").(int),
		d.Get("item").(string),
		d.Get("value").(string),
	)

	if v, ok := d.GetOk("range"); ok {
		rst += fmt.Sprintf("&range=%s", v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		rst += fmt.Sprintf("&start_time=%s", v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		rst += fmt.Sprintf("&end_time=%s", v.(string))
	}

	rawArray, ok := d.Get("vgw_id").([]interface{})
	if ok {
		for _, v := range rawArray {
			rst += fmt.Sprintf("&vgw_id=%s", v.(string))
		}
	}

	return rst
}

func dataSourceAttackLogStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/attack-detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAttackLogStatisticsQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW attack log statistics: %s", err)
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

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("data", flattenLogStatisticsData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogStatisticsData(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"app_count":         utils.PathSearch("app_count", respBody, nil),
		"attack_rule_count": utils.PathSearch("attack_rule_count", respBody, nil),
		"attack_type_count": utils.PathSearch("attack_type_count", respBody, nil),
		"count":             utils.PathSearch("count", respBody, nil),
		"dst_ip_count":      utils.PathSearch("dst_ip_count", respBody, nil),
		"dst_port_count":    utils.PathSearch("dst_port_count", respBody, nil),
		"end_time":          utils.PathSearch("end_time", respBody, nil),
		"records":           flattenDataRecords(utils.PathSearch("records", respBody, make([]interface{}, 0)).([]interface{})),
		"src_ip_count":      utils.PathSearch("src_ip_count", respBody, nil),
		"start_time":        utils.PathSearch("start_time", respBody, nil),
		"total":             utils.PathSearch("total", respBody, nil),
	}}
}

func flattenDataRecords(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		result = append(result, map[string]interface{}{
			"action":            utils.PathSearch("action", v, nil),
			"app":               utils.PathSearch("app", v, nil),
			"attack_rule":       utils.PathSearch("attack_rule", v, nil),
			"attack_rule_id":    utils.PathSearch("attack_rule_id", v, nil),
			"attack_type":       utils.PathSearch("attack_type", v, nil),
			"direction":         utils.PathSearch("direction", v, nil),
			"dst_ip":            utils.PathSearch("dst_ip", v, nil),
			"dst_port":          utils.PathSearch("dst_port", v, nil),
			"dst_region_id":     utils.PathSearch("dst_region_id", v, nil),
			"dst_region_name":   utils.PathSearch("dst_region_name", v, nil),
			"dst_province_id":   utils.PathSearch("dst_province_id", v, nil),
			"dst_province_name": utils.PathSearch("dst_province_name", v, nil),
			"dst_city_id":       utils.PathSearch("dst_city_id", v, nil),
			"dst_city_name":     utils.PathSearch("dst_city_name", v, nil),
			"event_time":        utils.PathSearch("event_time", v, nil),
			"level":             utils.PathSearch("level", v, nil),
			"protocol":          utils.PathSearch("protocol", v, nil),
			"source":            utils.PathSearch("source", v, nil),
			"src_ip":            utils.PathSearch("src_ip", v, nil),
			"real_ip":           utils.PathSearch("real_ip", v, nil),
			"tag":               utils.PathSearch("tag", v, nil),
			"src_port":          utils.PathSearch("src_port", v, nil),
			"src_region_id":     utils.PathSearch("src_region_id", v, nil),
			"src_region_name":   utils.PathSearch("src_region_name", v, nil),
			"src_province_id":   utils.PathSearch("src_province_id", v, nil),
			"src_province_name": utils.PathSearch("src_province_name", v, nil),
			"src_city_id":       utils.PathSearch("src_city_id", v, nil),
			"src_city_name":     utils.PathSearch("src_city_name", v, nil),
			"vgw_id":            utils.PathSearch("vgw_id", v, nil),
		})
	}

	return result
}

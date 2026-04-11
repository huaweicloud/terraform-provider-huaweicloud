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

// @API CFW POST /v1/{project_id}/cfw/{fw_instance_id}/logs
func DataSourceCfwLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogsRead,

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
			"start_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"filters": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"field": {
							Type:     schema.TypeString,
							Required: true,
						},
						"operator": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"log_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"next_date": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: logsRecordSchema(),
				},
			},
		},
	}
}

func logsRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bytes": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"direction": {
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
		"end_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"log_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"packets": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"protocol": {
			Type:     schema.TypeString,
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
		"start_time": {
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
		"sctp_verification_tag": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"sctp_is_handshake_flow": {
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
		"qos_channel_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"qos_channel_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"qos_drop_packets": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"qos_drop_bytes": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"qos_rule_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"qos_channel_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"action": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hit_time": {
			Type:     schema.TypeInt,
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
		"event_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"level": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"packet": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"source": {
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
	}
}

func buildLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time"),
		"end_time":   d.Get("end_time"),
		"log_type":   d.Get("log_type"),
		"type":       d.Get("type"),
	}

	if filtersInput := d.Get("filters").([]interface{}); len(filtersInput) > 0 {
		bodyParams["filters"] = buildLogsFiltersBodyParams(filtersInput)
	}

	if v, ok := d.GetOk("log_id"); ok && v.(string) != "" {
		bodyParams["log_id"] = v
	}
	if v, ok := d.GetOk("next_date"); ok {
		bodyParams["next_date"] = v
	}

	return bodyParams
}

func buildLogsFiltersBodyParams(filtersInput []interface{}) []interface{} {
	if len(filtersInput) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(filtersInput))
	for _, item := range filtersInput {
		filterInput, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		filter := map[string]interface{}{
			"field":    filterInput["field"],
			"operator": filterInput["operator"],
		}
		if values, ok := filterInput["values"].([]interface{}); ok && len(values) > 0 {
			filter["values"] = utils.ExpandToStringList(values)
		}

		result = append(result, filter)
	}

	return result
}

func buildLogsLimitAndOffsetBodyParams(bodyParams map[string]interface{}, limit, offset int) map[string]interface{} {
	// The `limit` is required.
	bodyParams["limit"] = limit
	if offset > 0 {
		bodyParams["offset"] = offset
	}

	return bodyParams
}

func dataSourceLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/{fw_instance_id}/logs"
		// The maximum `limit` is `1024`.
		limit  = 1024
		offset = 0
		result = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", d.Get("fw_instance_id").(string))
	logsBodyParams := utils.RemoveNil(buildLogsBodyParams(d))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         logsBodyParams,
	}

	for {
		requestOpt.JSONBody = buildLogsLimitAndOffsetBodyParams(logsBodyParams, limit, offset)
		resp, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving CFW logs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		recordsResp := utils.PathSearch("data.records", respBody, make([]interface{}, 0)).([]interface{})
		if len(recordsResp) == 0 || len(recordsResp) < limit {
			break
		}

		result = append(result, recordsResp...)
		offset += len(recordsResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("records", flattenLogsRecords(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenLogsRecords(recordsResp []interface{}) []interface{} {
	if len(recordsResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(recordsResp))
	for _, raw := range result {
		dstRegionID := utils.PathSearch("dst_region_id", raw, nil)
		var dstRegionIDStr interface{}
		if dstRegionID != nil {
			dstRegionIDStr = fmt.Sprint(dstRegionID)
		}

		result = append(result, map[string]interface{}{
			"app":                    utils.PathSearch("app", raw, nil),
			"bytes":                  utils.PathSearch("bytes", raw, nil),
			"direction":              utils.PathSearch("direction", raw, nil),
			"dst_host":               utils.PathSearch("dst_host", raw, nil),
			"dst_ip":                 utils.PathSearch("dst_ip", raw, nil),
			"dst_port":               utils.PathSearch("dst_port", raw, nil),
			"end_time":               utils.PathSearch("end_time", raw, nil),
			"log_id":                 utils.PathSearch("log_id", raw, nil),
			"packets":                utils.PathSearch("packets", raw, nil),
			"protocol":               utils.PathSearch("protocol", raw, nil),
			"src_ip":                 utils.PathSearch("src_ip", raw, nil),
			"src_port":               utils.PathSearch("src_port", raw, nil),
			"start_time":             utils.PathSearch("start_time", raw, nil),
			"dst_region_id":          dstRegionIDStr,
			"dst_region_name":        utils.PathSearch("dst_region_name", raw, nil),
			"dst_province_id":        utils.PathSearch("dst_province_id", raw, nil),
			"dst_province_name":      utils.PathSearch("dst_province_name", raw, nil),
			"dst_city_id":            utils.PathSearch("dst_city_id", raw, nil),
			"dst_city_name":          utils.PathSearch("dst_city_name", raw, nil),
			"src_region_id":          utils.PathSearch("src_region_id", raw, nil),
			"src_region_name":        utils.PathSearch("src_region_name", raw, nil),
			"src_province_id":        utils.PathSearch("src_province_id", raw, nil),
			"src_province_name":      utils.PathSearch("src_province_name", raw, nil),
			"src_city_id":            utils.PathSearch("src_city_id", raw, nil),
			"src_city_name":          utils.PathSearch("src_city_name", raw, nil),
			"vgw_id":                 utils.PathSearch("vgw_id", raw, nil),
			"sctp_verification_tag":  utils.PathSearch("sctp_verification_tag", raw, nil),
			"sctp_is_handshake_flow": utils.PathSearch("sctp_is_handshake_flow", raw, nil),
			"qos_rule_id":            utils.PathSearch("qos_rule_id", raw, nil),
			"qos_rule_name":          utils.PathSearch("qos_rule_name", raw, nil),
			"qos_channel_id":         utils.PathSearch("qos_channel_id", raw, nil),
			"qos_channel_name":       utils.PathSearch("qos_channel_name", raw, nil),
			"qos_drop_packets":       utils.PathSearch("qos_drop_packets", raw, nil),
			"qos_drop_bytes":         utils.PathSearch("qos_drop_bytes", raw, nil),
			"qos_rule_type":          utils.PathSearch("qos_rule_type", raw, nil),
			"qos_channel_type":       utils.PathSearch("qos_channel_type", raw, nil),
			"action":                 utils.PathSearch("action", raw, nil),
			"url":                    utils.PathSearch("url", raw, nil),
			"hit_time":               utils.PathSearch("hit_time", raw, nil),
			"rule_id":                utils.PathSearch("rule_id", raw, nil),
			"rule_name":              utils.PathSearch("rule_name", raw, nil),
			"rule_type":              utils.PathSearch("rule_type", raw, nil),
			"attack_rule":            utils.PathSearch("attack_rule", raw, nil),
			"attack_rule_id":         utils.PathSearch("attack_rule_id", raw, nil),
			"attack_type":            utils.PathSearch("attack_type", raw, nil),
			"event_time":             utils.PathSearch("event_time", raw, nil),
			"level":                  utils.PathSearch("level", raw, nil),
			"packet":                 utils.PathSearch("packet", raw, nil),
			"source":                 utils.PathSearch("source", raw, nil),
			"real_ip":                utils.PathSearch("real_ip", raw, nil),
			"tag":                    utils.PathSearch("tag", raw, nil),
		})
	}

	return result
}

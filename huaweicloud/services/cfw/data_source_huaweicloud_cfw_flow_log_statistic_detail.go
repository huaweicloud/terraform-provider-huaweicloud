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

// @API CFW GET /v1/{project_id}/cfw/logs/flow-detail
func DataSourceFlowLogStatisticDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowLogStatisticDetailRead,

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
			"value": {
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
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bytes": {
							Type:     schema.TypeFloat,
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
							Elem: &schema.Resource{
								Schema: flowLogStatisticDetailRecordSchema(),
							},
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
						"src_ip_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func flowLogStatisticDetailRecordSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"app": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"bytes": {
			Type:     schema.TypeFloat,
			Computed: true,
		},
		"dst_associate_instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_device_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_port": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dst_host": {
			Type:     schema.TypeString,
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
		"end_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"protocol": {
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
		"src_associate_instance_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_device_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"src_ip": {
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
		"start_time": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func buildFlowLogStatisticDetailQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&log_type=%v&item=%v&value=%v",
		d.Get("fw_instance_id"), d.Get("log_type"), d.Get("item"), d.Get("value"))

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

	return res
}

func dataSourceFlowLogStatisticDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/flow-detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildFlowLogStatisticDetailQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW flow log statistic detail: %s", err)
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
		d.Set("data", flattenFlowLogStatisticDetailData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFlowLogStatisticDetailData(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"app_count":      utils.PathSearch("app_count", raw, nil),
			"bytes":          utils.PathSearch("bytes", raw, nil),
			"dst_ip_count":   utils.PathSearch("dst_ip_count", raw, nil),
			"dst_port_count": utils.PathSearch("dst_port_count", raw, nil),
			"end_time":       utils.PathSearch("end_time", raw, nil),
			"records": flattenFlowLogStatisticDetailRecords(
				utils.PathSearch("records", raw, make([]interface{}, 0)).([]interface{})),
			"request_byte":  utils.PathSearch("request_byte", raw, nil),
			"response_byte": utils.PathSearch("response_byte", raw, nil),
			"sessions":      utils.PathSearch("sessions", raw, nil),
			"src_ip_count":  utils.PathSearch("src_ip_count", raw, nil),
			"start_time":    utils.PathSearch("start_time", raw, nil),
		},
	}
}

func flattenFlowLogStatisticDetailRecords(records []interface{}) []interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"app":                         utils.PathSearch("app", record, nil),
			"bytes":                       utils.PathSearch("bytes", record, nil),
			"dst_associate_instance_type": utils.PathSearch("dst_associate_instance_type", record, nil),
			"dst_device_name":             utils.PathSearch("dst_device_name", record, nil),
			"dst_ip":                      utils.PathSearch("dst_ip", record, nil),
			"dst_port":                    utils.PathSearch("dst_port", record, nil),
			"dst_host":                    utils.PathSearch("dst_host", record, nil),
			"dst_region_id":               utils.PathSearch("dst_region_id", record, nil),
			"dst_region_name":             utils.PathSearch("dst_region_name", record, nil),
			"end_time":                    utils.PathSearch("end_time", record, nil),
			"protocol":                    utils.PathSearch("protocol", record, nil),
			"request_byte":                utils.PathSearch("request_byte", record, nil),
			"response_byte":               utils.PathSearch("response_byte", record, nil),
			"sessions":                    utils.PathSearch("sessions", record, nil),
			"src_associate_instance_type": utils.PathSearch("src_associate_instance_type", record, nil),
			"src_device_name":             utils.PathSearch("src_device_name", record, nil),
			"src_ip":                      utils.PathSearch("src_ip", record, nil),
			"src_region_id":               utils.PathSearch("src_region_id", record, nil),
			"src_region_name":             utils.PathSearch("src_region_name", record, nil),
			"start_time":                  utils.PathSearch("start_time", record, nil),
		})
	}

	return result
}

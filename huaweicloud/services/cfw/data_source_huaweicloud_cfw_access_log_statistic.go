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

// @API CFW GET /v1/{project_id}/cfw/logs/access-top
func DataSourceAccessLogStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAccessLogStatisticRead,

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
			"rule_id": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: accessLogStatisticDataSchema(),
				},
			},
		},
	}
}

func accessLogStatisticDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"deny_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"deny_top_one_acl_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"deny_top_one_acl_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hit_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"in2out_deny_dst_ip_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"in2out_deny_dst_port_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"in2out_deny_dst_region_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"in2out_deny_src_ip_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"out2in_deny_dst_ip_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"out2in_deny_dst_port_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"out2in_deny_src_ip_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"out2in_deny_src_port_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"out2in_deny_src_region_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
		"permit_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"permit_top_one_acl_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"permit_top_one_acl_name": {
			Type:     schema.TypeString,
			Computed: true,
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
					"deny_access_top_counts": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"permit_access_top_counts": {
						Type:     schema.TypeInt,
						Computed: true,
					},
					"total_access_top_counts": {
						Type:     schema.TypeInt,
						Computed: true,
					},
				},
			},
		},
		"top_deny_rule_list": {
			Type:     schema.TypeList,
			Computed: true,
			Elem:     accessTopMemberVOSchema(),
		},
	}
}

func accessTopMemberVOSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"count": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"item": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildAccessLogStatisticQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?fw_instance_id=%v&item=%v",
		d.Get("fw_instance_id"), d.Get("item"))

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
	if rawArray, ok := d.Get("rule_id").([]interface{}); ok {
		for _, v := range rawArray {
			res = fmt.Sprintf("%s&rule_id=%s", res, v.(string))
		}
	}

	return res
}

func dataSourceAccessLogStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/cfw/logs/access-top"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildAccessLogStatisticQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW access log statistic: %s", err)
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
		d.Set("data", flattenAccessTopVOData(utils.PathSearch("data", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAccessTopVOData(dataResp interface{}) []interface{} {
	if dataResp == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"deny_count":            utils.PathSearch("deny_count", dataResp, nil),
		"deny_top_one_acl_id":   utils.PathSearch("deny_top_one_acl_id", dataResp, nil),
		"deny_top_one_acl_name": utils.PathSearch("deny_top_one_acl_name", dataResp, nil),
		"hit_count":             utils.PathSearch("hit_count", dataResp, nil),
		"in2out_deny_dst_ip_list": flattenAccessTopMemberVOs(
			utils.PathSearch("in2out_deny_dst_ip_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"in2out_deny_dst_port_list": flattenAccessTopMemberVOs(
			utils.PathSearch("in2out_deny_dst_port_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"in2out_deny_dst_region_list": flattenAccessTopMemberVOs(
			utils.PathSearch("in2out_deny_dst_region_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"in2out_deny_src_ip_list": flattenAccessTopMemberVOs(
			utils.PathSearch("in2out_deny_src_ip_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"out2in_deny_dst_ip_list": flattenAccessTopMemberVOs(
			utils.PathSearch("out2in_deny_dst_ip_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"out2in_deny_dst_port_list": flattenAccessTopMemberVOs(
			utils.PathSearch("out2in_deny_dst_port_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"out2in_deny_src_ip_list": flattenAccessTopMemberVOs(
			utils.PathSearch("out2in_deny_src_ip_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"out2in_deny_src_port_list": flattenAccessTopMemberVOs(
			utils.PathSearch("out2in_deny_src_port_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"out2in_deny_src_region_list": flattenAccessTopMemberVOs(
			utils.PathSearch("out2in_deny_src_region_list", dataResp, make([]interface{}, 0)).([]interface{})),
		"permit_count":            utils.PathSearch("permit_count", dataResp, nil),
		"permit_top_one_acl_id":   utils.PathSearch("permit_top_one_acl_id", dataResp, nil),
		"permit_top_one_acl_name": utils.PathSearch("permit_top_one_acl_name", dataResp, nil),
		"records": flattenAccessTopStatisticsRecords(
			utils.PathSearch("records", dataResp, make([]interface{}, 0)).([]interface{})),
		"top_deny_rule_list": flattenAccessTopMemberVOs(
			utils.PathSearch("top_deny_rule_list", dataResp, make([]interface{}, 0)).([]interface{})),
	}}
}

func flattenAccessTopMemberVOs(rawResp []interface{}) []interface{} {
	if len(rawResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawResp))
	for _, raw := range rawResp {
		result = append(result, map[string]interface{}{
			"count": utils.PathSearch("count", raw, nil),
			"item":  utils.PathSearch("item", raw, nil),
			"name":  utils.PathSearch("name", raw, nil),
		})
	}

	return result
}

func flattenAccessTopStatisticsRecords(records []interface{}) []interface{} {
	if len(records) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(records))
	for _, record := range records {
		result = append(result, map[string]interface{}{
			"agg_time":                 utils.PathSearch("agg_time", record, nil),
			"deny_access_top_counts":   utils.PathSearch("deny_access_top_counts", record, nil),
			"permit_access_top_counts": utils.PathSearch("permit_access_top_counts", record, nil),
			"total_access_top_counts":  utils.PathSearch("total_access_top_counts", record, nil),
		})
	}

	return result
}

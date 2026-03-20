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

// @API CFW GET /v1/{project_id}/report/{report_id}
func DataSourceReport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceReportRead,

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
			"report_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"report_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attack_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dst_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoDstIpElem(),
						},
						"ips_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoLevelElem(),
						},
						"rule": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoRuleElem(),
						},
						"src_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoSrcIpElem(),
						},
						"trend": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoTrendElem(),
						},
						"type": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataAttackInfoTypeElem(),
						},
					},
				},
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"internet_firewall": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataInternetFirewallEipElem(),
						},
						"in2out": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataInternetFirewallIn2outElem(),
						},
						"out2in": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataInternetFirewallOut2inElem(),
						},
						"overview": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataInternetFirewallOverviewElem(),
						},
						"traffic_trend": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataIntFirTraTrendElem(),
						},
					},
				},
			},
			"send_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"statistic_period": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"end_time": {
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
			"vpc_firewall": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallAppElem(),
						},
						"dst_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallDstIpElem(),
						},
						"overview": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallOverviewElem(),
						},
						"src_ip": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallSrcIpElem(),
						},
						"traffic_trend": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallTrafficTrendElem(),
						},
						"vpc": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     dataVpcFirewallVpcElem(),
						},
					},
				},
			},
		},
	}
}

func dataVpcFirewallOverviewElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"assets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"attack_event": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deny": {
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
			"traffic_peak": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"out_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"permit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agg_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"deny": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataAttackInfoDstIpElem() *schema.Resource {
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

func dataInternetFirewallOut2inElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dst_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"dst_port": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"src_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func dataVpcFirewallDstIpElem() *schema.Resource {
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

func dataVpcFirewallTrafficTrendElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"in_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"out_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"permit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agg_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"deny": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataAttackInfoTrendElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"in_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"out_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"permit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agg_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"deny": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataAttackInfoLevelElem() *schema.Resource {
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

func dataInternetFirewallIn2outElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"dst_port": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"src_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"dst_host": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
			"dst_ip": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func dataIntFirTraTrendElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"in_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"out_bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"permit": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agg_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bps": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"deny": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataVpcFirewallAppElem() *schema.Resource {
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

func dataVpcFirewallSrcIpElem() *schema.Resource {
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

func dataVpcFirewallVpcElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protected": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataAttackInfoSrcIpElem() *schema.Resource {
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

func dataAttackInfoTypeElem() *schema.Resource {
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

func dataAttackInfoRuleElem() *schema.Resource {
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

func dataInternetFirewallOverviewElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"nat": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"assets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"attack_event": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"deny": {
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
			"traffic_peak": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"out_bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"permit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"agg_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bps": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"deny": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataInternetFirewallEipElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protected": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"changed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"total": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildQueryReportQueryParams(d *schema.ResourceData) string {
	var (
		fwInstanceId    = d.Get("fw_instance_id").(string)
		reportProfileId = d.Get("report_profile_id").(string)
	)

	return fmt.Sprintf("?fw_instance_id=%s&report_profile_id=%s", fwInstanceId, reportProfileId)
}

func dataSourceReportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cfw"
		httpUrl = "v1/{project_id}/report/{report_id}"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{report_id}", d.Get("report_id").(string))
	requestPath += buildQueryReportQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW report: %s", err)
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
		d.Set("attack_info", flattenAttackInfo(utils.PathSearch("data.attack_info", respBody, nil))),
		d.Set("category", utils.PathSearch("data.category", respBody, nil)),
		d.Set("internet_firewall", flattenInternetFirewall(utils.PathSearch("data.internet_firewall", respBody, nil))),
		d.Set("send_time", utils.PathSearch("data.send_time", respBody, nil)),
		d.Set("statistic_period", flattenStatisticPeriod(utils.PathSearch("data.statistic_period", respBody, nil))),
		d.Set("vpc_firewall", flattenVpcFirewall(utils.PathSearch("data.vpc_firewall", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenVpcFirewall(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"app":           flattenItemArray(utils.PathSearch("app", respBody, make([]interface{}, 0)).([]interface{})),
			"dst_ip":        flattenItemArray(utils.PathSearch("dst_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"overview":      flattenOverviewObject(utils.PathSearch("overview", respBody, nil)),
			"src_ip":        flattenItemArray(utils.PathSearch("src_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"traffic_trend": flattenTrendArray(utils.PathSearch("traffic_trend", respBody, make([]interface{}, 0)).([]interface{})),
			"vpc":           flattenVpcFirewallVpc(utils.PathSearch("vpc", respBody, nil)),
		},
	}
}

func flattenVpcFirewallVpc(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"protected": flattenVpcFirewallVpcProtected(utils.PathSearch("protected", respBody, nil)),
			"total":     utils.PathSearch("total", respBody, nil),
		},
	}
}

func flattenVpcFirewallVpcProtected(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"changed": utils.PathSearch("changed", respBody, nil),
			"total":   utils.PathSearch("total", respBody, nil),
			"value":   utils.PathSearch("value", respBody, nil),
		},
	}
}

func flattenStatisticPeriod(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"end_time":   utils.PathSearch("end_time", respBody, nil),
			"start_time": utils.PathSearch("start_time", respBody, nil),
		},
	}
}

func flattenInternetFirewall(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"eip":           flattenEipObject(utils.PathSearch("eip", respBody, nil)),
			"in2out":        flattenInternetFirewallIn2out(utils.PathSearch("in2out", respBody, nil)),
			"out2in":        flattenInternetFirewallOut2in(utils.PathSearch("out2in", respBody, nil)),
			"overview":      flattenOverviewObject(utils.PathSearch("overview", respBody, nil)),
			"traffic_trend": flattenTrendArray(utils.PathSearch("traffic_trend", respBody, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenOverviewObject(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"access_policies": flattenOverviewAccessPolicies(utils.PathSearch("access_policies", respBody, nil)),
			"assets":          flattenOverviewAssets(utils.PathSearch("assets", respBody, nil)),
			"attack_event":    flattenOverviewAttackEvent(utils.PathSearch("attack_event", respBody, nil)),
			"traffic_peak":    flattenTrendObject(utils.PathSearch("traffic_peak", respBody, nil)),
		},
	}
}

func flattenTrendObject(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"agg_time": utils.PathSearch("agg_time", respBody, nil),
			"bps":      utils.PathSearch("bps", respBody, nil),
			"deny":     utils.PathSearch("deny", respBody, nil),
			"in_bps":   utils.PathSearch("in_bps", respBody, nil),
			"out_bps":  utils.PathSearch("out_bps", respBody, nil),
			"permit":   utils.PathSearch("permit", respBody, nil),
		},
	}
}

func flattenOverviewAttackEvent(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"changed": utils.PathSearch("changed", respBody, nil),
			"deny":    utils.PathSearch("deny", respBody, nil),
			"total":   utils.PathSearch("total", respBody, nil),
		},
	}
}

func flattenOverviewAssets(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"changed": utils.PathSearch("changed", respBody, nil),
			"total":   utils.PathSearch("total", respBody, nil),
			"value":   utils.PathSearch("value", respBody, nil),
		},
	}
}

func flattenOverviewAccessPolicies(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"changed": utils.PathSearch("changed", respBody, nil),
			"eip":     utils.PathSearch("eip", respBody, nil),
			"nat":     utils.PathSearch("nat", respBody, nil),
			"total":   utils.PathSearch("total", respBody, nil),
		},
	}
}

func flattenInternetFirewallOut2in(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"dst_ip":   flattenItemArray(utils.PathSearch("dst_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"dst_port": flattenItemArray(utils.PathSearch("dst_port", respBody, make([]interface{}, 0)).([]interface{})),
			"src_ip":   flattenItemArray(utils.PathSearch("src_ip", respBody, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenInternetFirewallIn2out(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"dst_host": flattenItemArray(utils.PathSearch("dst_host", respBody, make([]interface{}, 0)).([]interface{})),
			"dst_ip":   flattenItemArray(utils.PathSearch("dst_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"dst_port": flattenItemArray(utils.PathSearch("dst_port", respBody, make([]interface{}, 0)).([]interface{})),
			"src_ip":   flattenItemArray(utils.PathSearch("src_ip", respBody, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenEipObject(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"protected": flattenInternetFirewallEipProtected(utils.PathSearch("protected", respBody, nil)),
			"total":     utils.PathSearch("protectotalted", respBody, nil),
		},
	}
}

func flattenInternetFirewallEipProtected(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"changed": utils.PathSearch("changed", respBody, nil),
			"total":   utils.PathSearch("total", respBody, nil),
			"value":   utils.PathSearch("value", respBody, nil),
		},
	}
}

func flattenAttackInfo(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"dst_ip":   flattenItemArray(utils.PathSearch("dst_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"ips_mode": utils.PathSearch("ips_mode", respBody, nil),
			"level":    flattenItemArray(utils.PathSearch("level", respBody, make([]interface{}, 0)).([]interface{})),
			"rule":     flattenItemArray(utils.PathSearch("rule", respBody, make([]interface{}, 0)).([]interface{})),
			"src_ip":   flattenItemArray(utils.PathSearch("src_ip", respBody, make([]interface{}, 0)).([]interface{})),
			"trend":    flattenTrendArray(utils.PathSearch("trend", respBody, make([]interface{}, 0)).([]interface{})),
			"type":     flattenItemArray(utils.PathSearch("type", respBody, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenTrendArray(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, len(respArray))
	for i, item := range respArray {
		result[i] = map[string]interface{}{
			"agg_time": utils.PathSearch("agg_time", item, nil),
			"bps":      utils.PathSearch("bps", item, nil),
			"deny":     utils.PathSearch("deny", item, nil),
			"in_bps":   utils.PathSearch("in_bps", item, nil),
			"out_bps":  utils.PathSearch("out_bps", item, nil),
			"permit":   utils.PathSearch("permit", item, nil),
		}
	}
	return result
}

func flattenItemArray(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, len(respArray))
	for i, item := range respArray {
		result[i] = map[string]interface{}{
			"key":   utils.PathSearch("key", item, nil),
			"name":  utils.PathSearch("name", item, nil),
			"value": utils.PathSearch("value", item, nil),
		}
	}
	return result
}

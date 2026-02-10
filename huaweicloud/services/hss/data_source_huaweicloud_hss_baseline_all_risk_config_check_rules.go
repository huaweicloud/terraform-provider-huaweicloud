package hss

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

// @API HSS GET /v5/{project_id}/baseline/risk-config/check-rules
func DataSourceBaselineAllRiskConfigCheckRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineAllRiskConfigCheckRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"statistics_scan_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"check_rule_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"statistics_flag": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"total_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pass_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"processed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataSourceBaselineAllRiskConfigCheckRulesSchema(),
			},
		},
	}
}

func dataSourceBaselineAllRiskConfigCheckRulesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"tag": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_rule_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_type_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"failed_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scan_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"statistics_scan_result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_fix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_click": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"cancel_ignore_enable_click": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rule_params": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_param_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rule_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_value": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"range_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"range_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBaselineAllRiskConfigCheckRulesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("check_type"); ok {
		queryParams = fmt.Sprintf("%s&check_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("standard"); ok {
		queryParams = fmt.Sprintf("%s&standard=%v", queryParams, v)
	}
	if v, ok := d.GetOk("statistics_scan_result"); ok {
		queryParams = fmt.Sprintf("%s&statistics_scan_result=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_rule_name"); ok {
		queryParams = fmt.Sprintf("%s&check_rule_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("cluster_id"); ok {
		queryParams = fmt.Sprintf("%s&cluster_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("tag"); ok {
		queryParams = fmt.Sprintf("%s&tag=%v", queryParams, v)
	}
	if v, ok := d.GetOk("policy_group_id"); ok {
		queryParams = fmt.Sprintf("%s&policy_group_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("statistics_flag"); ok {
		queryParams = fmt.Sprintf("%s&statistics_flag=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceBaselineAllRiskConfigCheckRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "hss"
		epsId        = cfg.GetEnterpriseProjectID(d)
		result       = make([]interface{}, 0)
		offset       = 0
		httpUrl      = "v5/{project_id}/baseline/risk-config/check-rules"
		totalNum     float64
		passNum      float64
		failedNum    float64
		processedNum float64
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBaselineAllRiskConfigCheckRulesQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%v", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpts)
		if err != nil {
			return diag.Errorf("error retrieving HSS baseline all risk config check rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		totalNum = utils.PathSearch("total_num", respBody, float64(0)).(float64)
		passNum = utils.PathSearch("pass_num", respBody, float64(0)).(float64)
		failedNum = utils.PathSearch("failed_num", respBody, float64(0)).(float64)
		processedNum = utils.PathSearch("processed_num", respBody, float64(0)).(float64)
		dataListResp := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataListResp) == 0 {
			break
		}

		result = append(result, dataListResp...)
		offset += len(dataListResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("total_num", totalNum),
		d.Set("pass_num", passNum),
		d.Set("failed_num", failedNum),
		d.Set("processed_num", processedNum),
		d.Set("data_list", flattenBaselineAllRiskConfigCheckRulesDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineAllRiskConfigCheckRulesDataList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"tag":                        utils.PathSearch("tag", v, nil),
			"check_rule_name":            utils.PathSearch("check_rule_name", v, nil),
			"check_rule_id":              utils.PathSearch("check_rule_id", v, nil),
			"severity":                   utils.PathSearch("severity", v, nil),
			"check_type":                 utils.PathSearch("check_type", v, nil),
			"check_type_name":            utils.PathSearch("check_type_name", v, nil),
			"standard":                   utils.PathSearch("standard", v, nil),
			"host_num":                   utils.PathSearch("host_num", v, nil),
			"failed_num":                 utils.PathSearch("failed_num", v, nil),
			"scan_time":                  utils.PathSearch("scan_time", v, nil),
			"statistics_scan_result":     utils.PathSearch("statistics_scan_result", v, nil),
			"enable_fix":                 utils.PathSearch("enable_fix", v, nil),
			"enable_click":               utils.PathSearch("enable_click", v, nil),
			"cancel_ignore_enable_click": utils.PathSearch("cancel_ignore_enable_click", v, nil),
			"rule_params": flattenBaselineAllRiskConfigCheckRulesRuleParams(
				utils.PathSearch("rule_params", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenBaselineAllRiskConfigCheckRulesRuleParams(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"rule_param_id": utils.PathSearch("rule_param_id", v, nil),
			"rule_desc":     utils.PathSearch("rule_desc", v, nil),
			"default_value": utils.PathSearch("default_value", v, nil),
			"range_min":     utils.PathSearch("range_min", v, nil),
			"range_max":     utils.PathSearch("range_max", v, nil),
		})
	}

	return result
}

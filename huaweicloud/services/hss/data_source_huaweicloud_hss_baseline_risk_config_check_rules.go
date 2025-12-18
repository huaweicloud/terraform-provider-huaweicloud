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

// @API HSS GET /v5/{project_id}/baseline/risk-config/{check_name}/check-rules
func DataSourceBaselineRiskConfigCheckRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBaselineRiskConfigCheckRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"check_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standard": {
				Type:     schema.TypeString,
				Required: true,
			},
			"result_type": {
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
			"host_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"severity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"check_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standard": {
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
						"host_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"scan_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"not_enable_click_description": {
							Type:     schema.TypeString,
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
				},
			},
		},
	}
}

func buildBaselineRiskConfigCheckRulesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?standard=%v&limit=200", d.Get("standard"))

	if v, ok := d.GetOk("result_type"); ok {
		queryParams = fmt.Sprintf("%s&result_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("check_rule_name"); ok {
		queryParams = fmt.Sprintf("%s&check_rule_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("severity"); ok {
		queryParams = fmt.Sprintf("%s&severity=%v", queryParams, v)
	}
	if v, ok := d.GetOk("host_id"); ok {
		queryParams = fmt.Sprintf("%s&host_id=%v", queryParams, v)
	}
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceBaselineRiskConfigCheckRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v5/{project_id}/baseline/risk-config/{check_name}/check-rules"
		epsId     = cfg.GetEnterpriseProjectID(d)
		checkName = d.Get("check_name").(string)
		offset    = 0
		result    = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{check_name}", checkName)
	getPath += buildBaselineRiskConfigCheckRulesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving the checklist of a specified security configuration item: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenBaselineRiskConfigCheckRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBaselineRiskConfigCheckRules(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"severity":                     utils.PathSearch("severity", v, nil),
			"check_name":                   utils.PathSearch("check_name", v, nil),
			"check_type":                   utils.PathSearch("check_type", v, nil),
			"standard":                     utils.PathSearch("standard", v, nil),
			"check_rule_name":              utils.PathSearch("check_rule_name", v, nil),
			"check_rule_id":                utils.PathSearch("check_rule_id", v, nil),
			"host_num":                     utils.PathSearch("host_num", v, nil),
			"scan_result":                  utils.PathSearch("scan_result", v, nil),
			"status":                       utils.PathSearch("status", v, nil),
			"enable_fix":                   utils.PathSearch("enable_fix", v, nil),
			"enable_click":                 utils.PathSearch("enable_click", v, nil),
			"not_enable_click_description": utils.PathSearch("not_enable_click_description", v, nil),
			"rule_params": flattenBaselineRuleParams(
				utils.PathSearch("rule_params", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenBaselineRuleParams(dataResp []interface{}) []interface{} {
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

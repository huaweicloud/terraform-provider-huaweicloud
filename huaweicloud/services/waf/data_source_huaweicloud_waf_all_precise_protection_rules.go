package waf

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

// @API WAF GET /v1/{projectid}/waf/rule/custom
func DataSourceAllPreciseProtectionRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAllPreciseProtectionRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policyids": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"conditions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"index": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"logic_operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"contents": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"value_list_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"followed_action_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"start": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"terminal": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAllPreciseProtectionRulesQueryParams(d *schema.ResourceData, epsId string) string {
	res := "?pagesize=1000"
	if v, ok := d.GetOk("policyids"); ok {
		res = fmt.Sprintf("%s&policyids=%v", res, v)
	}
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	return res
}

func dataSourceAllPreciseProtectionRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		mErr        *multierror.Error
		httpUrl     = "v1/{projectid}/waf/rule/custom"
		allRules    []interface{}
		currentPage = 1
		epsId       = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{projectid}", client.ProjectID)
	requestPath += buildAllPreciseProtectionRulesQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithPage := fmt.Sprintf("%s&page=%d", requestPath, currentPage)
		resp, err := client.Request("GET", requestPathWithPage, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving WAF all precise protection rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		rulesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(rulesResp) == 0 {
			break
		}

		allRules = append(allRules, rulesResp...)
		currentPage++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("items", flattenAllPreciseProtectionRules(allRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAllPreciseProtectionRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"policyid":    utils.PathSearch("policyid", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"conditions": flattenPreciseProtectionRulesConditions(
				utils.PathSearch("conditions", v, make([]interface{}, 0)).([]interface{})),
			"action":    flattenPreciseProtectionRulesAction(utils.PathSearch("action", v, nil)),
			"priority":  utils.PathSearch("priority", v, nil),
			"timestamp": utils.PathSearch("timestamp", v, nil),
			"time":      utils.PathSearch("time", v, nil),
			"start":     utils.PathSearch("start", v, nil),
			"terminal":  utils.PathSearch("terminal", v, nil),
		})
	}
	return rst
}

func flattenPreciseProtectionRulesConditions(conditions []interface{}) []interface{} {
	if len(conditions) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(conditions))
	for _, v := range conditions {
		rst = append(rst, map[string]interface{}{
			"category":        utils.PathSearch("category", v, nil),
			"index":           utils.PathSearch("index", v, nil),
			"logic_operation": utils.PathSearch("logic_operation", v, nil),
			"contents":        utils.PathSearch("contents", v, nil),
			"value_list_id":   utils.PathSearch("value_list_id", v, nil),
		})
	}
	return rst
}

func flattenPreciseProtectionRulesAction(action interface{}) []interface{} {
	if action == nil {
		return nil
	}

	rst := map[string]interface{}{
		"category":           utils.PathSearch("category", action, nil),
		"followed_action_id": utils.PathSearch("followed_action_id", action, nil),
	}

	return []interface{}{rst}
}

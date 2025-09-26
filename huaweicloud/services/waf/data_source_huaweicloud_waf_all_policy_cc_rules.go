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

// @API WAF GET /v1/{projectid}/waf/rule/cc
func DataSourceAllPolicyCcRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAllPolicyCcRulesRead,

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
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"prefix": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"mode": {
							Type:     schema.TypeInt,
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
									"index": {
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
									"detail": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"response": {
													Type:     schema.TypeList,
													Computed: true,
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"content_type": {
																Type:     schema.TypeString,
																Computed: true,
															},
															"content": {
																Type:     schema.TypeString,
																Computed: true,
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"tag_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_index": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tag_condition": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"contents": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"limit_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"limit_period": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"unlock_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lock_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"domain_aggregation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"region_aggregation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timestamp": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildAllPolicyCcRulesQueryParams(d *schema.ResourceData, epsId string) string {
	res := "?pagesize=1000"
	if v, ok := d.GetOk("policyids"); ok {
		res = fmt.Sprintf("%s&policyids=%v", res, v)
	}
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	return res
}

func dataSourceAllPolicyCcRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "waf"
		httpUrl     = "v1/{projectid}/waf/rule/cc"
		epsId       = cfg.GetEnterpriseProjectID(d)
		result      = make([]interface{}, 0)
		currentPage = 1
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{projectid}", client.ProjectID)
	requestPath += buildAllPolicyCcRulesQueryParams(d, epsId)
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
			return diag.Errorf("error retrieving WAF all policy CC rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		rulesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(rulesResp) == 0 {
			break
		}

		result = append(result, rulesResp...)
		currentPage++
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenCcRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCcRules(rulesResp []interface{}) []interface{} {
	if len(rulesResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rulesResp))
	for _, rule := range rulesResp {
		result = append(result, map[string]interface{}{
			"name":               utils.PathSearch("name", rule, nil),
			"id":                 utils.PathSearch("id", rule, nil),
			"policyid":           utils.PathSearch("policyid", rule, nil),
			"url":                utils.PathSearch("url", rule, nil),
			"prefix":             utils.PathSearch("prefix", rule, nil),
			"mode":               utils.PathSearch("mode", rule, nil),
			"status":             utils.PathSearch("status", rule, nil),
			"conditions":         flattenConditions(utils.PathSearch("conditions", rule, make([]interface{}, 0)).([]interface{})),
			"action":             flattenRuleAction(utils.PathSearch("action", rule, nil)),
			"tag_type":           utils.PathSearch("tag_type", rule, nil),
			"tag_index":          utils.PathSearch("tag_index", rule, nil),
			"tag_condition":      flattenTagCondition(utils.PathSearch("tag_condition", rule, nil)),
			"limit_num":          utils.PathSearch("limit_num", rule, nil),
			"limit_period":       utils.PathSearch("limit_period", rule, nil),
			"unlock_num":         utils.PathSearch("unlock_num", rule, nil),
			"lock_time":          utils.PathSearch("lock_time", rule, nil),
			"domain_aggregation": utils.PathSearch("domain_aggregation", rule, nil),
			"region_aggregation": utils.PathSearch("region_aggregation", rule, nil),
			"description":        utils.PathSearch("description", rule, nil),
			"timestamp":          utils.PathSearch("timestamp", rule, nil),
		})
	}

	return result
}

func flattenConditions(conditionsResp []interface{}) []interface{} {
	if len(conditionsResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(conditionsResp))
	for _, condition := range conditionsResp {
		result = append(result, map[string]interface{}{
			"category":        utils.PathSearch("category", condition, nil),
			"logic_operation": utils.PathSearch("logic_operation", condition, nil),
			"contents":        utils.ExpandToStringList(utils.PathSearch("contents", condition, make([]interface{}, 0)).([]interface{})),
			"value_list_id":   utils.PathSearch("value_list_id", condition, nil),
			"index":           utils.PathSearch("index", condition, nil),
		})
	}

	return result
}

func flattenRuleAction(actionResp interface{}) []interface{} {
	if actionResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", actionResp, nil),
			"detail":   flattenActionDetail(utils.PathSearch("detail", actionResp, nil)),
		},
	}
}

func flattenActionDetail(detailResp interface{}) []interface{} {
	if detailResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"response": flattenActionResponse(utils.PathSearch("response", detailResp, nil)),
		},
	}
}

func flattenActionResponse(responseResp interface{}) []interface{} {
	if responseResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"content_type": utils.PathSearch("content_type", responseResp, nil),
			"content":      utils.PathSearch("content", responseResp, nil),
		},
	}
}

func flattenTagCondition(tagConditionResp interface{}) []interface{} {
	if tagConditionResp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", tagConditionResp, nil),
			"contents": utils.ExpandToStringList(utils.PathSearch("contents", tagConditionResp, make([]interface{}, 0)).([]interface{})),
		},
	}
}

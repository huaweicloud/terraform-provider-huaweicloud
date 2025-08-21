package aad

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AAD GET /v2/aad/policies/waf/frequency-control-rule
func DataSourceFrequencyControlRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFrequencyControlRulesRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
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
						"producer": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_time": {
							Type:     schema.TypeString,
							Computed: true,
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
						"mode": {
							Type:     schema.TypeString,
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
									"contents": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"index": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"logic_operation": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"unlock_num": {
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
						"captcha_lock_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"grayscale_time": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceFrequencyControlRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/policies/waf/frequency-control-rule"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?domain_name=%v", d.Get("domain_name"))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD frequency control rules: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("items",
			flattenFrequencyControlRulesItems(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenFrequencyControlRulesItems(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"producer":      utils.PathSearch("producer", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"url":           utils.PathSearch("url", v, nil),
			"limit_num":     utils.PathSearch("limit_num", v, nil),
			"limit_period":  utils.PathSearch("limit_period", v, nil),
			"lock_time":     utils.PathSearch("lock_time", v, nil),
			"tag_type":      utils.PathSearch("tag_type", v, nil),
			"tag_index":     utils.PathSearch("tag_index", v, nil),
			"tag_condition": flattenFrequencyControlRulesTagConditions(utils.PathSearch("tag_condition", v, nil)),
			"action":        flattenFrequencyControlRulesAction(utils.PathSearch("action", v, nil)),
			"domain_name":   utils.PathSearch("rule_id", v, nil),
			"overseas_type": utils.PathSearch("status", v, nil),
		})
	}

	return rst
}

func flattenFrequencyControlRulesTagConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", resp, nil),
			"contents": utils.ExpandToStringList(
				utils.PathSearch("contents", resp, make([]interface{}, 0)).([]interface{})),
		},
	}
}

func flattenFrequencyControlRulesAction(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", resp, nil),
			"detail":   flattenRulesActionDetail(utils.PathSearch("detail", resp, nil)),
		},
	}
}

func flattenRulesActionDetail(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"response": flattenRulesActionDetailResponse(utils.PathSearch("response", resp, nil)),
		},
	}
}

func flattenRulesActionDetailResponse(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"content_type": utils.PathSearch("content_type", resp, nil),
			"content":      utils.PathSearch("content", resp, nil),
		},
	}
}

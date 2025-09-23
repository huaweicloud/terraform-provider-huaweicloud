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

// @API AAD GET /v2/aad/policies/waf/custom-rule
func DataSourceCustomRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomRulesRead,

		Schema: map[string]*schema.Schema{
			"domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"overseas_type": {
				Type:     schema.TypeInt,
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
						"name": {
							Type:     schema.TypeString,
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
						"priority": {
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
								},
							},
						},
						"domain_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"overseas_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildCustomRulesQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?domain_name=%v", d.Get("domain_name"))
	queryParams += fmt.Sprintf("&overseas_type=%v", d.Get("overseas_type"))

	return queryParams
}

func dataSourceCustomRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "aad"
		httpUrl = "v2/aad/policies/waf/custom-rule"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AAD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildCustomRulesQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving AAD custom rules: %s", err)
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
			flattenCustomRulesItems(utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCustomRulesItems(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"name":          utils.PathSearch("name", v, nil),
			"time":          utils.PathSearch("time", v, nil),
			"start":         utils.PathSearch("start", v, nil),
			"terminal":      utils.PathSearch("terminal", v, nil),
			"priority":      utils.PathSearch("priority", v, nil),
			"conditions":    flattenCustomRulesConditions(utils.PathSearch("conditions", v, make([]interface{}, 0)).([]interface{})),
			"action":        flattenCustomRulesAction(utils.PathSearch("action", v, nil)),
			"domain_name":   utils.PathSearch("rule_id", v, nil),
			"overseas_type": utils.PathSearch("status", v, nil),
		})
	}

	return rst
}

func flattenCustomRulesConditions(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	var rst []interface{}
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"category":        utils.PathSearch("category", v, nil),
			"index":           utils.PathSearch("index", v, nil),
			"logic_operation": utils.PathSearch("logic_operation", v, nil),
			"contents":        utils.ExpandToStringList(utils.PathSearch("contents", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenCustomRulesAction(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", resp, nil),
		},
	}
}

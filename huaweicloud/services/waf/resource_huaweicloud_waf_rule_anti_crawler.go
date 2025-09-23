// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product WAF
// ---------------------------------------------------------------

package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/anticrawler
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/anticrawler
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}
func ResourceRuleAntiCrawler() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAntiCrawlerRuleCreate,
		UpdateContext: resourceAntiCrawlerRuleUpdate,
		ReadContext:   resourceAntiCrawlerRuleRead,
		DeleteContext: resourceAntiCrawlerRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the policy ID.`,
			},
			"protection_mode": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protection mode of WAF anti crawler rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the rule name.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the priority.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Elem:        ruleAntiCrawlerConditionsSchema(),
				Required:    true,
				Description: `Specifies the match condition list.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the rule description.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The rule status.`,
			},
		},
	}
}

func ruleAntiCrawlerConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"field": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the field type.`,
			},
			"logic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the logic for matching the condition.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the content of the condition.`,
			},
			"reference_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the reference table ID.`,
			},
		},
	}
}

func updateAntiCrawlerProtectionMode(client *golangsdk.ServiceClient, httpPath, protectionMode string) error {
	updateOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"anticrawler_type": protectionMode,
		},
	}

	_, err := client.Request("PUT", httpPath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating WAF anti crawler rule protection mode: %s", err)
	}
	return nil
}

func createAntiCrawlerRule(client *golangsdk.ServiceClient, httpPath string, d *schema.ResourceData) error {
	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         buildCreateOrUpdateRuleBodyParams(d),
	}

	createResp, err := client.Request("POST", httpPath, &createOpt)
	if err != nil {
		return fmt.Errorf("error creating WAF anti crawler rule: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return err
	}

	ruleId := utils.PathSearch("id", createRespBody, "").(string)
	if ruleId == "" {
		return fmt.Errorf("error creating WAF anti crawler rule: ID is not found in API response")
	}
	d.SetId(ruleId)
	return nil
}

func resourceAntiCrawlerRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/anticrawler"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	httpPath := client.Endpoint + httpUrl
	httpPath = strings.ReplaceAll(httpPath, "{project_id}", client.ProjectID)
	httpPath = strings.ReplaceAll(httpPath, "{policy_id}", d.Get("policy_id").(string))
	httpPath += buildQueryParams(d, cfg)

	// update WAF anti crawler protection mode first
	if err := updateAntiCrawlerProtectionMode(client, httpPath, d.Get("protection_mode").(string)); err != nil {
		return diag.FromErr(err)
	}

	if err := createAntiCrawlerRule(client, httpPath, d); err != nil {
		return diag.FromErr(err)
	}

	return resourceAntiCrawlerRuleRead(ctx, d, meta)
}

func buildCreateOrUpdateRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"type":        d.Get("protection_mode"),
		"priority":    d.Get("priority"),
		"description": d.Get("description"),
		"conditions":  buildRuleAntiCrawlerConditions(d.Get("conditions")),
	}
	return bodyParams
}

func buildRuleAntiCrawlerConditions(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"category":        raw["field"],
				"logic_operation": raw["logic"],
			}
			if content := raw["content"].(string); content != "" {
				rst[i]["contents"] = []string{content}
			}
			if referenceTableId := raw["reference_table_id"].(string); referenceTableId != "" {
				rst[i]["value_list_id"] = referenceTableId
			}
		}
		return rst
	}
	return nil
}

func resourceAntiCrawlerRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", d.Get("policy_id").(string))
	getPath = strings.ReplaceAll(getPath, "{rule_id}", d.Id())
	getPath += buildQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		// If the anti crawler rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF anti crawler rule")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("policy_id", utils.PathSearch("policyid", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("protection_mode", utils.PathSearch("type", getRespBody, nil)),
		d.Set("priority", utils.PathSearch("priority", getRespBody, nil)),
		d.Set("conditions", flattenRuleAntiCrawlerConditions(getRespBody)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRuleAntiCrawlerConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"field":              utils.PathSearch("category", v, nil),
			"logic":              utils.PathSearch("logic_operation", v, nil),
			"content":            utils.PathSearch("contents|[0]", v, nil),
			"reference_table_id": utils.PathSearch("value_list_id", v, nil),
		}
	}
	return rst
}

func resourceAntiCrawlerRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Get("policy_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
	updatePath += buildQueryParams(d, cfg)
	updateOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         buildCreateOrUpdateRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating WAF anti crawler rule: %s", err)
	}

	return resourceAntiCrawlerRuleRead(ctx, d, meta)
}

func resourceAntiCrawlerRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}/anticrawler/{rule_id}"
		product = "waf"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Get("policy_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	deletePath += buildQueryParams(d, cfg)
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the anti crawler rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF anti crawler rule")
	}
	return nil
}

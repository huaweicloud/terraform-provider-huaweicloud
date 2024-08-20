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
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}
// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/cc
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
func ResourceRuleCCProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleCCProtectionCreate,
		UpdateContext: resourceRuleCCProtectionUpdate,
		ReadContext:   resourceRuleCCProtectionRead,
		DeleteContext: resourceRuleCCProtectionDelete,
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
				Description: `Specifies the policy ID of WAF cc protection rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the rule name of WAF cc protection rule.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Elem:        ruleCCProtectionConditionsSchema(),
				Required:    true,
				Description: `Specifies the match condition list.`,
			},
			"protective_action": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the protective action taken when the number of requests reaches the upper limit.
The value can be **captcha**, **block**, **log** or **dynamic_block**.`,
			},
			"rate_limit_mode": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the rate limit mode.
Valid values are **ip**, **cookie**, **header**, **other**, **policy**, **domain**, **url**.`,
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Required: true,
				Description: `Specifies the number of requests allowed from a web visitor in a rate limiting period.
The value ranges from 1 to 2,147,483,647.`,
			},
			"limit_period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the rate limiting period. The value ranges from 1 to 3,600 in seconds.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID of WAF cc protection rule.`,
			},
			"block_page_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the type of the returned page.
The options are **application/json**, **text/html** and **text/xml**.`,
			},
			"page_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the content of the returned page.`,
			},
			"user_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user identifier.`,
			},
			"other_user_identifier": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the other user identifier.`,
			},
			"unlock_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the allowable frequency. The value ranges from 0 to 2,147,483,647.`,
			},
			"lock_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: `Specifies the lock time for resuming normal page access after blocking can be set.
The value ranges from 0 to 65,535 in seconds.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of WAF cc protection rule.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: `The status of a cc protection rule.`,
			},
			"request_aggregation": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable domain aggregation statistics. Default to false.`,
			},
			"all_waf_instances": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to enable global counting. Default to false.`,
			},
		},
	}
}

func ruleCCProtectionConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the field type.
The value can be **url**, **ip**, **ipv6**, **params**, **cookie**, **header** or **response_code**.`,
			},
			"logic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the condition matching logic.`,
			},
			"subfield": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the subfield of the condition.`,
			},
			"reference_table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the reference table id.`,
			},
			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the content of the match condition.`,
			},
		},
	}
	return &sc
}

func resourceRuleCCProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createRuleCCProtectionHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/cc"
		createRuleCCProtectionProduct = "waf"
	)
	createClient, err := cfg.NewServiceClient(createRuleCCProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	createPath := createClient.Endpoint + createRuleCCProtectionHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", createClient.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{policy_id}", fmt.Sprintf("%v", d.Get("policy_id")))
	createPath += buildQueryParams(d, cfg)

	createRuleCCProtectionOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createRuleCCProtectionOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateRuleCCProtectionBodyParams(d))
	createRuleCCProtectionResp, err := createClient.Request("POST", createPath, &createRuleCCProtectionOpt)
	if err != nil {
		return diag.Errorf("error creating RuleCCProtection: %s", err)
	}

	createRuleCCProtectionRespBody, err := utils.FlattenResponse(createRuleCCProtectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("id", createRuleCCProtectionRespBody)
	if err != nil {
		return diag.Errorf("error creating RuleCCProtection: ID is not found in API response")
	}
	d.SetId(id.(string))

	if d.Get("status").(int) == 0 {
		if err := updateRuleStatus(createClient, d, cfg, "cc"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRuleCCProtectionRead(ctx, d, meta)
}

func buildCreateOrUpdateRuleCCProtectionBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":               utils.ValueIgnoreEmpty(d.Get("name")),
		"tag_type":           utils.ValueIgnoreEmpty(d.Get("rate_limit_mode")),
		"tag_index":          utils.ValueIgnoreEmpty(d.Get("user_identifier")),
		"limit_num":          utils.ValueIgnoreEmpty(d.Get("limit_num")),
		"limit_period":       utils.ValueIgnoreEmpty(d.Get("limit_period")),
		"unlock_num":         utils.ValueIgnoreEmpty(d.Get("unlock_num")),
		"lock_time":          utils.ValueIgnoreEmpty(d.Get("lock_time")),
		"domain_aggregation": utils.ValueIgnoreEmpty(d.Get("request_aggregation")),
		"region_aggregation": utils.ValueIgnoreEmpty(d.Get("all_waf_instances")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"conditions":         buildRuleCCProtectionConditions(d.Get("conditions")),
		"action":             buildRuleCCProtectionAction(d),
		"tag_condition":      buildRuleCCProtectionTagCondition(d),
		"mode":               1,
	}
	return bodyParams
}

func buildRuleCCProtectionConditions(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"category":        utils.ValueIgnoreEmpty(raw["field"]),
				"logic_operation": utils.ValueIgnoreEmpty(raw["logic"]),
				"index":           utils.ValueIgnoreEmpty(raw["subfield"]),
				"contents":        buildCCProtectionContents(raw),
				"value_list_id":   utils.ValueIgnoreEmpty(raw["reference_table_id"]),
			}
		}
		return rst
	}
	return nil
}

func buildRuleCCProtectionAction(d *schema.ResourceData) map[string]interface{} {
	actionMap := map[string]interface{}{
		"category": d.Get("protective_action"),
	}
	if v, ok := d.GetOk("block_page_type"); ok {
		pageContent := d.Get("page_content")
		actionMap["detail"] = map[string]interface{}{
			"response": map[string]interface{}{
				"content_type": v,
				"content":      pageContent,
			},
		}
	}
	return actionMap
}

func buildRuleCCProtectionTagCondition(d *schema.ResourceData) map[string]interface{} {
	if v, ok := d.GetOk("other_user_identifier"); ok {
		return map[string]interface{}{
			"category": "referer",
			"contents": []string{v.(string)},
		}
	}
	return nil
}

func buildCCProtectionContents(raw map[string]interface{}) []string {
	var contents []string
	if content := utils.ValueIgnoreEmpty(raw["content"]); content != nil {
		contents = append(contents, content.(string))
	}
	return contents
}

func resourceRuleCCProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getRuleCCProtectionHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}"
		getRuleCCProtectionProduct = "waf"
	)
	getClient, err := cfg.NewServiceClient(getRuleCCProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	getPath := getClient.Endpoint + getRuleCCProtectionHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", getClient.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", fmt.Sprintf("%v", d.Get("policy_id")))
	getPath = strings.ReplaceAll(getPath, "{rule_id}", d.Id())
	getPath += buildQueryParams(d, cfg)

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getResp, err := getClient.Request("GET", getPath, &getOpt)

	if err != nil {
		// If the cc rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving RuleCCProtection")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("policy_id", utils.PathSearch("policyid", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("conditions", flattenRuleCCProtectionConditions(getRespBody)),
		d.Set("protective_action", utils.PathSearch("action.category", getRespBody, nil)),
		d.Set("block_page_type", utils.PathSearch("action.detail.response.content_type", getRespBody, nil)),
		d.Set("page_content", utils.PathSearch("action.detail.response.content", getRespBody, nil)),
		d.Set("rate_limit_mode", utils.PathSearch("tag_type", getRespBody, nil)),
		d.Set("user_identifier", utils.PathSearch("tag_index", getRespBody, nil)),
		d.Set("other_user_identifier", utils.PathSearch("tag_condition.contents|[0]", getRespBody, nil)),
		d.Set("limit_num", utils.PathSearch("limit_num", getRespBody, nil)),
		d.Set("limit_period", utils.PathSearch("limit_period", getRespBody, nil)),
		d.Set("unlock_num", utils.PathSearch("unlock_num", getRespBody, nil)),
		d.Set("lock_time", utils.PathSearch("lock_time", getRespBody, nil)),
		d.Set("request_aggregation", utils.PathSearch("domain_aggregation", getRespBody, nil)),
		d.Set("all_waf_instances", utils.PathSearch("region_aggregation", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRuleCCProtectionConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"field":              utils.PathSearch("category", v, nil),
			"subfield":           utils.PathSearch("index", v, nil),
			"logic":              utils.PathSearch("logic_operation", v, nil),
			"content":            utils.PathSearch("contents|[0]", v, nil),
			"reference_table_id": utils.PathSearch("value_list_id", v, nil),
		})
	}
	return rst
}

func resourceRuleCCProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateRuleCCProtectionClient, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	updateRuleCCProtectionChanges := []string{
		"name",
		"conditions",
		"protective_action",
		"block_page_type",
		"page_content",
		"rate_limit_mode",
		"user_identifier",
		"other_user_identifier",
		"limit_num",
		"limit_period",
		"unlock_num",
		"lock_time",
		"request_aggregation",
		"all_waf_instances",
		"description",
	}

	if d.HasChanges(updateRuleCCProtectionChanges...) {
		updateRuleCCProtectionHttpUrl := "v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}"

		policyID := fmt.Sprintf("%v", d.Get("policy_id"))
		updatePath := updateRuleCCProtectionClient.Endpoint + updateRuleCCProtectionHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", updateRuleCCProtectionClient.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{policy_id}", policyID)
		updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
		updatePath += buildQueryParams(d, cfg)

		updateRuleCCProtectionOpt := golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf8",
			},
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateRuleCCProtectionOpt.JSONBody = utils.RemoveNil(buildCreateOrUpdateRuleCCProtectionBodyParams(d))
		_, err := updateRuleCCProtectionClient.Request("PUT", updatePath, &updateRuleCCProtectionOpt)
		if err != nil {
			return diag.Errorf("error updating RuleCCProtection: %s", err)
		}
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(updateRuleCCProtectionClient, d, cfg, "cc"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRuleCCProtectionRead(ctx, d, meta)
}

func resourceRuleCCProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteRuleCCProtectionHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/cc/{rule_id}"
		deleteRuleCCProtectionProduct = "waf"
	)
	deleteRuleCCProtectionClient, err := cfg.NewServiceClient(deleteRuleCCProtectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating WAF Client: %s", err)
	}

	deletePath := deleteRuleCCProtectionClient.Endpoint + deleteRuleCCProtectionHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", deleteRuleCCProtectionClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", fmt.Sprintf("%v", d.Get("policy_id")))
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	deletePath += buildQueryParams(d, cfg)

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteRuleCCProtectionClient.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// If the cc rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting RuleCCProtection")
	}
	return nil
}

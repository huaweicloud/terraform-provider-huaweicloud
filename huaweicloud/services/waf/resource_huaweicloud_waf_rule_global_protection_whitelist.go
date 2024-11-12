// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product WAF
// ---------------------------------------------------------------

package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/ignore
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}
func ResourceRuleGlobalProtectionWhitelist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleGlobalProtectionWhitelistCreate,
		UpdateContext: resourceRuleGlobalProtectionWhitelistUpdate,
		ReadContext:   resourceRuleGlobalProtectionWhitelistRead,
		DeleteContext: resourceRuleGlobalProtectionWhitelistDelete,
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
				Description: `Specifies the policy ID of WAF global protection whitelist rule.`,
			},
			"domains": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				Description: `Specifies the protected domain name bound with the policy or manually enter a single
domain name corresponding to the wildcard domain name.`,
			},
			"ignore_waf_protection": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ignore waf protection rule.`,
			},
			"conditions": {
				Type:        schema.TypeList,
				Elem:        globalProtectionWhitelistConditionsSchema(),
				Required:    true,
				Description: `Specifies the match condition list.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID of WAF global protection whitelist rule.`,
			},
			"advanced_field": {
				Type:     schema.TypeString,
				Optional: true,
				Description: `Specifies the advanced field. The following fields are supported:
**params**, **cookie**, **header**, **body** and **multipart**.`,
			},
			"advanced_content": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the advanced content.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of WAF global protection whitelist rule.`,
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
				Description:  `Specifies the status of WAF global protection whitelist rule.`,
			},
		},
	}
}

func globalProtectionWhitelistConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"field": {
				Type:     schema.TypeString,
				Required: true,
				Description: `Specifies the field type. The value can be **ip**, **url**, **params**, **cookie** 
or **header**.`,
			},
			"logic": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the condition matching logic.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the content of the match condition.`,
			},
			"subfield": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the subfield of the condition.`,
			},
		},
	}
	return &sc
}

func resourceRuleGlobalProtectionWhitelistCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ignore"
		product  = "waf"
		policyID = d.Get("policy_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyID)
	requestPath += buildQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateGlobalProtectionWhitelistBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF rule global protection whitelist: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating WAF rule global protection whitelist: ID is not found in API response")
	}
	d.SetId(id)

	if d.Get("status").(int) == 0 {
		if err := updateRuleStatus(client, d, cfg, "ignore"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRuleGlobalProtectionWhitelistRead(ctx, d, meta)
}

func buildCreateOrUpdateGlobalProtectionWhitelistBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"rule":        utils.ValueIgnoreEmpty(d.Get("ignore_waf_protection")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"domain":      d.Get("domains"),
		"conditions":  buildGlobalProtectionWhitelistConditions(d.Get("conditions")),
		"mode":        1,
	}

	if v, ok := d.GetOk("advanced_field"); ok {
		advancedMap := map[string]interface{}{
			"index": v,
			// the WAF api parameter `contents` needs an empty array as a default value
			"contents": []string{},
		}
		if v1, ok1 := d.GetOk("advanced_content"); ok1 {
			advancedMap["contents"] = utils.ExpandToStringList([]interface{}{v1})
		}
		bodyParams["advanced"] = advancedMap
	}

	return bodyParams
}

func buildGlobalProtectionWhitelistConditions(rawParams interface{}) []map[string]interface{} {
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
				"contents":        buildGlobalProtectionWhitelistContents(raw),
				"index":           utils.ValueIgnoreEmpty(raw["subfield"]),
			}
		}
		return rst
	}
	return nil
}

func buildGlobalProtectionWhitelistContents(raw map[string]interface{}) []string {
	var contents []string
	if content := utils.ValueIgnoreEmpty(raw["content"]); content != nil {
		contents = append(contents, content.(string))
	}
	return contents
}

func resourceRuleGlobalProtectionWhitelistRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		mErr     *multierror.Error
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}"
		product  = "waf"
		policyID = d.Get("policy_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestPath += buildQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// If the rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF global protection whitelist rule")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("policy_id", utils.PathSearch("policyid", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("ignore_waf_protection", utils.PathSearch("rule", respBody, nil)),
		d.Set("conditions", flattenGlobalProtectionWhitelistConditions(respBody)),
		d.Set("advanced_field", utils.PathSearch("advanced.index", respBody, nil)),
		d.Set("advanced_content", utils.PathSearch("advanced.contents|[0]", respBody, nil)),
		d.Set("domains", utils.PathSearch("domain", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalProtectionWhitelistConditions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("conditions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"field":    utils.PathSearch("category", v, nil),
			"subfield": utils.PathSearch("index", v, nil),
			"logic":    utils.PathSearch("logic_operation", v, nil),
			"content":  utils.PathSearch("contents|[0]", v, nil),
		})
	}
	return rst
}

func resourceRuleGlobalProtectionWhitelistUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	updateRuleGlobalProtectionWhitelistChanges := []string{
		"policy_id",
		"domains",
		"ignore_waf_protection",
		"conditions",
		"advanced_field",
		"advanced_content",
		"description",
	}

	if d.HasChanges(updateRuleGlobalProtectionWhitelistChanges...) {
		var (
			httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}"
			policyID = d.Get("policy_id").(string)
		)

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyID)
		requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
		requestPath += buildQueryParams(d, cfg)
		requestOpt := golangsdk.RequestOpts{
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=utf8",
			},
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildCreateOrUpdateGlobalProtectionWhitelistBodyParams(d)),
		}

		_, err := client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating WAF global protection whitelist rule: %s", err)
		}
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(client, d, cfg, "ignore"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceRuleGlobalProtectionWhitelistRead(ctx, d, meta)
}

func resourceRuleGlobalProtectionWhitelistDelete(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ignore/{rule_id}"
		product  = "waf"
		policyID = d.Get("policy_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyID)
	requestPath = strings.ReplaceAll(requestPath, "{rule_id}", d.Id())
	requestPath += buildQueryParams(d, cfg)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		// If the rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF global protection whitelist rule")
	}
	return nil
}

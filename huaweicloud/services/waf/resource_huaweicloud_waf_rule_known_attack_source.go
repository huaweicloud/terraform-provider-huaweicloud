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

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/punishment
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}
func ResourceRuleKnownAttack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRuleKnownAttackCreate,
		UpdateContext: resourceRuleKnownAttackUpdate,
		ReadContext:   resourceRuleKnownAttackRead,
		DeleteContext: resourceRuleKnownAttackDelete,
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
			"block_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of WAF known attack source rule.`,
			},
			"block_time": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the blocking time in seconds.`,
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
				Description: `Specifies the description of WAF known attack source rule.`,
			},
		},
	}
}

func resourceRuleKnownAttackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/punishment"
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
		JSONBody:         buildRuleKnownAttackBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF known attack source rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating WAF known attack source rule: ID is not found in API response")
	}
	d.SetId(ruleId)

	return resourceRuleKnownAttackRead(ctx, d, meta)
}

func buildRuleKnownAttackBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"category":    d.Get("block_type"),
		"block_time":  d.Get("block_time"),
		"description": d.Get("description"),
	}
}

func resourceRuleKnownAttackRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		mErr     *multierror.Error
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}"
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
		// If the known attack source rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF known attack source rule")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("policy_id", utils.PathSearch("policyid", respBody, nil)),
		d.Set("block_time", utils.PathSearch("block_time", respBody, nil)),
		d.Set("block_type", utils.PathSearch("category", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRuleKnownAttackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}"
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
		JSONBody:         buildRuleKnownAttackBodyParams(d),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating WAF known attack source rule: %s", err)
	}

	return resourceRuleKnownAttackRead(ctx, d, meta)
}

func resourceRuleKnownAttackDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/punishment/{rule_id}"
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
		// If the known attack source rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF known attack source rule")
	}

	return nil
}

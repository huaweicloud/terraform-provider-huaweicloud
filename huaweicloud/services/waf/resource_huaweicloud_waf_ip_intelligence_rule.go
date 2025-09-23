package waf

import (
	"context"
	"fmt"
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

var nonUpdatableParamsIpIntelligenceRule = []string{"policy_id", "enterprise_project_id"}

// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/ip-reputation
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}
func ResourceIpIntelligenceRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpIntelligenceRuleCreate,
		ReadContext:   resourceIpIntelligenceRuleRead,
		UpdateContext: resourceIpIntelligenceRuleUpdate,
		DeleteContext: resourceIpIntelligenceRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIpIntelligenceRuleImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsIpIntelligenceRule),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			// Because the query API does not return the parameter, so `Computed` is not added.
			"policyname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"policyid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildCreateIpIntelligenceRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"tags":        utils.ExpandToStringList(d.Get("tags").([]interface{})),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"action":      utils.ValueIgnoreEmpty(buildRuleActionBodyParams(d.Get("action").([]interface{}))),
		"policyname":  utils.ValueIgnoreEmpty(d.Get("policyname")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func buildRuleActionBodyParams(rawAction []interface{}) map[string]interface{} {
	if len(rawAction) == 0 {
		return nil
	}

	ruleAction, ok := rawAction[0].(map[string]interface{})
	if !ok {
		return nil
	}

	actionParams := map[string]interface{}{
		"category": ruleAction["category"],
	}

	return actionParams
}

func resourceIpIntelligenceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		policyId      = d.Get("policy_id").(string)
		epsId         = cfg.GetEnterpriseProjectID(d)
		createHttpUrl = "v1/{project_id}/waf/policy/{policy_id}/ip-reputation"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{policy_id}", policyId)
	if epsId != "" {
		createPath = fmt.Sprintf("%s?enterprise_project_id=%s", createPath, epsId)
	}

	createOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateIpIntelligenceRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IP intelligence rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating IP intelligence rule: unable to find intelligence rule ID")
	}

	d.SetId(ruleId)

	return resourceIpIntelligenceRuleRead(ctx, d, meta)
}

func resourceIpIntelligenceRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Get("policy_id").(string)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	rule, err := GetIpIntelligenceRuleInfo(client, policyId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IP intelligence rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", rule, nil)),
		d.Set("tags", utils.PathSearch("tags", rule, nil)),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("action", flattenIntelligenceRuleAction(utils.PathSearch("action", rule, nil))),
		d.Set("description", utils.PathSearch("description", rule, nil)),
		d.Set("status", utils.PathSearch("status", rule, nil)),
		d.Set("policyid", utils.PathSearch("policyid", rule, nil)),
		d.Set("timestamp", utils.PathSearch("timestamp", rule, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetIpIntelligenceRuleInfo(client *golangsdk.ServiceClient, policyId, ruleId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", policyId)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", ruleId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenIntelligenceRuleAction(ruleAction interface{}) []interface{} {
	if ruleAction == nil {
		return nil
	}

	result := map[string]interface{}{
		"category": utils.PathSearch("category", ruleAction, nil),
	}

	return []interface{}{result}
}

func buildUpdateIpIntelligenceRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"type":        d.Get("type"),
		"tags":        utils.ExpandToStringList(d.Get("tags").([]interface{})),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"action":      buildRuleActionBodyParams(d.Get("action").([]interface{})),
		"policyname":  utils.ValueIgnoreEmpty(d.Get("policyname")),
		"description": d.Get("description"),
	}

	return bodyParams
}

func resourceIpIntelligenceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Get("policy_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", policyId)
	updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
	if epsId != "" {
		updatePath = fmt.Sprintf("%s?enterprise_project_id=%s", updatePath, epsId)
	}

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         buildUpdateIpIntelligenceRuleBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating IP intelligence rule: %s", err)
	}

	return resourceIpIntelligenceRuleRead(ctx, d, meta)
}

func resourceIpIntelligenceRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Get("policy_id").(string)
		epsId    = cfg.GetEnterpriseProjectID(d)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ip-reputation/{rule_id}"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())
	if epsId != "" {
		deletePath = fmt.Sprintf("%s?enterprise_project_id=%s", deletePath, epsId)
	}

	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IP intelligence rule, the error message")
	}

	return nil
}

func resourceIpIntelligenceRuleImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<policy_id>/<id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("policy_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

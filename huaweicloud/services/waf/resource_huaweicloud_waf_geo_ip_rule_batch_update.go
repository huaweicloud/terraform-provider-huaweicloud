package waf

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var geoIpRuleBatchUpdateNonUpdatableParams = []string{
	"policy_rule_ids",
	"policy_rule_ids.*.policy_id",
	"policy_rule_ids.*.rule_ids",
	"status",
	"name",
	"geoip",
	"white",
}

// @API WAF POST /v1/{project_id}/waf/rule/geoip/batch-update
func ResourceGeoIpRuleBatchUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeoIpRuleBatchUpdateCreate,
		ReadContext:   resourceGeoIpRuleBatchUpdateRead,
		UpdateContext: resourceGeoIpRuleBatchUpdateUpdate,
		DeleteContext: resourceGeoIpRuleBatchUpdateDelete,

		CustomizeDiff: config.FlexibleForceNew(geoIpRuleBatchUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level
				region will be used.`,
			},
			"geoip": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the geo location codes blocked by the rule, separated by '|'.`,
			},
			"policy_rule_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Description: `Specifies an array of policy and rule IDs, using to associate protection policies with
				corresponding rule sets`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the ID of the protection policy.`,
						},
						"rule_ids": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: `Specifies an array of rule IDs associated with the protection policy.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the status of the geo IP rule (0: disabled, 1: enabled).`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the geo IP rule.`,
			},
			"white": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the protection action (1: allow, 2: block).`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildGeoIpPolicyRuleIdsBodyParams(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	ruleParams := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		raw, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"policy_id": raw["policy_id"],
			"rule_ids":  utils.ExpandToStringList(raw["rule_ids"].([]interface{})),
		}
		ruleParams = append(ruleParams, params)
	}

	return ruleParams
}

func buildGeoIpRuleBatchUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"geoip":           d.Get("geoip"),
		"policy_rule_ids": buildGeoIpPolicyRuleIdsBodyParams(d.Get("policy_rule_ids").([]interface{})),
		"status":          d.Get("status"),
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"white":           utils.ValueIgnoreEmpty(d.Get("white")),
	}
}

func resourceGeoIpRuleBatchUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/geoip/batch-update"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildGeoIpRuleBatchUpdateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch updating WAF GEO IP rules: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceGeoIpRuleBatchUpdateRead(ctx, d, meta)
}

func resourceGeoIpRuleBatchUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceGeoIpRuleBatchUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a action resource.
	return nil
}

func resourceGeoIpRuleBatchUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to batch update WAF geo IP rules. 
Deleting this resource will not change the current geo IP rule configurations, but will only remove the resource 
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

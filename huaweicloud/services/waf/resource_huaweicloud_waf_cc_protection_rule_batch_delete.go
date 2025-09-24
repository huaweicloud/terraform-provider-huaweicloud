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

var ccRuleBatchDeleteNonUpdatableParams = []string{
	"policy_rule_ids",
	"policy_rule_ids.*.policy_id",
	"policy_rule_ids.*.rule_ids",
}

// @API WAF POST /v1/{project_id}/waf/rule/cc/batch-delete
func ResourceCcRuleBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCcRuleBatchDeleteCreate,
		ReadContext:   resourceCcRuleBatchDeleteRead,
		UpdateContext: resourceCcRuleBatchDeleteUpdate,
		DeleteContext: resourceCcRuleBatchDeleteDelete,

		CustomizeDiff: config.FlexibleForceNew(ccRuleBatchDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"policy_rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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

func buildCcRuleBatchDeleteBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_rule_ids": buildCcRuleBodyParams(d),
	}

	return bodyParams
}

func buildCcRuleBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("policy_rule_ids").([]interface{})
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

func resourceCcRuleBatchDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/cc/batch-delete"
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
		JSONBody: buildCcRuleBatchDeleteBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error deleting CC attack protection rules: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return nil
}

func resourceCcRuleBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceCcRuleBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a action resource.
	return nil
}

func resourceCcRuleBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}

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

// @API WAF POST /v1/{project_id}/waf/rule/whiteblackip/batch-update
func ResourceWafBatchUpdateWhiteblackipRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchUpdateWhiteblackipRulesCreate,
		ReadContext:   resourceWafBatchUpdateWhiteblackipRulesRead,
		UpdateContext: resourceWafBatchUpdateWhiteblackipRulesUpdate,
		DeleteContext: resourceWafBatchUpdateWhiteblackipRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"name",
			"white",
			"policy_rule_ids",
			"policy_rule_ids.*.policy_id",
			"policy_rule_ids.*.rule_ids",
			"addr",
			"description",
			"ip_group_id",
			"time_mode",
			"start",
			"terminal",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"white": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     buildBatchUpdateWhiteblackipRulesPolicyRuleIdsSchema(),
			},
			"addr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"terminal": {
				Type:     schema.TypeInt,
				Optional: true,
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

func buildBatchUpdateWhiteblackipRulesPolicyRuleIdsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func buildBatchUpdateWhiteblackipRulesPolicyRuleIds(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"policy_id": rawMap["policy_id"],
			"rule_ids":  rawMap["rule_ids"],
		})
	}

	return rst
}

func buildBatchUpdateWhiteblackipRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            d.Get("name"),
		"white":           d.Get("white"),
		"policy_rule_ids": buildBatchUpdateWhiteblackipRulesPolicyRuleIds(d.Get("policy_rule_ids").([]interface{})),
		"addr":            utils.ValueIgnoreEmpty(d.Get("addr")),
		"description":     utils.ValueIgnoreEmpty(d.Get("description")),
		"ip_group_id":     utils.ValueIgnoreEmpty(d.Get("ip_group_id")),
		"time_mode":       utils.ValueIgnoreEmpty(d.Get("time_mode")),
		"start":           utils.ValueIgnoreEmpty(d.Get("start")),
		"terminal":        utils.ValueIgnoreEmpty(d.Get("terminal")),
	}

	return bodyParams
}

func resourceWafBatchUpdateWhiteblackipRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/whiteblackip/batch-update"
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
		JSONBody: utils.RemoveNil(buildBatchUpdateWhiteblackipRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch updating WAF white-black IP rules: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(id)

	return resourceWafBatchUpdateWhiteblackipRulesRead(ctx, d, meta)
}

func resourceWafBatchUpdateWhiteblackipRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchUpdateWhiteblackipRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchUpdateWhiteblackipRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch update white-black IP rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

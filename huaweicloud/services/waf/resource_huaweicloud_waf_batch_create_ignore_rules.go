package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/rule/ignore
func ResourceWafBatchCreateIgnoreRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchCreateIgnoreRulesCreate,
		ReadContext:   resourceWafBatchCreateIgnoreRulesRead,
		UpdateContext: resourceWafBatchCreateIgnoreRulesUpdate,
		DeleteContext: resourceWafBatchCreateIgnoreRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"conditions",
			"conditions.*.category",
			"conditions.*.contents",
			"conditions.*.logic_operation",
			"conditions.*.check_all_indexes_logic",
			"conditions.*.index",
			"rule",
			"policy_ids",
			"domain",
			"advanced",
			"advanced.*.index",
			"advanced.*.contents",
			"description",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     batchCreateIgnoreRulesConditionsSchema(),
			},
			"rule": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"advanced": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     batchCreateIgnoreRulesAdvancedSchema(),
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
		},
	}
}

func batchCreateIgnoreRulesAdvancedSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"index": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contents": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func batchCreateIgnoreRulesConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contents": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"logic_operation": {
				Type:     schema.TypeString,
				Required: true,
			},
			"check_all_indexes_logic": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"index": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBatchCreateIgnoreRulesQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildBatchCreateIgnoreRulesConditionsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray, ok := d.Get("conditions").([]interface{})
	if !ok {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"category":                rawMap["category"],
			"contents":                rawMap["contents"],
			"logic_operation":         rawMap["logic_operation"],
			"check_all_indexes_logic": utils.ValueIgnoreEmpty(rawMap["check_all_indexes_logic"]),
			"index":                   utils.ValueIgnoreEmpty(rawMap["index"]),
		})
	}

	return rst
}

func buildBatchCreateIgnoreRulesAdvancedBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("advanced").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"index":    utils.ValueIgnoreEmpty(rawMap["index"]),
		"contents": utils.ValueIgnoreEmpty(rawMap["contents"]),
	}
}

func buildBatchCreateIgnoreRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":      d.Get("domain"),
		"conditions":  buildBatchCreateIgnoreRulesConditionsBodyParams(d),
		"mode":        1,
		"rule":        d.Get("rule"),
		"policy_ids":  d.Get("policy_ids"),
		"advanced":    buildBatchCreateIgnoreRulesAdvancedBodyParams(d),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}

	return bodyParams
}

func resourceWafBatchCreateIgnoreRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/ignore"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildBatchCreateIgnoreRulesQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildBatchCreateIgnoreRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch creating WAF ignore rules: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceWafBatchCreateIgnoreRulesRead(ctx, d, meta)
}

func resourceWafBatchCreateIgnoreRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateIgnoreRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateIgnoreRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch create WAF ignore rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

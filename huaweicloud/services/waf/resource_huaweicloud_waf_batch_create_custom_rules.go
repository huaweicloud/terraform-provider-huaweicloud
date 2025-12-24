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

// @API WAF POST /v1/{project_id}/waf/rule/custom
func ResourceWafBatchCreateCustomRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchCreateCustomRulesCreate,
		ReadContext:   resourceWafBatchCreateCustomRulesRead,
		UpdateContext: resourceWafBatchCreateCustomRulesUpdate,
		DeleteContext: resourceWafBatchCreateCustomRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"time",
			"conditions",
			"conditions.*.category",
			"conditions.*.index",
			"conditions.*.logic_operation",
			"conditions.*.contents",
			"conditions.*.value_list_id",
			"action",
			"action.*.category",
			"action.*.followed_action_id",
			"priority",
			"name",
			"policy_ids",
			"start",
			"terminal",
			"description",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"time": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     buildBatchCreateCustomRulesConditionsSchema(),
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     buildBatchCreateCustomRulesActionSchema(),
			},
			"priority": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
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
			"start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"terminal": {
				Type:     schema.TypeInt,
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
		},
	}
}

func buildBatchCreateCustomRulesActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"followed_action_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBatchCreateCustomRulesConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"index": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"logic_operation": {
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
			"value_list_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBatchCreateCustomRulesQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildBatchCreateCustomRulesConditionsBodyParam(d *schema.ResourceData) []map[string]interface{} {
	rawArray, ok := d.Get("conditions").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			return nil
		}

		rst = append(rst, map[string]interface{}{
			"category":        utils.ValueIgnoreEmpty(rawMap["category"]),
			"index":           utils.ValueIgnoreEmpty(rawMap["index"]),
			"logic_operation": utils.ValueIgnoreEmpty(rawMap["logic_operation"]),
			"contents":        utils.ValueIgnoreEmpty(rawMap["contents"]),
			"value_list_id":   utils.ValueIgnoreEmpty(rawMap["value_list_id"]),
		})
	}
	return rst
}

func buildBatchCreateCustomRulesActionBodyParam(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("action").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category":           rawMap["category"],
		"followed_action_id": utils.ValueIgnoreEmpty(rawMap["followed_action_id"]),
	}
}

func buildBatchCreateCustomRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"time":        d.Get("time"),
		"conditions":  buildBatchCreateCustomRulesConditionsBodyParam(d),
		"action":      buildBatchCreateCustomRulesActionBodyParam(d),
		"priority":    d.Get("priority"),
		"name":        d.Get("name"),
		"policy_ids":  d.Get("policy_ids"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"start":       utils.ValueIgnoreEmpty(d.Get("start")),
		"terminal":    utils.ValueIgnoreEmpty(d.Get("terminal")),
	}

	return bodyParams
}

func resourceWafBatchCreateCustomRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/custom"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildBatchCreateCustomRulesQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildBatchCreateCustomRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch creating WAF custom rules: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceWafBatchCreateCustomRulesRead(ctx, d, meta)
}

func resourceWafBatchCreateCustomRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateCustomRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateCustomRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch create WAF custom rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

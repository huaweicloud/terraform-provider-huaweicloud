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

// @API WAF POST /v1/{project_id}/waf/rule/cc
func ResourceWafBatchCreateCcRules() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafBatchCreateCcRulesCreate,
		ReadContext:   resourceWafBatchCreateCcRulesRead,
		UpdateContext: resourceWafBatchCreateCcRulesUpdate,
		DeleteContext: resourceWafBatchCreateCcRulesDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"name",
			"conditions",
			"conditions.*.category",
			"conditions.*.logic_operation",
			"conditions.*.contents",
			"conditions.*.value_list_id",
			"conditions.*.index",
			"action",
			"action.*.category",
			"action.*.detail",
			"action.*.detail.*.response",
			"action.*.detail.*.response.*.content_type",
			"action.*.detail.*.response.*.content",
			"tag_type",
			"limit_num",
			"limit_period",
			"policy_ids",
			"tag_index",
			"tag_condition",
			"tag_condition.*.category",
			"tag_condition.*.contents",
			"unlock_num",
			"lock_time",
			"domain_aggregation",
			"region_aggregation",
			"description",
			"enterprise_project_id",
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
			"conditions": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     buildBatchCreateCcRulesConditionsSchema(),
			},
			"action": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     buildBatchCreateCcRulesActionSchema(),
			},
			"tag_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"limit_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"limit_period": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"policy_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tag_index": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_condition": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     buildBatchCreateCcRulesTagConditionSchema(),
			},
			"unlock_num": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"lock_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"domain_aggregation": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"region_aggregation": {
				Type:     schema.TypeBool,
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

func buildBatchCreateCcRulesTagConditionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
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

func buildBatchCreateCcRulesActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"detail": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     buildBatchCreateCcRulesActionDetailSchema(),
			},
		},
	}
}

func buildBatchCreateCcRulesActionDetailSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"response": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     buildBatchCreateCcRulesActionDetailResponseSchema(),
			},
		},
	}
}

func buildBatchCreateCcRulesActionDetailResponseSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBatchCreateCcRulesConditionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"logic_operation": {
				Type:     schema.TypeString,
				Required: true,
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
			"index": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildBatchCreateCcRulesQueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}
	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func buildBatchCreateCcRulesConditionsBodyParam(d *schema.ResourceData) []map[string]interface{} {
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
			"category":        rawMap["category"],
			"logic_operation": rawMap["logic_operation"],
			"contents":        utils.ValueIgnoreEmpty(rawMap["contents"]),
			"value_list_id":   utils.ValueIgnoreEmpty(rawMap["value_list_id"]),
			"index":           utils.ValueIgnoreEmpty(rawMap["index"]),
		})
	}
	return rst
}

func buildBatchCreateCcRulesActionDetailResponseBodyParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"content_type": utils.ValueIgnoreEmpty(rawMap["content_type"]),
		"content":      utils.ValueIgnoreEmpty(rawMap["content"]),
	}
}

func buildBatchCreateCcRulesActionDetailBodyParam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"response": buildBatchCreateCcRulesActionDetailResponseBodyParam(rawMap["response"].([]interface{})),
	}
}

func buildBatchCreateCcRulesActionBodyParam(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("action").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category": rawMap["category"],
		"detail":   buildBatchCreateCcRulesActionDetailBodyParam(rawMap["detail"].([]interface{})),
	}
}

func buildBatchCreateCcRulesTagConditionBodyParam(d *schema.ResourceData) map[string]interface{} {
	rawArray, ok := d.Get("tag_condition").([]interface{})
	if !ok || len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category": utils.ValueIgnoreEmpty(rawMap["category"]),
		"contents": utils.ValueIgnoreEmpty(rawMap["contents"]),
	}
}

func buildBatchCreateCcRulesBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":               d.Get("name"),
		"mode":               1,
		"conditions":         buildBatchCreateCcRulesConditionsBodyParam(d),
		"action":             buildBatchCreateCcRulesActionBodyParam(d),
		"tag_type":           d.Get("tag_type"),
		"limit_num":          d.Get("limit_num"),
		"limit_period":       d.Get("limit_period"),
		"policy_ids":         d.Get("policy_ids"),
		"tag_index":          utils.ValueIgnoreEmpty(d.Get("tag_index")),
		"tag_condition":      buildBatchCreateCcRulesTagConditionBodyParam(d),
		"unlock_num":         utils.ValueIgnoreEmpty(d.Get("unlock_num")),
		"lock_time":          utils.ValueIgnoreEmpty(d.Get("lock_time")),
		"domain_aggregation": utils.ValueIgnoreEmpty(d.Get("domain_aggregation")),
		"region_aggregation": utils.ValueIgnoreEmpty(d.Get("region_aggregation")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceWafBatchCreateCcRulesCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/rule/cc"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildBatchCreateCcRulesQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildBatchCreateCcRulesBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error batch creating WAF CC rules: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceWafBatchCreateCcRulesRead(ctx, d, meta)
}

func resourceWafBatchCreateCcRulesRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateCcRulesUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is an action resource.
	return nil
}

func resourceWafBatchCreateCcRulesDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch create WAF cc rules. Deleting this resource
    will not remove the created rules, but will only remove the resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

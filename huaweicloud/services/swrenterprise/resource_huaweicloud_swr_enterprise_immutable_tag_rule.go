package swrenterprise

import (
	"context"
	"fmt"
	"strconv"
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

var enterpriseImmutableTagRuleNonUpdatableParams = []string{
	"instance_id", "namespace_name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules
// @API SWR GET /v2/{project_id}/instances/{instance_id}/immutabletagrules
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules/{immutable_rule_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules/{immutable_rule_id}
func ResourceSwrEnterpriseImmutableTagRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseImmutableTagRuleCreate,
		UpdateContext: resourceSwrEnterpriseImmutableTagRuleUpdate,
		ReadContext:   resourceSwrEnterpriseImmutableTagRuleRead,
		DeleteContext: resourceSwrEnterpriseImmutableTagRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseImmutableTagRuleNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the enterprise instance ID.`,
			},
			"namespace_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the namespace name.`,
			},
			"tag_selectors": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the repository version selector.`,
				Elem:        schemaSwrEnterpriseImmutableTagRuleRuleSelector(),
			},
			"scope_selectors": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the repository selectors.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the repository selector key.`,
						},
						"value": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Specifies the repository selector value.`,
							Elem:        schemaSwrEnterpriseImmutableTagRuleRuleSelector(),
						},
					},
				},
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the priority.`,
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the policy rule is disabled.`,
			},
			"action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the policy action.`,
			},
			"template": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the template type.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"namespace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the namespace ID`,
			},
			"immutable_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the policy ID`,
			},
		},
	}
}

func schemaSwrEnterpriseImmutableTagRuleRuleSelector() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the matching type.`,
			},
			"decoration": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the selector matching type.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pattern.`,
			},
			"extras": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the extra infos.`,
			},
		},
	}
}

func resourceSwrEnterpriseImmutableTagRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseImmutableTagRuleBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR immutable tag rule: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance immutable tag rule ID from the API response")
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseImmutableTagRuleRead(ctx, d, meta)
}

func buildCreateOrUpdateSwrEnterpriseImmutableTagRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"scope_selectors": buildSwrEnterpriseImmutableTagRuleRulesScopeSelectorsBodyParams(d.Get("scope_selectors")),
		"tag_selectors":   buildSwrEnterpriseImmutableTagRuleRulesTagSelectorBodyParams(d.Get("tag_selectors")),
		"priority":        utils.ValueIgnoreEmpty(d.Get("priority")),
		"disabled":        utils.ValueIgnoreEmpty(d.Get("disabled")),
		"action":          utils.ValueIgnoreEmpty(d.Get("action")),
		"template":        utils.ValueIgnoreEmpty(d.Get("template")),
	}

	return bodyParams
}

func buildSwrEnterpriseImmutableTagRuleRulesScopeSelectorsBodyParams(paramsList interface{}) map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make(map[string]interface{})
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				rst[param["key"].(string)] = buildSwrEnterpriseImmutableTagRuleRulesTagSelectorBodyParams(param["value"])
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseImmutableTagRuleRulesTagSelectorBodyParams(paramsList interface{}) []map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"kind":       param["kind"],
					"decoration": param["decoration"],
					"pattern":    param["pattern"],
					"extras":     param["extras"],
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func resourceSwrEnterpriseImmutableTagRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<immutable_rule_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/immutabletagrules?limit=100"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	offset := 0
	var rule interface{}
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%v", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error querying SWR instance immutable tag rule: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening SWR instance immutable tag rule response: %s", err)
		}

		rules := utils.PathSearch("immutable_rules", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(rules) == 0 {
			break
		}

		searchPath := fmt.Sprintf("immutable_rules[?namespace_name=='%s'&&id==`%s`]|[0]", namespaceName, id)
		rule = utils.PathSearch(searchPath, getRespBody, nil)
		if rule != nil {
			break
		}

		// offset must be the multiple of limit
		offset += 100
	}

	if rule == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error getting SWR immutable tag rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("namespace_name", namespaceName),
		d.Set("immutable_rule_id", id),
		d.Set("priority", utils.PathSearch("priority", rule, nil)),
		d.Set("disabled", utils.PathSearch("disabled", rule, nil)),
		d.Set("action", utils.PathSearch("action", rule, nil)),
		d.Set("template", utils.PathSearch("template", rule, nil)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", rule, nil)),
		d.Set("tag_selectors", flattenSwrEnterpriseImmutableTagRuleScopeRulesRuleSelector(
			utils.PathSearch("tag_selectors", rule, make([]interface{}, 0)))),
		d.Set("scope_selectors", flattenSwrEnterpriseImmutableTagRuleScopeRulesScopeSelectors(
			utils.PathSearch("scope_selectors", rule, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseImmutableTagRuleScopeRulesScopeSelectors(rawParams interface{}) []interface{} {
	if paramsMap, ok := rawParams.(map[string]interface{}); ok && len(paramsMap) > 0 {
		result := make([]interface{}, 0, len(paramsMap))
		for k, v := range paramsMap {
			m := map[string]interface{}{
				"key":   k,
				"value": flattenSwrEnterpriseImmutableTagRuleScopeRulesRuleSelector(v),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenSwrEnterpriseImmutableTagRuleScopeRulesRuleSelector(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"kind":       utils.PathSearch("kind", param, nil),
				"decoration": utils.PathSearch("decoration", param, nil),
				"pattern":    utils.PathSearch("pattern", param, nil),
				"extras":     utils.PathSearch("extras", param, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourceSwrEnterpriseImmutableTagRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules/{immutable_rule_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", d.Get("namespace_name").(string))
	updatePath = strings.ReplaceAll(updatePath, "{immutable_rule_id}", d.Get("immutable_rule_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseImmutableTagRuleBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR instance immutable tag rule: %s", err)
	}

	return resourceSwrEnterpriseImmutableTagRuleRead(ctx, d, meta)
}

func resourceSwrEnterpriseImmutableTagRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/immutabletagrules/{immutable_rule_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", d.Get("namespace_name").(string))
	deletePath = strings.ReplaceAll(deletePath, "{immutable_rule_id}", d.Get("immutable_rule_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR immutable tag rule")
	}

	return nil
}

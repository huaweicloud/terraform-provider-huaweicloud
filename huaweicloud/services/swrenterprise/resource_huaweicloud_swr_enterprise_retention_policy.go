package swrenterprise

import (
	"context"
	"encoding/json"
	"log"
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

var enterpriseRetentionPolicyNonUpdatableParams = []string{
	"instance_id", "namespace_name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}
func ResourceSwrEnterpriseRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseRetentionPolicyCreate,
		UpdateContext: resourceSwrEnterpriseRetentionPolicyUpdate,
		ReadContext:   resourceSwrEnterpriseRetentionPolicyRead,
		DeleteContext: resourceSwrEnterpriseRetentionPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseRetentionPolicyNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy name.`,
			},
			"algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the algorithm of policy.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether the policy is enabled.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the retention rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `Specifies the priority.`,
						},
						"action": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the policy action.`,
						},
						"template": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the template type.`,
						},
						// `params` type is Map<String, Object>
						"params": {
							Type:        schema.TypeMap,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the params.`,
						},
						"repo_scope_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the repo scope mode.`,
						},
						"tag_selectors": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Specifies the repository version selector.`,
							Elem:        schemaSwrEnterpriseRetentionPolicyRuleSelector(),
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
										Elem:        schemaSwrEnterpriseRetentionPolicyRuleSelector(),
									},
								},
							},
						},
						"disabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether the policy rule is disabled.`,
						},
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the retention policy rule ID.`,
						},
					},
				},
			},
			"trigger": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the trigger config.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the trigger type.`,
						},
						"trigger_settings": {
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: `Specifies the trigger settings.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cron": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: `Specifies the scheduled setting.`,
									},
								},
							},
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
			"namespace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the namespace ID`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the policy ID`,
			},
		},
	}
}

func schemaSwrEnterpriseRetentionPolicyRuleSelector() *schema.Resource {
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

func resourceSwrEnterpriseRetentionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseRetentionPolicyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR retention policy: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance retention policy ID from the API response")
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseRetentionPolicyRead(ctx, d, meta)
}

func buildCreateOrUpdateSwrEnterpriseRetentionPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      d.Get("name"),
		"algorithm": d.Get("algorithm"),
		"enabled":   d.Get("enabled"),
		"rules":     buildSwrEnterpriseRetentionPolicyScopeRulesBodyParams(d),
		"trigger":   buildSwrEnterpriseRetentionPolicyTriggerBodyParams(d),
	}

	return bodyParams
}

func buildSwrEnterpriseRetentionPolicyTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	if params := d.Get("trigger").([]interface{}); len(params) > 0 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"type":             param["type"],
				"trigger_settings": buildSwrEnterpriseRetentionPolicyTriggerSettingsBodyParams(param["trigger_settings"]),
			}

			return m
		}
	}

	return nil
}

func buildSwrEnterpriseRetentionPolicyTriggerSettingsBodyParams(rawParams interface{}) map[string]interface{} {
	if params := rawParams.([]interface{}); len(params) > 0 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"cron": utils.ValueIgnoreEmpty(param["cron"]),
			}

			return m
		}
	}

	return nil
}

func buildSwrEnterpriseRetentionPolicyScopeRulesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	if rules := d.Get("rules").([]interface{}); len(rules) > 0 {
		rst := make([]map[string]interface{}, 0, len(rules))
		for _, r := range rules {
			if rule, ok := r.(map[string]interface{}); ok {
				params := rule["params"].(map[string]interface{})
				for k, v := range params {
					params[k] = stringToJson(v.(string))
				}

				m := map[string]interface{}{
					"priority":        rule["priority"],
					"action":          rule["action"],
					"template":        rule["template"],
					"params":          params,
					"repo_scope_mode": rule["repo_scope_mode"],
					"disabled":        rule["disabled"],
					"scope_selectors": buildSwrEnterpriseRetentionPolicyRulesScopeSelectorsBodyParams(rule["scope_selectors"]),
					"tag_selectors":   buildSwrEnterpriseRetentionPolicyRulesTagSelectorBodyParams(rule["tag_selectors"]),
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func stringToJson(v string) interface{} {
	if v == "" {
		return nil
	}

	var data interface{}
	err := json.Unmarshal([]byte(v), &data)
	if err != nil {
		log.Printf("[DEBUG] Unable to parse JSON: %s", err)
		return v
	}

	return data
}

func buildSwrEnterpriseRetentionPolicyRulesScopeSelectorsBodyParams(paramsList interface{}) map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make(map[string]interface{})
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				rst[param["key"].(string)] = buildSwrEnterpriseRetentionPolicyRulesTagSelectorBodyParams(param["value"])
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseRetentionPolicyRulesTagSelectorBodyParams(paramsList interface{}) []map[string]interface{} {
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

func resourceSwrEnterpriseRetentionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<policy_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{namespace_name}", namespaceName)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR retention policy")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("namespace_name", namespaceName),
		d.Set("policy_id", id),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enabled", getRespBody, nil)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", getRespBody, nil)),
		d.Set("algorithm", utils.PathSearch("algorithm", getRespBody, nil)),
		d.Set("trigger", flattenSwrEnterpriseRetentionPolicyTrigger(getRespBody)),
		d.Set("rules", flattenSwrEnterpriseRetentionPolicyRetentionRules(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseRetentionPolicyTrigger(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("trigger", resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"type":             utils.PathSearch("type", param, nil),
			"trigger_settings": flattenSwrEnterpriseRetentionPolicyTriggerSettings(param),
		}
		return []interface{}{m}
	}

	return nil
}

func flattenSwrEnterpriseRetentionPolicyTriggerSettings(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("trigger_settings", resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"cron": utils.PathSearch("cron", param, nil),
		}
		return []interface{}{m}
	}

	return nil
}

func flattenSwrEnterpriseRetentionPolicyRetentionRules(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("rules", resp, make([]interface{}, 0))
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			params := utils.PathSearch("params", param, make(map[string]interface{})).(map[string]interface{})
			for k, v := range params {
				params[k] = utils.JsonToString(v)
			}

			m := map[string]interface{}{
				"id":              utils.PathSearch("id", param, nil),
				"priority":        utils.PathSearch("priority", param, nil),
				"action":          utils.PathSearch("action", param, nil),
				"template":        utils.PathSearch("template", param, nil),
				"params":          params,
				"repo_scope_mode": utils.PathSearch("repo_scope_mode", param, nil),
				"disabled":        utils.PathSearch("disabled", param, nil),
				"tag_selectors": flattenSwrEnterpriseRetentionPolicyScopeRulesRuleSelector(
					utils.PathSearch("tag_selectors", param, make([]interface{}, 0))),
				"scope_selectors": flattenSwrEnterpriseRetentionPolicyScopeRulesScopeSelectors(
					utils.PathSearch("scope_selectors", param, nil)),
			}

			result = append(result, m)
		}

		return result
	}

	return nil
}

func flattenSwrEnterpriseRetentionPolicyScopeRulesRuleSelector(rawParams interface{}) []interface{} {
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

func flattenSwrEnterpriseRetentionPolicyScopeRulesScopeSelectors(rawParams interface{}) []interface{} {
	if paramsMap, ok := rawParams.(map[string]interface{}); ok && len(paramsMap) > 0 {
		result := make([]interface{}, 0, len(paramsMap))
		for k, v := range paramsMap {
			m := map[string]interface{}{
				"key":   k,
				"value": flattenSwrEnterpriseRetentionPolicyScopeRulesRuleSelector(v),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourceSwrEnterpriseRetentionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", d.Get("namespace_name").(string))
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Get("policy_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseRetentionPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR instance retention policy: %s", err)
	}

	return resourceSwrEnterpriseRetentionPolicyRead(ctx, d, meta)
}

func resourceSwrEnterpriseRetentionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/retention/policies/{policy_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", d.Get("namespace_name").(string))
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Get("policy_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "SWR.400003"),
			"error deleting SWR retention policy")
	}

	return nil
}

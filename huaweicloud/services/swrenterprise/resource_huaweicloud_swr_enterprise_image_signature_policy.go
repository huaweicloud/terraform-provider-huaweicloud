package swrenterprise

import (
	"context"
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

var enterpriseImageSignaturePolicyNonUpdatableParams = []string{
	"instance_id", "namespace_name",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}
func ResourceSwrEnterpriseImageSignaturePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseImageSignaturePolicyCreate,
		UpdateContext: resourceSwrEnterpriseImageSignaturePolicyUpdate,
		ReadContext:   resourceSwrEnterpriseImageSignaturePolicyRead,
		DeleteContext: resourceSwrEnterpriseImageSignaturePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseImageSignaturePolicyNonUpdatableParams),

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
			"signature_method": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signature method.`,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signature algorithm.`,
			},
			"signature_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the signature key.`,
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
			"scope_rules": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the scope rules`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
										Elem:        schemaSwrEnterpriseImageSignaturePolicyRuleSelector(),
									},
								},
							},
						},
						"repo_scope_mode": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the repository select method.`,
						},
						"tag_selectors": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Specifies the repository version selector.`,
							Elem:        schemaSwrEnterpriseImageSignaturePolicyRuleSelector(),
						},
					},
				},
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the policy is enabled.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of policy.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creator`,
			},
			"namespace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the namespace ID`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
			"policy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the policy ID.`,
			},
		},
	}
}

func schemaSwrEnterpriseImageSignaturePolicyRuleSelector() *schema.Resource {
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

func resourceSwrEnterpriseImageSignaturePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseImageSignaturePolicyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR policy: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance policy ID from the API response")
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseImageSignaturePolicyRead(ctx, d, meta)
}

func buildCreateOrUpdateSwrEnterpriseImageSignaturePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                d.Get("name"),
		"trigger":             buildSwrEnterpriseImageSignaturePolicyTriggerBodyParams(d),
		"scope_rules":         buildSwrEnterpriseImageSignaturePolicyScopeRulesBodyParams(d),
		"signature_method":    d.Get("signature_method"),
		"signature_algorithm": d.Get("signature_algorithm"),
		"signature_key":       d.Get("signature_key"),
		"enabled":             d.Get("enabled"),
		"description":         d.Get("description"),
	}

	return bodyParams
}

func buildSwrEnterpriseImageSignaturePolicyTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	if params := d.Get("trigger").([]interface{}); len(params) == 1 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"type":             param["type"],
				"trigger_settings": buildSwrEnterpriseImageSignaturePolicyTriggerSettingsBodyParams(param["trigger_settings"]),
			}
			return m
		}
	}

	return nil
}

func buildSwrEnterpriseImageSignaturePolicyTriggerSettingsBodyParams(paramsList interface{}) map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) == 1 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"cron": param["cron"],
			}
			return m
		}
	}

	return nil
}

func buildSwrEnterpriseImageSignaturePolicyScopeRulesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	if params := d.Get("scope_rules").([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"scope_selectors": buildSwrEnterpriseImageSignaturePolicyScopeRulesScopeSelectorsBodyParams(param["scope_selectors"]),
					"repo_scope_mode": param["repo_scope_mode"],
					"tag_selectors":   buildSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelectorBodyParams(param["tag_selectors"]),
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseImageSignaturePolicyScopeRulesScopeSelectorsBodyParams(paramsList interface{}) map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make(map[string]interface{})
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				rst[param["key"].(string)] = buildSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelectorBodyParams(param["value"])
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelectorBodyParams(paramsList interface{}) []map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"kind":       param["kind"],
					"decoration": param["decoration"],
					"pattern":    param["pattern"],
					"extras":     utils.ValueIgnoreEmpty(param["extras"]),
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func resourceSwrEnterpriseImageSignaturePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}"
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
		return common.CheckDeletedDiag(d, err, "error retrieving SWR policy")
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
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("signature_method", utils.PathSearch("signature_method", getRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("signature_algorithm", getRespBody, nil)),
		d.Set("signature_key", utils.PathSearch("signature_key", getRespBody, nil)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", getRespBody, nil)),
		d.Set("creator", utils.PathSearch("creator", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("trigger", flattenSwrEnterpriseImageSignaturePolicyTrigger(getRespBody)),
		d.Set("scope_rules", flattenSwrEnterpriseImageSignaturePolicyScopeRules(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseImageSignaturePolicyTrigger(resp interface{}) []map[string]interface{} {
	rawParams := utils.PathSearch("trigger", resp, nil)
	if rawParams == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"type":             utils.PathSearch("type", rawParams, nil),
			"trigger_settings": flattenSwrEnterpriseImageSignaturePolicyTriggerSettings(rawParams),
		},
	}
}

func flattenSwrEnterpriseImageSignaturePolicyTriggerSettings(resp interface{}) []map[string]interface{} {
	rawParams := utils.PathSearch("trigger_settings", resp, nil)
	if rawParams == nil {
		return nil
	}
	return []map[string]interface{}{
		{
			"cron": utils.PathSearch("cron", rawParams, nil),
		},
	}
}

func flattenSwrEnterpriseImageSignaturePolicyScopeRules(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("scope_rules", resp, make([]interface{}, 0))
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"scope_selectors": flattenSwrEnterpriseImageSignaturePolicyScopeRulesScopeSelectors(
					utils.PathSearch("scope_selectors", param, nil)),
				"repo_scope_mode": utils.PathSearch("repo_scope_mode", param, nil),
				"tag_selectors": flattenSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelector(
					utils.PathSearch("tag_selectors", param, make([]interface{}, 0))),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelector(rawParams interface{}) []interface{} {
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

func flattenSwrEnterpriseImageSignaturePolicyScopeRulesScopeSelectors(rawParams interface{}) []interface{} {
	if paramsMap, ok := rawParams.(map[string]interface{}); ok && len(paramsMap) > 0 {
		result := make([]interface{}, 0, len(paramsMap))
		for k, v := range paramsMap {
			m := map[string]interface{}{
				"key":   k,
				"value": flattenSwrEnterpriseImageSignaturePolicyScopeRulesRuleSelector(v),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourceSwrEnterpriseImageSignaturePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", d.Get("namespace_name").(string))
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Get("policy_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseImageSignaturePolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR instance policy: %s", err)
	}

	return resourceSwrEnterpriseImageSignaturePolicyRead(ctx, d, meta)
}

func resourceSwrEnterpriseImageSignaturePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/signature/policies/{policy_id}"
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
		return common.CheckDeletedDiag(d, err, "error deleting SWR policy")
	}

	return nil
}

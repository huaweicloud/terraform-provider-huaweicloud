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

var enterpriseTriggerNonUpdatableParams = []string{
	"instance_id", "namespace_name", "targets.*.type", "targets.*.address_type", "targets.*.address",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies
// @API SWR GET /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}
func ResourceSwrEnterpriseTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseTriggerCreate,
		UpdateContext: resourceSwrEnterpriseTriggerUpdate,
		ReadContext:   resourceSwrEnterpriseTriggerRead,
		DeleteContext: resourceSwrEnterpriseTriggerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseTriggerNonUpdatableParams),

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
				Description: `Specifies the trigger name.`,
			},
			"targets": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the target params.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the trigger type.`,
						},
						"address_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the trigger address type.`,
						},
						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the trigger address.`,
						},
						"auth_header": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the auth header.`,
						},
						"skip_cert_verify": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: `Specifies whether to skip the verification of the certificate.`,
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
							Optional:    true,
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
										Elem:        schemaSwrEnterpriseTriggerRuleSelector(),
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
							Elem:        schemaSwrEnterpriseTriggerRuleSelector(),
						},
					},
				},
			},
			"event_types": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the event types of trigger.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the trigger is enabled.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of trigger.`,
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
			"trigger_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the trigger ID.`,
			},
		},
	}
}

func schemaSwrEnterpriseTriggerRuleSelector() *schema.Resource {
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
		},
	}
}

func resourceSwrEnterpriseTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	namespaceName := d.Get("namespace_name").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{namespace_name}", namespaceName)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseTriggerBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR trigger: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance trigger ID from the API response")
	}

	d.SetId(instanceId + "/" + namespaceName + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseTriggerRead(ctx, d, meta)
}

func buildCreateOrUpdateSwrEnterpriseTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"targets":     buildSwrEnterpriseTriggerTargetsBodyParams(d),
		"scope_rules": buildSwrEnterpriseTriggerScopeRulesBodyParams(d),
		"event_types": d.Get("event_types"),
		"enabled":     d.Get("enabled"),
		"description": d.Get("description"),
	}

	return bodyParams
}

func buildSwrEnterpriseTriggerTargetsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	if params := d.Get("targets").([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"type":             param["type"],
					"address_type":     param["address_type"],
					"address":          param["address"],
					"auth_header":      utils.ValueIgnoreEmpty(param["auth_header"]),
					"skip_cert_verify": utils.ValueIgnoreEmpty(param["skip_cert_verify"]),
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseTriggerScopeRulesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	if params := d.Get("scope_rules").([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"scope_selectors": buildSwrEnterpriseTriggerScopeRulesScopeSelectorsBodyParams(param["scope_selectors"]),
					"repo_scope_mode": param["repo_scope_mode"],
					"tag_selectors":   buildSwrEnterpriseTriggerScopeRulesRuleSelectorBodyParams(param["tag_selectors"]),
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseTriggerScopeRulesScopeSelectorsBodyParams(paramsList interface{}) map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make(map[string]interface{})
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				rst[param["key"].(string)] = buildSwrEnterpriseTriggerScopeRulesRuleSelectorBodyParams(param["value"])
			}
		}

		return rst
	}

	return nil
}

func buildSwrEnterpriseTriggerScopeRulesRuleSelectorBodyParams(paramsList interface{}) []map[string]interface{} {
	if params := paramsList.([]interface{}); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"kind":       param["kind"],
					"decoration": param["decoration"],
					"pattern":    param["pattern"],
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func resourceSwrEnterpriseTriggerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<namespace_name>/<trigger_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	namespaceName := parts[1]
	id := parts[2]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}"
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
		return common.CheckDeletedDiag(d, err, "error retrieving SWR trigger")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("namespace_name", namespaceName),
		d.Set("trigger_id", id),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enabled", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("event_types", utils.PathSearch("event_types", getRespBody, nil)),
		d.Set("namespace_id", utils.PathSearch("namespace_id", getRespBody, nil)),
		d.Set("creator", utils.PathSearch("creator", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
		d.Set("targets", flattenSwrEnterpriseTriggerTargets(getRespBody)),
		d.Set("scope_rules", flattenSwrEnterpriseTriggerScopeRules(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseTriggerTargets(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("targets", resp, make([]interface{}, 0))
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"type":             utils.PathSearch("type", param, nil),
				"address_type":     utils.PathSearch("address_type", param, nil),
				"address":          utils.PathSearch("address", param, nil),
				"auth_header":      utils.PathSearch("auth_header", param, nil),
				"skip_cert_verify": utils.PathSearch("skip_cert_verify", param, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenSwrEnterpriseTriggerScopeRules(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("scope_rules", resp, make([]interface{}, 0))
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"scope_selectors": flattenSwrEnterpriseTriggerScopeRulesScopeSelectors(
					utils.PathSearch("scope_selectors", param, nil)),
				"repo_scope_mode": utils.PathSearch("repo_scope_mode", param, nil),
				"tag_selectors": flattenSwrEnterpriseTriggerScopeRulesRuleSelector(
					utils.PathSearch("tag_selectors", param, make([]interface{}, 0))),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenSwrEnterpriseTriggerScopeRulesRuleSelector(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"kind":       utils.PathSearch("kind", param, nil),
				"decoration": utils.PathSearch("decoration", param, nil),
				"pattern":    utils.PathSearch("pattern", param, nil),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func flattenSwrEnterpriseTriggerScopeRulesScopeSelectors(rawParams interface{}) []interface{} {
	if paramsMap, ok := rawParams.(map[string]interface{}); ok && len(paramsMap) > 0 {
		result := make([]interface{}, 0, len(paramsMap))
		for k, v := range paramsMap {
			m := map[string]interface{}{
				"key":   k,
				"value": flattenSwrEnterpriseTriggerScopeRulesRuleSelector(v),
			}
			result = append(result, m)
		}
		return result
	}

	return nil
}

func resourceSwrEnterpriseTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{namespace_name}", d.Get("namespace_name").(string))
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Get("trigger_id").(string))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateSwrEnterpriseTriggerBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR instance trigger: %s", err)
	}

	return resourceSwrEnterpriseTriggerRead(ctx, d, meta)
}

func resourceSwrEnterpriseTriggerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/namespaces/{namespace_name}/webhook/policies/{policy_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{namespace_name}", d.Get("namespace_name").(string))
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Get("trigger_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR trigger")
	}

	return nil
}

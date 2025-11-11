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

var enterpriseReplicationPolicyNonUpdatableParams = []string{
	"instance_id",
}

// @API SWR POST /v2/{project_id}/instances/{instance_id}/replication/policies
// @API SWR GET /v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}
// @API SWR PUT /v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}
// @API SWR DELETE /v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}
func ResourceSwrEnterpriseReplicationPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrEnterpriseReplicationPolicyCreate,
		UpdateContext: resourceSwrEnterpriseReplicationPolicyUpdate,
		ReadContext:   resourceSwrEnterpriseReplicationPolicyRead,
		DeleteContext: resourceSwrEnterpriseReplicationPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(enterpriseReplicationPolicyNonUpdatableParams),

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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy name.`,
			},
			"repo_scope_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the repo scope mode.`,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Specifies whether the policy is enabled.`,
			},
			"filters": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: `Specifies the source resource filter.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the filter type.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the regular expression of the filter.`,
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
			"src_registry": {
				Type:         schema.TypeList,
				Optional:     true,
				Computed:     true,
				MaxItems:     1,
				Description:  `Specifies the source registry infos.`,
				Elem:         schemaSwrEnterpriseReplicationPolicyRegistry(),
				AtLeastOneOf: []string{"dest_registry"},
			},
			"dest_registry": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Specifies the destination registry infos.`,
				Elem:        schemaSwrEnterpriseReplicationPolicyRegistry(),
			},
			"dest_namespace": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the destination namespace name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of policy.`,
			},
			"override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether to override the repository.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"policy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the policy ID.`,
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
		},
	}
}

func schemaSwrEnterpriseReplicationPolicyRegistry() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the registry ID.`,
			},
		},
	}
}

func resourceSwrEnterpriseReplicationPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/policies"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateSwrEnterpriseReplicationPolicyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SWR replication policy: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := int(utils.PathSearch("id", createRespBody, float64(-1)).(float64))
	if id == -1 {
		return diag.Errorf("unable to find SWR instance replication policy ID from the API response")
	}

	d.SetId(instanceId + "/" + strconv.Itoa(id))

	return resourceSwrEnterpriseReplicationPolicyRead(ctx, d, meta)
}

func buildCreateSwrEnterpriseReplicationPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            d.Get("name"),
		"repo_scope_mode": d.Get("repo_scope_mode"),
		"enabled":         d.Get("enabled"),
		"filters":         buildSwrEnterpriseReplicationPolicyScopeRulesBodyParams(d),
		"trigger":         buildSwrEnterpriseReplicationPolicyTriggerBodyParams(d),
		"src_registry":    buildSwrEnterpriseReplicationPolicyRegistryBodyParams(d, "src_registry"),
		"dest_registry":   buildSwrEnterpriseReplicationPolicyRegistryBodyParams(d, "dest_registry"),
		"description":     d.Get("description"),
		"override":        d.Get("override"),
		"dest_namespace":  d.Get("dest_namespace"),
	}

	return bodyParams
}

func buildSwrEnterpriseReplicationPolicyTriggerBodyParams(d *schema.ResourceData) map[string]interface{} {
	if params := d.Get("trigger").([]interface{}); len(params) > 0 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"type":             param["type"],
				"trigger_settings": buildSwrEnterpriseReplicationPolicyTriggerSettingsBodyParams(param["trigger_settings"]),
			}

			return m
		}
	}

	return nil
}

func buildSwrEnterpriseReplicationPolicyRegistryBodyParams(d *schema.ResourceData, registry string) map[string]interface{} {
	if params := d.Get(registry).([]interface{}); len(params) > 0 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"id": param["id"],
			}

			return m
		}
	}

	return nil
}

func buildSwrEnterpriseReplicationPolicyTriggerSettingsBodyParams(rawParams interface{}) map[string]interface{} {
	if params := rawParams.([]interface{}); len(params) > 0 {
		if param, ok := params[0].(map[string]interface{}); ok {
			m := map[string]interface{}{
				"cron": param["cron"],
			}

			return m
		}
	}

	return nil
}

func buildSwrEnterpriseReplicationPolicyScopeRulesBodyParams(d *schema.ResourceData) []map[string]interface{} {
	if params := d.Get("filters").(*schema.Set).List(); len(params) > 0 {
		rst := make([]map[string]interface{}, 0, len(params))
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				m := map[string]interface{}{
					"type":  param["type"],
					"value": param["value"],
				}
				rst = append(rst, m)
			}
		}

		return rst
	}

	return nil
}

func resourceSwrEnterpriseReplicationPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<instance_id>/<policy_id>', but got '%s'", d.Id())
	}
	instanceId := parts[0]
	id := parts[1]

	getHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving SWR replication policy")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("enabled", utils.PathSearch("enabled", getRespBody, nil)),
		d.Set("dest_namespace", utils.PathSearch("dest_namespace", getRespBody, nil)),
		d.Set("repo_scope_mode", utils.PathSearch("repo_scope_mode", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("override", utils.PathSearch("override", getRespBody, nil)),
		d.Set("trigger", flattenSwrEnterpriseReplicationPolicyTrigger(getRespBody)),
		d.Set("filters", flattenSwrEnterpriseReplicationPolicyReplicationRules(getRespBody)),
		d.Set("src_registry", flattenSwrEnterpriseReplicationPolicyRegistry(getRespBody, "src_registry")),
		d.Set("dest_registry", flattenSwrEnterpriseReplicationPolicyRegistry(getRespBody, "dest_registry")),
		d.Set("policy_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_at", getRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_at", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSwrEnterpriseReplicationPolicyTrigger(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("trigger", resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"type":             utils.PathSearch("type", param, nil),
			"trigger_settings": flattenSwrEnterpriseReplicationPolicyTriggerSettings(param),
		}
		return []interface{}{m}
	}

	return nil
}

func flattenSwrEnterpriseReplicationPolicyTriggerSettings(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("trigger_settings", resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"cron": utils.PathSearch("cron", param, nil),
		}
		return []interface{}{m}
	}

	return nil
}

func flattenSwrEnterpriseReplicationPolicyRegistry(resp interface{}, registry string) []interface{} {
	rawParams := utils.PathSearch(registry, resp, nil)
	if param, ok := rawParams.(map[string]interface{}); ok {
		m := map[string]interface{}{
			"id": utils.PathSearch("id", param, nil),
		}
		return []interface{}{m}
	}

	return nil
}

func flattenSwrEnterpriseReplicationPolicyReplicationRules(resp interface{}) []interface{} {
	rawParams := utils.PathSearch("filters", resp, make([]interface{}, 0))
	if paramsList, ok := rawParams.([]interface{}); ok && len(paramsList) > 0 {
		result := make([]interface{}, 0, len(paramsList))
		for _, param := range paramsList {
			m := map[string]interface{}{
				"type":  utils.PathSearch("type", param, nil),
				"value": utils.PathSearch("value", param, nil),
			}

			result = append(result, m)
		}

		return result
	}

	return nil
}

func resourceSwrEnterpriseReplicationPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", strconv.Itoa(d.Get("policy_id").(int)))
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateSwrEnterpriseReplicationPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SWR instance replication policy: %s", err)
	}

	return resourceSwrEnterpriseReplicationPolicyRead(ctx, d, meta)
}

func buildUpdateSwrEnterpriseReplicationPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":            d.Get("name"),
		"repo_scope_mode": d.Get("repo_scope_mode"),
		"enabled":         d.Get("enabled"),
		"filters":         buildSwrEnterpriseReplicationPolicyScopeRulesBodyParams(d),
		"trigger":         buildSwrEnterpriseReplicationPolicyTriggerBodyParams(d),
		"description":     d.Get("description"),
		"override":        d.Get("override"),
		"dest_namespace":  d.Get("dest_namespace"),
	}

	// the direction can be changed, one of them is computed
	src := buildSwrEnterpriseReplicationPolicyRegistryBodyParams(d, "src_registry")
	dest := buildSwrEnterpriseReplicationPolicyRegistryBodyParams(d, "dest_registry")

	switch {
	case !d.HasChanges("src_registry", "dest_registry"), (d.HasChange("src_registry") && d.HasChange("dest_registry")):
		bodyParams["src_registry"] = src
		bodyParams["dest_registry"] = dest
	case d.HasChange("src_registry"):
		bodyParams["src_registry"] = src
	case d.HasChange("dest_registry"):
		bodyParams["dest_registry"] = dest
	}

	return bodyParams
}

func resourceSwrEnterpriseReplicationPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("swr", region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/replication/policies/{policy_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", strconv.Itoa(d.Get("policy_id").(int)))
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SWR replication policy")
	}

	return nil
}

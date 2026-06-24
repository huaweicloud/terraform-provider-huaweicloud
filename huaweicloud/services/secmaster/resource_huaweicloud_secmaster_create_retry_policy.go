package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/policys
func ResourceCreateRetryPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateRetryPolicyCreate,
		UpdateContext: resourceCreateRetryPolicyUpdate,
		ReadContext:   resourceCreateRetryPolicyRead,
		DeleteContext: resourceCreateRetryPolicyDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"workspace_id",
			"action_type",
			"version",
			"retry_list",
			"block_age.*.is_block_ageing",
			"block_age.*.block_ageing",
			"block_target",
			"defense_policy_list.*.defense_connection_id",
			"defense_policy_list.*.defense_connection_name",
			"defense_policy_list.*.defense_connection_region_id",
			"defense_policy_list.*.defense_connection_region_name",
			"defense_policy_list.*.defense_type",
			"defense_policy_list.*.target_enterprise_id",
			"defense_policy_list.*.target_enterprise_name",
			"defense_policy_list.*.target_project_id",
			"defense_policy_list.*.target_project_name",
			"description",
			"labels",
			"policy_category",
			"policy_type.*.policy_type",
			"policy_direction",
			"account_scope",
			"eps_scope",
			"region_scope",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"action_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"block_age": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     buildBlockAgeSchema(),
			},
			"block_target": {
				Type:     schema.TypeString,
				Required: true,
			},
			"defense_policy_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     buildDefensePolicyListSchema(),
			},
			"policy_category": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_type": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     buildPolicyTypeSchema(),
			},
			"retry_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_direction": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"account_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"eps_scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"region_scope": {
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

func buildBlockAgeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"is_block_ageing": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"block_ageing": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildDefensePolicyListSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"defense_connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"defense_connection_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"defense_connection_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"defense_connection_region_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"defense_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_enterprise_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_enterprise_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"target_project_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildPolicyTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func buildBlockAgeBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"is_block_ageing": rawMap["is_block_ageing"],
		"block_ageing":    utils.ValueIgnoreEmpty(rawMap["block_ageing"]),
	}
}

func buildDefensePolicyListBodyParams(rawArray []interface{}) []map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		params := map[string]interface{}{
			"defense_connection_id":          rawMap["defense_connection_id"],
			"defense_connection_name":        utils.ValueIgnoreEmpty(rawMap["defense_connection_name"]),
			"defense_connection_region_id":   utils.ValueIgnoreEmpty(rawMap["defense_connection_region_id"]),
			"defense_connection_region_name": utils.ValueIgnoreEmpty(rawMap["defense_connection_region_name"]),
			"defense_type":                   utils.ValueIgnoreEmpty(rawMap["defense_type"]),
			"target_enterprise_id":           utils.ValueIgnoreEmpty(rawMap["target_enterprise_id"]),
			"target_enterprise_name":         utils.ValueIgnoreEmpty(rawMap["target_enterprise_name"]),
			"target_project_id":              utils.ValueIgnoreEmpty(rawMap["target_project_id"]),
			"target_project_name":            utils.ValueIgnoreEmpty(rawMap["target_project_name"]),
		}

		result = append(result, params)
	}

	return result
}

func buildPolicyTypeBodyParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"policy_type": rawMap["policy_type"],
	}
}

func buildCreateRetryPolicyBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"data_object": map[string]interface{}{
			"retry_list":          utils.ValueIgnoreEmpty(d.Get("retry_list")),
			"block_age":           buildBlockAgeBodyParams(d.Get("block_age").([]interface{})),
			"block_target":        d.Get("block_target"),
			"defense_policy_list": buildDefensePolicyListBodyParams(d.Get("defense_policy_list").([]interface{})),
			"description":         utils.ValueIgnoreEmpty(d.Get("description")),
			"labels":              utils.ValueIgnoreEmpty(d.Get("labels")),
			"policy_category":     d.Get("policy_category"),
			"policy_type":         buildPolicyTypeBodyParams(d.Get("policy_type").([]interface{})),
			"region_id":           region,
			"policy_direction":    utils.ValueIgnoreEmpty(d.Get("policy_direction")),
			"account_scope":       utils.ValueIgnoreEmpty(d.Get("account_scope")),
			"eps_scope":           utils.ValueIgnoreEmpty(d.Get("eps_scope")),
			"region_scope":        utils.ValueIgnoreEmpty(d.Get("region_scope")),
		},
	}
}

func resourceCreateRetryPolicyCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/policys"
		workspaceId = d.Get("workspace_id").(string)
		actionType  = d.Get("action_type").(string)
		version     = d.Get("version").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)
	requestPath = fmt.Sprintf("%s?action_type=%s", requestPath, actionType)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type":        "application/json;charset=UTF-8",
			"X-Secmaster-Version": version,
		},
		JSONBody: utils.RemoveNil(buildCreateRetryPolicyBodyParams(d, region)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating/retrying SecMaster policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening response body: %s", err)
	}

	data := utils.PathSearch("data", respBody, "").(string)
	if data == "" {
		return diag.Errorf("unable to find `data` in response body")
	}

	d.SetId(data)

	return nil
}

func resourceCreateRetryPolicyRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCreateRetryPolicyUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCreateRetryPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource using to create or retry policy. Deleting this resource will not
	 change the status of the currently SecMaster policy resource, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

package hss

import (
	"context"
	"errors"
	"fmt"
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

// @API HSS GET /v5/{project_id}/custom/rule/config/detail
// @API HSS POST /v5/{project_id}/custom/rule/config/operate
// @API HSS POST /v5/{project_id}/custom/rule/config
// @API HSS GET /v5/{project_id}/custom/rule/config
// @API HSS PUT /v5/{project_id}/custom/rule/config
// @API HSS DELETE /v5/{project_id}/custom/rule/config
func ResourceCustomRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCustomRuleCreate,
		ReadContext:   resourceCustomRuleRead,
		UpdateContext: resourceCustomRuleUpdate,
		DeleteContext: resourceCustomRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			// The API works fine when the rule name is changed, but the changes don't take effect.
			// A special note has been added to the documentation that the provider will not restrict the API's behavior.
			"rule_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"custom_rule_value_info": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     customRuleValueInfoSchema(),
			},
			// This fields defaults to false in API definition.
			"is_all_host": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// The value of the `agent_ids` field does not always meet expectations.
			// For example, a response timeout, a change in value, etc.
			// Therefore, an additional attribute field is provided to record its value.
			"agent_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Defaults to `1`, means enabled. `0` means disabled.
			"rule_status": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"host_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"agent_ids_attr": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func customRuleValueInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"rule_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"hash_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auto_block": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"rule_values": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildCustomRuleValueInfoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"rule_type":   rawMap["rule_type"],
		"hash_type":   rawMap["hash_type"],
		"auto_block":  rawMap["auto_block"],
		"rule_values": rawMap["rule_values"],
	}
}

func buildAgentIdsParam(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	return utils.ExpandToStringList(rawArray)
}

func buildCreateCustomRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"rule_name":              d.Get("rule_name"),
		"custom_rule_value_info": buildCustomRuleValueInfoParams(d.Get("custom_rule_value_info").([]interface{})),
		"is_all_host":            d.Get("is_all_host"),
		"agent_ids":              buildAgentIdsParam(d.Get("agent_ids").([]interface{})),
	}
}

func updateCustomRuleRunStatus(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v5/{project_id}/custom/rule/config/operate"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		JSONBody: map[string]interface{}{
			"enable": d.Get("rule_status"),
			"rule_id_list": []string{
				d.Id(),
			},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceCustomRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/custom/rule/config"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCustomRuleBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating HSS custom rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("rule_id", respBody, "").(string)
	if id == "" {
		return diag.FromErr(errors.New("error creating HSS custom rule: empty rule_id in response"))
	}
	d.SetId(id)

	if d.Get("rule_status").(int) == 0 {
		// close the custom rule
		if err := updateCustomRuleRunStatus(client, d); err != nil {
			return diag.Errorf("error updating HSS custom rule status in creation operation: %s", err)
		}
	}

	return resourceCustomRuleRead(ctx, d, meta)
}

// We need to use both `GetCustomRuleByDetail` and `GetCustomRuleByList` to get the complete information of a custom rule.
func GetCustomRuleByList(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/custom/rule/config"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?rule_id=%s", ruleId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	target := utils.PathSearch("data_list|[0]", respBody, nil)
	if target == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return target, nil
}

func GetCustomRuleByDetail(client *golangsdk.ServiceClient, ruleId string) (interface{}, error) {
	requestPath := client.Endpoint + "v5/{project_id}/custom/rule/config/detail"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?rule_id=%s", ruleId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCustomRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	targetFromList, err := GetCustomRuleByList(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving custom rule from list API")
	}

	targetFromDetail, err := GetCustomRuleByDetail(client, d.Id())
	if err != nil {
		return diag.Errorf("error retrieving custom rule from detail API: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rule_name", utils.PathSearch("rule_name", targetFromList, nil)),
		d.Set("is_all_host", utils.PathSearch("is_all_host", targetFromList, nil)),
		d.Set("rule_status", utils.PathSearch("rule_status", targetFromList, nil)),
		d.Set("create_time", utils.PathSearch("create_time", targetFromList, nil)),
		d.Set("update_time", utils.PathSearch("update_time", targetFromList, nil)),
		d.Set("host_num", utils.PathSearch("host_num", targetFromList, nil)),
		d.Set("custom_rule_value_info", flattenCustomRuleValueInfo(targetFromDetail)),
		d.Set("agent_ids_attr", utils.PathSearch("agent_ids", targetFromDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCustomRuleValueInfo(respBody interface{}) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"rule_type":   utils.PathSearch("rule_type", respBody, nil),
			"hash_type":   utils.PathSearch("hash_type", respBody, nil),
			"auto_block":  utils.PathSearch("auto_block", respBody, nil),
			"rule_values": utils.PathSearch("rule_values", respBody, nil),
		},
	}
}

func buildUpdateCustomRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return utils.RemoveNil(map[string]interface{}{
		"rule_id":                d.Id(),
		"rule_name":              d.Get("rule_name"),
		"custom_rule_value_info": buildCustomRuleValueInfoParams(d.Get("custom_rule_value_info").([]interface{})),
		"is_all_host":            d.Get("is_all_host"),
		"agent_ids":              buildAgentIdsParam(d.Get("agent_ids").([]interface{})),
	})
}

func resourceCustomRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	if d.HasChanges("rule_name", "custom_rule_value_info", "is_all_host", "agent_ids") {
		requestPath := client.Endpoint + "v5/{project_id}/custom/rule/config"
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateCustomRuleBodyParams(d)),
		}

		_, err := client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating HSS custom rule: %s", err)
		}
	}

	if d.HasChange("rule_status") {
		if err := updateCustomRuleRunStatus(client, d); err != nil {
			return diag.Errorf("error updating HSS custom rule status in update operation: %s", err)
		}
	}

	return resourceCustomRuleRead(ctx, d, meta)
}

func resourceCustomRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/{project_id}/custom/rule/config"
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"rule_id_list": []string{d.Id()},
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS custom rule: %s", err)
	}

	return nil
}

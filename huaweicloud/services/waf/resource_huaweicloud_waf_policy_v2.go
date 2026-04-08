package waf

import (
	"context"
	"fmt"
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

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}
// @API WAF PATCH /v1/{project_id}/waf/policy/{policy_id}
// @API WAF POST /v1/{project_id}/waf/policy
func ResourceWafPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafPolicyV2Create,
		ReadContext:   resourceWafPolicyV2Read,
		UpdateContext: resourceWafPolicyV2Update,
		DeleteContext: resourceWafPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"log_action_replaced",
			"enterprise_project_id",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API documentation specifies the type as Boolean.
			// Has no response value.
			// Cannot be update.
			"log_action_replaced": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Cannot be update.
			// Has no response value.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			// The API documentation specifies the type as Boolean.
			"full_detection": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"robot_action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     policyV2RobotActionSchema(),
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     policyV2ActionSchema(),
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     policyV2OptionsSchema(),
			},
			"hosts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"bind_host": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     policyV2BindHostSchema(),
			},
			// The structure of this field is complex, so a set operation will not be performed;
			// instead, a separate attribute field (`extend_attribute`) will be provided.
			"extend": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"extend_attribute": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func policyV2RobotActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func policyV2ActionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"followed_action_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

// These fields are classified as boolean in the API documentation.
func policyV2OptionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"webattack": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"common": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crawler": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"crawler_engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crawler_scanner": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crawler_script": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"crawler_other": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"webshell": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cc": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"custom": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"whiteblackip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"geoip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ignore": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"privacy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"antitamper": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"antileakage": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bot_enable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"modulex_enabled": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func policyV2BindHostSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"waf_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildWafPolicyV2QueryParams(epsId string) string {
	if epsId == "" {
		return ""
	}

	return fmt.Sprintf("?enterprise_project_id=%s", epsId)
}

func convertStringToBool(stringValue string) interface{} {
	if stringValue == "" {
		return nil
	}

	boolValue, err := strconv.ParseBool(stringValue)
	if err != nil {
		log.Printf("[ERROR] error converting string %s to boolean: %s", stringValue, err)
		return nil
	}

	return boolValue
}

func buildCreateWafPolicyV2BodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                d.Get("name"),
		"log_action_replaced": convertStringToBool(d.Get("log_action_replaced").(string)),
	}
}

func buildRobotActionBodyparam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category": utils.ValueIgnoreEmpty(rawMap["category"]),
	}
}

func buildPolicyV2ActionBodyparam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"category":           utils.ValueIgnoreEmpty(rawMap["category"]),
		"followed_action_id": utils.ValueIgnoreEmpty(rawMap["followed_action_id"]),
	}
}

func buildOptionsBodyparam(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"webattack":       convertStringToBool(rawMap["webattack"].(string)),
		"common":          convertStringToBool(rawMap["common"].(string)),
		"crawler":         convertStringToBool(rawMap["crawler"].(string)),
		"crawler_engine":  convertStringToBool(rawMap["crawler_engine"].(string)),
		"crawler_scanner": convertStringToBool(rawMap["crawler_scanner"].(string)),
		"crawler_script":  convertStringToBool(rawMap["crawler_script"].(string)),
		"crawler_other":   convertStringToBool(rawMap["crawler_other"].(string)),
		"webshell":        convertStringToBool(rawMap["webshell"].(string)),
		"cc":              convertStringToBool(rawMap["cc"].(string)),
		"custom":          convertStringToBool(rawMap["custom"].(string)),
		"whiteblackip":    convertStringToBool(rawMap["whiteblackip"].(string)),
		"geoip":           convertStringToBool(rawMap["geoip"].(string)),
		"ignore":          convertStringToBool(rawMap["ignore"].(string)),
		"privacy":         convertStringToBool(rawMap["privacy"].(string)),
		"antitamper":      convertStringToBool(rawMap["antitamper"].(string)),
		"antileakage":     convertStringToBool(rawMap["antileakage"].(string)),
		"bot_enable":      convertStringToBool(rawMap["bot_enable"].(string)),
		"modulex_enabled": convertStringToBool(rawMap["modulex_enabled"].(string)),
	}
}

func buildHostsBodyparam(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	return utils.ExpandToStringList(rawArray)
}

func buildBindHostBodyparam(rawArray []interface{}) interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		rst = append(rst, map[string]interface{}{
			"id":       utils.ValueIgnoreEmpty(rawMap["id"]),
			"hostname": utils.ValueIgnoreEmpty(rawMap["hostname"]),
			"waf_type": utils.ValueIgnoreEmpty(rawMap["waf_type"]),
			"mode":     utils.ValueIgnoreEmpty(rawMap["mode"]),
		})
	}
	return rst
}

func buildExtendBodyparam(rawMap map[string]interface{}) interface{} {
	if len(rawMap) == 0 {
		return nil
	}

	return utils.ExpandToStringMap(rawMap)
}

func buildUpdateWafPolicyV2BodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":           d.Get("name"),
		"level":          utils.ValueIgnoreEmpty(d.Get("level")), // A zero value for this field is meaningless.
		"full_detection": convertStringToBool(d.Get("full_detection").(string)),
		"robot_action":   buildRobotActionBodyparam(d.Get("robot_action").([]interface{})),
		"action":         buildPolicyV2ActionBodyparam(d.Get("action").([]interface{})),
		"options":        buildOptionsBodyparam(d.Get("options").([]interface{})),
		"hosts":          buildHostsBodyparam(d.Get("hosts").([]interface{})),
		"bind_host":      buildBindHostBodyparam(d.Get("bind_host").([]interface{})),
		"extend":         buildExtendBodyparam(d.Get("extend").(map[string]interface{})),
	}
}

func updateWafPolicyV2(client *golangsdk.ServiceClient, d *schema.ResourceData, epsId string) error {
	requestPath := client.Endpoint + "v1/{project_id}/waf/policy/{policy_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Id())
	requestPath += buildWafPolicyV2QueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildUpdateWafPolicyV2BodyParams(d)),
	}

	_, err := client.Request("PATCH", requestPath, &requestOpt)
	return err
}

func resourceWafPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildWafPolicyV2QueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildCreateWafPolicyV2BodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating WAF policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating WAF policy: ID is not found in API response")
	}
	d.SetId(id)

	if err := updateWafPolicyV2(client, d, epsId); err != nil {
		return diag.Errorf("error updating WAF policy in creation: %s", err)
	}

	return resourceWafPolicyV2Read(ctx, d, meta)
}

func convertBoolValueToString(respValue interface{}) string {
	if respValue == nil {
		return ""
	}

	boolValue, ok := respValue.(bool)
	if !ok {
		log.Printf("[ERROR] the response value %v is not boolean value", respValue)
		return ""
	}

	if boolValue {
		return "true"
	}
	return "false"
}

func flattenPolicyV2RobotAction(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	mapRst := map[string]interface{}{
		"category": utils.PathSearch("category", respBody, nil),
	}

	return []interface{}{mapRst}
}

func flattenPolicyV2Action(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	mapRst := map[string]interface{}{
		"category":           utils.PathSearch("category", respBody, nil),
		"followed_action_id": utils.PathSearch("followed_action_id", respBody, nil),
	}

	return []interface{}{mapRst}
}

func flattenPolicyV2Options(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	mapRst := map[string]interface{}{
		"webattack":       convertBoolValueToString(utils.PathSearch("webattack", respBody, nil)),
		"common":          convertBoolValueToString(utils.PathSearch("common", respBody, nil)),
		"crawler":         convertBoolValueToString(utils.PathSearch("crawler", respBody, nil)),
		"crawler_engine":  convertBoolValueToString(utils.PathSearch("crawler_engine", respBody, nil)),
		"crawler_scanner": convertBoolValueToString(utils.PathSearch("crawler_scanner", respBody, nil)),
		"crawler_script":  convertBoolValueToString(utils.PathSearch("crawler_script", respBody, nil)),
		"crawler_other":   convertBoolValueToString(utils.PathSearch("crawler_other", respBody, nil)),
		"webshell":        convertBoolValueToString(utils.PathSearch("webshell", respBody, nil)),
		"cc":              convertBoolValueToString(utils.PathSearch("cc", respBody, nil)),
		"custom":          convertBoolValueToString(utils.PathSearch("custom", respBody, nil)),
		"whiteblackip":    convertBoolValueToString(utils.PathSearch("whiteblackip", respBody, nil)),
		"geoip":           convertBoolValueToString(utils.PathSearch("geoip", respBody, nil)),
		"ignore":          convertBoolValueToString(utils.PathSearch("ignore", respBody, nil)),
		"privacy":         convertBoolValueToString(utils.PathSearch("privacy", respBody, nil)),
		"antitamper":      convertBoolValueToString(utils.PathSearch("antitamper", respBody, nil)),
		"antileakage":     convertBoolValueToString(utils.PathSearch("antileakage", respBody, nil)),
		"bot_enable":      convertBoolValueToString(utils.PathSearch("bot_enable", respBody, nil)),
		"modulex_enabled": convertBoolValueToString(utils.PathSearch("modulex_enabled", respBody, nil)),
	}

	return []interface{}{mapRst}
}

func flattenPolicyV2BindHost(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"id":       utils.PathSearch("id", v, nil),
			"hostname": utils.PathSearch("hostname", v, nil),
			"waf_type": utils.PathSearch("waf_type", v, nil),
			"mode":     utils.PathSearch("mode", v, nil),
		})
	}

	return rst
}

func GetWafPolicyV2(client *golangsdk.ServiceClient, policyId string, epsId string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/waf/policy/{policy_id}"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", policyId)
	requestPath += buildWafPolicyV2QueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceWafPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	respBody, err := GetWafPolicyV2(client, d.Id(), epsId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("level", utils.PathSearch("level", respBody, nil)),
		d.Set("full_detection", convertBoolValueToString(utils.PathSearch("full_detection", respBody, nil))),
		d.Set("robot_action", flattenPolicyV2RobotAction(utils.PathSearch("robot_action", respBody, nil))),
		d.Set("action", flattenPolicyV2Action(utils.PathSearch("action", respBody, nil))),
		d.Set("options", flattenPolicyV2Options(utils.PathSearch("options", respBody, nil))),
		d.Set("hosts", utils.ExpandToStringList(utils.PathSearch("hosts", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("bind_host", flattenPolicyV2BindHost(utils.PathSearch("bind_host", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("extend_attribute", utils.PathSearch("extend", respBody, nil)),
		d.Set("timestamp", utils.PathSearch("timestamp", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if err := updateWafPolicyV2(client, d, epsId); err != nil {
		return diag.Errorf("error updating WAF policy in update operation: %s", err)
	}

	return resourceWafPolicyV2Read(ctx, d, meta)
}

func resourceWafPolicyV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/policy/{policy_id}"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Id())
	requestPath += buildWafPolicyV2QueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting WAF policy: %s", err)
	}

	return nil
}

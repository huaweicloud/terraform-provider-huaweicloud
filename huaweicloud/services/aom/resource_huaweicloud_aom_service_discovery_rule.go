package aom

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var basicNameRule = schema.Schema{
	Type:     schema.TypeList,
	Required: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"cmdLineHash", "cmdLine", "env", "str",
				}, false),
			},
			"args": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	},
}

// @API AOM GET /v1/{project_id}/inv/servicediscoveryrules
// @API AOM PUT /v1/{project_id}/inv/servicediscoveryrules
// @API AOM DELETE /v1/{project_id}/inv/servicediscoveryrules
func ResourceServiceDiscoveryRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceDiscoveryRuleCreateOrUpdate,
		ReadContext:   resourceServiceDiscoveryRuleRead,
		UpdateContext: resourceServiceDiscoveryRuleCreateOrUpdate,
		DeleteContext: resourceServiceDiscoveryRuleDelete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"discovery_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"check_content": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"check_mode": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"contain", "equals",
							}, false),
						},
						"check_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"cmdLine", "env", "scope",
							}, false),
						},
					},
				},
			},
			"name_rules": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_name_rule":     &basicNameRule,
						"application_name_rule": &basicNameRule,
					},
				},
			},
			"log_file_suffix": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"discovery_rule_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"is_default_rule": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"detect_log_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  9999,
			},
			"log_path_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name_type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"cmdLineHash",
							}, false),
						},
						"args": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"value": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildServiceDiscoveryRuleBodyParams(d *schema.ResourceData, pid string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":        d.Get("rule_id"),
		"enable":    d.Get("discovery_rule_enabled"),
		"eventName": "aom_inventory_rules_event",
		"name":      d.Get("name"),
		"projectid": pid,
		"desc":      d.Get("description"),
		"spec":      buildServiceDiscoveryRuleBodyParamsSpec(d),
	}

	return map[string]interface{}{
		"appRules": []interface{}{bodyParams},
	}
}

func buildServiceDiscoveryRuleBodyParamsSpec(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"appType":       d.Get("service_type"),
		"detectLog":     strconv.FormatBool(d.Get("detect_log_enabled").(bool)),
		"discoveryRule": buildDiscoveryRuleOpts(d.Get("discovery_rules").([]interface{})),
		"isDefaultRule": strconv.FormatBool(d.Get("is_default_rule").(bool)),
		"isDetect":      "false",
		"logFileFix":    d.Get("log_file_suffix"),
		"logPathRule":   buildBasicRuleOpts(d.Get("log_path_rules").([]interface{})),
		"nameRule":      buildNameRuleOpts(d.Get("name_rules").([]interface{})),
		"priority":      d.Get("priority"),
	}

	return bodyParams
}

func buildDiscoveryRuleOpts(rawRules []interface{}) []interface{} {
	discoveryRules := make([]interface{}, 0, len(rawRules))
	for _, v := range rawRules {
		rawRule := v.(map[string]interface{})
		discoveryRules = append(discoveryRules, map[string]interface{}{
			"checkContent": rawRule["check_content"],
			"checkType":    rawRule["check_type"],
			"checkMode":    rawRule["check_mode"],
		})
	}
	return discoveryRules
}

func buildBasicRuleOpts(rawRules []interface{}) []interface{} {
	if len(rawRules) == 0 {
		return nil
	}

	logPathRules := make([]interface{}, 0, len(rawRules))
	for _, v := range rawRules {
		rawRule := v.(map[string]interface{})
		logPathRules = append(logPathRules, map[string]interface{}{
			"args":     rawRule["args"],
			"nameType": rawRule["name_type"],
			"value":    rawRule["value"],
		})
	}
	return logPathRules
}

func buildNameRuleOpts(rawRules []interface{}) interface{} {
	if len(rawRules) != 1 {
		return nil
	}
	raw := rawRules[0].(map[string]interface{})
	return map[string]interface{}{
		"appNameRule":         buildBasicRuleOpts(raw["service_name_rule"].([]interface{})),
		"applicationNameRule": buildBasicRuleOpts(raw["application_name_rule"].([]interface{})),
	}
}

func resourceServiceDiscoveryRuleCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createHttpUrl := "v1/{project_id}/inv/servicediscoveryrules"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildServiceDiscoveryRuleBodyParams(d, client.ProjectID)),
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error AOM service discovery rule: %s", err)
	}

	d.SetId(d.Get("name").(string))

	// wait for the configuration to take effect
	// lintignore:R018
	time.Sleep(30 * time.Second)

	return resourceServiceDiscoveryRuleRead(ctx, d, meta)
}

func resourceServiceDiscoveryRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	rule, err := GetServiceDiscoveryRule(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM service discovery rule")
	}

	isDefaultRule, _ := strconv.ParseBool(utils.PathSearch("spec.isDefaultRule", rule, "").(string))
	detectLogEnabled, _ := strconv.ParseBool(utils.PathSearch("spec.detectLog", rule, "").(string))
	createdAt, _ := strconv.ParseInt(utils.PathSearch("createTime", rule, "").(string), 10, 64)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("name", rule, nil)),
		d.Set("rule_id", utils.PathSearch("id", rule, nil)),
		d.Set("discovery_rule_enabled", utils.PathSearch("enable", rule, nil)),
		d.Set("is_default_rule", isDefaultRule),
		d.Set("log_file_suffix", utils.PathSearch("spec.logFileFix", rule, nil)),
		d.Set("service_type", utils.PathSearch("spec.appType", rule, nil)),
		d.Set("detect_log_enabled", detectLogEnabled),
		d.Set("priority", utils.PathSearch("spec.priority", rule, nil)),
		d.Set("discovery_rules", flattenDiscoveryRules(
			utils.PathSearch("spec.discoveryRule", rule, make([]interface{}, 0)).([]interface{}))),
		d.Set("log_path_rules", flattenBasicRules(
			utils.PathSearch("spec.logPathRule", rule, make([]interface{}, 0)).([]interface{}))),
		d.Set("name_rules", flattenNameRulesRules(utils.PathSearch("spec.nameRule", rule, nil))),
		d.Set("description", utils.PathSearch("desc", rule, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(createdAt/1000, false)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting AOM service discovery rule fields: %s", err)
	}

	return nil
}

func GetServiceDiscoveryRule(client *golangsdk.ServiceClient, name string) (interface{}, error) {
	getHttpUrl := "v1/{project_id}/inv/servicediscoveryrules"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting the service discovery rule: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening the response: %s", err)
	}

	searchPath := fmt.Sprintf("appRules[?name=='%s']|[0]", name)
	rule := utils.PathSearch(searchPath, getRespBody, nil)
	if rule == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return rule, nil
}

func resourceServiceDiscoveryRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("aom", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteHttpUrl := "v1/{project_id}/inv/servicediscoveryrules?appRulesIds={rule_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Get("rule_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	// return 200 even deleting a non exist rule
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting the service discovery rule")
	}

	return nil
}

func flattenDiscoveryRules(rule []interface{}) []map[string]interface{} {
	var discoveryRules []map[string]interface{}
	for _, pairObject := range rule {
		discoveryRule := make(map[string]interface{})
		discoveryRule["check_content"] = utils.PathSearch("checkContent", pairObject, nil)
		discoveryRule["check_mode"] = utils.PathSearch("checkMode", pairObject, nil)
		discoveryRule["check_type"] = utils.PathSearch("checkType", pairObject, nil)

		discoveryRules = append(discoveryRules, discoveryRule)
	}
	return discoveryRules
}

func flattenBasicRules(rule []interface{}) []map[string]interface{} {
	var logPathRules []map[string]interface{}
	for _, pairObject := range rule {
		logPathRule := make(map[string]interface{})
		logPathRule["name_type"] = utils.PathSearch("nameType", pairObject, nil)
		logPathRule["args"] = utils.PathSearch("args", pairObject, nil)
		logPathRule["value"] = utils.PathSearch("value", pairObject, nil)

		logPathRules = append(logPathRules, logPathRule)
	}
	return logPathRules
}

func flattenNameRulesRules(rule interface{}) []interface{} {
	nameRule := map[string]interface{}{
		"service_name_rule": flattenBasicRules(
			utils.PathSearch("appNameRule", rule, make([]interface{}, 0)).([]interface{})),
		"application_name_rule": flattenBasicRules(
			utils.PathSearch("applicationNameRule", rule, make([]interface{}, 0)).([]interface{})),
	}
	return []interface{}{nameRule}
}

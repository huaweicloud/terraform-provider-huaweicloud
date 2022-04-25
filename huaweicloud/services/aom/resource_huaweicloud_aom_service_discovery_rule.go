package aom

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	aom "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"

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

func ResourceServiceDiscoveryRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceDiscoveryRuleCreateOrUpdate,
		ReadContext:   resourceServiceDiscoveryRuleRead,
		UpdateContext: resourceServiceDiscoveryRuleCreateOrUpdate,
		DeleteContext: resourceServiceDiscoveryRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				ValidateFunc: validation.All(
					validation.StringLenBetween(4, 63),
					validation.StringMatch(regexp.MustCompile("^[a-z]([a-z0-9-]*[a-z0-9])?$"),
						"The name can only consist of lowercase letters, digits and hyphens (-), "+
							"and it must start with a lowercase letter but cannot end with a hyphen (-)."),
				),
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
			"rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDiscoveryRuleOpts(rawRules []interface{}) []aom.DiscoveryRule {
	discoveryRules := make([]aom.DiscoveryRule, len(rawRules))
	for i, v := range rawRules {
		rawRule := v.(map[string]interface{})
		discoveryRules[i].CheckContent = utils.ExpandToStringList(rawRule["check_content"].([]interface{}))
		discoveryRules[i].CheckType = rawRule["check_type"].(string)
		discoveryRules[i].CheckMode = rawRule["check_mode"].(string)
	}
	return discoveryRules
}

func buildLogPathRuleOpts(rawRules []interface{}) *[]aom.LogPathRule {
	logPathRules := make([]aom.LogPathRule, len(rawRules))
	for i, v := range rawRules {
		rawRule := v.(map[string]interface{})
		logPathRules[i].Args = utils.ExpandToStringList(rawRule["args"].([]interface{}))
		logPathRules[i].NameType = rawRule["name_type"].(string)
		logPathRules[i].Value = utils.ExpandToStringList(rawRule["value"].([]interface{}))
	}
	return &logPathRules
}

func buildNameRuleOpts(rawRules []interface{}) *aom.NameRule {
	if len(rawRules) != 1 {
		return nil
	}
	raw := rawRules[0].(map[string]interface{})
	rawAppNameRule := raw["service_name_rule"].([]interface{})
	rawApplicationNameRule := raw["application_name_rule"].([]interface{})

	appNameRules := make([]aom.AppNameRule, len(rawAppNameRule))
	for i, v := range rawAppNameRule {
		rawRule := v.(map[string]interface{})
		appNameRules[i].Args = utils.ExpandToStringList(rawRule["args"].([]interface{}))
		appNameRules[i].NameType = rawRule["name_type"].(string)
		appNameRules[i].Value = utils.ExpandToStringListPointer(rawRule["value"].([]interface{}))
	}

	applicationNameRule := make([]aom.ApplicationNameRule, len(rawApplicationNameRule))
	for i, v := range rawApplicationNameRule {
		rawRule := v.(map[string]interface{})
		applicationNameRule[i].Args = utils.ExpandToStringList(rawRule["args"].([]interface{}))
		applicationNameRule[i].NameType = rawRule["name_type"].(string)
		applicationNameRule[i].Value = utils.ExpandToStringListPointer(rawRule["value"].([]interface{}))
	}

	nameRules := aom.NameRule{
		AppNameRule:         appNameRules,
		ApplicationNameRule: applicationNameRule,
	}

	return &nameRules
}

func FilterRules(allRules []aom.AppRules, name string) (*aom.AppRules, error) {
	for _, rule := range allRules {
		if rule.Name == name {
			return &rule, nil
		}
	}
	return nil, golangsdk.ErrDefault404{}
}

func resourceServiceDiscoveryRuleCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	createOpts := aom.AppRules{
		Id:        d.Get("rule_id").(string),
		Enable:    d.Get("discovery_rule_enabled").(bool),
		EventName: "aom_inventory_rules_event",
		Name:      d.Get("name").(string),
		Projectid: config.HwClient.ProjectID,
		Spec: &aom.AppRulesSpec{
			AppType:       d.Get("service_type").(string),
			DetectLog:     strconv.FormatBool(d.Get("detect_log_enabled").(bool)),
			DiscoveryRule: buildDiscoveryRuleOpts(d.Get("discovery_rules").([]interface{})),
			IsDefaultRule: strconv.FormatBool(d.Get("is_default_rule").(bool)),
			IsDetect:      "false",
			LogFileFix:    utils.ExpandToStringList(d.Get("log_file_suffix").([]interface{})),
			LogPathRule:   buildLogPathRuleOpts(d.Get("log_path_rules").([]interface{})),
			NameRule:      buildNameRuleOpts(d.Get("name_rules").([]interface{})),
			Priority:      int32(d.Get("priority").(int)),
		},
	}

	log.Printf("[DEBUG] Create or update %s Options: %#v", createOpts.Name, createOpts)

	createReq := aom.AddOrUpdateServiceDiscoveryRulesRequest{
		Body: &aom.AppRulesBody{
			AppRules: &[]aom.AppRules{createOpts},
		},
	}
	_, err = client.AddOrUpdateServiceDiscoveryRules(&createReq)
	if err != nil {
		return diag.Errorf("error creating or update AOM service discovery rule %s: %s", createOpts.Name, err)
	}

	d.SetId(createOpts.Name)

	return resourceServiceDiscoveryRuleRead(ctx, d, meta)
}

func resourceServiceDiscoveryRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	response, err := client.ListServiceDiscoveryRules(&aom.ListServiceDiscoveryRulesRequest{})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM service discovery rule")
	}

	allRules := *response.AppRules

	rule, err := FilterRules(allRules, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM service discovery rule")
	}

	log.Printf("[DEBUG] Retrieved AOM service discovery rule %s: %#v", d.Id(), rule)

	isDefaultRule, _ := strconv.ParseBool(rule.Spec.IsDefaultRule)
	detectLogEnabled, _ := strconv.ParseBool(rule.Spec.DetectLog)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", rule.Name),
		d.Set("rule_id", rule.Id),
		d.Set("discovery_rule_enabled", rule.Enable),
		d.Set("is_default_rule", isDefaultRule),
		d.Set("log_file_suffix", rule.Spec.LogFileFix),
		d.Set("service_type", rule.Spec.AppType),
		d.Set("detect_log_enabled", detectLogEnabled),
		d.Set("priority", rule.Spec.Priority),
		d.Set("discovery_rules", flattenDiscoveryRules(rule.Spec.DiscoveryRule)),
		d.Set("log_path_rules", flattenLogPathRulesRules(rule.Spec.LogPathRule)),
		d.Set("name_rules", flattenNameRulesRules(rule.Spec.NameRule)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting AOM service discovery rule fields: %s", err)
	}

	return nil
}

func resourceServiceDiscoveryRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	deleteReq := aom.DeleteserviceDiscoveryRulesRequest{
		AppRulesIds: []string{d.Get("rule_id").(string)},
	}
	_, err = client.DeleteserviceDiscoveryRules(&deleteReq)
	if err != nil {
		return diag.Errorf("error deleting AOM service discovery rule %s: %s", d.Id(), err)
	}
	return nil
}

func flattenDiscoveryRules(rule []aom.DiscoveryRule) []map[string]interface{} {
	var discoveryRules []map[string]interface{}
	for _, pairObject := range rule {
		discoveryRule := make(map[string]interface{})
		discoveryRule["check_content"] = pairObject.CheckContent
		discoveryRule["check_mode"] = pairObject.CheckMode
		discoveryRule["check_type"] = pairObject.CheckType

		discoveryRules = append(discoveryRules, discoveryRule)
	}
	return discoveryRules
}

func flattenLogPathRulesRules(rule *[]aom.LogPathRule) []map[string]interface{} {
	var logPathRules []map[string]interface{}
	for _, pairObject := range *rule {
		logPathRule := make(map[string]interface{})
		logPathRule["name_type"] = pairObject.NameType
		logPathRule["args"] = pairObject.Args
		logPathRule["value"] = pairObject.Value

		logPathRules = append(logPathRules, logPathRule)
	}
	return logPathRules
}

func flattenNameRulesRules(rule *aom.NameRule) []interface{} {
	var appNameRules []map[string]interface{}
	for _, pairObject := range rule.AppNameRule {
		appNameRule := make(map[string]interface{})
		appNameRule["name_type"] = pairObject.NameType
		appNameRule["args"] = pairObject.Args
		appNameRule["value"] = pairObject.Value

		appNameRules = append(appNameRules, appNameRule)
	}

	var applicationNameRules []map[string]interface{}
	for _, pairObject := range rule.AppNameRule {
		applicationNameRule := make(map[string]interface{})
		applicationNameRule["name_type"] = pairObject.NameType
		applicationNameRule["args"] = pairObject.Args
		applicationNameRule["value"] = pairObject.Value

		applicationNameRules = append(applicationNameRules, applicationNameRule)
	}
	nameRule := map[string]interface{}{
		"service_name_rule":     appNameRules,
		"application_name_rule": applicationNameRules,
	}
	return []interface{}{nameRule}
}

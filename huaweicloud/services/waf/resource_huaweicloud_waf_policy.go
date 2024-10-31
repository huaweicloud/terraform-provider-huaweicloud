package waf

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}
// @API WAF PATCH /v1/{project_id}/waf/policy/{policy_id}
// @API WAF POST /v1/{project_id}/waf/policy
func ResourceWafPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafPolicyCreate,
		ReadContext:   resourceWafPolicyRead,
		UpdateContext: resourceWafPolicyUpdate,
		DeleteContext: resourceWafPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"full_detection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"protection_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"robot_action": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"level": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"deep_inspection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"header_inspection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"shiro_decryption_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     policyOptionSchema(),
			},
			"bind_hosts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     policyBindHostSchema(),
			},
		},
	}
}

func policyOptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"basic_web_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"general_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"crawler_engine": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"crawler_scanner": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"crawler_script": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"crawler_other": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"webshell": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cc_attack_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"precise_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"blacklist": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"data_masking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"false_alarm_masking": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"web_tamper_protection": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"geolocation_access_control": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"information_leakage_prevention": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"bot_enable": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"known_attack_source": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"anti_crawler": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"crawler": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "schema: Deprecated",
			},
		},
	}
	return &sc
}

func policyBindHostSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"waf_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceWafPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}
	policy, err := policies.Create(wafClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF policy: %s", err)
	}
	d.SetId(policy.Id)

	// Update policy
	if err := updatePolicy(wafClient, d, cfg); err != nil {
		return diag.FromErr(err)
	}

	return resourceWafPolicyRead(ctx, d, meta)
}

func updatePolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, cfg *config.Config) error {
	updateOpts := policies.UpdateOpts{
		Name:                d.Get("name").(string),
		Level:               d.Get("level").(int),
		FullDetection:       utils.Bool(d.Get("full_detection").(bool)),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Options:             buildUpdatePolicyOption(d),
		Extend:              buildExtend(d),
	}

	if v, ok := d.GetOk("protection_mode"); ok {
		updateOpts.Action = &policies.Action{
			Category: v.(string),
		}
	}

	if v, ok := d.GetOk("robot_action"); ok {
		updateOpts.RobotAction = &policies.Action{
			Category: v.(string),
		}
	}
	_, err := policies.Update(client, d.Id(), updateOpts).Extract()
	return err
}

// Due to API reasons, these three fields need to be edited through the `extend` structure.
func buildExtend(d *schema.ResourceData) map[string]string {
	extendObj := map[string]interface{}{
		"deep_decode":             d.Get("deep_inspection"),
		"check_all_headers":       d.Get("header_inspection"),
		"shiro_rememberMe_enable": d.Get("shiro_decryption_check"),
	}

	extendJsonString := utils.JsonToString(extendObj)
	if extendJsonString == "" {
		return nil
	}

	return map[string]string{
		"extend": extendJsonString,
	}
}

func buildUpdatePolicyOption(d *schema.ResourceData) *policies.PolicyOption {
	// if not specified, make all the switch off
	defaultOptionMap := make(map[string]bool)
	if rawArray, ok := d.Get("options").([]interface{}); ok && len(rawArray) > 0 {
		if rawMap, rawOk := rawArray[0].(map[string]interface{}); rawOk {
			for k, v := range rawMap {
				defaultOptionMap[k] = v.(bool)
			}
		}
	}

	return &policies.PolicyOption{
		Webattack:      utils.Bool(defaultOptionMap["basic_web_protection"]),
		Common:         utils.Bool(defaultOptionMap["general_check"]),
		CrawlerEngine:  utils.Bool(defaultOptionMap["crawler_engine"]),
		CrawlerScanner: utils.Bool(defaultOptionMap["crawler_scanner"]),
		CrawlerScript:  utils.Bool(defaultOptionMap["crawler_script"]),
		CrawlerOther:   utils.Bool(defaultOptionMap["crawler_other"]),
		Webshell:       utils.Bool(defaultOptionMap["webshell"]),
		Cc:             utils.Bool(defaultOptionMap["cc_attack_protection"]),
		Custom:         utils.Bool(defaultOptionMap["precise_protection"]),
		Whiteblackip:   utils.Bool(defaultOptionMap["blacklist"]),
		Ignore:         utils.Bool(defaultOptionMap["false_alarm_masking"]),
		Privacy:        utils.Bool(defaultOptionMap["data_masking"]),
		Antitamper:     utils.Bool(defaultOptionMap["web_tamper_protection"]),
		GeoIP:          utils.Bool(defaultOptionMap["geolocation_access_control"]),
		Antileakage:    utils.Bool(defaultOptionMap["information_leakage_prevention"]),
		BotEnable:      utils.Bool(defaultOptionMap["bot_enable"]),
		FollowedAction: utils.Bool(defaultOptionMap["known_attack_source"]),
		Anticrawler:    utils.Bool(defaultOptionMap["anti_crawler"]),
	}
}

func resourceWafPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	n, err := policies.GetWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		// If the policy does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF policy")
	}

	// the extend struct value example is: "extend": "{\"deep_decode\":true}"
	var extendRespBody interface{}
	if extendValue, ok := n.Extend["extend"]; ok {
		if err := json.Unmarshal([]byte(extendValue), &extendRespBody); err != nil {
			log.Printf("[WARN] error flatten extend map: %s", err)
		}
	}
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("full_detection", n.FullDetection),
		d.Set("protection_mode", n.Action.Category),
		d.Set("level", n.Level),
		d.Set("robot_action", n.RobotAction.Category),
		d.Set("options", flattenOptions(n.Options)),
		d.Set("bind_hosts", flattenBindHosts(n.BindHosts)),
		d.Set("deep_inspection", utils.PathSearch("deep_decode", extendRespBody, false)),
		d.Set("header_inspection", utils.PathSearch("check_all_headers", extendRespBody, false)),
		d.Set("shiro_decryption_check", utils.PathSearch("shiro_rememberMe_enable", extendRespBody, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOptions(policyOption policies.PolicyOption) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"basic_web_protection":           policyOption.Webattack,
			"general_check":                  policyOption.Common,
			"crawler_engine":                 policyOption.CrawlerEngine,
			"crawler_scanner":                policyOption.CrawlerScanner,
			"crawler_script":                 policyOption.CrawlerScript,
			"crawler_other":                  policyOption.CrawlerOther,
			"webshell":                       policyOption.Webshell,
			"cc_attack_protection":           policyOption.Cc,
			"precise_protection":             policyOption.Custom,
			"blacklist":                      policyOption.Whiteblackip,
			"false_alarm_masking":            policyOption.Ignore,
			"data_masking":                   policyOption.Privacy,
			"web_tamper_protection":          policyOption.Antitamper,
			"geolocation_access_control":     policyOption.GeoIP,
			"information_leakage_prevention": policyOption.Antileakage,
			"bot_enable":                     policyOption.BotEnable,
			"known_attack_source":            policyOption.FollowedAction,
			"anti_crawler":                   policyOption.Anticrawler,
			"crawler":                        policyOption.Crawler,
		},
	}
}

func flattenBindHosts(bindHosts []policies.BindHost) []map[string]interface{} {
	rst := make([]map[string]interface{}, len(bindHosts))
	for i, host := range bindHosts {
		rst[i] = map[string]interface{}{
			"id":       host.Id,
			"hostname": host.Hostname,
			"waf_type": host.WafType,
			"mode":     host.Mode,
		}
	}
	return rst
}

func resourceWafPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if err := updatePolicy(wafClient, d, cfg); err != nil {
		return diag.Errorf("error updating WAF policy: %s", err)
	}
	return resourceWafPolicyRead(ctx, d, meta)
}

func resourceWafPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	err = policies.DeleteWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).ExtractErr()
	if err != nil {
		// If the policy does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF policy")
	}
	return nil
}

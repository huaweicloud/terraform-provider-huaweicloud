package waf

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	ProtectionModeLog   = "log"
	ProtectionModeBlock = "block"
)

func ResourceWafPolicyV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafPolicyV1Create,
		ReadContext:   resourceWafPolicyV1Read,
		UpdateContext: resourceWafPolicyV1Update,
		DeleteContext: resourceWafPolicyV1Delete,
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
			"protection_mode": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					ProtectionModeLog, ProtectionModeBlock,
				}, false),
			},
			"level": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"options": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"basic_web_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"general_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_engine": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_scanner": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_script": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"crawler_other": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"webshell": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"cc_attack_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"precise_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"blacklist": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"data_masking": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"false_alarm_masking": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"web_tamper_protection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"full_detection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceWafPolicyV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	d.Set("name", policy.Name)

	level := d.Get("level").(int)
	protectionMode := d.Get("protection_mode").(string)

	// Get the policy details, then check if need to update 'protection_mode' or 'level.
	resourceWafPolicyV1Read(ctx, d, meta)

	// If the vlaue of 'protection_mode' or 'level' is not equal to the value returned by the server,
	// then update the policy.
	err = checkAndUpdateDefaultVal(wafClient, d, protectionMode, level, cfg)
	if err != nil {
		return diag.Errorf("the Waf Policy was created successfully, "+
			"but failed to update protection_mode or level : %s", err)
	} else {
		d.Set("protection_mode", protectionMode)
		d.Set("level", level)
	}

	return resourceWafPolicyV1Read(ctx, d, meta)
}

// checkAndUpdateDefaultVal check the vlaue of 'protection_mode' or 'level' is not equal to
// the value returned by the server.
// If the value is not equal to the value returned by the server, call the update API to make changes.
func checkAndUpdateDefaultVal(wafClient *golangsdk.ServiceClient, d *schema.ResourceData,
	protectionMode string, level int, conf *config.Config) error {
	needUpdate := false
	updateOpts := policies.UpdateOpts{
		EnterpriseProjectId: common.GetEnterpriseProjectID(d, conf),
	}

	if strings.Compare(d.Get("protection_mode").(string), protectionMode) != 0 && len(protectionMode) != 0 {
		needUpdate = true
		updateOpts.Action = &policies.Action{
			Category: protectionMode,
		}
	}
	if d.Get("level").(int) != level && level != 0 {
		needUpdate = true
		updateOpts.Level = level
	}

	if needUpdate {
		_, err := policies.Update(wafClient, d.Id(), updateOpts).Extract()
		return err
	}
	return nil
}

func resourceWafPolicyV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	n, err := policies.GetWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF policy")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("level", n.Level),
		d.Set("protection_mode", n.Action.Category),
		d.Set("full_detection", n.FullDetection),
		d.Set("options", flattenOptions(n)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOptions(policy *policies.Policy) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"basic_web_protection":  policy.Options.Webattack,
			"general_check":         policy.Options.Common,
			"crawler":               policy.Options.Crawler,
			"crawler_engine":        policy.Options.CrawlerEngine,
			"crawler_scanner":       policy.Options.CrawlerScanner,
			"crawler_script":        policy.Options.CrawlerScript,
			"crawler_other":         policy.Options.CrawlerOther,
			"webshell":              policy.Options.Webshell,
			"cc_attack_protection":  policy.Options.Cc,
			"precise_protection":    policy.Options.Custom,
			"blacklist":             policy.Options.Whiteblackip,
			"false_alarm_masking":   policy.Options.Ignore,
			"data_masking":          policy.Options.Privacy,
			"web_tamper_protection": policy.Options.Antitamper,
		},
	}
}

func resourceWafPolicyV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChanges("name", "level", "protection_mode") {
		updateOpts := policies.UpdateOpts{
			Name:  d.Get("name").(string),
			Level: d.Get("level").(int),
			Action: &policies.Action{
				Category: d.Get("protection_mode").(string),
			},
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}

		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF policy: %s", err)
		}
	}
	return resourceWafPolicyV1Read(ctx, d, meta)
}

func resourceWafPolicyV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	err = policies.DeleteWithEpsID(wafClient, d.Id(), cfg.GetEnterpriseProjectID(d)).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting WAF policy: %s", err)
	}
	return nil
}

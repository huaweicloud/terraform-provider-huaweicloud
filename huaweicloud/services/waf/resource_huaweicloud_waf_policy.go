package waf

import (
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	PROTECTION_MODE_LOG   = "log"
	PROTECTION_MODE_BLOCK = "block"
)

func ResourceWafPolicyV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafPolicyV1Create,
		Read:   resourceWafPolicyV1Read,
		Update: resourceWafPolicyV1Update,
		Delete: resourceWafPolicyV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
					PROTECTION_MODE_LOG, PROTECTION_MODE_BLOCK,
				}, false),
			},
			"level": {
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 3),
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

func resourceWafPolicyV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud WAF client: %s", err)
	}

	createOpts := policies.CreateOpts{
		Name: d.Get("name").(string),
	}
	policy, err := policies.Create(wafClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating waf policy: %s", err)
	}

	logp.Printf("[DEBUG] Waf policy created: %#v", policy)
	d.SetId(policy.Id)
	d.Set("name", policy.Name)

	level := d.Get("level").(int)
	protectionMode := d.Get("protection_mode").(string)

	// Get the policy details, then check if need to update 'protection_mode' or 'level.
	err = resourceWafPolicyV1Read(d, meta)

	// If the vlaue of 'protection_mode' or 'level' is not equal to the value returned by the server,
	// then update the policy.
	err = checkAndUpdateDefaultVal(wafClient, d, protectionMode, level)
	if err != nil {
		return fmtp.Errorf("the Waf Policy was created successfully, "+
			"but failed to update protection_mode or level : %s", err)
	} else {
		d.Set("protection_mode", protectionMode)
		d.Set("level", level)
	}

	return err
}

// checkAndUpdateDefaultVal check the vlaue of 'protection_mode' or 'level' is not equal to
// the value returned by the server.
// If the value is not equal to the value returned by the server, call the update API to make changes.
func checkAndUpdateDefaultVal(wafClient *golangsdk.ServiceClient, d *schema.ResourceData,
	protectionMode string, level int) error {

	needUpdate := false
	updateOpts := policies.UpdateOpts{}

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
		logp.Printf("[DEBUG] update default protection_mode or level: %#v", updateOpts)
		_, err := policies.Update(wafClient, d.Id(), updateOpts).Extract()
		return err
	}
	return nil
}

func resourceWafPolicyV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	n, err := policies.Get(wafClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Waf Policy")
	}

	d.Set("region", config.GetRegion(d))
	d.Set("name", n.Name)
	d.Set("level", n.Level)
	d.Set("protection_mode", n.Action.Category)
	d.Set("full_detection", n.FullDetection)

	options := []map[string]interface{}{
		{
			"basic_web_protection":  n.Options.Webattack,
			"general_check":         n.Options.Common,
			"crawler":               n.Options.Crawler,
			"crawler_engine":        n.Options.CrawlerEngine,
			"crawler_scanner":       n.Options.CrawlerScanner,
			"crawler_script":        n.Options.CrawlerScript,
			"crawler_other":         n.Options.CrawlerOther,
			"webshell":              n.Options.Webshell,
			"cc_attack_protection":  n.Options.Cc,
			"precise_protection":    n.Options.Custom,
			"blacklist":             n.Options.Whiteblackip,
			"false_alarm_masking":   n.Options.Ignore,
			"data_masking":          n.Options.Privacy,
			"web_tamper_protection": n.Options.Antitamper,
		},
	}
	d.Set("options", options)
	return nil
}

func resourceWafPolicyV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	if d.HasChanges("name", "level", "protection_mode") {
		updateOpts := policies.UpdateOpts{
			Name:  d.Get("name").(string),
			Level: d.Get("level").(int),
			Action: &policies.Action{
				Category: d.Get("protection_mode").(string),
			},
		}

		logp.Printf("[DEBUG] updateOpts: %#v", updateOpts)
		_, err = policies.Update(wafClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating WAF Policy: %s", err)
		}
	}
	return resourceWafPolicyV1Read(d, meta)
}

func resourceWafPolicyV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	err = policies.Delete(wafClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting WAF Policy: %s", err)
	}

	d.SetId("")
	return nil
}

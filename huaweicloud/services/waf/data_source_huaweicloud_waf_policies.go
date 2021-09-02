package waf

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// DataSourceWafPoliciesV1 the function is used for data source 'huaweicloud_waf_policies'.
func DataSourceWafPoliciesV1() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceWafPoliciesV1Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeInt,
							Computed: true,
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
				},
			},
		},
	}
}

func dataSourceWafPoliciesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud WAF client: %s", err)
	}

	listOpts := policies.ListPolicyOpts{
		Name: d.Get("name").(string),
	}
	rst, err := policies.ListPolicy(wafClient, listOpts)
	if err != nil {
		return fmtp.Errorf("Unable to retrieve waf policies: %s", err)
	}
	logp.Printf("[DEBUG] Get a list of policies: %#v.", rst)

	if len(rst.Items) == 0 {
		return fmtp.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	// The IDs of all policies in the list, and get it hashcode set to the value of schema id.
	ids := make([]string, 0, len(rst.Items))
	policies := make([]map[string]interface{}, 0, len(rst.Items))

	for _, p := range rst.Items {
		options := []map[string]interface{}{
			{
				"basic_web_protection":  p.Options.Webattack,
				"general_check":         p.Options.Common,
				"crawler":               p.Options.Crawler,
				"crawler_engine":        p.Options.CrawlerEngine,
				"crawler_scanner":       p.Options.CrawlerScanner,
				"crawler_script":        p.Options.CrawlerScript,
				"crawler_other":         p.Options.CrawlerOther,
				"webshell":              p.Options.Webshell,
				"cc_attack_protection":  p.Options.Cc,
				"precise_protection":    p.Options.Custom,
				"blacklist":             p.Options.Whiteblackip,
				"false_alarm_masking":   p.Options.Ignore,
				"data_masking":          p.Options.Privacy,
				"web_tamper_protection": p.Options.Antitamper,
			},
		}
		plc := map[string]interface{}{
			"id":              p.Id,
			"name":            p.Name,
			"protection_mode": p.Action.Category,
			"level":           p.Level,
			"options":         options,
			"full_detection":  p.FullDetection,
		}
		policies = append(policies, plc)
		ids = append(ids, p.Id)
	}

	d.SetId(hashcode.Strings(ids))
	if err = d.Set("policies", policies); err != nil {
		return fmtp.Errorf("error setting WAF policy fields: %s", err)
	}

	return nil
}

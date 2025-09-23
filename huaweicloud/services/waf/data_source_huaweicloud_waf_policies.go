package waf

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf_hw/v1/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/policy
func DataSourceWafPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafPoliciesRead,

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
			"enterprise_project_id": {
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
						"full_detection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"protection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"robot_action": {
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
							Elem:     dataSourcePolicyOptionSchema(),
						},
						"bind_hosts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     policyBindHostSchema(),
						},
						"deep_inspection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"header_inspection": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"shiro_decryption_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourcePolicyOptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"basic_web_protection": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"general_check": {
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
			"geolocation_access_control": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"information_leakage_prevention": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"bot_enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"known_attack_source": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"anti_crawler": {
				Type:     schema.TypeBool,
				Computed: true,
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

func dataSourceWafPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	listOpts := policies.ListPolicyOpts{
		Name:                d.Get("name").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		Pagesize:            100,
	}
	policyList, err := policies.List(wafClient, listOpts)
	if err != nil {
		return diag.Errorf("error retrieving WAF policies, %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("policies", flattenListPolicies(policyList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPolicies(policyList []policies.Policy) []interface{} {
	rst := make([]interface{}, len(policyList))
	for i, n := range policyList {
		// the extend struct value example is: "extend": "{\"deep_decode\":true}"
		var extendRespBody interface{}
		if extendValue, ok := n.Extend["extend"]; ok {
			if err := json.Unmarshal([]byte(extendValue), &extendRespBody); err != nil {
				log.Printf("[WARN] error flatten extend map: %s", err)
			}
		}

		rst[i] = map[string]interface{}{
			"id":                     n.Id,
			"name":                   n.Name,
			"full_detection":         n.FullDetection,
			"protection_mode":        n.Action.Category,
			"level":                  n.Level,
			"robot_action":           n.RobotAction.Category,
			"options":                flattenOptions(n.Options),
			"bind_hosts":             flattenBindHosts(n.BindHosts),
			"deep_inspection":        utils.PathSearch("deep_decode", extendRespBody, false),
			"header_inspection":      utils.PathSearch("check_all_headers", extendRespBody, false),
			"shiro_decryption_check": utils.PathSearch("shiro_rememberMe_enable", extendRespBody, false),
		}
	}
	return rst
}

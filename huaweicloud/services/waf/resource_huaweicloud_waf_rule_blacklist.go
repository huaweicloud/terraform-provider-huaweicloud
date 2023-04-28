package waf

import (
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/whiteblackip_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const (
	// PROTECTION_ACTION_BLOCK block the request
	PROTECTION_ACTION_BLOCK = 0
	// PROTECTION_ACTION_ALLOW allow the request
	PROTECTION_ACTION_ALLOW = 1
	// PROTECTION_ACTION_LOG log the request only
	PROTECTION_ACTION_LOG = 2
)

func ResourceWafRuleBlackListV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleBlackListCreate,
		Read:   resourceWafRuleBlackListRead,
		Update: resourceWafRuleBlackListUpdate,
		Delete: resourceWafRuleBlackListDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Required",
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"ip_address", "address_group_id"},
			},
			"address_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  PROTECTION_ACTION_BLOCK,
				ValidateFunc: validation.IntInSlice([]int{
					PROTECTION_ACTION_BLOCK, PROTECTION_ACTION_ALLOW, PROTECTION_ACTION_LOG,
				}),
			},
			"address_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"address_group_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceWafRuleBlackListCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	epsID := config.GetEnterpriseProjectID(d)
	createOpts := rules.CreateOpts{
		White: d.Get("action").(int),
	}
	if v, ok := d.GetOk("name"); ok {
		createOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("ip_address"); ok {
		createOpts.Addr = v.(string)
	}
	if v, ok := d.GetOk("address_group_id"); ok {
		createOpts.IPGroupID = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		createOpts.Description = v.(string)
	}

	logp.Printf("[DEBUG] WAF black list rule creating opts: %#v", createOpts)
	rule, err := rules.CreateWithEpsId(wafClient, createOpts, policyID, epsID).Extract()
	if err != nil {
		return fmtp.Errorf("error creating WAF black list rule: %s", err)
	}
	logp.Printf("[DEBUG] WAF black list rule created: %#v", rule)
	// After the creation is successful, set the value id of the schema.
	d.SetId(rule.Id)

	return resourceWafRuleBlackListRead(d, meta)
}

func resourceWafRuleBlackListRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.GetWithEpsId(wafClient, policyID, d.Id(), config.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "WAF Black List Rule")
	}
	logp.Printf("[DEBUG] fetching WAF black list rule: %#v", n)

	d.SetId(n.Id)
	mErr := multierror.Append(nil,
		d.Set("policy_id", n.PolicyID),
		d.Set("name", n.Name),
		d.Set("ip_address", n.Addr),
		d.Set("description", n.Description),
		d.Set("action", n.White),
	)

	ipGroup := n.IPGroup
	if ipGroup != nil {
		mErr = multierror.Append(mErr,
			d.Set("address_group_id", ipGroup.ID),
			d.Set("address_group_name", ipGroup.Name),
			d.Set("address_group_size", ipGroup.Size),
		)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("address_group_id", ""),
			d.Set("address_group_name", ""),
			d.Set("address_group_size", 0),
		)
	}
	if mErr.ErrorOrNil() != nil {
		return fmtp.Errorf("error setting WAF rule blacklist fields: %s", mErr)
	}
	return nil
}

func resourceWafRuleBlackListUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	if d.HasChanges("name", "ip_address", "address_group_id", "description", "action") {
		updateOpts := rules.UpdateOpts{
			White:       utils.Int(d.Get("action").(int)),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		if v, ok := d.GetOk("ip_address"); ok {
			updateOpts.Addr = v.(string)
		}
		if v, ok := d.GetOk("address_group_id"); ok {
			updateOpts.IPGroupID = v.(string)
		}

		logp.Printf("[DEBUG] updating blacklist and whitelist rule, updateOpts: %#v", updateOpts)

		policyID := d.Get("policy_id").(string)
		epsID := config.GetEnterpriseProjectID(d)
		_, err = rules.UpdateWithEpsId(wafClient, updateOpts, policyID, d.Id(), epsID).Extract()
		if err != nil {
			return fmtp.Errorf("error updating HuaweiCloud WAF Blacklist and Whitelist Rule: %s", err)
		}
	}

	return resourceWafRuleBlackListRead(d, meta)
}

func resourceWafRuleBlackListDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteWithEpsId(wafClient, policyID, d.Id(), config.GetEnterpriseProjectID(d)).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting HuaweiCloud WAF Blacklist and Whitelist Rule: %s", err)
	}

	d.SetId("")
	return nil
}

// resourceWafRulesImport query the rules from HuaweiCloud and imports them to Terraform.
// It is a common function in waf and is also called by other rule resources.
func resourceWafRulesImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for WAF rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	policyID := parts[0]
	ruleID := parts[1]

	d.SetId(ruleID)
	d.Set("policy_id", policyID)

	return []*schema.ResourceData{d}, nil
}

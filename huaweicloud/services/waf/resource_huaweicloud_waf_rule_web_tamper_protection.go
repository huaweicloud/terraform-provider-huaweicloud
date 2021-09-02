/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"github.com/hashicorp/go-multierror"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/webtamperprotection_rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceWafRuleWebTamperProtectionV1 manages the resources for web tamper protection rules
func ResourceWafRuleWebTamperProtectionV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleWebTamperProtectionCreate,
		Read:   resourceWafRuleWebTamperProtectionRead,
		Delete: resourceWafRuleWebTamperProtectionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceWafRulesImport,
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
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

// resourceWafRuleWebTamperProtectionCreate create rules
func resourceWafRuleWebTamperProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}
	// create options
	createOpts := rules.CreateOpts{
		Hostname: d.Get("domain").(string),
		Url:      d.Get("path").(string),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Web Tamper Protection Rule: %s", err)
	}

	logp.Printf("[DEBUG] WAF web tamper protection rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleWebTamperProtectionRead(d, meta)
}

// resourceWafRuleWebTamperProtectionRead read rules from HuaweiCloud by id and policyid
func resourceWafRuleWebTamperProtectionRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "WAF Web Tamper Protection Rule")
	}

	d.SetId(n.Id)
	mErr := multierror.Append(nil,
		d.Set("policy_id", n.PolicyID),
		d.Set("domain", n.Hostname),
		d.Set("path", n.Url),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

// resourceWafRuleWebTamperProtectionDelete delete the rules from HuaweiCloud by id
func resourceWafRuleWebTamperProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting HuaweiCloud WAF Web Tamper Protection Rule: %s", err)
	}

	d.SetId("")
	return nil
}

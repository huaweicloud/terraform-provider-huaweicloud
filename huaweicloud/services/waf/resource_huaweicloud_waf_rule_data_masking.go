/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/datamasking_rules"
)

const (
	FIELD_POSITION_HEADER = "header"
	FIELD_POSITION_PARAMS = "params"
	FIELD_POSITION_COOKIE = "cookie"
	FIELD_POSITION_FORM   = "form"
)

// ResourceWafRuleDataMaskingV1 the resource of managing a WAF Data Masking Rule within HuaweiCloud.
func ResourceWafRuleDataMaskingV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceWafRuleDataMaskingCreate,
		Read:   resourceWafRuleDataMaskingRead,
		Update: resourceWafRuleDataMaskingUpdate,
		Delete: resourceWafRuleDataMaskingDelete,
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
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					FIELD_POSITION_HEADER, FIELD_POSITION_PARAMS, FIELD_POSITION_COOKIE, FIELD_POSITION_FORM,
				}, false),
			},
			"subfield": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// resourceWafRuleDataMaskingCreate create a rule
func resourceWafRuleDataMaskingCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	createOpts := rules.CreateOpts{
		Path:     d.Get("path").(string),
		Category: d.Get("field").(string),
		Index:    d.Get("subfield").(string),
	}

	logp.Printf("[DEBUG] WAF Data Masking Rule creating opts: %#v", createOpts)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Data Masking Rule: %s", err)
	}

	logp.Printf("[DEBUG] WAF data masking rule created: %#v", rule)
	d.SetId(rule.Id)

	return resourceWafRuleDataMaskingRead(d, meta)
}

// resourceWafRuleDataMaskingRead get rule detail from HuaweiCloud by id and policy_id
func resourceWafRuleDataMaskingRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.Get(wafClient, policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "WAF Data Masking Rule")
	}
	logp.Printf("[DEBUG] fetching WAF data masking rule: %#v", n)

	d.SetId(n.Id)
	d.Set("path", n.Path)
	d.Set("field", n.Category)
	d.Set("subfield", n.Index)

	return nil
}

// resourceWafRuleDataMaskingUpdate update the existing rules.
// Supported fields: path, field, subfield
func resourceWafRuleDataMaskingUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF Client: %s", err)
	}

	if d.HasChanges("path", "field", "subfield") {
		policyID := d.Get("policy_id").(string)
		updateOpts := rules.UpdateOpts{
			Path:     d.Get("path").(string),
			Category: d.Get("field").(string),
			Index:    d.Get("subfield").(string),
		}

		logp.Printf("[DEBUG] WAF Data Masking Rule updating opts: %#v", updateOpts)
		_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("error updating HuaweiCloud WAF Data Masking Rule: %s", err)
		}
	}

	return resourceWafRuleDataMaskingRead(d, meta)
}

// resourceWafRuleDataMaskingDelete delete the rules from HuaweiCloud by id
func resourceWafRuleDataMaskingDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	wafClient, err := config.WafV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.Delete(wafClient, policyID, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("error deleting HuaweiCloud WAF Data Masking Rule: %s", err)
	}

	d.SetId("")
	return nil
}

/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/webtamperprotection_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceWafRuleWebTamperProtectionV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafRuleWebTamperProtectionCreate,
		ReadContext:   resourceWafRuleWebTamperProtectionRead,
		DeleteContext: resourceWafRuleWebTamperProtectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFRuleImportState,
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
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceWafRuleWebTamperProtectionCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createOpts := rules.CreateOpts{
		Hostname:            d.Get("domain").(string),
		Url:                 d.Get("path").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF web tamper protection rule: %s", err)
	}
	d.SetId(rule.Id)

	return resourceWafRuleWebTamperProtectionRead(ctx, d, meta)
}

func resourceWafRuleWebTamperProtectionRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	epsID := cfg.GetEnterpriseProjectID(d)
	n, err := rules.GetWithEpsID(wafClient, policyID, d.Id(), epsID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving WAF web tamper protection rule")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_id", n.PolicyID),
		d.Set("domain", n.Hostname),
		d.Set("path", n.Url),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafRuleWebTamperProtectionDelete(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	epsID := cfg.GetEnterpriseProjectID(d)
	err = rules.DeleteWithEpsID(wafClient, policyID, d.Id(), epsID).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting WAF web tamper protection rule: %s", err)
	}
	return nil
}

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

// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/antitamper
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/antitamper/{rule_id}
// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/antitamper/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
func ResourceWafRuleWebTamperProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafRuleWebTamperProtectionCreate,
		UpdateContext: resourceWafRuleWebTamperProtectionUpdate,
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
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
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
		Description:         d.Get("description").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	policyID := d.Get("policy_id").(string)
	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF web tamper protection rule: %s", err)
	}
	d.SetId(rule.Id)

	if d.Get("status").(int) == 0 {
		if err := updateRuleStatus(wafClient, d, cfg, "antitamper"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceWafRuleWebTamperProtectionRead(ctx, d, meta)
}

func resourceWafRuleWebTamperProtectionUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(wafClient, d, cfg, "antitamper"); err != nil {
			return diag.FromErr(err)
		}
	}
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
		// If the web tamper protection rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF web tamper protection rule")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_id", n.PolicyID),
		d.Set("domain", n.Hostname),
		d.Set("path", n.Url),
		d.Set("description", n.Description),
		d.Set("status", n.Status),
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
		// If the web tamper protection rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF web tamper protection rule")
	}
	return nil
}

/*
 Copyright (c) Huawei Technologies Co., Ltd. 2021. All rights reserved.
*/

package waf

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/datamasking_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

const (
	fieldPositionHeader = "header"
	fieldPositionParams = "params"
	fieldPositionCookie = "cookie"
	fieldPositionForm   = "form"
)

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/privacy/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/privacy/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/privacy/{rule_id}
// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/privacy
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
func ResourceWafRuleDataMasking() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafRuleDataMaskingCreate,
		ReadContext:   resourceWafRuleDataMaskingRead,
		UpdateContext: resourceWafRuleDataMaskingUpdate,
		DeleteContext: resourceWafRuleDataMaskingDelete,
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
			"path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"field": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					fieldPositionHeader, fieldPositionParams, fieldPositionCookie, fieldPositionForm,
				}, false),
			},
			"subfield": {
				Type:     schema.TypeString,
				Required: true,
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
				Computed: true,
			},
		},
	}
}

func resourceWafRuleDataMaskingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	createOpts := rules.CreateOpts{
		Path:                d.Get("path").(string),
		Category:            d.Get("field").(string),
		Index:               d.Get("subfield").(string),
		Description:         d.Get("description").(string),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	rule, err := rules.Create(wafClient, policyID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF data masking rule: %s", err)
	}
	d.SetId(rule.Id)

	if d.Get("status").(int) == 0 {
		if err := updateRuleStatus(wafClient, d, cfg, "privacy"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceWafRuleDataMaskingRead(ctx, d, meta)
}

func resourceWafRuleDataMaskingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	epsID := cfg.GetEnterpriseProjectID(d)
	n, err := rules.GetWithEpsID(wafClient, policyID, d.Id(), epsID).Extract()
	if err != nil {
		// If the data masking rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF data masking rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_id", n.PolicyID),
		d.Set("path", n.Path),
		d.Set("field", n.Category),
		d.Set("subfield", n.Index),
		d.Set("description", n.Description),
		d.Set("status", n.Status),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafRuleDataMaskingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChanges("path", "field", "subfield", "description") {
		policyID := d.Get("policy_id").(string)
		updateOpts := rules.UpdateOpts{
			Path:                d.Get("path").(string),
			Category:            d.Get("field").(string),
			Index:               d.Get("subfield").(string),
			Description:         d.Get("description").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}

		_, err = rules.Update(wafClient, policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF data masking rule: %s", err)
		}
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(wafClient, d, cfg, "privacy"); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceWafRuleDataMaskingRead(ctx, d, meta)
}

func resourceWafRuleDataMaskingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteWithEpsID(wafClient, policyID, d.Id(), cfg.GetEnterpriseProjectID(d)).ExtractErr()
	if err != nil {
		// If the data masking rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF data masking rule")
	}
	return nil
}

package waf

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	rules "github.com/chnsz/golangsdk/openstack/waf_hw/v1/whiteblackip_rules"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF DELETE /v1/{project_id}/waf/policy/{policy_id}/whiteblackip/{rule_id}
// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/whiteblackip/{rule_id}
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/whiteblackip/{rule_id}
// @API WAF POST /v1/{project_id}/waf/policy/{policy_id}/whiteblackip
// @API WAF PUT /v1/{project_id}/waf/policy/{policy_id}/{rule_type}/{rule_id}/status
func ResourceWafRuleBlackList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafRuleBlackListCreate,
		ReadContext:   resourceWafRuleBlackListRead,
		UpdateContext: resourceWafRuleBlackListUpdate,
		DeleteContext: resourceWafRuleBlackListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFRuleImportState,
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
			},
			"status": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1,
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

func resourceWafRuleBlackListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	epsID := cfg.GetEnterpriseProjectID(d)
	createOpts := rules.CreateOpts{
		White:       d.Get("action").(int),
		Name:        d.Get("name").(string),
		Addr:        d.Get("ip_address").(string),
		IPGroupID:   d.Get("address_group_id").(string),
		Description: d.Get("description").(string),
	}

	rule, err := rules.CreateWithEpsId(wafClient, createOpts, policyID, epsID).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF blacklist and whitelist rule: %s", err)
	}
	d.SetId(rule.Id)

	if d.Get("status").(int) == 0 {
		if err := updateRuleStatus(wafClient, d, cfg, "whiteblackip"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceWafRuleBlackListRead(ctx, d, meta)
}

func resourceWafRuleBlackListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	wafClient, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	n, err := rules.GetWithEpsId(wafClient, policyID, d.Id(), cfg.GetEnterpriseProjectID(d)).Extract()
	if err != nil {
		// If the blacklist and whitelist rule does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF blacklist and whitelist rule")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policy_id", n.PolicyID),
		d.Set("name", n.Name),
		d.Set("ip_address", n.Addr),
		d.Set("description", n.Description),
		d.Set("action", n.White),
		d.Set("status", n.Status),
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
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWafRuleBlackListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	updateChanges := []string{
		"name",
		"ip_address",
		"address_group_id",
		"description",
		"action",
	}

	if d.HasChanges(updateChanges...) {
		updateOpts := rules.UpdateOpts{
			White:       utils.Int(d.Get("action").(int)),
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Addr:        d.Get("ip_address").(string),
			IPGroupID:   d.Get("address_group_id").(string),
		}
		policyID := d.Get("policy_id").(string)
		epsID := cfg.GetEnterpriseProjectID(d)
		_, err = rules.UpdateWithEpsId(wafClient, updateOpts, policyID, d.Id(), epsID).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF blacklist and whitelist rule: %s", err)
		}
	}

	if d.HasChange("status") {
		if err := updateRuleStatus(wafClient, d, cfg, "whiteblackip"); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceWafRuleBlackListRead(ctx, d, meta)
}

func resourceWafRuleBlackListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	wafClient, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	err = rules.DeleteWithEpsId(wafClient, policyID, d.Id(), cfg.GetEnterpriseProjectID(d)).ExtractErr()
	if err != nil {
		// If the blacklist and whitelist rule does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF blacklist and whitelist rule")
	}
	return nil
}

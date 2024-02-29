package iec

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/iec/v1/firewalls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IEC PUT /v1/firewalls/{firewall_id}/update_firewall_rules
// @API IEC GET /v1/firewalls/{firewall_id}
func ResourceNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkACLRuleCreate,
		ReadContext:   resourceNetworkACLRuleRead,
		UpdateContext: resourceNetworkACLRuleUpdate,
		DeleteContext: resourceNetworkACLRuleDelete,

		Schema: map[string]*schema.Schema{
			"network_acl_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ingress", "egress",
				}, true),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "any",
				ValidateFunc: validation.StringInSlice([]string{
					"tcp", "udp", "icmp", "any",
				}, true),
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "allow",
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny",
				}, true),
			},
			"ip_version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     4,
				Description: "schema: Computed",
			},
			"source_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0/0",
			},
			"destination_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.0.0.0/0",
			},
			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildNetworkACLRule(d *schema.ResourceData, operateType, ruleID string) firewalls.ReqFirewallRulesOpts {
	enabled := d.Get("enabled").(bool)
	ruleOpts := firewalls.ReqFirewallRulesOpts{
		Description: d.Get("description").(string),
		Action:      d.Get("action").(string),
		IPVersion:   d.Get("ip_version").(int),
		Protocol:    d.Get("protocol").(string),
		SrcIPAddr:   d.Get("source_ip_address").(string),
		DstIPAddr:   d.Get("destination_ip_address").(string),
		SrcPort:     d.Get("source_port").(string),
		DstPort:     d.Get("destination_port").(string),
		Enabled:     &enabled,
		OperateType: operateType,
	}
	if operateType != "add" {
		ruleOpts.ID = ruleID
	}
	return ruleOpts
}

func resourceNetworkACLRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return diag.Errorf("error retrieving IEC network ACL %s: %s", aclID, err)
	}

	var oldRules, newRules []string
	var opts firewalls.UpdateRuleOpts
	var ruleOpts firewalls.ReqPolicyOpts
	ruleOpts.FirewallRules = &[]firewalls.ReqFirewallRulesOpts{
		buildNetworkACLRule(d, "add", ""),
	}
	if d.Get("direction").(string) == "ingress" {
		oldRules = getFirewallRuleIDs(fwGroup.IngressFWPolicy)
		ruleOpts.PolicyID = fwGroup.IngressFWPolicy.ID
		opts.ReqFirewallInPolicy = &ruleOpts
	} else {
		oldRules = getFirewallRuleIDs(fwGroup.EgressFWPolicy)
		ruleOpts.PolicyID = fwGroup.EgressFWPolicy.ID
		opts.ReqFirewallOutPolicy = &ruleOpts
	}

	log.Printf("[DEBUG] create IEC network ACL rule: %#v", opts)
	fwGroup, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating IEC network ACL rule: %s", err)
	}

	if d.Get("direction").(string) == "ingress" {
		newRules = getFirewallRuleIDs(fwGroup.IngressFWPolicy)
	} else {
		newRules = getFirewallRuleIDs(fwGroup.EgressFWPolicy)
	}

	ruleID := getNewFirewallRuleID(oldRules, newRules)
	if ruleID == "" {
		return diag.Errorf("error creating IEC network ACL rule: not found")
	}

	log.Printf("[DEBUG] create IEC network ACL rule with id %s", ruleID)
	d.SetId(ruleID)

	return resourceNetworkACLRuleRead(ctx, d, meta)
}

func resourceNetworkACLRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	var mErr *multierror.Error
	if err != nil {
		return diag.Errorf("error creating IEC client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "IEC network ACL")
	}

	var fwPolicy firewalls.RespPolicyEntity
	if d.Get("direction").(string) == "ingress" {
		fwPolicy = fwGroup.IngressFWPolicy
	} else {
		fwPolicy = fwGroup.EgressFWPolicy
	}

	ruleID := d.Id()
	ruleEntity := getFirewallRuleEntity(fwPolicy, ruleID)
	if ruleEntity.ID == "" {
		d.SetId("")
		log.Printf("[WARN] the IEC network ACL rule: %s can not be found", ruleID)
		return nil
	}

	log.Printf("[DEBUG] retrieve IEC network ACL rule %s: %#v", ruleID, ruleEntity)
	mErr = multierror.Append(
		mErr,
		d.Set("policy_id", fwPolicy.ID),
		d.Set("description", ruleEntity.Description),
		d.Set("enabled", ruleEntity.Enabled),
		d.Set("action", ruleEntity.Action),
		d.Set("ip_version", ruleEntity.IPVersion),
		d.Set("source_ip_address", ruleEntity.SrcIPAddr),
		d.Set("destination_ip_address", ruleEntity.DstIPAddr),
		d.Set("source_port", ruleEntity.SrcPort),
		d.Set("destination_port", ruleEntity.DstPort),
	)

	if ruleEntity.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", ruleEntity.Protocol)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkACLRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return diag.Errorf("error retrieving IEC network ACL %s: %s", aclID, err)
	}

	var opts firewalls.UpdateRuleOpts
	var ruleOpts firewalls.ReqPolicyOpts
	ruleOpts.FirewallRules = &[]firewalls.ReqFirewallRulesOpts{
		buildNetworkACLRule(d, "modify", d.Id()),
	}

	if d.Get("direction").(string) == "ingress" {
		ruleOpts.PolicyID = fwGroup.IngressFWPolicy.ID
		opts.ReqFirewallInPolicy = &ruleOpts
	} else {
		ruleOpts.PolicyID = fwGroup.EgressFWPolicy.ID
		opts.ReqFirewallOutPolicy = &ruleOpts
	}

	log.Printf("[DEBUG] updating IEC network ACL rule %s: %#v", d.Id(), opts)
	_, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()

	if err != nil {
		return diag.Errorf("error updating IEC network ACL rule: %s", err)
	}

	return resourceNetworkACLRuleRead(ctx, d, meta)
}

func resourceNetworkACLRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	iecClient, err := conf.IECV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return diag.Errorf("error retrieving IEC network ACL %s: %s", aclID, err)
	}

	var opts firewalls.UpdateRuleOpts
	var ruleOpts firewalls.ReqPolicyOpts
	ruleOpts.FirewallRules = &[]firewalls.ReqFirewallRulesOpts{
		buildNetworkACLRule(d, "delete", d.Id()),
	}
	if d.Get("direction").(string) == "ingress" {
		ruleOpts.PolicyID = fwGroup.IngressFWPolicy.ID
		opts.ReqFirewallInPolicy = &ruleOpts
	} else {
		ruleOpts.PolicyID = fwGroup.EgressFWPolicy.ID
		opts.ReqFirewallOutPolicy = &ruleOpts
	}

	log.Printf("[DEBUG] destroy IEC network ACL rule: %s", d.Id())
	_, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()
	if err != nil {
		return diag.Errorf("error deleting IEC network ACL rule: %s", err)
	}

	d.SetId("")
	return nil
}

func getFirewallRuleIDs(fwPolicy firewalls.RespPolicyEntity) []string {
	rawRules := fwPolicy.FirewallRules
	ruleIDs := make([]string, len(rawRules))
	for i, val := range rawRules {
		ruleIDs[i] = val.ID
	}
	return ruleIDs
}

func getFirewallRuleEntity(fwPolicy firewalls.RespPolicyEntity, ruleID string) firewalls.RespFirewallRulesEntity {
	for _, val := range fwPolicy.FirewallRules {
		if val.ID == ruleID {
			return val
		}
	}
	return firewalls.RespFirewallRulesEntity{}
}

func getNewFirewallRuleID(old []string, new []string) string {
	ruleMap := make(map[string]int)
	for _, v := range old {
		ruleMap[v] = 1
	}

	for _, v := range new {
		if ruleMap[v] == 0 {
			return v
		}
	}
	return ""
}

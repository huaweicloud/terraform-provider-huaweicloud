package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/firewalls"
)

func resourceIecNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecNetworkACLRuleCreate,
		Read:   resourceIecNetworkACLRuleRead,
		Update: resourceIecNetworkACLRuleUpdate,
		Delete: resourceIecNetworkACLRuleDelete,

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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
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

func resourceIecNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving IEC network acl %s: %s", aclID, err)
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

	log.Printf("[DEBUG] Create IEC Network ACL rule: %#v", opts)
	fwGroup, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating IEC Network ACL rule: %s", err)
	}

	if d.Get("direction").(string) == "ingress" {
		newRules = getFirewallRuleIDs(fwGroup.IngressFWPolicy)
	} else {
		newRules = getFirewallRuleIDs(fwGroup.EgressFWPolicy)
	}

	ruleID := getNewFirewallRuleID(oldRules, newRules)
	if ruleID == "" {
		return fmt.Errorf("Error creating IEC Network ACL rule: not found")
	}

	log.Printf("[DEBUG] Create Network IEC ACL rule with id %s", ruleID)
	d.SetId(ruleID)

	return resourceIecNetworkACLRuleRead(d, meta)
}

func resourceIecNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "iec network acl")
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
		log.Printf("[WARN] the IEC Network ACL rule: %s can not be found", ruleID)
		return nil
	}

	log.Printf("[DEBUG] Retrieve IEC Network ACL rule %s: %#v", ruleID, ruleEntity)
	d.Set("policy_id", fwPolicy.ID)
	d.Set("description", ruleEntity.Description)
	d.Set("enabled", ruleEntity.Enabled)
	d.Set("action", ruleEntity.Action)
	d.Set("ip_version", ruleEntity.IPVersion)
	d.Set("source_ip_address", ruleEntity.SrcIPAddr)
	d.Set("destination_ip_address", ruleEntity.DstIPAddr)
	d.Set("source_port", ruleEntity.SrcPort)
	d.Set("destination_port", ruleEntity.DstPort)

	if ruleEntity.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", ruleEntity.Protocol)
	}

	return nil
}

func resourceIecNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving IEC network acl %s: %s", aclID, err)
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

	log.Printf("[DEBUG] Updating IEC Network ACL rule %s: %#v", d.Id(), opts)
	fwGroup, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()

	if err != nil {
		return fmt.Errorf("Error updating IEC Network ACL rule: %s", err)
	}

	return resourceIecNetworkACLRuleRead(d, meta)
}

func resourceIecNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	aclID := d.Get("network_acl_id").(string)
	fwGroup, err := firewalls.Get(iecClient, aclID).Extract()
	if err != nil {
		return fmt.Errorf("Error retrieving IEC network acl %s: %s", aclID, err)
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

	log.Printf("[DEBUG] Destroy IEC Network ACL rule: %s", d.Id())
	fwGroup, err = firewalls.UpdateRule(iecClient, aclID, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error deleting IEC Network ACL rule: %s", err)
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

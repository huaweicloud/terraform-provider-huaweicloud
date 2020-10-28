package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func resourceNetworkACLRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkACLRuleCreate,
		Read:   resourceNetworkACLRuleRead,
		Update: resourceNetworkACLRuleUpdate,
		Delete: resourceNetworkACLRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"tcp", "udp", "icmp", "any",
				}, true),
			},
			"action": {
				Type:     schema.TypeString,
				Required: true,
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
			},
			"destination_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
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
		},
	}
}

func resourceNetworkACLRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	ipVersion := normalizeNetworkACLRuleIPVersion(d.Get("ip_version").(int))
	protocol := normalizeNetworkACLRuleProtocol(d.Get("protocol").(string))

	ruleConfiguration := rules.CreateOpts{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Action:               d.Get("action").(string),
		IPVersion:            ipVersion,
		Protocol:             protocol,
		SourceIPAddress:      d.Get("source_ip_address").(string),
		DestinationIPAddress: d.Get("destination_ip_address").(string),
		SourcePort:           d.Get("source_port").(string),
		DestinationPort:      d.Get("destination_port").(string),
		Enabled:              &enabled,
	}

	log.Printf("[DEBUG] Create Network ACL rule: %#v", ruleConfiguration)
	rule, err := rules.Create(fwClient, ruleConfiguration).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Network ACL rule with id %s", rule.ID)
	d.SetId(rule.ID)

	return resourceNetworkACLRuleRead(d, meta)
}

func resourceNetworkACLRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	rule, err := rules.Get(fwClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Network ACL rule")
	}

	log.Printf("[DEBUG] Retrieve HuaweiCloud Network ACL rule %s: %#v", d.Id(), rule)

	d.Set("action", rule.Action)
	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("ip_version", rule.IPVersion)
	d.Set("source_ip_address", rule.SourceIPAddress)
	d.Set("destination_ip_address", rule.DestinationIPAddress)
	d.Set("source_port", rule.SourcePort)
	d.Set("destination_port", rule.DestinationPort)
	d.Set("enabled", rule.Enabled)

	if rule.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", rule.Protocol)
	}

	return nil
}

func resourceNetworkACLRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	var updateOpts rules.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("protocol") {
		protocol := d.Get("protocol").(string)
		updateOpts.Protocol = &protocol
	}
	if d.HasChange("action") {
		action := d.Get("action").(string)
		updateOpts.Action = &action
	}
	if d.HasChange("ip_version") {
		ipVersion := normalizeNetworkACLRuleIPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}
	if d.HasChange("source_ip_address") {
		sourceIPAddress := d.Get("source_ip_address").(string)
		updateOpts.SourceIPAddress = &sourceIPAddress
	}
	if d.HasChange("source_port") {
		sourcePort := d.Get("source_port").(string)
		updateOpts.SourcePort = &sourcePort
	}
	if d.HasChange("destination_ip_address") {
		destinationIPAddress := d.Get("destination_ip_address").(string)
		updateOpts.DestinationIPAddress = &destinationIPAddress
	}
	if d.HasChange("destination_port") {
		destinationPort := d.Get("destination_port").(string)
		updateOpts.DestinationPort = &destinationPort
	}
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	log.Printf("[DEBUG] Updating Network ACL rule %s: %#v", d.Id(), updateOpts)
	err = rules.Update(fwClient, d.Id(), updateOpts).Err
	if err != nil {
		return err
	}

	return resourceNetworkACLRuleRead(d, meta)
}

func resourceNetworkACLRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	fwClient, err := config.fwV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud fw client: %s", err)
	}

	rule, err := rules.Get(fwClient, d.Id()).Extract()
	if err != nil {
		return err
	}

	policyID, err := assignedPolicyID(fwClient, rule.ID)
	if err != nil {
		return err
	}
	if policyID != "" {
		_, err := policies.RemoveRule(fwClient, policyID, rule.ID).Extract()
		if err != nil {
			return err
		}
	}

	log.Printf("[DEBUG] Destroy Network ACL rule: %s", d.Id())
	return rules.Delete(fwClient, d.Id()).Err
}

func normalizeNetworkACLRuleIPVersion(ipv int) golangsdk.IPVersion {
	// Determine the IP Version
	var ipVersion golangsdk.IPVersion
	switch ipv {
	case 4:
		ipVersion = golangsdk.IPv4
	case 6:
		ipVersion = golangsdk.IPv6
	}

	return ipVersion
}

func normalizeNetworkACLRuleProtocol(p string) rules.Protocol {
	var protocol rules.Protocol
	switch p {
	case "any":
		protocol = rules.ProtocolAny
	case "icmp":
		protocol = rules.ProtocolICMP
	case "tcp":
		protocol = rules.ProtocolTCP
	case "udp":
		protocol = rules.ProtocolUDP
	}

	return protocol
}

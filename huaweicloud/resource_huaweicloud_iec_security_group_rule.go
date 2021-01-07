package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/security/rules"
)

func resourceIecSecurityGroupRule() *schema.Resource {

	return &schema.Resource{
		Create: resourceIecSecurityGroupRuleV1Create,
		Read:   resourceIecSecurityGroupRuleV1Read,
		Delete: resourceIecSecurityGroupRuleV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"egress", "ingress",
				}, true),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"icmp", "tcp", "udp", "gre",
				}, true),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port_range_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"port_range_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"ethertype": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"IPv4", "IPv6"}, true),
				Default:      "IPv4",
			},
			"remote_ip_prefix": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"remote_group_id"},
			},
			"remote_group_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"remote_ip_prefix"},
			},
		},
	}
}

func resourceIecSecurityGroupRuleV1Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	if d.Get("protocol").(string) != "icmp" && d.Get("port_range_min").(int) > d.Get("port_range_max").(int) {
		return fmt.Errorf("The value of `port_range_min` can not be greater than the value of `port_range_max`")
	}

	sgRule := common.ReqSecurityGroupRuleEntity{
		Direction:       d.Get("direction").(string),
		SecurityGroupID: d.Get("security_group_id").(string),
		Description:     d.Get("description").(string),
		EtherType:       d.Get("ethertype").(string),
		Protocol:        d.Get("protocol").(string),
		RemoteIPPrefix:  d.Get("remote_ip_prefix").(string),
		RemoteGroupID:   d.Get("remote_group_id").(string),
	}

	if d.Get("protocol").(string) != "icmp" {
		sgRule.PortRangeMin = d.Get("port_range_min")
		sgRule.PortRangeMax = d.Get("port_range_max")
	}

	createOpts := rules.CreateOpts{
		SecurityGroupRule: &sgRule,
	}

	rule, err := rules.Create(iecClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC Security Group Rule: %s", err)
	}

	d.SetId(rule.SecurityGroupRule.ID)
	return resourceIecSecurityGroupRuleV1Read(d, meta)
}

func resourceIecSecurityGroupRuleV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	rule, err := rules.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "HuaweiCloud IEC Security Group Rule")
	}

	d.Set("description", rule.SecurityGroupRule.Description)
	d.Set("direction", rule.SecurityGroupRule.Direction)
	d.Set("ethertype", rule.SecurityGroupRule.EtherType)
	d.Set("protocol", rule.SecurityGroupRule.Protocol)
	d.Set("security_group_id", rule.SecurityGroupRule.SecurityGroupID)
	d.Set("remote_ip_prefix", rule.SecurityGroupRule.RemoteIPPrefix)
	d.Set("remote_group_id", rule.SecurityGroupRule.RemoteGroupID)

	if rule.SecurityGroupRule.PortRangeMin.(string) != "" {
		d.Set("port_range_min", rule.SecurityGroupRule.PortRangeMin)
	}
	if rule.SecurityGroupRule.PortRangeMax.(string) != "" {
		d.Set("port_range_max", rule.SecurityGroupRule.PortRangeMax)
	}

	return nil
}

func resourceIecSecurityGroupRuleV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSecurityGroupRuleDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      8 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud IEC Security Group Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSecurityGroupRuleDelete(client *golangsdk.ServiceClient, ruleID string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		rule, err := rules.Get(client, ruleID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Security Group Rule %s", ruleID)
				return rule, "DELETED", nil
			}
			return err, "ACTIVE", err
		}

		err = rules.Delete(client, ruleID).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Security Group Rule %s", ruleID)
				return rule, "DELETED", nil
			}
			return rule, "ACTIVE", err
		}
		log.Printf("[DEBUG] HuaweiCloud IEC Security Group Rule %s still active.\n", ruleID)
		return rule, "ACTIVE", nil
	}
}

package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func ResourceNetworkingPortV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingPortV2Create,
		Read:   resourceNetworkingPortV2Read,
		Update: resourceNetworkingPortV2Update,
		Delete: resourceNetworkingPortV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"no_security_groups": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"allowed_address_pairs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      allowedAddressPairsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mac_address": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"extra_dhcp_option": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_version": {
							Type:     schema.TypeInt,
							Default:  4,
							Optional: true,
						},
					},
				},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_fixed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceNetworkingPortV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var securityGroups []string
	v := d.Get("security_group_ids")
	securityGroups = resourcePortSecurityGroupsV2(v.(*schema.Set))
	noSecurityGroups := d.Get("no_security_groups").(bool)

	// Check and make sure an invalid security group configuration wasn't given.
	if noSecurityGroups && len(securityGroups) > 0 {
		return fmt.Errorf("Cannot have both no_security_groups and security_group_ids set")
	}

	createOpts := PortCreateOpts{
		ports.CreateOpts{
			Name:                d.Get("name").(string),
			AdminStateUp:        resourcePortAdminStateUpV2(d),
			NetworkID:           d.Get("network_id").(string),
			MACAddress:          d.Get("mac_address").(string),
			TenantID:            d.Get("tenant_id").(string),
			DeviceOwner:         d.Get("device_owner").(string),
			DeviceID:            d.Get("device_id").(string),
			FixedIPs:            resourcePortFixedIpsV2(d),
			AllowedAddressPairs: resourceAllowedAddressPairsV2(d),
		},
		MapValueSpecs(d),
	}

	if noSecurityGroups {
		securityGroups = []string{}
		createOpts.SecurityGroups = &securityGroups
	}

	// Only set SecurityGroups if one was specified.
	// Otherwise this would mimic the no_security_groups action.
	if len(securityGroups) > 0 {
		createOpts.SecurityGroups = &securityGroups
	}

	// Declare a finalCreateOpts interface to hold either the
	// base create options or the extended DHCP options.
	var finalCreateOpts ports.CreateOptsBuilder
	finalCreateOpts = createOpts

	dhcpOpts := d.Get("extra_dhcp_option").(*schema.Set)
	if dhcpOpts.Len() > 0 {
		finalCreateOpts = extradhcpopts.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			ExtraDHCPOpts:     expandNetworkingPortDHCPOptsV2Create(dhcpOpts),
		}
	}

	log.Printf("[DEBUG] huaweicloud_networking_port_v2 create options: %#v", finalCreateOpts)

	// Create a Neutron port and set extra DHCP options if they're specified.
	var p struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err = ports.Create(networkingClient, finalCreateOpts).ExtractInto(&p)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Neutron port: %s", err)
	}
	log.Printf("[INFO] Network ID: %s", p.ID)

	log.Printf("[DEBUG] Waiting for HuaweiCloud Neutron Port (%s) to become available.", p.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkPortActive(networkingClient, p.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(p.ID)

	return resourceNetworkingPortV2Read(d, meta)
}

func resourceNetworkingPortV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var p struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}

	err = ports.Get(networkingClient, d.Id()).ExtractInto(&p)
	if err != nil {
		return CheckDeleted(d, err, "port")
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", d.Id(), p)

	d.Set("name", p.Name)
	d.Set("admin_state_up", p.AdminStateUp)
	d.Set("network_id", p.NetworkID)
	d.Set("mac_address", p.MACAddress)
	d.Set("tenant_id", p.TenantID)
	d.Set("device_owner", p.DeviceOwner)
	d.Set("device_id", p.DeviceID)

	// Create a slice of all returned Fixed IPs.
	// This will be in the order returned by the API,
	// which is usually alpha-numeric.
	var ips []string
	for _, ipObject := range p.FixedIPs {
		ips = append(ips, ipObject.IPAddress)
	}
	d.Set("all_fixed_ips", ips)

	// Set all security groups.
	// This can be different from what the user specified since
	// the port can have the "default" group automatically applied.
	d.Set("all_security_group_ids", p.SecurityGroups)

	// Convert AllowedAddressPairs to list of map
	var pairs []map[string]interface{}
	for _, pairObject := range p.AllowedAddressPairs {
		pair := make(map[string]interface{})
		pair["ip_address"] = pairObject.IPAddress
		pair["mac_address"] = pairObject.MACAddress
		pairs = append(pairs, pair)
	}
	d.Set("allowed_address_pairs", pairs)
	d.Set("extra_dhcp_option", flattenNetworkingPortDHCPOptsV2(p.ExtraDHCPOptsExt))

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingPortV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	v := d.Get("security_group_ids").(*schema.Set)
	securityGroups := resourcePortSecurityGroupsV2(v)
	noSecurityGroups := d.Get("no_security_groups").(bool)

	// Check and make sure an invalid security group configuration wasn't given.
	if noSecurityGroups && len(securityGroups) > 0 {
		return fmt.Errorf("Cannot have both no_security_groups and security_group_ids set")
	}

	var hasChange bool
	var updateOpts ports.UpdateOpts

	if d.HasChange("allowed_address_pairs") {
		hasChange = true
		aap := resourceAllowedAddressPairsV2(d)
		updateOpts.AllowedAddressPairs = &aap
	}

	if d.HasChange("no_security_groups") {
		if noSecurityGroups {
			hasChange = true
			v := []string{}
			updateOpts.SecurityGroups = &v
		}
	}

	if d.HasChange("security_group_ids") {
		hasChange = true
		sgs := d.Get("security_group_ids").(*schema.Set)
		securityGroups := resourcePortSecurityGroupsV2(sgs)
		updateOpts.SecurityGroups = &securityGroups
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		updateOpts.AdminStateUp = resourcePortAdminStateUpV2(d)
	}

	if d.HasChange("device_owner") {
		hasChange = true
		updateOpts.DeviceOwner = d.Get("device_owner").(string)
	}

	if d.HasChange("device_id") {
		hasChange = true
		updateOpts.DeviceID = d.Get("device_id").(string)
	}

	if d.HasChange("fixed_ip") {
		hasChange = true
		updateOpts.FixedIPs = resourcePortFixedIpsV2(d)
	}

	if hasChange {
		log.Printf("[DEBUG] Updating Port %s with options: %+v", d.Id(), updateOpts)

		_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud Neutron Port: %s", err)
		}
	}

	// Next, perform any dhcp option changes.
	if d.HasChange("extra_dhcp_option") {
		o, n := d.GetChange("extra_dhcp_option")
		oldDHCPOpts := o.(*schema.Set)
		newDHCPOpts := n.(*schema.Set)

		// Delete all old DHCP options, regardless of if they still exist.
		// If they do still exist, they will be re-added below.
		if oldDHCPOpts.Len() != 0 {
			deleteExtraDHCPOpts := expandNetworkingPortDHCPOptsV2Delete(oldDHCPOpts)
			dhcpUpdateOpts := extradhcpopts.UpdateOptsExt{
				UpdateOptsBuilder: &ports.UpdateOpts{},
				ExtraDHCPOpts:     deleteExtraDHCPOpts,
			}

			log.Printf("[DEBUG] Deleting old DHCP opts for huaweicloud_networking_port_v2 %s", d.Id())
			_, err = ports.Update(networkingClient, d.Id(), dhcpUpdateOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error updating HuaweiCloud Neutron Port: %s", err)
			}
		}

		// Add any new DHCP options and re-add previously set DHCP options.
		if newDHCPOpts.Len() != 0 {
			updateExtraDHCPOpts := expandNetworkingPortDHCPOptsV2Update(newDHCPOpts)
			dhcpUpdateOpts := extradhcpopts.UpdateOptsExt{
				UpdateOptsBuilder: &ports.UpdateOpts{},
				ExtraDHCPOpts:     updateExtraDHCPOpts,
			}

			log.Printf("[DEBUG] Updating huaweicloud_networking_port_v2 %s with options: %#v", d.Id(), dhcpUpdateOpts)
			_, err = ports.Update(networkingClient, d.Id(), dhcpUpdateOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error updating huaweicloud_networking_port_v2 %s: %s", d.Id(), err)
			}
		}
	}

	return resourceNetworkingPortV2Read(d, meta)
}

func resourceNetworkingPortV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkPortDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Neutron Network: %s", err)
	}

	d.SetId("")
	return nil
}

func resourcePortSecurityGroupsV2(v *schema.Set) []string {
	var securityGroups []string
	for _, v := range v.List() {
		securityGroups = append(securityGroups, v.(string))
	}
	return securityGroups
}

func resourcePortFixedIpsV2(d *schema.ResourceData) interface{} {
	rawIP := d.Get("fixed_ip").([]interface{})

	if len(rawIP) == 0 {
		return nil
	}

	ip := make([]ports.IP, len(rawIP))
	for i, raw := range rawIP {
		rawMap := raw.(map[string]interface{})
		ip[i] = ports.IP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		}
	}
	return ip
}

func resourceAllowedAddressPairsV2(d *schema.ResourceData) []ports.AddressPair {
	// ports.AddressPair
	rawPairs := d.Get("allowed_address_pairs").(*schema.Set).List()

	pairs := make([]ports.AddressPair, len(rawPairs))
	for i, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		pairs[i] = ports.AddressPair{
			IPAddress:  rawMap["ip_address"].(string),
			MACAddress: rawMap["mac_address"].(string),
		}
	}
	return pairs
}

func resourcePortAdminStateUpV2(d *schema.ResourceData) *bool {
	value := false

	if raw, ok := d.GetOk("admin_state_up"); ok && raw == true {
		value = true
	}

	return &value
}

func allowedAddressPairsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-%s", m["ip_address"].(string), m["mac_address"].(string)))

	return hashcode.String(buf.String())
}

func waitForNetworkPortActive(networkingClient *golangsdk.ServiceClient, portId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, err := ports.Get(networkingClient, portId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud Neutron Port: %+v", p)
		if p.Status == "DOWN" || p.Status == "ACTIVE" {
			return p, "ACTIVE", nil
		}

		return p, p.Status, nil
	}
}

func waitForNetworkPortDelete(networkingClient *golangsdk.ServiceClient, portId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud Neutron Port %s", portId)

		p, err := ports.Get(networkingClient, portId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Port %s", portId)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		err = ports.Delete(networkingClient, portId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Port %s", portId)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloud Port %s still active.\n", portId)
		return p, "ACTIVE", nil
	}
}

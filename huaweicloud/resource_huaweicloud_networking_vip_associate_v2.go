package huaweicloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func resourceNetworkingVIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingVIPAssociateV2Create,
		Read:   resourceNetworkingVIPAssociateV2Read,
		Delete: resourceNetworkingVIPAssociateV2Delete,

		Schema: map[string]*schema.Schema{
			"vip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_ids": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"vip_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vip_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func parseNetworkingVIPAssociateID(id string) (string, []string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", nil, fmt.Errorf("Unable to determine vip association ID")
	}

	vipid := idParts[0]
	portids := idParts[1:]

	return vipid, portids, nil
}

func resourceNetworkingPortIDs(d *schema.ResourceData) []string {
	rawPortIDs := d.Get("port_ids").(*schema.Set).List()
	portids := make([]string, len(rawPortIDs))
	for i, raw := range rawPortIDs {
		portids[i] = raw.(string)
	}
	return portids
}

func resourceNetworkingVIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	vipid := d.Get("vip_id").(string)
	portids := resourceNetworkingPortIDs(d)

	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	// port by port
	fauxid := fmt.Sprintf("%s", vipid)
	for _, portid := range portids {
		// First get the port information
		fauxid = fmt.Sprintf("%s/%s", fauxid, portid)
		port, err := ports.Get(networkingClient, portid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "port")
		}

		ipaddress := ""
		if len(port.FixedIPs) > 0 {
			ipaddress = port.FixedIPs[0].IPAddress
		}
		if len(ipaddress) == 0 {
			return fmt.Errorf("IPAddress is empty, Error associate vip: %#v", port)
		}

		// Then get the vip information
		vip, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "vip")
		}

		// Finnaly associate vip to port
		// Update VIP AllowedAddressPairs
		isfound := false
		for _, raw := range vip.AllowedAddressPairs {
			if ipaddress == raw.IPAddress {
				isfound = true
				break
			}
		}

		// If IP Address is found, not to update VIP
		if !isfound {
			pairs := make([]ports.AddressPair, len(vip.AllowedAddressPairs)+1)
			for i, raw := range vip.AllowedAddressPairs {
				pairs[i] = ports.AddressPair{
					IPAddress:  raw.IPAddress,
					MACAddress: raw.MACAddress,
				}
			}
			pairs[len(vip.AllowedAddressPairs)] = ports.AddressPair{
				IPAddress: ipaddress,
			}
			associateOpts := ports.UpdateOpts{
				AllowedAddressPairs: &pairs,
			}

			log.Printf("[DEBUG] VIP Associate %s with options: %#v", vipid, associateOpts)
			_, err = ports.Update(networkingClient, vipid, associateOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error associate vip: %s", err)
			}
		}

		// Update Port AllowedAddressPairs
		portspairs := make([]ports.AddressPair, 1)
		portspairs[0] = ports.AddressPair{
			IPAddress: "1.1.1.1/0",
		}
		portUpdateOpts := ports.UpdateOpts{
			AllowedAddressPairs: &portspairs,
		}

		log.Printf("[DEBUG] Port Update %s with options: %#v", vipid, portUpdateOpts)
		_, err = ports.Update(networkingClient, portid, portUpdateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error update port: %s", err)
		}
	}

	// There's no assciate vip id, therefore a faux ID will be used.
	d.SetId(fauxid)

	return resourceNetworkingVIPAssociateV2Read(d, meta)
}

func resourceNetworkingVIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	// Obtain relevant info from parsing the ID
	vipid, portids, err := parseNetworkingVIPAssociateID(d.Id())
	if err != nil {
		return err
	}

	// First see if the port still exists
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	// Then try to do this by querying the vip API.
	vip, err := ports.Get(networkingClient, vipid).Extract()
	if err != nil {
		return CheckDeleted(d, err, "vip")
	}

	// port by port
	newportids := make(map[string]string)
	for _, portid := range portids {
		p, err := ports.Get(networkingClient, portid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "port")
		}

		for _, ip := range p.FixedIPs {
			for _, addresspair := range vip.AllowedAddressPairs {
				if ip.IPAddress == addresspair.IPAddress {
					// Associated
					newportids[portid] = portid
					break
				}
			}
		}
	}

	// if no port is associated
	if len(newportids) == 0 {
		d.SetId("")
		return nil
	}

	// convert results from map to array
	newresults := make([]string, len(newportids))
	var index = 0
	for newvalue := range newportids {
		newresults[index] = newvalue
		index++
	}

	// Set the attributes pulled from the composed resource ID
	d.Set("vip_id", vipid)
	d.Set("port_ids", newresults)
	d.Set("vip_subnet_id", vip.FixedIPs[0].SubnetID)
	d.Set("vip_ip_address", vip.FixedIPs[0].IPAddress)

	return nil
}

func resourceNetworkingVIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	id := d.Id()
	vipid, portids, err := parseNetworkingVIPAssociateID(id)
	if err != nil {
		return err
	}

	// port by port
	for _, portid := range portids {
		// First get the port information
		port, err := ports.Get(networkingClient, portid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "port")
		}

		ipaddress := ""
		if len(port.FixedIPs) > 0 {
			ipaddress = port.FixedIPs[0].IPAddress
		}
		if len(ipaddress) == 0 {
			return fmt.Errorf("IPAddress is empty, Error disassociate vip: %#v", port)
		}

		// Then get the vip information
		vip, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			return CheckDeleted(d, err, "vip")
		}

		// Update VIP AllowedAddressPairs
		isfound := false
		for _, raw := range vip.AllowedAddressPairs {
			if ipaddress == raw.IPAddress {
				isfound = true
				break
			}
		}

		// If IP Address is found, need to update VIP
		if isfound {
			pairs := make([]ports.AddressPair, len(vip.AllowedAddressPairs)-1)
			i := 0
			for _, raw := range vip.AllowedAddressPairs {
				if ipaddress != raw.IPAddress {
					pairs[i] = ports.AddressPair{
						IPAddress:  raw.IPAddress,
						MACAddress: raw.MACAddress,
					}
					i++
				}
			}
			disassociateOpts := ports.UpdateOpts{
				AllowedAddressPairs: &pairs,
			}

			log.Printf("[DEBUG] VIP Disassociate %s with options: %#v", vipid, disassociateOpts)
			_, err = ports.Update(networkingClient, vipid, disassociateOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error disassociate vip: %s", err)
			}
		}
	}

	d.SetId("")
	log.Printf("[DEBUG] Successfully disassociate vip (%s)", id)
	return nil
}

package huaweicloud

// This set of code handles all functions required to configure networking
// on an huaweicloud_compute_instance_v2 resource.
//
// This is a complicated task because it's not possible to obtain all
// information in a single API call. In fact, it even traverses multiple
// HuaweiCloud services.
//
// The end result, from the user's point of view, is a structured set of
// understandable network information within the instance resource.

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/servers"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

// InstanceNIC is a structured representation of a servers.Server virtual NIC
type InstanceNIC struct {
	NetworkID string
	PortID    string
	FixedIPv4 string
	FixedIPv6 string
	MAC       string
	Fetched   bool
}

// InstanceNetwork represents a collection of network information that a
// Terraform instance needs to satisfy all network information requirements.
type InstanceNetwork struct {
	UUID          string
	Name          string
	Port          string
	FixedIP       string
	AccessNetwork bool
}

// expandInstanceNetworks builds a []servers.Network for use in creating an Instance.
func expandInstanceNetworks(d *schema.ResourceData) ([]servers.Network, error) {
	var instanceNetworks []servers.Network

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		nic := v.(map[string]interface{})
		network := servers.Network{
			UUID:    nic["uuid"].(string),
			Port:    nic["port"].(string),
			FixedIP: nic["fixed_ip_v4"].(string),
		}
		if network.UUID == "" && network.Port == "" {
			return nil, fmt.Errorf(
				"At least one of network.uuid or network.port must be set.")
		}
		instanceNetworks = append(instanceNetworks, network)
	}

	log.Printf("[DEBUG] expand Instance Networks opts: %#v", instanceNetworks)
	return instanceNetworks, nil
}

// getInstanceAddresses parses a server.Server's Address field into a structured
// InstanceNIC list struct.
func getInstanceAddresses(d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]InstanceNIC, error) {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	allInstanceNics := make([]InstanceNIC, 0)
	var networkID string
	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			// Skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			p, err := ports.Get(networkingClient, addr.PortID).Extract()
			if err != nil {
				networkID = ""
				log.Printf("[DEBUG] getInstanceAddresses: failed to fetch port %s", addr.PortID)
			} else {
				networkID = p.NetworkID
			}

			instanceNIC := InstanceNIC{
				NetworkID: networkID,
				PortID:    addr.PortID,
				MAC:       addr.MacAddr,
			}
			if addr.Version == "6" {
				instanceNIC.FixedIPv6 = addr.Addr
			} else {
				instanceNIC.FixedIPv4 = addr.Addr
			}

			allInstanceNics = append(allInstanceNics, instanceNIC)
		}
	}

	log.Printf("[DEBUG] get all of the Instance Addresses: %#v", allInstanceNics)

	return allInstanceNics, nil
}

// getAllInstanceNetworks loops through the networks defined in the Terraform
// configuration
func getAllInstanceNetworks(d *schema.ResourceData) []InstanceNetwork {
	var instanceNetworks []InstanceNetwork

	networks := d.Get("network").([]interface{})
	for _, v := range networks {
		nic := v.(map[string]interface{})
		network := InstanceNetwork{
			UUID:          nic["uuid"].(string),
			Port:          nic["port"].(string),
			FixedIP:       nic["fixed_ip_v4"].(string),
			AccessNetwork: nic["access_network"].(bool),
		}
		instanceNetworks = append(instanceNetworks, network)
	}

	log.Printf("[DEBUG] get all of the Instance Networks: %#v", instanceNetworks)
	return instanceNetworks
}

// flattenInstanceNetworks collects instance network information from different
// sources and aggregates it all together into a map array.
func flattenInstanceNetworks(
	d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]map[string]interface{}, error) {

	allInstanceNetworks := getAllInstanceNetworks(d)
	allInstanceNics, _ := getInstanceAddresses(d, meta, server)

	networks := []map[string]interface{}{}
	// Loop through all networks and addresses, merge relevant address details.
	for _, instanceNetwork := range allInstanceNetworks {
		for i := range allInstanceNics {
			isExist := false
			nic := &allInstanceNics[i]
			// seem port as the unique key
			if instanceNetwork.Port != "" && instanceNetwork.Port == nic.PortID {
				nic.Fetched = true
				isExist = true
			} else if instanceNetwork.UUID == nic.NetworkID && nic.Fetched == false {
				// Only use one NIC since it's possible the user defined another NIC
				// on this same network in another Terraform network block.
				nic.Fetched = true
				isExist = true
			}

			if isExist {
				v := map[string]interface{}{
					"uuid":           nic.NetworkID,
					"port":           nic.PortID,
					"fixed_ip_v4":    nic.FixedIPv4,
					"fixed_ip_v6":    nic.FixedIPv6,
					"mac":            nic.MAC,
					"access_network": instanceNetwork.AccessNetwork,
				}
				networks = append(networks, v)
				break
			}
		}
	}

	log.Printf("[DEBUG] flatten Instance Networks: %#v", networks)
	return networks, nil
}

// getInstanceAccessAddresses determines the best IP address to communicate
// with the instance. It does this by looping through all networks and looking
// for a valid IP address. Priority is given to a network that was flagged as
// an access_network.
func getInstanceAccessAddresses(
	d *schema.ResourceData, networks []map[string]interface{}) (string, string) {

	var hostv4, hostv6 string

	// Loop through all networks
	// If the network has a valid fixed v4 or fixed v6 address
	// and hostv4 or hostv6 is not set, set hostv4/hostv6.
	// If the network is an "access_network" overwrite hostv4/hostv6.
	for _, n := range networks {
		var accessNetwork bool

		if an, ok := n["access_network"].(bool); ok && an {
			accessNetwork = true
		}

		if fixedIPv4, ok := n["fixed_ip_v4"].(string); ok && fixedIPv4 != "" {
			if hostv4 == "" || accessNetwork {
				hostv4 = fixedIPv4
			}
		}

		if fixedIPv6, ok := n["fixed_ip_v6"].(string); ok && fixedIPv6 != "" {
			if hostv6 == "" || accessNetwork {
				hostv6 = fixedIPv6
			}
		}
	}

	log.Printf("[DEBUG] compute instance Network Access Addresses: %s, %s", hostv4, hostv6)

	return hostv4, hostv6
}

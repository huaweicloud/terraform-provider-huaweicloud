package ecs

// This set of code handles all functions required to configure networking
// on an huaweicloud_compute_instance resource.
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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// InstanceNIC is a structured representation of a servers.Server virtual NIC
type InstanceNIC struct {
	NetworkID       string
	PortID          string
	FixedIPv4       string
	FixedIPv6       string
	MAC             string
	SourceDestCheck bool
	Fetched         bool
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

// getInstanceAddresses parses a server.Server's Address field into a structured
// InstanceNIC list struct.
func getInstanceAddresses(d *schema.ResourceData, meta interface{}, server *cloudservers.CloudServer) ([]InstanceNIC, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	networkingClient, err := cfg.NetworkingV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating networking client: %s", err)
	}

	var lastPort string
	allInstanceNics := make([]InstanceNIC, 0)
	for _, addresses := range server.Addresses {
		for _, addr := range addresses {
			// skip if not fixed ip
			if addr.Type != "fixed" {
				continue
			}

			// IPv4 nic and IPv6 nic have the same port ID, and
			// they are continuous in the array, so skip one of them.
			if lastPort == addr.PortID {
				continue
			}

			lastPort = addr.PortID
			p, err := ports.Get(networkingClient, addr.PortID)
			if err != nil {
				log.Printf("[WARN] getInstanceAddresses: failed to fetch port %s", addr.PortID)
				continue
			}

			instanceNIC := InstanceNIC{
				NetworkID:       p.NetworkId,
				PortID:          addr.PortID,
				MAC:             addr.MacAddr,
				SourceDestCheck: len(p.AllowedAddressPairs) == 0,
			}

			for _, portIP := range p.FixedIps {
				if portIP.IpAddress == "" {
					continue
				}

				if utils.IsIPv4Address(portIP.IpAddress) {
					instanceNIC.FixedIPv4 = portIP.IpAddress
				} else {
					instanceNIC.FixedIPv6 = portIP.IpAddress
				}
			}

			allInstanceNics = append(allInstanceNics, instanceNIC)
		}
	}

	log.Printf("[DEBUG] get all of the Instance Addresses from cloud: %#v", allInstanceNics)

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

	log.Printf("[DEBUG] get all of the Instance Networks from config: %#v", instanceNetworks)
	return instanceNetworks
}

// flattenInstanceNetworks collects instance network information from different
// sources and aggregates it all together into a map array.
func flattenInstanceNetworks(d *schema.ResourceData, meta interface{},
	server *cloudservers.CloudServer) ([]map[string]interface{}, error) {
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
			} else if instanceNetwork.UUID == nic.NetworkID && !nic.Fetched {
				// Only use one NIC since it's possible the user defined another NIC
				// on this same network in another Terraform network block.
				nic.Fetched = true
				isExist = true
			}

			if isExist {
				v := map[string]interface{}{
					"uuid":              nic.NetworkID,
					"port":              nic.PortID,
					"fixed_ip_v4":       nic.FixedIPv4,
					"fixed_ip_v6":       nic.FixedIPv6,
					"ipv6_enable":       nic.FixedIPv6 != "",
					"source_dest_check": nic.SourceDestCheck,
					"mac":               nic.MAC,
					"access_network":    instanceNetwork.AccessNetwork,
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
func getInstanceAccessAddresses(networks []map[string]interface{}) (ipv4Addr, ipv6Addr string) {
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
			if ipv4Addr == "" || accessNetwork {
				ipv4Addr = fixedIPv4
			}
		}

		if fixedIPv6, ok := n["fixed_ip_v6"].(string); ok && fixedIPv6 != "" {
			if ipv6Addr == "" || accessNetwork {
				ipv6Addr = fixedIPv6
			}
		}
	}

	log.Printf("[DEBUG] compute instance Network Access Addresses: %s, %s", ipv4Addr, ipv6Addr)
	return
}

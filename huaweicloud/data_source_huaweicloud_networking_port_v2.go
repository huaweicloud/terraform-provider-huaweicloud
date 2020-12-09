package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/ports"
)

func DataSourceNetworkingPortV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingPortV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"device_owner": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

func dataSourceNetworkingPortV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	listOpts := ports.ListOpts{}

	if v, ok := d.GetOk("port_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("device_owner"); ok {
		listOpts.DeviceOwner = v.(string)
	}

	if v, ok := d.GetOk("mac_address"); ok {
		listOpts.MACAddress = v.(string)
	}

	if v, ok := d.GetOk("device_id"); ok {
		listOpts.DeviceID = v.(string)
	}

	allPages, err := ports.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list huaweicloud_networking_ports_v2: %s", err)
	}

	var allPorts []ports.Port

	err = ports.ExtractPortsInto(allPages, &allPorts)
	if err != nil {
		return fmt.Errorf("Unable to retrieve huaweicloud_networking_ports_v2: %s", err)
	}

	if len(allPorts) == 0 {
		return fmt.Errorf("No huaweicloud_networking_port_v2 found")
	}

	var portsList []ports.Port

	// Filter returned Fixed IPs by a "fixed_ip".
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			for _, ipObject := range p.FixedIPs {
				if v.(string) == ipObject.IPAddress {
					portsList = append(portsList, p)
				}
			}
		}
		if len(portsList) == 0 {
			log.Printf("No huaweicloud_networking_port_v2 found after the 'fixed_ip' filter")
			return fmt.Errorf("No huaweicloud_networking_port_v2 found")
		}
	} else {
		portsList = allPorts
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	if len(securityGroups) > 0 {
		var sgPorts []ports.Port
		for _, p := range portsList {
			for _, sg := range p.SecurityGroups {
				if strSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}
		if len(sgPorts) == 0 {
			log.Printf("[DEBUG] No huaweicloud_networking_port_v2 found after the 'security_group_ids' filter")
			return fmt.Errorf("No huaweicloud_networking_port_v2 found")
		}
		portsList = sgPorts
	}

	if len(portsList) > 1 {
		return fmt.Errorf("More than one huaweicloud_networking_port_v2 found (%d)", len(portsList))
	}

	port := portsList[0]

	log.Printf("[DEBUG] Retrieved huaweicloud_networking_port_v2 %s: %+v", port.ID, port)
	d.SetId(port.ID)

	d.Set("port_id", port.ID)
	d.Set("name", port.Name)
	d.Set("admin_state_up", port.AdminStateUp)
	d.Set("network_id", port.NetworkID)
	d.Set("tenant_id", port.TenantID)
	d.Set("project_id", port.ProjectID)
	d.Set("device_owner", port.DeviceOwner)
	d.Set("mac_address", port.MACAddress)
	d.Set("device_id", port.DeviceID)
	d.Set("region", GetRegion(d, config))
	d.Set("all_security_group_ids", port.SecurityGroups)
	d.Set("all_fixed_ips", expandNetworkingPortFixedIPToStringSlice(port.FixedIPs))

	return nil
}

func expandNetworkingPortFixedIPToStringSlice(fixedIPs []ports.IP) []string {
	s := make([]string, len(fixedIPs))
	for i, fixedIP := range fixedIPs {
		s[i] = fixedIP.IPAddress
	}

	return s
}

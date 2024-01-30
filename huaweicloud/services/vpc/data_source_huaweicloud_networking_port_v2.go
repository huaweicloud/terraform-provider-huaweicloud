package vpc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPC GET /v2.0/ports
func DataSourceNetworkingPortV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingPortV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			// marked as computed
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
			},
			"device_owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
			},
			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "schema: Computed",
			},

			// Computed
			"all_allowed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_fixed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_security_group_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			// marked as deprecated
			"admin_state_up": {
				Type:       schema.TypeBool,
				Optional:   true,
				Computed:   true,
				Deprecated: "this field is deprecated",
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "this field is deprecated",
			},
			"project_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "this field is deprecated",
			},
		},
	}
}

func getNetworkingPortOpts(d *schema.ResourceData) ports.ListOpts {
	listOpts := ports.ListOpts{}

	if v, ok := d.GetOk("port_id"); ok {
		listOpts.ID = v.(string)
	}
	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}
	if v, ok := d.GetOk("mac_address"); ok {
		listOpts.MACAddress = v.(string)
	}
	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}
	if v, ok := d.GetOk("admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
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
	if v, ok := d.GetOk("device_id"); ok {
		listOpts.DeviceID = v.(string)
	}

	return listOpts
}

func dataSourceNetworkingPortV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	listOpts := getNetworkingPortOpts(d)

	allPages, err := ports.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to list networking ports v2: %s", err)
	}

	var allPorts []ports.Port

	err = ports.ExtractPortsInto(allPages, &allPorts)
	if err != nil {
		return diag.Errorf("unable to retrieve networking ports v2: %s", err)
	}

	if len(allPorts) == 0 {
		return diag.Errorf("no networking port found")
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
			log.Printf("No networking port found after the 'fixed_ip' filter")
			return diag.Errorf("no networking port found")
		}
	} else {
		portsList = allPorts
	}

	sgRaw := d.Get("security_group_ids").(*schema.Set).List()
	securityGroups := utils.ExpandToStringList(sgRaw)
	if len(securityGroups) > 0 {
		var sgPorts []ports.Port
		for _, p := range portsList {
			for _, sg := range p.SecurityGroups {
				if utils.StrSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}
		if len(sgPorts) == 0 {
			log.Printf("[DEBUG] No networking port found after the 'security_group_ids' filter")
			return diag.Errorf("no networking port found")
		}
		portsList = sgPorts
	}

	if len(portsList) > 1 {
		return diag.Errorf("more than one networking port found (%d)", len(portsList))
	}

	port := portsList[0]

	log.Printf("[DEBUG] Retrieved networking port %s: %+v", port.ID, port)
	d.SetId(port.ID)

	d.Set("port_id", port.ID)
	d.Set("name", port.Name)
	d.Set("status", port.Status)
	d.Set("admin_state_up", port.AdminStateUp)
	d.Set("network_id", port.NetworkID)
	d.Set("mac_address", port.MACAddress)
	d.Set("device_owner", port.DeviceOwner)
	d.Set("device_id", port.DeviceID)
	d.Set("region", config.GetRegion(d))
	d.Set("all_security_group_ids", port.SecurityGroups)
	d.Set("all_allowed_ips", expandNetworkingPortAllowedAddressPairToStringSlice(port.AllowedAddressPairs))
	d.Set("all_fixed_ips", expandNetworkingPortFixedIPToStringSlice(port.FixedIPs))

	return nil
}

func expandNetworkingPortAllowedAddressPairToStringSlice(addressPairs []ports.AddressPair) []string {
	s := make([]string, len(addressPairs))
	for i, addressPair := range addressPairs {
		s[i] = addressPair.IPAddress
	}

	return s
}

func expandNetworkingPortFixedIPToStringSlice(fixedIPs []ports.IP) []string {
	s := make([]string, len(fixedIPs))
	for i, fixedIP := range fixedIPs {
		s[i] = fixedIP.IPAddress
	}

	return s
}

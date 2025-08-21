package deprecated

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/subnets"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingSubnetV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingSubnetV2Create,
		Read:   resourceNetworkingSubnetV2Read,
		Update: resourceNetworkingSubnetV2Update,
		Delete: resourceNetworkingSubnetV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_vpc_subnet resource instead",

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
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"allocation_pools": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:     schema.TypeString,
							Required: true,
						},
						"end": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"gateway_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"no_gateway"},
				Optional:      true,
				Computed:      true,
			},
			"no_gateway": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"gateway_ip"},
				Optional:      true,
				Default:       false,
			},
			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  4,
				ForceNew: true,
			},
			"enable_dhcp": {
				Type:         schema.TypeBool,
				Optional:     true,
				Default:      true,
				ValidateFunc: utils.ValidateTrueOnly,
			},
			"dns_nameservers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"host_routes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"ipv6_address_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateSubnetV2IPv6Mode,
			},
			"ipv6_ra_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: utils.ValidateSubnetV2IPv6Mode,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceNetworkingSubnetV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	createOpts := SubnetCreateOpts{
		subnets.CreateOpts{
			NetworkID:       d.Get("network_id").(string),
			CIDR:            d.Get("cidr").(string),
			Name:            d.Get("name").(string),
			TenantID:        d.Get("tenant_id").(string),
			IPv6AddressMode: d.Get("ipv6_address_mode").(string),
			IPv6RAMode:      d.Get("ipv6_ra_mode").(string),
			AllocationPools: resourceSubnetAllocationPoolsV2(d),
			DNSNameservers:  resourceSubnetDNSNameserversV2(d),
			HostRoutes:      resourceSubnetHostRoutesV2(d),
			EnableDHCP:      nil,
		},
		MapValueSpecs(d),
	}

	if v, ok := d.GetOk("gateway_ip"); ok {
		gatewayIP := v.(string)
		createOpts.GatewayIP = &gatewayIP
	}

	noGateway := d.Get("no_gateway").(bool)
	if noGateway {
		gatewayIP := ""
		createOpts.GatewayIP = &gatewayIP
	}

	enableDHCP := d.Get("enable_dhcp").(bool)
	createOpts.EnableDHCP = &enableDHCP

	if v, ok := d.GetOk("ip_version"); ok {
		ipVersion := resourceNetworkingSubnetV2DetermineIPVersion(v.(int))
		createOpts.IPVersion = ipVersion
	}

	s, err := subnets.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Neutron subnet: %s", err)
	}

	logp.Printf("[DEBUG] Waiting for Subnet (%s) to become available", s.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSubnetActive(networkingClient, s.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(s.ID)

	logp.Printf("[DEBUG] Created Subnet %s: %#v", s.ID, s)
	return resourceNetworkingSubnetV2Read(d, meta)
}

func resourceNetworkingSubnetV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	s, err := subnets.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "subnet")
	}

	logp.Printf("[DEBUG] Retrieved Subnet %s: %#v", d.Id(), s)

	d.Set("network_id", s.NetworkID)
	d.Set("cidr", s.CIDR)
	d.Set("ip_version", s.IPVersion)
	d.Set("name", s.Name)
	d.Set("tenant_id", s.TenantID)
	d.Set("dns_nameservers", s.DNSNameservers)
	d.Set("enable_dhcp", s.EnableDHCP)
	d.Set("network_id", s.NetworkID)
	d.Set("ipv6_address_mode", s.IPv6AddressMode)
	d.Set("ipv6_ra_mode", s.IPv6RAMode)

	// Set the host_routes
	var hostRoutes []map[string]interface{} = make([]map[string]interface{}, len(s.HostRoutes))
	for i, v := range s.HostRoutes {
		routes := make(map[string]interface{})
		routes["destination_cidr"] = v.DestinationCIDR
		routes["next_hop"] = v.NextHop
		hostRoutes[i] = routes
	}
	if err = d.Set("host_routes", hostRoutes); err != nil {
		return fmtp.Errorf("Saving host_routes failed: %s", err)
	}

	// Set the allocation_pools
	var allocationPools []map[string]interface{}
	for _, v := range s.AllocationPools {
		pool := make(map[string]interface{})
		pool["start"] = v.Start
		pool["end"] = v.End

		allocationPools = append(allocationPools, pool)
	}
	d.Set("allocation_pools", allocationPools)

	// Set the subnet's Gateway IP.
	gatewayIP := s.GatewayIP
	d.Set("gateway_ip", s.GatewayIP)

	// Based on the subnet's Gateway IP, set `no_gateway` accordingly.
	if gatewayIP == "" {
		d.Set("no_gateway", true)
	} else {
		d.Set("no_gateway", false)
	}

	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceNetworkingSubnetV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("gateway_ip") {
		updateOpts.GatewayIP = nil
		if v, ok := d.GetOk("gateway_ip"); ok {
			gatewayIP := v.(string)
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("no_gateway") {
		if d.Get("no_gateway").(bool) {
			gatewayIP := ""
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("dns_nameservers") {
		updateOpts.DNSNameservers = resourceSubnetDNSNameserversV2(d)
	}

	if d.HasChange("host_routes") {
		updateOpts.HostRoutes = resourceSubnetHostRoutesV2(d)
	}

	if d.HasChange("enable_dhcp") {
		v := d.Get("enable_dhcp").(bool)
		updateOpts.EnableDHCP = &v
	}

	if d.HasChange("allocation_pools") {
		updateOpts.AllocationPools = resourceSubnetAllocationPoolsV2(d)
	}

	logp.Printf("[DEBUG] Updating Subnet %s with options: %+v", d.Id(), updateOpts)

	_, err = subnets.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud Neutron Subnet: %s", err)
	}

	return resourceNetworkingSubnetV2Read(d, meta)
}

func resourceNetworkingSubnetV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSubnetDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Neutron Subnet: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceSubnetAllocationPoolsV2(d *schema.ResourceData) []subnets.AllocationPool {
	rawAPs := d.Get("allocation_pools").([]interface{})
	aps := make([]subnets.AllocationPool, len(rawAPs))
	for i, raw := range rawAPs {
		rawMap := raw.(map[string]interface{})
		aps[i] = subnets.AllocationPool{
			Start: rawMap["start"].(string),
			End:   rawMap["end"].(string),
		}
	}
	return aps
}

func resourceSubnetDNSNameserversV2(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_nameservers").(*schema.Set)
	dnsn := make([]string, rawDNSN.Len())
	for i, raw := range rawDNSN.List() {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

func resourceSubnetHostRoutesV2(d *schema.ResourceData) []subnets.HostRoute {
	rawHR := d.Get("host_routes").([]interface{})
	hr := make([]subnets.HostRoute, len(rawHR))
	for i, raw := range rawHR {
		rawMap := raw.(map[string]interface{})
		hr[i] = subnets.HostRoute{
			DestinationCIDR: rawMap["destination_cidr"].(string),
			NextHop:         rawMap["next_hop"].(string),
		}
	}
	return hr
}

func resourceNetworkingSubnetV2DetermineIPVersion(v int) golangsdk.IPVersion {
	var ipVersion golangsdk.IPVersion
	switch v {
	case 4:
		ipVersion = golangsdk.IPv4
	case 6:
		ipVersion = golangsdk.IPv6
	}

	return ipVersion
}

func waitForSubnetActive(networkingClient *golangsdk.ServiceClient, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := subnets.Get(networkingClient, subnetId).Extract()
		if err != nil {
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Neutron Subnet: %+v", s)
		return s, "ACTIVE", nil
	}
}

func waitForSubnetDelete(networkingClient *golangsdk.ServiceClient, subnetId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Subnet %s.\n", subnetId)

		s, err := subnets.Get(networkingClient, subnetId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Subnet %s", subnetId)
				return s, "DELETED", nil
			}
			return s, "ACTIVE", err
		}

		err = subnets.Delete(networkingClient, subnetId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Subnet %s", subnetId)
				return s, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				return s, "ACTIVE", nil
			}
			return s, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Subnet %s still active.\n", subnetId)
		return s, "ACTIVE", nil
	}
}

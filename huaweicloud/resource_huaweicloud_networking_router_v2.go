package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func resourceNetworkingRouterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRouterV2Create,
		Read:   resourceNetworkingRouterV2Read,
		Update: resourceNetworkingRouterV2Update,
		Delete: resourceNetworkingRouterV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_vpc resource instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"distributed": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"external_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_snat": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"external_fixed_ip": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
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

func resourceNetworkingRouterV2Create(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(*config.Config)
	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	createOpts := RouterCreateOpts{
		routers.CreateOpts{
			Name:     d.Get("name").(string),
			TenantID: d.Get("tenant_id").(string),
		},
		MapValueSpecs(d),
	}

	if asuRaw, ok := d.GetOk("admin_state_up"); ok {
		asu := asuRaw.(bool)
		createOpts.AdminStateUp = &asu
	}

	if dRaw, ok := d.GetOk("distributed"); ok {
		d := dRaw.(bool)
		createOpts.Distributed = &d
	}

	// Gateway settings
	var externalNetworkID string
	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
	}

	if externalNetworkID != "" {
		gatewayInfo := routers.GatewayInfo{
			NetworkID: externalNetworkID,
		}
		createOpts.GatewayInfo = &gatewayInfo
	}

	if esRaw, ok := d.GetOk("enable_snat"); ok {
		if externalNetworkID == "" {
			return fmtp.Errorf("setting enable_snat requires external_network_id to be set")
		}
		es := esRaw.(bool)
		createOpts.GatewayInfo.EnableSNAT = &es
	}

	externalFixedIPs := resourceRouterExternalFixedIPsV2(d)
	if len(externalFixedIPs) > 0 {
		if externalNetworkID == "" {
			return fmtp.Errorf("setting an external_fixed_ip requires external_network_id to be set")
		}
		createOpts.GatewayInfo.ExternalFixedIPs = externalFixedIPs
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := routers.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Neutron router: %s", err)
	}
	logp.Printf("[INFO] Router ID: %s", n.ID)

	logp.Printf("[DEBUG] Waiting for HuaweiCloud Neutron Router (%s) to become available", n.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForRouterActive(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(n.ID)

	return resourceNetworkingRouterV2Read(d, meta)
}

func resourceNetworkingRouterV2Read(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(*config.Config)
	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	n, err := routers.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Neutron Router: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved Router %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("admin_state_up", n.AdminStateUp)
	d.Set("distributed", n.Distributed)
	d.Set("tenant_id", n.TenantID)
	d.Set("region", GetRegion(d, conf))

	// Gateway settings
	d.Set("external_network_id", n.GatewayInfo.NetworkID)
	d.Set("enable_snat", n.GatewayInfo.EnableSNAT)

	var externalFixedIPs []map[string]string
	for _, v := range n.GatewayInfo.ExternalFixedIPs {
		externalFixedIPs = append(externalFixedIPs, map[string]string{
			"subnet_id":  v.SubnetID,
			"ip_address": v.IPAddress,
		})
	}

	if err = d.Set("external_fixed_ip", externalFixedIPs); err != nil {
		logp.Printf("[DEBUG] unable to set external_fixed_ip: %s", err)
	}

	return nil
}

func resourceNetworkingRouterV2Update(d *schema.ResourceData, meta interface{}) error {
	routerId := d.Id()
	config.MutexKV.Lock(routerId)
	defer config.MutexKV.Unlock(routerId)

	conf := meta.(*config.Config)
	networkingClient, err := conf.NetworkingV2Client(GetRegion(d, conf))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var updateOpts routers.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Gateway settings
	var updateGatewaySettings bool
	var externalNetworkID string
	gatewayInfo := routers.GatewayInfo{}

	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
	}

	if externalNetworkID != "" {
		gatewayInfo.NetworkID = externalNetworkID
	}

	if d.HasChange("external_network_id") {
		updateGatewaySettings = true
	}

	if d.HasChange("enable_snat") {
		updateGatewaySettings = true
		if externalNetworkID == "" {
			return fmtp.Errorf("setting enable_snat requires external_network_id to be set")
		}

		enableSNAT := d.Get("enable_snat").(bool)
		gatewayInfo.EnableSNAT = &enableSNAT
	}

	if d.HasChange("external_fixed_ip") {
		updateGatewaySettings = true

		externalFixedIPs := resourceRouterExternalFixedIPsV2(d)
		gatewayInfo.ExternalFixedIPs = externalFixedIPs
		if len(externalFixedIPs) > 0 {
			if externalNetworkID == "" {
				return fmtp.Errorf("setting an external_fixed_ip requires external_network_id to be set")
			}
		}
	}

	if updateGatewaySettings {
		updateOpts.GatewayInfo = &gatewayInfo
	}

	logp.Printf("[DEBUG] Updating Router %s with options: %+v", d.Id(), updateOpts)

	_, err = routers.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating HuaweiCloud Neutron Router: %s", err)
	}

	return resourceNetworkingRouterV2Read(d, meta)
}

func resourceNetworkingRouterV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForRouterDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      8 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Neutron Router: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForRouterActive(networkingClient *golangsdk.ServiceClient, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := routers.Get(networkingClient, routerId).Extract()
		if err != nil {
			return nil, r.Status, err
		}

		logp.Printf("[DEBUG] HuaweiCloud Neutron Router: %+v", r)
		return r, r.Status, nil
	}
}

func waitForRouterDelete(networkingClient *golangsdk.ServiceClient, routerId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Router %s.\n", routerId)

		r, err := routers.Get(networkingClient, routerId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Router %s", routerId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = routers.Delete(networkingClient, routerId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Router %s", routerId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Router %s still active.\n", routerId)
		return r, "ACTIVE", nil
	}
}

func resourceRouterExternalFixedIPsV2(d *schema.ResourceData) []routers.ExternalFixedIP {
	var externalFixedIPs []routers.ExternalFixedIP
	eFIPs := d.Get("external_fixed_ip").([]interface{})

	for _, eFIP := range eFIPs {
		v := eFIP.(map[string]interface{})
		fip := routers.ExternalFixedIP{
			SubnetID:  v["subnet_id"].(string),
			IPAddress: v["ip_address"].(string),
		}
		externalFixedIPs = append(externalFixedIPs, fip)
	}

	return externalFixedIPs
}

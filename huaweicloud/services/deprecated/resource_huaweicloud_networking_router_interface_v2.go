package deprecated

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/layer3/routers"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingRouterInterfaceV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingRouterInterfaceV2Create,
		Read:   resourceNetworkingRouterInterfaceV2Read,
		Delete: resourceNetworkingRouterInterfaceV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use VPC resources instead",

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
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterInterfaceV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	createOpts := routers.AddInterfaceOpts{
		SubnetID: d.Get("subnet_id").(string),
		PortID:   d.Get("port_id").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := routers.AddInterface(networkingClient, d.Get("router_id").(string), createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Neutron router interface: %s", err)
	}
	logp.Printf("[INFO] Router interface Port ID: %s", n.PortID)

	logp.Printf("[DEBUG] Waiting for Router Interface (%s) to become available", n.PortID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForRouterInterfaceActive(networkingClient, n.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(n.PortID)

	return resourceNetworkingRouterInterfaceV2Read(d, meta)
}

func resourceNetworkingRouterInterfaceV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	n, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmtp.Errorf("Error retrieving HuaweiCloud Neutron Router Interface: %s", err)
	}

	logp.Printf("[DEBUG] Retrieved Router Interface %s: %+v", d.Id(), n)

	d.Set("router_id", n.DeviceID)
	d.Set("port_id", n.ID)

	// Set the subnet ID by looking at thet port's FixedIPs.
	// If there's more than one FixedIP, do not set the subnet
	// as it's not possible to confidently determine which subnet
	// belongs to this interface. However, that situation should
	// not happen.
	if len(n.FixedIPs) == 1 {
		d.Set("subnet_id", n.FixedIPs[0].SubnetID)
	}

	d.Set("region", config.GetRegion(d))

	return nil
}

func resourceNetworkingRouterInterfaceV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForRouterInterfaceDelete(networkingClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Neutron Router Interface: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForRouterInterfaceActive(networkingClient *golangsdk.ServiceClient, rId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := ports.Get(networkingClient, rId).Extract()
		if err != nil {
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Neutron Router Interface: %+v", r)
		return r, r.Status, nil
	}
}

func waitForRouterInterfaceDelete(networkingClient *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		routerId := d.Get("router_id").(string)
		routerInterfaceId := d.Id()

		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Router Interface %s.", routerInterfaceId)

		removeOpts := routers.RemoveInterfaceOpts{
			SubnetID: d.Get("subnet_id").(string),
			PortID:   d.Get("port_id").(string),
		}

		r, err := ports.Get(networkingClient, routerInterfaceId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Router Interface %s", routerInterfaceId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		_, err = routers.RemoveInterface(networkingClient, routerId, removeOpts).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud Router Interface %s.", routerInterfaceId)
				return r, "DELETED", nil
			}
			if _, ok := err.(golangsdk.ErrDefault409); ok {
				logp.Printf("[DEBUG] Router Interface %s is still in use.", routerInterfaceId)
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Router Interface %s is still active.", routerInterfaceId)
		return r, "ACTIVE", nil
	}
}

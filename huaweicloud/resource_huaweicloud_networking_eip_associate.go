package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingFloatingIPAssociateV2Create,
		Read:   resourceNetworkingFloatingIPAssociateV2Read,
		Delete: resourceNetworkingFloatingIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingFloatingIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	floatingIP := d.Get("public_ip").(string)
	portID := d.Get("port_id").(string)

	floatingIPID, err := resourceNetworkingFloatingIPAssociateV2IP2ID(networkingClient, floatingIP)
	if err != nil {
		return fmtp.Errorf("Unable to get ID of floating IP: %s", err)
	}

	updateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}

	logp.Printf("[DEBUG] Floating IP Associate Create Options: %#v", updateOpts)
	_, err = floatingips.Update(networkingClient, floatingIPID, updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error associating floating IP %s to port %s: %s",
			floatingIPID, portID, err)
	}

	d.SetId(floatingIPID)

	return resourceNetworkFloatingIPV2Read(d, meta)
}

func resourceNetworkingFloatingIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	floatingIP, err := floatingips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "floating IP")
	}

	d.Set("public_ip", floatingIP.FloatingIP)
	d.Set("port_id", floatingIP.PortID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingFloatingIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	portID := d.Get("port_id").(string)
	updateOpts := floatingips.UpdateOpts{
		PortID: nil,
	}

	logp.Printf("[DEBUG] Floating IP Delete Options: %#v", updateOpts)
	_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error disassociating floating IP %s from port %s: %s",
			d.Id(), portID, err)
	}

	return nil
}

func resourceNetworkingFloatingIPAssociateV2IP2ID(client *golangsdk.ServiceClient, floatingIP string) (string, error) {
	listOpts := floatingips.ListOpts{
		FloatingIP: floatingIP,
	}

	allPages, err := floatingips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		return "", err
	}

	if len(allFloatingIPs) != 1 {
		return "", fmtp.Errorf("unable to determine the ID of %s", floatingIP)
	}

	return allFloatingIPs[0].ID, nil
}

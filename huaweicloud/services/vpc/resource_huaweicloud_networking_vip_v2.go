package vpc

import (
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingVIPV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingVIPV2Create,
		Read:   resourceNetworkingVIPV2Read,
		Update: resourceNetworkingVIPV2Update,
		Delete: resourceNetworkingVIPV2Delete,
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
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"device_owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnet_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"ip_version"},
				Deprecated:    "use ip_version instead",
			},
		},
	}
}

func resourceNetworkingVIPV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	networkingClient, err := config.NetworkingV2Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud VPC client: %s", err)
	}

	networkID := d.Get("network_id").(string)
	n, err := subnets.Get(vpcClient, networkID).Extract()
	if err != nil {
		return fmtp.Errorf("Error retrieving Huaweicloud Subnet %s: %s", networkID, err)
	}

	// validate the value of subnet_id
	subnetID := d.Get("subnet_id").(string)
	if subnetID != "" && subnetID != n.SubnetId && subnetID != n.IPv6SubnetId {
		return fmtp.Errorf("Error invalid value of subnet_id %s, expect to %s or %s",
			subnetID, n.SubnetId, n.IPv6SubnetId)
	}

	// update subnetID for IPv6
	ipVersion := d.Get("ip_version").(int)
	if ipVersion == 6 {
		if n.IPv6SubnetId == "" {
			return fmtp.Errorf("the subnet does not support IPv6, please enable IPv6 first")
		}
		subnetID = n.IPv6SubnetId
	} else {
		subnetID = n.SubnetId
	}

	fixip := make([]ports.IP, 1)
	fixip[0] = ports.IP{
		SubnetID:  subnetID,
		IPAddress: d.Get("ip_address").(string),
	}

	createOpts := ports.CreateOpts{
		Name:        d.Get("name").(string),
		DeviceOwner: "neutron:VIP_PORT",
		NetworkID:   networkID,
		FixedIPs:    fixip,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	vip, err := ports.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Network VIP: %s", err)
	}
	logp.Printf("[DEBUG] Waiting for HuaweiCloud Network VIP (%s) to become available.", vip.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForNetworkVIPActive(networkingClient, vip.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	d.SetId(vip.ID)

	return resourceNetworkingVIPV2Read(d, meta)
}

func resourceNetworkingVIPV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	vip, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "vip")
	}

	logp.Printf("[DEBUG] Retrieved VIP %s: %+v", d.Id(), vip)

	d.SetId(vip.ID)
	// Computed values
	d.Set("network_id", vip.NetworkID)
	if len(vip.FixedIPs) > 0 {
		addr := vip.FixedIPs[0].IPAddress
		d.Set("ip_address", addr)
		d.Set("subnet_id", vip.FixedIPs[0].SubnetID)
		if utils.IsIPv4Address(addr) {
			d.Set("ip_version", 4)
		} else {
			d.Set("ip_version", 6)
		}
	} else {
		d.Set("subnet_id", "")
		d.Set("ip_address", "")
	}

	d.Set("name", vip.Name)
	d.Set("status", vip.Status)
	d.Set("device_owner", vip.DeviceOwner)
	d.Set("mac_address", vip.MACAddress)

	return nil
}

func resourceNetworkingVIPV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	if d.HasChange("name") {
		updateOpts := ports.UpdateOpts{
			Name: d.Get("name").(string),
		}
		logp.Printf("[DEBUG] Updating networking vip %s with options: %#v", d.Id(), updateOpts)

		_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud networking vip: %s", err)
		}
	}

	return resourceNetworkingVIPV2Read(d, meta)
}

func resourceNetworkingVIPV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkVIPDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud Network VIP: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForNetworkVIPActive(networkingClient *golangsdk.ServiceClient, vipid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			return nil, "", err
		}

		logp.Printf("[DEBUG] HuaweiCloud Network Port: %+v", p)
		if p.Status == "DOWN" || p.Status == "ACTIVE" {
			return p, "ACTIVE", nil
		}

		return p, p.Status, nil
	}
}

func waitForNetworkVIPDelete(networkingClient *golangsdk.ServiceClient, vipid string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		logp.Printf("[DEBUG] Attempting to delete HuaweiCloud Network VIP %s", vipid)

		p, err := ports.Get(networkingClient, vipid).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud VIP %s", vipid)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		err = ports.Delete(networkingClient, vipid).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] Successfully deleted HuaweiCloud VIP %s", vipid)
				return p, "DELETED", nil
			}
			return p, "ACTIVE", err
		}

		logp.Printf("[DEBUG] HuaweiCloud VIP %s still active.\n", vipid)
		return p, "ACTIVE", nil
	}
}

package vpc

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	"github.com/chnsz/golangsdk/openstack/networking/v1/subnets"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceNetworkingVip() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingVipCreate,
		ReadContext:   resourceNetworkingVipRead,
		UpdateContext: resourceNetworkingVipUpdate,
		DeleteContext: resourceNetworkingVipDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
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

			// Computed
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

			// Deprecated
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

func resourceNetworkingVipCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC network v1 client: %s", err)
	}

	networkId := d.Get("network_id").(string)
	n, err := subnets.Get(client, networkId).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error retrieving HuaweiCloud subnet by network ID (%s): %s", networkId, err)
	}

	// Check whether the subnet ID entered by the user belongs to the same subnet as the network ID.
	subnetId := d.Get("subnet_id").(string)
	if subnetId != "" && subnetId != n.SubnetId && subnetId != n.IPv6SubnetId {
		return fmtp.DiagErrorf("The subnet ID does not belong to the subnet where the network ID is located.")
	}

	// Pre-check for subnet network, the virtual IP of IPv6 must be established on the basis that the subnet supports
	// IPv6.
	if d.Get("ip_version").(int) == 6 {
		if n.IPv6SubnetId == "" {
			return fmtp.DiagErrorf("The subnet does not support IPv6, please enable IPv6 first.")
		}
		subnetId = n.IPv6SubnetId
	} else {
		subnetId = n.SubnetId
	}

	opts := ports.CreateOpts{
		Name:        d.Get("name").(string),
		DeviceOwner: "neutron:VIP_PORT",
		NetworkId:   networkId,
		FixedIps: []ports.FixedIp{
			{
				SubnetId:  subnetId,
				IpAddress: d.Get("ip_address").(string),
			},
		},
	}

	logp.Printf("[DEBUG] Updating network VIP (%s) with options: %#v", d.Id(), opts)
	vip, err := ports.Create(client, opts)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud network VIP: %s", err)
	}
	logp.Printf("[DEBUG] Waiting for HuaweiCloud network VIP (%s) to become available.", vip.ID)
	d.SetId(vip.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"DOWN", "ACTIVE"},
		Refresh:    waitForNetworkVipStateRefresh(client, vip.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceNetworkingVipRead(ctx, d, meta)
}

func setVipNetworkParams(d *schema.ResourceData, port *ports.Port) error {
	if len(port.FixedIps) > 0 {
		addr := port.FixedIps[0].IpAddress
		var ipVersion int
		if utils.IsIPv4Address(addr) {
			ipVersion = 4
		} else {
			ipVersion = 6
		}
		mErr := multierror.Append(nil,
			d.Set("ip_address", addr),
			d.Set("subnet_id", port.FixedIps[0].SubnetId),
			d.Set("ip_version", ipVersion),
		)
		return mErr.ErrorOrNil()
	}
	return nil
}

// For VIP ports, the status will always be 'DOWN'.
func parseNetworkVipStatus(status string) string {
	if status == "DOWN" {
		return "ACTIVE"
	}
	return status
}

func resourceNetworkingVipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC network v1 client: %s", err)
	}

	vip, err := ports.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "VPC network VIP")
	}

	logp.Printf("[DEBUG] Retrieved VIP %s: %+v", d.Id(), vip)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", vip.Name),
		d.Set("status", parseNetworkVipStatus(vip.Status)),
		d.Set("device_owner", vip.DeviceOwner),
		d.Set("mac_address", vip.MacAddress),
		d.Set("network_id", vip.NetworkId),
		setVipNetworkParams(d, vip),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNetworkingVipUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC network v1 client: %s", err)
	}

	opts := ports.UpdateOpts{
		Name: d.Get("name").(string),
	}
	logp.Printf("[DEBUG] Updating network VIP (%s) with options: %#v", d.Id(), opts)

	_, err = ports.Update(client, d.Id(), opts)
	if err != nil {
		return fmtp.DiagErrorf("Error updating HuaweiCloud networking VIP: %s", err)
	}

	return resourceNetworkingVipRead(ctx, d, meta)
}

func resourceNetworkingVipDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC network v1 client: %s", err)
	}

	err = ports.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Network VIP: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"DOWN", "ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForNetworkVipStateRefresh(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmtp.DiagErrorf("Error deleting HuaweiCloud Network VIP: %s", err)
	}

	d.SetId("")

	return nil
}

func waitForNetworkVipStateRefresh(networkingClient *golangsdk.ServiceClient, vipId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := ports.Get(networkingClient, vipId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				logp.Printf("[DEBUG] The network VIP (%s) has been deleted.", vipId)
				return resp, "DELETED", nil
			}
			return nil, "ERROR", err
		}
		logp.Printf("[DEBUG] The status of the network VIP is: %s", resp.Status)
		return resp, resp.Status, nil
	}
}

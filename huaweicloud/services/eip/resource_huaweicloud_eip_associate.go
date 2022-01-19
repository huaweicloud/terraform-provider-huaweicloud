package eip

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// ResourceEIPAssociate is the ipml for huaweicloud_vpc_eip_associate resource
func ResourceEIPAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEIPAssociateCreate,
		ReadContext:   resourceEIPAssociateRead,
		DeleteContext: resourceEIPAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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

func resourceEIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	publicIP := d.Get("public_ip").(string)
	publicID, err := getEIPByAddress(vpcClient, publicIP)
	if err != nil {
		return fmtp.DiagErrorf("Unable to get ID of public IP %s: %s", publicIP, err)
	}

	portID := d.Get("port_id").(string)
	err = bindPort(vpcClient, publicID, portID, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmtp.DiagErrorf("Error associating EIP %s to port %s: %s", publicID, portID, err)
	}

	d.SetId(publicID)
	return resourceEIPAssociateRead(ctx, d, meta)
}

func resourceEIPAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	eIP, err := eips.Get(vpcClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "EIP")
	}

	d.Set("region", region)
	d.Set("public_ip", eIP.PublicAddress)
	d.Set("port_id", eIP.PortID)

	return nil
}

func resourceEIPAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	portID := d.Get("port_id").(string)
	err = unbindPort(vpcClient, d.Id(), portID, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return fmtp.DiagErrorf("Error disassociating EIP %s from port %s: %s",
			d.Id(), portID, err)
	}

	return nil
}

func getEIPByAddress(client *golangsdk.ServiceClient, address string) (string, error) {
	listOpts := eips.ListOpts{
		PublicIp: []string{address},
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}
	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return "", fmtp.Errorf("Unable to retrieve EIPs: %s ", err)
	}

	if len(allEips) != 1 {
		return "", fmtp.Errorf("unable to determine the ID of %s", address)
	}

	return allEips[0].ID, nil
}

func bindPort(client *golangsdk.ServiceClient, eipID, portID string, timeout time.Duration) error {
	logp.Printf("[DEBUG] Bind EIP %s to port %s", eipID, portID)
	return actionOnPort(client, eipID, portID, timeout)
}

func unbindPort(client *golangsdk.ServiceClient, eipID, portID string, timeout time.Duration) error {
	logp.Printf("[DEBUG] Unbind EIP %s from port: %s", eipID, portID)
	return actionOnPort(client, eipID, "", timeout)
}

func actionOnPort(client *golangsdk.ServiceClient, eipID, portID string, timeout time.Duration) error {
	updateOpts := eips.UpdateOpts{
		PortID: portID,
	}
	_, err := eips.Update(client, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}

	return waitForEIPActive(client, eipID, timeout)
}

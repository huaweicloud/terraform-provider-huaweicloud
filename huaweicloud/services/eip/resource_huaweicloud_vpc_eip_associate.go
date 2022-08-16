package eip

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/chnsz/golangsdk/openstack/networking/v1/ports"
	"github.com/chnsz/golangsdk/pagination"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// ResourceEIPAssociate is the impl for huaweicloud_vpc_eip_associate resource
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
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
				ExactlyOneOf: []string{"port_id"},
				RequiredWith: []string{"network_id"},
			},
			"network_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"port_id"},
			},

			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func eipAssociateRefreshFunc(client *golangsdk.ServiceClient, id string, portId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := eips.Get(client, id).Extract()
		if err != nil {
			return nil, "ERROR", err
		}
		if resp.PortID == "" {
			return resp, "STARTING", nil
		}
		if resp.PortID == portId {
			return resp, "COMPLETED", nil
		}
		return nil, "ERROR", fmt.Errorf("the EIP is already bound to the port %s", resp.PortID)
	}
}

// waitForStateCompleted is a method that continuously queries the resource state and judges whether the resource has
// reached the expected state.
// Parameters:
//
//	context.Context: Context detail.
//	resource.StateRefreshFunc: Function method to get the validation object.
//	time.Duration: Maximum waiting time.
func waitForStateCompleted(ctx context.Context, f resource.StateRefreshFunc, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"STARTING"},
		Target:       []string{"COMPLETED"},
		Refresh:      f,
		Timeout:      t,
		Delay:        5 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceEIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	publicIP := d.Get("public_ip").(string)
	epsID := "all_granted_eps"
	publicID, err := common.GetEipIDbyAddress(vpcClient, publicIP, epsID)
	if err != nil {
		return fmtp.DiagErrorf("Unable to get ID of public IP %s: %s", publicIP, err)
	}

	var portID string
	if v, ok := d.GetOk("port_id"); ok {
		portID = v.(string)
	} else {
		networkID := d.Get("network_id").(string)
		fixedIP := d.Get("fixed_ip").(string)
		portID, err = getPortbyFixedIP(vpcClient, networkID, fixedIP)
		if err != nil {
			return fmtp.DiagErrorf("Unable to get port ID of %s: %s", fixedIP, err)
		}
	}

	// The maximum timeout of excution methods for associate EIP.
	t := d.Timeout(schema.TimeoutCreate)
	err = bindPort(vpcClient, publicID, portID, t)
	if err != nil {
		return fmtp.DiagErrorf("Error associating EIP %s to port %s: %s", publicID, portID, err)
	}

	d.SetId(publicID)
	err = waitForStateCompleted(ctx, eipAssociateRefreshFunc(vpcClient, publicID, portID), t)
	if err != nil {
		return diag.Errorf("error waiting for EIP association to complete: %s", err)
	}

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
		return common.CheckDeletedDiag(d, err, "Error fetching EIP")
	}

	if eIP.PortID == "" {
		diags := diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  fmt.Sprintf("the EIP %s(%s) is not associated", eIP.PublicAddress, d.Id()),
			},
		}

		d.SetId("")
		return diags
	}

	associatedPort, err := ports.Get(vpcClient, eIP.PortID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error fetching port")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("port_id", eIP.PortID),
		d.Set("public_ip", eIP.PublicAddress),
		d.Set("fixed_ip", eIP.PrivateAddress),
		d.Set("network_id", associatedPort.NetworkId),
		d.Set("mac_address", associatedPort.MacAddress),
		d.Set("status", NormalizeEIPStatus(eIP.Status)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting eip associate fields: %s", err)
	}

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

func getPortbyFixedIP(client *golangsdk.ServiceClient, networkID, fixedIP string) (string, error) {
	var portID string

	listOpts := ports.ListOpts{
		NetworkID: networkID,
		FixedIps:  []string{fmt.Sprintf("ip_address=%s", fixedIP)},
	}

	pager := ports.List(client, listOpts)
	err := pager.EachPage(func(page pagination.Page) (b bool, err error) {
		portList, err := ports.ExtractPorts(page)
		if err != nil {
			return false, err
		}
		for _, item := range portList {
			for _, addr := range item.FixedIps {
				if addr.IpAddress == fixedIP {
					portID = item.ID
					return false, nil
				}
			}
		}
		return true, nil
	})

	if err != nil {
		return "", fmtp.Errorf("Unable to list ports: %s", err)
	}

	if portID == "" {
		return "", fmtp.Errorf("can not find %s in subnet %s", fixedIP, networkID)
	}

	return portID, nil
}

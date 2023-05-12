package ecs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceComputeInterfaceAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeInterfaceAttachCreate,
		ReadContext:   resourceComputeInterfaceAttachRead,
		UpdateContext: resourceComputeInterfaceAttachUpdate,
		DeleteContext: resourceComputeInterfaceAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"port_id", "network_id"},
			},
			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"security_group_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"source_dest_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"mac": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func updateInterfacePort(client *golangsdk.ServiceClient, portId string, securityGroupIds []string,
	sourceDestCheckEnabled bool) error {
	opts := ports.UpdateOpts{
		AllowedAddressPairs: nil,
		SecurityGroups:      &securityGroupIds,
	}
	if !sourceDestCheckEnabled {
		// Update the allowed-address-pairs of the port to 1.1.1.1/0
		// to disable the source/destination check
		portpairs := []ports.AddressPair{
			{
				IPAddress: "1.1.1.1/0",
			},
		}
		opts.AllowedAddressPairs = &portpairs
	}

	_, err := ports.Update(client, portId, opts).Extract()
	return err
}

func resourceComputeInterfaceAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}
	nicClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating VPC v2.0 client: %s", err)
	}

	var portId string
	if v, ok := d.GetOk("port_id"); ok {
		portId = v.(string)
	}

	var networkId string
	if v, ok := d.GetOk("network_id"); ok {
		networkId = v.(string)
	}

	// For some odd reason the API takes an array of IPs, but you can only have one element in the array.
	var fixedIPs []attachinterfaces.FixedIP
	if v, ok := d.GetOk("fixed_ip"); ok {
		fixedIPs = append(fixedIPs, attachinterfaces.FixedIP{IPAddress: v.(string)})
	}

	attachOpts := attachinterfaces.CreateOpts{
		PortID:    portId,
		NetworkID: networkId,
		FixedIPs:  fixedIPs,
	}

	log.Printf("[DEBUG] compute interface attach options: %#v", attachOpts)
	instanceId := d.Get("instance_id").(string)
	attachment, err := attachinterfaces.Create(computeClient, instanceId, attachOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	portID := attachment.PortID
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeInterfaceAttachAttachFunc(computeClient, instanceId, portID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error creating attaching interface to compute instance %s: %s", instanceId, err)
	}

	// Use the instance ID and port ID as the resource ID.
	id := fmt.Sprintf("%s/%s", instanceId, portID)
	d.SetId(id)

	var (
		securityGroupIds       = d.Get("security_group_ids").([]interface{})
		sourceDestCheckEnabled = d.Get("source_dest_check").(bool)
	)
	err = updateInterfacePort(nicClient, portID, utils.ExpandToStringList(securityGroupIds), sourceDestCheckEnabled)
	if err != nil {
		return diag.Errorf("error updating VPC port (%s): %s", portID, err)
	}

	return resourceComputeInterfaceAttachRead(ctx, d, meta)
}

func resourceComputeInterfaceAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	computeClient, err := cfg.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}
	networkingClient, err := cfg.NetworkingV2Client(region)
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	instanceId, portId, err := computeInterfaceAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving compute interface attach")
	}

	var (
		securitygroups []string
		ipAddress      string
		macAddress     string
		sdCheck        bool
	)

	if len(attachment.FixedIPs) > 0 {
		ipAddress = attachment.FixedIPs[0].IPAddress
	}
	if port, err := ports.Get(networkingClient, attachment.PortID).Extract(); err == nil {
		macAddress = port.MACAddress
		securitygroups = port.SecurityGroups
		sdCheck = len(port.AllowedAddressPairs) == 0
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceId),
		d.Set("port_id", attachment.PortID),
		d.Set("network_id", attachment.NetID),
		d.Set("fixed_ip", ipAddress),
		d.Set("mac", macAddress),
		d.Set("security_group_ids", securitygroups),
		d.Set("source_dest_check", sdCheck),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComputeInterfaceAttachUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	nicClient, err := cfg.NetworkingV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating networking client: %s", err)
	}

	var (
		portId                 = d.Get("port_id").(string)
		securityGroupIds       = d.Get("security_group_ids").([]interface{})
		sourceDestCheckEnabled = d.Get("source_dest_check").(bool)
	)
	err = updateInterfacePort(nicClient, portId, utils.ExpandToStringList(securityGroupIds), sourceDestCheckEnabled)
	if err != nil {
		return diag.Errorf("error updating VPC port (%s): %s", portId, err)
	}

	return resourceComputeInterfaceAttachRead(ctx, d, meta)
}

func resourceComputeInterfaceAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating compute client: %s", err)
	}

	instanceId, portId, err := computeInterfaceAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeInterfaceAttachDetachFunc(computeClient, instanceId, portId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error detaching interface from compute instance %s: %s", instanceId, err)
	}

	return nil
}

func computeInterfaceAttachAttachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		va, err := attachinterfaces.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "ATTACHING", nil
			}
			return va, "", err
		}

		return va, "ATTACHED", nil
	}
}

func computeInterfaceAttachDetachFunc(
	computeClient *golangsdk.ServiceClient, instanceId, attachmentId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to detach interface attachment %s from instance %s",
			attachmentId, instanceId)

		va, err := attachinterfaces.Get(computeClient, instanceId, attachmentId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}
			return va, "", err
		}

		err = attachinterfaces.Delete(computeClient, instanceId, attachmentId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return va, "DETACHED", nil
			}

			if _, ok := err.(golangsdk.ErrDefault400); ok {
				return nil, "", nil
			}

			return nil, "", err
		}

		log.Printf("[DEBUG] compute interface attachment %s is still active.", attachmentId)
		return nil, "", nil
	}
}

func computeInterfaceAttachParseID(id string) (instanceID, portID string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		err = fmt.Errorf("unable to parse the resource ID, must be <instance_id>/<port_id> format")
		return
	}

	instanceID = idParts[0]
	portID = idParts[1]
	return
}

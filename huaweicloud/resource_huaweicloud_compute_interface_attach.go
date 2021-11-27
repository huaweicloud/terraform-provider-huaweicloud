package huaweicloud

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/attachinterfaces"
	"github.com/chnsz/golangsdk/openstack/networking/v2/ports"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeInterfaceAttachV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeInterfaceAttachV2Create,
		Read:   resourceComputeInterfaceAttachV2Read,
		Update: resourceComputeInterfaceAttachV2Update,
		Delete: resourceComputeInterfaceAttachV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceComputeInterfaceAttachV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

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

	logp.Printf("[DEBUG] huaweicloud_compute_interface_attach attach options: %#v", attachOpts)
	attachment, err := attachinterfaces.Create(computeClient, instanceId, attachOpts).Extract()
	if err != nil {
		return err
	}

	portID := attachment.PortID
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeInterfaceAttachV2AttachFunc(computeClient, instanceId, portID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmtp.Errorf("Error creating huaweicloud_compute_interface_attach %s: %s", instanceId, err)
	}

	// Use the instance ID and port ID as the resource ID.
	id := fmt.Sprintf("%s/%s", instanceId, portID)

	logp.Printf("[DEBUG] Created huaweicloud_compute_interface_attach %s: %#v", id, attachment)

	d.SetId(id)

	if sourceDestCheck := d.Get("source_dest_check").(bool); !sourceDestCheck {
		nicClient, err := config.NetworkingV2Client(GetRegion(d, config))
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
		}
		if err := disableSourceDestCheck(nicClient, portID); err != nil {
			return fmtp.Errorf("Error disable source dest check on port(%s) of instance(%s) failed: %s", portID, d.Id(), err)
		}
	}

	return resourceComputeInterfaceAttachV2Read(d, meta)
}

func resourceComputeInterfaceAttachV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	instanceId, portId, err := computeInterfaceAttachV2ParseID(d.Id())
	if err != nil {
		return err
	}

	attachment, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving huaweicloud_compute_interface_attach")
	}

	logp.Printf("[DEBUG] Retrieved huaweicloud_compute_interface_attach %s: %#v", d.Id(), attachment)

	d.Set("instance_id", instanceId)
	d.Set("port_id", attachment.PortID)
	d.Set("network_id", attachment.NetID)
	d.Set("region", GetRegion(d, config))

	if len(attachment.FixedIPs) > 0 {
		firstAddress := attachment.FixedIPs[0].IPAddress
		d.Set("fixed_ip", firstAddress)
	}

	if port, err := ports.Get(networkingClient, attachment.PortID).Extract(); err == nil {
		d.Set("source_dest_check", len(port.AllowedAddressPairs) == 0)
		d.Set("mac", port.MACAddress)
	}

	return nil
}

func resourceComputeInterfaceAttachV2Update(d *schema.ResourceData, meta interface{}) error {
	var err error

	config := meta.(*config.Config)
	nicClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	portID := d.Get("port_id").(string)
	if sourceDestCheck := d.Get("source_dest_check").(bool); !sourceDestCheck {
		err = disableSourceDestCheck(nicClient, portID)
	} else {
		err = enableSourceDestCheck(nicClient, portID)
	}

	if err != nil {
		return fmtp.Errorf("Error updating source_dest_check on port(%s) of instance(%s) failed: %s", portID, d.Id(), err)
	}

	return resourceComputeInterfaceAttachV2Read(d, meta)
}

func resourceComputeInterfaceAttachV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instanceId, portId, err := computeInterfaceAttachV2ParseID(d.Id())
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeInterfaceAttachV2DetachFunc(computeClient, instanceId, portId),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 5 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmtp.Errorf("Error detaching huaweicloud_compute_interface_attach %s: %s", d.Id(), err)
	}

	return nil
}

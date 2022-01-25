package huaweicloud

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	bandwidthsv1 "github.com/chnsz/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/chnsz/golangsdk/openstack/networking/v2/bandwidths"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

const publicIPv6Type string = "5_dualStack"

func ResourceComputeFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeEIPAssociateCreate,
		Read:   resourceComputeEIPAssociateRead,
		Delete: resourceComputeEIPAssociateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceComputeEIPAssociateImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			"public_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPv4Address,
				ExactlyOneOf: []string{"bandwidth_id"},
			},
			"bandwidth_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"fixed_ip"},
			},
			"fixed_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateFunc:     validation.IsIPAddress,
				DiffSuppressFunc: utils.SuppressComputedFixedWhenFloatingIp,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeEIPAssociateCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var publicID string
	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	if _, ok := d.GetOk("bandwidth_id"); ok {
		// fixed_ip must be a valid IPv6 address when combining with bandwidth_id
		if utils.IsIPv4Address(fixedIP) {
			return fmtp.Errorf("the fixed_ip must be a valid IPv6 address, got: %s", fixedIP)
		}
	}

	// get port id
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error getting port id of compute instance: %s", err)
	}

	if v, ok := d.GetOk("public_ip"); ok {
		vpcClient, err := config.NetworkingV1Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud VPC client: %s", err)
		}

		// get EIP id
		eipAddr := v.(string)
		publicID = eipAddr
		pAddress, err := getFloatingIPbyAddress(vpcClient, eipAddr)
		if err != nil {
			return fmtp.Errorf("Error getting EIP: %s", err)
		}

		eipID := pAddress.ID
		err = bindPortToEIP(vpcClient, eipID, portID)
		if err != nil {
			return fmtp.Errorf("Error associating port %s to EIP: %s", portID, err)
		}
	} else {
		bwClient, err := config.NetworkingV2Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating bandwidth v2.0 client: %s", err)
		}

		bwID := d.Get("bandwidth_id").(string)
		publicID = bwID
		err = insertPortToBandwidth(bwClient, bwID, portID)
		if err != nil {
			return fmtp.Errorf("Error associating IPv6 port %s to bandwidth: %s", portID, err)
		}
	}

	id := fmt.Sprintf("%s/%s/%s", publicID, instanceID, privateIP)
	d.SetId(id)

	return resourceComputeEIPAssociateRead(d, meta)
}

func resourceComputeEIPAssociateRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	vpcClient, err := config.NetworkingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating VPC client: %s", err)
	}

	var associated bool
	var publicID string
	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id of compute instance
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error getting port id of compute instance: %s", err)
	}

	if v, ok := d.GetOk("public_ip"); ok {
		eipAddr := v.(string)
		publicID = eipAddr
		eipInfo, err := getFloatingIPbyAddress(vpcClient, eipAddr)
		if err != nil {
			if eipInfo != nil {
				logp.Printf("[WARN] can not find the EIP by %s", eipAddr)
				d.SetId("")
				return nil
			}
			return err
		}

		if eipInfo.PortID == portID {
			associated = true
		}
	} else {
		bwID := d.Get("bandwidth_id").(string)
		publicID = bwID
		band, err := bandwidthsv1.Get(vpcClient, bwID).Extract()
		if err != nil {
			return common.CheckDeleted(d, err, "bandwidth")
		}

		for _, ipInfo := range band.PublicipInfo {
			if ipInfo.PublicipId == portID {
				associated = true
				break
			}
		}
	}

	if !associated {
		logp.Printf("[WARN] the resource is not associated with the specified EIP or bandwidth")
		d.SetId("")
	}

	id := fmt.Sprintf("%s/%s/%s", publicID, instanceID, privateIP)
	d.SetId(id)

	// Set the attributes pulled from the composed resource ID
	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", privateIP)
	d.Set("port_id", portID)
	d.Set("region", region)

	return nil
}

func resourceComputeEIPAssociateDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id of compute instance
	portID, _, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error getting port id of compute instance: %s", err)
	}

	if v, ok := d.GetOk("public_ip"); ok {
		vpcClient, err := config.NetworkingV1Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating HuaweiCloud VPC client: %s", err)
		}

		eipAddr := v.(string)
		eipInfo, err := getFloatingIPbyAddress(vpcClient, eipAddr)
		if err != nil {
			return fmtp.Errorf("Error getting EIP: %s", err)
		}

		err = unbindPortFromEIP(vpcClient, eipInfo.ID, portID)
		if err != nil {
			return fmtp.Errorf("Error disassociating Floating IP: %s", err)
		}
	} else {
		bwClient, err := config.NetworkingV2Client(region)
		if err != nil {
			return fmtp.Errorf("Error creating bandwidth v2.0 client: %s", err)
		}

		bwID := d.Get("bandwidth_id").(string)
		err = removePortFromBandwidth(bwClient, bwID, portID)
		if err != nil {
			return fmtp.Errorf("Error associating IPv6 port %s to bandwidth: %s", portID, err)
		}
	}

	return nil
}

func getComputeInstancePortIDbyFixedIP(d *schema.ResourceData, config *config.Config,
	instanceId, fixedIP string) (string, string, error) {

	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return "", "", fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instance, err := cloudservers.Get(computeClient, instanceId).Extract()
	if err != nil {
		return "", "", err
	}

	var portID, privateIP string
	for _, networkAddresses := range instance.Addresses {
		for _, address := range networkAddresses {
			if address.Type == "fixed" {
				if fixedIP == "" || address.Addr == fixedIP {
					portID = address.PortID
					privateIP = address.Addr
					break
				}
			}
		}
		if portID != "" {
			break
		}
	}

	if portID == "" {
		return "", "", fmt.Errorf("the port ID does not exist")
	}
	return portID, privateIP, nil
}

func getFloatingIPbyAddress(client *golangsdk.ServiceClient, floatingIP string) (*eips.PublicIp, error) {
	listOpts := &eips.ListOpts{
		PublicIp: []string{floatingIP},
	}

	pages, err := eips.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return nil, fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) != 1 {
		return &eips.PublicIp{}, fmtp.Errorf("can not find the EIP by %s", floatingIP)
	}

	return &allEips[0], nil
}

func insertPortToBandwidth(client *golangsdk.ServiceClient, bwID, portID string) error {
	insertOpts := bandwidths.BandWidthInsertOpts{
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: publicIPv6Type,
			},
		},
	}

	logp.Printf("[DEBUG] Insert port %s to bandwidth %s", portID, bwID)
	_, err := bandwidths.Insert(client, bwID, insertOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error inserting %s into bandwidth %s: %s", portID, bwID, err)
	}
	return nil
}

func removePortFromBandwidth(client *golangsdk.ServiceClient, bwID, portID string) error {
	removalChargeMode := "bandwidth"
	removalSize := 5
	removeOpts := bandwidths.BandWidthRemoveOpts{
		ChargeMode: removalChargeMode,
		Size:       &removalSize,
		PublicipInfo: []bandwidths.PublicIpInfoID{
			{
				PublicIPID:   portID,
				PublicIPType: publicIPv6Type,
			},
		},
	}

	logp.Printf("[DEBUG] Remove port %s from bandwidth %s", portID, bwID)
	err := bandwidths.Remove(client, bwID, removeOpts).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error removing %s from bandwidth: %s", portID, err)
	}
	return nil
}

func bindPortToEIP(client *golangsdk.ServiceClient, eipID, portID string) error {
	logp.Printf("[DEBUG] Bind port %s to EIP %s", portID, eipID)
	return actionOnPort(client, eipID, portID)
}

func unbindPortFromEIP(client *golangsdk.ServiceClient, eipID, portID string) error {
	logp.Printf("[DEBUG] Unbind port %s from EIP: %s", portID, eipID)
	return actionOnPort(client, eipID, "")
}

func actionOnPort(client *golangsdk.ServiceClient, eipID, portID string) error {
	updateOpts := eips.UpdateOpts{
		PortID: portID,
	}
	_, err := eips.Update(client, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}

	return nil
}

func parseComputeFloatingIPAssociateID(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 3 && len(idParts) != 2 {
		return "", "", "",
			fmtp.Errorf("Unable to parse the resource ID, must be <eip address or bandwidth_id>/<instance_id>/<fixed_ip> format")
	}

	publicID := idParts[0]
	instanceID := idParts[1]

	var fixedIP string
	if len(idParts) == 3 {
		fixedIP = idParts[2]
	}

	return publicID, instanceID, fixedIP, nil
}

func resourceComputeEIPAssociateImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	publicID, instanceID, fixedIP, err := parseComputeFloatingIPAssociateID(d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", fixedIP)
	parsedIP := net.ParseIP(publicID)
	if parsedIP != nil {
		d.Set("public_ip", publicID)
	} else {
		d.Set("bandwidth_id", publicID)
	}

	return []*schema.ResourceData{d}, nil
}

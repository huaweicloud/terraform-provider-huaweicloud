package huaweicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/chnsz/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceComputeFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFloatingIPAssociateV2Create,
		Read:   resourceComputeFloatingIPAssociateV2Read,
		Delete: resourceComputeFloatingIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressComputedFixedWhenFloatingIp,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeFloatingIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	floatingIP := d.Get("public_ip").(string)
	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error get port id of compute instance: %s", err)
	}

	// get EIP id
	pAddress, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}
	floatingID := pAddress.ID

	// Associate EIP to compute instance
	associateOpts := eips.UpdateOpts{
		PortID: portID,
	}
	logp.Printf("[DEBUG] EIP Associate Options: %#v", associateOpts)

	_, err = eips.Update(vpcClient, floatingID, associateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error associating EIP: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceID, privateIP)
	d.SetId(id)

	logp.Printf("[DEBUG] Waiting for eip associate to instance %s", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{fmt.Sprintf("floating/%s", floatingIP)},
		Refresh:    resourceComputeFloatingIPAssociateStateRefreshFunc(d, config, instanceID, floatingIP),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmtp.Errorf(
			"Error waiting for eip associate to instance(%s): %s",
			id, err)
	}

	return resourceComputeFloatingIPAssociateV2Read(d, meta)
}

func resourceComputeFloatingIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)

	floatingIP, instanceID, fixedIP, err := parseComputeFloatingIPAssociateId(d.Id())
	if err != nil {
		return err
	}

	eipInfo, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}

	// get port id of compute instance
	portID, privateIP, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error get port id of compute instance: %s", err)
	}

	var associated bool
	if eipInfo.PortID == portID {
		associated = true
	}

	if !associated {
		d.SetId("")
	}

	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceID, privateIP)
	d.SetId(id)

	// Set the attributes pulled from the composed resource ID
	d.Set("public_ip", floatingIP)
	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", privateIP)
	d.Set("port_id", portID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeFloatingIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud VPC client: %s", err)
	}

	// get EIP id
	floatingIP := d.Get("public_ip").(string)
	eipInfo, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}
	floatingID := eipInfo.ID

	disassociateOpts := eips.UpdateOpts{PortID: ""}
	logp.Printf("[DEBUG] EIP Disssociate Options: %#v", disassociateOpts)

	_, err = eips.Update(vpcClient, floatingID, disassociateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error disassociating Floating IP: %s", err)
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

func getFloatingIPbyAddress(d *schema.ResourceData, config *config.Config, floatingIP string) (eips.PublicIp, error) {
	vpcClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return eips.PublicIp{}, fmtp.Errorf("Error creating VPC client: %s", err)
	}

	listOpts := &eips.ListOpts{
		PublicIp: []string{floatingIP},
	}

	pages, err := eips.List(vpcClient, listOpts).AllPages()
	if err != nil {
		return eips.PublicIp{}, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return eips.PublicIp{}, fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) != 1 {
		return eips.PublicIp{}, fmtp.Errorf("can not find the EIP %s", floatingIP)
	}

	return allEips[0], nil
}

func resourceComputeFloatingIPAssociateStateRefreshFunc(d *schema.ResourceData, config *config.Config, instanceId, floatingIP string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		computeClient, err := config.ComputeV1Client(GetRegion(d, config))
		if err != nil {
			return computeClient, "", fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
		}

		instance, err := cloudservers.Get(computeClient, instanceId).Extract()
		if err != nil {
			return instance, "", err
		}

		for _, networkAddresses := range instance.Addresses {
			for _, address := range networkAddresses {
				if address.Type == "floating" && address.Addr == floatingIP {
					return instance, fmt.Sprintf("floating/%s", floatingIP), nil
				}
			}
		}

		return instance, "", nil
	}
}

func parseComputeFloatingIPAssociateId(id string) (string, string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 3 && len(idParts) != 2 {
		return "", "", "", fmtp.Errorf("Unable to parse the resource ID, must be <eip>/<instance_id>/<fixed_ip> format")
	}

	floatingIP := idParts[0]
	instanceID := idParts[1]

	var fixedIP string
	if len(idParts) == 3 {
		fixedIP = idParts[2]
	}

	return floatingIP, instanceID, fixedIP, nil
}

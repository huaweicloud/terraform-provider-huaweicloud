package huaweicloud

import (
	"strings"
	"time"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
		},
	}
}

func resourceComputeFloatingIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	floatingIP := d.Get("public_ip").(string)
	instanceID := d.Get("instance_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	// get port id
	portID, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error get port id of compute instance: %s", err)
	}

	// get floating_ip id
	Eip, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}
	floatingID := Eip.ID

	// Associate Eip to compute instance
	associateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}
	logp.Printf("[DEBUG] Associate Options: %#v", associateOpts)

	_, err = floatingips.Update(networkingClient, floatingID, associateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error associating Floating IP: %s", err)
	}

	id := fmtp.Sprintf("%s/%s/%s", floatingIP, instanceID, fixedIP)
	d.SetId(id)

	logp.Printf("[DEBUG] Waiting for eip associate to instance %s", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{fmtp.Sprintf("floating/%s", floatingIP), fmtp.Sprintf("fixed/%s", floatingIP)},
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

	Eip, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}

	// get port id of compute instance
	portID, err := getComputeInstancePortIDbyFixedIP(d, config, instanceID, fixedIP)
	if err != nil {
		return fmtp.Errorf("Error get port id of compute instance: %s", err)
	}

	var associated bool
	if Eip.PortID == portID {
		associated = true
	}

	if !associated {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	d.Set("public_ip", floatingIP)
	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", fixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeFloatingIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	// get floating_ip id
	floatingIP := d.Get("public_ip").(string)
	Eip, err := getFloatingIPbyAddress(d, config, floatingIP)
	if err != nil {
		return fmtp.Errorf("Error get eip: %s", err)
	}
	floatingID := Eip.ID

	disassociateOpts := floatingips.UpdateOpts{}
	logp.Printf("[DEBUG] Disssociate Options: %#v", disassociateOpts)

	_, err = floatingips.Update(networkingClient, floatingID, disassociateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error disassociating Floating IP: %s", err)
	}

	return nil
}

func getComputeInstancePortIDbyFixedIP(d *schema.ResourceData, config *config.Config, instanceId, fixedIP string) (string, error) {
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return "", fmtp.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	instance, err := cloudservers.Get(computeClient, instanceId).Extract()
	if err != nil {
		return "", err
	}

	var portId string
	if fixedIP != "" {
		for _, networkAddresses := range instance.Addresses {
			for _, address := range networkAddresses {
				if address.Type == "fixed" && address.Addr == fixedIP {
					portId = address.PortID
					break
				}
			}
			if portId != "" {
				break
			}
		}
	} else {
		// return the first port if fixedIP not specified
		for _, networkAddresses := range instance.Addresses {
			for _, address := range networkAddresses {
				portId = address.PortID
				break
			}
			if portId != "" {
				break
			}
		}
	}
	return portId, nil
}

func getFloatingIPbyAddress(d *schema.ResourceData, config *config.Config, floatingIP string) (eips.PublicIp, error) {
	networkingV1Client, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return eips.PublicIp{}, fmtp.Errorf("Error creating networking client: %s", err)
	}

	listOpts := &eips.ListOpts{
		PublicIp: floatingIP,
	}

	pages, err := eips.List(networkingV1Client, listOpts).AllPages()
	if err != nil {
		return eips.PublicIp{}, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return eips.PublicIp{}, fmtp.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return eips.PublicIp{}, fmtp.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return eips.PublicIp{}, fmtp.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
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
					return instance, fmtp.Sprintf("floating/%s", floatingIP), nil
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

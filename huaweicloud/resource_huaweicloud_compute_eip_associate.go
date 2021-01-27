package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/ecs/v1/cloudservers"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/layer3/floatingips"
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

			"floating_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"public_ip"},
				Deprecated:    "use public_ip instead",
			},
			"public_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"floating_ip"},
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fixed_ip": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				DiffSuppressFunc: suppressComputedFixedWhenFloatingIp,
			},
		},
	}
}

func resourceComputeFloatingIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	floating_ip, fip_ok := d.GetOk("floating_ip")
	public_ip, pip_ok := d.GetOk("public_ip")
	if !fip_ok && !pip_ok {
		return fmt.Errorf("One of floating_ip or public_ip must be configured")
	}
	fixedIP := d.Get("fixed_ip").(string)
	instanceId := d.Get("instance_id").(string)

	var floatingIP string
	if fip_ok {
		floatingIP = floating_ip.(string)
	} else {
		floatingIP = public_ip.(string)
	}

	config := meta.(*Config)

	// get port id
	portId, err := resourceComputeFloatingIPAssociateV2GetPortId(d, config, instanceId, fixedIP)

	if err != nil {
		return fmt.Errorf("Error get port id of compute instance: %s", err)
	}

	// get floating_ip id
	Eip, err := resourceComputeFloatingIPAssociateV2GetEip(d, config, floatingIP)
	if err != nil {
		return fmt.Errorf("Error get eip: %s", err)
	}

	floatingIpId := Eip.ID

	// Associate Eip to compute instance
	associateOpts := floatingips.UpdateOpts{
		PortID: &portId,
	}
	log.Printf("[DEBUG] Associate Options: %#v", associateOpts)

	networkingClientAssociate, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	_, err = floatingips.Update(networkingClientAssociate, floatingIpId, associateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error associating Floating IP: %s", err)
	}

	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceId, fixedIP)
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for eip associate to instance %s", id)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{fmt.Sprintf("floating/%s", floatingIP), fmt.Sprintf("fixed/%s", floatingIP)},
		Refresh:    resourceComputeFloatingIPAssociateStateRefreshFunc(d, config, instanceId, floatingIP),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf(
			"Error waiting for eip associate to instance(%s): %s",
			id, err)
	}

	return resourceComputeFloatingIPAssociateV2Read(d, meta)
}

func resourceComputeFloatingIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	floatingIP, instanceId, fixedIP, err := parseComputeFloatingIPAssociateId(d.Id())
	if err != nil {
		return err
	}

	Eip, err := resourceComputeFloatingIPAssociateV2GetEip(d, config, floatingIP)
	if err != nil {
		return fmt.Errorf("Error get eip: %s", err)
	}

	// get port id of compute instance
	portId, err := resourceComputeFloatingIPAssociateV2GetPortId(d, config, instanceId, fixedIP)

	if err != nil {
		return fmt.Errorf("Error get port id of compute instance: %s", err)
	}

	var associated bool

	if Eip.PortID == portId {
		associated = true
	}

	if !associated {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	if _, ok := d.GetOk("floating_ip"); ok {
		d.Set("floating_ip", floatingIP)
	} else {
		d.Set("public_ip", floatingIP)
	}
	d.Set("instance_id", instanceId)
	d.Set("fixed_ip", fixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeFloatingIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	var floatingIP string
	if _, ok := d.GetOk("floating_ip"); ok {
		floatingIP = d.Get("floating_ip").(string)
	} else {
		floatingIP = d.Get("public_ip").(string)
	}

	// get floating_ip id
	Eip, err := resourceComputeFloatingIPAssociateV2GetEip(d, config, floatingIP)
	if err != nil {
		return fmt.Errorf("Error get eip: %s", err)
	}
	floatingIpId := Eip.ID

	disassociateOpts := floatingips.UpdateOpts{}
	log.Printf("[DEBUG] Disssociate Options: %#v", disassociateOpts)

	networkingClientAssociate, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud network client: %s", err)
	}

	_, err = floatingips.Update(networkingClientAssociate, floatingIpId, disassociateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error disassociating Floating IP: %s", err)
	}

	return nil
}

func resourceComputeFloatingIPAssociateV2GetPortId(d *schema.ResourceData, config *Config, instanceId, fixedIP string) (string, error) {
	computeClient, err := config.ComputeV1Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
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

func resourceComputeFloatingIPAssociateV2GetEip(d *schema.ResourceData, config *Config, floatingIP string) (eips.PublicIp, error) {
	networkingClientFloatingIp, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return eips.PublicIp{}, fmt.Errorf("Error creating networking client: %s", err)
	}

	listOpts := &eips.ListOpts{
		PublicIp: floatingIP,
	}

	pages, err := eips.List(networkingClientFloatingIp, listOpts).AllPages()
	if err != nil {
		return eips.PublicIp{}, err
	}

	allEips, err := eips.ExtractPublicIPs(pages)
	if err != nil {
		return eips.PublicIp{}, fmt.Errorf("Unable to retrieve eips: %s ", err)
	}

	if len(allEips) < 1 {
		return eips.PublicIp{}, fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allEips) > 1 {
		return eips.PublicIp{}, fmt.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	return allEips[0], nil
}

func resourceComputeFloatingIPAssociateStateRefreshFunc(d *schema.ResourceData, config *Config, instanceId, floatingIP string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		computeClient, err := config.ComputeV1Client(GetRegion(d, config))
		if err != nil {
			return computeClient, "", fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
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
	if len(idParts) < 3 {
		return "", "", "", fmt.Errorf("Unable to determine floating ip association ID")
	}

	floatingIP := idParts[0]
	instanceId := idParts[1]
	fixedIP := idParts[2]

	return floatingIP, instanceId, fixedIP, nil
}

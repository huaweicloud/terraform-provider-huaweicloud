package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/ports"
)

func resourceIecVipV1() *schema.Resource {

	return &schema.Resource{
		Create: resourceIecVIPV1Create,
		Read:   resourceIecVIPV1Read,
		Delete: resourceIecVIPV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"subnet_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ips": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceIecVIPV1Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	createOpts := ports.CreateOpts{
		NetworkId:   d.Get("subnet_id").(string),
		DeviceOwner: "neutron:VIP_PORT",
	}

	p, err := ports.Create(iecClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC port: %s", err)
	}
	log.Printf("[INFO] Network ID: %s", p.ID)
	log.Printf("[DEBUG] Waiting for HuaweiCloud IEC Port (%s) to become available.", p.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitingForIECVIPActive(iecClient, p.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForState()
	d.SetId(p.ID)
	return resourceIecVIPV1Read(d, meta)
}

func resourceIecVIPV1Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	n, err := ports.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving Huaweicloud IEC VPC")
	}

	d.Set("name", n.Name)
	d.Set("subnet_id", n.NetworkID)
	d.Set("mac_address", n.MacAddress)

	ipsSet := make([]map[string]interface{}, len(n.FixedIPs))
	for index, ipObj := range n.FixedIPs {
		ipsSet[index] = map[string]interface{}{
			"subnet_id":  ipObj.SubnetId,
			"ip_address": ipObj.IpAddress,
		}
	}
	d.Set("fixed_ips", ipsSet)

	return nil
}

func resourceIecVIPV1Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitingForIECVIPDelete(iecClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud IEC Network: %s", err)
	}
	d.SetId("")
	return nil
}

func waitingForIECVIPActive(client *golangsdk.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		p, err := ports.Get(client, portID).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud Neutron Port: %+v", p)
		if p.Status == "DOWN" || p.Status == "ACTIVE" {
			return p, "ACTIVE", nil
		}

		return p, p.Status, nil
	}
}

func waitingForIECVIPDelete(client *golangsdk.ServiceClient, portID string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud IEC Port %s", portID)
		port, err := ports.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Port %s", portID)
				return port, "DELETED", nil
			}
			return port, "ACTIVATE", err
		}
		err = ports.Delete(client, portID).ExtractErr()

		// remote service will return code 204 when delete success
		if err == nil {
			log.Printf("[DEBUG] Successfully deleted HuaweiCloud IEC Port %s", portID)
			return port, "DELETED", nil
		}

		log.Printf("[DEBUG] HuaweiCloud IEC Port %s still active.\n", portID)
		return port, "ACTIVE", nil
	}
}

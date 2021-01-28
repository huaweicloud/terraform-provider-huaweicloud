package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/huaweicloud/golangsdk"
	iec_common "github.com/huaweicloud/golangsdk/openstack/iec/v1/common"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/ports"
)

func resourceIecVipV1() *schema.Resource {

	return &schema.Resource{
		Create: resourceIecVIPV1Create,
		Update: resourceIecVIPV1Update,
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
			"port_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"mac_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fixed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allowed_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func getIecVipPortIDs(d *schema.ResourceData) []string {
	rawPortIDs := d.Get("port_ids").([]interface{})
	portids := make([]string, len(rawPortIDs))
	for i, raw := range rawPortIDs {
		portids[i] = raw.(string)
	}
	return portids
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

	// associate ports with the vip
	portids := getIecVipPortIDs(d)
	if len(portids) > 0 {
		if err = updateIecVipAssociate(iecClient, p.ID, portids); err != nil {
			return err
		}
	}

	return resourceIecVIPV1Read(d, meta)
}

func resourceIecVIPV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	vip, err := ports.Get(iecClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving Huaweicloud IEC VPC")
	}

	d.Set("name", vip.Name)
	d.Set("subnet_id", vip.NetworkID)
	d.Set("mac_address", vip.MacAddress)

	allIPs := make([]string, len(vip.FixedIPs))
	for i, ipObj := range vip.FixedIPs {
		allIPs[i] = ipObj.IpAddress
	}
	d.Set("fixed_ips", allIPs)

	allPortAddrs := make([]string, len(vip.AllowedAddressPairs))
	for i, pair := range vip.AllowedAddressPairs {
		allPortAddrs[i] = pair.IpAddress
	}
	d.Set("allowed_addresses", allPortAddrs)

	return nil
}

func resourceIecVIPV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	portids := getIecVipPortIDs(d)
	if err = updateIecVipAssociate(iecClient, d.Id(), portids); err != nil {
		return err
	}

	return resourceIecVIPV1Read(d, meta)
}

func resourceIecVIPV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iecClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IEC client: %s", err)
	}

	if len(getIecVipPortIDs(d)) > 0 {
		// disassociate ports
		if err := updateIecVipAssociate(iecClient, d.Id(), []string{}); err != nil {
			return err
		}
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

func updateIecVipAssociate(client *golangsdk.ServiceClient, vipID string, portIDs []string) error {
	allAddrs := make([]string, len(portIDs))
	action := "associate"
	if len(portIDs) == 0 {
		action = "disassociate"
	}

	// check the port id and get ip address
	for i, portid := range portIDs {
		port, err := ports.Get(client, portid).Extract()
		if err != nil {
			return fmt.Errorf("Error fetching port %s: %s", portid, err)
		}

		if len(port.FixedIPs) > 0 {
			allAddrs[i] = port.FixedIPs[0].IpAddress
		} else {
			return fmt.Errorf("port %s has no ip address, Error associate it", portid)
		}
	}

	// construct allowed address pairs
	allowedPairs := make([]iec_common.AllowedAddressPair, len(allAddrs))
	for i, addr := range allAddrs {
		allowedPairs[i] = iec_common.AllowedAddressPair{
			IpAddress: addr,
		}
	}
	// associate/disassociate ports with the vip
	associateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &allowedPairs,
	}
	log.Printf("[DEBUG] VIP %s %s with options: %#v", action, vipID, associateOpts)
	_, err := ports.Update(client, vipID, associateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error %s vip: %s", action, err)
	}

	// Update the allowed-address-pairs of the port to 1.1.1.1/0
	// to disable the source/destination check
	portpairs := make([]iec_common.AllowedAddressPair, 1)
	portpairs[0] = iec_common.AllowedAddressPair{
		IpAddress: "1.1.1.1/0",
	}
	portUpdateOpts := ports.UpdateOpts{
		AllowedAddressPairs: &portpairs,
	}

	for _, portid := range portIDs {
		_, err = ports.Update(client, portid, portUpdateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error update port %s: %s", portid, err)
		}
	}

	return nil
}

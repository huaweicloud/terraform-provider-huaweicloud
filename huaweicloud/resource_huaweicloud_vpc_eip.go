package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/golangsdk/openstack/networking/v1/eips"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func ResourceVpcEIPV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcEIPV1Create,
		Read:   resourceVpcEIPV1Read,
		Update: resourceVpcEIPV1Update,
		Delete: resourceVpcEIPV1Delete,
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
			"publicip": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"ip_address": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"port_id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"bandwidth": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"share_type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"charge_mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceVpcEIPV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	createOpts := EIPCreateOpts{
		eips.ApplyOpts{
			IP:        resourcePublicIP(d),
			Bandwidth: resourceBandWidth(d),
		},
		MapValueSpecs(d),
	}

	epsID := GetEnterpriseProjectID(d, config)

	if epsID != "" {
		createOpts.EnterpriseProjectID = epsID
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	eIP, err := eips.Apply(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error allocating EIP: %s", err)
	}

	log.Printf("[DEBUG] Waiting for EIP %#v to become available.", eIP)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForEIPActive(networkingClient, eIP.ID, timeout)
	if err != nil {
		return fmt.Errorf(
			"Error waiting for EIP (%s) to become ready: %s",
			eIP.ID, err)
	}

	err = bindToPort(d, eIP.ID, networkingClient, timeout)
	if err != nil {
		return fmt.Errorf("Error binding eip:%s to port: %s", eIP.ID, err)
	}

	d.SetId(eIP.ID)

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	eIP, err := eips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "eIP")
	}
	bandWidth, err := bandwidths.Get(networkingClient, eIP.BandwidthID).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching bandwidth: %s", err)
	}

	// Set public ip
	publicIP := []map[string]string{
		{
			"type":       eIP.Type,
			"ip_address": eIP.PublicAddress,
			"port_id":    eIP.PortID,
		},
	}
	d.Set("publicip", publicIP)

	// Set bandwidth
	bW := []map[string]interface{}{
		{
			"name":        bandWidth.Name,
			"size":        eIP.BandwidthSize,
			"id":          eIP.BandwidthID,
			"share_type":  eIP.BandwidthShareType,
			"charge_mode": bandWidth.ChargeMode,
		},
	}
	d.Set("bandwidth", bW)
	d.Set("address", eIP.PublicAddress)
	d.Set("region", GetRegion(d, config))
	d.Set("enterprise_project_id", eIP.EnterpriseProjectID)

	return nil
}

func resourceVpcEIPV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	// Update bandwidth change
	if d.HasChange("bandwidth") {
		var updateOpts bandwidths.UpdateOpts

		newBWList := d.Get("bandwidth").([]interface{})
		newMap := newBWList[0].(map[string]interface{})
		updateOpts.Size = newMap["size"].(int)
		updateOpts.Name = newMap["name"].(string)

		log.Printf("[DEBUG] Bandwidth Update Options: %#v", updateOpts)

		eIP, err := eips.Get(networkingClient, d.Id()).Extract()
		if err != nil {
			return CheckDeleted(d, err, "eIP")
		}
		_, err = bandwidths.Update(networkingClient, eIP.BandwidthID, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating bandwidth: %s", err)
		}

	}

	// Update publicip change
	if d.HasChange("publicip") {
		var updateOpts eips.UpdateOpts

		newIPList := d.Get("publicip").([]interface{})
		newMap := newIPList[0].(map[string]interface{})
		updateOpts.PortID = newMap["port_id"].(string)

		log.Printf("[DEBUG] PublicIP Update Options: %#v", updateOpts)
		_, err = eips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating publicip: %s", err)
		}
	}

	return resourceVpcEIPV1Read(d, meta)
}

func resourceVpcEIPV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating VPC client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = unbindToPort(d, d.Id(), networkingClient, timeout)
	if err != nil {
		log.Printf("[WARN] Error trying to unbind eip %s :%s", d.Id(), err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEIPDelete(networkingClient, d.Id()),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting EIP: %s", err)
	}

	d.SetId("")

	return nil
}

func getEIPStatus(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] EIP: %+v", e)
		if e.Status == "DOWN" || e.Status == "ACTIVE" {
			return e, "ACTIVE", nil
		}

		return e, "", nil
	}
}

func waitForEIPDelete(networkingClient *golangsdk.ServiceClient, eId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete EIP %s.\n", eId)

		e, err := eips.Get(networkingClient, eId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted EIP %s", eId)
				return e, "DELETED", nil
			}
			return e, "ACTIVE", err
		}

		err = eips.Delete(networkingClient, eId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted EIP %s", eId)
				return e, "DELETED", nil
			}
			return e, "ACTIVE", err
		}

		log.Printf("[DEBUG] EIP %s still active.\n", eId)
		return e, "ACTIVE", nil
	}
}

func resourcePublicIP(d *schema.ResourceData) eips.PublicIpOpts {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})

	publicip := eips.PublicIpOpts{
		Type:    rawMap["type"].(string),
		Address: rawMap["ip_address"].(string),
	}
	return publicip
}

func resourceBandWidth(d *schema.ResourceData) eips.BandwidthOpts {
	bandwidthRaw := d.Get("bandwidth").([]interface{})
	rawMap := bandwidthRaw[0].(map[string]interface{})

	bandwidth := eips.BandwidthOpts{
		Id:         rawMap["id"].(string),
		Name:       rawMap["name"].(string),
		Size:       rawMap["size"].(int),
		ShareType:  rawMap["share_type"].(string),
		ChargeMode: rawMap["charge_mode"].(string),
	}
	return bandwidth
}

func bindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	port_id, ok := rawMap["port_id"]
	if !ok || port_id == "" {
		return nil
	}

	pd := port_id.(string)
	log.Printf("[DEBUG] Bind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: pd}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func unbindToPort(d *schema.ResourceData, eipID string, networkingClient *golangsdk.ServiceClient, timeout time.Duration) error {
	publicIPRaw := d.Get("publicip").([]interface{})
	rawMap := publicIPRaw[0].(map[string]interface{})
	port_id, ok := rawMap["port_id"]
	if !ok || port_id == "" {
		return nil
	}

	pd := port_id.(string)
	log.Printf("[DEBUG] Unbind eip:%s to port: %s", eipID, pd)

	updateOpts := eips.UpdateOpts{PortID: ""}
	_, err := eips.Update(networkingClient, eipID, updateOpts).Extract()
	if err != nil {
		return err
	}
	return waitForEIPActive(networkingClient, eipID, timeout)
}

func waitForEIPActive(networkingClient *golangsdk.ServiceClient, eipID string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    getEIPStatus(networkingClient, eipID),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	return err
}

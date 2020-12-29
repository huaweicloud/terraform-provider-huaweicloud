package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/iec/v1/subnets"
)

func resourceIecSubnetDNSListV1(d *schema.ResourceData) []string {
	rawDNSN := d.Get("dns_list").([]interface{})
	dnsn := make([]string, len(rawDNSN))
	for i, raw := range rawDNSN {
		dnsn[i] = raw.(string)
	}
	return dnsn
}

func resourceIecSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceIecSubnetV1Create,
		Read:   resourceIecSubnetV1Read,
		Update: resourceIecSubnetV1Update,
		Delete: resourceIecSubnetV1Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"site_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"gateway_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIP,
			},
			"dhcp_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"dns_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"site_info": {
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

func resourceIecSubnetV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.IECV1Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	dhcp := d.Get("dhcp_enable").(bool)
	createOpts := subnets.CreateOpts{
		Name:       d.Get("name").(string),
		Cidr:       d.Get("cidr").(string),
		VpcID:      d.Get("vpc_id").(string),
		SiteID:     d.Get("site_id").(string),
		GatewayIP:  d.Get("gateway_ip").(string),
		DhcpEnable: &dhcp,
		DNSList:    resourceIecSubnetDNSListV1(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	n, err := subnets.Create(subnetClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC subnets: %s", err)
	}

	d.SetId(n.ID)
	log.Printf("[DEBUG] Waiting for IEC subnets (%s) to become active", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"UNKNOWN"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForIecSubnetStatus(subnetClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for IEC subnets (%s) to become ACTIVE: %s",
			n.ID, stateErr)
	}

	return resourceIecSubnetV1Read(d, config)
}

func resourceIecSubnetV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	n, err := subnets.Get(subnetClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving Huaweicloud IEC subnets")
	}

	log.Printf("[DEBUG] IEC subnets %s: %+v", d.Id(), n)

	d.Set("name", n.Name)
	d.Set("cidr", n.Cidr)
	d.Set("vpc_id", n.VpcID)
	d.Set("site_id", n.SiteID)
	d.Set("gateway_ip", n.GatewayIP)
	d.Set("dhcp_enable", n.DhcpEnable)
	d.Set("dns_list", n.DNSList)
	d.Set("site_info", n.SiteInfo)
	d.Set("status", n.Status)

	return nil
}

func resourceIecSubnetV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	var updateOpts subnets.UpdateOpts

	// name is mandatory while updating subnets
	updateOpts.Name = d.Get("name").(string)

	if d.HasChange("dhcp_enable") {
		dhcp := d.Get("dhcp_enable").(bool)
		updateOpts.DhcpEnable = &dhcp
	}
	if d.HasChange("dns_list") {
		updateOpts.DNSList = resourceSubnetDNSListV1(d)
	}

	_, err = subnets.Update(subnetClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud IEC subnets: %s", err)
	}

	return resourceIecSubnetV1Read(d, meta)
}

func resourceIecSubnetV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	subnetClient, err := config.IECV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud IEC client: %s", err)
	}

	err = subnets.Delete(subnetClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting Huaweicloud IEC subnets")
	}

	// waiting for subnets to become deleted
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE", "UNKNOWN"},
		Target:     []string{"DELETED"},
		Refresh:    waitForIecSubnetStatus(subnetClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, stateErr := stateConf.WaitForState()
	if stateErr != nil {
		return fmt.Errorf(
			"Error waiting for IEC subnets (%s) to become deleted: %s",
			d.Id(), stateErr)
	}

	d.SetId("")
	return nil
}

func waitForIecSubnetStatus(subnetClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := subnets.Get(subnetClient, id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud IEC subnets %s", id)
				return n, "DELETED", nil
			}
			return n, "ERROR", err
		}

		return n, n.Status, nil
	}
}

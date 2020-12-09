package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/peerings"
)

func ResourceVpcPeeringConnectionV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVPCPeeringV2Create,
		Read:   resourceVPCPeeringV2Read,
		Update: resourceVPCPeeringV2Update,
		Delete: resourceVPCPeeringV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{ //request and response parameters
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateString64WithChinese,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"peer_tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVPCPeeringV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.NetworkingV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc Peering Connection Client: %s", err)
	}

	requestvpcinfo := peerings.VpcInfo{
		VpcId: d.Get("vpc_id").(string),
	}

	acceptvpcinfo := peerings.VpcInfo{
		VpcId:    d.Get("peer_vpc_id").(string),
		TenantId: d.Get("peer_tenant_id").(string),
	}

	createOpts := peerings.CreateOpts{
		Name:           d.Get("name").(string),
		RequestVpcInfo: requestvpcinfo,
		AcceptVpcInfo:  acceptvpcinfo,
	}

	n, err := peerings.Create(peeringClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Vpc Peering Connection: %s", err)
	}

	log.Printf("[INFO] Vpc Peering Connection ID: %s", n.ID)

	log.Printf("[INFO] Waiting for Huaweicloud Vpc Peering Connection(%s) to become available", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"CREATING"},
		Target:     []string{"PENDING_ACCEPTANCE", "ACTIVE"},
		Refresh:    waitForVpcPeeringActive(peeringClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	d.SetId(n.ID)

	return resourceVPCPeeringV2Read(d, meta)

}

func resourceVPCPeeringV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud   Vpc Peering Connection Client: %s", err)
	}

	n, err := peerings.Get(peeringClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Huaweicloud Vpc Peering Connection: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("status", n.Status)
	d.Set("vpc_id", n.RequestVpcInfo.VpcId)
	d.Set("peer_vpc_id", n.AcceptVpcInfo.VpcId)
	d.Set("peer_tenant_id", n.AcceptVpcInfo.TenantId)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceVPCPeeringV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	peeringClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud  Vpc Peering Connection Client: %s", err)
	}

	var updateOpts peerings.UpdateOpts

	updateOpts.Name = d.Get("name").(string)

	_, err = peerings.Update(peeringClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud Vpc Peering Connection: %s", err)
	}

	return resourceVPCPeeringV2Read(d, meta)
}

func resourceVPCPeeringV2Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	peeringClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud  Vpc Peering Connection Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForVpcPeeringDelete(peeringClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud Vpc Peering Connection: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVpcPeeringActive(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := peerings.Get(peeringClient, peeringId).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "PENDING_ACCEPTANCE" || n.Status == "ACTIVE" {
			return n, n.Status, nil
		}

		return n, "CREATING", nil
	}
}

func waitForVpcPeeringDelete(peeringClient *golangsdk.ServiceClient, peeringId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := peerings.Get(peeringClient, peeringId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		err = peerings.Delete(peeringClient, peeringId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud vpc peering connection %s", peeringId)
				return r, "DELETED", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "ACTIVE", nil
				}
			}
			return r, "ACTIVE", err
		}

		return r, "ACTIVE", nil
	}
}

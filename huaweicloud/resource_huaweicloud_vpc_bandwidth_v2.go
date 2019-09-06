package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/huaweicloud/golangsdk"
	bandwidthsv1 "github.com/huaweicloud/golangsdk/openstack/networking/v1/bandwidths"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/bandwidths"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVpcBandWidthV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceVpcBandWidthV2Create,
		Read:   resourceVpcBandWidthV2Read,
		Update: resourceVpcBandWidthV2Update,
		Delete: resourceVpcBandWidthV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: validateIntegerInRange(5, 2000),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceVpcBandWidthV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	networkingV1Client, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	size := d.Get("size").(int)
	createOpts := bandwidths.CreateOpts{
		Name:                d.Get("name").(string),
		Size:                &size,
		EnterpriseProjectId: d.Get("enterprise_project_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	b, err := bandwidths.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Bandwidth: %s", err)
	}

	log.Printf("[DEBUG] Waiting for Bandwidth (%s) to become available.", b.ID)
	stateConf := &resource.StateChangeConf{
		Target:     []string{"NORMAL"},
		Pending:    []string{"CREATING"},
		Refresh:    waitForBandwidth(networkingV1Client, b.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf(
			"Error waiting for Bandwidth (%s) to become ACTIVE for creation: %s",
			b.ID, err)
	}
	d.SetId(b.ID)

	return resourceVpcBandWidthV2Read(d, meta)
}

func resourceVpcBandWidthV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	var bandwidthOpts bandwidths.Bandwidth

	if d.HasChange("name") {
		bandwidthOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("size") {
		bandwidthOpts.Size = d.Get("size").(int)
	}

	if bandwidthOpts != (bandwidths.Bandwidth{}) {
		updateOpts := bandwidths.UpdateOpts{
			Bandwidth: bandwidthOpts,
		}
		_, err := bandwidths.Update(networkingClient, d.Id(), updateOpts)
		if err != nil {
			return fmt.Errorf("Error updating Huaweicloud BandWidth (%s): %s", d.Id(), err)
		}
	}

	return resourceVpcBandWidthV2Read(d, meta)
}

func resourceVpcBandWidthV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	b, err := bandwidthsv1.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "bandwidth")
	}

	d.Set("name", b.Name)
	d.Set("size", b.Size)
	d.Set("enterprise_project_id", b.EnterpriseProjectID)

	return nil
}

func resourceVpcBandWidthV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	networkingV1Client, err := config.networkingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating networking client: %s", err)
	}

	err = bandwidths.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Bandwidth: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForBandwidth(networkingV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      3 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Bandwidth: %s", err)
	}

	d.SetId("")

	return nil
}

func waitForBandwidth(networkingClient *golangsdk.ServiceClient, Id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		b, err := bandwidthsv1.Get(networkingClient, Id).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return b, "DELETED", nil
			}
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud Bandwidth (%s) current status: %s", b.ID, b.Status)
		return b, b.Status, nil
	}
}

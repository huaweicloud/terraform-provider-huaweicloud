package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/compute/v2/extensions/floatingips"
)

func ResourceComputeFloatingIPV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeFloatingIPV2Create,
		Read:   resourceComputeFloatingIPV2Read,
		Update: nil,
		Delete: resourceComputeFloatingIPV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_vpc_eip resource instead",

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"pool": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "admin_external_net",
			},

			"address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeFloatingIPV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	createOpts := &floatingips.CreateOpts{
		Pool: d.Get("pool").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newFip, err := floatingips.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Floating IP: %s", err)
	}

	d.SetId(newFip.ID)

	return resourceComputeFloatingIPV2Read(d, meta)
}

func resourceComputeFloatingIPV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	fip, err := floatingips.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "floating ip")
	}

	log.Printf("[DEBUG] Retrieved Floating IP %s: %+v", d.Id(), fip)

	d.Set("pool", fip.Pool)
	d.Set("instance_id", fip.InstanceID)
	d.Set("address", fip.IP)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", GetRegion(d, config))

	return nil
}

func FloatingIPV2StateRefreshFunc(computeClient *golangsdk.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := floatingips.Get(computeClient, d.Id()).Extract()
		if err != nil {
			err = CheckDeleted(d, err, "Floating IP")
			if err != nil {
				return s, "", err
			} else {
				log.Printf("[DEBUG] Successfully deleted Floating IP %s", d.Id())
				return s, "DELETED", nil
			}
		}

		log.Printf("[DEBUG] Floating IP %s still active.\n", d.Id())
		return s, "ACTIVE", nil
	}
}

func resourceComputeFloatingIPV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud compute client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete Floating IP %s.\n", d.Id())

	if err := floatingips.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Error deleting Floating IP: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    FloatingIPV2StateRefreshFunc(computeClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Floating IP: %s", err)
	}

	return nil
}

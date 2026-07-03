package deprecated

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/compute/v2/extensions/floatingips"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
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
	cfg := meta.(*config.Config)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	createOpts := &floatingips.CreateOpts{
		Pool: d.Get("pool").(string),
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newFip, err := floatingips.Create(computeClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("error creating floating IP: %s", err)
	}

	d.SetId(newFip.ID)

	return resourceComputeFloatingIPV2Read(d, meta)
}

func resourceComputeFloatingIPV2Read(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	fip, err := floatingips.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "floating IP")
	}

	log.Printf("[DEBUG] Retrieved Floating IP %s: %+v", d.Id(), fip)

	d.Set("pool", fip.Pool)
	d.Set("instance_id", fip.InstanceID)
	d.Set("address", fip.IP)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", cfg.GetRegion(d))

	return nil
}

func FloatingIPV2StateRefreshFunc(computeClient *golangsdk.ServiceClient, d *schema.ResourceData) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		s, err := floatingips.Get(computeClient, d.Id()).Extract()
		if err != nil {
			err = common.CheckDeleted(d, err, "floating IP")
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
	cfg := meta.(*config.Config)
	computeClient, err := cfg.ComputeV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating compute client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete Floating IP %s.\n", d.Id())

	if err := floatingips.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("error deleting floating IP: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    FloatingIPV2StateRefreshFunc(computeClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting floating IP: %s", err)
	}

	return nil
}

package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/hw_snatrules"
)

func resourceNatSnatRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNatSnatRuleV2Create,
		Read:   resourceNatSnatRuleV2Read,
		Delete: resourceNatSnatRuleV2Delete,

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
			"nat_gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"source_type": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1),
				Optional:     true,
				ForceNew:     true,
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"cidr"},
			},
			"cidr": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"floating_ip_id": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressSnatFiplistDiffs,
			},
			"floating_ip_address": {
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

func resourceNatSnatRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	sourceType := d.Get("source_type").(int)
	if sourceType == 1 {
		if _, ok := d.GetOk("network_id"); ok {
			return fmt.Errorf("source_type and network_id is incompatible in the Direct Connect scenario (source_type=1)")
		}
	}

	createOpts := &hw_snatrules.CreateOpts{
		NatGatewayID: d.Get("nat_gateway_id").(string),
		FloatingIPID: d.Get("floating_ip_id").(string),
		NetworkID:    d.Get("network_id").(string),
		Cidr:         d.Get("cidr").(string),
		SourceType:   sourceType,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	snatRule, err := hw_snatrules.Create(natV2Client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creatting Snat Rule: %s", err)
	}

	log.Printf("[DEBUG] Waiting for HuaweiCloud Snat Rule (%s) to become available.", snatRule.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    waitForSnatRuleActive(natV2Client, snatRule.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Snat Rule: %s", err)
	}

	d.SetId(snatRule.ID)

	return resourceNatSnatRuleV2Read(d, meta)
}

func resourceNatSnatRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	snatRule, err := hw_snatrules.Get(natV2Client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Snat Rule")
	}

	d.Set("nat_gateway_id", snatRule.NatGatewayID)
	d.Set("floating_ip_id", snatRule.FloatingIPID)
	d.Set("floating_ip_address", snatRule.FloatingIPAddress)
	d.Set("source_type", snatRule.SourceType)
	d.Set("network_id", snatRule.NetworkID)
	d.Set("cidr", snatRule.Cidr)
	d.Set("status", snatRule.Status)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNatSnatRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	natV2Client, err := config.natV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud nat client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSnatRuleDelete(natV2Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud Snat Rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSnatRuleActive(natV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := hw_snatrules.Get(natV2Client, nId).Extract()
		if err != nil {
			return nil, "", err
		}

		log.Printf("[DEBUG] HuaweiCloud Snat Rule: %+v", n)
		if n.Status == "ACTIVE" {
			return n, "ACTIVE", nil
		}

		return n, "", nil
	}
}

func waitForSnatRuleDelete(natV2Client *golangsdk.ServiceClient, nId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete HuaweiCloud Snat Rule %s.\n", nId)

		n, err := hw_snatrules.Get(natV2Client, nId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		err = hw_snatrules.Delete(natV2Client, nId).ExtractErr()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted HuaweiCloud Snat Rule %s", nId)
				return n, "DELETED", nil
			}
			return n, "ACTIVE", err
		}

		log.Printf("[DEBUG] HuaweiCloud Snat Rule %s still active.\n", nId)
		return n, "ACTIVE", nil
	}
}

package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/whitelists"
)

func ResourceWhitelistV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceWhitelistV2Create,
		Read:   resourceWhitelistV2Read,
		Update: resourceWhitelistV2Update,
		Delete: resourceWhitelistV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enable_whitelist": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"whitelist": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressLBWhitelistDiffs,
			},
		},
	}
}

func resourceWhitelistV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	enableWhitelist := d.Get("enable_whitelist").(bool)
	createOpts := whitelists.CreateOpts{
		TenantId:        d.Get("tenant_id").(string),
		ListenerId:      d.Get("listener_id").(string),
		EnableWhitelist: &enableWhitelist,
		Whitelist:       d.Get("whitelist").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	wl, err := whitelists.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud Whitelist: %s", err)
	}

	d.SetId(wl.ID)
	return resourceWhitelistV2Read(d, meta)
}

func resourceWhitelistV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	wl, err := whitelists.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "whitelist")
	}

	log.Printf("[DEBUG] Retrieved whitelist %s: %#v", d.Id(), wl)

	d.SetId(wl.ID)
	d.Set("tenant_id", wl.TenantId)
	d.Set("listener_id", wl.ListenerId)
	d.Set("enable_whitelist", wl.EnableWhitelist)
	d.Set("whitelist", wl.Whitelist)

	return nil
}

func resourceWhitelistV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts whitelists.UpdateOpts
	if d.HasChange("enable_whitelist") {
		ew := d.Get("enable_whitelist").(bool)
		updateOpts.EnableWhitelist = &ew
	}
	if d.HasChange("whitelist") {
		updateOpts.Whitelist = d.Get("whitelist").(string)
	}

	log.Printf("[DEBUG] Updating whitelist %s with options: %#v", d.Id(), updateOpts)
	_, err = whitelists.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Unable to update whitelist %s: %s", d.Id(), err)
	}

	return resourceWhitelistV2Read(d, meta)
}

func resourceWhitelistV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete whitelist %s", d.Id())
	err = whitelists.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud whitelist: %s", err)
	}
	d.SetId("")
	return nil
}

package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/l7policies"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/pools"
)

func ResourceL7PolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceL7PolicyV2Create,
		Read:   resourceL7PolicyV2Read,
		Update: resourceL7PolicyV2Update,
		Delete: resourceL7PolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceL7PolicyV2Import,
		},

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

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_POOL", "REDIRECT_TO_LISTENER",
				}, true),
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"redirect_pool_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_listener_id"},
				Optional:      true,
			},

			"redirect_listener_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_pool_id"},
				Optional:      true,
			},

			"admin_state_up": {
				Type:         schema.TypeBool,
				Default:      true,
				Optional:     true,
				ValidateFunc: validateTrueOnly,
			},
		},
	}
}

func resourceL7PolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Assign some required variables for use in creation.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectListenerID := d.Get("redirect_listener_id").(string)

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectListenerID, redirectPoolID)
	if err != nil {
		return fmt.Errorf("Unable to create L7 Policy: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateOpts{
		TenantID:           d.Get("tenant_id").(string),
		Name:               d.Get("name").(string),
		Description:        d.Get("description").(string),
		Action:             l7policies.Action(action),
		ListenerID:         listenerID,
		RedirectPoolID:     redirectPoolID,
		RedirectListenerID: redirectListenerID,
		AdminStateUp:       &adminStateUp,
	}

	if v, ok := d.GetOk("position"); ok {
		createOpts.Position = int32(v.(int))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Make sure the associated pool is active before proceeding.
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return fmt.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(lbClient, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return err
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(lbClient, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create L7 Policy")
	var l7Policy *l7policies.L7Policy
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		l7Policy, err = l7policies.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating L7 Policy: %s", err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV2Read(d, meta)
}

func resourceL7PolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	d.Set("action", l7Policy.Action)
	d.Set("description", l7Policy.Description)
	d.Set("tenant_id", l7Policy.TenantID)
	d.Set("name", l7Policy.Name)
	d.Set("position", int(l7Policy.Position))
	d.Set("redirect_listener_id", l7Policy.RedirectListenerID)
	d.Set("redirect_pool_id", l7Policy.RedirectPoolID)
	d.Set("region", GetRegion(d, config))
	d.Set("admin_state_up", l7Policy.AdminStateUp)

	return nil
}

func resourceL7PolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Assign some required variables for use in updating.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectListenerID := d.Get("redirect_listener_id").(string)

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("redirect_pool_id") {
		redirectPoolID = d.Get("redirect_pool_id").(string)

		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectListenerID, redirectPoolID)
	if err != nil {
		return err
	}

	// Make sure the pool is active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return fmt.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(lbClient, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return err
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve parent listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve L7 Policy: %s: %s", d.Id(), err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(lbClient, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Unable to update L7 Policy %s: %s", d.Id(), err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	return resourceL7PolicyV2Read(d, meta)
}

func resourceL7PolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the listener.
	listener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmt.Errorf("Unable to retrieve parent listener (%s) for the L7 Policy: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Unable to retrieve L7 Policy")
	}

	// Wait for Listener to become active before continuing.
	err = waitForLBV2Listener(lbClient, listener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return CheckDeleted(d, err, "Error deleting L7 Policy")
	}

	err = waitForLBV2L7Policy(lbClient, listener, l7Policy, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return err
	}

	return nil
}

func resourceL7PolicyV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return nil, CheckDeleted(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s during the import: %#v", d.Id(), l7Policy)

	if l7Policy.ListenerID != "" {
		d.Set("listener_id", l7Policy.ListenerID)
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err := getListenerIDForL7Policy(lbClient, d.Id())
		if err != nil {
			return nil, err
		}
		d.Set("listener_id", listenerID)
	}

	return []*schema.ResourceData{d}, nil
}

func checkL7PolicyAction(action, redirectListenerID, redirectPoolID string) error {
	if action == "REDIRECT_TO_POOL" && redirectListenerID != "" {
		return fmt.Errorf("redirect_listener_id must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_LISTENER" && redirectPoolID != "" {
		return fmt.Errorf("redirect_pool_id must be empty when action is set to %s", action)
	}

	return nil
}

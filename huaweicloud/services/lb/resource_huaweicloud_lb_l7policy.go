package lb

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v2/l7policies"
	"github.com/chnsz/golangsdk/openstack/elb/v2/listeners"
	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API ELB GET /v2/{project_id}/elb/pools/{pool_id}
// @API ELB GET /v2/{project_id}/elb/listeners/{listener_id}
// @API ELB POST /v2/{project_id}/elb/l7policies
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB GET /v2/{project_id}/elb/loadbalancers
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}/statuses
// @API ELB GET /v2/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB PUT /v2/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB DELETE /v2/{project_id}/elb/l7policies/{l7policy_id}
func ResourceL7PolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV2Create,
		ReadContext:   resourceL7PolicyV2Read,
		UpdateContext: resourceL7PolicyV2Update,
		DeleteContext: resourceL7PolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7PolicyV2Import,
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
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
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
				Type:         schema.TypeString,
				ExactlyOneOf: []string{"redirect_listener_id"},
				Optional:     true,
			},

			"redirect_listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7PolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Assign some required variables for use in creation.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectListenerID := d.Get("redirect_listener_id").(string)

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectListenerID, redirectPoolID)
	if err != nil {
		return fmtp.DiagErrorf("Unable to create L7 Policy: %s", err)
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

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Make sure the associated pool is active before proceeding.
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, lbClient, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Attempting to create L7 Policy")
	l7Policy, err := l7policies.Create(lbClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating L7 Policy: %s", err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV2Read(ctx, d, meta)
}

func resourceL7PolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving L7 Policy")
	}

	logp.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("action", l7Policy.Action),
		d.Set("description", l7Policy.Description),
		d.Set("tenant_id", l7Policy.TenantID),
		d.Set("name", l7Policy.Name),
		d.Set("position", int(l7Policy.Position)),
		d.Set("redirect_listener_id", l7Policy.RedirectListenerID),
		d.Set("redirect_pool_id", l7Policy.RedirectPoolID),
		d.Set("admin_state_up", l7Policy.AdminStateUp),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting L7 Policy fields: %s", err)
	}

	return nil
}

func resourceL7PolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
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
		return diag.FromErr(err)
	}

	// Make sure the pool is active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return fmtp.DiagErrorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBV2Pool(ctx, lbClient, pool.ID, "ACTIVE", lbPendingStatuses, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve parent listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve L7 Policy: %s: %s", d.Id(), err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, parentListener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Unable to update L7 Policy %s: %s", d.Id(), err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV2Read(ctx, d, meta)
}

func resourceL7PolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the listener.
	listener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve parent listener (%s) for the L7 Policy: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Unable to retrieve L7 Policy")
	}

	// Wait for Listener to become active before continuing.
	err = waitForLBV2Listener(ctx, lbClient, listener.ID, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting L7 Policy")
	}

	err = waitForLBV2L7Policy(ctx, lbClient, listener, l7Policy, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7PolicyV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return nil, common.CheckDeleted(d, err, "L7 Policy")
	}

	logp.Printf("[DEBUG] Retrieved L7 Policy %s during the import: %#v", d.Id(), l7Policy)

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
		return fmtp.Errorf("redirect_listener_id must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_LISTENER" && redirectPoolID != "" {
		return fmtp.Errorf("redirect_pool_id must be empty when action is set to %s", action)
	}

	return nil
}

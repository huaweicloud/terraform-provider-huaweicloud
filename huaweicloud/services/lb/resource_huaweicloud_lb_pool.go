package lb

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourcePoolV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolV2Create,
		ReadContext:   resourcePoolV2Read,
		UpdateContext: resourcePoolV2Update,
		DeleteContext: resourcePoolV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP",
				}, false),
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				AtLeastOneOf: []string{
					"listener_id",
				},
			},

			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP",
				}, false),
			},

			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
						},

						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
					},
				},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourcePoolV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	var persistence pools.SessionPersistence
	if p, ok := d.GetOk("persistence"); ok {
		pV := (p.([]interface{}))[0].(map[string]interface{})

		persistence = pools.SessionPersistence{
			Type: pV["type"].(string),
		}

		if persistence.Type == "APP_COOKIE" {
			if pV["cookie_name"].(string) == "" {
				return fmtp.DiagErrorf(
					"Persistence cookie_name needs to be set if using 'APP_COOKIE' persistence type")
			}
			persistence.CookieName = pV["cookie_name"].(string)

			if pV["timeout"].(int) != 0 {
				return fmtp.DiagErrorf(
					"Persistence timeout is invalid when type is set to 'APP_COOKIE'")
			}
		} else {
			if pV["cookie_name"].(string) != "" {
				return fmtp.DiagErrorf(
					"Persistence cookie_name can only be set if using 'APP_COOKIE' persistence type")
			}
		}

		persistence.PersistenceTimeout = pV["timeout"].(int)
	}

	createOpts := pools.CreateOpts{
		TenantID:       d.Get("tenant_id").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Protocol:       pools.Protocol(d.Get("protocol").(string)),
		LoadbalancerID: d.Get("loadbalancer_id").(string),
		ListenerID:     d.Get("listener_id").(string),
		LBMethod:       pools.LBMethod(d.Get("lb_method").(string)),
		AdminStateUp:   &adminStateUp,
	}

	// Must omit if not set
	if persistence != (pools.SessionPersistence{}) {
		createOpts.Persistence = &persistence
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	lbID := createOpts.LoadbalancerID
	listenerID := createOpts.ListenerID
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if listenerID != "" {
		// Wait for Listener to become active before continuing
		err = waitForLBV2Listener(ctx, lbClient, listenerID, "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	pool, err := pools.Create(lbClient, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating pool: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		// Pool exists by now so we can ask for lbID
		err = waitForLBV2viaPool(ctx, lbClient, pool.ID, "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(pool.ID)

	return resourcePoolV2Read(ctx, d, meta)
}

func resourcePoolV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb v2 client: %s", err)
	}

	pool, err := pools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving member")
	}

	logp.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("lb_method", pool.LBMethod),
		d.Set("protocol", pool.Protocol),
		d.Set("description", pool.Description),
		d.Set("tenant_id", pool.TenantID),
		d.Set("admin_state_up", pool.AdminStateUp),
		d.Set("name", pool.Name),
	)

	if len(pool.Loadbalancers) != 0 {
		mErr = multierror.Append(mErr, d.Set("loadbalancer_id", pool.Loadbalancers[0].ID))
	}

	if len(pool.Listeners) != 0 {
		mErr = multierror.Append(mErr, d.Set("listener_id", pool.Listeners[0].ID))
	}

	if pool.Persistence.Type != "" {
		var persistence []map[string]interface{} = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["cookie_name"] = pool.Persistence.CookieName
		params["type"] = pool.Persistence.Type
		params["timeout"] = pool.Persistence.PersistenceTimeout
		persistence[0] = params
		mErr = multierror.Append(mErr, d.Set("persistence", persistence))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting pool fields: %s", err)
	}
	return nil
}

func resourcePoolV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = pools.LBMethod(d.Get("lb_method").(string))
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(ctx, lbClient, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = pools.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Unable to update pool %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active before continuing
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		err = waitForLBV2viaPool(ctx, lbClient, d.Id(), "ACTIVE", timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolV2Read(ctx, d, meta)
}

func resourcePoolV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	timeout := d.Timeout(schema.TimeoutDelete)
	lbID := d.Get("loadbalancer_id").(string)
	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	logp.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = pools.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if lbID != "" {
		err = waitForLBV2LoadBalancer(ctx, lbClient, lbID, "ACTIVE", nil, timeout)
	} else {
		// Wait for Pool to delete
		err = waitForLBV2Pool(ctx, lbClient, d.Id(), "DELETED", nil, timeout)
	}
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

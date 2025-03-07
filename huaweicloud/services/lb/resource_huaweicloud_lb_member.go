package lb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v2/pools"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API ELB POST /v2/{project_id}/elb/pools/{pool_id}/members
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB GET /v2/{project_id}/elb/listeners/{listener_id}
// @API ELB GET /v2/{project_id}/elb/pools/{pool_id}
// @API ELB GET /v2/{project_id}/elb/pools/{pool_id}/members/{memeber_id}
// @API ELB PUT /v2/{project_id}/elb/pools/{pool_id}/members/{memeber_id}
// @API ELB DELETE /v2/{project_id}/elb/pools/{pool_id}/members/{memeber_id}
func ResourceMemberV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMemberV2Create,
		ReadContext:   resourceMemberV2Read,
		UpdateContext: resourceMemberV2Update,
		DeleteContext: resourceMemberV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMemberV2Import,
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

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				ForceNew:   true,
				Deprecated: "tenant_id is deprecated",
			},

			"address": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"weight": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "the IPv4 subnet ID of the subnet in which to access the member",
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"backend_server_status": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"operating_status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// deprecated
			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "schema: Deprecated",
			},
		},
	}
}

func resourceMemberV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := pools.CreateMemberOpts{
		Name:         d.Get("name").(string),
		TenantID:     d.Get("tenant_id").(string),
		Address:      d.Get("address").(string),
		ProtocolPort: d.Get("protocol_port").(int),
		Weight:       d.Get("weight").(int),
		AdminStateUp: &adminStateUp,
	}

	// Must omit if not set
	if v, ok := d.GetOk("subnet_id"); ok {
		createOpts.SubnetID = v.(string)
	}

	// Wait for LB to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	member, err := pools.CreateMember(lbClient, poolID, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating member: %s", err)
	}

	// Wait for LB to become ACTIVE again
	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(member.ID)

	return resourceMemberV2Read(ctx, d, meta)
}

func resourceMemberV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	member, err := pools.GetMember(lbClient, d.Get("pool_id").(string), d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving member")
	}

	logp.Printf("[DEBUG] Retrieved member %s: %#v", d.Id(), member)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", member.Name),
		d.Set("weight", member.Weight),
		d.Set("tenant_id", member.TenantID),
		d.Set("subnet_id", member.SubnetID),
		d.Set("address", member.Address),
		d.Set("protocol_port", member.ProtocolPort),
		d.Set("operating_status", member.OperatingStatus),
		d.Set("backend_server_status", member.AdminStateUp),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting member fields: %s", err)
	}

	return nil
}

func resourceMemberV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts pools.UpdateMemberOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("weight") {
		updateOpts.Weight = d.Get("weight").(int)
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Wait for LB to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Updating member %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = pools.UpdateMember(lbClient, poolID, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Unable to update member %s: %s", d.Id(), err)
	}

	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMemberV2Read(ctx, d, meta)
}

func resourceMemberV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Wait for Pool to become active before continuing
	poolID := d.Get("pool_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Attempting to delete member %s", d.Id())
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = pools.DeleteMember(lbClient, poolID, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting member: %s", err)
	}

	// Wait for LB to become ACTIVE
	err = waitForLBV2viaPool(ctx, lbClient, poolID, "ACTIVE", timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMemberV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("invalid format specified for member. Format must be <pool_id>/<member_id>")
		return nil, err
	}

	poolID := parts[0]
	memberID := parts[1]

	d.SetId(memberID)
	d.Set("pool_id", poolID)

	return []*schema.ResourceData{d}, nil
}

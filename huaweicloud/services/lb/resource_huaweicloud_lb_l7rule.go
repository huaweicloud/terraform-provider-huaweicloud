package lb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/elb/v2/l7policies"
	"github.com/chnsz/golangsdk/openstack/elb/v2/listeners"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// @API ELB GET /v2/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB GET /v2/{project_id}/elb/listeners/{listener_id}
// @API ELB POST /v2/{project_id}/elb/l7policies/{l7policy_id}/rules
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}
// @API ELB GET /v2/{project_id}/elb/loadbalancers
// @API ELB GET /v2/{project_id}/elb/loadbalancers/{loadbalancer_id}/statuses
// @API ELB GET /v2/{project_id}/elb/l7policies/{l7policy_id}/rules/{rule_id}
// @API ELB PUT /v2/{project_id}/elb/l7policies/{l7policy_id}/rules/{rule_id}
// @API ELB DELETE /v2/{project_id}/elb/l7policies/{l7policy_id}/rules/{rule_id}
func ResourceL7RuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7RuleV2Create,
		ReadContext:   resourceL7RuleV2Read,
		UpdateContext: resourceL7RuleV2Update,
		DeleteContext: resourceL7RuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7RuleV2Import,
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

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"compare_type": {
				Type:     schema.TypeString,
				Required: true,
			},

			"l7policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
			},

			"key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7RuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Assign some required variables for use in creation.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := ""
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)
	compareType := d.Get("compare_type").(string)
	adminStateUp := d.Get("admin_state_up").(bool)

	// Ensure the right combination of options have been specified.
	err = checkL7RuleType(ruleType, key)
	if err != nil {
		return fmtp.DiagErrorf("Unable to create L7 Rule: %s", err)
	}

	createOpts := l7policies.CreateRuleOpts{
		TenantID:     d.Get("tenant_id").(string),
		RuleType:     l7policies.RuleType(ruleType),
		CompareType:  l7policies.CompareType(compareType),
		Value:        d.Get("value").(string),
		Key:          key,
		AdminStateUp: &adminStateUp,
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Attempting to create L7 Rule")
	l7Rule, err := l7policies.CreateRule(lbClient, l7policyID, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating L7 Rule: %s", err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId(l7Rule.ID)
	d.Set("listener_id", listenerID)

	return resourceL7RuleV2Read(ctx, d, meta)
}

func resourceL7RuleV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7policyID := d.Get("l7policy_id").(string)

	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving L7 Rule")
	}

	logp.Printf("[DEBUG] Retrieved L7 Rule %s: %#v", d.Id(), l7Rule)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("l7policy_id", l7policyID),
		d.Set("type", l7Rule.RuleType),
		d.Set("compare_type", l7Rule.CompareType),
		d.Set("tenant_id", l7Rule.TenantID),
		d.Set("value", l7Rule.Value),
		d.Set("key", l7Rule.Key),
		d.Set("admin_state_up", l7Rule.AdminStateUp),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting L7 Rule fields: %s", err)
	}

	return nil
}

func resourceL7RuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Assign some required variables for use in updating.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)

	// Key should always be set
	updateOpts := l7policies.UpdateRuleOpts{
		Key: &key,
	}

	if d.HasChange("compare_type") {
		updateOpts.CompareType = l7policies.CompareType(d.Get("compare_type").(string))
	}
	if d.HasChange("value") {
		updateOpts.Value = d.Get("value").(string)
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	err = checkL7RuleType(ruleType, key)
	if err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to get parent L7 Policy: %s", err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to get L7 Rule: %s", err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Updating L7 Rule %s with options: %#v", d.Id(), updateOpts)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err := l7policies.UpdateRule(lbClient, l7policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmtp.DiagErrorf("Unable to update L7 Rule %s: %s", d.Id(), err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBV2L7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7RuleV2Read(ctx, d, meta)
}

func resourceL7RuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud elb client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve parent listener (%s) for the L7 Rule: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve parent L7 Policy (%s) for the L7 Rule: %s", l7policyID, err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Unable to retrieve L7 Rule")
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBV2L7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", lbPendingStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	logp.Printf("[DEBUG] Attempting to delete L7 Rule %s", d.Id())
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = l7policies.DeleteRule(lbClient, l7policyID, d.Id()).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "Error deleting L7 Rule")
	}

	err = waitForLBV2L7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "DELETED", lbPendingDeleteStatuses, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7RuleV2Import(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for L7 Rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	config := meta.(*config.Config)
	lbClient, err := config.LoadBalancerClient(config.GetRegion(d))
	if err != nil {
		return nil, fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	listenerID := ""
	l7policyID := parts[0]
	l7ruleID := parts[1]

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return nil, fmtp.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaSv2 extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return nil, err
		}
	}

	d.SetId(l7ruleID)
	d.Set("l7policy_id", l7policyID)
	d.Set("listener_id", listenerID)

	return []*schema.ResourceData{d}, nil
}

func checkL7RuleType(ruleType, key string) error {
	keyRequired := []string{"COOKIE", "HEADER"}
	if utils.StrSliceContains(keyRequired, ruleType) && key == "" {
		return fmtp.Errorf("key attribute is required, when the L7 Rule type is %s", strings.Join(keyRequired, " or "))
	} else if !utils.StrSliceContains(keyRequired, ruleType) && key != "" {
		return fmtp.Errorf("key attribute must not be used, when the L7 Rule type is not %s", strings.Join(keyRequired, " or "))
	}
	return nil
}

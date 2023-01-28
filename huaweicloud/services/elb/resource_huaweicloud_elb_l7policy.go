package elb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceL7PolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV3Create,
		ReadContext:   resourceL7PolicyV3Read,
		UpdateContext: resourceL7PolicyV3Update,
		DeleteContext: resourceL7PolicyV3Delete,
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

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"redirect_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceL7PolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	createOpts := l7policies.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Action:         "REDIRECT_TO_POOL",
		ListenerID:     d.Get("listener_id").(string),
		RedirectPoolID: d.Get("redirect_pool_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	l7Policy, err := l7policies.Create(lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating L7 Policy: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	// Wait for L7 Policy to become active before continuing
	err = waitForElbV3Policy(ctx, lbClient, l7Policy.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func resourceL7PolicyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	mErr := multierror.Append(nil,
		d.Set("description", l7Policy.Description),
		d.Set("name", l7Policy.Name),
		d.Set("listener_id", l7Policy.ListenerID),
		d.Set("redirect_pool_id", l7Policy.RedirectPoolID),
		d.Set("region", config.GetRegion(d)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB l7policy fields: %s", err)
	}

	return nil
}

func resourceL7PolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

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
		redirectPoolID := d.Get("redirect_pool_id").(string)
		updateOpts.RedirectPoolID = &redirectPoolID
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update L7 Policy %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Policy(ctx, lbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func resourceL7PolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating elb client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting L7 Policy")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Policy(ctx, lbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForElbV3Policy(ctx context.Context, elbClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {

	log.Printf("[DEBUG] Waiting for policy %s to become %s", id, target)

	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3PolicyRefreshFunc(elbClient, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: policy %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for policy %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3PolicyRefreshFunc(elbClient *golangsdk.ServiceClient,
	id string) resource.StateRefreshFunc {

	return func() (interface{}, string, error) {
		policy, err := l7policies.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return policy, policy.ProvisioningStatus, nil
	}
}

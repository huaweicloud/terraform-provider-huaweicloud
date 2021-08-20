package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/elb/v3/l7policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceL7PolicyV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceL7PolicyV3Create,
		Read:   resourceL7PolicyV3Read,
		Update: resourceL7PolicyV3Update,
		Delete: resourceL7PolicyV3Delete,

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

func resourceL7PolicyV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	createOpts := l7policies.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Action:         "REDIRECT_TO_POOL",
		ListenerID:     d.Get("listener_id").(string),
		RedirectPoolID: d.Get("redirect_pool_id").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	l7Policy, err := l7policies.Create(lbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating L7 Policy: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	// Wait for L7 Policy to become active before continuing
	err = waitForElbV3Policy(lbClient, l7Policy.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV3Read(d, meta)
}

func resourceL7PolicyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "L7 Policy")
	}

	logp.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	d.Set("description", l7Policy.Description)
	d.Set("name", l7Policy.Name)
	d.Set("listener_id", l7Policy.ListenerID)
	d.Set("redirect_pool_id", l7Policy.RedirectPoolID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceL7PolicyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
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

	logp.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to update L7 Policy %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Policy(lbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	return resourceL7PolicyV3Read(d, meta)
}

func resourceL7PolicyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	lbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting L7 Policy")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Policy(lbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}

func waitForElbV3Policy(elbClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {

	logp.Printf("[DEBUG] Waiting for policy %s to become %s", id, target)

	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3PolicyRefreshFunc(elbClient, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmtp.Errorf("Error: policy %s not found: %s", id, err)
			}
		}
		return fmtp.Errorf("Error waiting for policy %s to become %s: %s", id, target, err)
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

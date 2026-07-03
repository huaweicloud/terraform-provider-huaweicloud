package deprecated

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/networking/v2/extensions/fwaas_v2/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceFWPolicyV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceFWPolicyV2Create,
		Read:   resourceFWPolicyV2Read,
		Update: resourceFWPolicyV2Update,
		Delete: resourceFWPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		DeprecationMessage: "use huaweicloud_network_acl resource instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
			"audited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tenant_id": {
				Type:       schema.TypeString,
				Optional:   true,
				ForceNew:   true,
				Computed:   true,
				Deprecated: "tenant_id is deprecated",
			},
			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceFWPolicyV2Create(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	fwClient, err := cfg.FwV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating fw client: %s", err)
	}

	v := d.Get("rules").([]interface{})

	log.Printf("[DEBUG] Rules found : %#v", v)
	log.Printf("[DEBUG] Rules count : %d", len(v))

	rules := make([]string, len(v))
	for i, v := range v {
		rules[i] = v.(string)
	}

	audited := d.Get("audited").(bool)

	opts := PolicyCreateOpts{
		policies.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Audited:     &audited,
			TenantID:    d.Get("tenant_id").(string),
			Rules:       rules,
		},
		MapValueSpecs(d),
	}

	if r, ok := d.GetOk("shared"); ok {
		shared := r.(bool)
		opts.Shared = &shared
	}

	log.Printf("[DEBUG] Create firewall policy: %#v", opts)

	policy, err := policies.Create(fwClient, opts).Extract()
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Firewall policy created: %#v", policy)

	d.SetId(policy.ID)

	return resourceFWPolicyV2Read(d, meta)
}

func resourceFWPolicyV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Retrieve information about firewall policy: %s", d.Id())

	cfg := meta.(*config.Config)
	fwClient, err := cfg.FwV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating fw client: %s", err)
	}

	policy, err := policies.Get(fwClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "FW policy")
	}

	log.Printf("[DEBUG] Read firewall policy %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("shared", policy.Shared)
	d.Set("audited", policy.Audited)
	d.Set("tenant_id", policy.TenantID)
	if err := d.Set("rules", policy.Rules); err != nil {
		return fmt.Errorf("error saving rules to state for firewall policy (%s): %s", d.Id(), err)
	}
	d.Set("region", cfg.GetRegion(d))

	return nil
}

func resourceFWPolicyV2Update(d *schema.ResourceData, meta interface{}) error {
	cfg := meta.(*config.Config)
	fwClient, err := cfg.FwV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating fw client: %s", err)
	}

	opts := policies.UpdateOpts{}

	if d.HasChange("name") {
		opts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		opts.Description = d.Get("description").(string)
	}

	if d.HasChange("rules") {
		v := d.Get("rules").([]interface{})

		log.Printf("[DEBUG] Rules found : %#v", v)
		log.Printf("[DEBUG] Rules count : %d", len(v))

		rules := make([]string, len(v))
		for i, v := range v {
			rules[i] = v.(string)
		}
		opts.Rules = rules
	}

	log.Printf("[DEBUG] Updating firewall policy with id %s: %#v", d.Id(), opts)

	err = policies.Update(fwClient, d.Id(), opts).Err
	if err != nil {
		return err
	}

	return resourceFWPolicyV2Read(d, meta)
}

func resourceFWPolicyV2Delete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Destroy firewall policy: %s", d.Id())

	cfg := meta.(*config.Config)
	fwClient, err := cfg.FwV2Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating fw client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    waitForFirewallPolicyDeletion(fwClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(context.Background()); err != nil {
		return err
	}

	return nil
}

func waitForFirewallPolicyDeletion(fwClient *golangsdk.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := policies.Delete(fwClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(golangsdk.ErrDefault409); ok {
			// This error usually means that the policy is attached
			// to a firewall. At this point, the firewall is probably
			// being delete. So, we retry a few times.
			return nil, "ACTIVE", nil
		}

		return nil, "ACTIVE", err
	}
}

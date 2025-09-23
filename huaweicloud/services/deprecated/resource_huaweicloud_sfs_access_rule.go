package deprecated

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API SFS POST /v2/{project_id}/shares/{id}/action
func ResourceSFSAccessRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSAccessRuleV2Create,
		ReadContext:   resourceSFSAccessRuleV2Read,
		DeleteContext: resourceSFSAccessRuleV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceSFSAccessRuleV2Import,
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
			"sfs_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "rw",
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "cert",
			},
			"access_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceSFSAccessRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	shareID := d.Get("sfs_id").(string)
	grantAccessOpts := shares.GrantAccessOpts{
		AccessLevel: d.Get("access_level").(string),
		AccessType:  d.Get("access_type").(string),
		AccessTo:    d.Get("access_to").(string),
	}

	log.Printf("[DEBUG] Applied access rule to share file %s, opts: %#v", shareID, grantAccessOpts)
	grant, err := shares.GrantAccess(sfsClient, shareID, grantAccessOpts).ExtractAccess()
	if err != nil {
		return diag.Errorf("error creating access rule to share file: %s", err)
	}

	d.SetId(grant.ID)
	// wait access rule to become active
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"new", "queued_to_apply", "applying"},
		Target:     []string{"active"},
		Refresh:    waitForSFSAccessStatus(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error creating SFS access rule: %s", err)
	}

	return resourceSFSAccessRuleV2Read(ctx, d, meta)
}

func resourceSFSAccessRuleV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	shareID := d.Get("sfs_id").(string)
	rules, err := shares.ListAccessRights(sfsClient, shareID).ExtractAccessRights()
	if err != nil {
		return diag.Errorf("error retrieving SFS access rules: %s", err)
	}

	for _, rule := range rules {
		if rule.ID == d.Id() {
			mErr := multierror.Append(nil,
				d.Set("access_to", rule.AccessTo),
				d.Set("access_type", rule.AccessType),
				d.Set("access_level", rule.AccessLevel),
				d.Set("status", rule.State),
			)

			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	// the access rule was not found
	log.Printf("[WARN] access rule (%s) of share file %s was not exist!", d.Id(), shareID)
	d.SetId("")
	return nil
}

func resourceSFSAccessRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	shareID := d.Get("sfs_id").(string)
	deleteAccessOpts := shares.DeleteAccessOpts{AccessID: d.Id()}
	deny := shares.DeleteAccess(sfsClient, shareID, deleteAccessOpts)
	if deny.Err != nil {
		return common.CheckDeletedDiag(d, deny.Err, "Error deleting access rule")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "queued_to_deny", "denying"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSAccessStatus(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error deleting SFS access rule: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSAccessStatus(sfsClient *golangsdk.ServiceClient, shareID, ruleID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		rules, err := shares.ListAccessRights(sfsClient, shareID).ExtractAccessRights()
		if err != nil {
			log.Printf("[WARN] list access rules error")
			return nil, "error", err
		}

		for _, rule := range rules {
			if rule.ID == ruleID {
				log.Printf("[DEBUG] find access rule %s, state: %s", ruleID, rule.State)
				return rule, rule.State, nil
			}
		}

		// the rule was not found, seem as deleted
		log.Printf("[DEBUG] could not find the access rule %s", ruleID)
		return shares.AccessRight{}, "deleted", nil
	}
}

func resourceSFSAccessRuleV2Import(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	arr := strings.Split(d.Id(), "/")
	shareID := arr[0]
	ruleID := arr[1]
	d.Set("sfs_id", shareID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}

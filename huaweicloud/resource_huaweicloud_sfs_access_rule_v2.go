package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func resourceSFSAccessRuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSFSAccessRuleV2Create,
		Read:   resourceSFSAccessRuleV2Read,
		Delete: resourceSFSAccessRuleV2Delete,

		Importer: &schema.ResourceImporter{
			State: resourceSFSAccessRuleV2Import,
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

func resourceSFSAccessRuleV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Client: %s", err)
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
		return fmt.Errorf("Error creating access rule to share file: %s", err)
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
	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS access rule: %s", err)
	}

	return resourceSFSAccessRuleV2Read(d, meta)
}

func resourceSFSAccessRuleV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Client: %s", err)
	}

	shareID := d.Get("sfs_id").(string)
	rules, err := shares.ListAccessRights(sfsClient, shareID).ExtractAccessRights()
	if err != nil {
		return fmt.Errorf("Error retrieving HuaweiCloud Shares rules: %s", err)
	}

	for _, rule := range rules {
		if rule.ID == d.Id() {
			d.Set("access_to", rule.AccessTo)
			d.Set("access_type", rule.AccessType)
			d.Set("access_level", rule.AccessLevel)
			d.Set("status", rule.State)
			return nil
		}
	}

	// the access rule was not found
	log.Printf("[WARN] access rule (%s) of share file %s was not exist!", d.Id(), shareID)
	d.SetId("")
	return nil
}

func resourceSFSAccessRuleV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud SFS Client: %s", err)
	}

	shareID := d.Get("sfs_id").(string)
	deleteAccessOpts := shares.DeleteAccessOpts{AccessID: d.Id()}
	deny := shares.DeleteAccess(sfsClient, shareID, deleteAccessOpts)
	if deny.Err != nil {
		return CheckDeleted(d, deny.Err, "Error deleting access rule")
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "queued_to_deny", "denying"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSAccessStatus(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud SFS access rule: %s", err)
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

func resourceSFSAccessRuleV2Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	arr := strings.Split(d.Id(), "/")
	shareID := arr[0]
	ruleID := arr[1]
	d.Set("sfs_id", shareID)
	d.SetId(ruleID)

	return []*schema.ResourceData{d}, nil
}

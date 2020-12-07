package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/sfs/v2/shares"
)

func resourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceSFSFileSystemV2Create,
		Read:   resourceSFSFileSystemV2Read,
		Update: resourceSFSFileSystemV2Update,
		Delete: resourceSFSFileSystemV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_proto": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "NFS",
			},
			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"access_level": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"access_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": tagsSchema(),
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"share_access_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rule_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"export_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_rule_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_to": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceSFSFileSystemV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))

	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud File Share Client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		IsPublic:         d.Get("is_public").(bool),
		Metadata:         resourceSFSMetadataV2(d, config),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	create, err := shares.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud File Share: %s", err)
	}

	d.SetId(create.ID)
	log.Printf("[INFO] Share ID: %s", create.Name)

	log.Printf("[DEBUG] Waiting for Huaweicloud SFS File Share (%s) to be become available", create.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    waitForSFSFileActive(sfsClient, create.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, StateErr := stateConf.WaitForState()
	if StateErr != nil {
		return fmt.Errorf("Error waiting for Share File (%s) to become available: %s ", d.Id(), StateErr)
	}

	// specified the "access_to" field, apply first access rule to share file
	if _, ok := d.GetOk("access_to"); ok {
		grantAccessOpts := shares.GrantAccessOpts{
			AccessTo: d.Get("access_to").(string),
		}

		if _, ok := d.GetOk("access_level"); ok {
			grantAccessOpts.AccessLevel = d.Get("access_level").(string)
		} else {
			grantAccessOpts.AccessLevel = "rw"
		}

		if _, ok := d.GetOk("access_type"); ok {
			grantAccessOpts.AccessType = d.Get("access_type").(string)
		} else {
			grantAccessOpts.AccessType = "cert"
		}

		grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()
		if accessErr != nil {
			return fmt.Errorf("Error applying access rule to share file : %s", accessErr)
		}

		log.Printf("[DEBUG] Applied access rule (%s) to share file %s", grant.ID, d.Id())
		d.Set("share_access_id", grant.ID)
	}

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(sfsClient, "sfs", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of sfs %s: %s", d.Id(), tagErr)
		}
	}

	return resourceSFSFileSystemV2Read(d, meta)
}

func resourceSFSFileSystemV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud File Share Client: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Huaweicloud Shares: %s", err)
	}

	d.Set("name", n.Name)
	d.Set("share_proto", n.ShareProto)
	d.Set("size", n.Size)
	d.Set("description", n.Description)
	d.Set("share_type", n.ShareType)
	d.Set("is_public", n.IsPublic)
	d.Set("availability_zone", n.AvailabilityZone)
	d.Set("region", GetRegion(d, config))
	d.Set("export_location", n.ExportLocation)
	d.Set("host", n.Host)
	d.Set("enterprise_project_id", n.Metadata["enterprise_project_id"])

	// NOTE: This tries to remove system metadata.
	md := make(map[string]string)
	var sys_keys = [2]string{"enterprise_project_id", "share_used"}

OUTER:
	for key, val := range n.Metadata {
		if strings.HasPrefix(key, "#sfs") {
			continue
		}
		for i := range sys_keys {
			if key == sys_keys[i] {
				continue OUTER
			}
		}
		md[key] = val
	}
	d.Set("metadata", md)

	// list access rules
	rules, err := shares.ListAccessRights(sfsClient, d.Id()).ExtractAccessRights()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Huaweicloud Shares rules: %s", err)
	}

	var ruleExist bool
	accessID := d.Get("share_access_id").(string)
	allAccessRules := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		acessRule := map[string]interface{}{
			"access_rule_id": rule.ID,
			"access_level":   rule.AccessLevel,
			"access_type":    rule.AccessType,
			"access_to":      rule.AccessTo,
			"status":         rule.State,
		}
		allAccessRules = append(allAccessRules, acessRule)

		// find share_access_id
		if accessID != "" && rule.ID == accessID {
			d.Set("access_rule_status", rule.State)
			d.Set("access_to", rule.AccessTo)
			d.Set("access_type", rule.AccessType)
			d.Set("access_level", rule.AccessLevel)
			ruleExist = true
		}
	}

	if accessID != "" && !ruleExist {
		log.Printf("[WARN] access rule (%s) of share file %s was not exist!", accessID, d.Id())
		d.Set("share_access_id", "")
	}
	d.Set("access_rules", allAccessRules)

	if len(rules) != 0 {
		d.Set("status", n.Status)
	} else {
		// The file system is not bind with any VPC.
		d.Set("status", "unavailable")
	}

	// set tags
	resourceTags, err := tags.Get(sfsClient, "sfs", d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error fetching tags of sfs: %s", err)
	}
	tagmap := tagsToMap(resourceTags.Tags)
	d.Set("tags", tagmap)

	return nil
}

func resourceSFSFileSystemV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud Share File Client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := shares.UpdateOpts{
			DisplayName:        d.Get("name").(string),
			DisplayDescription: d.Get("description").(string),
		}
		_, err = shares.Update(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating Huaweicloud Share File: %s", err)
		}
	}

	if d.HasChanges("access_to", "access_level", "access_type") {
		ruleID := d.Get("share_access_id").(string)
		if ruleID != "" {
			deleteAccessOpts := shares.DeleteAccessOpts{AccessID: ruleID}
			deny := shares.DeleteAccess(sfsClient, d.Id(), deleteAccessOpts)
			if deny.Err != nil {
				return fmt.Errorf("Error changing access rules for share file : %s", deny.Err)
			}
			d.Set("share_access_id", "")
		}

		if _, ok := d.GetOk("access_to"); ok {
			grantAccessOpts := shares.GrantAccessOpts{
				AccessLevel: d.Get("access_level").(string),
				AccessType:  d.Get("access_type").(string),
				AccessTo:    d.Get("access_to").(string),
			}

			log.Printf("[DEBUG] Grant Access Rules: %#v", grantAccessOpts)
			grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()
			if accessErr != nil {
				return fmt.Errorf("Error changing access rules for share file : %s", accessErr)
			}
			d.Set("share_access_id", grant.ID)
		}
	}

	if d.HasChange("size") {
		old, newsize := d.GetChange("size")
		if old.(int) < newsize.(int) {
			expandOpts := shares.ExpandOpts{OSExtend: shares.OSExtendOpts{NewSize: newsize.(int)}}
			expand := shares.Expand(sfsClient, d.Id(), expandOpts)
			if expand.Err != nil {
				return fmt.Errorf("Error Expanding Huaweicloud Share File size: %s", expand.Err)
			}
		} else {
			shrinkOpts := shares.ShrinkOpts{OSShrink: shares.OSShrinkOpts{NewSize: newsize.(int)}}
			shrink := shares.Shrink(sfsClient, d.Id(), shrinkOpts)
			if shrink.Err != nil {
				return fmt.Errorf("Error Shrinking Huaweicloud Share File size: %s", shrink.Err)
			}
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(sfsClient, d, "sfs", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of sfs:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceSFSFileSystemV2Read(d, meta)
}

func resourceSFSFileSystemV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.SfsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud Shared File Client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSFileDelete(sfsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud Share File: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSFileActive(sfsClient *golangsdk.ServiceClient, shareID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := shares.Get(sfsClient, shareID).Extract()
		if err != nil {
			return nil, "", err
		}

		if n.Status == "error" {
			return n, n.Status, nil
		}

		return n, n.Status, nil
	}
}

func waitForSFSFileDelete(sfsClient *golangsdk.ServiceClient, shareId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {

		r, err := shares.Get(sfsClient, shareId).Extract()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud shared File %s", shareId)
				return r, "deleted", nil
			}
			return r, "available", err
		}
		err = shares.Delete(sfsClient, shareId).ExtractErr()

		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[INFO] Successfully deleted Huaweicloud shared File %s", shareId)
				return r, "deleted", nil
			}
			if errCode, ok := err.(golangsdk.ErrUnexpectedResponseCode); ok {
				if errCode.Actual == 409 {
					return r, "available", nil
				}
			}
			return r, "available", err
		}

		return r, r.Status, nil
	}
}

func resourceSFSMetadataV2(d *schema.ResourceData, config *Config) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}

	epsID := GetEnterpriseProjectID(d, config)

	if epsID != "" {
		m["enterprise_project_id"] = epsID
	}

	return m
}

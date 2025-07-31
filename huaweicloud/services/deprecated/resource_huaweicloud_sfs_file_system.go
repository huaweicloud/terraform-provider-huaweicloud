package deprecated

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/common/tags"
	"github.com/chnsz/golangsdk/openstack/sfs/v2/shares"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFS POST /v2/{project_id}/shares/{id}/action
// @API SFS DELETE /v2/{project_id}/shares/{id}
// @API SFS GET /v2/{project_id}/shares/{id}
// @API SFS PUT /v2/{project_id}/shares/{id}
// @API SFS POST /v2/{project_id}/shares
// @API SFS POST /v2/{project_id}/sfs/{id}/tags/action
// @API SFS GET /v2/{project_id}/sfs/{id}/tags
func ResourceSFSFileSystemV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSFSFileSystemV2Create,
		ReadContext:   resourceSFSFileSystemV2Read,
		UpdateContext: resourceSFSFileSystemV2Update,
		DeleteContext: resourceSFSFileSystemV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.MergeDefaultTags(),

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
				Computed: true,
			},
			"access_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"access_to": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
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

func resourceSFSFileSystemV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))

	if err != nil {
		return diag.Errorf("error creating SFS client: %s", err)
	}

	createOpts := shares.CreateOpts{
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		IsPublic:         d.Get("is_public").(bool),
		Metadata:         resourceSFSMetadataV2(d, cfg),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	create, err := shares.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating SFS file system: %s", err)
	}

	d.SetId(create.ID)

	log.Printf("[DEBUG] Waiting for SFS file system (%s) to be available", create.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"creating"},
		Target:     []string{"available"},
		Refresh:    waitForSFSFileRefresh(sfsClient, create.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, stateErr := stateConf.WaitForStateContext(ctx)
	if stateErr != nil {
		return diag.Errorf("error waiting for SFS file system (%s) to be available: %s ", d.Id(), stateErr)
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
			return diag.Errorf("error applying access rule to SFS file system : %s", accessErr)
		}

		log.Printf("[DEBUG] Applied access rule (%s) to SFS file system %s", grant.ID, d.Id())
		d.Set("share_access_id", grant.ID)
	}

	// create tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := utils.ExpandResourceTags(tagRaw)
		if tagErr := tags.Create(sfsClient, "sfs", d.Id(), taglist).ExtractErr(); tagErr != nil {
			return diag.Errorf("error setting tags for sfs %s: %s", d.Id(), tagErr)
		}
	}

	return resourceSFSFileSystemV2Read(ctx, d, meta)
}

func resourceSFSFileSystemV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	sfsClient, err := cfg.SfsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	n, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("error retrieving SFS file system: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", n.Name),
		d.Set("share_proto", n.ShareProto),
		d.Set("size", n.Size),
		d.Set("description", n.Description),
		d.Set("is_public", n.IsPublic),
		d.Set("availability_zone", n.AvailabilityZone),
		d.Set("region", region),
		d.Set("export_location", n.ExportLocation),
		d.Set("enterprise_project_id", n.Metadata["enterprise_project_id"]),
	)

	// NOTE: only support the following metadata key
	var metaKeys = [3]string{"#sfs_crypt_key_id", "#sfs_crypt_domain_id", "#sfs_crypt_alias"}
	md := make(map[string]string)

	for key, val := range n.Metadata {
		for i := range metaKeys {
			if key == metaKeys[i] {
				md[key] = val
				break
			}
		}
	}
	mErr = multierror.Append(mErr, d.Set("metadata", md))

	// list access rules
	rules, err := shares.ListAccessRights(sfsClient, d.Id()).ExtractAccessRights()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			d.SetId("")
			return nil
		}
		return diag.Errorf("error retrieving SFS file system: %s", err)
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
			mErr = multierror.Append(mErr,
				d.Set("access_rule_status", rule.State),
				d.Set("access_to", rule.AccessTo),
				d.Set("access_type", rule.AccessType),
				d.Set("access_level", rule.AccessLevel),
			)

			ruleExist = true
		}
	}

	if accessID != "" && !ruleExist {
		log.Printf("[WARN] access rule (%s) of SFS file system %s was not exist!", accessID, d.Id())
		mErr = multierror.Append(mErr, d.Set("share_access_id", ""))
	}
	mErr = multierror.Append(mErr, d.Set("access_rules", allAccessRules))

	if len(rules) != 0 {
		mErr = multierror.Append(mErr, d.Set("status", n.Status))
	} else {
		// The file system is not bind with any VPC.
		mErr = multierror.Append(mErr, d.Set("status", "unavailable"))
	}

	// set tags
	if resourceTags, err := tags.Get(sfsClient, "sfs", d.Id()).Extract(); err == nil {
		tagmap := utils.TagsToMap(resourceTags.Tags)
		if err := d.Set("tags", tagmap); err != nil {
			return diag.Errorf("error saving tags to state for SFS file system (%s): %s", d.Id(), err)
		}
	} else {
		log.Printf("[WARN] Error fetching tags of SFS file system (%s): %s", d.Id(), err)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSFSFileSystemV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS Client: %s", err)
	}

	if d.HasChanges("name", "description") {
		updateOpts := shares.UpdateOpts{
			DisplayName:        d.Get("name").(string),
			DisplayDescription: d.Get("description").(string),
		}
		_, err = shares.Update(sfsClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating SFS file system: %s", err)
		}
	}

	if d.HasChanges("access_to", "access_level", "access_type") {
		if ruleID, ok := d.GetOk("share_access_id"); ok {
			deleteAccessOpts := shares.DeleteAccessOpts{AccessID: ruleID.(string)}
			deny := shares.DeleteAccess(sfsClient, d.Id(), deleteAccessOpts)
			if deny.Err != nil {
				return diag.Errorf("error changing access rules for SFS file system: %s", deny.Err)
			}
			d.Set("share_access_id", "")
		}

		if v, ok := d.GetOk("access_to"); ok {
			grantAccessOpts := shares.GrantAccessOpts{
				AccessTo: v.(string),
			}

			if v, ok := d.GetOk("access_level"); ok {
				grantAccessOpts.AccessLevel = v.(string)
			} else {
				grantAccessOpts.AccessLevel = "rw"
			}

			if v, ok := d.GetOk("access_type"); ok {
				grantAccessOpts.AccessType = v.(string)
			} else {
				grantAccessOpts.AccessType = "cert"
			}

			log.Printf("[DEBUG] Grant Access Rules: %#v", grantAccessOpts)
			grant, accessErr := shares.GrantAccess(sfsClient, d.Id(), grantAccessOpts).ExtractAccess()
			if accessErr != nil {
				return diag.Errorf("error changing access rules for share file : %s", accessErr)
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
				return diag.Errorf("error expanding SFS file system size: %s", expand.Err)
			}
		} else {
			shrinkOpts := shares.ShrinkOpts{OSShrink: shares.OSShrinkOpts{NewSize: newsize.(int)}}
			shrink := shares.Shrink(sfsClient, d.Id(), shrinkOpts)
			if shrink.Err != nil {
				return diag.Errorf("error shrinking SFS file system size: %s", shrink.Err)
			}
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := utils.UpdateResourceTags(sfsClient, d, "sfs", d.Id())
		if tagErr != nil {
			return diag.Errorf("error updating tags of sfs:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceSFSFileSystemV2Read(ctx, d, meta)
}

func resourceSFSFileSystemV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	sfsClient, err := cfg.SfsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SFS client: %s", err)
	}

	err = shares.Delete(sfsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("error deleting SFS file system: %s", err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{"available", "deleting"},
		Target:     []string{"deleted"},
		Refresh:    waitForSFSFileRefresh(sfsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("timeout waiting for SFS file system deletion to complete %s", err)
	}

	d.SetId("")
	return nil
}

func waitForSFSFileRefresh(sfsClient *golangsdk.ServiceClient, shareID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := shares.Get(sfsClient, shareID).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return n, "deleted", nil
			}
			return nil, "", err
		}

		return n, n.Status, nil
	}
}

func resourceSFSMetadataV2(d *schema.ResourceData, cfg *config.Config) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}

	epsID := cfg.GetEnterpriseProjectID(d)

	if epsID != "" {
		m["enterprise_project_id"] = epsID
	}

	return m
}

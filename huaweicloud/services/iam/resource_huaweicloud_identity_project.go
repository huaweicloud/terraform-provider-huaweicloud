package iam

import (
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityProjectV3Create,
		Read:   resourceIdentityProjectV3Read,
		Update: resourceIdentityProjectV3Update,
		Delete: resourceIdentityProjectV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceIdentityProjectV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	createOpts := projects.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	project, err := projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud project: %s", err)
	}

	d.SetId(project.ID)

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	project, err := projects.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "project")
	}

	logp.Printf("[DEBUG] Retrieved Huaweicloud project: %#v", project)

	d.Set("name", project.Name)
	d.Set("description", project.Description)
	d.Set("parent_id", project.ParentID)
	d.Set("enabled", project.Enabled)

	return nil
}

func resourceIdentityProjectV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	var hasChange bool
	var updateOpts projects.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = description
	}

	if hasChange {
		_, err := projects.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating Huaweicloud project: %s", err)
		}
	}

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	err = projects.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting Huaweicloud project: %s", err)
	}

	return nil
}

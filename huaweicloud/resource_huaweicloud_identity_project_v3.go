package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/projects"
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceIdentityProjectV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	createOpts := projects.CreateOpts{
		Description: d.Get("description").(string),
		DomainID:    d.Get("domain_id").(string),
		Name:        d.Get("name").(string),
		ParentID:    d.Get("parent_id").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	project, err := projects.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud project: %s", err)
	}

	d.SetId(project.ID)

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	project, err := projects.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "project")
	}

	log.Printf("[DEBUG] Retrieved Huaweicloud project: %#v", project)

	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("is_domain", project.IsDomain)
	d.Set("name", project.Name)
	d.Set("parent_id", project.ParentID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityProjectV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud identity client: %s", err)
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
			return fmt.Errorf("Error updating Huaweicloud project: %s", err)
		}
	}

	return resourceIdentityProjectV3Read(d, meta)
}

func resourceIdentityProjectV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud identity client: %s", err)
	}

	err = projects.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting Huaweicloud project: %s", err)
	}

	return nil
}

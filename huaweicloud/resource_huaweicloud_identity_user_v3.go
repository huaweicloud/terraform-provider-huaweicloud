package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/users"
)

var userOptions = map[users.Option]string{
	users.IgnoreChangePasswordUponFirstUse: "ignore_change_password_upon_first_use",
	users.IgnorePasswordExpiry:             "ignore_password_expiry",
	users.IgnoreLockoutFailureAttempts:     "ignore_lockout_failure_attempts",
	users.MultiFactorAuthEnabled:           "multi_factor_auth_enabled",
}

func ResourceIdentityUserV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityUserV3Create,
		Read:   resourceIdentityUserV3Read,
		Update: resourceIdentityUserV3Update,
		Delete: resourceIdentityUserV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"default_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceIdentityUserV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	createOpts := users.CreateOpts{
		DefaultProjectID: d.Get("default_project_id").(string),
		DomainID:         d.Get("domain_id").(string),
		Enabled:          &enabled,
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Add password here so it wouldn't go in the above log entry
	createOpts.Password = d.Get("password").(string)

	user, err := users.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud user: %s", err)
	}

	d.SetId(user.ID)

	return resourceIdentityUserV3Read(d, meta)
}

func resourceIdentityUserV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	user, err := users.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "user")
	}

	log.Printf("[DEBUG] Retrieved HuaweiCloud user: %#v", user)

	d.Set("default_project_id", user.DefaultProjectID)
	d.Set("domain_id", user.DomainID)
	d.Set("enabled", user.Enabled)
	d.Set("name", user.Name)
	d.Set("description", user.Description)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceIdentityUserV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	var hasChange bool
	var updateOpts users.UpdateOpts

	if d.HasChange("default_project_id") {
		hasChange = true
		updateOpts.DefaultProjectID = d.Get("default_project_id").(string)
	}

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if hasChange {
		log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if d.HasChange("password") {
		hasChange = true
		updateOpts.Password = d.Get("password").(string)
	}

	if hasChange {
		_, err := users.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating HuaweiCloud user: %s", err)
		}
	}

	return resourceIdentityUserV3Read(d, meta)
}

func resourceIdentityUserV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = users.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud user: %s", err)
	}

	return nil
}

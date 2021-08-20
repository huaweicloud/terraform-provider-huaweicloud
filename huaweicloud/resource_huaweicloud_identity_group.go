package huaweicloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/groups"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceIdentityGroupV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIdentityGroupV3Create,
		Read:   resourceIdentityGroupV3Read,
		Update: resourceIdentityGroupV3Update,
		Delete: resourceIdentityGroupV3Delete,
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
		},
	}
}

func resourceIdentityGroupV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	createOpts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)

	group, err := groups.Create(identityClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud Group: %s", err)
	}

	d.SetId(group.ID)

	return resourceIdentityGroupV3Read(d, meta)
}

func resourceIdentityGroupV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	group, err := groups.Get(identityClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "group")
	}

	logp.Printf("[DEBUG] Retrieved HuaweiCloud Group: %#v", group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceIdentityGroupV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	var hasChange bool
	var updateOpts groups.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		updateOpts.Description = d.Get("description").(string)
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if hasChange {
		logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
	}

	if hasChange {
		_, err := groups.Update(identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmtp.Errorf("Error updating HuaweiCloud group: %s", err)
		}
	}

	return resourceIdentityGroupV3Read(d, meta)
}

func resourceIdentityGroupV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	err = groups.Delete(identityClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud group: %s", err)
	}

	return nil
}

package huaweicloud

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityRoleV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "display_name"},
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				AtLeastOneOf: []string{"name", "display_name"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policy": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityRoleV3Read performs the role lookup.
func dataSourceIdentityRoleV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	listOpts := roles.ListOpts{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var role roles.Role
	allPages, err := roles.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query roles: %s", err)
	}

	allRoles, err := roles.ExtractRoles(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve roles: %s", err)
	}

	if len(allRoles) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allRoles) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allRoles)
		return fmt.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}
	role = allRoles[0]

	log.Printf("[DEBUG] Single Role found: %s", role.ID)
	return dataSourceIdentityRoleV3Attributes(d, config, &role)
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityRoleV3Attributes(d *schema.ResourceData, config *config.Config, role *roles.Role) error {
	log.Printf("[DEBUG] huaweicloud_identity_role_v3 details: %#v", role)

	d.SetId(role.ID)
	d.Set("name", role.Name)
	d.Set("description", role.Description)
	d.Set("display_name", role.DisplayName)
	d.Set("catalog", role.Catalog)
	d.Set("type", role.Type)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmt.Errorf("Error marshalling policy: %s", err)
	}
	d.Set("policy", string(policy))

	return nil
}

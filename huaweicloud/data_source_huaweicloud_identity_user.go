package huaweicloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/users"
)

func dataSourceIdentityUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIdentityUserRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// dataSourceIdentityUserRead performs the iam user lookup.
func dataSourceIdentityUserRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	listOpts := users.ListOpts{
		Name:     d.Get("name").(string),
		DomainID: d.Get("domain_id").(string),
		Enabled:  &enabled,
	}
	log.Printf("[DEBUG] List Options: %#v", listOpts)

	allPages, err := users.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query users: %s", err)
	}

	allUsers, err := users.ExtractUsers(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve users: %s", err)
	}

	if len(allUsers) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allUsers) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allUsers)
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria, or set `most_recent` attribute to true.")
	}
	var user users.User
	user = allUsers[0]

	d.SetId(user.ID)
	d.Set("name", user.Name)
	d.Set("domain_id", user.DomainID)
	d.Set("enable", user.Enabled)
	d.Set("description", user.Description)

	return nil
}

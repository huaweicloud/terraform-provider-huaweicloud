package iam

import (
	"context"

	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func DataSourceIdentityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityGroupV3Read,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
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
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"password_expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"password_strength": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityGroupV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	listOpts := groups.ListOpts{
		Name: d.Get("name").(string),
	}

	allPages, err := groups.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to query groups: %s", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)

	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve groups: %s", err)
	}

	conditions := map[string]interface{}{}

	if v, ok := d.GetOk("id"); ok {
		conditions["id"] = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		conditions["description"] = v.(string)
	}

	var foundGroups []groups.Group

	for _, group := range allGroups {
		if groupsFilter(group, conditions) {
			foundGroups = append(foundGroups, group)
		}
	}

	if len(foundGroups) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(foundGroups) > 1 {
		return fmtp.DiagErrorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	group := foundGroups[0]

	d.SetId(group.ID)

	mErr := multierror.Append(nil,
		d.Set("domain_id", group.DomainID),
		d.Set("name", group.Name),
		d.Set("description", group.Description),
	)
	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("error setting identity group fields: %s", err)
	}

	// get users of this group
	allUsers, err := groups.ListUsers(identityClient, d.Id()).Extract()

	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve users: %s", err)
	}

	if len(allUsers) > 0 {
		users := make([]map[string]interface{}, 0, len(allUsers))
		for _, userInGroup := range allUsers {
			user := setUserAttributes(userInGroup)
			users = append(users, user)
		}

		d.Set("users", users)
	}

	return nil
}

func groupsFilter(group groups.Group, conditions map[string]interface{}) bool {
	if v, ok := conditions["id"]; ok && v != group.ID {
		return false
	}
	if v, ok := conditions["description"]; ok && v != group.Description {
		return false
	}
	return true
}

func setUserAttributes(userInGroup groups.User) map[string]interface{} {
	user := make(map[string]interface{})

	user["name"] = userInGroup.Name
	user["id"] = userInGroup.Id
	user["description"] = userInGroup.Description
	user["enabled"] = userInGroup.Enabled
	user["password_expires_at"] = userInGroup.PasswordExpiresAt
	user["password_status"] = userInGroup.PwdStatus
	user["password_strength"] = userInGroup.PwdStrength

	return user
}

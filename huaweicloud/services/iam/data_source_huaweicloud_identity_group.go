package iam

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v3/groups
// @API IAM GET /v3/groups/{group_id}/users
func DataSourceIdentityGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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

func DataSourceIdentityGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	listOpts := groups.ListOpts{
		Name: d.Get("name").(string),
	}

	allPages, err := groups.List(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to query IAM groups: %s", err)
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return diag.Errorf("unable to extract IAM groups: %s", err)
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
		if filterGroups(group, conditions) {
			foundGroups = append(foundGroups, group)
		}
	}

	if len(foundGroups) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(foundGroups) > 1 {
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}

	group := foundGroups[0]
	log.Printf("[DEBUG] retrieve IAM group: %#v", group)

	d.SetId(group.ID)
	mErr := multierror.Append(nil,
		d.Set("domain_id", group.DomainID),
		d.Set("name", group.Name),
		d.Set("description", group.Description),
		d.Set("users", flattenUsersInGroup(identityClient, group.ID)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM group fields: %s", err)
	}

	return nil
}

func filterGroups(group groups.Group, conditions map[string]interface{}) bool {
	if v, ok := conditions["id"]; ok && v != group.ID {
		return false
	}
	if v, ok := conditions["description"]; ok && v != group.Description {
		return false
	}
	return true
}

func flattenUsersInGroup(client *golangsdk.ServiceClient, groupID string) []map[string]interface{} {
	// get users in this group
	allUsers, err := groups.ListUsers(client, groupID).Extract()
	if err != nil {
		log.Printf("[WARN] unable to retrieve users in group %s: %s", groupID, err)
		return nil
	}

	if len(allUsers) > 0 {
		users := make([]map[string]interface{}, len(allUsers))
		for i, userInGroup := range allUsers {
			users[i] = map[string]interface{}{
				"id":                  userInGroup.Id,
				"name":                userInGroup.Name,
				"description":         userInGroup.Description,
				"enabled":             userInGroup.Enabled,
				"password_expires_at": userInGroup.PasswordExpiresAt,
				"password_status":     userInGroup.PwdStatus,
				"password_strength":   userInGroup.PwdStrength,
			}
		}
		return users
	}

	return nil
}

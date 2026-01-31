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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The name of the identity group.`,
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of the identity group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the identity group.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain the group belongs to.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The users the group contains.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the IAM user.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the IAM user.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the IAM user.`,
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the IAM user is enabled.`,
						},
						"password_expires_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the password will expire.`,
						},
						"password_status": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `The password status.`,
						},
						"password_strength": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The password strength.`,
						},
					},
				},
			},
		},
	}
}

func listIdentityGroups(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]groups.Group, error) {
	listOpts := groups.ListOpts{
		Name: d.Get("name").(string),
	}

	allPages, err := groups.List(client, listOpts).AllPages()
	if err != nil {
		return nil, err
	}

	allGroups, err := groups.ExtractGroups(allPages)
	if err != nil {
		return nil, err
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

	return foundGroups, nil
}

func DataSourceIdentityGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	foundGroups, err := listIdentityGroups(d, identityClient)
	if err != nil {
		return diag.Errorf("error listing IAM groups: %s", err)
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
	allUsers, err := groups.ListUsers(client, groupID).Extract()
	if err != nil {
		log.Printf("[WARN] unable to retrieve users in group %s: %s", groupID, err)
		return nil
	}

	result := make([]map[string]interface{}, 0, len(allUsers))
	for _, userInGroup := range allUsers {
		result = append(result, map[string]interface{}{
			"id":                  userInGroup.Id,
			"name":                userInGroup.Name,
			"description":         userInGroup.Description,
			"enabled":             userInGroup.Enabled,
			"password_expires_at": userInGroup.PasswordExpiresAt,
			"password_status":     userInGroup.PwdStatus,
			"password_strength":   userInGroup.PwdStrength,
		})
	}

	return result
}

package iam

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/groups"
	"github.com/chnsz/golangsdk/openstack/identity/v3/users"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityUsers is the impl of data/huaweicloud_identity_users
func DataSourceIdentityUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityUsersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
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
						"password_status": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"password_expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.IdentityV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	listOpts := users.ListOpts{
		Name:    d.Get("name").(string),
		Enabled: utils.Bool(d.Get("enabled").(bool)),
	}
	pages, err := users.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("error retrieving IAM user list: %v", err)
	}
	userList, err := users.ExtractUsers(pages)
	if err != nil {
		return diag.Errorf("error extracting IAM user objects: %v", err)
	}

	result, ids := flattenIAMUserList(client, userList)

	d.SetId(hashcode.Strings(ids))
	return diag.FromErr(d.Set("users", result))
}

func flattenIAMUserList(client *golangsdk.ServiceClient, userList []users.User) ([]map[string]interface{}, []string) {
	if len(userList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(userList))
	ids := make([]string, len(userList))
	for i, val := range userList {
		ids[i] = val.ID
		result[i] = map[string]interface{}{
			"id":                  val.ID,
			"name":                val.Name,
			"enabled":             val.Enabled,
			"description":         val.Description,
			"password_status":     val.PasswordStatus,
			"password_expires_at": val.PasswordExpiresAt.Format(time.RFC3339),
		}
		if groupNames, err := getUserOwnGroups(client, val.ID); err == nil {
			result[i]["groups"] = groupNames
		} else {
			log.Printf("[WARN] faied to get groups to which the user belongs: %s", err)
		}
	}
	return result, ids
}

func getUserOwnGroups(client *golangsdk.ServiceClient, id string) ([]string, error) {
	pages, err := users.ListGroups(client, id).AllPages()
	if err != nil {
		return nil, fmt.Errorf("error retrieving IAM user %s groups: %v", id, err)
	}
	allGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return nil, fmt.Errorf("error extracting IAM user %s group objects: %v", id, err)
	}

	groupNames := make([]string, len(allGroups))
	for i, g := range allGroups {
		groupNames[i] = g.Name
	}
	return groupNames, nil
}

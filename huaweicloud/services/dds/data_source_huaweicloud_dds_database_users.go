package dds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dds/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS GET /v3/{project_id}/instances/{instance_id}/db-user/detail
func DateSourceDDSDatabaseUser() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSDatabaseUserRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Optional: true,
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
						"db_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roles": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"db_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"privileges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     privilegeSchemaResource(),
						},
						"inherited_privileges": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     privilegeSchemaResource(),
						},
					},
				},
			},
		},
	}
}

func dataSourceDDSDatabaseUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	opts := users.ListOpts{
		DbName: d.Get("db_name").(string),
		Name:   d.Get("name").(string),
	}
	resp, err := users.List(client, d.Get("instance_id").(string), opts)
	if err != nil {
		return diag.Errorf("error retrieving user list: %s", err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("users", flattenDatabaseUsers(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDatabaseUsers(userList []users.UserResp) []map[string]interface{} {
	if len(userList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(userList))
	for i, user := range userList {
		result[i] = map[string]interface{}{
			"db_name":              user.DbName,
			"name":                 user.Name,
			"roles":                flattenDatabaseRoles(user.Roles),
			"privileges":           flattenDatabasePrivileges(user.Privileges),
			"inherited_privileges": flattenDatabasePrivileges(user.InheritedPrivileges),
		}
	}
	return result
}

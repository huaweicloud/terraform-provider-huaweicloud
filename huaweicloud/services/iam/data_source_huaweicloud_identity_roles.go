package iam

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v3/roles
func DataSourceIdentityRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityRolesRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"roles": {
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
						"display_name": {
							Type:     schema.TypeString,
							Computed: true,
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
				},
			},
		},
	}
}

// dataSourceIdentityRolesRead performs the role lookup.
func dataSourceIdentityRolesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	listOpts := roles.ListOpts{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)
	allPages, err := roles.ListWithPages(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to query IAM roles: %s", err)
	}

	allRoles, err := roles.ExtractOffsetRoles(allPages)
	if err != nil {
		return diag.Errorf("unable to retrieve IAM roles: %s", err)
	}

	if len(allRoles) < 1 {
		return diag.Errorf("your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	roleList := make([]map[string]interface{}, 0, len(allRoles))
	for _, role := range allRoles {
		roleMap, err := flattenIdentityRole(&role)
		if err != nil {
			return diag.FromErr(err)
		}
		roleList = append(roleList, roleMap)
	}

	d.SetId("identity_roles")
	if err := d.Set("roles", roleList); err != nil {
		return diag.Errorf("error setting roles: %s", err)
	}

	return nil
}

// flattenIdentityRole converts a Role struct into a map for Terraform state.
func flattenIdentityRole(role *roles.Role) (map[string]interface{}, error) {
	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":           role.ID,
		"name":         role.Name,
		"display_name": role.DisplayName,
		"description":  role.Description,
		"catalog":      role.Catalog,
		"type":         role.Type,
		"policy":       string(policy),
	}, nil
}

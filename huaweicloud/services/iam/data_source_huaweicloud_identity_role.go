package iam

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API IAM GET /v3/roles
// @API IAM GET /v3/roles/{role_id}
func DataSourceIdentityRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "display_name", "role_id"},
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "display_name", "role_id"},
			},
			"role_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name", "display_name"},
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

// dataSourceIdentityRoleRead performs the role lookup.
func dataSourceIdentityRoleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	roleId := d.Get("role_id").(string)
	if roleId != "" {
		role, err := roles.Get(identityClient, roleId).Extract()
		if err != nil {
			return diag.Errorf("error fetching role details: %s", err)
		}
		d.SetId(role.ID)
		return dataSourceIdentityRoleAttributes(d, role)
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

	if len(allRoles) > 1 {
		return diag.Errorf("your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}
	role := allRoles[0]
	log.Printf("[DEBUG] retrieve IAM role: %#v", role)

	d.SetId(role.ID)
	return dataSourceIdentityRoleAttributes(d, &role)
}

// dataSourceIdentityRoleAttributes populates the fields of an Role resource.
func dataSourceIdentityRoleAttributes(d *schema.ResourceData, role *roles.Role) diag.Diagnostics {
	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return diag.Errorf("error marshaling the policy of IAM role: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("role_id", role.ID),
		d.Set("name", role.Name),
		d.Set("description", role.Description),
		d.Set("display_name", role.DisplayName),
		d.Set("catalog", role.Catalog),
		d.Set("type", role.Type),
		d.Set("policy", string(policy)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting IAM role fields: %s", err)
	}
	return nil
}

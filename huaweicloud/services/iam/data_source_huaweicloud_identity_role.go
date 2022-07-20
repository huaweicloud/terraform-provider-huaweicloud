package iam

import (
	"context"
	"encoding/json"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceIdentityRoleV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityRoleV3Read,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "display_name"},
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
func dataSourceIdentityRoleV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	listOpts := roles.ListOpts{
		Name:        d.Get("name").(string),
		DisplayName: d.Get("display_name").(string),
	}

	logp.Printf("[DEBUG] List Options: %#v", listOpts)

	var role roles.Role
	allPages, err := roles.List(identityClient, listOpts).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to query roles: %s", err)
	}

	allRoles, err := roles.ExtractRoles(allPages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve roles: %s", err)
	}

	if len(allRoles) < 1 {
		return fmtp.DiagErrorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allRoles) > 1 {
		logp.Printf("[DEBUG] Multiple results found: %#v", allRoles)
		return fmtp.DiagErrorf("Your query returned more than one result. " +
			"Please try a more specific search criteria.")
	}
	role = allRoles[0]

	logp.Printf("[DEBUG] Single Role found: %s", role.ID)
	return dataSourceIdentityRoleV3Attributes(ctx, d, config, &role)
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityRoleV3Attributes(_ context.Context, d *schema.ResourceData, config *config.Config, role *roles.Role) diag.Diagnostics {
	logp.Printf("[DEBUG] huaweicloud_identity_role_v3 details: %#v", role)

	d.SetId(role.ID)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmtp.DiagErrorf("Error marshalling policy: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", role.Name),
		d.Set("description", role.Description),
		d.Set("display_name", role.DisplayName),
		d.Set("catalog", role.Catalog),
		d.Set("type", role.Type),
		d.Set("policy", string(policy)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity custom role fields: %s", err)
	}

	return nil
}

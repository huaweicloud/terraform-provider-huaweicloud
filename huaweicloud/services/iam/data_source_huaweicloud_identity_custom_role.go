package iam

import (
	"context"
	"encoding/json"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/policies"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func DataSourceIdentityCustomRole() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityCustomRoleRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"name", "id"},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"references": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"catalog": {
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

func dataSourceIdentityCustomRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	identityClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	allPages, err := policies.List(identityClient).AllPages()
	if err != nil {
		return fmtp.DiagErrorf("Unable to query roles: %s", err)
	}

	roles, err := policies.ExtractPageRoles(allPages)
	if err != nil {
		return fmtp.DiagErrorf("Unable to retrieve roles: %s", err)
	}

	conditions := map[string]interface{}{}

	if v, ok := d.GetOk("name"); ok {
		conditions["name"] = v.(string)
	}
	if v, ok := d.GetOk("id"); ok {
		conditions["id"] = v.(string)
	}
	if v, ok := d.GetOk("domain_id"); ok {
		conditions["domain_id"] = v.(string)
	}
	if v, ok := d.GetOk("references"); ok {
		conditions["references"] = v.(int)
	}
	if v, ok := d.GetOk("description"); ok {
		conditions["description"] = v.(string)
	}
	if v, ok := d.GetOk("type"); ok {
		conditions["type"] = v.(string)
	}

	var allRoles []policies.Role

	for _, role := range roles {
		if rolesFilter(role, conditions) {
			allRoles = append(allRoles, role)
		}
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
	role := allRoles[0]

	return dataSourceIdentityCustomRoleAttributes(ctx, d, config, &role)
}

// dataSourceIdentityRoleV3Attributes populates the fields of an Role resource.
func dataSourceIdentityCustomRoleAttributes(_ context.Context, d *schema.ResourceData, config *config.Config, role *policies.Role) diag.Diagnostics {
	logp.Printf("[DEBUG] huaweicloud_identity_role details: %#v", role)

	d.SetId(role.ID)

	policy, err := json.Marshal(role.Policy)
	if err != nil {
		return fmtp.DiagErrorf("Error marshalling policy: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", role.Name),
		d.Set("domain_id", role.DomainId),
		d.Set("references", role.References),
		d.Set("catalog", role.Catalog),
		d.Set("description", role.Description),
		d.Set("type", role.Type),
		d.Set("policy", string(policy)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity custom role fields: %s", err)
	}

	return nil
}

func rolesFilter(role policies.Role, conditions map[string]interface{}) bool {
	if v, ok := conditions["name"]; ok && v != role.Name {
		return false
	}
	if v, ok := conditions["id"]; ok && v != role.ID {
		return false
	}
	if v, ok := conditions["domain_id"]; ok && v != role.DomainId {
		return false
	}
	if v, ok := conditions["references"]; ok && v != role.References {
		return false
	}
	if v, ok := conditions["description"]; ok && v != role.Description {
		return false
	}
	if v, ok := conditions["type"]; ok && v != role.Type {
		return false
	}
	return true
}

package iam

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
)

// @API IAM GET /v3/roles
func DataSourceIdentityPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityPermissionsRead,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "system",
				ValidateFunc: validation.StringInSlice([]string{"system", "system-role", "system-policy", "custom"}, false),
			},
			"scope_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"catalog": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"permissions": {
				Type:     schema.TypeList,
				Elem:     iamPermissionSchema(),
				Computed: true,
			},
		},
	}
}

func iamPermissionSchema() *schema.Resource {
	sc := schema.Resource{
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
			"description_cn": {
				Type:     schema.TypeString,
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
	return &sc
}

// dataSourceIdentityPermissionsRead performs the permissions lookup
func dataSourceIdentityPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	identityClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating identity client: %s", err)
	}

	listOpts := roles.ListOpts{
		DisplayName: d.Get("name").(string),
		Catalog:     d.Get("catalog").(string),
		Type:        d.Get("scope_type").(string),
	}

	roleType := d.Get("type").(string)
	if roleType == "custom" {
		if cfg.DomainID == "" {
			return diag.Errorf("the domain_id must be specified in the provider configuration")
		}

		// if DomainID is specified, only custom policies of the account will be returned.
		listOpts.DomainID = cfg.DomainID
	} else {
		// trim "system-" to get the PermissionType
		prefix := "system-"
		if strings.HasPrefix(roleType, prefix) {
			listOpts.PermissionType = roleType[len(prefix):]
		}
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)
	allPages, err := roles.ListWithPages(identityClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("unable to query IAM permissions: %s", err)
	}

	rawPermissions, err := roles.ExtractOffsetRoles(allPages)
	if err != nil {
		return diag.Errorf("unable to retrieve IAM permissions: %s", err)
	}

	allPermissions, ids := flattenIAMRoleList(rawPermissions)

	d.SetId(hashcode.Strings(ids))
	mErr := multierror.Append(nil,
		d.Set("permissions", allPermissions),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenIAMRoleList(roleList []roles.Role) ([]map[string]interface{}, []string) {
	if len(roleList) < 1 {
		return nil, nil
	}

	result := make([]map[string]interface{}, len(roleList))
	ids := make([]string, len(roleList))
	for i, val := range roleList {
		ids[i] = val.ID

		policy, err := json.Marshal(val.Policy)
		if err != nil {
			log.Printf("[WARN] failed to marshal the policy of IAM permission %s: %s", val.ID, err)
		}

		result[i] = map[string]interface{}{
			"id":             val.ID,
			"name":           val.DisplayName,
			"description":    val.Description,
			"description_cn": val.DescriptionCN,
			"catalog":        val.Catalog,
			"policy":         string(policy),
		}
	}
	return result, ids
}

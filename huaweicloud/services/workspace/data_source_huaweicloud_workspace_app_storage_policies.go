package workspace

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v1/{project_id}/storages-policy/actions/list-statements
func DataSourceAppStoragePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppStoragePoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the custom storage permission policies are located.",
			},
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the storage permission policy.",
						},
						"server_actions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The collection of permissions that server can use to access storage.",
						},
						"client_actions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The collection of permissions that client can use to access storage.",
						},
					},
				},
			},
		},
	}
}

func flattenAppStoragePolicies(policies []interface{}) []interface{} {
	result := make([]interface{}, 0, len(policies))

	for _, val := range policies {
		result = append(result, map[string]interface{}{
			"id":             utils.PathSearch("policy_statement_id", val, nil),
			"server_actions": utils.PathSearch("roam_actions", val, make([]interface{}, 0)),
			"client_actions": utils.PathSearch("actions", val, make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourceAppStoragePoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	policies, err := listAppStoragePermissionPolicies(client)
	if err != nil {
		// API error already formated in the list method.
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenAppStoragePolicies(policies)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("unable to setting data source fields of the storage permission policies: %s", err)
	}
	return nil
}

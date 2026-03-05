package organizations

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations GET /v1/organizations/roots
func DataSourceOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationRead,
		Schema: map[string]*schema.Schema{
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the organization.`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the organization's management account.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the organization was created.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the root.`,
			},
			"root_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the root.`,
			},
			"root_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The urn of the root.`,
			},
			"root_tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The key/value pairs attached to the root.`,
			},
			"enabled_policy_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of enabled Organizations policy types in the Organization Root.`,
			},
		},
	}
}

func dataSourceOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	getOrganizationClient, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	respBody, err := GetOrganization(getOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}

	getRootRespBody, err := getRoot(getOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}

	organizationId := utils.PathSearch("organization.id", respBody, "").(string)
	if organizationId == "" {
		return diag.Errorf("unable to find the organization ID from the API response")
	}

	d.SetId(organizationId)

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)
	policyTypes := utils.PathSearch("roots|[0].policy_types[?status=='enabled'].type", getRootRespBody,
		make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("urn", utils.PathSearch("organization.urn", respBody, nil)),
		d.Set("master_account_id", utils.PathSearch("organization.management_account_id", respBody, nil)),
		d.Set("master_account_name", utils.PathSearch("organization.management_account_name", respBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at", respBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name", getRootRespBody, nil)),
		d.Set("root_urn", utils.PathSearch("roots|[0].urn", getRootRespBody, nil)),
		d.Set("enabled_policy_types", policyTypes),
	)

	tagMap, err := getTags(getOrganizationClient, rootType, rootId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of root (%s): %s", rootId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("root_tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

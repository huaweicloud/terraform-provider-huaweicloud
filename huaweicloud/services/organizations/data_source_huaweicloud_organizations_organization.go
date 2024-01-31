// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

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
				Description: `Indicates the uniform resource name of the organization.`,
			},
			"master_account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the organization's management account.`,
			},
			"master_account_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the organization's management account.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the organization was created.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the root.`,
			},
			"root_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the root.`,
			},
			"root_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the urn of the root.`,
			},
			"root_tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the key/value attached to the root.`,
			},
			"enabled_policy_types": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Indicates the list of enabled Organizations policy types in the Organization Root.`,
			},
		},
	}
}

func dataSourceOrganizationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganization: Query Organizations organization
	var (
		getOrganizationProduct = "organizations"
	)
	getOrganizationClient, err := cfg.NewServiceClient(getOrganizationProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationRespBody, err := getOrganization(getOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}

	getRootRespBody, err := getRoot(getOrganizationClient)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("organization.id", getOrganizationRespBody, "")
	d.SetId(id.(string))

	rootId := utils.PathSearch("roots|[0].id", getRootRespBody, "").(string)

	policyTypes := utils.PathSearch("roots|[0].policy_types[?status=='enabled'].type", getRootRespBody,
		make([]interface{}, 0)).([]interface{})

	mErr = multierror.Append(
		mErr,
		d.Set("urn", utils.PathSearch("organization.urn", getOrganizationRespBody, nil)),
		d.Set("master_account_id", utils.PathSearch("organization.management_account_id",
			getOrganizationRespBody, nil)),
		d.Set("master_account_name", utils.PathSearch("organization.management_account_name",
			getOrganizationRespBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at",
			getOrganizationRespBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name",
			getRootRespBody, nil)),
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

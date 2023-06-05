// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceOrganizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrganizationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the organizations.`,
			},
			"account_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the organization's management account.`,
			},
			"account_name": {
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
				Description: `Indicates the key/value to attach to the organization.`,
			},
		},
	}
}

func dataSourceOrganizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getOrganizations: Query Organizations
	var (
		getOrganizationsProduct = "organizations"
	)
	getOrganizationsClient, err := cfg.NewServiceClient(getOrganizationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations Client: %s", err)
	}

	getOrganizationsRespBody, diagErr := getOrganizations(d, getOrganizationsClient)
	if diagErr != nil {
		return diagErr
	}

	getOrganizationsRootRespBody, diagErr := getOrganizationsRoot(d, getOrganizationsClient)
	if diagErr != nil {
		return diagErr
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	rootId := utils.PathSearch("roots|[0].id", getOrganizationsRootRespBody, "").(string)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("urn", utils.PathSearch("organization.urn", getOrganizationsRespBody, nil)),
		d.Set("account_id", utils.PathSearch("organization.management_account_id",
			getOrganizationsRespBody, nil)),
		d.Set("account_name", utils.PathSearch("organization.management_account_name",
			getOrganizationsRespBody, nil)),
		d.Set("created_at", utils.PathSearch("organization.created_at",
			getOrganizationsRespBody, nil)),
		d.Set("root_id", rootId),
		d.Set("root_name", utils.PathSearch("roots|[0].name",
			getOrganizationsRootRespBody, nil)),
		d.Set("root_urn", utils.PathSearch("roots|[0].urn", getOrganizationsRootRespBody, nil)),
	)

	tagMap, err := getTags(getOrganizationsClient, rootType, rootId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations root (%s): %s", rootId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("root_tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

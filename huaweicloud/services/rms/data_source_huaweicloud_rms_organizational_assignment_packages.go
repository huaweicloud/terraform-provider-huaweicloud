// Generated by PMS #88
package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceRmsOrganizationalAssignmentPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRmsOrganizationalAssignmentPackagesRead,

		Schema: map[string]*schema.Schema{
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the organization.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the organizational assignment package name.`,
			},
			"package_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the organizational assignment package ID.`,
			},
			"packages": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: `The list of organizational assignment packages.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The organizational assignment package name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The organizational assignment package ID.`,
						},
						"organization_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the organization.`,
						},
						"owner_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the organizational assignment package.`,
						},
						"org_conformance_pack_urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique identifier of the organizational assignment package.`,
						},
						"vars_structure": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: `The parameters of the organizational assignment package.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"var_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of a parameter.`,
									},
									"var_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The value of a parameter.`,
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the organizational assignment package.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the organizational assignment package.`,
						},
					},
				},
			},
		},
	}
}

type OrganizationalAssignmentPackagesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newOrganizationalAssignmentPackagesDSWrapper(d *schema.ResourceData, meta interface{}) *OrganizationalAssignmentPackagesDSWrapper {
	return &OrganizationalAssignmentPackagesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRmsOrganizationalAssignmentPackagesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newOrganizationalAssignmentPackagesDSWrapper(d, meta)
	lisOrgConPacRst, err := wrapper.ListOrganizationConformancePacks()
	if err != nil {
		return diag.FromErr(err)
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	err = wrapper.listOrganizationConformancePacksToSchema(lisOrgConPacRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CONFIG GET /v1/resource-manager/organizations/{organization_id}/conformance-packs
func (w *OrganizationalAssignmentPackagesDSWrapper) ListOrganizationConformancePacks() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rms")
	if err != nil {
		return nil, err
	}

	d := w.ResourceData
	uri := "/v1/resource-manager/organizations/{organization_id}/conformance-packs"
	uri = strings.ReplaceAll(uri, "{organization_id}", d.Get("organization_id").(string))
	params := map[string]any{
		"organization_conformance_pack_id": w.Get("package_id"),
		"conformance_pack_name":            w.Get("name"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		MarkerPager("organization_conformance_packs", "page_info.next_marker", "marker").
		Request().
		Result()
}

func (w *OrganizationalAssignmentPackagesDSWrapper) listOrganizationConformancePacksToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("packages", schemas.SliceToList(body.Get("organization_conformance_packs"),
			func(pac gjson.Result) any {
				return map[string]any{
					"name":                     pac.Get("org_conformance_pack_name").Value(),
					"id":                       pac.Get("org_conformance_pack_id").Value(),
					"organization_id":          pac.Get("organization_id").Value(),
					"owner_id":                 pac.Get("owner_id").Value(),
					"org_conformance_pack_urn": pac.Get("org_conformance_pack_urn").Value(),
					"vars_structure": schemas.SliceToList(pac.Get("vars_structure"),
						func(varStr gjson.Result) any {
							return map[string]any{
								"var_key":   varStr.Get("var_key").Value(),
								"var_value": utils.JsonToString(varStr.Get("var_value").Value()),
							}
						},
					),
					"created_at": pac.Get("created_at").Value(),
					"updated_at": pac.Get("updated_at").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

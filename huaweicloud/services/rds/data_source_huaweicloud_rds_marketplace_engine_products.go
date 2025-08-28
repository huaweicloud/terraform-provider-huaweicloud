package rds

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

func DataSourceRdsMarketplaceEngineProducts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsMarketplaceEngineProductsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bp_domain_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"marketplace_engine_products": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bp_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bp_domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"image_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_license_agreement": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agreements": {
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
									"language": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"version": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"provision_url": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type MarketplaceEngineProductsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newMarketplaceEngineProductsDSWrapper(d *schema.ResourceData, meta interface{}) *MarketplaceEngineProductsDSWrapper {
	return &MarketplaceEngineProductsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceRdsMarketplaceEngineProductsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newMarketplaceEngineProductsDSWrapper(d, meta)
	lisMarEngProRst, err := wrapper.ListMarketplaceEngineProducts()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listMarketplaceEngineProductsToSchema(lisMarEngProRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API RDS GET /v3/{project_id}/business-partner/{bp_domain_id}
func (w *MarketplaceEngineProductsDSWrapper) ListMarketplaceEngineProducts() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "rds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/business-partner/{bp_domain_id}"
	uri = strings.ReplaceAll(uri, "{bp_domain_id}", w.Get("bp_domain_id").(string))
	params := map[string]any{
		"engine_id": w.Get("engine_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		OffsetPager("marketplace_engine_products", "offset", "limit", 0).
		Request().
		Result()
}

func (w *MarketplaceEngineProductsDSWrapper) listMarketplaceEngineProductsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("marketplace_engine_products", schemas.SliceToList(body.Get("marketplace_engine_products"),
			func(marEngProducts gjson.Result) any {
				return map[string]any{
					"agreements": schemas.SliceToList(marEngProducts.Get("agreements"),
						func(agreements gjson.Result) any {
							return map[string]any{
								"id":            agreements.Get("id").Value(),
								"language":      agreements.Get("language").Value(),
								"name":          agreements.Get("name").Value(),
								"provision_url": agreements.Get("provision_url").Value(),
								"version":       agreements.Get("version").Value(),
							}
						},
					),
					"bp_domain_id":           marEngProducts.Get("bp_domain_id").Value(),
					"bp_name":                marEngProducts.Get("bp_name").Value(),
					"engine_id":              marEngProducts.Get("engine_id").Value(),
					"engine_version":         marEngProducts.Get("engine_version").Value(),
					"image_id":               marEngProducts.Get("image_id").Value(),
					"instance_mode":          marEngProducts.Get("instance_mode").Value(),
					"product_id":             marEngProducts.Get("product_id").Value(),
					"spec_code":              marEngProducts.Get("spec_code").Value(),
					"user_license_agreement": marEngProducts.Get("user_license_agreement").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

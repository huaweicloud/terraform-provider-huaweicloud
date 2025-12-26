package cce

import (
	"context"

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

// @API CCE GET /api/v2/flavor/specifications
func DataSourceCCEFlavorSpecifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEFlavorSpecificationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_flavor_specs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_capacity": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"is_sold_out": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"is_support_multi_az": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"available_master_flavors": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"azs": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"az_fault_domains": {
										Type:     schema.TypeMap,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeList,
											Elem: &schema.Schema{Type: schema.TypeString},
										},
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

type FlavorSpecificationsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newCCEFlavorSpecificationsDSWrapper(d *schema.ResourceData, meta interface{}) *FlavorSpecificationsDSWrapper {
	return &FlavorSpecificationsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCCEFlavorSpecificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newCCEFlavorSpecificationsDSWrapper(d, meta)
	flavorSpecificationsRst, err := wrapper.ShowFlavorSpecifications()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showFlavorSpecificationsToSchema(flavorSpecificationsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /api/v2/flavor/specifications
func (w *FlavorSpecificationsDSWrapper) ShowFlavorSpecifications() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/api/v2/flavor/specifications"
	params := map[string]any{
		"clusterType": w.Get("cluster_type"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *FlavorSpecificationsDSWrapper) showFlavorSpecificationsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("cluster_flavor_specs", schemas.ObjectToList(body.Get("clusterFlavorSpecs"),
			func(values gjson.Result) any {
				return map[string]any{
					"name":                values.Get("name").Value(),
					"node_capacity":       values.Get("nodeCapacity").Value(),
					"is_sold_out":         values.Get("isSoldOut").Value(),
					"is_support_multi_az": values.Get("isSupportMultiAZ").Value(),
					"available_master_flavors": schemas.ObjectToList(body.Get("availableMasterFlavors"),
						func(values gjson.Result) any {
							return map[string]any{
								"name":             values.Get("name").Value(),
								"azs":              schemas.SliceToStrList(values.Get("azs")),
								"az_fault_domains": schemas.MapConverter(values.Get("azFaultDomains"), schemas.SliceToStrList),
							}
						},
					),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

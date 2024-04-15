package cc

import (
	"context"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCcCentralNetworkCapabilities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworkCapabilitiesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"capability": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the capability of the central network.`,
			},
			"capabilities": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Central network capability list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `ID of the account that the central network belongs to.`,
						},
						"capability": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The capability of the central network.`,
						},
						"specifications": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"support_ipv6_regions": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The regions that support IPv6.`,
									},
									"is_support": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: `Whether the central network supports the capability.`,
									},
									"size_range": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"min": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The minimum size.`,
												},
												"max": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: `The maximum size.`,
												},
											},
										},
										Description: `The range of the size.`,
									},
									"support_dscp_regions": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The regions that support DSCP.`,
									},
									"support_sites": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The supported sites.`,
									},
									"support_regions": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The supported regions.`,
									},
									"charge_mode": {
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Computed:    true,
										Description: `The charge mode.`,
									},
								},
							},
							Description: `The specifications of the central network.`,
						},
					},
				},
			},
		},
	}
}

// @API CC GET /v3/{domain_id}/gcn/capabilities
func dataSourceCcCentralNetworkCapabilitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cc", region)

	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	httpUrl := "v3/{domain_id}/gcn/capabilities"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{domain_id}", cfg.DomainID)

	params := buildCentralNetworkCapabilitiesQueryParams(d, cfg)
	path += params

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", path, &opt)
	if err != nil {
		return diag.Errorf("error creating bandwidth package: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	capabilities := utils.PathSearch("capabilities", respBody, make([]interface{}, 0))

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("capabilities", flattenCentralNetworkCapabilities(capabilities.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCentralNetworkCapabilitiesQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	var params string
	if capability, ok := d.GetOk("capability"); ok {
		params += "?capability=" + capability.(string)
	}
	return params
}

func flattenCentralNetworkCapabilities(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))

	for _, item := range resp {
		specifications := utils.PathSearch("specifications", item, nil)
		rst = append(rst, map[string]interface{}{
			"domain_id":      utils.PathSearch("domain_id", item, nil),
			"capability":     utils.PathSearch("capability", item, nil),
			"specifications": flattenCentralNetworkCapabilitiesSpecifications(specifications),
		})
	}
	return rst
}

func flattenCentralNetworkCapabilitiesSpecifications(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	item := map[string]interface{}{}
	specifications := resp.(map[string]interface{})
	support_ipv6_regions := utils.PathSearch("support_ipv6_regions", specifications, nil)
	is_support := utils.PathSearch("is_support", specifications, nil)
	size_range := utils.PathSearch("size_range", specifications, nil)
	support_dscp_regions := utils.PathSearch("support_dscp_regions", specifications, nil)
	support_sites := utils.PathSearch("support_sites", specifications, nil)
	support_regions := utils.PathSearch("support_regions", specifications, nil)
	charge_mode := utils.PathSearch("charge_mode", specifications, nil)

	if support_ipv6_regions != nil {
		item["support_ipv6_regions"] = support_ipv6_regions
	}
	if is_support != nil {
		item["is_support"] = is_support
	}
	if size_range != nil {
		item["size_range"] = flattenCentralNetworkCapabilitiesSpecificationsSizeRange(size_range)
	}
	if support_dscp_regions != nil {
		item["support_dscp_regions"] = support_dscp_regions
	}
	if support_sites != nil {
		item["support_sites"] = support_sites
	}
	if support_regions != nil {
		item["support_regions"] = support_regions
	}
	if charge_mode != nil {
		item["charge_mode"] = charge_mode
	}
	rst = append(rst, item)
	return rst
}

func flattenCentralNetworkCapabilitiesSpecificationsSizeRange(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, 1)
	item := map[string]interface{}{}
	item["min"] = utils.PathSearch("min", resp, nil)
	item["max"] = utils.PathSearch("max", resp, nil)

	rst = append(rst, item)
	return rst
}

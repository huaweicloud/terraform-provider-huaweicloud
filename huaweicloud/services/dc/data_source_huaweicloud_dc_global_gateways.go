package dc

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/global-dc-gateways
func DataSourceDcGlobalGateways() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcGlobalGatewaysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"fields": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the list of fields to be displayed.`,
			},
			"sort_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting field.`,
			},
			"sort_dir": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting order of returned results.`,
			},
			"global_gateway_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource IDs for querying instances.`,
			},
			"names": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the resource names for querying instances.`,
			},
			"enterprise_project_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the enterprise project IDs for querying instances.`,
			},
			"site_network_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the site network IDs.`,
			},
			"cloud_connection_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the cloud connection IDs.`,
			},
			"statuses": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the statuses by which instances are filtered.`,
			},
			"global_center_network_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the central network IDs.`,
			},
			"gateways": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the global DC gateways.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"locales": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The locale address description information.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en_us": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The region name in English.`,
									},
									"zh_cn": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The region name in Chinese.`,
									},
								},
							},
						},
						"current_peer_link_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of peer links allowed on a global DC gateway.`,
						},
						"tags": common.TagsComputedSchema(),
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the global DC gateway.`,
						},
						"location_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The location where the underlying device of the global DC gateway is deployed.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the global DC gateway.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the global DC gateway.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID that the global DC gateway belongs to.`,
						},
						"global_center_network_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the central network that the global DC gateway is added to.`,
						},
						"bgp_asn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The BGP ASN of the global DC gateway.`,
						},
						"address_family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address family of the global DC gateway.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The global DC gateway ID.`,
						},
						"available_peer_link_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of peer links that can be created for a global DC gateway.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the global DC gateway was created.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the global DC gateway was updated.`,
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cause of the failure to create the global DC gateway.`,
						},
					},
				},
			},
		},
	}
}

func buildQueryStringParams(queryKey string, queryValues []interface{}) string {
	rst := ""
	for _, val := range queryValues {
		if queryValue, ok := val.(string); ok && queryValue != "" {
			rst += fmt.Sprintf("&%s=%s", queryKey, queryValue)
		}
	}

	return rst
}

func buildDcGlobalGatewaysQueryParams(d *schema.ResourceData) string {
	rst := "?limit=2000"
	if v, ok := d.GetOk("fields"); ok {
		rst += buildQueryStringParams("fields", v.([]interface{}))
	}

	if v, ok := d.GetOk("sort_key"); ok {
		rst += fmt.Sprintf("&sort_key=%v", v)
	}

	if v, ok := d.GetOk("sort_dir"); ok {
		rst += fmt.Sprintf("&sort_dir=%v", v)
	}

	if v, ok := d.GetOk("global_gateway_ids"); ok {
		rst += buildQueryStringParams("id", v.([]interface{}))
	}

	if v, ok := d.GetOk("names"); ok {
		rst += buildQueryStringParams("name", v.([]interface{}))
	}

	if v, ok := d.GetOk("enterprise_project_ids"); ok {
		rst += buildQueryStringParams("enterprise_project_id", v.([]interface{}))
	}

	if v, ok := d.GetOk("site_network_ids"); ok {
		rst += buildQueryStringParams("site_network_id", v.([]interface{}))
	}

	if v, ok := d.GetOk("cloud_connection_ids"); ok {
		rst += buildQueryStringParams("cloud_connection_id", v.([]interface{}))
	}

	if v, ok := d.GetOk("statuses"); ok {
		rst += buildQueryStringParams("status", v.([]interface{}))
	}

	if v, ok := d.GetOk("global_center_network_ids"); ok {
		rst += buildQueryStringParams("global_center_network_id", v.([]interface{}))
	}

	return rst
}

func buildDcGlobalGatewaysQueryParamsWithMarker(requestPath, marker string) string {
	if marker == "" {
		return requestPath
	}

	return fmt.Sprintf("%s&marker=%s", requestPath, marker)
}

func dataSourceDcGlobalGatewaysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v3/{project_id}/dcaas/global-dc-gateways"
		product     = "dc"
		marker      = ""
		allGateways []interface{}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDcGlobalGatewaysQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := buildDcGlobalGatewaysQueryParamsWithMarker(requestPath, marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DC global gateways: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		gateways := utils.PathSearch("global_dc_gateways", respBody, make([]interface{}, 0)).([]interface{})
		if len(gateways) == 0 {
			break
		}

		allGateways = append(allGateways, gateways...)

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("gateways", flattenDcGlobalGatewaysAttribute(allGateways)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDcGlobalGatewaysAttribute(allGateways []interface{}) []interface{} {
	if len(allGateways) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allGateways))
	for _, v := range allGateways {
		rst = append(rst, map[string]interface{}{
			"name":                      utils.PathSearch("name", v, nil),
			"location_name":             utils.PathSearch("location_name", v, nil),
			"status":                    utils.PathSearch("status", v, nil),
			"description":               utils.PathSearch("description", v, nil),
			"enterprise_project_id":     utils.PathSearch("enterprise_project_id", v, nil),
			"global_center_network_id":  utils.PathSearch("global_center_network_id", v, nil),
			"bgp_asn":                   utils.PathSearch("bgp_asn", v, nil),
			"address_family":            utils.PathSearch("address_family", v, nil),
			"id":                        utils.PathSearch("id", v, nil),
			"available_peer_link_count": utils.PathSearch("available_peer_link_count", v, nil),
			"created_time":              utils.PathSearch("created_time", v, nil),
			"updated_time":              utils.PathSearch("updated_time", v, nil),
			"reason":                    utils.PathSearch("reason", v, nil),
			"current_peer_link_count":   utils.PathSearch("current_peer_link_count", v, nil),
			"tags":                      utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
			"locales":                   flattenDcGlobalGatewayLocales(utils.PathSearch("locales", v, nil)),
		})
	}

	return rst
}

func flattenDcGlobalGatewayLocales(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	localeMap := map[string]interface{}{
		"en_us": utils.PathSearch("en_us", respBody, nil),
		"zh_cn": utils.PathSearch("zh_cn", respBody, nil),
	}

	return []interface{}{localeMap}
}

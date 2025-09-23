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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DC GET /v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links
func DataSourceDcGlobalGatewayPeerLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcGlobalGatewayPeerLinksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"global_dc_gateway_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the global DC gateway ID.`,
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
			"peer_link_ids": {
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
			"peer_links": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the peer links.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The peer link ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the peer link.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the peer link.`,
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cause of the failure to add the peer link.`,
						},
						"global_dc_gateway_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the global DC gateway that the peer link is added for.`,
						},
						"bandwidth_info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The bandwidth information.`,
							Elem:        datasourceBandwidthInfoSchema(),
						},
						"peer_site": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The site of the peer link.`,
							Elem:        datasourcePeerSiteSchema(),
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the peer link.`,
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the peer link was added.`,
						},
						"updated_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the peer link was updated.`,
						},
						"create_owner": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The cloud service where the peer link is used.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance associated with the peer link.`,
						},
					},
				},
			},
		},
	}
}

func datasourceBandwidthInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bandwidth_size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The bandwidth size.`,
			},
			"gcb_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The global connection bandwidth ID.`,
			},
		},
	}
}

func datasourcePeerSiteSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of enterprise router that the global DC gateway is attached to.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The project ID of the enterprise router that the global DC gateway is attached to.`,
			},
			"region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region ID of the enterprise router that the global DC gateway is attached to.`,
			},
			"link_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The connection ID of the peer gateway at the peer site.`,
			},
			"site_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The site information of the global DC gateway.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the peer gateway.`,
			},
		},
	}
}

func buildDcGlobalGatewayPeerLinksQueryParams(d *schema.ResourceData) string {
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

	if v, ok := d.GetOk("peer_link_ids"); ok {
		rst += buildQueryStringParams("id", v.([]interface{}))
	}

	if v, ok := d.GetOk("names"); ok {
		rst += buildQueryStringParams("name", v.([]interface{}))
	}

	return rst
}

func buildPeerLinksQueryParamsWithMarker(requestPath, marker string) string {
	if marker == "" {
		return requestPath
	}

	return fmt.Sprintf("%s&marker=%s", requestPath, marker)
}

func dataSourceDcGlobalGatewayPeerLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v3/{project_id}/dcaas/global-dc-gateways/{global_dc_gateway_id}/peer-links"
		product   = "dc"
		marker    = ""
		allValues []interface{}
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{global_dc_gateway_id}", d.Get("global_dc_gateway_id").(string))
	requestPath += buildDcGlobalGatewayPeerLinksQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		requestPathWithMarker := buildPeerLinksQueryParamsWithMarker(requestPath, marker)
		resp, err := client.Request("GET", requestPathWithMarker, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DC global gateway peer links: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		peerLinks := utils.PathSearch("peer_links", respBody, make([]interface{}, 0)).([]interface{})
		if len(peerLinks) == 0 {
			break
		}

		allValues = append(allValues, peerLinks...)

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
		d.Set("peer_links", flattenDcGlobalGatewayPeerLinksAttribute(allValues)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDcGlobalGatewayPeerLinksAttribute(allValues []interface{}) []interface{} {
	if len(allValues) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(allValues))
	for _, v := range allValues {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("id", v, nil),
			"name":                 utils.PathSearch("name", v, nil),
			"description":          utils.PathSearch("description", v, nil),
			"reason":               utils.PathSearch("reason", v, nil),
			"global_dc_gateway_id": utils.PathSearch("global_dc_gateway_id", v, nil),
			"bandwidth_info":       flattenDatasourceBandwidthInfoAttribute(v),
			"peer_site":            flattenDatasourcePeerSiteAttribute(v),
			"status":               utils.PathSearch("status", v, nil),
			"created_time":         utils.PathSearch("created_time", v, nil),
			"updated_time":         utils.PathSearch("updated_time", v, nil),
			"create_owner":         utils.PathSearch("create_owner", v, nil),
			"instance_id":          utils.PathSearch("instance_id", v, nil),
		})
	}

	return rst
}

func flattenDatasourceBandwidthInfoAttribute(respBody interface{}) []interface{} {
	bandwidthInfo := utils.PathSearch("bandwidth_info", respBody, nil)
	if bandwidthInfo == nil {
		return nil
	}

	rawMap := map[string]interface{}{
		"bandwidth_size": utils.PathSearch("bandwidth_size", bandwidthInfo, nil),
		"gcb_id":         utils.PathSearch("gcb_id", bandwidthInfo, nil),
	}
	return []interface{}{rawMap}
}

func flattenDatasourcePeerSiteAttribute(respBody interface{}) []interface{} {
	peerSite := utils.PathSearch("peer_site", respBody, nil)
	if peerSite == nil {
		return nil
	}

	rawMap := map[string]interface{}{
		"gateway_id": utils.PathSearch("gateway_id", peerSite, nil),
		"project_id": utils.PathSearch("project_id", peerSite, nil),
		"region_id":  utils.PathSearch("region_id", peerSite, nil),
		"link_id":    utils.PathSearch("link_id", peerSite, nil),
		"site_code":  utils.PathSearch("site_code", peerSite, nil),
		"type":       utils.PathSearch("type", peerSite, nil),
	}
	return []interface{}{rawMap}
}

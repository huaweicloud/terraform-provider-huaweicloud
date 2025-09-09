package cc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CC POST /v3/{domain_id}/gcn/central-networks/filter
func DataSourceCcCentralNetworksByTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCcCentralNetworksByTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     centralNetworksByTagsTagsSchema(),
			},
			"central_networks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksSchema(),
			},
		},
	}
}

func centralNetworksByTagsTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"values": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksTagsSchema(),
			},
			"default_plane_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_associate_route_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"auto_propagate_route_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"planes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksPlanesSchema(),
			},
			"er_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksErInstancesSchema(),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksTagsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksPlanesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"associate_er_tables": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksPlanesAssociateErTablesSchema(),
			},
			"exclude_er_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksPlanesExcludeErConnectionsSchema(),
			},
			"is_full_mesh": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksPlanesAssociateErTablesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_router_table_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksPlanesExcludeErConnectionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"exclude_er_instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     centralNetworksByTagsCentralNetworksPlanesExcludeErInstancesSchema(),
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksPlanesExcludeErInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func centralNetworksByTagsCentralNetworksErInstancesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enterprise_router_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"asn": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"site_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCcCentralNetworksByTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{domain_id}/gcn/central-networks/filter"
		product = "cc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CC client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = buildCcCentralNetworksByTagsQueryParams(d)

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.FromErr(err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("central_networks", flattenCentralNetworksByTagsCentralNetworksResponseBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCcCentralNetworksByTagsQueryParams(d *schema.ResourceData) map[string]interface{} {
	tags := d.Get("tags").([]interface{})
	if len(tags) == 0 {
		return nil
	}

	params := make([]map[string]interface{}, 0, len(tags))
	for _, v := range tags {
		raw := v.(map[string]interface{})
		params = append(params, map[string]interface{}{
			"key":    raw["key"],
			"values": raw["values"],
		})
	}
	bodyParams := map[string]interface{}{
		"tags": params,
	}

	return bodyParams
}

func flattenCentralNetworksByTagsCentralNetworksResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("central_networks", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})

	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                           utils.PathSearch("id", v, nil),
			"name":                         utils.PathSearch("name", v, nil),
			"description":                  utils.PathSearch("description", v, nil),
			"state":                        utils.PathSearch("state", v, nil),
			"enterprise_project_id":        utils.PathSearch("enterprise_project_id", v, nil),
			"tags":                         utils.PathSearch("tags", v, nil),
			"default_plane_id":             utils.PathSearch("default_plane_id", v, nil),
			"auto_associate_route_enabled": utils.PathSearch("auto_associate_route_enabled", v, nil),
			"auto_propagate_route_enabled": utils.PathSearch("auto_propagate_route_enabled", v, nil),
			"planes":                       flattenCentralNetworksPlanes(v),
			"er_instances":                 flattenCentralNetworksErInstances(v),
			"created_at":                   utils.PathSearch("created_at", v, nil),
			"updated_at":                   utils.PathSearch("updated_at", v, nil),
		})
	}

	return rst
}

func flattenCentralNetworksPlanes(resp interface{}) []interface{} {
	curJson := utils.PathSearch("planes", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                     utils.PathSearch("id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"associate_er_tables":    flattenCentralNetworksPlanesAssociateErTables(v),
			"exclude_er_connections": flattenCentralNetworksPlanesExcludeErConnections(v),
			"is_full_mesh":           utils.PathSearch("is_full_mesh", v, nil),
		})
	}
	return rst
}

func flattenCentralNetworksPlanesAssociateErTables(resp interface{}) []interface{} {
	curJson := utils.PathSearch("associate_er_tables", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_id":                 utils.PathSearch("project_id", v, nil),
			"region_id":                  utils.PathSearch("region_id", v, nil),
			"enterprise_router_id":       utils.PathSearch("enterprise_router_id", v, nil),
			"enterprise_router_table_id": utils.PathSearch("enterprise_router_table_id", v, nil),
		})
	}
	return rst
}

func flattenCentralNetworksPlanesExcludeErConnections(resp interface{}) []interface{} {
	curJson := utils.PathSearch("exclude_er_connections", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"exclude_er_instances": flattenCentralNetworksPlanesExcludeErInstances(v),
		})
	}
	return rst
}

func flattenCentralNetworksPlanesExcludeErInstances(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"project_id":           utils.PathSearch("project_id", v, nil),
			"region_id":            utils.PathSearch("region_id", v, nil),
			"enterprise_router_id": utils.PathSearch("enterprise_router_id", v, nil),
		})
	}
	return rst
}

func flattenCentralNetworksErInstances(resp interface{}) []interface{} {
	curJson := utils.PathSearch("er_instances", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"enterprise_router_id": utils.PathSearch("enterprise_router_id", v, nil),
			"project_id":           utils.PathSearch("project_id", v, nil),
			"region_id":            utils.PathSearch("region_id", v, nil),
			"asn":                  utils.PathSearch("asn", v, nil),
			"site_code":            utils.PathSearch("site_code", v, nil),
		})
	}
	return rst
}

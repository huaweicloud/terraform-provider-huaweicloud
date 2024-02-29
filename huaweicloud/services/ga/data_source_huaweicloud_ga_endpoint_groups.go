package ga

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GA GET /v1/endpoint-groups
func DataSourceEndpointGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceEndpointGroupsRead,
		Schema: map[string]*schema.Schema{
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the endpoint group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the endpoint group.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the endpoint group.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the listener to which the endpoint group belongs.",
			},
			"endpoint_groups": {
				Type:        schema.TypeList,
				Elem:        endpointGroupsSchema(),
				Computed:    true,
				Description: "The list of the endpoint groups.",
			},
		},
	}
}

func endpointGroupsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the endpoint group.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the endpoint group.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the endpoint group.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the endpoint group.",
			},
			"traffic_dial_percentage": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The percentage of traffic distributed to the endpoint group.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region where the endpoint group belongs.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the listener to which the endpoint group belongs.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the endpoint group.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the endpoint group.",
			},
		},
	}
	return &sc
}

func dataSourceEndpointGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listEndpointGroups: Query the list of endpoint groups
	var (
		listEndpointGroupsHttpUrl = "v1/endpoint-groups"
		listEndpointGroupsProduct = "ga"
	)
	listEndpointGroupsClient, err := cfg.NewServiceClient(listEndpointGroupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	listEndpointGroupsPath := listEndpointGroupsClient.Endpoint + listEndpointGroupsHttpUrl

	listEndpointGroupsqueryParams := buildListEndpointGroupsQueryParams(d)
	listEndpointGroupsPath += listEndpointGroupsqueryParams

	listEndpointGroupsResp, err := pagination.ListAllItems(
		listEndpointGroupsClient,
		"marker",
		listEndpointGroupsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving endpoint groups")
	}

	listEndpointGroupsRespJson, err := json.Marshal(listEndpointGroupsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listEndpointGroupsRespBody interface{}
	err = json.Unmarshal(listEndpointGroupsRespJson, &listEndpointGroupsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("endpoint_groups", flattenListEndpointGroupsResponseBody(listEndpointGroupsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEndpointGroupsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("endpoint_groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                      utils.PathSearch("id", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"status":                  utils.PathSearch("status", v, nil),
			"description":             utils.PathSearch("description", v, nil),
			"traffic_dial_percentage": utils.PathSearch("traffic_dial_percentage", v, float64(0)),
			"region_id":               utils.PathSearch("region_id", v, nil),
			"listener_id":             utils.PathSearch("listeners[0].id", v, nil),
			"created_at":              utils.PathSearch("created_at", v, nil),
			"updated_at":              utils.PathSearch("updated_at", v, nil),
		})
	}
	return rst
}

func buildListEndpointGroupsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("endpoint_group_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("listener_id"); ok {
		res = fmt.Sprintf("%s&listener_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

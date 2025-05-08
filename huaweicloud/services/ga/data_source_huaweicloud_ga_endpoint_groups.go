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
			"frozen_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The frozen details of cloud services or resources.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of a cloud service or resource.`,
						},
						"effect": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The status of the resource after being forzen.`,
						},
						"scene": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The service scenario.`,
						},
					},
				},
			},
		},
	}
	return &sc
}

func dataSourceEndpointGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/endpoint-groups"
		product = "ga"
		mErr    *multierror.Error
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += buildListEndpointGroupsQueryParams(d)
	resp, err := pagination.ListAllItems(
		client,
		"marker",
		requestPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving GA endpoint groups: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("endpoint_groups", flattenListEndpointGroupsResponseBody(respBody)),
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
			"frozen_info":             flattenEndpointGroupsFrozenInfo(utils.PathSearch("frozen_info", v, nil)),
		})
	}
	return rst
}

func flattenEndpointGroupsFrozenInfo(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	frozenInfo := map[string]interface{}{
		"status": utils.PathSearch("status", resp, nil),
		"effect": utils.PathSearch("effect", resp, nil),
		"scene":  utils.PathSearch("scene", resp, []string{}),
	}

	return []map[string]interface{}{frozenInfo}
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

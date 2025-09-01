package coc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API COC POST /v1/instances/batches
func DataSourceCocInstanceBatches() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCocInstanceBatchesRead,

		Schema: map[string]*schema.Schema{
			"batch_strategy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_instances": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cloud_service_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"custom_attributes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:     schema.TypeString,
										Required: true,
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
						"properties": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"fixed_ip": {
										Type:     schema.TypeString,
										Required: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"zone_id": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_index": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_instances": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cloud_service_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"region_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"custom_attributes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
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
										},
									},
									"properties": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"host_name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"fixed_ip": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"region_id": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"zone_id": {
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
				},
			},
		},
	}
}

func dataSourceCocInstanceBatchesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("coc", region)
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	instanceBatches, err := queryInstanceBatches(client, d)
	if err != nil {
		return diag.Errorf("error querying instance batches: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		nil,
		d.Set("data", flattenCocInstanceBatches(
			utils.PathSearch("data", instanceBatches, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func queryInstanceBatches(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	queryHttpUrl := "v1/instances/batches"
	queryPath := client.Endpoint + queryHttpUrl
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildQueryInstanceBatchesBodyParams(d)),
	}

	queryResp, err := client.Request("POST", queryPath, &queryOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying COC instance batches: %s", err)
	}

	queryRespBody, err := utils.FlattenResponse(queryResp)
	if err != nil {
		return nil, fmt.Errorf("error querying COC instance batches: %s", err)
	}

	return queryRespBody, nil
}

func buildQueryInstanceBatchesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"batch_strategy":   d.Get("batch_strategy"),
		"target_instances": buildQueryInstanceBatchesTargetInstancesBodyParams(d.Get("target_instances")),
	}

	return bodyParams
}

func buildQueryInstanceBatchesTargetInstancesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"resource_id": raw["resource_id"],
				"provider":    raw["cloud_service_name"],
				"region_id":   raw["region_id"],
				"type":        raw["type"],
				"custom_attributes": buildQueryInstanceBatchesTargetInstancesCustomAttributesBodyParams(
					raw["custom_attributes"]),
				"properties": buildQueryInstanceBatchesTargetInstancesPropertiesBodyParams(raw["properties"]),
			}
		}
		return params
	}

	return nil
}

func buildQueryInstanceBatchesTargetInstancesCustomAttributesBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"key":   raw["key"],
				"value": raw["value"],
			}
		}
		return params
	}

	return nil
}

func buildQueryInstanceBatchesTargetInstancesPropertiesBodyParams(rawParam interface{}) map[string]interface{} {
	if rawArray, ok := rawParam.([]interface{}); ok {
		if len(rawArray) != 1 {
			return nil
		}

		raw := rawArray[0].(map[string]interface{})
		param := map[string]interface{}{
			"host_name": raw["host_name"],
			"fixed_ip":  raw["fixed_ip"],
			"region_id": raw["region_id"],
			"zone_id":   raw["zone_id"],
		}

		return param
	}

	return nil
}

func flattenCocInstanceBatches(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"batch_index": utils.PathSearch("batch_index", raw, nil),
				"target_instances": flattenCocInstanceBatchesTargetInstances(
					utils.PathSearch("target_instances", raw, nil)),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocInstanceBatchesTargetInstances(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"resource_id":        utils.PathSearch("resource_id", raw, nil),
				"cloud_service_name": utils.PathSearch("provider", raw, nil),
				"region_id":          utils.PathSearch("region_id", raw, nil),
				"type":               utils.PathSearch("type", raw, nil),
				"custom_attributes": flattenCocInstanceBatchesTargetInstancesCustomAttributes(
					utils.PathSearch("custom_attributes", raw, nil)),
				"properties": flattenCocInstanceBatchesTargetInstancesProperties(
					utils.PathSearch("properties", raw, nil)),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocInstanceBatchesTargetInstancesCustomAttributes(rawParams interface{}) []interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) == 0 {
			return nil
		}
		rst := make([]interface{}, 0, len(paramsList))
		for _, params := range paramsList {
			raw := params.(map[string]interface{})
			m := map[string]interface{}{
				"key":   utils.PathSearch("key", raw, nil),
				"value": utils.PathSearch("value", raw, nil),
			}
			rst = append(rst, m)
		}

		return rst
	}

	return nil
}

func flattenCocInstanceBatchesTargetInstancesProperties(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"host_name": utils.PathSearch("host_name", param, nil),
			"fixed_ip":  utils.PathSearch("fixed_ip", param, nil),
			"region_id": utils.PathSearch("region_id", param, nil),
			"zone_id":   utils.PathSearch("zone_id", param, nil),
		},
	}

	return rst
}

package ces

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

// @API CES POST /V1.0/{project_id}/batch-query-metric-data
func DataSourceMultipleMetricsData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultipleMetricsDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `Specifies the metric data.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the namespace of a service.`,
						},
						"metric_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `Specifies the metric ID.`,
						},
						"dimensions": {
							Type:        schema.TypeList,
							Required:    true,
							Description: `Specifies metric dimensions.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the dimension.`,
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: `Specifies the dimension value.`,
									},
								},
							},
						},
					},
				},
			},
			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time of the query.`,
			},
			"to": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time of the query.`,
			},
			"period": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies how often Cloud Eye aggregates data.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the data rollup method. The field does not affect the query result of raw data. (The period is **1**.)`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The metric data.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The namespace of a service.`,
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The metric ID.`,
						},
						"dimensions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The metric dimensions.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dimension.`,
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The dimension value.`,
									},
								},
							},
						},
						"unit": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The metric unit.`,
						},
						"datapoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The metric data list.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"average": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The average value of metric data within a rollup period.`,
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The maximum value of metric data within a rollup period.`,
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The minimum value of metric data within a rollup period.`,
									},
									"sum": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The sum of metric data within a rollup period.`,
									},
									"variance": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The variance of metric data within a rollup period.`,
									},
									"timestamp": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The time when the metric is collected. The time is a UNIX timestamp and the unit is ms.`,
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

func dataSourceMultipleMetricsDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	metricsData, err := getMultipleMetricsData(client, d)
	if err != nil {
		return diag.Errorf("error retrieving CES metrics data: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data", flattenMultipleMetricsData(metricsData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getMultipleMetricsData(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	httpUrl := "V1.0/{project_id}/batch-query-metric-data"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	params, err := buildMultipleMetricsDataBodyParams(d)
	if err != nil {
		return nil, fmt.Errorf("error building multiple metrics data body params: %s", err)
	}
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         params,
	}

	resp, err := client.Request("POST", path, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	metricsData := utils.PathSearch("metrics", respBody, make([]interface{}, 0)).([]interface{})
	return metricsData, nil
}

func buildMultipleMetricsDataBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	fromTime, err := utils.FormatUTCTimeStamp(d.Get("from").(string))
	if err != nil {
		return nil, err
	}

	toTime, err := utils.FormatUTCTimeStamp(d.Get("to").(string))
	if err != nil {
		return nil, err
	}

	param := map[string]interface{}{
		"metrics": buildMultipleMetricsDataMetricsBodyParams(d.Get("metrics").([]interface{})),
		"from":    fromTime * 1000,
		"to":      toTime * 1000,
		"period":  d.Get("period").(string),
		"filter":  d.Get("filter").(string),
	}

	return param, nil
}

func buildMultipleMetricsDataMetricsBodyParams(metrics []interface{}) []map[string]interface{} {
	if len(metrics) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(metrics))
	for i, metric := range metrics {
		m := metric.(map[string]interface{})
		result[i] = map[string]interface{}{
			"namespace":   m["namespace"],
			"metric_name": m["metric_name"],
			"dimensions":  buildMultipleMetricsDataDimensionsBodyParams(m["dimensions"].([]interface{})),
		}
	}
	return result
}

func buildMultipleMetricsDataDimensionsBodyParams(dimensions []interface{}) []map[string]interface{} {
	if len(dimensions) == 0 {
		return nil
	}
	result := make([]map[string]interface{}, len(dimensions))
	for i, dimension := range dimensions {
		d := dimension.(map[string]interface{})
		result[i] = map[string]interface{}{
			"name":  d["name"],
			"value": d["value"],
		}
	}
	return result
}

func flattenMultipleMetricsData(metricsData []interface{}) []interface{} {
	if len(metricsData) == 0 {
		return nil
	}

	result := make([]interface{}, len(metricsData))
	for i, data := range metricsData {
		result[i] = map[string]interface{}{
			"namespace":   utils.PathSearch("namespace", data, nil),
			"metric_name": utils.PathSearch("metric_name", data, nil),
			"dimensions":  flattenMultipleMetricsDataDimensions(utils.PathSearch("dimensions", data, make([]interface{}, 0)).([]interface{})),
			"datapoints":  flattenMultipleMetricsDataPoints(utils.PathSearch("datapoints", data, make([]interface{}, 0)).([]interface{})),
			"unit":        utils.PathSearch("unit", data, nil),
		}
	}
	return result
}

func flattenMultipleMetricsDataDimensions(dimensions []interface{}) []interface{} {
	if len(dimensions) == 0 {
		return nil
	}

	result := make([]interface{}, len(dimensions))
	for i, dimension := range dimensions {
		result[i] = map[string]interface{}{
			"name":  utils.PathSearch("name", dimension, nil),
			"value": utils.PathSearch("value", dimension, nil),
		}
	}
	return result
}

func flattenMultipleMetricsDataPoints(dataPoints []interface{}) []interface{} {
	if len(dataPoints) == 0 {
		return nil
	}

	result := make([]interface{}, len(dataPoints))
	for i, dataPoint := range dataPoints {
		result[i] = map[string]interface{}{
			"average":   utils.PathSearch("average", dataPoint, nil),
			"sum":       utils.PathSearch("sum", dataPoint, nil),
			"variance":  utils.PathSearch("variance", dataPoint, nil),
			"max":       utils.PathSearch("max", dataPoint, nil),
			"min":       utils.PathSearch("min", dataPoint, nil),
			"timestamp": utils.PathSearch("timestamp", dataPoint, nil),
		}
	}
	return result
}

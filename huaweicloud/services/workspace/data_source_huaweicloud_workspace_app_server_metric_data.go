package workspace

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

// @API Workspace GET /v1/{project_id}/app-servers/server-metric-data/{server_id}
func DataSourceAppServerMetricData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppServerMetricDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the server metric data are located.`,
			},
			"server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the server.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The namespace of the service.`,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the monitoring metric.`,
			},
			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the query data, in RFC3339 format.`,
			},
			"to": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the query data, in RFC3339 format.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The granularity of monitoring data.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data aggregation method.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of server metric data.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the monitoring metric.`,
						},
						"dimension_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dimension value.`,
						},
						"datapoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of metric data points.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"average": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The average value of the metric data within the aggregation period.`,
									},
									"max": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The maximum value of the metric data within the aggregation period.`,
									},
									"min": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The minimum value of the metric data within the aggregation period.`,
									},
									"sum": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The sum value of the metric data within the aggregation period.`,
									},
									"variance": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The variance of the metric data within the aggregation period.`,
									},
									"collection_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The collection time of the metric, in RFC3339 format.`,
									},
									"unit": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The unit of the metric.`,
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

func buildAppServerMetricDataQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?namespace=%v&metric_name=%v&from=%v&to=%v&period=%v&filter=%v",
		d.Get("namespace"),
		d.Get("metric_name"),
		utils.ConvertTimeStrToNanoTimestamp(d.Get("from").(string)),
		utils.ConvertTimeStrToNanoTimestamp(d.Get("to").(string)),
		d.Get("period"),
		d.Get("filter"),
	)
}

func listAppServerMetricData(client *golangsdk.ServiceClient, d *schema.ResourceData, serverId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/app-servers/server-metric-data/{server_id}"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{server_id}", serverId)
	listPath += buildAppServerMetricDataQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func dataSourceAppServerMetricDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		serverId = d.Get("server_id").(string)
	)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	respBody, err := listAppServerMetricData(client, d, serverId)
	if err != nil {
		return diag.Errorf("error querying APP server (%s) metric data: %s", serverId, err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("metrics", flattenAppServerMetricDataMetrics(utils.PathSearch("server_metrics",
			respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAppServerMetricDataMetrics(metrics []interface{}) []map[string]interface{} {
	if len(metrics) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		result = append(result, map[string]interface{}{
			"metric_name":     utils.PathSearch("metric_name", metric, nil),
			"dimension_value": utils.PathSearch("dimension_value", metric, nil),
			"datapoints": flattenAppServerMetricDataMetricDatapoints(utils.PathSearch("datapoints",
				metric, make([]interface{}, 0)).([]interface{})),
		})
	}

	return result
}

func flattenAppServerMetricDataMetricDatapoints(datapoints []interface{}) []map[string]interface{} {
	if len(datapoints) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(datapoints))
	for _, v := range datapoints {
		result = append(result, map[string]interface{}{
			"average":  utils.PathSearch("average", v, nil),
			"max":      utils.PathSearch("max", v, nil),
			"min":      utils.PathSearch("min", v, nil),
			"sum":      utils.PathSearch("sum", v, nil),
			"variance": utils.PathSearch("variance", v, nil),
			"collection_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("timestamp",
				v, float64(0)).(float64))/1000, false),
			"unit": utils.PathSearch("unit", v, nil),
		})
	}

	return result
}

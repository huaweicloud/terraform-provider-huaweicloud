package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/metric-data
func DataSourceInstanceMetricData() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstanceMetricDataRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the dedicated instance is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the dedicated instance.`,
			},
			"dim": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The dimension of the metric data.`,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the metric data.`,
			},
			"from": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time of the metric data, UNIX timestamp in milliseconds.`,
			},
			"to": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time of the metric data, UNIX timestamp in milliseconds.`,
			},
			"period": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The granularity of the metric data.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data aggregation method of the metric data.`,
			},

			// Attributes.
			"datapoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the metric data points that matched the filter parameters.`,
				Elem:        dataPointsSchema(),
			},
		},
	}
}

func dataPointsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"average": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The average value of the metric data within the aggregation period.`,
			},
			"max": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum value of the metric data within the aggregation period.`,
			},
			"min": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum value of the metric data within the aggregation period.`,
			},
			"sum": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The sum value of the metric data within the aggregation period.`,
			},
			"variance": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The variance value of the metric data within the aggregation period.`,
			},
			"timestamp": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The collection time of the metric data, UNIX timestamp in milliseconds.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit of the metric data.`,
			},
		},
	}
}

func flattenInstanceMetricDatapoints(datapoints []interface{}) []interface{} {
	if len(datapoints) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(datapoints))
	for _, datapoint := range datapoints {
		result = append(result, map[string]interface{}{
			"average":   utils.PathSearch("average", datapoint, nil),
			"max":       utils.PathSearch("max", datapoint, nil),
			"min":       utils.PathSearch("min", datapoint, nil),
			"sum":       utils.PathSearch("sum", datapoint, nil),
			"variance":  utils.PathSearch("variance", datapoint, nil),
			"timestamp": utils.PathSearch("timestamp", datapoint, nil),
			"unit":      utils.PathSearch("unit", datapoint, nil),
		})
	}

	return result
}

func buildInstanceMetricDataQueryParams(d *schema.ResourceData) string {
	res := ""

	res = fmt.Sprintf("%s&dim=%v", res, d.Get("dim"))
	res = fmt.Sprintf("%s&metric_name=%v", res, d.Get("metric_name"))
	res = fmt.Sprintf("%s&from=%v", res, d.Get("from"))
	res = fmt.Sprintf("%s&to=%v", res, d.Get("to"))
	res = fmt.Sprintf("%s&period=%v", res, d.Get("period"))
	res = fmt.Sprintf("%s&filter=%v", res, d.Get("filter"))

	return "?" + res[1:]
}

func dataSourceInstanceMetricDataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/metric-data"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath += buildInstanceMetricDataQueryParams(d)

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return diag.Errorf("error querying instance metric data: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("datapoints", flattenInstanceMetricDatapoints(
			utils.PathSearch("datapoints", respBody, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

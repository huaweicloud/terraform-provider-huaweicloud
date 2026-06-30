package secmaster

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/metrics
func DataSourceMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     metricsDataSchema(),
			},
		},
	}
}

func metricsDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_dimension": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cache_ttl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"report_period": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"effective_column": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"max_query_range": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"derived_metrics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     metricsDerivedMetricsSchema(),
			},
			"compound_expression": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metric_format": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     metricsFormatSchema(),
			},
			"metric_expand_dim": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     metricsExpandDimSchema(),
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func metricsDerivedMetricsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"metric_dimension": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_query_range": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"date_start": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_end": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"date_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"query_function": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func metricsFormatSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_param": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"data_param": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func metricsExpandDimSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"functions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/metrics"
		product = "secmaster"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving metrics: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenMetrics(utils.PathSearch(
			"data", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMetrics(metrics []interface{}) []interface{} {
	if len(metrics) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(metrics))
	for _, v := range metrics {
		rst = append(rst, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"id":               utils.PathSearch("id", v, nil),
			"metric_type":      utils.PathSearch("metric_type", v, nil),
			"data_type":        utils.PathSearch("data_type", v, nil),
			"metric_dimension": utils.PathSearch("metric_dimension", v, nil),
			"cache_ttl":        utils.PathSearch("cache_ttl", v, nil),
			"report_period":    utils.PathSearch("report_period", v, nil),
			"is_built_in":      utils.PathSearch("is_built_in", v, nil),
			"effective_column": utils.PathSearch("effective_column", v, nil),
			"max_query_range":  utils.PathSearch("max_query_range", v, nil),
			"derived_metrics": flattenMetricsDerivedMetrics(
				utils.PathSearch("derived_metrics", v, make([]interface{}, 0)).([]interface{})),
			"compound_expression": utils.PathSearch("compound_expression", v, nil),
			"metric_format": flattenMetricsFormat(
				utils.PathSearch("metric_format", v, make([]interface{}, 0)).([]interface{})),
			"metric_expand_dim": flattenMetricsExpandDim(
				utils.PathSearch("metric_expand_dim", v, nil)),
			"version": utils.PathSearch("version", v, nil),
		})
	}

	return rst
}

func flattenMetricsDerivedMetrics(derivedMetrics []interface{}) []interface{} {
	if len(derivedMetrics) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(derivedMetrics))
	for _, v := range derivedMetrics {
		rst = append(rst, map[string]interface{}{
			"metric_dimension": utils.PathSearch("metric_dimension", v, nil),
			"max_query_range":  utils.PathSearch("max_query_range", v, nil),
			"date_start":       utils.PathSearch("date_start", v, nil),
			"date_end":         utils.PathSearch("date_end", v, nil),
			"date_format":      utils.PathSearch("date_format", v, nil),
			"query_type":       utils.PathSearch("query_type", v, nil),
			"query_function":   utils.PathSearch("query_function", v, nil),
		})
	}

	return rst
}

func flattenMetricsFormat(metricFormats []interface{}) []interface{} {
	if len(metricFormats) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(metricFormats))
	for _, v := range metricFormats {
		rst = append(rst, map[string]interface{}{
			"data":    utils.PathSearch("data", v, nil),
			"display": utils.PathSearch("display", v, nil),
			"display_param": utils.ExpandToStringMap(
				utils.PathSearch("display_param", v, make(map[string]interface{})).(map[string]interface{})),
			"data_param": utils.ExpandToStringMap(
				utils.PathSearch("data_param", v, make(map[string]interface{})).(map[string]interface{})),
		})
	}

	return rst
}

func flattenMetricsExpandDim(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"labels": utils.ExpandToStringList(utils.PathSearch(
				"labels", raw, make([]interface{}, 0)).([]interface{})),
			"functions": utils.ExpandToStringList(utils.PathSearch(
				"functions", raw, make([]interface{}, 0)).([]interface{})),
		},
	}
}

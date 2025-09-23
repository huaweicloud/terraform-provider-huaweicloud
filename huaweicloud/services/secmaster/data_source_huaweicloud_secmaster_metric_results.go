package secmaster

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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/sa/metrics/hits
func DataSourceMetricResults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMetricResultsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the data source. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID.`,
			},
			"metric_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the metrics IDs.`,
			},
			"timespan": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time range for querying metrics.`,
			},
			"cache": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether the cache is enabled.`,
			},
			"field_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the indicator card IDs.`,
			},
			"params": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the parameter list of the metric.`,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"interactive_params": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `Specifies the interactive parameters.`,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{Type: schema.TypeString},
				},
			},
			"metric_results": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The metric results.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The metric ID.`,
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The statistical labels of the metric.`,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"data_rows": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `All statistical data of the metric.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"data_row": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
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

func dataSourceMetricResultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	// Query the list of SecMaster metrics.
	resp, err := getMetricResults(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("metric_results", flattenMetricResultsResponse(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getMetricResults(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	listMetricResultsHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/sa/metrics/hits"
	listMetricResultsPath := client.Endpoint + listMetricResultsHttpUrl
	listMetricResultsPath = strings.ReplaceAll(listMetricResultsPath, "{project_id}", client.ProjectID)
	listMetricResultsPath = strings.ReplaceAll(listMetricResultsPath, "{workspace_id}", d.Get("workspace_id").(string))
	listMetricResultsPath += buildListMetricResultsQueryParams(d)

	listMetricResultsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listMetricResultsOpt.JSONBody = utils.RemoveNil(buildListMetricResultsBodyParams(d))

	listMetricResultsResp, err := client.Request("POST", listMetricResultsPath, &listMetricResultsOpt)
	if err != nil {
		return nil, fmt.Errorf("error querying SecMaster metric results: %s", err)
	}
	listMetricResultsRespBody, err := utils.FlattenResponse(listMetricResultsResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening SecMaster metric results: %s", err)
	}
	metricResults := listMetricResultsRespBody.([]interface{})

	return metricResults, nil
}

func buildListMetricResultsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("timespan"); ok {
		res = fmt.Sprintf("%s&timespan=%v", res, v)
	}

	if v, ok := d.GetOk("cache"); ok {
		res = fmt.Sprintf("%s&cache=%v", res, v)
	}

	if res != "" {
		return fmt.Sprintf("?%s", res[1:])
	}

	return ""
}

func buildListMetricResultsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"metric_ids":         d.Get("metric_ids"),
		"field_ids":          utils.ValueIgnoreEmpty(d.Get("field_ids")),
		"params":             utils.ValueIgnoreEmpty(d.Get("params")),
		"interactive_params": utils.ValueIgnoreEmpty(d.Get("interactive_params")),
	}

	return bodyParam
}

func flattenMetricResultsResponse(rawParams []interface{}) []interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	metrics := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		metrics[i] = map[string]interface{}{
			"id":        utils.PathSearch("metric_id", v, nil),
			"labels":    utils.PathSearch("result.labels", v, nil),
			"data_rows": flattenResultResponse(utils.PathSearch("result", v, nil)),
		}
	}
	return metrics
}

func flattenResultResponse(params interface{}) []interface{} {
	if params == nil {
		return nil
	}

	dataRows := utils.PathSearch("datarows", params, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(dataRows))
	for i, v := range dataRows {
		result[i] = map[string]interface{}{
			"data_row": expandToStringList(v.([]interface{})),
		}
	}
	return result
}

// Convert all elements in the array to character transfer type.
func expandToStringList(v []interface{}) []string {
	s := make([]string, len(v))
	for i, val := range v {
		s[i] = fmt.Sprintf("%v", val)
	}

	return s
}

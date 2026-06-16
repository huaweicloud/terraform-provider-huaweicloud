package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v2/{project_id}/design/biz-metrics
func DataSourceArchitectureBusinessMetrics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureBusinessMetricsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the business metrics are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the business metrics belong.`,
			},

			// Optional parameters.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name or code of the business metric to be fuzzy queried.`,
			},
			"create_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The creator of the business metric to be queried.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The owner of the business metric to be queried.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The publishing status of the business metric to be queried.`,
			},
			"biz_catalog_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The process architecture ID to which the business metric belongs.`,
			},

			// Attributes.
			"metrics": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureBusinessMetricsElem(),
				Description: `The list of business metrics that matched filter parameters.`,
			},
		},
	}
}

func dataArchitectureBusinessMetricsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the business metric.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the business metric.`,
			},
			"biz_catalog_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process architecture ID.`,
			},
			"time_filters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The statistical frequency.`,
			},
			"interval_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The refresh frequency.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of person responsible for the indicator.`,
			},
			"owner_department": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The indicator management department name.`,
			},
			"destination": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The purpose of setting.`,
			},
			"definition": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The indicator definition.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The calculation formula.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description.`,
			},
			"apply_scenario": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The application scenarios.`,
			},
			"technical_metric": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The related technical indicators.`,
			},
			"measure": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The measurement object.`,
			},
			"dimensions": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The statistical dimension.`,
			},
			"general_filters": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The statistical caliber and modifiers.`,
			},
			"data_origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data sources.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit of measurement.`,
			},
			"name_alias": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The indicator alias.`,
			},
			"code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The indicator encoding.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the business metric.`,
			},
			"biz_catalog_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process architecture path.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator of the business metric.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The editor of the business metric.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the metric, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the metric, in RFC3339 format.`,
			},
			"technical_metric_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The related technical indicator type.`,
			},
			"technical_metric_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The related technical indicator name.`,
			},
			"l1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject domain grouping Chinese name.`,
			},
			"l2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subject field Chinese name.`,
			},
			"l3": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business object Chinese name.`,
			},
			"biz_metric": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The business indicator synchronization status.`,
			},
			"summary_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The synchronize statistics status.`,
			},
		},
	}
}

func buildArchitectureBusinessMetricsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&name=%v", res, v)
	}
	if v, ok := d.GetOk("create_by"); ok {
		res = fmt.Sprintf("%s&create_by=%v", res, v)
	}
	if v, ok := d.GetOk("owner"); ok {
		res = fmt.Sprintf("%s&owner=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("biz_catalog_id"); ok {
		res = fmt.Sprintf("%s&biz_catalog_id=%v", res, v)
	}

	return res
}

func listArchitectureBusinessMetrics(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/design/biz-metrics?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath += buildArchitectureBusinessMetricsQueryParams(d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		metrics := utils.PathSearch("data.value.records", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, metrics...)

		if len(metrics) < limit {
			break
		}
		offset += len(metrics)
	}

	return result, nil
}

func flattenArchitectureBusinessMetrics(metrics []interface{}) []map[string]interface{} {
	if len(metrics) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metrics))
	for _, metric := range metrics {
		technicalMetricResp := utils.PathSearch("technical_metric", metric, "").(string)
		technicalMetric := utils.StringToInt(&technicalMetricResp)

		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", metric, nil),
			"name":                  utils.PathSearch("name", metric, nil),
			"biz_catalog_id":        utils.PathSearch("biz_catalog_id", metric, nil),
			"time_filters":          utils.PathSearch("time_filters", metric, nil),
			"interval_type":         utils.PathSearch("interval_type", metric, nil),
			"owner":                 utils.PathSearch("owner", metric, nil),
			"owner_department":      utils.PathSearch("owner_department", metric, nil),
			"destination":           utils.PathSearch("destination", metric, nil),
			"definition":            utils.PathSearch("definition", metric, nil),
			"expression":            utils.PathSearch("expression", metric, nil),
			"description":           utils.PathSearch("remark", metric, nil),
			"apply_scenario":        utils.PathSearch("apply_scenario", metric, nil),
			"technical_metric":      technicalMetric,
			"measure":               utils.PathSearch("measure", metric, nil),
			"dimensions":            utils.PathSearch("dimensions", metric, nil),
			"general_filters":       utils.PathSearch("general_filters", metric, nil),
			"data_origin":           utils.PathSearch("data_origin", metric, nil),
			"unit":                  utils.PathSearch("unit", metric, nil),
			"name_alias":            utils.PathSearch("name_alias", metric, nil),
			"code":                  utils.PathSearch("code", metric, nil),
			"status":                utils.PathSearch("status", metric, nil),
			"biz_catalog_path":      utils.PathSearch("biz_catalog_path", metric, nil),
			"created_by":            utils.PathSearch("create_by", metric, nil),
			"updated_by":            utils.PathSearch("update_by", metric, nil),
			"created_at":            utils.PathSearch("create_time", metric, nil),
			"updated_at":            utils.PathSearch("update_time", metric, nil),
			"technical_metric_type": utils.PathSearch("technical_metric_type", metric, nil),
			"technical_metric_name": utils.PathSearch("technical_metric_name", metric, nil),
			"l1":                    utils.PathSearch("l1", metric, nil),
			"l2":                    utils.PathSearch("l2", metric, nil),
			"l3":                    utils.PathSearch("l3", metric, nil),
			"biz_metric":            utils.PathSearch("biz_metric", metric, nil),
			"summary_status":        utils.PathSearch("summary_status", metric, nil),
		})
	}
	return result
}

func dataSourceArchitectureBusinessMetricsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	metrics, err := listArchitectureBusinessMetrics(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture business metrics: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("metrics", flattenArchitectureBusinessMetrics(metrics)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

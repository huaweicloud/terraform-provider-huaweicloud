package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var nonUpdatableParamsMetric = []string{
	"workspace_id",
	"metric_type",
	"data_type",
	"is_built_in",
	"version",
}

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/metrics
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}
func ResourceMetric() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetricCreate,
		ReadContext:   resourceMetricRead,
		UpdateContext: resourceMetricUpdate,
		DeleteContext: resourceMetricDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceMetricImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParamsMetric),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metric_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cache_ttl": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric_dimension": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"report_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"effective_column": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_query_range": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"derived_metrics": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"metric_dimension": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"max_query_range": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"date_start": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"date_end": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"date_format": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"query_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"query_function": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"compound_expression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"metric_format": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"data": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"display": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"display_param": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"data_param": {
							Type:     schema.TypeMap,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"metric_expand_dim": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"labels": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"functions": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDerivedMetricsBodyParams(rawList []interface{}) []interface{} {
	if len(rawList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawList))
	for _, v := range rawList {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"metric_dimension": utils.ValueIgnoreEmpty(item["metric_dimension"]),
			"max_query_range":  utils.ValueIgnoreEmpty(item["max_query_range"]),
			"date_start":       utils.ValueIgnoreEmpty(item["date_start"]),
			"date_end":         utils.ValueIgnoreEmpty(item["date_end"]),
			"date_format":      utils.ValueIgnoreEmpty(item["date_format"]),
			"query_type":       utils.ValueIgnoreEmpty(item["query_type"]),
			"query_function":   utils.ValueIgnoreEmpty(item["query_function"]),
		})
	}

	return result
}

func buildMetricFormatBodyParams(rawList []interface{}) []interface{} {
	if len(rawList) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rawList))
	for _, v := range rawList {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"data":          utils.ValueIgnoreEmpty(item["data"]),
			"display":       utils.ValueIgnoreEmpty(item["display"]),
			"display_param": utils.ValueIgnoreEmpty(item["display_param"]),
			"data_param":    utils.ValueIgnoreEmpty(item["data_param"]),
		})
	}

	return result
}

func buildMetricExpandDimBodyParams(rawList []interface{}) interface{} {
	if len(rawList) == 0 {
		return nil
	}

	item, ok := rawList[0].(map[string]interface{})
	if !ok {
		return nil
	}

	return map[string]interface{}{
		"labels":    utils.ValueIgnoreEmpty(item["labels"]),
		"functions": utils.ValueIgnoreEmpty(item["functions"]),
	}
}

func buildCreateMetricBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"metric_type": d.Get("metric_type"),
		"data_type":   d.Get("data_type"),
		"cache_ttl":   d.Get("cache_ttl"),
		"is_built_in": d.Get("is_built_in"),
	}

	if v, ok := d.GetOk("metric_dimension"); ok {
		bodyParams["metric_dimension"] = v
	}
	if v, ok := d.GetOk("report_period"); ok {
		bodyParams["report_period"] = v
	}
	if v, ok := d.GetOk("effective_column"); ok {
		bodyParams["effective_column"] = v
	}
	if v, ok := d.GetOk("max_query_range"); ok {
		bodyParams["max_query_range"] = v
	}
	if v, ok := d.GetOk("derived_metrics"); ok {
		bodyParams["derived_metrics"] = buildDerivedMetricsBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("compound_expression"); ok {
		bodyParams["compound_expression"] = v
	}
	if v, ok := d.GetOk("metric_format"); ok {
		bodyParams["metric_format"] = buildMetricFormatBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("metric_expand_dim"); ok {
		bodyParams["metric_expand_dim"] = buildMetricExpandDimBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("version"); ok {
		bodyParams["version"] = v
	}

	return bodyParams
}

func resourceMetricCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		workspaceId   = d.Get("workspace_id").(string)
		createHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/metrics"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{workspace_id}", workspaceId)

	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateMetricBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster metric: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	metricId := utils.PathSearch("id", respBody, "").(string)
	if metricId == "" {
		return diag.Errorf("error creating SecMaster metric: unable to find metric ID")
	}

	d.SetId(metricId)

	return resourceMetricRead(ctx, d, meta)
}

func GetMetricInfo(client *golangsdk.ServiceClient, workspaceId, metricId string) (interface{}, error) {
	getPath := client.Endpoint + "v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{workspace_id}", workspaceId)
	getPath = strings.ReplaceAll(getPath, "{metric_id}", metricId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func flattenDerivedMetrics(respBody interface{}) []map[string]interface{} {
	derivedMetricsRaw := utils.PathSearch("derived_metrics", respBody, make([]interface{}, 0))
	derivedMetricsList, ok := derivedMetricsRaw.([]interface{})
	if !ok || len(derivedMetricsList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(derivedMetricsList))
	for _, v := range derivedMetricsList {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"metric_dimension": utils.PathSearch("metric_dimension", item, nil),
			"max_query_range":  utils.PathSearch("max_query_range", item, nil),
			"date_start":       utils.PathSearch("date_start", item, nil),
			"date_end":         utils.PathSearch("date_end", item, nil),
			"date_format":      utils.PathSearch("date_format", item, nil),
			"query_type":       utils.PathSearch("query_type", item, nil),
			"query_function":   utils.PathSearch("query_function", item, nil),
		})
	}

	return result
}

func flattenMetricFormat(respBody interface{}) []map[string]interface{} {
	metricFormatRaw := utils.PathSearch("metric_format", respBody, make([]interface{}, 0))
	metricFormatList, ok := metricFormatRaw.([]interface{})
	if !ok || len(metricFormatList) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(metricFormatList))
	for _, v := range metricFormatList {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		result = append(result, map[string]interface{}{
			"data":          utils.PathSearch("data", item, nil),
			"display":       utils.PathSearch("display", item, nil),
			"display_param": utils.PathSearch("display_param", item, nil),
			"data_param":    utils.PathSearch("data_param", item, nil),
		})
	}

	return result
}

func flattenMetricExpandDim(respBody interface{}) []map[string]interface{} {
	metricExpandDimRaw := utils.PathSearch("metric_expand_dim", respBody, nil)
	if metricExpandDimRaw == nil {
		return nil
	}

	labelsRaw := utils.PathSearch("labels", metricExpandDimRaw, make([]interface{}, 0))
	functionsRaw := utils.PathSearch("functions", metricExpandDimRaw, make([]interface{}, 0))

	labels, _ := labelsRaw.([]interface{})
	functions, _ := functionsRaw.([]interface{})

	if len(labels) == 0 && len(functions) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"labels":    labels,
			"functions": functions,
		},
	}
}

func resourceMetricRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	respBody, err := GetMetricInfo(client, workspaceId, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SecMaster.20097017"),
			"error retrieving SecMaster metric",
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("metric_type", utils.PathSearch("metric_type", respBody, nil)),
		d.Set("data_type", utils.PathSearch("data_type", respBody, nil)),
		d.Set("cache_ttl", utils.PathSearch("cache_ttl", respBody, nil)),
		d.Set("metric_dimension", utils.PathSearch("metric_dimension", respBody, nil)),
		d.Set("report_period", utils.PathSearch("report_period", respBody, nil)),
		d.Set("is_built_in", utils.PathSearch("is_built_in", respBody, nil)),
		d.Set("effective_column", utils.PathSearch("effective_column", respBody, nil)),
		d.Set("max_query_range", utils.PathSearch("max_query_range", respBody, nil)),
		d.Set("derived_metrics", flattenDerivedMetrics(respBody)),
		d.Set("compound_expression", utils.PathSearch("compound_expression", respBody, nil)),
		d.Set("metric_format", flattenMetricFormat(respBody)),
		d.Set("metric_expand_dim", flattenMetricExpandDim(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateMetricBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"metric_type": d.Get("metric_type"),
		"data_type":   d.Get("data_type"),
		"cache_ttl":   d.Get("cache_ttl"),
		"is_built_in": d.Get("is_built_in"),
	}

	if v, ok := d.GetOk("metric_dimension"); ok {
		bodyParams["metric_dimension"] = v
	}
	if v, ok := d.GetOk("report_period"); ok {
		bodyParams["report_period"] = v
	}
	if v, ok := d.GetOk("effective_column"); ok {
		bodyParams["effective_column"] = v
	}
	if v, ok := d.GetOk("max_query_range"); ok {
		bodyParams["max_query_range"] = v
	}
	if v, ok := d.GetOk("derived_metrics"); ok {
		bodyParams["derived_metrics"] = buildDerivedMetricsBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("compound_expression"); ok {
		bodyParams["compound_expression"] = v
	}
	if v, ok := d.GetOk("metric_format"); ok {
		bodyParams["metric_format"] = buildMetricFormatBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("metric_expand_dim"); ok {
		bodyParams["metric_expand_dim"] = buildMetricExpandDimBodyParams(v.([]interface{}))
	}
	if v, ok := d.GetOk("version"); ok {
		bodyParams["version"] = v
	}

	return bodyParams
}

func resourceMetricUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{workspace_id}", workspaceId)
	updatePath = strings.ReplaceAll(updatePath, "{metric_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateMetricBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating SecMaster metric: %s", err)
	}

	return resourceMetricRead(ctx, d, meta)
}

func resourceMetricDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/soc/metrics/{metric_id}"
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{workspace_id}", workspaceId)
	deletePath = strings.ReplaceAll(deletePath, "{metric_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"content-type": "application/json;charset=UTF-8"},
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)

	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "SecMaster.20097017"),
			"error deleting SecMaster metric",
		)
	}

	return nil
}

func resourceMetricImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<metric_id>', but got '%s'",
			importedId)
	}

	d.SetId(parts[1])

	mErr := multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

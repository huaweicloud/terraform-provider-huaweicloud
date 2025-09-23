package ces

import (
	"context"
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

var widgetNonUpdatableParams = []string{"dashboard_id", "view"}

// @API CES POST /v2/{project_id}/dashboards/{dashboard_id}/widgets
// @API CES GET /v2/{project_id}/widgets/{widget_id}
// @API CES POST /v2/{project_id}/widgets/batch-update
// @API CES DELETE /v2/{project_id}/widgets/{widget_id}
func ResourceDashboardWidget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDashboardWidgetCreate,
		UpdateContext: resourceDashboardWidgetUpdate,
		ReadContext:   resourceDashboardWidgetRead,
		DeleteContext: resourceDashboardWidgetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(widgetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dashboard ID.`,
			},
			"title": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dashboard widget title.`,
			},
			"metrics": {
				Type:        schema.TypeList,
				Elem:        widgetMetricOptsSchema(),
				Required:    true,
				Description: `Specifies the metric list.`,
			},
			"view": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the graph type.`,
			},
			"metric_display_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies how many metrics will be displayed on one widget.`,
			},
			"location": {
				Type:        schema.TypeList,
				Elem:        widgetLocationOptsSchema(),
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the dashboard widget coordinates.`,
			},
			"properties": {
				Type:        schema.TypeList,
				Elem:        widgetPropertiesOptsSchema(),
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: `Specifies additional information`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the metric unit.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `When the dashboard widget was created.`,
			},
		},
	}
}

func widgetMetricOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the cloud service dimension.`,
			},
			"dimensions": {
				Type:        schema.TypeList,
				Elem:        widgetDimensionOptsSchema(),
				Required:    true,
				MaxItems:    1,
				Description: `Specifies the dimension list.`,
			},
			"metric_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the metric name.`,
			},
			"alias": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the alias list of metrics.`,
			},
		},
	}
	return &sc
}

func widgetLocationOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"top": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the grids between the widget and the top of the dashboard.`,
			},
			"left": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the grids between the widget and the left side of the dashboard.`,
			},
			"width": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the dashboard widget width.`,
			},
			"height": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the dashboard widget height.`,
			},
		},
	}
	return &sc
}

func widgetPropertiesOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"top_n": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the top n resources sorted by a metric.`,
			},
			"filter": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies how metric data is aggregated.`,
			},
			"order": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies how top n resources by a metric are sorted on a dashboard widget.`,
			},
		},
	}
	return &sc
}

func widgetDimensionOptsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the dimension name.`,
			},
			"filter_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource type.`,
			},
			"values": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the dimension value list.`,
			},
		},
	}
	return &sc
}

func resourceDashboardWidgetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createWidgetHttpUrl = "v2/{project_id}/dashboards/{dashboard_id}/widgets"
		createWidgetProduct = "ces"
	)
	createWidgetClient, err := cfg.NewServiceClient(createWidgetProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	createWidgetPath := createWidgetClient.Endpoint + createWidgetHttpUrl
	createWidgetPath = strings.ReplaceAll(createWidgetPath, "{project_id}", createWidgetClient.ProjectID)
	createWidgetPath = strings.ReplaceAll(createWidgetPath, "{dashboard_id}", d.Get("dashboard_id").(string))

	createWidgetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	params := utils.RemoveNil(buildCreateWidgetBodyParams(d))
	createWidgetOpt.JSONBody = []map[string]interface{}{params}
	createWidgetResp, err := createWidgetClient.Request("POST", createWidgetPath, &createWidgetOpt)
	if err != nil {
		return diag.Errorf("error creating CES dashboard widget: %s", err)
	}

	createWidgetRespBody, err := utils.FlattenResponse(createWidgetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("widget_ids|[0]", createWidgetRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CES dashboard widget: ID is not found in API response")
	}
	d.SetId(id)

	if val, ok := d.GetOk("unit"); ok {
		updateWidgetPath := createWidgetClient.Endpoint + "v2/{project_id}/widgets/batch-update"
		updateWidgetPath = strings.ReplaceAll(updateWidgetPath, "{project_id}", createWidgetClient.ProjectID)
		updateWidgetOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		}
		param := map[string]interface{}{
			"widget_id": d.Id(),
			"unit":      val,
		}
		updateWidgetOpt.JSONBody = []interface{}{param}
		_, err := createWidgetClient.Request("POST", updateWidgetPath, &updateWidgetOpt)
		if err != nil {
			return diag.Errorf("error updating CES dashboard widget: %s", err)
		}
	}

	return resourceDashboardWidgetRead(ctx, d, meta)
}

func buildCreateWidgetBodyParams(d *schema.ResourceData) map[string]interface{} {
	// Currently, the threshold related functions are invalid.
	return map[string]interface{}{
		"title":               d.Get("title"),
		"metrics":             buildWidgetMetrics(d.Get("metrics")),
		"view":                d.Get("view"),
		"metric_display_mode": d.Get("metric_display_mode"),
		"location":            buildWidgetLocation(d.Get("location")),
		"properties":          buildWidgetProperties(d.Get("properties")),
		"threshold_enabled":   false,
	}
}

func buildWidgetMetrics(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		params := make([]map[string]interface{}, 0, len(rawArray))
		for _, rawMetric := range rawArray {
			metric := rawMetric.(map[string]interface{})
			params = append(params, map[string]interface{}{
				"namespace":   metric["namespace"],
				"dimensions":  buildWidgetDimensions(metric["dimensions"]),
				"metric_name": metric["metric_name"],
				"alias":       utils.ValueIgnoreEmpty(metric["alias"]),
			})
		}
		return params
	}
	return nil
}

func buildWidgetDimensions(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		return map[string]interface{}{
			"name":        raw["name"],
			"filter_type": raw["filter_type"],
			"values":      utils.ValueIgnoreEmpty(raw["values"]),
		}
	}
	return nil
}

func buildWidgetLocation(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		return map[string]interface{}{
			"top":    raw["top"],
			"left":   raw["left"],
			"width":  raw["width"],
			"height": raw["height"],
		}
	}
	return nil
}

func buildWidgetProperties(rawParams interface{}) map[string]interface{} {
	if rawProperties, ok := rawParams.([]interface{}); ok {
		if len(rawProperties) == 0 {
			return nil
		}
		raw := rawProperties[0].(map[string]interface{})
		return map[string]interface{}{
			"filter": utils.ValueIgnoreEmpty(raw["filter"]),
			"topN":   utils.ValueIgnoreEmpty(raw["top_n"]),
			"order":  utils.ValueIgnoreEmpty(raw["order"]),
		}
	}
	return nil
}

func resourceDashboardWidgetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getWidgetHttpUrl = "v2/{project_id}/widgets/{widget_id}"
		getWidgetProduct = "ces"
	)
	getWidgetClient, err := cfg.NewServiceClient(getWidgetProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	getWidgetPath := getWidgetClient.Endpoint + getWidgetHttpUrl
	getWidgetPath = strings.ReplaceAll(getWidgetPath, "{project_id}", getWidgetClient.ProjectID)
	getWidgetPath = strings.ReplaceAll(getWidgetPath, "{widget_id}", d.Id())

	getWidgetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	getWidgetResp, err := getWidgetClient.Request("GET", getWidgetPath, &getWidgetOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CES dashboard widget")
	}

	getWidgetRespBody, err := utils.FlattenResponse(getWidgetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	createdAt := utils.PathSearch("create_time", getWidgetRespBody, float64(0)).(float64) / 1000

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("title", utils.PathSearch("title", getWidgetRespBody, nil)),
		d.Set("metrics", flattenWidgetMetrics(getWidgetRespBody)),
		d.Set("view", utils.PathSearch("view", getWidgetRespBody, nil)),
		d.Set("metric_display_mode", utils.PathSearch("metric_display_mode", getWidgetRespBody, nil)),
		d.Set("location", flattenWidgetLocation(getWidgetRespBody)),
		d.Set("properties", flattenWidgetProperties(getWidgetRespBody)),
		d.Set("unit", utils.PathSearch("unit", getWidgetRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createdAt), false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenWidgetMetrics(resp interface{}) []interface{} {
	curJson := utils.PathSearch("metrics", resp, nil)

	if curArray, ok := curJson.([]interface{}); ok {
		params := make([]interface{}, 0, len(curArray))
		for _, metric := range curArray {
			params = append(params, map[string]interface{}{
				"namespace":   utils.PathSearch("namespace", metric, nil),
				"dimensions":  flattenWidgetDimensions(metric),
				"metric_name": utils.PathSearch("metric_name", metric, nil),
				"alias":       utils.PathSearch("alias", metric, nil),
			})
		}
		return params
	}
	return nil
}

func flattenWidgetLocation(resp interface{}) []interface{} {
	curJson := utils.PathSearch("location", resp, nil)
	if curJson == nil {
		return nil
	}

	params := map[string]interface{}{
		"top":    utils.PathSearch("top", curJson, nil),
		"left":   utils.PathSearch("left", curJson, nil),
		"width":  utils.PathSearch("width", curJson, nil),
		"height": utils.PathSearch("height", curJson, nil),
	}
	return []interface{}{params}
}

func flattenWidgetProperties(resp interface{}) []interface{} {
	curJson := utils.PathSearch("properties", resp, nil)
	if curJson == nil {
		return nil
	}

	params := map[string]interface{}{
		"filter": utils.PathSearch("filter", curJson, nil),
		"top_n":  utils.PathSearch("topN", curJson, 0),
		"order":  utils.PathSearch("order", curJson, nil),
	}
	return []interface{}{params}
}

func flattenWidgetDimensions(resp interface{}) []interface{} {
	curJson := utils.PathSearch("dimensions", resp, nil)
	if curJson == nil {
		return nil
	}

	params := map[string]interface{}{
		"name":        utils.PathSearch("name", curJson, nil),
		"filter_type": utils.PathSearch("filter_type", curJson, nil),
		"values":      utils.PathSearch("values", curJson, nil),
	}
	return []interface{}{params}
}

func resourceDashboardWidgetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateWidgetClient, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	updateWidgetPath := updateWidgetClient.Endpoint + "v2/{project_id}/widgets/batch-update"
	updateWidgetPath = strings.ReplaceAll(updateWidgetPath, "{project_id}", updateWidgetClient.ProjectID)
	updateWidgetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	param := utils.RemoveNil(buildUpdateWidgetBodyParams(d))
	updateWidgetOpt.JSONBody = []interface{}{param}
	_, err = updateWidgetClient.Request("POST", updateWidgetPath, &updateWidgetOpt)
	if err != nil {
		return diag.Errorf("error updating CES dashboard widget: %s", err)
	}

	return resourceDashboardWidgetRead(ctx, d, meta)
}

func buildUpdateWidgetBodyParams(d *schema.ResourceData) map[string]interface{} {
	param := map[string]interface{}{
		"widget_id": d.Id(),
	}

	if d.HasChange("title") {
		param["title"] = d.Get("title")
	}

	if d.HasChange("metrics") {
		param["metrics"] = buildWidgetMetrics(d.Get("metrics"))
	}

	if d.HasChange("metric_display_mode") {
		param["metric_display_mode"] = d.Get("metric_display_mode")
	}

	if d.HasChange("location") {
		param["location"] = buildWidgetLocation(d.Get("location"))
	}

	if d.HasChange("properties") {
		param["properties"] = buildWidgetProperties(d.Get("properties"))
	}

	if d.HasChange("unit") {
		param["unit"] = d.Get("unit")
	}
	return param
}

func resourceDashboardWidgetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteWidgetHttpUrl = "v2/{project_id}/widgets/{widget_id}"
		deleteWidgetProduct = "ces"
	)
	deleteWidgetClient, err := cfg.NewServiceClient(deleteWidgetProduct, region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	deleteWidgetPath := deleteWidgetClient.Endpoint + deleteWidgetHttpUrl
	deleteWidgetPath = strings.ReplaceAll(deleteWidgetPath, "{project_id}", deleteWidgetClient.ProjectID)
	deleteWidgetPath = strings.ReplaceAll(deleteWidgetPath, "{widget_id}", d.Id())

	deleteWidgetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	_, err = deleteWidgetClient.Request("DELETE", deleteWidgetPath, &deleteWidgetOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CES dashboard widget")
	}

	return nil
}

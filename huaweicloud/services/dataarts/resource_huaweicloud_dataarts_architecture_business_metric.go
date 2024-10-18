// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DataArts
// ---------------------------------------------------------------

package dataarts

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio POST /v2/{project_id}/design/biz-metrics
// @API DataArtsStudio GET /v2/{project_id}/design/biz-metrics/{id}
// @API DataArtsStudio PUT /v2/{project_id}/design/biz-metrics
// @API DataArtsStudio DELETE /v2/{project_id}/design/biz-metrics
func ResourceBusinessMetric() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBusinessMetricCreate,
		UpdateContext: resourceBusinessMetricUpdate,
		ReadContext:   resourceBusinessMetricRead,
		DeleteContext: resourceBusinessMetricDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDataArtsStudioImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of DataArts Studio workspace.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name.`,
			},
			"biz_catalog_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the process architecture ID.`,
			},
			"time_filters": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the statistical frequency.`,
			},
			"interval_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the refresh frequency.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of person responsible for the indicator.`,
			},
			"owner_department": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator management department name.`,
			},
			"destination": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the purpose of setting.`,
			},
			"definition": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator definition.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the calculation formula.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
			"apply_scenario": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application scenarios.`,
			},
			"technical_metric": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the related technical indicators.`,
			},
			"measure": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the measurement object.`,
			},
			"dimensions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the statistical dimension.`,
			},
			"general_filters": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the statistical caliber and modifiers.`,
			},
			"data_origin": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the data sources.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the unit of measurement.`,
			},
			"name_alias": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the indicator alias.`,
			},
			"code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the indicator encoding.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status.`,
			},
			"biz_catalog_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The process architecture path.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The editor.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
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

func resourceBusinessMetricCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/biz-metrics"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateBusinessMetricBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Architecture business metric: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	metricId := utils.PathSearch("data.value.id", createRespBody, "").(string)
	if metricId == "" {
		return diag.Errorf("unable to find the DataArts Architecture business metric ID from the API response")
	}
	d.SetId(metricId)

	return resourceBusinessMetricRead(ctx, d, meta)
}

func buildCreateOrUpdateBusinessMetricBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":             d.Get("name"),
		"biz_catalog_id":   d.Get("biz_catalog_id"),
		"time_filters":     d.Get("time_filters"),
		"interval_type":    d.Get("interval_type"),
		"owner":            d.Get("owner"),
		"owner_department": d.Get("owner_department"),
		"destination":      d.Get("destination"),
		"definition":       d.Get("definition"),
		"expression":       d.Get("expression"),
		"remark":           utils.ValueIgnoreEmpty(d.Get("description")),
		"apply_scenario":   utils.ValueIgnoreEmpty(d.Get("apply_scenario")),
		"technical_metric": utils.ValueIgnoreEmpty(d.Get("technical_metric")),
		"measure":          utils.ValueIgnoreEmpty(d.Get("measure")),
		"dimensions":       utils.ValueIgnoreEmpty(d.Get("dimensions")),
		"general_filters":  utils.ValueIgnoreEmpty(d.Get("general_filters")),
		"data_origin":      utils.ValueIgnoreEmpty(d.Get("data_origin")),
		"unit":             utils.ValueIgnoreEmpty(d.Get("unit")),
		"name_alias":       utils.ValueIgnoreEmpty(d.Get("name_alias")),
		"code":             utils.ValueIgnoreEmpty(d.Get("code")),
	}
}

func resourceBusinessMetricRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getResp, err := readBusinessMetric(client, d)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.3903"),
			"error retrieving DataArts Architecture business metric")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataValueBody := utils.PathSearch("data.value", getRespBody, make(map[string]interface{}))
	technicalMetricResp := utils.PathSearch("technical_metric", dataValueBody, "").(string)
	technicalMetric := utils.StringToInt(&technicalMetricResp)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", dataValueBody, nil)),
		d.Set("biz_catalog_id", utils.PathSearch("biz_catalog_id", dataValueBody, nil)),
		d.Set("time_filters", utils.PathSearch("time_filters", dataValueBody, nil)),
		d.Set("interval_type", utils.PathSearch("interval_type", dataValueBody, nil)),
		d.Set("owner", utils.PathSearch("owner", dataValueBody, nil)),
		d.Set("owner_department", utils.PathSearch("owner_department", dataValueBody, nil)),
		d.Set("destination", utils.PathSearch("destination", dataValueBody, nil)),
		d.Set("definition", utils.PathSearch("definition", dataValueBody, nil)),
		d.Set("expression", utils.PathSearch("expression", dataValueBody, nil)),
		d.Set("description", utils.PathSearch("remark", dataValueBody, nil)),
		d.Set("apply_scenario", utils.PathSearch("apply_scenario", dataValueBody, nil)),
		d.Set("technical_metric", technicalMetric),
		d.Set("measure", utils.PathSearch("measure", dataValueBody, nil)),
		d.Set("dimensions", utils.PathSearch("dimensions", dataValueBody, nil)),
		d.Set("general_filters", utils.PathSearch("general_filters", dataValueBody, nil)),
		d.Set("data_origin", utils.PathSearch("data_origin", dataValueBody, nil)),
		d.Set("unit", utils.PathSearch("unit", dataValueBody, nil)),
		d.Set("name_alias", utils.PathSearch("name_alias", dataValueBody, nil)),
		d.Set("code", utils.PathSearch("code", dataValueBody, nil)),
		d.Set("status", utils.PathSearch("status", dataValueBody, nil)),
		d.Set("biz_catalog_path", utils.PathSearch("biz_catalog_path", dataValueBody, nil)),
		d.Set("created_by", utils.PathSearch("create_by", dataValueBody, nil)),
		d.Set("updated_by", utils.PathSearch("update_by", dataValueBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", dataValueBody, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", dataValueBody, nil)),
		d.Set("technical_metric_type", utils.PathSearch("technical_metric_type", dataValueBody, nil)),
		d.Set("technical_metric_name", utils.PathSearch("technical_metric_name", dataValueBody, nil)),
		d.Set("l1", utils.PathSearch("l1", dataValueBody, nil)),
		d.Set("l2", utils.PathSearch("l2", dataValueBody, nil)),
		d.Set("l3", utils.PathSearch("l3", dataValueBody, nil)),
		d.Set("biz_metric", utils.PathSearch("biz_metric", dataValueBody, nil)),
		d.Set("summary_status", utils.PathSearch("summary_status", dataValueBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func readBusinessMetric(client *golangsdk.ServiceClient, d *schema.ResourceData) (*http.Response, error) {
	getPath := client.Endpoint + "v2/{project_id}/design/biz-metrics/{id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}

	return client.Request("GET", getPath, &getOpt)
}

func resourceBusinessMetricUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/biz-metrics"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateBody := utils.RemoveNil(buildCreateOrUpdateBusinessMetricBodyParams(d))
	updateBody["id"] = d.Id()
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         updateBody,
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DataArts Architecture business metric: %s", err)
	}
	return resourceBusinessMetricRead(ctx, d, meta)
}

func resourceBusinessMetricDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/design/biz-metrics"
		product = "dataarts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
		JSONBody:         buildDeleteBusinessMetricBodyParams(d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting DataArts Architecture business metric: %s", err)
	}

	// Successful deletion API call does not guarantee that the resource is successfully deleted.
	// Call the details API to confirm that the resource has been successfully deleted.
	_, err = readBusinessMetric(client, d)
	if err == nil {
		return diag.Errorf("error deleting DataArts Architecture business metric: the business metric still exists")
	}

	return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", "DLG.3903"),
		"error deleting DataArts Architecture business metric")
}

func buildDeleteBusinessMetricBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ids": []string{d.Id()},
	}
	return bodyParams
}

func resourceDataArtsStudioImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	partLength := len(parts)

	if partLength == 2 {
		d.SetId(parts[1])
		return []*schema.ResourceData{d}, d.Set("workspace_id", parts[0])
	}
	return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<id>")
}

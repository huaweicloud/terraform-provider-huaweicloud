package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	IndicatorNotExistsCode = "SecMaster.20030005"
)

// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/indicators
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/indicators
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/indicators/{id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/indicators/{id}
func ResourceIndicator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIndicatorCreate,
		UpdateContext: resourceIndicatorUpdate,
		ReadContext:   resourceIndicatorRead,
		DeleteContext: resourceIndicatorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIndicatorImportState,
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
				Description: `Specifies the ID of the workspace to which the indicator belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator name.`,
			},
			"type": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        indicatorIndicatorTypeSchema(),
				Required:    true,
				Description: `Specifies the indicator type.`,
			},
			"threat_degree": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator type.`,
			},
			"data_source": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        indicatorDataSourceSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the data source of the indicator.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator status.`,
			},
			"confidence": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the confidence of the indicator.`,
			},
			"first_occurrence_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the first occurrence time of the indicator.`,
			},
			"last_occurrence_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the last occurrence time of the indicator.`,
			},
			"granularity": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the granularity of the indicator.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the value of the indicator.`,
			},
			"invalid": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether the indicator is invalid.`,
			},
			"labels": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the labels of the indicator in comma-separated string.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The create time`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time`,
			},
		},
	}
}

func indicatorIndicatorTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the category.`,
			},
			"indicator_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator type.`,
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the indicator type ID.`,
			},
		},
	}
	return &sc
}

func indicatorDataSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the data source type.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product name.`,
			},
			"product_feature": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product feature.`,
			},
		},
	}
	return &sc
}

func resourceIndicatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIndicator: Create a SecMaster indicator.
	var (
		createIndicatorHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/indicators"
		createIndicatorProduct = "secmaster"
	)
	createIndicatorClient, err := cfg.NewServiceClient(createIndicatorProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createIndicatorPath := createIndicatorClient.Endpoint + createIndicatorHttpUrl
	createIndicatorPath = strings.ReplaceAll(createIndicatorPath, "{project_id}", createIndicatorClient.ProjectID)
	createIndicatorPath = strings.ReplaceAll(createIndicatorPath, "{workspace_id}", d.Get("workspace_id").(string))

	createIndicatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	createOpts, err := buildCreateIndicatorBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	createIndicatorOpt.JSONBody = utils.RemoveNil(createOpts)
	createIndicatorResp, err := createIndicatorClient.Request("POST", createIndicatorPath, &createIndicatorOpt)
	if err != nil {
		return diag.Errorf("error creating Indicator: %s", err)
	}

	createIndicatorRespBody, err := utils.FlattenResponse(createIndicatorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createIndicatorRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster indicator: ID is not found in API response")
	}
	d.SetId(id)

	return resourceIndicatorRead(ctx, d, meta)
}

func buildCreateIndicatorBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	dataObject := map[string]interface{}{
		"name":             d.Get("name"),
		"indicator_type":   buildIndicatorTypeOpts(d.Get("type")),
		"verdict":          d.Get("threat_degree"),
		"data_source":      buildIndicatorDataSourceOpts(d, cfg),
		"status":           d.Get("status"),
		"confidence":       d.Get("confidence"),
		"labels":           utils.ValueIgnoreEmpty(d.Get("labels")),
		"granular_marking": d.Get("granularity"),
		"value":            d.Get("value"),
		"environment":      buildEnvironmentOpts(d, cfg),
		"defanged":         d.Get("invalid"),
	}

	firstOccurrenceTimeWithZ, err := formatInputTime(d.Get("first_occurrence_time").(string))
	if err != nil {
		return nil, err
	}

	dataObject["first_report_time"] = firstOccurrenceTimeWithZ

	lastOccurrenceTimeWithZ, err := formatInputTime(d.Get("last_occurrence_time").(string))
	if err != nil {
		return nil, err
	}

	dataObject["last_report_time"] = lastOccurrenceTimeWithZ

	bodyParams := map[string]interface{}{
		"data_object": dataObject,
	}
	return bodyParams, nil
}

func buildIndicatorTypeOpts(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"category":       utils.ValueIgnoreEmpty(raw["category"]),
			"indicator_type": utils.ValueIgnoreEmpty(raw["indicator_type"]),
			"id":             utils.ValueIgnoreEmpty(raw["id"]),
		}
		return params
	}
	return nil
}

func buildIndicatorDataSourceOpts(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	rawArray := d.Get("data_source").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	region := cfg.GetRegion(d)

	raw, ok := rawArray[0].(map[string]interface{})
	if !ok {
		return nil
	}

	params := map[string]interface{}{
		"domain_id":       cfg.DomainID,
		"project_id":      cfg.GetProjectID(region),
		"region_id":       region,
		"product_feature": utils.ValueIgnoreEmpty(raw["product_feature"]),
		"product_name":    utils.ValueIgnoreEmpty(raw["product_name"]),
		"source_type":     utils.ValueIgnoreEmpty(raw["source_type"]),
	}
	return params
}

func resourceIndicatorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIndicator: Query the SecMaster indicator detail
	var (
		getIndicatorHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/indicators/{id}"
		getIndicatorProduct = "secmaster"
	)
	getIndicatorClient, err := cfg.NewServiceClient(getIndicatorProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getIndicatorPath := getIndicatorClient.Endpoint + getIndicatorHttpUrl
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{project_id}", getIndicatorClient.ProjectID)
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{workspace_id}", d.Get("workspace_id").(string))
	getIndicatorPath = strings.ReplaceAll(getIndicatorPath, "{id}", d.Id())

	getIndicatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getIndicatorResp, err := getIndicatorClient.Request("GET", getIndicatorPath, &getIndicatorOpt)

	if err != nil {
		// SecMaster.20010001: workspace ID not found
		// SecMaster.20030005: the incident not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", IndicatorNotExistsCode)
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster indicator")
	}

	getIndicatorRespBody, err := utils.FlattenResponse(getIndicatorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataObject := utils.PathSearch("data.data_object", getIndicatorRespBody, nil)
	if dataObject == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving SecMaster indicator")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", dataObject, nil)),
		d.Set("type", flattenGetIndicatorResponseBodyIndicatorType(dataObject)),
		d.Set("threat_degree", utils.PathSearch("verdict", dataObject, nil)),
		d.Set("data_source", flattenGetIndicatorResponseBodyDataSource(dataObject)),
		d.Set("status", utils.PathSearch("status", dataObject, nil)),
		d.Set("confidence", utils.PathSearch("confidence", dataObject, nil)),
		d.Set("labels", utils.PathSearch("labels", dataObject, nil)),
		d.Set("first_occurrence_time", utils.PathSearch("first_report_time", dataObject, nil)),
		d.Set("last_occurrence_time", utils.PathSearch("last_report_time", dataObject, nil)),
		d.Set("granularity", utils.PathSearch("granular_marking", dataObject, nil)),
		d.Set("value", utils.PathSearch("value", dataObject, nil)),
		d.Set("invalid", utils.PathSearch("defanged", dataObject, nil)),
		d.Set("created_at", utils.PathSearch("create_time", dataObject, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", dataObject, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetIndicatorResponseBodyIndicatorType(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("indicator_type", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"category":       utils.PathSearch("category", curJson, nil),
			"indicator_type": utils.PathSearch("indicator_type", curJson, nil),
			"id":             utils.PathSearch("id", curJson, nil),
		},
	}
	return rst
}

func flattenGetIndicatorResponseBodyDataSource(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("data_source", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"source_type":     utils.PathSearch("source_type", curJson, nil),
			"product_name":    utils.PathSearch("product_name", curJson, nil),
			"product_feature": utils.PathSearch("product_feature", curJson, nil),
		},
	}
	return rst
}

func resourceIndicatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateIndicatorChanges := []string{
		"name",
		"type",
		"threat_degree",
		"status",
		"confidence",
		"labels",
		"first_occurrence_time",
		"last_occurrence_time",
		"granularity",
		"value",
		"environment",
		"invalid",
	}

	if d.HasChanges(updateIndicatorChanges...) {
		// updateIndicator: Update the configuration of SecMaster indicator
		var (
			updateIndicatorHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/indicators/{id}"
			updateIndicatorProduct = "secmaster"
		)
		updateIndicatorClient, err := cfg.NewServiceClient(updateIndicatorProduct, region)
		if err != nil {
			return diag.Errorf("error creating SecMaster client: %s", err)
		}

		updateIndicatorPath := updateIndicatorClient.Endpoint + updateIndicatorHttpUrl
		updateIndicatorPath = strings.ReplaceAll(updateIndicatorPath, "{project_id}", updateIndicatorClient.ProjectID)
		updateIndicatorPath = strings.ReplaceAll(updateIndicatorPath, "{workspace_id}", d.Get("workspace_id").(string))
		updateIndicatorPath = strings.ReplaceAll(updateIndicatorPath, "{id}", d.Id())

		updateIndicatorOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		}

		updateOpts, err := buildUpdateIndicatorBodyParams(d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}

		updateIndicatorOpt.JSONBody = utils.RemoveNil(updateOpts)
		_, err = updateIndicatorClient.Request("PUT", updateIndicatorPath, &updateIndicatorOpt)
		if err != nil {
			return diag.Errorf("error updating SecMaster indicator: %s", err)
		}
	}
	return resourceIndicatorRead(ctx, d, meta)
}

func buildUpdateIndicatorBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	dataObject := map[string]interface{}{
		"name":             d.Get("name"),
		"indicator_type":   buildIndicatorTypeOpts(d.Get("type")),
		"verdict":          d.Get("threat_degree"),
		"status":           d.Get("status"),
		"confidence":       d.Get("confidence"),
		"labels":           utils.ValueIgnoreEmpty(d.Get("labels")),
		"granular_marking": d.Get("granularity"),
		"value":            d.Get("value"),
		"environment":      buildEnvironmentOpts(d, cfg),
		"defanged":         d.Get("invalid"),
	}

	if v, ok := d.GetOk("first_occurrence_time"); ok {
		firstOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["first_report_time"] = firstOccurrenceTimeWithZ
	}

	if v, ok := d.GetOk("last_occurrence_time"); ok {
		lastOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["last_report_time"] = lastOccurrenceTimeWithZ
	}

	bodyParams := map[string]interface{}{
		"data_object": dataObject,
	}
	return bodyParams, nil
}

func resourceIndicatorDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteIndicator: Delete an existing SecMaster indicator
	var (
		deleteIndicatorHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/indicators"
		deleteIndicatorProduct = "secmaster"
	)
	deleteIndicatorClient, err := cfg.NewServiceClient(deleteIndicatorProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deleteIndicatorPath := deleteIndicatorClient.Endpoint + deleteIndicatorHttpUrl
	deleteIndicatorPath = strings.ReplaceAll(deleteIndicatorPath, "{project_id}", deleteIndicatorClient.ProjectID)
	deleteIndicatorPath = strings.ReplaceAll(deleteIndicatorPath, "{workspace_id}", d.Get("workspace_id").(string))

	deleteIndicatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	deleteIndicatorOpt.JSONBody = utils.RemoveNil(buildDeleteIndicatorBodyParams(d))
	_, err = deleteIndicatorClient.Request("DELETE", deleteIndicatorPath, &deleteIndicatorOpt)
	if err != nil {
		return diag.Errorf("error deleting SecMaster indicator: %s", err)
	}

	return nil
}

func buildDeleteIndicatorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"batch_ids": []string{d.Id()},
	}
	return bodyParams
}

func resourceIndicatorImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<indicator_id>")
	}

	d.SetId(parts[1])

	mErr := multierror.Append(d.Set("workspace_id", parts[0]))

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

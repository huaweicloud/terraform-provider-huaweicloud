// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SecMaster
// ---------------------------------------------------------------

package secmaster

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	AlertNotExistsCode = "SecMaster.20030005"
)

// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/alerts
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/alerts
// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/alerts/{id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/alerts/{id}
func ResourceAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlertCreate,
		UpdateContext: resourceAlertUpdate,
		ReadContext:   resourceAlertRead,
		DeleteContext: resourceAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAlertImportState,
		},

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
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the alert.`,
			},
			"type": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     alertAlertTypeSchema(),
				Required: true,
			},
			"data_source": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     alertDataSourceSchema(),
				Required: true,
				ForceNew: true,
			},
			"severity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert severity.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert status.`,
			},
			"stage": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert stage.`,
			},
			"verification_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert verification status.`,
			},
			"first_occurrence_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the first occurrence time of the indicator.`,
			},
			"last_occurrence_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the last occurrence time of the indicator.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the owner name of the alert.`,
			},
			"debugging_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether it's a debugging data.`,
			},
			"labels": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the labels of the alert in comma-separated string.`,
			},
			"close_reason": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the close reason.`,
			},
			"close_comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the close comment.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The created time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updated time.`,
			},
		},
	}
}

func alertAlertTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the category.`,
			},
			"alert_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the alert type`,
			},
		},
	}
	return &sc
}

func alertDataSourceSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"product_feature": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product feature.`,
			},
			"product_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the product name.`,
			},
			"source_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the source type.`,
			},
		},
	}
	return &sc
}

func resourceAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAlert: Create a SecMaster alert.
	var (
		createAlertHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts"
		createAlertProduct = "secmaster"
	)
	createAlertClient, err := cfg.NewServiceClient(createAlertProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	createAlertPath := createAlertClient.Endpoint + createAlertHttpUrl
	createAlertPath = strings.ReplaceAll(createAlertPath, "{project_id}", createAlertClient.ProjectID)
	createAlertPath = strings.ReplaceAll(createAlertPath, "{workspace_id}", d.Get("workspace_id").(string))

	createAlertOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createOpts, err := buildCreateAlertBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	createAlertOpt.JSONBody = utils.RemoveNil(createOpts)
	createAlertResp, err := createAlertClient.Request("POST", createAlertPath, &createAlertOpt)
	if err != nil {
		return diag.Errorf("error creating Alert: %s", err)
	}

	createAlertRespBody, err := utils.FlattenResponse(createAlertResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.id", createAlertRespBody)
	if err != nil {
		return diag.Errorf("error creating Alert: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceAlertRead(ctx, d, meta)
}

func buildCreateAlertBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	dataObject := map[string]interface{}{
		"title":              d.Get("name"),
		"description":        d.Get("description"),
		"alert_type":         buildAlertTypeOpts(d.Get("type")),
		"data_source":        buildAlertDataSourceOpts(d, cfg),
		"severity":           d.Get("severity"),
		"handle_status":      d.Get("status"),
		"ipdrr_phase":        d.Get("stage"),
		"verification_state": d.Get("verification_status"),
		"owner":              utils.ValueIngoreEmpty(d.Get("owner")),
		"simulation":         utils.ValueIngoreEmpty(d.Get("debugging_data")),
		"labels":             utils.ValueIngoreEmpty(d.Get("labels")),
		"close_reason":       utils.ValueIngoreEmpty(d.Get("close_reason")),
		"close_comment":      utils.ValueIngoreEmpty(d.Get("close_comment")),
		"environment":        buildEnvironmentOpts(d, cfg),
		"domain_id":          cfg.DomainID,
		"region_id":          cfg.GetRegion(d),
	}

	if v, ok := d.GetOk("first_occurrence_time"); ok {
		firstOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["first_observed_time"] = firstOccurrenceTimeWithZ
	}

	if v, ok := d.GetOk("last_occurrence_time"); ok {
		lastOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["last_observed_time"] = lastOccurrenceTimeWithZ
	}

	bodyParams := map[string]interface{}{
		"data_object": dataObject,
	}
	return bodyParams, nil
}

func buildAlertTypeOpts(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw, ok := rawArray[0].(map[string]interface{})
		if !ok {
			return nil
		}

		params := map[string]interface{}{
			"category":   utils.ValueIngoreEmpty(raw["category"]),
			"alert_type": utils.ValueIngoreEmpty(raw["alert_type"]),
		}
		return params
	}
	return nil
}

func buildAlertDataSourceOpts(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
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
		"product_feature": utils.ValueIngoreEmpty(raw["product_feature"]),
		"product_name":    utils.ValueIngoreEmpty(raw["product_name"]),
		"source_type":     utils.ValueIngoreEmpty(raw["source_type"]),
	}
	return params
}

func resourceAlertRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAlert: Query the SecMaster alert detail
	var (
		getAlertHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts/{id}"
		getAlertProduct = "secmaster"
	)
	getAlertClient, err := cfg.NewServiceClient(getAlertProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	getAlertPath := getAlertClient.Endpoint + getAlertHttpUrl
	getAlertPath = strings.ReplaceAll(getAlertPath, "{project_id}", getAlertClient.ProjectID)
	getAlertPath = strings.ReplaceAll(getAlertPath, "{workspace_id}", d.Get("workspace_id").(string))
	getAlertPath = strings.ReplaceAll(getAlertPath, "{id}", d.Id())

	getAlertOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	getAlertResp, err := getAlertClient.Request("GET", getAlertPath, &getAlertOpt)

	if err != nil {
		if hasErrorCode(err, AlertNotExistsCode) {
			err = golangsdk.ErrDefault404{}
		}

		return common.CheckDeletedDiag(d, err, "error retrieving Alert")
	}

	getAlertRespBody, err := utils.FlattenResponse(getAlertResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataObject := utils.PathSearch("data.data_object", getAlertRespBody, nil)
	if dataObject == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving alert")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("title", dataObject, nil)),
		d.Set("description", utils.PathSearch("description", dataObject, nil)),
		d.Set("type", flattenGetAlertResponseBodyAlertType(dataObject)),
		d.Set("data_source", flattenGetAlertResponseBodyDataSource(dataObject)),
		d.Set("severity", utils.PathSearch("severity", dataObject, nil)),
		d.Set("status", utils.PathSearch("handle_status", dataObject, nil)),
		d.Set("stage", utils.PathSearch("ipdrr_phase", dataObject, nil)),
		d.Set("verification_status", utils.PathSearch("verification_state", dataObject, nil)),
		d.Set("first_occurrence_time", utils.PathSearch("first_observed_time", dataObject, nil)),
		d.Set("last_occurrence_time", utils.PathSearch("last_observed_time", dataObject, nil)),
		d.Set("owner", utils.PathSearch("owner", dataObject, nil)),
		d.Set("debugging_data", utils.PathSearch("simulation", dataObject, nil)),
		d.Set("labels", utils.PathSearch("labels", dataObject, nil)),
		d.Set("close_comment", utils.PathSearch("close_comment", dataObject, nil)),
		d.Set("close_reason", utils.PathSearch("close_reason", dataObject, nil)),
		d.Set("created_at", utils.PathSearch("create_time", dataObject, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", dataObject, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetAlertResponseBodyAlertType(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("alert_type", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing alert_type from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"category":   utils.PathSearch("category", curJson, nil),
			"alert_type": utils.PathSearch("alert_type", curJson, nil),
		},
	}
	return rst
}

func flattenGetAlertResponseBodyDataSource(resp interface{}) []interface{} {
	var rst []interface{}
	curJson, err := jmespath.Search("data_source", resp)
	if err != nil {
		log.Printf("[ERROR] error parsing data_source from response= %#v", resp)
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"product_feature": utils.PathSearch("product_feature", curJson, nil),
			"product_name":    utils.PathSearch("product_name", curJson, nil),
			"source_type":     utils.PathSearch("source_type", curJson, nil),
		},
	}
	return rst
}

func resourceAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// updateAlert: Update the configuration of SecMaster alert
	var (
		updateAlertHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts/{id}"
		updateAlertProduct = "secmaster"
	)
	updateAlertClient, err := cfg.NewServiceClient(updateAlertProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	updateAlertPath := updateAlertClient.Endpoint + updateAlertHttpUrl
	updateAlertPath = strings.ReplaceAll(updateAlertPath, "{project_id}", updateAlertClient.ProjectID)
	updateAlertPath = strings.ReplaceAll(updateAlertPath, "{workspace_id}", d.Get("workspace_id").(string))
	updateAlertPath = strings.ReplaceAll(updateAlertPath, "{id}", d.Id())

	updateAlertOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	updateOpts, err := buildUpdateAlertBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}

	updateAlertOpt.JSONBody = utils.RemoveNil(updateOpts)
	_, err = updateAlertClient.Request("PUT", updateAlertPath, &updateAlertOpt)
	if err != nil {
		return diag.Errorf("error updating Alert: %s", err)
	}

	return resourceAlertRead(ctx, d, meta)
}

func buildUpdateAlertBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	dataObject := map[string]interface{}{
		"title":              d.Get("name"),
		"description":        d.Get("description"),
		"alert_type":         buildAlertTypeOpts(d.Get("type")),
		"severity":           d.Get("severity"),
		"handle_status":      d.Get("status"),
		"ipdrr_phase":        d.Get("stage"),
		"verification_state": d.Get("verification_status"),
		"owner":              utils.ValueIngoreEmpty(d.Get("owner")),
		"labels":             utils.ValueIngoreEmpty(d.Get("labels")),
		"close_reason":       utils.ValueIngoreEmpty(d.Get("close_reason")),
		"close_comment":      utils.ValueIngoreEmpty(d.Get("close_comment")),
		"domain_id":          cfg.DomainID,
		"region_id":          cfg.GetRegion(d),
	}

	if v, ok := d.GetOk("first_occurrence_time"); ok {
		firstOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["first_observed_time"] = firstOccurrenceTimeWithZ
	}

	if v, ok := d.GetOk("last_occurrence_time"); ok {
		lastOccurrenceTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["last_observed_time"] = lastOccurrenceTimeWithZ
	}

	bodyParams := map[string]interface{}{
		"data_object": dataObject,
		"batch_ids":   []string{d.Id()},
	}
	return bodyParams, nil
}

func resourceAlertDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAlert: Delete an existing SecMaster alert
	var (
		deleteAlertHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/alerts"
		deleteAlertProduct = "secmaster"
	)
	deleteAlertClient, err := cfg.NewServiceClient(deleteAlertProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deleteAlertPath := deleteAlertClient.Endpoint + deleteAlertHttpUrl
	deleteAlertPath = strings.ReplaceAll(deleteAlertPath, "{project_id}", deleteAlertClient.ProjectID)
	deleteAlertPath = strings.ReplaceAll(deleteAlertPath, "{workspace_id}", d.Get("workspace_id").(string))

	deleteAlertOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	deleteAlertOpt.JSONBody = utils.RemoveNil(buildDeleteAlertBodyParams(d))
	_, err = deleteAlertClient.Request("DELETE", deleteAlertPath, &deleteAlertOpt)
	if err != nil {
		return diag.Errorf("error deleting Alert: %s", err)
	}

	return nil
}

func buildDeleteAlertBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"batch_ids": []string{d.Id()},
	}
	return bodyParams
}

func resourceAlertImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<alert_id>")
	}

	d.SetId(parts[1])

	err := d.Set("workspace_id", parts[0])
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

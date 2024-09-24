package secmaster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	GetIncidentNotFound = "SecMaster.20030005"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/incidents/{id}
// @API SecMaster PUT /v1/{project_id}/workspaces/{workspace_id}/soc/incidents/{id}
// @API SecMaster DELETE /v1/{project_id}/workspaces/{workspace_id}/soc/incidents
// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/soc/incidents
func ResourceIncident() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIncidentCreate,
		UpdateContext: resourceIncidentUpdate,
		ReadContext:   resourceIncidentRead,
		DeleteContext: resourceIncidentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIncidentImportState,
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
				Description: `Specifies the ID of the workspace to which the incident belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the incident name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the incident description.`,
			},
			"type": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        IncidentTypeSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the incident type configuration.`,
			},
			"level": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the incident level.`,
			},
			"status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the incident status.`,
			},
			"data_source": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        IncidentDataSourceSchema(),
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the data source configuration.`,
			},
			"first_occurrence_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the first occurrence time of the incident.`,
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the user name of the owner.`,
			},
			"last_occurrence_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the last occurrence time of the incident.`,
			},
			"planned_closure_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the planned closure time of the incident.`,
			},
			"verification_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the verification status.`,
			},
			"stage": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the stage of the incident.`,
			},
			"debugging_data": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether it's debugging data.`,
			},
			"labels": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the labels.`,
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
			"creator": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creator.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func IncidentTypeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"category": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the category.`,
			},
			"incident_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the incident type.`,
			},
		},
	}
	return &sc
}

func IncidentDataSourceSchema() *schema.Resource {
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

const (
	standardTimeFormat = "2006-01-02T15:04:05.000-07:00"
)

func resourceIncidentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createIncident: Create a SecMaster incident.
	var (
		createIncidentHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/incidents"
		createIncidentProduct = "secmaster"
	)
	createIncidentClient, err := cfg.NewServiceClient(createIncidentProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster Client: %s", err)
	}

	createIncidentPath := createIncidentClient.Endpoint + createIncidentHttpUrl
	createIncidentPath = strings.ReplaceAll(createIncidentPath, "{project_id}", createIncidentClient.ProjectID)
	createIncidentPath = strings.ReplaceAll(createIncidentPath, "{workspace_id}", fmt.Sprintf("%v", d.Get("workspace_id")))

	createIncidentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createOpts, err := buildIncidentBodyParams(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	createIncidentOpt.JSONBody = utils.RemoveNil(createOpts)
	createIncidentResp, err := createIncidentClient.Request("POST", createIncidentPath, &createIncidentOpt)
	if err != nil {
		return diag.Errorf("error creating SecMaster incident: %s", err)
	}

	createIncidentRespBody, err := utils.FlattenResponse(createIncidentResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createIncidentRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SecMaster incident: ID is not found in API response")
	}
	d.SetId(id)

	return resourceIncidentRead(ctx, d, meta)
}

func buildIncidentBodyParams(d *schema.ResourceData, cfg *config.Config) (map[string]interface{}, error) {
	dataObject := map[string]interface{}{
		"title":              utils.ValueIgnoreEmpty(d.Get("name")),
		"description":        utils.ValueIgnoreEmpty(d.Get("description")),
		"incident_type":      buildIncidentRequestBodyType(d.Get("type")),
		"severity":           utils.ValueIgnoreEmpty(d.Get("level")),
		"handle_status":      utils.ValueIgnoreEmpty(d.Get("status")),
		"owner":              utils.ValueIgnoreEmpty(d.Get("owner")),
		"data_source":        buildIncidentRequestBodyDataSource(d, cfg),
		"verification_state": utils.ValueIgnoreEmpty(d.Get("verification_status")),
		"ipdrr_phase":        utils.ValueIgnoreEmpty(d.Get("stage")),
		"simulation":         utils.ValueIgnoreEmpty(d.Get("debugging_data")),
		"labels":             utils.ValueIgnoreEmpty(d.Get("labels")),
		"close_reason":       utils.ValueIgnoreEmpty(d.Get("close_reason")),
		"close_comment":      utils.ValueIgnoreEmpty(d.Get("close_comment")),
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

	if v, ok := d.GetOk("planned_closure_time"); ok {
		plannedClosureTimeWithZ, err := formatInputTime(v.(string))
		if err != nil {
			return nil, err
		}

		dataObject["sla"] = plannedClosureTimeWithZ
	}

	bodyParams := map[string]interface{}{
		"data_object": dataObject,
	}
	return bodyParams, nil
}

func formatInputTime(timeFromSchema string) (string, error) {
	inputTimeFormat := "2006-01-02T15:04:05.000-0700"
	standardTime, err := time.Parse(standardTimeFormat, timeFromSchema)
	if err != nil {
		return "", fmt.Errorf("error parsing time to standard time: %s", err)
	}

	inputTimeWithoutZ := standardTime.Format(inputTimeFormat)
	inputTimeWithZ := inputTimeWithoutZ[:23] + "Z" + inputTimeWithoutZ[23:]

	return inputTimeWithZ, nil
}

func buildIncidentRequestBodyType(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"category":      utils.ValueIgnoreEmpty(raw["category"]),
			"incident_type": utils.ValueIgnoreEmpty(raw["incident_type"]),
		}
		return params
	}
	return nil
}

func buildIncidentRequestBodyDataSource(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
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

func buildEnvironmentOpts(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	region := cfg.GetRegion(d)
	return map[string]interface{}{
		"domain_id":  cfg.DomainID,
		"project_id": cfg.GetProjectID(region),
		"region_id":  region,
	}
}

func resourceIncidentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// getIncident: Query the SecMaster incident detail
	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster Client: %s", err)
	}

	dataObject, err := GetIncident(client, d.Get("workspace_id").(string), d.Id())
	if err != nil {
		// SecMaster.20010001: workspace ID not found
		// SecMaster.20030005: the incident not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		err = common.ConvertExpected400ErrInto404Err(err, "code", GetIncidentNotFound)
		return common.CheckDeletedDiag(d, err, "error retrieving SecMaster incident")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("title", dataObject, nil)),
		d.Set("description", utils.PathSearch("description", dataObject, nil)),
		d.Set("type", flattenGetIncidentResponseBodyType(dataObject)),
		d.Set("level", utils.PathSearch("severity", dataObject, nil)),
		d.Set("status", utils.PathSearch("handle_status", dataObject, nil)),
		d.Set("owner", utils.PathSearch("owner", dataObject, nil)),
		d.Set("data_source", flattenGetIncidentResponseBodyDataSource(dataObject)),
		d.Set("first_occurrence_time", utils.PathSearch("first_observed_time", dataObject, nil)),
		d.Set("last_occurrence_time", utils.PathSearch("last_observed_time", dataObject, nil)),
		d.Set("verification_status", utils.PathSearch("verification_state", dataObject, nil)),
		d.Set("stage", utils.PathSearch("ipdrr_phase", dataObject, nil)),
		d.Set("debugging_data", utils.PathSearch("simulation", dataObject, nil)),
		d.Set("labels", utils.PathSearch("labels", dataObject, nil)),
		d.Set("close_reason", utils.PathSearch("close_reason", dataObject, nil)),
		d.Set("close_comment", utils.PathSearch("close_comment", dataObject, nil)),
		d.Set("creator", utils.PathSearch("creator", dataObject, nil)),
		d.Set("created_at", utils.PathSearch("create_time", dataObject, nil)),
		d.Set("updated_at", utils.PathSearch("update_time", dataObject, nil)),
	)

	plannedClosureTime := utils.PathSearch("sla", dataObject, "").(string)
	if plannedClosureTime != "" {
		outputTimeFormat := "2006-01-02T15:04:05.000-0700"
		plannedClosureTimeWithoutZ := fmt.Sprintf(plannedClosureTime[:23] + plannedClosureTime[24:])
		plannedClosureTime, err := time.Parse(outputTimeFormat, plannedClosureTimeWithoutZ)
		if err != nil {
			return diag.Errorf("error parsing planned_closure_time: %s", err)
		}

		mErr = multierror.Append(
			mErr,
			d.Set("planned_closure_time", plannedClosureTime.Format(standardTimeFormat)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetIncidentResponseBodyType(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("incident_type", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"category":      utils.PathSearch("category", curJson, nil),
			"incident_type": utils.PathSearch("incident_type", curJson, nil),
		},
	}
	return rst
}

func flattenGetIncidentResponseBodyDataSource(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("data_source", resp, nil)
	if curJson == nil {
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

func resourceIncidentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateIncidentChanges := []string{
		"name",
		"description",
		"type",
		"level",
		"status",
		"owner",
		"first_occurrence_time",
		"last_occurrence_time",
		"planned_closure_time",
		"verification_status",
		"stage",
		"debugging_data",
		"labels",
		"close_reason",
		"close_comment",
	}

	if d.HasChanges(updateIncidentChanges...) {
		// updateIncident: Update the configuration of SecMaster incident
		var (
			updateIncidentHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/incidents/{id}"
			updateIncidentProduct = "secmaster"
		)
		updateIncidentClient, err := cfg.NewServiceClient(updateIncidentProduct, region)
		if err != nil {
			return diag.Errorf("error creating SecMaster client: %s", err)
		}

		updateIncidentPath := updateIncidentClient.Endpoint + updateIncidentHttpUrl
		updateIncidentPath = strings.ReplaceAll(updateIncidentPath, "{project_id}", updateIncidentClient.ProjectID)
		updateIncidentPath = strings.ReplaceAll(updateIncidentPath, "{workspace_id}", fmt.Sprintf("%v", d.Get("workspace_id")))
		updateIncidentPath = strings.ReplaceAll(updateIncidentPath, "{id}", d.Id())

		updateIncidentOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateOpts, err := buildIncidentBodyParams(d, cfg)
		if err != nil {
			return diag.FromErr(err)
		}
		updateIncidentOpt.JSONBody = utils.RemoveNil(updateOpts)
		_, err = updateIncidentClient.Request("PUT", updateIncidentPath, &updateIncidentOpt)
		if err != nil {
			return diag.Errorf("error updating SecMaster incident: %s", err)
		}
	}

	return resourceIncidentRead(ctx, d, meta)
}

func resourceIncidentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	// deleteIncident: Delete an existing SecMaster incident
	var (
		deleteIncidentHttpUrl = "v1/{project_id}/workspaces/{workspace_id}/soc/incidents"
		deleteIncidentProduct = "secmaster"
	)
	deleteIncidentClient, err := cfg.NewServiceClient(deleteIncidentProduct, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	deleteIncidentPath := deleteIncidentClient.Endpoint + deleteIncidentHttpUrl
	deleteIncidentPath = strings.ReplaceAll(deleteIncidentPath, "{project_id}", deleteIncidentClient.ProjectID)
	deleteIncidentPath = strings.ReplaceAll(deleteIncidentPath, "{workspace_id}", fmt.Sprintf("%v", workspaceID))

	deleteIncidentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	deleteIncidentOpt.JSONBody = utils.RemoveNil(buildDeleteIncidentBodyParams(d))
	_, err = deleteIncidentClient.Request("DELETE", deleteIncidentPath, &deleteIncidentOpt)
	if err != nil {
		// "SecMaster.20010001": workspace ID not found
		err = common.ConvertExpected403ErrInto404Err(err, "code", WorkspaceNotFound)
		return common.CheckDeletedDiag(d, err, "error deleting SecMaster incident")
	}

	dataObject, _ := GetIncident(deleteIncidentClient, d.Get("workspace_id").(string), d.Id())
	if dataObject != nil {
		return diag.Errorf("error deleting SecMaster incident, the incident still exists")
	}

	return nil
}

func buildDeleteIncidentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"batch_ids": utils.ValueIgnoreEmpty(d.Id()),
	}
	return bodyParams
}

func GetIncident(client *golangsdk.ServiceClient, workspaceId, id string) (interface{}, error) {
	// getIncident: Query the SecMaster incident detail
	getIncidentHttpUrl := "v1/{project_id}/workspaces/{workspace_id}/soc/incidents/{id}"
	getIncidentPath := client.Endpoint + getIncidentHttpUrl
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{project_id}", client.ProjectID)
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{workspace_id}", workspaceId)
	getIncidentPath = strings.ReplaceAll(getIncidentPath, "{id}", id)

	getIncidentOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getIncidentResp, err := client.Request("GET", getIncidentPath, &getIncidentOpt)
	if err != nil {
		return nil, err
	}

	getIncidentRespBody, err := utils.FlattenResponse(getIncidentResp)
	if err != nil {
		return nil, err
	}

	dataObject := utils.PathSearch("data.data_object", getIncidentRespBody, nil)
	if dataObject == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return dataObject, nil
}

func resourceIncidentImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <workspace_id>/<incident_id>")
	}

	d.SetId(parts[1])

	mErr := multierror.Append(d.Set("workspace_id", parts[0]))

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

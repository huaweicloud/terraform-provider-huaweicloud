package geminidb

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

var geminiDBParameterTemplateNonUpdatableParams = []string{
	"instance_id",
	"datastore",
	"datastore.*.type",
	"datastore.*.version",
	"datastore.*.mode",
}

// @API GaussDBforNoSQL POST /v3/{project_id}/configurations
// @API GaussDBforNoSQL PUT /v3/{project_id}/configurations/{config_id}
// @API GaussDBforNoSQL GET /v3/{project_id}/configurations/{config_id}
// @API GaussDBforNoSQL DELETE /v3/{project_id}/configurations/{config_id}
func ResourceGeminiDBParameterTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBParameterTemplateCreate,
		UpdateContext: resourceGeminiDBParameterTemplateUpdate,
		ReadContext:   resourceGeminiDBParameterTemplateRead,
		DeleteContext: resourceGeminiDBParameterTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(geminiDBParameterTemplateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     geminiDBParameterTemplateDatastoreSchema(),
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"datastore_version_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBParameterTemplateParameterSchema(),
			},
		},
	}
}

func geminiDBParameterTemplateDatastoreSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func geminiDBParameterTemplateParameterSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restart_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGeminiDBParameterTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateGeminiDBParameterTemplateBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GeminiDB parameter template: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	configurationID := utils.PathSearch("configuration.id", createRespBody, "").(string)
	if configurationID == "" {
		return diag.Errorf("error creating GeminiDB parameter template: configuration ID is not found in API response")
	}

	d.SetId(configurationID)

	return resourceGeminiDBParameterTemplateRead(ctx, d, meta)
}

func buildCreateGeminiDBParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"datastore":   buildGeminiDBParameterTemplateDatastore(d.Get("datastore").([]interface{})),
		"values":      utils.ValueIgnoreEmpty(d.Get("values")),
		"instance_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
	}
	return bodyParams
}

func buildGeminiDBParameterTemplateDatastore(datastore []interface{}) map[string]interface{} {
	if len(datastore) == 0 {
		return nil
	}

	datastoreMap := datastore[0].(map[string]interface{})
	result := map[string]interface{}{
		"type":    datastoreMap["type"],
		"version": datastoreMap["version"],
	}

	if mode, ok := datastoreMap["mode"]; ok {
		result["mode"] = mode
	}

	return result
}

func resourceGeminiDBParameterTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB parameter template")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("datastore_version_name", utils.PathSearch("datastore_version_name", getRespBody, nil)),
		d.Set("datastore_name", utils.PathSearch("datastore_name", getRespBody, nil)),
		d.Set("mode", utils.PathSearch("mode", getRespBody, nil)),
		d.Set("created", utils.PathSearch("created", getRespBody, nil)),
		d.Set("updated", utils.PathSearch("updated", getRespBody, nil)),
		d.Set("configuration_parameters", flattenGeminiDBParameterTemplateParameters(getRespBody)),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting GeminiDB parameter template fields: %s", err)
	}

	return nil
}

func resourceGeminiDBParameterTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{config_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildUpdateGeminiDBParameterTemplateBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating GeminiDB parameter template: %s", err)
	}

	return resourceGeminiDBParameterTemplateRead(ctx, d, meta)
}

func buildUpdateGeminiDBParameterTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
		"values":      utils.ValueIgnoreEmpty(d.Get("values")),
	}
	return bodyParams
}

func resourceGeminiDBParameterTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{config_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting GeminiDB parameter template: %s", err)
	}

	return nil
}

func flattenGeminiDBParameterTemplateParameters(resp interface{}) []map[string]interface{} {
	parametersRaw := utils.PathSearch("configuration_parameters", resp, nil)
	if parametersRaw == nil {
		return nil
	}

	parametersSlice, ok := parametersRaw.([]interface{})
	if !ok {
		return nil
	}

	parameters := make([]map[string]interface{}, 0, len(parametersSlice))
	for _, paramRaw := range parametersSlice {
		param, ok := paramRaw.(map[string]interface{})
		if !ok {
			continue
		}

		paramMap := map[string]interface{}{
			"name":             utils.PathSearch("name", param, nil),
			"value":            utils.PathSearch("value", param, nil),
			"restart_required": utils.PathSearch("restart_required", param, nil),
			"readonly":         utils.PathSearch("readonly", param, nil),
			"value_range":      utils.PathSearch("value_range", param, nil),
			"type":             utils.PathSearch("type", param, nil),
			"description":      utils.PathSearch("description", param, nil),
		}
		parameters = append(parameters, paramMap)
	}

	return parameters
}

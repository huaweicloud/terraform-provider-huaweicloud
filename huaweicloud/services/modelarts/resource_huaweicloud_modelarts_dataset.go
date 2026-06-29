package modelarts

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type datasetSourceType int

const (
	datasetSourceTypeOBS datasetSourceType = 0
	datasetSourceTypeDWS datasetSourceType = 1
	datasetSourceTypeDLI datasetSourceType = 2
	datasetSourceTypeMRS datasetSourceType = 4
)

var (
	datasetNonUpdatableParams = []string{
		"type",
		"output_path",
		"data_source",
		"data_source.*.data_type",
		"data_source.*.path",
		"data_source.*.with_column_header",
		"data_source.*.queue_name",
		"data_source.*.database_name",
		"data_source.*.table_name",
		"data_source.*.user_name",
		"data_source.*.password",
		"data_source.*.cluster_id",
		"schemas",
		"schemas.*.type",
		"schemas.*.name",
		"import_labeled_enabled",
		"label_format",
		"label_format.*.type",
		"label_format.*.text_label_separator",
		"label_format.*.label_separator",
	}
	datasetNotFoundErrCodes = []string{
		"ModelArts.4352",
	}
)

// @API ModelArts DELETE /v2/{project_id}/datasets/{dataset_id}
// @API ModelArts GET /v2/{project_id}/datasets/{dataset_id}
// @API ModelArts PUT /v2/{project_id}/datasets/{dataset_id}
// @API ModelArts POST /v2/{project_id}/datasets
func ResourceDataset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasetCreate,
		ReadContext:   resourceDatasetRead,
		UpdateContext: resourceDatasetUpdate,
		DeleteContext: resourceDatasetDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(datasetNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the dataset is located.`,
			},

			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the dataset.`,
			},
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The type of the dataset.`,
			},
			"output_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The OBS storage path that used to store output files.`,
			},
			"data_source": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem:     dataSourceSchemaResource(),
				Description: `The data sources which be used to imported the source data (such as pictures/files/audio,
etc.) in this directory and subdirectories to the dataset.`,
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the dataset.`,
			},
			"schemas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The field type of the schema.`,
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The field name of the schema.`,
						},
					},
				},
				Description: `The schema configurations of the dataset.`,
			},
			"import_labeled_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether to enable the import labeled features.`,
			},
			"label_format": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "1",
							Description: `The type of the label format.`,
						},
						"text_label_separator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The separator between text and label.`,
						},
						"label_separator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The separator between label and label.`,
						},
					},
				},
				Description: `The custom format of labeled features when import labeled files for Text classification.`,
			},
			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the label.`,
						},
						"property_color": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The color of the label.`,
						},
						"property_shape": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The shape of the label.`,
						},
						"property_shortcut": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `The shortcut of the label.`,
						},
					},
				},
				Description: `The labels of the dataset.`,
			},

			// Attributes.
			"data_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data format of the dataset.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the dataset, in RFC3339 format.`,
			},
			"status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The status of the dataset.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func dataSourceSchemaResource() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: `The type of the data source.`,
			},
			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The OBS storage path or MRS HDFS path.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the DWS/MRS cluster.`,
			},
			"database_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DWS/DLI database.`,
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DWS/DLI table.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DWS database user.`,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The password of the DWS database user.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the DLI queue.`,
			},
			"with_column_header": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: `Whether the data contains table header when the type of dataset is 400 (table type).`,
			},
		},
	}

	return &nodeResource
}

func buildDatasetDataSource(dataSources []interface{}) ([]interface{}, error) {
	var (
		result = make([]interface{}, 0, len(dataSources))
	)
	for _, dataSource := range dataSources {
		dataType := utils.PathSearch("data_type", dataSource, 0).(int)

		item := map[string]interface{}{
			"data_type":          dataType,
			"with_column_header": utils.PathSearch("with_column_header", dataSource, true).(bool),
		}
		switch dataType {
		case int(datasetSourceTypeOBS):
			path := utils.PathSearch("path", dataSource, "").(string)
			if path == "" {
				return nil, errors.New("when import data from OBS, path is required")
			}
			item["data_path"] = path
		case int(datasetSourceTypeDWS):
			clusterId := utils.PathSearch("cluster_id", dataSource, "").(string)
			databaseName := utils.PathSearch("database_name", dataSource, "").(string)
			tableName := utils.PathSearch("table_name", dataSource, "").(string)
			userName := utils.PathSearch("user_name", dataSource, "").(string)
			password := utils.PathSearch("password", dataSource, "").(string)
			if clusterId == "" || databaseName == "" || tableName == "" || userName == "" || password == "" {
				return nil, errors.New("when import data from DWS, both cluster_id, database_name, table_name, user_name and password are required")
			}
			item["source_info"] = map[string]interface{}{
				"cluster_id":    clusterId,
				"database_name": databaseName,
				"table_name":    tableName,
				"user_name":     userName,
				"user_password": password,
			}
		case int(datasetSourceTypeDLI):
			queueName := utils.PathSearch("queue_name", dataSource, "").(string)
			databaseName := utils.PathSearch("database_name", dataSource, "").(string)
			tableName := utils.PathSearch("table_name", dataSource, "").(string)
			if queueName == "" || databaseName == "" || tableName == "" {
				return nil, errors.New("when import data from DLI, both queue_name, database_name and table_name are required")
			}
			item["source_info"] = map[string]interface{}{
				"queue_name":    queueName,
				"database_name": databaseName,
				"table_name":    tableName,
			}
		case int(datasetSourceTypeMRS):
			clusterId := utils.PathSearch("cluster_id", dataSource, "").(string)
			path := utils.PathSearch("path", dataSource, "").(string)
			if clusterId == "" || path == "" {
				return nil, errors.New("when import data from MRS, both cluster_id and path are required")
			}
			item["source_info"] = map[string]interface{}{
				"cluster_id": clusterId,
				"input":      path,
			}
		}
		result = append(result, item)
	}
	return result, nil
}

func buildDatasetSchema(datasetType int, schemas []interface{}) ([]interface{}, error) {
	if datasetType == 400 && len(schemas) < 1 {
		return nil, errors.New("the schema cannot be empty if type is 400 (table type)")
	}

	result := make([]interface{}, len(schemas))
	for i, item := range schemas {
		result[i] = map[string]interface{}{
			"schema_id": i + 1,
			"type":      utils.PathSearch("type", item, ""),
			"name":      utils.PathSearch("name", item, ""),
		}
	}
	return result, nil
}

func buildDatasetLabelFormat(labelFormats []interface{}) map[string]interface{} {
	if len(labelFormats) < 1 {
		return nil
	}

	return map[string]interface{}{
		"label_type":            utils.PathSearch("type", labelFormats[0], ""),
		"text_label_separator":  utils.PathSearch("text_label_separator", labelFormats[0], ""),
		"text_sample_separator": utils.PathSearch("label_separator", labelFormats[0], ""),
	}
}

func buildDatasetLabels(labels []interface{}) []interface{} {
	if len(labels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"name": utils.PathSearch("name", label, ""),
			"property": utils.RemoveNil(map[string]interface{}{
				"@modelarts:color":         utils.ValueIgnoreEmpty(utils.PathSearch("property_color", label, "")),
				"@modelarts:default_shape": utils.ValueIgnoreEmpty(utils.PathSearch("property_shape", label, "")),
				"@modelarts:shortcut":      utils.ValueIgnoreEmpty(utils.PathSearch("property_shortcut", label, "")),
			}),
		})
	}
	return result
}

func buildCreateDatasetBodyParams(d *schema.ResourceData) (map[string]interface{}, error) {
	datasetType := d.Get("type").(int)
	dataSources, err := buildDatasetDataSource(d.Get("data_source").([]interface{}))
	if err != nil {
		return nil, err
	}
	schemas, err := buildDatasetSchema(datasetType, d.Get("schemas").([]interface{}))
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{
		// Fixed values.
		"work_path_type": 0,
		"import_data":    true,
		// Required parameters.
		"dataset_name": d.Get("name").(string),
		"dataset_type": datasetType,
		"work_path":    d.Get("output_path").(string),
		"data_sources": dataSources,
		// Optional parameters.
		"description":        d.Get("description").(string),
		"import_annotations": utils.Bool(d.Get("import_labeled_enabled").(bool)),
		"schema":             schemas,
		"label_format":       buildDatasetLabelFormat(d.Get("label_format").([]interface{})),
		"labels":             buildDatasetLabels(d.Get("labels").([]interface{})),
	}
	return result, nil
}

func createDataset(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	httpUrl := "v2/{project_id}/datasets"

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	bodyParams, err := buildCreateDatasetBodyParams(d)
	if err != nil {
		return nil, err
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(bodyParams),
	}

	requestResp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func waitForDatasetStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, datasetId string, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshDatasetStatusFunc(client, datasetId, []string{"1"}),
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 20 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts dataset (%s) status to be completed: %s", datasetId, err)
	}
	return nil
}

func GetDatasetById(client *golangsdk.ServiceClient, datasetId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/datasets/{dataset_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{dataset_id}", datasetId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(resp)
}

func refreshDatasetStatusFunc(client *golangsdk.ServiceClient, datasetId string, targets []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := GetDatasetById(client, datasetId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && len(targets) < 1 {
				return "RESOURCE_NOT_FOUND", "COMPLETED", nil
			}
			return respBody, "ERROR", err
		}

		unexpectedStatus := []string{
			"4", // EXCEPTION
		}

		status := fmt.Sprintf("%v", utils.PathSearch("status", respBody, float64(0)).(float64))
		if utils.StrSliceContains(unexpectedStatus, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		if utils.StrSliceContains(targets, status) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceDatasetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := createDataset(client, d)
	if err != nil {
		return diag.Errorf("error creating ModelArts datasets: %s", err)
	}

	datasetId := utils.PathSearch("dataset_id", respBody, "").(string)
	if datasetId == "" {
		return diag.Errorf("unable to find the ModelArts dataset ID from the API response")
	}
	d.SetId(datasetId)

	err = waitForDatasetStatusCompleted(ctx, client, datasetId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDatasetRead(ctx, d, meta)
}

func flattenDatasetDataSource(dataSources []interface{}) []interface{} {
	if len(dataSources) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(dataSources))
	for _, dataSource := range dataSources {
		dataType := int(utils.PathSearch("data_type", dataSource, float64(0)).(float64))
		elem := map[string]interface{}{
			"data_type":          dataType,
			"cluster_id":         utils.PathSearch("source_info.cluster_id", dataSource, "").(string),
			"queue_name":         utils.PathSearch("source_info.queue_name", dataSource, "").(string),
			"path":               utils.PathSearch("data_path", dataSource, "").(string),
			"with_column_header": utils.PathSearch("with_column_header", dataSource, true).(bool),
			"database_name":      utils.PathSearch("source_info.database_name", dataSource, "").(string),
			"table_name":         utils.PathSearch("source_info.table_name", dataSource, "").(string),
			"user_name":          utils.PathSearch("source_info.user_name", dataSource, "").(string),
		}
		if dataType == int(datasetSourceTypeMRS) {
			elem["path"] = utils.PathSearch("source_info.input", dataSource, "").(string)
		}
		result = append(result, elem)
	}
	return result
}

func flattenDatasetSchema(schemas []interface{}) []interface{} {
	if len(schemas) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(schemas))

	for _, item := range schemas {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("type", item, "").(string),
			"name": utils.PathSearch("name", item, "").(string),
		})
	}
	return result
}

func flattenDatasetLabels(labels []interface{}) []interface{} {
	if len(labels) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(labels))
	for _, label := range labels {
		result = append(result, map[string]interface{}{
			"name":              utils.PathSearch("name", label, "").(string),
			"property_color":    utils.PathSearch("property.@modelarts:color", label, "").(string),
			"property_shape":    utils.PathSearch("property.@modelarts:default_shape", label, "").(string),
			"property_shortcut": utils.PathSearch("property.@modelarts:shortcut", label, "").(string),
		})
	}
	return result
}

func resourceDatasetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	detail, err := GetDatasetById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", datasetNotFoundErrCodes...),
			"error retrieving ModelArts dataset")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		// Required parameters.
		d.Set("name", utils.PathSearch("dataset_name", detail, nil)),
		d.Set("type", utils.PathSearch("dataset_type", detail, nil)),
		d.Set("output_path", utils.PathSearch("work_path", detail, nil)),
		d.Set("data_source", flattenDatasetDataSource(utils.PathSearch("data_sources", detail, make([]interface{}, 0)).([]interface{}))),
		// Optional parameters.
		d.Set("description", utils.PathSearch("description", detail, nil)),
		d.Set("schemas", flattenDatasetSchema(utils.PathSearch("schema", detail, make([]interface{}, 0)).([]interface{}))),
		d.Set("labels", flattenDatasetLabels(utils.PathSearch("labels", detail, make([]interface{}, 0)).([]interface{}))),
		d.Set("data_format", utils.PathSearch("data_format", detail, nil)),
		d.Set("import_labeled_enabled", utils.PathSearch("import_data", detail, nil)),
		d.Set("created_at", utils.FormatTimeStampUTC(int64(utils.PathSearch("create_time", detail, float64(0)).(float64))/1000)),
		d.Set("status", int(utils.PathSearch("status", detail, float64(0)).(float64))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDatasetBodyParams(d *schema.ResourceData) map[string]interface{} {
	result := map[string]interface{}{
		// Required parameters.
		"dataset_name": d.Get("name").(string),
		// Optional parameters.
		"description": d.Get("description").(string),
	}

	if d.HasChange("labels") {
		oldVal, newVal := d.GetChange("labels")
		result["add_labels"] = utils.ValueIgnoreEmpty(buildDatasetLabels(newVal.([]interface{})))
		result["delete_labels"] = utils.ValueIgnoreEmpty(buildDatasetLabels(oldVal.([]interface{})))
	}

	return result
}

func updateDataset(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl   = "v2/{project_id}/datasets/{dataset_id}"
		datasetId = d.Id()
	)

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{dataset_id}", datasetId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDatasetBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", updatePath, &opt)
	return err
}

func resourceDatasetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = updateDataset(client, d)
	if err != nil {
		return diag.Errorf("error updating ModelArts dataset: %s", err)
	}

	return resourceDatasetRead(ctx, d, meta)
}

func deleteDataset(client *golangsdk.ServiceClient, datasetId string) error {
	httpUrl := "v2/{project_id}/datasets/{dataset_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{dataset_id}", datasetId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &opt)
	return err
}

func resourceDatasetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	err = deleteDataset(client, d.Id())
	if err != nil {
		return diag.Errorf("error deleting ModelArts dataset: %s", err)
	}
	return nil
}

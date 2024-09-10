package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts DELETE /v2/{project_id}/datasets/{datasetId}
// @API ModelArts GET /v2/{project_id}/datasets/{datasetId}
// @API ModelArts PUT /v2/{project_id}/datasets/{datasetId}
// @API ModelArts POST /v2/{project_id}/datasets
func ResourceDataset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatasetCreate,
		ReadContext:   resourceDatasetRead,
		UpdateContext: resourceDatasetUpdate,
		DeleteContext: resourceDatasetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

			"type": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 3, 100, 101, 102, 200, 201, 202, 400, 600, 900}),
			},

			"output_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"data_source": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     dataSourceSchemaResource(),
			},

			"schemas": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{"String", "Short", "Int", "Long", "Double",
								"Float", "Byte", "Date", "Timestamp", "Boolean"}, false),
						},

						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"import_labeled_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"label_format": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     true,
							Default:      "1",
							ValidateFunc: validation.StringInSlice([]string{"0", "1"}, false),
						},

						"text_label_separator": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},

						"label_separator": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Description: "It is required only the dataType=100",
			},

			"labels": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"property_color": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},

						"property_shape": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ValidateFunc: validation.StringInSlice([]string{"bndbox", "polygon", "circle", "line",
								"dashed", "point", "polyline"}, false),
						},

						"property_shortcut": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"data_format": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func dataSourceSchemaResource() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      0,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2, 4}),
			},

			"path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"with_column_header": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},

			"queue_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"database_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"table_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"user_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}

	return &nodeResource
}

func resourceDatasetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	opts, err := buildCreateParamter(d)
	if err != nil {
		return diag.FromErr(err)
	}
	rst, err := dataset.Create(client, *opts)
	if err != nil {
		return diag.Errorf("error creating ModelArts datasets: %s", err)
	}

	d.SetId(rst.DatasetId)

	err = waitingforDatasetCreated(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceDatasetRead(ctx, d, meta)
}

func resourceDatasetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	detail, err := dataset.Get(client, d.Id(), dataset.GetOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDatasetErrorToError404(err), "error retrieving ModelArts dataset")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.DatasetName),
		d.Set("type", detail.DatasetType),
		d.Set("data_format", detail.DataFormat),
		d.Set("output_path", detail.WorkPath),
		d.Set("description", detail.Description),
		setDataSourcesToState(d, detail.DataSources[0]),
		setSchemaToState(d, detail.Schema),
		d.Set("import_labeled_enabled", detail.ImportData),
		setLabelsToState(d, detail.Labels),
		d.Set("created_at", utils.FormatTimeStampUTC(int64(detail.CreateTime)/1000)),
		d.Set("status", detail.Status),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDatasetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	_, err = dataset.Get(client, d.Id(), dataset.GetOpts{})
	if err != nil {
		return common.CheckDeletedDiag(d, parseDatasetErrorToError404(err), "error retrieving ModelArts dataset")
	}

	desc := d.Get("description").(string)
	updateParams := dataset.UpdateOpts{
		DatasetName: d.Get("name").(string),
		Description: &desc,
	}

	if d.HasChange("labels") {
		o, n := d.GetChange("labels")
		updateParams.AddLabels = buildLabelsParamter(n)
		updateParams.DeleteLabels = buildLabelsParamter(o)
	}

	rst := dataset.Update(client, d.Id(), updateParams)
	if rst.Err != nil {
		return diag.Errorf("update ModelArts dataset=%s failed, error: %s", d.Id(), err)
	}

	return resourceDatasetRead(ctx, d, meta)
}

func resourceDatasetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	dErr := dataset.Delete(client, d.Id())
	if dErr.Err != nil {
		return common.CheckDeletedDiag(d, parseDatasetErrorToError404(err), "Delete ModelArts dataset failed")
	}

	return nil
}

func waitingforDatasetCreated(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"0", "5", "6", "7", "8"},
		Target:  []string{"1"},
		Refresh: func() (interface{}, string, error) {
			resp, err := dataset.Get(client, id, dataset.GetOpts{})
			if err != nil {
				return nil, "", err
			}
			return resp, fmt.Sprint(resp.Status), nil
		},
		Timeout:      timeout,
		PollInterval: 20 * time.Second,
		Delay:        20 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts dataset (%s) to be created: %s", id, err)
	}
	return nil
}

func buildCreateParamter(d *schema.ResourceData) (*dataset.CreateOpts, error) {
	dataSources, err := buildDataSourceParamter(d)
	if err != nil {
		return nil, err
	}

	dataType := d.Get("type").(int)
	schemas, err := buildSchemaParamter(d.Get("schemas"), dataType)
	if err != nil {
		return nil, err
	}

	rst := dataset.CreateOpts{
		DatasetName:       d.Get("name").(string),
		DatasetType:       dataType,
		WorkPathType:      0,
		DataSources:       dataSources,
		WorkPath:          d.Get("output_path").(string),
		Description:       d.Get("description").(string),
		ImportData:        true,
		ImportAnnotations: utils.Bool(d.Get("import_labeled_enabled").(bool)),
		Schema:            schemas,
		Labels:            buildLabelsParamter(d.Get("labels")),
	}

	if format := buildLabelFormatParamter(d); format != nil {
		rst.LabelFormat = *format
	}

	return &rst, nil
}

func buildLabelFormatParamter(d *schema.ResourceData) *dataset.LabelFormat {
	if v, ok := d.GetOk("label_format"); ok {
		configRaw := v.([]interface{})[0].(map[string]interface{})
		configs := dataset.LabelFormat{
			LabelType:           configRaw["type"].(string),
			TextLabelSeparator:  configRaw["text_label_separator"].(string),
			TextSampleSeparator: configRaw["label_separator"].(string),
		}
		return &configs
	}
	return nil
}

func buildDataSourceParamter(d *schema.ResourceData) (dataSources []dataset.DataSource, err error) {
	item := d.Get("data_source").([]interface{})[0].(map[string]interface{})

	dataType := item["data_type"].(int)
	dataSource := dataset.DataSource{
		DataType:         dataType,
		WithColumnHeader: utils.Bool(item["with_column_header"].(bool)),
	}

	path := item["path"].(string)
	// OBS check
	if dataType == 0 {
		if path == "" {
			err = fmt.Errorf("when import data from OBS, path is required")
			return
		}
		dataSource.DataPath = path
	}

	clusterId := item["cluster_id"].(string)
	databaseName := item["database_name"].(string)
	tableName := item["table_name"].(string)
	userName := item["user_name"].(string)
	password := item["password"].(string)

	// DWS check
	if dataType == 1 {
		if clusterId == "" || databaseName == "" || tableName == "" || userName == "" || password == "" {
			err = fmt.Errorf("when import data from DWS, cluster_id, database_name, table_name, user_name and" +
				" password are required")
			return
		}

		dataSource.SourceInfo.ClusterId = clusterId
		dataSource.SourceInfo.DatabaseName = databaseName
		dataSource.SourceInfo.TableName = tableName
		dataSource.SourceInfo.UserName = userName
		dataSource.SourceInfo.UserPassword = password
	}

	queueName := item["queue_name"].(string)

	// DLI check
	if dataType == 2 {
		if queueName == "" || databaseName == "" || tableName == "" {
			err = fmt.Errorf("when import data from DLI, queue_name, database_name and table_name are required")
		}

		dataSource.SourceInfo.QueueName = queueName
		dataSource.SourceInfo.DatabaseName = databaseName
		dataSource.SourceInfo.TableName = tableName
	}

	// MRS check
	if dataType == 4 {
		if clusterId == "" || path == "" {
			err = fmt.Errorf("when import data from MRS, cluster_id and path are required")
		}

		dataSource.SourceInfo.ClusterId = clusterId
		dataSource.SourceInfo.Input = path
	}

	dataSources = append(dataSources, dataSource)
	return
}

func buildLabelsParamter(v interface{}) (labels []dataset.Label) {
	if v != nil {
		configRaw := v.([]interface{})
		for _, item := range configRaw {
			tmp := item.(map[string]interface{})
			labels = append(labels, dataset.Label{
				Name: tmp["name"].(string),
				Property: dataset.LabelProperty{
					Color:        tmp["property_color"].(string),
					DefaultShape: tmp["property_shape"].(string),
					Shortcut:     tmp["property_shortcut"].(string),
				},
			})
		}
	}
	return
}

func buildSchemaParamter(v interface{}, dataType int) (schemas []dataset.Field, err error) {
	if v != nil && dataType == 400 {
		configRaw := v.([]interface{})
		if len(configRaw) == 0 {
			err = fmt.Errorf("the schema cannot be empty if type is 400(Table type)")
			return
		}
		for i, item := range configRaw {
			tmp := item.(map[string]interface{})
			schemas = append(schemas, dataset.Field{
				SchemaId: i + 1,
				Name:     tmp["name"].(string),
				Type:     tmp["type"].(string),
			})
		}
	}
	return
}

func parseDatasetErrorToError404(respErr error) error {
	var apiError dataset.CreateResp
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && (apiError.ErrorCode == "ModelArts.4352") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}

func setDataSourcesToState(d *schema.ResourceData, ds dataset.DataSource) error {
	result := make([]interface{}, 1)
	item := map[string]interface{}{
		"data_type":          ds.DataType,
		"path":               ds.DataPath,
		"with_column_header": ds.WithColumnHeader,
	}
	// the API lost some info: queue_name,database_name,table_name,user_name,password,cluster_id,input
	if ds.DataType == 4 {
		item["path"] = ds.SourceInfo.Input
	}

	result[0] = item
	return d.Set("data_source", result)
}

func setLabelsToState(d *schema.ResourceData, labels []dataset.Label) error {
	result := make([]interface{}, len(labels))
	for i, v := range labels {
		result[i] = map[string]interface{}{
			"name":              v.Name,
			"property_color":    v.Property.Color,
			"property_shape":    v.Property.DefaultShape,
			"property_shortcut": v.Property.Shortcut,
		}
	}
	return d.Set("labels", result)
}

func setSchemaToState(d *schema.ResourceData, in []dataset.Field) error {
	result := make([]interface{}, len(in))
	for i, v := range in {
		result[i] = map[string]interface{}{
			"name": v.Name,
			"type": v.Type,
		}
	}
	return d.Set("schemas", result)
}

package dms

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceDmsKafkaSmartConnectTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaSmartConnectTaskCreate,
		ReadContext:   resourceDmsKafkaSmartConnectTaskRead,
		DeleteContext: resourceDmsKafkaSmartConnectTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaSmartConnectTaskImportState,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"connector_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the connector id of the kafka instance.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the source type of the smart connect task.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the name of the smart connect task.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the destination type of the smart connect task.`,
			},
			"obs_destination_descriptor": {
				Type:        schema.TypeList,
				Elem:        obsDestinationDescriptorSchema(),
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the destination parameters of the smart connect task.`,
			},
			"topics": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the topics of the smart connect task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the smart connect task.`,
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the creation time of the smart connect task.`,
			},
		},
	}
}

func obsDestinationDescriptorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the access key of the smart connect task.`,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the secret key of the smart connect task.`,
			},
			"consumer_strategy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the consumer strategy of the smart connect task.`,
			},
			"deliver_time_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the deliver time interval of the smart connect task.`,
			},
			"obs_bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the obs bucket name of the smart connect task.`,
			},
			"partition_format": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the partiton format of the smart connect task.`,
			},
			"obs_path": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the obs path of the smart connect task.`,
			},
			"destination_file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "TEXT",
				Description: `Specifies the destination file type of the smart connect task.`,
			},
			"record_delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "\n",
				Description: `Specifies the record delimiter of the smart connect task.`,
			},
			"topics": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the topics of the task.`,
			},
			"topics_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the topics regular expression of the smart connect task.`,
			},
		},
	}
	return &sc
}

func resourceDmsKafkaSmartConnectTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	// createKafkaSmartConnectTask: create DMS kafka smart connect task
	var (
		createKafkaSmartConnectTaskHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks"
		createKafkaSmartConnectTaskProduct = "dms"
	)
	createKafkaSmartConnectTaskClient, err := cfg.NewServiceClient(createKafkaSmartConnectTaskProduct, region)

	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	connectorID := d.Get("connector_id").(string)
	createKafkaSmartConnectTaskPath := createKafkaSmartConnectTaskClient.Endpoint + createKafkaSmartConnectTaskHttpUrl
	createKafkaSmartConnectTaskPath = strings.ReplaceAll(createKafkaSmartConnectTaskPath, "{project_id}",
		createKafkaSmartConnectTaskClient.ProjectID)
	createKafkaSmartConnectTaskPath = strings.ReplaceAll(createKafkaSmartConnectTaskPath, "{connector_id}", connectorID)

	createKafkaSmartConnectTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createKafkaSmartConnectTaskOpt.JSONBody = utils.RemoveNil(buildCreateKafkaSmartConnectTaskBodyParams(d))

	createKafkaSmartConnectTaskResp, createErr := createKafkaSmartConnectTaskClient.Request("POST",
		createKafkaSmartConnectTaskPath, &createKafkaSmartConnectTaskOpt)

	if createErr != nil {
		return diag.Errorf("error creating DMS kafka smart connect task: %v", createErr)
	}

	kafkaSmartConnectTaskRespBody, err := utils.FlattenResponse(createKafkaSmartConnectTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("task_id", kafkaSmartConnectTaskRespBody, nil)
	d.SetId(taskID.(string))
	return resourceDmsKafkaSmartConnectTaskRead(ctx, d, meta)
}

func buildCreateKafkaSmartConnectTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	obsDestinationDescriptorStruct, err := buildObsDestinationDescriptorStruct(d)
	if err != nil {
		log.Printf("[DEBUG] obsDestinationDescriptorConfig: %v", err)
	}

	bodyParams := map[string]interface{}{
		"source_type":                d.Get("source_type"),
		"task_name":                  d.Get("task_name"),
		"destination_type":           d.Get("destination_type"),
		"obs_destination_descriptor": obsDestinationDescriptorStruct,
	}
	return bodyParams
}

// obs_destination_descriptor struct
type ObsDestinationDescriptorStruct struct {
	Topics              string `json:"topics,omitempty"`
	TopicsRegex         string `json:"topics_regex,omitempty"`
	ConsumerStrategy    string `json:"consumer_strategy,omitempty"`
	DestinationFileType string `json:"destination_file_type"`
	AccessKey           string `json:"access_key,omitempty"`
	SecretKey           string `json:"secret_key,omitempty"`
	ObsBucketName       string `json:"obs_bucket_name,omitempty"`
	ObsPath             string `json:"obs_path,omitempty"`
	PartitionFormat     string `json:"partition_format,omitempty"`
	RecordDelimiter     string `json:"record_delimiter,omitempty"`
	DeliverTimeInterval int    `json:"deliver_time_interval"`
}

func buildObsDestinationDescriptorStruct(d *schema.ResourceData) (*ObsDestinationDescriptorStruct, diag.Diagnostics) {
	var obsDestinationDescriptorStruct *ObsDestinationDescriptorStruct
	obsDestinationDescriptorRaw := d.Get("obs_destination_descriptor").([]interface{})
	if len(obsDestinationDescriptorRaw) == 1 {
		if v, ok := obsDestinationDescriptorRaw[0].(map[string]interface{}); ok {
			obsDestinationDescriptorStruct = &ObsDestinationDescriptorStruct{
				ConsumerStrategy:    v["consumer_strategy"].(string),
				DestinationFileType: v["destination_file_type"].(string),
				AccessKey:           v["access_key"].(string),
				SecretKey:           v["secret_key"].(string),
				ObsBucketName:       v["obs_bucket_name"].(string),
				ObsPath:             v["obs_path"].(string),
				PartitionFormat:     v["partition_format"].(string),
				RecordDelimiter:     v["record_delimiter"].(string),
				DeliverTimeInterval: v["deliver_time_interval"].(int),
			}
			// one of topics and topics_regex is required
			topicsRaw := v["topics"].(string)
			topicsRegexRaw := v["topics_regex"].(string)
			if topicsRaw == "" && topicsRegexRaw == "" {
				return nil, diag.Errorf("error building destination descriptor config: topics or topics_regex is required")
			}
			if v["topics"] != "" {
				obsDestinationDescriptorStruct.Topics = v["topics"].(string)
			} else {
				obsDestinationDescriptorStruct.TopicsRegex = v["topics_regex"].(string)
			}
		}
	}
	return obsDestinationDescriptorStruct, nil
}

func resourceDmsKafkaSmartConnectTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getKafkaSmartConnectTask: query DMS kafka smart connect task
	var (
		getKafkaSmartConnectTaskHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks/{task_id}"
		getKafkaSmartConnectTaskProduct = "dms"
	)
	getKafkaSmartConnectTaskClient, err := cfg.NewServiceClient(getKafkaSmartConnectTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	connectorID := d.Get("connector_id").(string)
	getKafkaSmartConnectTaskPath := getKafkaSmartConnectTaskClient.Endpoint + getKafkaSmartConnectTaskHttpUrl
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{project_id}",
		getKafkaSmartConnectTaskClient.ProjectID)
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{connector_id}", connectorID)
	getKafkaSmartConnectTaskPath = strings.ReplaceAll(getKafkaSmartConnectTaskPath, "{task_id}", d.Id())

	getKafkaSmartConnectTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getKafkaSmartConnectTaskResp, err := getKafkaSmartConnectTaskClient.Request("GET", getKafkaSmartConnectTaskPath,
		&getKafkaSmartConnectTaskOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka smart connect task")
	}

	getKafkaSmartConnectTaskRespBody, respBodyerr := utils.FlattenResponse(getKafkaSmartConnectTaskResp)
	if respBodyerr != nil {
		return diag.FromErr(respBodyerr)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("connector_id", connectorID),
		d.Set("task_name", utils.PathSearch("task_name", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("destination_type", utils.PathSearch("destination_type", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("created_at", utils.PathSearch("create_time", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("topics", utils.PathSearch("topics", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("obs_destination_descriptor", flattenObsDestinationDescriptor(
			utils.PathSearch("obs_destination_descriptor", getKafkaSmartConnectTaskRespBody, nil).(map[string]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenObsDestinationDescriptor(obsDestinationDescriptorStruct map[string]interface{}) []map[string]interface{} {
	var obsDestinationDescriptorList []map[string]interface{}
	if obsDestinationDescriptorStruct != nil {
		obsDestinationDescriptorList = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["topics"] = utils.PathSearch("topics", obsDestinationDescriptorStruct, nil)
		params["topics_regex"] = utils.PathSearch("topics_regex", obsDestinationDescriptorStruct, nil)
		params["consumer_strategy"] = utils.PathSearch("consumer_strategy", obsDestinationDescriptorStruct, nil)
		params["destination_file_type"] = utils.PathSearch("destination_file_type", obsDestinationDescriptorStruct, nil)
		params["obs_bucket_name"] = utils.PathSearch("obs_bucket_name", obsDestinationDescriptorStruct, nil)
		params["obs_path"] = utils.PathSearch("obs_path", obsDestinationDescriptorStruct, nil)
		params["partition_format"] = utils.PathSearch("partition_format", obsDestinationDescriptorStruct, nil)
		params["record_delimiter"] = utils.PathSearch("record_delimiter", obsDestinationDescriptorStruct, nil)
		params["deliver_time_interval"] = utils.PathSearch("deliver_time_interval", obsDestinationDescriptorStruct, nil).(float64)
		obsDestinationDescriptorList[0] = params
	}
	return obsDestinationDescriptorList
}

func resourceDmsKafkaSmartConnectTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteKafkaSmartConnectTask: delete DMS kafka smart connect task
	var (
		deleteKafkaSmartConnectTaskHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks/{task_id}"
		deleteKafkaSmartConnectTaskProduct = "dms"
	)
	deleteKafkaSmartConnectTaskClient, err := cfg.NewServiceClient(deleteKafkaSmartConnectTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS Client: %s", err)
	}

	connectorID := d.Get("connector_id").(string)
	deleteKafkaSmartConnectTaskPath := deleteKafkaSmartConnectTaskClient.Endpoint + deleteKafkaSmartConnectTaskHttpUrl
	deleteKafkaSmartConnectTaskPath = strings.ReplaceAll(deleteKafkaSmartConnectTaskPath, "{project_id}",
		deleteKafkaSmartConnectTaskClient.ProjectID)
	deleteKafkaSmartConnectTaskPath = strings.ReplaceAll(deleteKafkaSmartConnectTaskPath, "{connector_id}", connectorID)
	deleteKafkaSmartConnectTaskPath = strings.ReplaceAll(deleteKafkaSmartConnectTaskPath, "{task_id}", d.Id())

	deleteKafkaSmartConnectTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, deleteErr := deleteKafkaSmartConnectTaskClient.Request("DELETE",
		deleteKafkaSmartConnectTaskPath, &deleteKafkaSmartConnectTaskOpt)

	if deleteErr != nil {
		return diag.Errorf("error deleting DMS kafka smart connect task: %v", err)
	}

	d.SetId("")

	return nil
}

// resourceDmsKafkaSmartConnectTaskImportState is used to import an id with format <connector_id>/<id>
func resourceDmsKafkaSmartConnectTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <connector_id>/<id>")
	}

	d.Set("connector_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

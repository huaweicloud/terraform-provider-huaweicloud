package kafka

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka POST /v2/{project_id}/connectors/{connector_id}/sink-tasks
// @API Kafka GET /v2/{project_id}/connectors/{connector_id}/sink-tasks
// @API Kafka DELETE /v2/{project_id}/connectors/{connector_id}/sink-tasks/{task_id}
func ResourceDmsKafkaSmartConnectTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaSmartConnectTaskCreate,
		ReadContext:   resourceDmsKafkaSmartConnectTaskRead,
		DeleteContext: resourceDmsKafkaSmartConnectTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaSmartConnectTaskImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
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
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the smart connect task.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the destination type of the smart connect task.`,
			},
			"access_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the access key used to access the OBS bucket.`,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `Specifies the secret access key used to access the OBS bucket.`,
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
				ForceNew:    true,
				Description: `Specifies the destination file type of the smart connect task.`,
			},
			"record_delimiter": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the record delimiter of the smart connect task.`,
			},
			"topics": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  `Specifies the topics of the task.`,
				ExactlyOneOf: []string{"topics", "topics_regex"},
			},
			"topics_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the topics regular expression of the smart connect task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the smart connect task.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time of the smart connect task.`,
			},
		},
	}
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

	createKafkaSmartConnectTaskResp, err := createKafkaSmartConnectTaskClient.Request("POST",
		createKafkaSmartConnectTaskPath, &createKafkaSmartConnectTaskOpt)

	if err != nil {
		return diag.Errorf("error creating DMS kafka smart connect task: %v", err)
	}

	kafkaSmartConnectTaskRespBody, err := utils.FlattenResponse(createKafkaSmartConnectTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("task_id", kafkaSmartConnectTaskRespBody, nil)
	if taskID == nil {
		return diag.Errorf("error retrieving DMS kafka smart connect task id: the task id is nil.")
	}
	d.SetId(taskID.(string))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING"},
		Refresh:      smartConnectTaskStateRefreshFunc(createKafkaSmartConnectTaskClient, connectorID, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the smart connect task (%s) to be done: %s", d.Id(), err)
	}

	return resourceDmsKafkaSmartConnectTaskRead(ctx, d, meta)
}

func buildCreateKafkaSmartConnectTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_type":                d.Get("source_type"),
		"task_name":                  d.Get("task_name"),
		"destination_type":           d.Get("destination_type"),
		"obs_destination_descriptor": buildObsDestinationDescriptorStruct(d),
	}
	return bodyParams
}

func buildObsDestinationDescriptorStruct(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"consumer_strategy":     d.Get("consumer_strategy"),
		"destination_file_type": d.Get("destination_file_type"),
		"access_key":            d.Get("access_key"),
		"secret_key":            d.Get("secret_key"),
		"obs_bucket_name":       d.Get("obs_bucket_name"),
		"obs_path":              d.Get("obs_path"),
		"partition_format":      d.Get("partition_format"),
		"record_delimiter":      d.Get("record_delimiter"),
		"deliver_time_interval": d.Get("deliver_time_interval"),
		"topics":                d.Get("topics"),
		"topics_regex":          d.Get("topics_regex"),
	}
	return rst
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

	getKafkaSmartConnectTaskRespBody, err := utils.FlattenResponse(getKafkaSmartConnectTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	obsDestinationDescriptorStruct := utils.PathSearch("obs_destination_descriptor", getKafkaSmartConnectTaskRespBody, nil)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("connector_id", connectorID),
		d.Set("task_name", utils.PathSearch("task_name", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("destination_type", utils.PathSearch("destination_type", getKafkaSmartConnectTaskRespBody, nil)),
		d.Set("topics", utils.PathSearch("topics", obsDestinationDescriptorStruct, nil)),
		d.Set("topics_regex", utils.PathSearch("topics_regex", obsDestinationDescriptorStruct, nil)),
		d.Set("consumer_strategy", utils.PathSearch("consumer_strategy", obsDestinationDescriptorStruct, nil)),
		d.Set("destination_file_type", utils.PathSearch("destination_file_type", obsDestinationDescriptorStruct, nil)),
		d.Set("obs_bucket_name", utils.PathSearch("obs_bucket_name", obsDestinationDescriptorStruct, nil)),
		d.Set("obs_path", utils.PathSearch("obs_path", obsDestinationDescriptorStruct, nil)),
		d.Set("partition_format", utils.PathSearch("partition_format", obsDestinationDescriptorStruct, nil)),
		d.Set("record_delimiter", utils.PathSearch("record_delimiter", obsDestinationDescriptorStruct, nil)),
		d.Set("deliver_time_interval", utils.PathSearch("deliver_time_interval", obsDestinationDescriptorStruct, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getKafkaSmartConnectTaskRespBody, float64(0)).(float64)/1000), false)),
		d.Set("status", utils.PathSearch("status", getKafkaSmartConnectTaskRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDmsKafkaSmartConnectTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, err = deleteKafkaSmartConnectTaskClient.Request("DELETE", deleteKafkaSmartConnectTaskPath, &deleteKafkaSmartConnectTaskOpt)

	if err != nil {
		return diag.Errorf("error deleting DMS kafka smart connect task: %v", err)
	}

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

func smartConnectTaskStateRefreshFunc(client *golangsdk.ServiceClient, connectorID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// getSmartConnectTask: query smart connect task
		var (
			getSmartConnectTaskHttpUrl = "v2/{project_id}/connectors/{connector_id}/sink-tasks/{task_id}"
		)

		getSmartConnectTaskPath := client.Endpoint + getSmartConnectTaskHttpUrl
		getSmartConnectTaskPath = strings.ReplaceAll(getSmartConnectTaskPath, "{project_id}", client.ProjectID)
		getSmartConnectTaskPath = strings.ReplaceAll(getSmartConnectTaskPath, "{connector_id}", connectorID)
		getSmartConnectTaskPath = strings.ReplaceAll(getSmartConnectTaskPath, "{task_id}", taskID)

		getSmartConnectTaskOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		getSmartConnectTaskResp, err := client.Request("GET", getSmartConnectTaskPath, &getSmartConnectTaskOpt)
		if err != nil {
			return nil, "QUERY ERROR", err
		}

		smartConnectTaskRespBody, err := utils.FlattenResponse(getSmartConnectTaskResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		status := utils.PathSearch("status", smartConnectTaskRespBody, "").(string)
		return smartConnectTaskRespBody, status, nil
	}
}

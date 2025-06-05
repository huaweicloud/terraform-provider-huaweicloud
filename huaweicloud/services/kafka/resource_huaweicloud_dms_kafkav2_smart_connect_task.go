package kafka

import (
	"context"
	"fmt"
	"reflect"
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

// @API Kafka POST /v2/{project_id}/instances/{instance_id}/connector/tasks
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}
// @API Kafka DELETE /v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}
func ResourceDmsKafkav2SmartConnectTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkav2SmartConnectTaskCreate,
		ReadContext:   resourceDmsKafkav2SmartConnectTaskRead,
		DeleteContext: resourceDmsKafkav2SmartConnectTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkav2SmartConnectTaskImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the kafka instance ID.`,
			},
			"task_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the smart connect task name.`,
			},
			"topics": {
				Type:         schema.TypeSet,
				Optional:     true,
				ForceNew:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
				Description:  `Specifies the topics list of the task.`,
				ExactlyOneOf: []string{"topics", "topics_regex"},
			},
			"topics_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the topics regular expression of the smart connect task.`,
			},
			"start_later": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to start a task later.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the source type of the smart connect task.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the destination type of the smart connect task.`,
			},

			// Source task configuration
			"source_task": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Specifies the source configuration of a smart connect task.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// connection for kafka to kafka
						"current_instance_alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the current Kafka instance alias.`,
						},
						"peer_instance_alias": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the peer Kafka instance alias.`,
						},
						"peer_instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the peer Kafka instance ID.`,
						},
						"peer_instance_address": {
							Type:        schema.TypeSet,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the peer Kafka instance address.`,
						},
						"security_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the peer Kafka authentication.`,
						},
						"sasl_mechanism": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the peer Kafka authentication mode.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the peer Kafka username.`,
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: `Specifies the peer Kafka user password.`,
						},

						// task configuration
						"direction": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the sync direction.`,
						},
						"sync_consumer_offsets_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies whether to sync the consumption progress.`,
						},
						"replication_factor": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the number of topic replicas.`,
						},
						"task_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the number of data replication tasks.`,
						},
						"rename_topic_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies whether to rename the topic.`,
						},
						"provenance_header_enabled": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies whether the message header contains the message source.`,
						},
						"consumer_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the start offset.`,
						},
						"compression_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							ForceNew:    true,
							Description: `Specifies the compression algorithm to use for copying messages.`,
						},
						"topics_mapping": {
							Type:        schema.TypeSet,
							Optional:    true,
							ForceNew:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `Specifies the topic mapping, which is used to customize the target topic name.`,
						},
					},
				},
			},

			"destination_task": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: `Specifies the source configuration of a smart connect task.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Kafka to OBS
						"access_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the access key used to access the OBS bucket.`,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Sensitive:   true,
							Description: `Specifies the secret access key used to access the OBS bucket.`,
						},
						"agency_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "schema: Internal",
						},
						"consumer_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the consumer strategy of the smart connect task.`,
						},
						"deliver_time_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the deliver time interval of the smart connect task.`,
						},
						"obs_bucket_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies the obs bucket name of the smart connect task.`,
						},
						"partition_format": {
							Type:        schema.TypeString,
							Optional:    true,
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
						"store_keys": {
							Type:        schema.TypeBool,
							Optional:    true,
							ForceNew:    true,
							Description: `Specifies whether to store keys.`,
						},
						"obs_part_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the size of each file to be uploaded.`,
						},
						"flush_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the flush size.`,
						},
						"timezone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time zone.`,
						},
						"schema_generator_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the schema generator class.`,
						},
						"partitioner_class": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the partitioner class.`,
						},
						"key_converter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the key converter.`,
						},
						"value_converter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the value converter.`,
						},
						"kv_delimiter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the kv delimiter.`,
						},
					},
				},
			},

			// Computed
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

func resourceDmsKafkav2SmartConnectTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/connector/tasks"
	createTaskPath := client.Endpoint + createTaskHttpUrl
	createTaskPath = strings.ReplaceAll(createTaskPath, "{project_id}", client.ProjectID)
	createTaskPath = strings.ReplaceAll(createTaskPath, "{instance_id}", instanceID)
	createTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateKafkav2SmartConnectTaskBodyParams(d)),
	}

	createTaskResp, err := client.Request("POST", createTaskPath, &createTaskOpt)
	if err != nil {
		return diag.Errorf("error creating DMS kafka smart connect task: %v", err)
	}
	createTaskRespBody, err := utils.FlattenResponse(createTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("id", createTaskRespBody, nil)
	if taskID == nil {
		return diag.Errorf("error retrieving DMS kafka smart connect task ID: id is nil.")
	}
	d.SetId(taskID.(string))

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"CREATING"},
		Target:       []string{"RUNNING", "WAITING"},
		Refresh:      kafkav2SmartConnectTaskStateRefreshFunc(client, instanceID, d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the smart connect task (%s) to be done: %s", d.Id(), err)
	}

	return resourceDmsKafkav2SmartConnectTaskRead(ctx, d, meta)
}

func buildCreateKafkav2SmartConnectTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"task_name":    d.Get("task_name"),
		"topics":       utils.ValueIgnoreEmpty(changeListToStringWithCommasSplit(d.Get("topics").(*schema.Set).List())),
		"topics_regex": utils.ValueIgnoreEmpty(d.Get("topics_regex")),
		"start_later":  utils.ValueIgnoreEmpty(d.Get("start_later")),
		"source_type":  utils.ValueIgnoreEmpty(d.Get("source_type")),
		"sink_type":    utils.ValueIgnoreEmpty(d.Get("destination_type")),
		"source_task":  utils.ValueIgnoreEmpty(buildSourceTaskRequestBody(d.Get("source_task").([]interface{}))),
		"sink_task":    utils.ValueIgnoreEmpty(buildSinkTaskRequestBody(d.Get("destination_task").([]interface{}))),
	}
	return bodyParams
}

func buildSourceTaskRequestBody(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := map[string]interface{}{
		"current_cluster_name": utils.ValueIgnoreEmpty(params["current_instance_alias"]),
		"cluster_name":         utils.ValueIgnoreEmpty(params["peer_instance_alias"]),
		"instance_id":          utils.ValueIgnoreEmpty(params["peer_instance_id"]),
		"bootstrap_servers": utils.ValueIgnoreEmpty(changeListToStringWithCommasSplit(
			params["peer_instance_address"].(*schema.Set).List())),
		"security_protocol":             utils.ValueIgnoreEmpty(params["security_protocol"]),
		"user_name":                     utils.ValueIgnoreEmpty(params["user_name"]),
		"password":                      utils.ValueIgnoreEmpty(params["password"]),
		"sasl_mechanism":                utils.ValueIgnoreEmpty(params["sasl_mechanism"]),
		"direction":                     utils.ValueIgnoreEmpty(params["direction"]),
		"sync_consumer_offsets_enabled": utils.ValueIgnoreEmpty(params["sync_consumer_offsets_enabled"]),
		"replication_factor":            utils.ValueIgnoreEmpty(params["replication_factor"]),
		"task_num":                      utils.ValueIgnoreEmpty(params["task_num"]),
		"rename_topic_enabled":          utils.ValueIgnoreEmpty(params["rename_topic_enabled"]),
		"provenance_header_enabled":     utils.ValueIgnoreEmpty(params["provenance_header_enabled"]),
		"consumer_strategy":             utils.ValueIgnoreEmpty(params["consumer_strategy"]),
		"compression_type":              utils.ValueIgnoreEmpty(params["compression_type"]),
		"topics_mapping": utils.ValueIgnoreEmpty(changeListToStringWithCommasSplit(
			params["topics_mapping"].(*schema.Set).List())),
	}
	return rst
}

func changeListToStringWithCommasSplit(params []interface{}) string {
	strArray := make([]string, 0, len(params))
	for _, param := range params {
		strArray = append(strArray, param.(string))
	}
	return strings.Join(strArray, ",")
}

func buildSinkTaskRequestBody(rawParams []interface{}) map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	params := rawParams[0].(map[string]interface{})
	rst := map[string]interface{}{
		"consumer_strategy":     utils.ValueIgnoreEmpty(params["consumer_strategy"]),
		"access_key":            utils.ValueIgnoreEmpty(params["access_key"]),
		"secret_key":            utils.ValueIgnoreEmpty(params["secret_key"]),
		"agency_name":           utils.ValueIgnoreEmpty(params["agency_name"]),
		"obs_bucket_name":       utils.ValueIgnoreEmpty(params["obs_bucket_name"]),
		"partition_format":      utils.ValueIgnoreEmpty(params["partition_format"]),
		"deliver_time_interval": utils.ValueIgnoreEmpty(params["deliver_time_interval"]),
		"obs_path":              utils.ValueIgnoreEmpty(params["obs_path"]),
		"record_delimiter":      utils.ValueIgnoreEmpty(params["record_delimiter"]),
		"destination_file_type": utils.ValueIgnoreEmpty(params["destination_file_type"]),
		"store_keys":            utils.ValueIgnoreEmpty(params["store_keys"]),
	}

	return rst
}

func resourceDmsKafkav2SmartConnectTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	getTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}"
	getTaskPath := client.Endpoint + getTaskHttpUrl
	getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", client.ProjectID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{instance_id}", instanceID)
	getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", d.Id())
	getTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getTaskResp, err := client.Request("GET", getTaskPath, &getTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DMS kafka smart connect task")
	}
	getTaskRespBody, err := utils.FlattenResponse(getTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// in order to be compatible with old API, does not set source_type and destination_type
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("task_name", utils.PathSearch("task_name", getTaskRespBody, nil)),
		d.Set("topics", flattenStringWithCommaSplitToSlice(utils.PathSearch("topics", getTaskRespBody, "").(string))),
		d.Set("topics_regex", utils.PathSearch("topics_regex", getTaskRespBody, nil)),
		d.Set("source_type", utils.PathSearch("source_type", getTaskRespBody, nil)),
		d.Set("destination_type", utils.PathSearch("sink_type", getTaskRespBody, nil)),
		d.Set("source_task", flattenSourceTaskResponse(d, utils.PathSearch("source_task", getTaskRespBody, nil))),
		d.Set("destination_task", flattenSinkTaskResponse(d, utils.PathSearch("sink_task", getTaskRespBody, nil))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getTaskRespBody, float64(0)).(float64)/1000), false)),
		d.Set("status", utils.PathSearch("status", getTaskRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSourceTaskResponse(d *schema.ResourceData, rawParams interface{}) []map[string]interface{} {
	if reflect.DeepEqual(rawParams, map[string]interface{}{}) {
		return nil
	}
	rst := make([]map[string]interface{}, 1)
	params := map[string]interface{}{
		"current_instance_alias":        utils.PathSearch("current_cluster_name", rawParams, nil),
		"peer_instance_alias":           utils.PathSearch("cluster_name", rawParams, nil),
		"peer_instance_id":              utils.PathSearch("instance_id", rawParams, nil),
		"peer_instance_address":         strings.Split(utils.PathSearch("bootstrap_servers", rawParams, "").(string), ","),
		"security_protocol":             utils.PathSearch("security_protocol", rawParams, nil),
		"user_name":                     utils.PathSearch("user_name", rawParams, nil),
		"password":                      d.Get("source_task.0.password"),
		"sasl_mechanism":                utils.PathSearch("sasl_mechanism", rawParams, nil),
		"direction":                     utils.PathSearch("direction", rawParams, nil),
		"sync_consumer_offsets_enabled": utils.PathSearch("sync_consumer_offsets_enabled", rawParams, false),
		"replication_factor":            utils.PathSearch("replication_factor", rawParams, 0),
		"task_num":                      utils.PathSearch("task_num", rawParams, 0),
		"rename_topic_enabled":          utils.PathSearch("rename_topic_enabled", rawParams, false),
		"provenance_header_enabled":     utils.PathSearch("provenance_header_enabled", rawParams, false),
		"consumer_strategy":             utils.PathSearch("consumer_strategy", rawParams, nil),
		"compression_type":              utils.PathSearch("compression_type", rawParams, nil),
		"topics_mapping":                flattenStringWithCommaSplitToSlice(utils.PathSearch("topics_mapping", rawParams, "").(string)),
	}
	rst[0] = params
	return rst
}

func flattenSinkTaskResponse(d *schema.ResourceData, rawParams interface{}) []map[string]interface{} {
	if reflect.DeepEqual(rawParams, map[string]interface{}{}) {
		return nil
	}
	rst := make([]map[string]interface{}, 1)
	params := map[string]interface{}{
		"access_key":             d.Get("destination_task.0.access_key"),
		"secret_key":             d.Get("destination_task.0.secret_key"),
		"consumer_strategy":      utils.PathSearch("consumer_strategy", rawParams, nil),
		"destination_file_type":  utils.PathSearch("destination_file_type", rawParams, nil),
		"obs_bucket_name":        utils.PathSearch("obs_bucket_name", rawParams, nil),
		"obs_path":               utils.PathSearch("obs_path", rawParams, nil),
		"partition_format":       utils.PathSearch("partition_format", rawParams, nil),
		"record_delimiter":       utils.PathSearch("record_delimiter", rawParams, nil),
		"deliver_time_interval":  utils.PathSearch("deliver_time_interval", rawParams, 0),
		"store_keys":             utils.PathSearch("store_keys", rawParams, false),
		"obs_part_size":          utils.PathSearch("obs_part_size", rawParams, 0),
		"flush_size":             utils.PathSearch("flush_size", rawParams, 0),
		"timezone":               utils.PathSearch("timezone", rawParams, nil),
		"schema_generator_class": utils.PathSearch("schema_generator_class", rawParams, nil),
		"partitioner_class":      utils.PathSearch("partitioner_class", rawParams, nil),
		"key_converter":          utils.PathSearch("key_converter", rawParams, nil),
		"value_converter":        utils.PathSearch("value_converter", rawParams, nil),
		"kv_delimiter":           utils.PathSearch("kv_delimiter", rawParams, nil),
	}

	rst[0] = params
	return rst
}

func flattenStringWithCommaSplitToSlice(s string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, ",")
}

func resourceDmsKafkav2SmartConnectTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	deleteTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}"
	deleteTaskPath := client.Endpoint + deleteTaskHttpUrl
	deleteTaskPath = strings.ReplaceAll(deleteTaskPath, "{project_id}", client.ProjectID)
	deleteTaskPath = strings.ReplaceAll(deleteTaskPath, "{instance_id}", instanceID)
	deleteTaskPath = strings.ReplaceAll(deleteTaskPath, "{task_id}", d.Id())
	deleteTaskOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deleteTaskPath, &deleteTaskOpt)
	if err != nil {
		return diag.Errorf("error deleting DMS kafka smart connect task: %v", err)
	}

	return nil
}

// resourceDmsKafkav2SmartConnectTaskImportState is used to import an id with format <instance_id>/<task_id>
func resourceDmsKafkav2SmartConnectTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<task_id>")
	}

	d.Set("instance_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}

func kafkav2SmartConnectTaskStateRefreshFunc(client *golangsdk.ServiceClient, instanceID, taskID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getTaskHttpUrl := "v2/{project_id}/instances/{instance_id}/connector/tasks/{task_id}"
		getTaskPath := client.Endpoint + getTaskHttpUrl
		getTaskPath = strings.ReplaceAll(getTaskPath, "{project_id}", client.ProjectID)
		getTaskPath = strings.ReplaceAll(getTaskPath, "{instance_id}", instanceID)
		getTaskPath = strings.ReplaceAll(getTaskPath, "{task_id}", taskID)
		getTaskOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		getTaskResp, err := client.Request("GET", getTaskPath, &getTaskOpt)
		if err != nil {
			return nil, "QUERY ERROR", err
		}
		getTaskRespBody, err := utils.FlattenResponse(getTaskResp)
		if err != nil {
			return nil, "PARSE ERROR", err
		}

		status := utils.PathSearch("status", getTaskRespBody, "").(string)
		return getTaskRespBody, status, nil
	}
}

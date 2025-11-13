package rocketmq

import (
	"context"
	"encoding/json"
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

// @API RocketMQ POST /v2/{project_id}/instances/{instance_id}/metadata
// @API RocketMQ GET /v2/{project_id}/instances/{instance_id}/metadata
// @API RocketMQ DELETE /v2/{project_id}/instances/{instance_id}/metadata
func ResourceDmsRocketmqMigrationTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRocketmqMigrationTaskCreate,
		ReadContext:   resourceRocketmqMigrationTaskRead,
		DeleteContext: resourceRocketmqMigrationTaskDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRocketmqMigrationTaskImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(50 * time.Minute),
			Delete: schema.DefaultTimeout(15 * time.Minute),
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
				Description: `Specifies the ID of the RocketMQ instance.`,
			},
			"overwrite": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies whether to overwrite configurations with the same name.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the migration task.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the migration task type.`,
			},
			"topic_configs": {
				Type:        schema.TypeList,
				Elem:        topicSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the topic metadata.`,
			},
			"subscription_groups": {
				Type:        schema.TypeList,
				Elem:        consumerGroupSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the consumer group metadata.`,
			},
			"vhosts": {
				Type:        schema.TypeList,
				Elem:        vhostSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the virtual hosts metadata.`,
			},
			"queues": {
				Type:        schema.TypeList,
				Elem:        queueSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the queue metadata.`,
			},
			"exchanges": {
				Type:        schema.TypeList,
				Elem:        exchangeSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the exchange metadata.`,
			},
			"bindings": {
				Type:        schema.TypeList,
				Elem:        bindingSchema(),
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the binding metadata.`,
			},
			"start_date": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the start time of the migration task.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the status of the migration task.`,
			},
		},
	}
}

func topicSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"topic_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the topic name.`,
			},
			"order": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether a message is an ordered message.`,
			},
			"perm": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of permission.`,
			},
			"read_queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of read queues.`,
			},
			"write_queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of write queues.`,
			},
			"topic_filter_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the filter type of a topic.`,
			},
			"topic_sys_flag": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the system flag of a topic.`,
			},
		},
	}
	return &sc
}

func consumerGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of a consumer group.`,
			},
			"consume_broadcast_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable broadcast.`,
			},
			"consume_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable consumption.`,
			},
			"consume_from_min_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable consumption from the earliest offset.`,
			},
			"notify_consumerids_changed_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to notify changes of consumer IDs.`,
			},
			"retry_max_times": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the maximum number of consumption retries.`,
			},
			"retry_queue_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of retry queues.`,
			},
			"which_broker_when_consume_slow": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the broker selected for slow consumption.`,
			},
		},
	}
	return &sc
}

func vhostSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the virtual host name.`,
			},
		},
	}
	return &sc
}

func queueSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vhost": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the virtual host name.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the queue name.`,
			},
			"durable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable data persistence.`,
			},
		},
	}
	return &sc
}

func exchangeSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vhost": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the virtual host name.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the switch name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the exchange type.`,
			},
			"durable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies whether to enable data persistence.`,
			},
		},
	}
	return &sc
}

func bindingSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"vhost": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the virtual host name.`,
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the message source.`,
			},
			"destination": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the message target.`,
			},
			"destination_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the message target type.`,
			},
			"routing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the routing key.`,
			},
		},
	}
	return &sc
}

func resourceRocketmqMigrationTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createRocketmqMigrationTask: create RocketMQ migration task
	var (
		createRocketmqMigrationTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/metadata"
		createRocketmqMigrationTaskProduct = "dms"
	)
	createRocketmqMigrationTaskClient, err := cfg.NewServiceClient(createRocketmqMigrationTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	createParams := fmt.Sprintf("?name=%v&overwrite=%v&type=%v", d.Get("name"), d.Get("overwrite"), d.Get("type"))
	createRocketmqMigrationTaskPath := createRocketmqMigrationTaskClient.Endpoint + createRocketmqMigrationTaskHttpUrl
	createRocketmqMigrationTaskPath = strings.ReplaceAll(createRocketmqMigrationTaskPath, "{project_id}",
		createRocketmqMigrationTaskClient.ProjectID)
	createRocketmqMigrationTaskPath = strings.ReplaceAll(createRocketmqMigrationTaskPath, "{instance_id}", instanceID)
	createRocketmqMigrationTaskPath += createParams
	createRocketmqMigrationTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	createRocketmqMigrationTaskOpt.JSONBody = utils.RemoveNil(buildCreateMigrationTaskBodyParams(d))

	err = waitForInstanceStatusCompleted(ctx, createRocketmqMigrationTaskClient, instanceID, []string{"RUNNING"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the status of the RocketMQ instance (%s) to be RUNNING: %s", instanceID, err)
	}

	createRocketmqMigrationTaskResp, err := createRocketmqMigrationTaskClient.Request("POST", createRocketmqMigrationTaskPath,
		&createRocketmqMigrationTaskOpt)
	if err != nil {
		return diag.Errorf("error creating migration task: %s", err)
	}
	createRocketmqMigrationTaskRespBody, err := utils.FlattenResponse(createRocketmqMigrationTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskId := utils.PathSearch("task_id", createRocketmqMigrationTaskRespBody, nil)
	if taskId == nil {
		return diag.Errorf("error creating migration task: task ID is not found in API response")
	}

	d.SetId(taskId.(string))

	return resourceRocketmqMigrationTaskRead(ctx, d, meta)
}

func buildCreateMigrationTaskBodyParams(d *schema.ResourceData) map[string]interface{} {
	// RocketMQ to RocketMQ
	if d.Get("type") == "rocketmq" {
		topicConfigList := d.Get("topic_configs").([]interface{})
		subscriptionGroupList := d.Get("subscription_groups").([]interface{})
		bodyParams := map[string]interface{}{
			"topicConfigTable":       buildTopicBodyParams(topicConfigList),
			"subscriptionGroupTable": buildGroupBodyParams(subscriptionGroupList),
		}
		return bodyParams
	}
	// RabbitMQ to RocketMQ
	if d.Get("type") == "rabbitToRocket" {
		bodyParams := map[string]interface{}{
			"vhosts":    d.Get("vhosts"),
			"queues":    d.Get("queues"),
			"exchanges": d.Get("exchanges"),
			"bindings":  d.Get("bindings"),
		}
		return bodyParams
	}
	return nil
}

func buildTopicBodyParams(topicConfigList []interface{}) map[string]interface{} {
	topicConfigTable := map[string]interface{}{}
	// topic metadata
	for _, topicConfig := range topicConfigList {
		topicName := utils.PathSearch("topic_name", topicConfig, "").(string)
		topicConfigTable[topicName] = map[string]interface{}{
			"order":           utils.PathSearch("order", topicConfig, nil),
			"perm":            utils.PathSearch("perm", topicConfig, nil),
			"readQueueNums":   utils.PathSearch("read_queue_num", topicConfig, nil),
			"topicFilterType": utils.PathSearch("topic_filter_type", topicConfig, nil),
			"topicName":       topicName,
			"topicSysFlag":    utils.PathSearch("topic_sys_flag", topicConfig, nil),
			"writeQueueNums":  utils.PathSearch("write_queue_num", topicConfig, nil),
		}
	}
	return topicConfigTable
}

func buildGroupBodyParams(subscriptionGroupList []interface{}) map[string]interface{} {
	subscriptionGroupTable := map[string]interface{}{}
	// group metadata
	for _, subscriptionGroup := range subscriptionGroupList {
		groupName := utils.PathSearch("group_name", subscriptionGroup, "").(string)
		subscriptionGroupTable[groupName] = map[string]interface{}{
			"consumeBroadcastEnable":         utils.PathSearch("consume_broadcast_enable", subscriptionGroup, nil),
			"consumeEnable":                  utils.PathSearch("consume_enable", subscriptionGroup, nil),
			"consumeFromMinEnable":           utils.PathSearch("consume_from_min_enable", subscriptionGroup, nil),
			"notifyConsumerIdsChangedEnable": utils.PathSearch("notify_consumerids_changed_enable", subscriptionGroup, nil),
			"groupName":                      groupName,
			"retryMaxTimes":                  utils.PathSearch("retry_max_times", subscriptionGroup, nil),
			"retryQueueNums":                 utils.PathSearch("retry_queue_num", subscriptionGroup, nil),
			"whichBrokerWhenConsumeSlow":     utils.PathSearch("which_broker_when_consume_slow", subscriptionGroup, nil),
		}
	}
	return subscriptionGroupTable
}

func resourceRocketmqMigrationTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getRocketmqMigrationTask: query RocketMQ migration task
	var (
		getRocketmqMigrationTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/metadata"
		getRocketmqMigrationTaskProduct = "dms"
	)
	getRocketmqMigrationTaskClient, err := cfg.NewServiceClient(getRocketmqMigrationTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}
	instanceID := d.Get("instance_id").(string)
	getRocketmqMigrationTaskPath := getRocketmqMigrationTaskClient.Endpoint + getRocketmqMigrationTaskHttpUrl
	getRocketmqMigrationTaskPath = strings.ReplaceAll(getRocketmqMigrationTaskPath, "{project_id}", getRocketmqMigrationTaskClient.ProjectID)
	getRocketmqMigrationTaskPath = strings.ReplaceAll(getRocketmqMigrationTaskPath, "{instance_id}", instanceID)
	getRocketmqMigrationTaskPath += fmt.Sprintf("?id=%s&type=vhost", d.Id())

	getRocketmqMigrationTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	getRocketmqMigrationTaskResp, err := getRocketmqMigrationTaskClient.Request("GET", getRocketmqMigrationTaskPath, &getRocketmqMigrationTaskOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DMS.00405011"),
			"error retrieving RocketMQ migration task")
	}

	getRocketmqMigrationTaskRespBody, err := utils.FlattenResponse(getRocketmqMigrationTaskResp)
	if err != nil {
		return diag.FromErr(err)
	}
	jsonContent := utils.PathSearch("json_content", getRocketmqMigrationTaskRespBody, "")
	var jsonData interface{}
	err = json.Unmarshal([]byte(jsonContent.(string)), &jsonData)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("name", utils.PathSearch("name", getRocketmqMigrationTaskRespBody, nil)),
		d.Set("start_date", utils.PathSearch("start_date", getRocketmqMigrationTaskRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRocketmqMigrationTaskRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRocketmqMigrationTaskRespBody, nil)),
		d.Set("topic_configs", getTopicConfigs(jsonData)),
		d.Set("subscription_groups", getGroupConfigs(jsonData)),
		d.Set("vhosts", utils.PathSearch("vhosts", jsonData, nil)),
		d.Set("queues", utils.PathSearch("queues", jsonData, nil)),
		d.Set("exchanges", utils.PathSearch("exchanges", jsonData, nil)),
		d.Set("bindings", utils.PathSearch("bindings", jsonData, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getTopicConfigs(jsonData interface{}) []interface{} {
	topicConfigs := utils.PathSearch("topicConfigTable", jsonData, make(map[string]interface{}, 0)).(map[string]interface{})
	if len(topicConfigs) < 1 {
		return nil
	}
	topicList := make([]interface{}, 0, len(topicConfigs))
	for _, topicConfig := range topicConfigs {
		topicList = append(topicList, map[string]interface{}{
			"order":             utils.PathSearch("order", topicConfig, nil),
			"perm":              utils.PathSearch("perm", topicConfig, nil),
			"read_queue_num":    utils.PathSearch("readQueueNums", topicConfig, nil),
			"topic_filter_type": utils.PathSearch("topicFilterType", topicConfig, nil),
			"topic_name":        utils.PathSearch("topicName", topicConfig, nil),
			"topic_sys_flag":    utils.PathSearch("topicSysFlag", topicConfig, nil),
			"write_queue_num":   utils.PathSearch("writeQueueNums", topicConfig, nil),
		})
	}
	return topicList
}

func getGroupConfigs(jsonData interface{}) []interface{} {
	groupConfigs := utils.PathSearch("subscriptionGroupTable", jsonData, make(map[string]interface{}, 0)).(map[string]interface{})
	if len(groupConfigs) < 1 {
		return nil
	}
	groupList := make([]interface{}, 0, len(groupConfigs))
	for _, groupConfig := range groupConfigs {
		groupList = append(groupList, map[string]interface{}{
			"consume_broadcast_enable":          utils.PathSearch("consumeBroadcastEnable", groupConfig, nil),
			"consume_enable":                    utils.PathSearch("consumeEnable", groupConfig, nil),
			"consume_from_min_enable":           utils.PathSearch("consumeFromMinEnable", groupConfig, nil),
			"group_name":                        utils.PathSearch("groupName", groupConfig, nil),
			"notify_consumerids_changed_enable": utils.PathSearch("notifyConsumerIdsChangedEnable", groupConfig, nil),
			"retry_max_times":                   utils.PathSearch("retryMaxTimes", groupConfig, nil),
			"retry_queue_num":                   utils.PathSearch("retryQueueNums", groupConfig, nil),
			"which_broker_when_consume_slow":    utils.PathSearch("whichBrokerWhenConsumeSlow", groupConfig, nil),
		})
	}
	return groupList
}

func resourceRocketmqMigrationTaskDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteRocketmqMigrationTask: delete RocketMQ migration task
	var (
		deleteRocketmqMigrationTaskHttpUrl = "v2/{project_id}/instances/{instance_id}/metadata"
		deleteRocketmqMigrationTaskProduct = "dms"
	)
	deleteRocketmqMigrationTaskClient, err := cfg.NewServiceClient(deleteRocketmqMigrationTaskProduct, region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}
	instanceID := d.Get("instance_id").(string)
	deleteRocketmqMigrationTaskPath := deleteRocketmqMigrationTaskClient.Endpoint + deleteRocketmqMigrationTaskHttpUrl
	deleteRocketmqMigrationTaskPath = strings.ReplaceAll(deleteRocketmqMigrationTaskPath, "{project_id}", deleteRocketmqMigrationTaskClient.ProjectID)
	deleteRocketmqMigrationTaskPath = strings.ReplaceAll(deleteRocketmqMigrationTaskPath, "{instance_id}", instanceID)

	deleteRocketmqMigrationTaskOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	deleteRocketmqMigrationTaskOpt.JSONBody = map[string]interface{}{"task_ids": []string{d.Id()}}

	// DELETED means the RocketMQ instance does not exist.
	err = waitForInstanceStatusCompleted(
		ctx,
		deleteRocketmqMigrationTaskClient,
		instanceID,
		[]string{"RUNNING", "DELETED"},
		d.Timeout(schema.TimeoutDelete),
	)
	if err != nil {
		return diag.Errorf("error waiting for the status of the RocketMQ instance (%s) to be RUNNING: %s", instanceID, err)
	}

	_, err = deleteRocketmqMigrationTaskClient.Request("DELETE", deleteRocketmqMigrationTaskPath, &deleteRocketmqMigrationTaskOpt)
	if err != nil {
		return diag.Errorf("error deleting RocketMQ migration task: %s", err)
	}

	return resourceRocketmqMigrationTaskRead(ctx, d, meta)
}

// resourceRocketmqMigrationTaskImportState is used to import an id with format <instance_id>/<id>
func resourceRocketmqMigrationTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}

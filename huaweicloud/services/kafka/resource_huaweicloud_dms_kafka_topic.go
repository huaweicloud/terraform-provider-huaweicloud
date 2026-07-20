package kafka

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	topicNonUpdatableParams = []string{
		"instance_id",
		"name",
		"replicas",
	}
)

// @API Kafka POST /v2/{project_id}/instances/{instance_id}/topics
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/topics
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/topics
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/topics/delete
func ResourceTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTopicCreate,
		ReadContext:   resourceTopicRead,
		UpdateContext: resourceTopicUpdate,
		DeleteContext: resourceTopicDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceTopicImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(topicNonUpdatableParams),
			func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
				if d.HasChange("partitions") {
					oldValue, newValue := d.GetChange("partitions")
					if oldValue.(int) > newValue.(int) {
						return fmt.Errorf("only support to add partitions")
					}
				}
				return nil
			},
		),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the topic is located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the DMS kafka instance to which the topic belongs.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the topic.`,
			},
			"partitions": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `The partition number.`,
			},

			// Optional parameters.
			"new_partition_brokers": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: `The integers list of brokers for new partitions.`,
			},
			"replicas": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The replica number.`,
			},
			"aging_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `The aging time in hours.`,
			},
			"sync_replication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable synchronous replication.`,
			},
			"sync_flushing": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable synchronous flushing.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the topic.`,
			},
			"configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The configuration name.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The configuration value.`,
						},
					},
				},
				Description: `The other topic configurations.`,
			},
			"policies_only": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether this policy is the default policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The topic type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The topic create time.`,
			},

			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					"Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.",
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildCreateTopicBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"id":        d.Get("name"),
		"partition": d.Get("partitions"),
		// Optional parameters.
		"replication":         utils.ValueIgnoreEmpty(d.Get("replicas")),
		"retention_time":      utils.ValueIgnoreEmpty(d.Get("aging_time")),
		"sync_replication":    d.Get("sync_replication"),
		"sync_message_flush":  d.Get("sync_flushing"),
		"topic_desc":          utils.ValueIgnoreEmpty(d.Get("description")),
		"topic_other_configs": buildTopicConfigs(d.Get("configs").(*schema.Set).List()),
	}
}

func buildTopicConfigs(params []interface{}) []map[string]interface{} {
	if len(params) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(params))
	for _, v := range params {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}

	return result
}

func createTopic(client *golangsdk.ServiceClient, instanceId string, bodyParams map[string]interface{}) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/topics"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(bodyParams),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(createResp)
}

func resourceTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	v, err := createTopic(client, instanceId, buildCreateTopicBodyParams(d))
	if err != nil {
		return diag.Errorf("error creating topic for the kafka instance (%s): %s", instanceId, err)
	}

	// use topic name as the resource ID
	topicName, _ := utils.PathSearch("name", v, "").(string)
	if topicName == "" {
		return diag.Errorf("unable to find the topic name from the API response")
	}

	d.SetId(topicName)

	// wait for topic create complete
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"SUCCESS"},
		Refresh:      kafkaTopicCreateRefreshFunc(client, instanceId, topicName),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for topic (%s) of the kafka instance (%s) to be created: %s", topicName, instanceId, err)
	}

	return resourceTopicRead(ctx, d, meta)
}

func listTopics(client *golangsdk.ServiceClient, instanceId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/topics?limit={limit}"
		// The limit maximum value is 200, default is 50.
		limit  = 200
		offset = 0
		result = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		topics := utils.PathSearch("topics", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, topics...)
		if len(topics) < limit {
			break
		}

		offset += len(topics)
	}

	return result, nil
}

func GetTopicByName(client *golangsdk.ServiceClient, instanceId, topicName string) (interface{}, error) {
	allTopics, err := listTopics(client, instanceId)
	if err != nil {
		return nil, err
	}

	topic := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0]", topicName), allTopics, nil)
	if topic == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/instances/{instance_id}/topics",
				RequestId: "NONE",
				Body:      []byte("the topic not exist"),
			},
		}
	}

	return topic, nil
}

func kafkaTopicCreateRefreshFunc(client *golangsdk.ServiceClient, instanceId, topicName string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		topic, err := GetTopicByName(client, instanceId, topicName)
		if err != nil {
			if errCode, ok := err.(golangsdk.ErrDefault404); ok {
				if reflect.DeepEqual(errCode.Body, []byte("the topic not exist")) {
					return topic, "PENDING", nil
				}
				return false, "QUERY ERROR", err
			}

			return nil, "QUERY ERROR", err
		}

		return topic, "SUCCESS", nil
	}
}

func resourceTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		topicName  = d.Id()
	)
	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topic, err := GetTopicByName(client, instanceId, topicName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving topic (%s) of the instance (%s)", topicName, instanceId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("name", utils.PathSearch("name", topic, nil)),
		d.Set("partitions", utils.PathSearch("partition", topic, nil)),
		// Optional parameters.
		d.Set("replicas", utils.PathSearch("replication", topic, nil)),
		d.Set("aging_time", utils.PathSearch("retention_time", topic, nil)),
		d.Set("sync_replication", utils.PathSearch("sync_replication", topic, nil)),
		d.Set("sync_flushing", utils.PathSearch("sync_message_flush", topic, nil)),
		d.Set("configs", flattenTopicConfigs(utils.PathSearch("topic_other_configs", topic, make([]interface{}, 0)).([]interface{}))),
		d.Set("description", utils.PathSearch("topic_desc", topic, nil)),
		d.Set("policies_only", utils.PathSearch("policiesOnly", topic, nil)),
		d.Set("type", setTopicType(int(utils.PathSearch("topic_type", topic, float64(0)).(float64)))),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("created_at", topic, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func setTopicType(topicValue int) string {
	topicType := "common topic"
	if topicValue == 1 {
		topicType = "system topic"
	}

	return topicType
}

func flattenTopicConfigs(params []interface{}) []map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(params))
	for _, val := range params {
		result = append(result, map[string]interface{}{
			"name":  utils.PathSearch("name", val, nil),
			"value": utils.PathSearch("value", val, nil),
		})
	}
	return result
}

func buildUpdateTopicBodyParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"id": d.Get("name"),
	}

	if d.HasChange("partitions") {
		params["new_partition_numbers"] = d.Get("partitions").(int)
		params["new_partition_brokers"] = utils.ValueIgnoreEmpty(d.Get("new_partition_brokers").(*schema.Set).List())
	}

	if d.HasChange("aging_time") {
		params["retention_time"] = d.Get("aging_time")
	}

	if d.HasChange("sync_replication") {
		params["sync_replication"] = d.Get("sync_replication")
	}

	if d.HasChange("sync_flushing") {
		params["sync_message_flush"] = d.Get("sync_flushing")
	}

	if d.HasChange("description") {
		params["topic_desc"] = d.Get("description")
	}

	if d.HasChange("configs") {
		params["topic_other_configs"] = buildTopicConfigs(d.Get("configs").(*schema.Set).List())
	}

	return map[string]interface{}{
		"topics": []map[string]interface{}{
			params,
		},
	}
}

func updateTopic(client *golangsdk.ServiceClient, instanceId string, bodyParams map[string]interface{}) error {
	httpUrl := "v2/{project_id}/instances/{instance_id}/topics"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(bodyParams),
		OkCodes:          []int{204},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		instanceId := d.Get("instance_id").(string)
		err = updateTopic(client, instanceId, buildUpdateTopicBodyParams(d))
		if err != nil {
			return diag.Errorf("error updating topic (%s) of the kafka instance (%s): %s", d.Id(), instanceId, err)
		}
	}

	return resourceTopicRead(ctx, d, meta)
}

func deleteTopic(client *golangsdk.ServiceClient, instanceId string, topicName string) (interface{}, error) {
	httpUrl := "v2/{project_id}/instances/{instance_id}/topics/delete"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"topics": []string{topicName},
		},
	}

	resp, err := client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceTopicDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topicName := d.Id()
	instanceId := d.Get("instance_id").(string)
	respBody, err := deleteTopic(client, instanceId, topicName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting topic (%s) of the kafka instance (%s)", topicName, instanceId))
	}

	// The API will return success even deleting a non-existent topic.
	success := utils.PathSearch(fmt.Sprintf("topics[?id=='%s']|[0].success", topicName), respBody, false).(bool)
	if !success {
		return diag.Errorf("error deleting topic (%s) of the kafka instance (%s)", topicName, instanceId)
	}

	return nil
}

func resourceTopicImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<topic_name>', but got '%s'", importId)
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}

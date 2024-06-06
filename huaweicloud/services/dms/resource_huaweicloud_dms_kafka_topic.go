package dms

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/topics"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceDmsKafkaTopic implements the resource of "huaweicloud_dms_kafka_topic"
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/topics/delete
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/topics
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/topics
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/topics
func ResourceDmsKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaTopicCreate,
		ReadContext:   resourceDmsKafkaTopicRead,
		UpdateContext: resourceDmsKafkaTopicUpdate,
		DeleteContext: resourceDmsKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaTopicImport,
		},

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			if d.HasChange("partitions") {
				oldValue, newValue := d.GetChange("partitions")
				if oldValue.(int) > newValue.(int) {
					return fmt.Errorf("only support to add partitions")
				}
			}
			return nil
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"partitions": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"replicas": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"aging_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sync_replication": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"sync_flushing": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"new_partition_brokers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configs": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"policies_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dmsV2Client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createOpts := &topics.CreateOps{
		Name:             d.Get("name").(string),
		Partition:        d.Get("partitions").(int),
		Replication:      d.Get("replicas").(int),
		RetentionTime:    d.Get("aging_time").(int),
		SyncReplication:  d.Get("sync_replication").(bool),
		SyncMessageFlush: d.Get("sync_flushing").(bool),
		Description:      d.Get("description").(string),
		Configs:          buildKafkaTopicParameters(d.Get("configs").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	instanceID := d.Get("instance_id").(string)
	v, err := topics.Create(dmsV2Client, instanceID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating DMS kafka topic: %s", err)
	}

	// use topic name as the resource ID
	d.SetId(v.Name)
	return resourceDmsKafkaTopicRead(ctx, d, meta)
}

func buildKafkaTopicParameters(params []interface{}) []topics.ConfigParam {
	paramList := make([]topics.ConfigParam, len(params))
	for i, v := range params {
		paramList[i] = topics.ConfigParam{
			Name:  v.(map[string]interface{})["name"].(string),
			Value: v.(map[string]interface{})["value"].(string),
		}
	}

	return paramList
}

func resourceDmsKafkaTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dmsV2Client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	allTopics, err := topics.List(dmsV2Client, instanceID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS kafka topic")
	}

	topicID := d.Id()
	isFound := false
	var found topics.Topic
	for _, item := range allTopics {
		if item.Name == topicID {
			found = item
			isFound = true
			break
		}
	}

	if !isFound {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] DMS kafka topic %s: %+v", d.Id(), found)

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", found.Name),
		d.Set("partitions", found.Partition),
		d.Set("replicas", found.Replication),
		d.Set("aging_time", found.RetentionTime),
		d.Set("sync_replication", found.SyncReplication),
		d.Set("sync_flushing", found.SyncMessageFlush),
		d.Set("configs", flattenConfigs(found.Configs)),
		d.Set("description", found.Description),
		d.Set("created_at", utils.FormatTimeStampRFC3339(found.CreatedAt/1000, false)),
		d.Set("policies_only", found.PoliciesOnly),
		d.Set("type", setTopicType(found.TopicType)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func setTopicType(topicValue int) string {
	topicType := "common topic"
	if topicValue == 1 {
		topicType = "system topic"
	}

	return topicType
}

func flattenConfigs(params []topics.ConfigParam) []map[string]interface{} {
	if len(params) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(params))
	for i, val := range params {
		result[i] = map[string]interface{}{
			"name":  val.Name,
			"value": val.Value,
		}
	}
	return result
}

func buildBrokerList(params []interface{}) []int {
	paramList := make([]int, len(params))
	for i, v := range params {
		paramList[i] = v.(int)
	}

	return paramList
}

func resourceDmsKafkaTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dmsV2Client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	updateItem := topics.UpdateItem{
		Name: d.Get("name").(string),
	}
	// Set value if it changed.
	if d.HasChange("partitions") {
		newPartition := d.Get("partitions").(int)
		updateItem.Partition = &newPartition
		updateItem.NewPartitionBrokers = buildBrokerList(d.Get("new_partition_brokers").(*schema.Set).List())
	}
	if d.HasChange("aging_time") {
		retentionTime := d.Get("aging_time").(int)
		updateItem.RetentionTime = &retentionTime
	}
	if d.HasChange("sync_replication") {
		syncReplication := d.Get("sync_replication").(bool)
		updateItem.SyncReplication = &syncReplication
	}
	if d.HasChange("sync_flushing") {
		syncMessageFlush := d.Get("sync_flushing").(bool)
		updateItem.SyncMessageFlush = &syncMessageFlush
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateItem.Description = &description
	}
	if d.HasChange("configs") {
		updateItem.Configs = buildKafkaTopicParameters(d.Get("configs").(*schema.Set).List())
	}

	updateOpts := topics.UpdateOpts{
		Topics: []topics.UpdateItem{
			updateItem,
		},
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)
	instanceID := d.Get("instance_id").(string)
	err = topics.Update(dmsV2Client, instanceID, updateOpts).Err
	if err != nil {
		return diag.Errorf("error updating DMS kafka topic: %s", err)
	}

	return resourceDmsKafkaTopicRead(ctx, d, meta)
}

func resourceDmsKafkaTopicDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	dmsV2Client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	topicID := d.Id()
	instanceID := d.Get("instance_id").(string)
	response, err := topics.Delete(dmsV2Client, instanceID, []string{topicID}).Extract()
	if err != nil {
		return diag.Errorf("error deleting DMS kafka topic: %s", err)
	}

	var success bool
	for _, item := range response {
		if item.Name == topicID {
			success = item.Success
			break
		}
	}
	if !success {
		return diag.Errorf("error deleting DMS kafka topic")
	}

	d.SetId("")
	return nil
}

// resourceDmsKafkaTopicImport query the rules from HuaweiCloud and imports them to Terraform.
// It is a common function in waf and is also called by other rule resources.
func resourceDmsKafkaTopicImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for DMS kafka topic. Format must be <instance id>/<topic name>")
		return nil, err
	}

	instanceID := parts[0]
	topicID := parts[1]

	d.SetId(topicID)
	err := d.Set("instance_id", instanceID)

	return []*schema.ResourceData{d}, err
}

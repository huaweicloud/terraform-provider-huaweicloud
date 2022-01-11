package dms

import (
	"context"
	"strings"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/topics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// ResourceDmsKafkaTopic implements the resource of "huaweicloud_dms_kafka_topic"
func ResourceDmsKafkaTopic() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaTopicCreate,
		ReadContext:   resourceDmsKafkaTopicRead,
		UpdateContext: resourceDmsKafkaTopicUpdate,
		DeleteContext: resourceDmsKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceDmsKafkaTopicImport,
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
		},
	}
}

func resourceDmsKafkaTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS client: %s", err)
	}

	createOpts := &topics.CreateOps{
		Name:             d.Get("name").(string),
		Partition:        d.Get("partitions").(int),
		Replication:      d.Get("replicas").(int),
		RetentionTime:    d.Get("aging_time").(int),
		SyncReplication:  d.Get("sync_replication").(bool),
		SyncMessageFlush: d.Get("sync_flushing").(bool),
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	instanceID := d.Get("instance_id").(string)
	v, err := topics.Create(dmsV2Client, instanceID, createOpts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS kafka topic: %s", err)
	}

	// use topic name as the resource ID
	d.SetId(v.Name)
	return resourceDmsKafkaTopicRead(ctx, d, meta)
}

func resourceDmsKafkaTopicRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	allTopics, err := topics.List(dmsV2Client, instanceID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "DMS kafka topic")
	}

	topicID := d.Id()
	var found *topics.Topic
	for _, item := range allTopics {
		if item.Name == topicID {
			found = &item
			break
		}
	}

	if found == nil {
		d.SetId("")
		return nil
	}

	logp.Printf("[DEBUG] DMS kafka topic %s: %+v", d.Id(), found)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", found.Name),
		d.Set("partitions", found.Partition),
		d.Set("replicas", found.Replication),
		d.Set("aging_time", found.RetentionTime),
		d.Set("sync_replication", found.SyncReplication),
		d.Set("sync_flushing", found.SyncMessageFlush),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.FromErr(mErr)
	}

	return nil
}

func resourceDmsKafkaTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS client: %s", err)
	}

	updateItem := topics.UpdateItem{
		Name: d.Get("name").(string),
	}
	// Set value if it changed.
	if d.HasChange("partitions") {
		newPartition := d.Get("partitions").(int)
		updateItem.Partition = &newPartition
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

	updateOpts := topics.UpdateOpts{
		Topics: []topics.UpdateItem{
			updateItem,
		},
	}

	logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
	instanceID := d.Get("instance_id").(string)
	err = topics.Update(dmsV2Client, instanceID, updateOpts).Err
	if err != nil {
		return fmtp.DiagErrorf("error updating HuaweiCloud DMS kafka topic: %s", err)
	}

	return resourceDmsKafkaTopicRead(ctx, d, meta)
}

func resourceDmsKafkaTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating HuaweiCloud DMS client: %s", err)
	}

	topicID := d.Id()
	instanceID := d.Get("instance_id").(string)
	response, err := topics.Delete(dmsV2Client, instanceID, []string{topicID}).Extract()
	if err != nil {
		return fmtp.DiagErrorf("error deleting DMS kafka topic: %s", err)
	}

	var success bool
	for _, item := range response {
		if item.Name == topicID {
			success = item.Success
			break
		}
	}
	if !success {
		return fmtp.DiagErrorf("error deleting DMS kafka topic")
	}

	d.SetId("")
	return nil
}

// resourceDmsKafkaTopicImport query the rules from HuaweiCloud and imports them to Terraform.
// It is a common function in waf and is also called by other rule resources.
func resourceDmsKafkaTopicImport(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmtp.Errorf("Invalid format specified for DMS kafka topic. Format must be <instance id>/<topic name>")
		return nil, err
	}

	instanceID := parts[0]
	topicID := parts[1]

	d.SetId(topicID)
	d.Set("instance_id", instanceID)

	return []*schema.ResourceData{d}, nil
}

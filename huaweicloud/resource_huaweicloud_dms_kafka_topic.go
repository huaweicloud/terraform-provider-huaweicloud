package huaweicloud

import (
	"strings"

	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/topics"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

// ResourceDmsKafkaTopic implements the resource of "huaweicloud_dms_kafka_topic"
func ResourceDmsKafkaTopic() *schema.Resource {
	return &schema.Resource{
		Create: resourceDmsKafkaTopicCreate,
		Read:   resourceDmsKafkaTopicRead,
		Update: resourceDmsKafkaTopicUpdate,
		Delete: resourceDmsKafkaTopicDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDmsKafkaTopicImport,
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

func resourceDmsKafkaTopicCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DMS client: %s", err)
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
		return fmtp.Errorf("error creating HuaweiCloud DMS kafka topic: %s", err)
	}

	// use topic name as the resource ID
	d.SetId(v.Name)
	return resourceDmsKafkaTopicRead(d, meta)
}

func resourceDmsKafkaTopicRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	allTopics, err := topics.List(dmsV2Client, instanceID).Extract()
	if err != nil {
		return CheckDeleted(d, err, "DMS kafka topic")
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
		d.Set("region", GetRegion(d, config)),
		d.Set("name", found.Name),
		d.Set("partitions", found.Partition),
		d.Set("replicas", found.Replication),
		d.Set("aging_time", found.RetentionTime),
		d.Set("sync_replication", found.SyncReplication),
		d.Set("sync_flushing", found.SyncMessageFlush),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}

	return nil
}

func resourceDmsKafkaTopicUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DMS client: %s", err)
	}

	newPartition := d.Get("partitions").(int)
	retentionTime := d.Get("aging_time").(int)
	syncReplication := d.Get("sync_replication").(bool)
	syncMessageFlush := d.Get("sync_flushing").(bool)

	updateOpts := topics.UpdateOpts{
		Topics: []topics.UpdateItem{
			{
				Name:             d.Get("name").(string),
				Partition:        &newPartition,
				RetentionTime:    &retentionTime,
				SyncMessageFlush: &syncMessageFlush,
				SyncReplication:  &syncReplication,
			},
		},
	}

	logp.Printf("[DEBUG] Update Options: %#v", updateOpts)
	instanceID := d.Get("instance_id").(string)
	err = topics.Update(dmsV2Client, instanceID, updateOpts).Err
	if err != nil {
		return fmtp.Errorf("error updating HuaweiCloud DMS kafka topic: %s", err)
	}

	return resourceDmsKafkaTopicRead(d, meta)
}

func resourceDmsKafkaTopicDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	dmsV2Client, err := config.DmsV2Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("error creating HuaweiCloud DMS client: %s", err)
	}

	topicID := d.Id()
	instanceID := d.Get("instance_id").(string)
	response, err := topics.Delete(dmsV2Client, instanceID, []string{topicID}).Extract()
	if err != nil {
		return fmtp.Errorf("error deleting DMS kafka topic: %s", err)
	}

	var success bool
	for _, item := range response {
		if item.Name == topicID {
			success = item.Success
			break
		}
	}
	if !success {
		return fmtp.Errorf("error deleting DMS kafka topic")
	}

	d.SetId("")
	return nil
}

// resourceDmsKafkaTopicImport query the rules from HuaweiCloud and imports them to Terraform.
// It is a common function in waf and is also called by other rule resources.
func resourceDmsKafkaTopicImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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

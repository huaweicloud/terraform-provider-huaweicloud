---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_group_topic_batch_delete"
description: |-
  Use this resource to batch delete subscribed topics under the specified consumer group within HuaweiCloud.
---

# huaweicloud_dms_kafka_consumer_group_topic_batch_delete

Use this resource to batch delete subscribed topics under the specified consumer group within HuaweiCloud.

-> 1. This resource will permanently delete the consumer offset in the topics.
   <br>2. This resource is only a one-time action resource for batch deleting subscribed topics under the consumer group.
   Deleting this resource will not clear the corresponding request record, but will only remove the resource
   information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "consumer_group_id" {}
variable "topic_names" {
  type = list(string)
}

resource "huaweicloud_dms_kafka_consumer_group_topic_batch_delete" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_id
  topics      = var.topic_names
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the topics to be deleted in batches
  are located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance.

* `group` - (Required, String, NonUpdatable) Specifies the ID of the consumer group.

* `topics` - (Required, List, NonUpdatable) Specifies the list of topic names to be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `result` - The result of the batch delete operation.  
 The [result](#kafka_consumer_group_topic_batch_delete_result) structure is documented below.

<a name="kafka_consumer_group_topic_batch_delete_result"></a>
The `result` block supports:

* `name` - The name of the topic.

* `success` - Whether the topic was deleted successfully.

* `error_code` - The error code if the topic deletion failed.

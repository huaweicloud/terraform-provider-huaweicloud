---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_topic_message_batch_delete"
description: |-
  Use this resource to delete messages under the specified topics within HuaweiCloud.
---

# huaweicloud_dms_kafka_topic_message_batch_delete

Use this resource to delete messages under the specified topics within HuaweiCloud.

-> This resource is only a one-time action resource for deleting topic messages in batches. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from the
   tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}
variable "partitions" {
  type = list(object({
    partition = number
    offset    = number
  }))
}

resource "huaweicloud_dms_kafka_topic_message_batch_delete" "test" {
  instance_id = var.instance_id
  topic       = var.topic_name

  dynamic "partitions" {
    for_each = var.partitions

    content {
      partition = partitions.value.partition
      offset    = partitions.value.offset
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the topic messages to be deleted are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Kafka instance.

* `topic` - (Required, String, NonUpdatable) Specifies the name of the topic.

* `partitions` - (Required, List, NonUpdatable) Specifies the partition configuration list to which the messages
  to be deleted belong.  
  The [partitions](#kafka_topic_message_delete_partitions) structure is documented below.

<a name="kafka_topic_message_delete_partitions"></a>
The `partitions` block supports:

* `partition` - (Required, Int, NonUpdatable) Specifies the number of the partition.

* `offset` - (Required, Int, NonUpdatable) Specifies the offset of the message to be deleted.

 -> The data after the earliest offset and before this offset will be deleted. For example, if the earliest offset
   is `2` and the entered offset is `5`, the messages whose offset ranges from `2` to `4` will be deleted.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `result` - The result of the message delete operation.  
 The [result](#kafka_topic_message_delete_result) structure is documented below.

<a name="kafka_topic_message_delete_result"></a>
The `result` block supports:

* `partition` - The number of the partition.

* `result` - The operation result.

* `error_code` - The error code if the operation failed.

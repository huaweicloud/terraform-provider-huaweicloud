---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_consumer_group_message_offsets"
description: |-
  Use this data source to get the message offset list under the specified consumer group within HuaweiCloud.
---

# huaweicloud_dms_kafka_consumer_group_message_offsets

Use this data source to get the message offset list under the specified consumer group within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "consumer_group_name" {}
variable "topic_name" {}

data "huaweicloud_dms_kafka_consumer_group_message_offsets" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_name
  topic       = var.topic_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the consumer group message offsets are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the Kafka instance.

* `group` - (Required, String) Specifies the name of the consumer group.

* `topic` - (Required, String) Specifies the name of the topic.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `message_offsets` - The list of consumer group message offsets.  
  The [message_offsets](#kafka_consumer_group_message_offsets_struct) structure is documented below.

<a name="kafka_consumer_group_message_offsets_struct"></a>
The `message_offsets` block supports:

* `partition` - The name of the partition.

* `message_current_offset` - The current offset of the message.

* `message_log_start_offset` - The start offset of the message.

* `message_log_end_offset` - The end offset of the message.

* `consumer_id` - The consumer ID of the consumed message.

* `host` - The consumer address of the consumed message.

* `client_id` - The ID of the client.

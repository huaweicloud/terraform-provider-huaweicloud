---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_dead_letter_messages"
description: |-
  Use this data source to get the list of RocketMQ dead letter messages within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_dead_letter_messages

Use this data source to get the list of RocketMQ dead letter messages within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "topic_name" {}
variable "msg_id_list" {
  type = list(string)
}

data "huaweicloud_dms_rocketmq_dead_letter_messages" "test" {
  instance_id  = var.instance_id
  topic        = var.topic_name
  msg_id_list  = var.msg_id_list
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dead letter messages are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `topic` - (Required, String) Specifies the name of the topic to which the dead letter messages belong.

* `msg_id_list` - (Required, List) Specifies the list of dead letter message IDs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `messages` - All dead letter messages that match the filter parameters.  
  The [messages](#dead_letter_messages) structure is documented below.

<a name="dead_letter_messages"></a>
The `messages` block supports:

* `msg_id` - The ID of the dead letter message.

* `instance_id` - The ID of the RocketMQ instance.

* `topic` - The name of the topic to which the dead letter message belongs.

* `store_time` - The time when the dead letter message was stored, in RFC3339 format.

* `born_time` - The time when the dead letter message was generated, in RFC3339 format.

* `reconsume_times` - The number of times the message has been retried.

* `body` - The body of the message.

* `body_crc` - The checksum of the message body.

* `store_size` - The storage size of the message.

* `born_host` - The IP address of the host that generated the message.

* `store_host` - The IP address of the host that stored the message.

* `queue_id` - The ID of the queue.

* `queue_offset` - The offset in the queue.

* `property_list` - The list of message properties.  
  The [property_list](#dead_letter_messages_property_list) structure is documented below.

<a name="dead_letter_messages_property_list"></a>
The `property_list` block supports:

* `name` - The name of the property.

* `value` - The value of the property.

---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_messages"
description: |-
  Use this data source to get the list of RocketMQ instance messages.
---

# huaweicloud_dms_rocketmq_messages

Use this data source to get the list of RocketMQ instance messages.

## Example Usage

### Query message by topic

```hcl
variable "instance_id" {}
variable "topic" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  start_time  = var.start_time
  end_time    = var.end_time
}
```

### Query message by key

```hcl
variable "instance_id" {}
variable "topic" {}
variable "key" {}

data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  key         = var.key
}
```

### Query message by message ID

```hcl
variable "instance_id" {}
variable "topic" {}
variable "message_id" {}

data "huaweicloud_dms_rocketmq_messages" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  message_id  = var.message_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `topic` - (Required, String) Specifies the topic name.

* `start_time` - (Optional, String) Specifies the start time, a Unix timestamp in millisecond.

* `end_time` - (Optional, String) Specifies the end time, a Unix timestamp in millisecond.

* `key` - (Optional, String) Specifies the message key.

* `message_id` - (Optional, String) Specifies the message ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `messages` - Indicates the message list.
  The [messages](#attrblock--messages) structure is documented below.

<a name="attrblock--messages"></a>
The `messages` block supports:

* `message_id` - Indicates the message ID.

* `body` - Indicates the message body. Only return when querying message by message ID.

* `body_crc` - Indicates the message body checksum.

* `property_list` - Indicates the property list.
  The [property_list](#attrblock--messages--property_list) structure is documented below.

* `queue_id` - Indicates the queue ID.

* `queue_offset` - Indicates the offset in the queue.

* `reconsume_times` - Indicates the number of retry times.

* `born_host` - Indicates the IP address of the host that generates the message.

* `store_host` - Indicates the IP address of the host that stores the message.

* `store_size` - Indicates the storage size.

* `born_time` - Indicates the message generated time.

* `store_time` - Indicates the message stored time.

<a name="attrblock--messages--property_list"></a>
The `property_list` block supports:

* `name` - Indicates the property name.

* `value` - Indicates the property value.

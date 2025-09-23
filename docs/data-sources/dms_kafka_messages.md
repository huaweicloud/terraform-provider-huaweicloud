---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_messages"
description: |-
  Use this data source to get the list of Kafka messages.
---

# huaweicloud_dms_kafka_messages

Use this data source to get the list of Kafka messages.

## Example Usage

### Query messages by creation time

```hcl
variable "instance_id" {}
variable "topic" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dms_kafka_messages" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  start_time  = var.start_time
  end_time    = var.end_time
}
```

### Query messages by content's keyword, a maximum of 10 messages can be returned

```hcl
variable "instance_id" {}
variable "topic" {}
variable "start_time" {}
variable "end_time" {}
variable "keyword" {}

data "huaweicloud_dms_kafka_messages" "test" {
  instance_id = var.instance_id
  topic       = var.topic
  start_time  = var.start_time
  end_time    = var.end_time
  keyword     = var.keyword
}
```

### Query messages content by offset

```hcl
variable "instance_id" {}
variable "topic" {}
variable "partition" {}
variable "message_offset" {}

data "huaweicloud_dms_kafka_messages" "test" {
  instance_id    = var.instance_id
  topic          = var.topic
  partition      = var.partition
  message_offset = var.message_offset
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `topic` - (Required, String) Specifies the topic name.

* `start_time` - (Optional, String) Specifies the start time, a Unix timestamp in millisecond.
  This parameter is mandatory when you query the message creation time.

* `end_time` - (Optional, String) Specifies the end time, a Unix timestamp in millisecond.
  This parameter is mandatory when you query the message creation time.

* `download` - (Optional, Bool) Whether download is required.
  If it is **false**, the big message will be truncated. Defaults to **false**.

* `message_offset` - (Optional, String) Specifies the message offset.
  This parameter is mandatory when you query the message content by offset.

* `partition` - (Optional, String) Specifies the partition.
  This parameter is mandatory when you query the message content by offset.

* `keyword` - (Optional, String) Specifies the keyword.
  If it's specified, a maximum of **10** messages can be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `messages` - Indicates the message list.

  The [messages](#messages_struct) structure is documented below.

<a name="messages_struct"></a>
The `messages` block supports:

* `key` - Indicates the message key.

* `value` - Indicates the message content.

* `timestamp` - Indicates the message production time.

* `huge_message` - Indicates the big data flag.

* `message_offset` - Indicates the message offset.

* `message_id` - Indicates the message ID.

* `partition` - Indicates the partition where the message is located.

* `size` - Indicates the message size.

* `app_id` - Indicates the application ID.

* `tag` - Indicates the message label.

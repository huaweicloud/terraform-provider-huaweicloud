---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_message_offset_reset"
description: |-
  Use this resource to reset the message offset under the specified consumer group within HuaweiCloud.
---

# huaweicloud_dms_kafka_message_offset_reset

Use this resource to reset the message offset under the specified consumer group within HuaweiCloud.

-> 1. Before using this resource to reset the consumption progress, you must stop the client of the reset consumer group.
   <br>2. This resource is only a one-time action resource for resetting the consumption progress of the consumer group.
   Deleting this resource will not clear the corresponding request record, but will only remove the resource information
   from the tfstate file.

## Example Usage

### Reset message offset for all topic with timestamp

```hcl
variable "instance_id" {}
variable "consumer_group_name" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_name
  partition   = -1
  timestamp   = "0"
}
```

### Reset message offset for all partition under specific topic with timestamp

```hcl
variable "instance_id" {}
variable "consumer_group_name" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = var.instance_id
  group       = var.consumer_group_name
  topic       = var.topic_name
  partition   = -1
  timestamp   = "0"
}
```

### Reset message offset for all partition under specific topic with message offset

```hcl
variable "instance_id" {}
variable "consumer_group_name" {}
variable "topic_name" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id    = var.instance_id
  group          = var.consumer_group_name
  topic          = var.topic_name
  partition      = 0
  message_offset = "0"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the consumption progress is to be reset is located.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the Kafka instance.  
  Changing this creates a new resource.

* `group` - (Required, String, ForceNew) Specifies the name of the consumer group.  
  Changing this creates a new resource.

* `partition` - (Required, Int, ForceNew) Specifies the partiton number.
  
  -> If value is `-1`, it means to reset offset of all partitions.

* `topic` - (Optional, String, ForceNew) Specifies name of the topic.  
  Changing this creates a new resource.

  -> When `topic` is not specified, it means to reset offset of all topics.

* `message_offset` - (Optional, String, ForceNew) Specifies the offset to reset the consumption progress.
  + If this offset is earlier than the current earliest offset, the offset will be reset to the earliest offset.
  + If this offset is later than the current largest offset, the offset will be reset to the latest offset.

* `timestamp` - (Optional, String, ForceNew) Specifies the time that the offset is to be reset to.
  The value is a Unix timestamp, in millisecond.
  + If this time is earlier than the current earliest timestamp, the offset will be reset to the earliest timestamp.
  + If this time is later than the current largest timestamp, the offset will be reset to the latest timestamp.

  Changing this creates a new resource.

-> Exactly one of `message_offset` and `timestamp` should be specified, when `topic` is empty, only support to reset
  with `timestamp`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

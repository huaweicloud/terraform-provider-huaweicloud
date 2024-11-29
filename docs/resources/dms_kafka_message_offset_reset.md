---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_message_offset_reset"
description: |-
  Manage DMS kafka message offset reset resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_message_offset_reset

Manage DMS kafka message offset reset resource within HuaweiCloud.

## Example Usage

### Reset message offset for all topic with timestamp

```hcl
variable "instance_id" {}
variable "group" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = var.instance_id
  group       = var.group
  topic       = ""
  partition   = -1
  timestamp   = 0
}
```

### Reset message offset for all partition under specific topic with timestamp

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id = var.instance_id
  group       = var.group
  topic       = var.topic
  partition   = -1
  timestamp   = 0
}
```

### Reset message offset for all partition under specific topic with message offset

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}

resource "huaweicloud_dms_kafka_message_offset_reset" "test" {
  instance_id    = var.instance_id
  group          = var.group
  topic          = var.topic
  partition      = -1
  message_offset = 0
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID.
  Changing this creates a new resource.

* `group` - (Required, String, ForceNew) Specifies the group name.
  Changing this creates a new resource.

* `partition` - (Required, Int, ForceNew) Specifies the partiton number.
  + If value is **-1**, reset all partitions. When `topic` is empty, only support reset all partitions.
  + If value is specific number, reset that partiton only.

  Changing this creates a new resource.

* `topic` - (Optional, String, ForceNew) Specifies the topic name. If it is empty, reset all topic.
  Changing this creates a new resource.

* `message_offset` - (Optional, String, ForceNew) Specifies the message offset.
  + If this offset is earlier than the current earliest offset, the offset will be reset to the earliest offset.
  + If this offset is later than the current largest offset, the offset will be reset to the latest offset.

  Changing this creates a new resource.

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

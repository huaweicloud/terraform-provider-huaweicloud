---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_message_offset_reset"
description: |-
  Manages a DMS RocketMQ message offset reset resource within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_message_offset_reset

Manages a DMS RocketMQ message offset reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "group" {}
variable "topic" {}

resource "huaweicloud_dms_rocketmq_message_offset_reset" "test" {
  instance_id = var.instance_id
  group       = var.group
  topic       = var.topic
  timestamp   = 0
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

* `topic` - (Required, String, ForceNew) Specifies the topic name.
  Changing this creates a new resource.

* `timestamp` - (Required, String, ForceNew) Specifies the timestamp.
  + If it is specified as **0**, reset to earliset.
  + If it is specified as **-1**, reset to latest.
  + If it is specified as a timestamp in milliseconds, reset to specific time.

  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_kafka_smart_connect_task_action"
description: |-
  Manage DMS kafka smart connect task action resource within HuaweiCloud.
---

# huaweicloud_dms_kafka_smart_connect_task_action

Manage DMS kafka smart connect task action resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "task_id" {}
variable "action" {}

resource "huaweicloud_dms_kafka_smart_connect_task_action" "test" {
  instance_id = var.instance_id
  task_id     = var.task_id
  action      = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the kafka instance ID.
  Changing this parameter will create a new resource.

* `task_id` - (Required, String, ForceNew) Specifies the smart connect task ID.
  Changing this parameter will create a new resource.

* `action` - (Required, String) Specifies the smart connect task action.
  Valid values are:
  + **pause**: Pause the task from running status.
  + **resume**: Resume the task from paused status.
  + **start**: Start the task from waiting status.
  + **restart**: Restart the job from pausd or running status.

## Timeouts

This resource provides the following timeout configuration options:

* `update` - Default is 20 minutes.

* `create` - Default is 20 minutes.

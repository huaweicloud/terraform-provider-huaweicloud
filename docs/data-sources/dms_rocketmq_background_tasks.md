---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_background_tasks"
description: |-
  Use this data source to get the background task list under the specified RocketMQ instance within HuaweiCloud.
---

# huaweicloud_dms_rocketmq_background_tasks

Use this data source to query the background task list under the specified RocketMQ instance within HuaweiCloud.

## Example Usage

### Query all background tasks

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_background_tasks" "test" {
  instance_id = var.instance_id
}
```

### Query background tasks within a specified time range

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}

data "huaweicloud_dms_rocketmq_background_tasks" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the background tasks are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RocketMQ instance.

* `begin_time` - (Optional, String) Specifies the start time of the background tasks, in UTC format.  
  The format is `YYYYMMDDHHmmss`.

* `end_time` - (Optional, String) Specifies the end time of the background tasks, in UTC format.  
  The format is `YYYYMMDDHHmmss`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the background tasks under the specified RocketMQ instance.  
  The [tasks](#data_rocketmq_background_tasks) structure is documented below.

<a name="data_rocketmq_background_tasks"></a>
The `tasks` block supports:

* `id` - The ID of the background task.

* `name` - The name of the background task.

* `user_name` - The username of the user who executed the background task.

* `user_id` - The ID of the user who executed the background task.

* `params` - The parameters of the background task.

* `status` - The status of the background task.

* `created_at` - The creation time of the background task, in UTC format.

* `updated_at` - The latest update time of the background task, in UTC format.

---
subcategory: "Distributed Message Service (DMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dms_rocketmq_migration_tasks"
description: |-
  Use this data source to get the list of RocketMQ instance's migration tasks.
---

# huaweicloud_dms_rocketmq_migration_tasks

Use this data source to get the list of RocketMQ instance's migration tasks.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dms_rocketmq_migration_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the RocketMQ instance ID.

* `task_id` - (Optional, String) Specifies the RocketMQ migration task ID.

* `type` - (Optional, String) Specifies the RocketMQ migration task type.
  Valid values are **rocketmq** and **rabbitToRocket**.

* `name` - (Optional, String) Specifies the RocketMQ migration task name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the list of metadata migration tasks.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the ID of a metadata migration task.

* `type` - Indicates the metadata migration task type.

* `name` - Indicates the name of a metadata migration task.

* `start_date` - Indicates the start time of a metadata migration task.

* `status` - Indicates the status of a metadata migration task.

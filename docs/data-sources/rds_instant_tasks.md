---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instant_tasks"
description: |-
  Use this data source to query RDS instant tasks.
---

# huaweicloud_rds_instant_tasks

Use this data source to query RDS instant tasks.

## Example Usage

```hcl
data "huaweicloud_rds_instant_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `task_id` - (Optional, String) Specifies the ID of the task.

* `name` - (Optional, String) Specifies the name of the task.

* `status` - (Optional, String) Specifies the status of the task. Value options: **Running**, **Completed**, **Failed**.

* `instance_id` - (Optional, String) Specifies the ID of the instance.

* `order_id` - (Optional, String) Specifies the ID of the order.

* `start_time` - (Optional, String) Specifies the start time in UTC timestamp format (milliseconds since epoch).
  `end_time` is mandatory if it is not empty.

* `end_time` - (Optional, String) Specifies the end time in UTC timestamp format (milliseconds since epoch).
  `start_time` is mandatory if it is not empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the list of instant tasks.

  The [tasks](#tasks_struct) structure is documented below.

* `actions` - Indicates the list of task names.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `instance_id` - Indicates the ID of the instance.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the instance status. The value can be: **BUILD**, **CREATE FAIL**, **ACTIVE**, **FAILED**,
  **FROZEN**, **MODIFYING**, **REBOOTING**, **RESTORING**, **MODIFYING INSTANCE TYPE**, **SWITCHOVER**, **MIGRATING**,
  **BACKING UP**, **MODIFYING DATABASE PORT**.

* `process` - Indicates the task execution process. The execution progress (such as "60", indicating the task execution
  progress is 60%) is displayed only when the task is being executed. Otherwise, "" is returned.

* `order_id` - Indicates the ID of the order.

* `status` - Indicates the task execution status. The value can be:
  + **Running**: Indicates the task is being executed.
  + **Completed**: Indicates the task is successfully executed.
  + **Failed**: Indicates the task fails to be executed.

* `fail_reason` - Indicates the error information displayed when a task failed.

* `create_time` - Indicates the creation time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the end time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.

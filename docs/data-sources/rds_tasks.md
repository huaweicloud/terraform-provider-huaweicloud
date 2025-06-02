---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_tasks"
description: |-
  Use this data source to get the list of asynchronous tasks for a specified
  RDS instance within a specified time range.
---

# huaweicloud_rds_tasks

Use this data source to get the list of asynchronous tasks for a specified
RDS instance within a specified time range.

## Example Usage

```hcl
variable "instance_id" {}
variable "start_time" {}

data "huaweicloud_rds_tasks" "test" {
  instance_id = var.instance_id
  start_time  = var.start_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `start_time` - (Required, String) Specifies the start time in UTC timestamp format
  (milliseconds since epoch).

* `end_time` - (Optional, String) Specifies the end time in UTC timestamp format
  (milliseconds since epoch).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the task information.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `status` - Indicates the task execution status.
  + **Running**: Indicates the task is being executed.
  + **Completed**: Indicates the task is successfully executed.
  + **Failed**: Indicates the task fails to be executed.

* `created` - Indicates the creation time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.
  T is the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `ended` - Indicates the end time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format. T is
  the separator between the calendar and the hourly notation of time. Z indicates the time
  zone offset. For example, in the Beijing time zone, the time zone offset is shown as +0800.

* `process` - Indicates the task execution progress. The execution progress (such as "60",
  indicating the task execution progress is 60%) is displayed only when the task is being
  executed. Otherwise, "" is returned.

* `instance` - Indicates the information of the DB instance on which the task is executed.

  The [instance](#jobs_instance_struct) structure is documented below.

* `task_detail` - Indicates the displayed information varies depending on the tasks.

* `fail_reason` - Indicates the error information displayed when a task failed.

* `entities` - Indicates the displayed information varies depending on the tasks.

<a name="jobs_instance_struct"></a>
The 'instance' block supports:

* `id` - Indicates the DB instance ID.

* `name` - Indicates the DB instance name.

---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_schedule_task_executions"
description: |-
  Use this data source to get execution list under specified schedule task within HuaweiCloud.
---

# huaweicloud_workspace_app_schedule_task_executions

Use this data source to get execution list under specified schedule task within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_workspace_app_schedule_task_executions" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the schedule task is located.  
  If omitted, the provider-level region will be used.

* `task_id` - (Required, String) Specifies the ID of the schedule task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `executions` - The list of schedule task executions.
  The [executions](#schedule_task_executions_attr) structure is documented below.

<a name="schedule_task_executions_attr"></a>
The `executions` block supports:

* `id` - The ID of the schedule task execution record.

* `task_id` - The ID of the schedule task.

* `task_type` - The type of the schedule task.

* `scheduled_type` - The execution cycle of the schedule task.
  + **FIXED_TIME**
  + **DAY**
  + **WEEK**
  + **MONTH**

* `status` - The status of the schedule task execution.
  + **SUCCESS**
  + **FAILED**

* `total_count` - The total number of subtasks.

* `success_count` - The number of successful subtasks.

* `failed_count` - The number of failed subtasks.

* `time_zone` - The timezone of the schedule task.

* `begin_time` - The begin time of the schedule task execution, in UTC format.

* `create_time` - The creation time of the schedule task, in UTC format.

* `end_time` - The end time of the schedule task execution, in UTC format.

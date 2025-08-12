---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_schedule_task_execute_details"
description: |-
  Use this data source to get the details of Workspace APP schedule task executions within HuaweiCloud.
---

# huaweicloud_workspace_app_schedule_task_execute_details

Use this data source to get the details of Workspace APP schedule task executions within HuaweiCloud.

## Example Usage

```hcl
variable "execute_history_id" {}

data "huaweicloud_workspace_app_schedule_task_execute_details" "test" {
  execute_history_id = var.execute_history_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the schedule task execution details are located.  
  If omitted, the provider-level region will be used.

* `execute_history_id` - (Required, String) Specifies the ID of the schedule task execution record.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `execute_details` - The list of the sub-task execution details.
  The [execute_details](#app_schedule_task_execute_details) structure is documented below.

<a name="app_schedule_task_execute_details"></a>
The `execute_details` block supports:

* `id` - The ID of the sub-task to be executed.

* `execute_id` - The ID of the schedule task execution record.

* `server_id` - The ID of the server being operated.

* `server_name` - The name of the server being operated.

* `status` - The status of schedule task execution.
  + **SUCCESS**
  + **FAILED**

* `task_type` - The type of the schedule task.
  + **RESTART_SERVER**
  + **START_SERVER**
  + **STOP_SERVER**
  + **REINSTALL_OS**

* `time_zone` - The timezone of the schedule task.

* `begin_time` - The start time of the sub-task, in UTC format.

* `end_time` - The end time of the sub-task, in UTC format.

* `result_code` - The error code when the task execution fails.

* `result_message` - The reason of task failure.

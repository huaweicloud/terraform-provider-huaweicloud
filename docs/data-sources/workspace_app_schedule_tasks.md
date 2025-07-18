---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_schedule_tasks"
description: |-
  Use this data source to get the schedule task list of the Workspace APP within HuaweiCloud.
---

# huaweicloud_workspace_app_schedule_tasks

Use this data source to get the schedule task list of the Workspace APP within HuaweiCloud.

## Example Usage

### Query all schedule tasks

```hcl
data "huaweicloud_workspace_app_schedule_tasks" "test" {}
```

### Query schedule task by specified name

```hcl
variable "task_name" {}

data "huaweicloud_workspace_app_schedule_tasks" "test" {
  task_name = var.task_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the schedule tasks are located.  
  If omitted, the provider-level region will be used.

* `task_name` - (Optional, String) Specifies the name of the schedule task.

* `task_type` - (Optional, String) Specifies the type of the schedule task.  
  The valid values are as follows:
  + **RESTART_SERVER**
  + **START_SERVER**
  + **STOP_SERVER**
  + **REINSTALL_OS**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - All schedule tasks that match the filter parameters.  
  The [tasks](#data_schedule_tasks) structure is documented below.

<a name="data_schedule_tasks"></a>
The `tasks` block supports:

* `id` - The ID of the schedule task.

* `task_name` - The name of the schedule task.

* `scheduled_type` - The execution cycle of the schedule task.
  + **FIXED_TIME**
  + **DAY**
  + **WEEK**
  + **MONTH**

* `task_type` - The type of the schedule task.

* `scheduled_time` - The scheduled time of the schedule task.  
  The format is `HH:MM:SS`.

* `day_interval` - The interval in days of the schedule task.

* `week_list` - The days of the week of the schedule task.

* `month_list` - The month of the schedule task.

* `date_list` - The days of the month of the schedule task.

* `scheduled_date` - The fixed time of the schedule task.  
  The format is `YYYY-MM-DD`.

* `time_zone` - The time zone of the schedule task.

* `expire_time` - The expire time of the schedule task.

* `description` - The description of the schedule task.

* `is_enable` - Whether the schedule task is enabled.

* `task_cron` - The cron expression of the schedule task.

* `next_execution_time` - The next execution time of the schedule task.

* `created_at` - The create time of the schedule task, in RFC3339 format.

* `updated_at` - The latest update time of the schedule task, in RFC3339 format.

* `last_status` - The last execution status of the schedule task.
  + **RUNNING**
  + **SUCCESS**
  + **FAILED**

---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_scheduled_tasks"
description: |-
  Use this data source to get the list of COC scheduled tasks.
---

# huaweicloud_coc_scheduled_tasks

Use this data source to get the list of COC scheduled tasks.

## Example Usage

```hcl
data "huaweicloud_coc_scheduled_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `task_id` - (Optional, String) Specifies the task ID.

* `task_name` - (Optional, String) Specifies the task name.

* `scheduled_type` - (Optional, String) Specifies the scheduled task execution strategy.
  The valid values are as follows:
  + **PERIODIC**: Periodic execution.
  + **ONCE**: Single execution.
  + **CRON**: Execute according to a CRON expression.

* `task_type` - (Optional, String) Specifies the scheduled task type.
  The valid values are as follows:
  + **SCRIPT**: Script.
  + **RUNBOOK**: Job.

* `associated_task_type` - (Optional, String) Specifies the type of the associated task.
  The valid values are as follows:
  + **CUSTOMIZATION**: Custom script/job.
  + **COMMUNAL**: Public script/job.

* `risk_level` - (Optional, String) Specifies the risk level of the scheduled task.
  The valid values can be **HIGH**, **MEDIUM** or **LOW**.

* `created_by` - (Optional, String) Specifies the IAM user ID of the person who created the scheduled task.

* `reviewer` - (Optional, String) Specifies the IAM user ID of the scheduled task approver.

* `reviewer_user_name` - (Optional, String) Specifies the IAM user nickname of the scheduled task approver.

* `approve_status` - (Optional, String) Specifies the approval status of the scheduled task.
  The valid values are as follows:
  + **PASSED**: Passed.
  + **PENDING**: Pending approval.
  + **REJECTED**: Rejected.

* `last_execution_status` - (Optional, String) Specifies the most recent execution status of a scheduled task.
  The valid values are as follows:
  + **READY**: Ready for execution.
  + **PROCESSING**: Executing.
  + **ABNORMAL**: Abnormal.
  + **PAUSED**: Paused.
  + **CANCELED**: Cancelled.
  + **FINISHED**: Successful.

* `last_execution_start_time` - (Optional, Int) Specifies the timestamp of the query start time for the most recent
  scheduled task execution. UTC timestamp in milliseconds.

* `last_execution_end_time` - (Optional, Int) Specifies the timestamp of the query end time for the most recent
  scheduled task execution. UTC timestamp in milliseconds.

* `region_id` - (Optional, String) Specifies the region ID.

* `resource_id` - (Optional, String) Specifies the resource ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scheduled_tasks` - Indicates the list of scheduled tasks.

  The [scheduled_tasks](#scheduled_tasks_struct) structure is documented below.

<a name="scheduled_tasks_struct"></a>
The `scheduled_tasks` block supports:

* `id` - Indicates the task ID.

* `enterprise_project_id` - Indicates the enterprise project ID.

* `name` - Indicates the task name.

* `scheduled_type` - Indicates the task execution strategy.

* `task_type` - Indicates the reference task type.

* `associated_task_type` - Indicates the properties (public/custom) of the task associated with the task.

* `risk_level` - Indicates the risk level of the task.

* `created_by` - Indicates the task creator.

* `update_by` - Indicates the task modifier.

* `created_user_name` - Indicates the nickname of the task creator.

* `reviewer` - Indicates the task approver.

* `reviewer_user_name` - Indicates the nickname of the task approver.

* `approve_status` - Indicates the task approval status.

* `last_execution_time` - Indicates the timestamp of the most recent task execution time.

* `last_execution_status` - Indicates the most recent execution status of a task.

* `execution_count` - Indicates the number of times a task is executed.

* `enabled` - Indicates whether the task is enabled.

* `created_time` - Indicates the task create time.

* `modified_time` - Indicates the task update time.

* `region_id` - Indicates the region to which the task belongs.

* `associated_task_name` - Indicates the script/job name associated with the task.

* `associated_task_name_en` - Indicates the English name of the script or job associated with the task.

---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_scheduled_task"
description: |-
  Manages a COC scheduled task resource within HuaweiCloud.
---

# huaweicloud_coc_scheduled_task

Manages a COC scheduled task resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "time_zone" {}
variable "single_scheduled_time" {}
variable "associated_task_id" {}
variable "associated_task_name" {}

resource "huaweicloud_coc_scheduled_task" "test" {
  name       = var.name
  version_no = "1.0.0"
  trigger_time {
    time_zone             = var.time_zone
    policy                = "ONCE"
    single_scheduled_time = var.single_scheduled_time
  }
  task_type            = "SCRIPT"
  associated_task_id   = var.associated_task_id
  associated_task_type = "CUSTOMIZATION"
  associated_task_name = var.associated_task_name
  risk_level           = "LOW"
  
  input_param = {
    key = "value"
  }
  target_instances {
    target_selection = "NONE"
    order_no         = 1
  }
  enable_approve              = false
  enable_message_notification = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the task name.

* `version_no` - (Required, String) Specifies the version number.

* `trigger_time` - (Required, List) Specifies the execution time details of the scheduled task.

  The [trigger_time](#trigger_time_struct) structure is documented below.

* `task_type` - (Required, String) Specifies the task type. The default value is **SCRIPT**.
  Values can be as follows:
  + **SCRIPT**: Script.
  + **RUNBOOK**: Job.

* `associated_task_id` - (Required, String) Specifies the associated task ID (script ID or job ID).

* `associated_task_type` - (Required, String) Specifies the type of associated task. The default value is **CUSTOMIZATION**.
  Values can be as follows:
  + **CUSTOMIZATION**: Custom script/job.
  + **COMMUNAL**: Public script/job.

* `associated_task_name` - (Required, String) Specifies the name of the associated task (script name/job name).

* `risk_level` - (Required, String) Specifies the risk level of the scheduled task. The default value is **HIGH**.
  Values can be as follows:
  + **HIGH**: High risk.
  + **MEDIUM**: Medium risk.
  + **LOW**: Low risk.

* `input_param` - (Required, Map) Specifies the task execution parameters.

* `target_instances` - (Required, List) Specifies the target instance information.

  The [target_instances](#target_instances_struct) structure is documented below.

  -> **NOTE:** If the task is associated with a script, this array has only one element, which specifies the resource
  instance that the script operates on.
  <br>If the task is associated with a job and the target instance is **SAME**, this array has only one element, which
  specifies the resource instance that the job operates on.
  <br>If the task is associated with a job and the target instance is **DIFF**, the number of elements in this array is
  equal to the number of steps in the associated job. That is, each task under each step specifies the resource instance
  to operate on.

* `enable_approve` - (Required, Bool) Specifies whether to enable manual review.

* `enable_message_notification` - (Required, Bool) Specifies whether to enable message notifications.

* `ticket_infos` - (Optional, List) Specifies the change management related configuration information.

  The [ticket_infos](#ticket_infos_struct) structure is documented below.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `agency_name` - (Optional, String) Specifies the delegate name. The default value is **ServiceAgencyForCOC**.

* `associated_task_name_en` - (Optional, String) Specifies the English name of the associated task.

* `associated_task_enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the associated task.

* `runbook_instance_mode` - (Optional, String) Specifies the target instance mode. The default value is **SAME**.
  Values can be as follows:
  + **SAME**: All job steps are associated with the same instance resource.
  + **DIFF**: Each task in each job step is independently associated with an instance resource.

  -> When the associated task type is a job, this value is required. It specifies the instance resource association method.

* `reviewer_notification` - (Optional, List) Specifies the reviewer notification information.

  The [reviewer_notification](#message_notification_struct) structure is documented below.

* `reviewer_user_name` - (Optional, String) Specifies the reviewer nickname.

* `message_notification` - (Optional, List) Specifies the message notification information.

  The [message_notification](#message_notification_struct) structure is documented below.

* `enabled` - (Optional, String) Specifies whether the task is enabled.

<a name="trigger_time_struct"></a>
The `trigger_time` block supports:

* `time_zone` - (Required, String) Specifies the time zone. The default value is **Asia/Shanghai**.

* `policy` - (Required, String) Specifies the scheduled task execution strategy. The default value is **ONCE**.
  Values can be as follows:
  + **PERIODIC**: Periodic execution.
  + **ONCE**: Single execution.
  + **CRON**: Execute according to a CRON expression.

* `single_scheduled_time` - (Optional, Int) Specifies the execution time of a single scheduled task.

  -> This value is required if the scheduled task execution policy is **ONCE**. UTC timestamp in milliseconds.

* `periodic_scheduled_time` - (Optional, String) Specifies the daily execution time for periodic scheduled tasks.

  -> This value is required if the scheduled task execution policy is **PERIODIC**. A 24-hour time string. For example,
  if a task is to be executed at 5:30 PM on a given day, it would be **17:30:00**.

* `period` - (Optional, String) Specifies the specific week list for periodic scheduled tasks.

  -> This value is required if the scheduled task execution policy is **PERIODIC**. Days of the week are separated by
  commas. For example, Sunday is **1** and Monday is **3**. If a task is to be executed on Monday, Wednesday, Thursday,
  and Sunday, the value would be **1,2,4,5**.

* `cron` - (Optional, String) Specifies the specific CRON expression value for executing a scheduled task based on a
  CRON expression.

  -> This value is required if the scheduled task execution policy is **CRON**. A valid CRON expression is sufficient.
  For example, to execute a task at 10:15 AM every day, use `0 15 10 ? * *`.

* `scheduled_close_time` - (Optional, Int) Specifies the deadline for scheduled task execution.

  -> This value is required if the scheduled task execution policy is **PERIODIC** or **CRON**. It specifies the
  timestamp of the scheduled task rule deadline. UTC timestamp in milliseconds.

<a name="target_instances_struct"></a>
The `target_instances` block supports:

* `target_selection` - (Required, String) Specifies the target instance selection method. The default value is **MANUAL**.
  Values can be as follows:
  + **ALL**: All instances.
  + **MANUAL**: Manual selection.
  + **NONE**: No instance specified.

* `order_no` - (Required, Int) Specifies the step number.
  + When a scheduled task is associated with a script, this value is **1**.
  + When a scheduled task is associated with a job, this value indicates the step number within the job.

* `target_resource` - (Optional, List) Specifies the target instance query condition. The default value is empty.

  The [target_resource](#target_instances_target_resource_struct) structure is documented below.

  -> This field is required if the `target_selection` is **ALL**.

* `target_instances` - (Optional, String) Specifies the instance information.

  -> This field is required if the `target_selection` is **MANUAL**.

* `batch_strategy` - (Optional, String) Specifies the instance batching strategy.
  Values can be as follows:
  + **AUTO_BATCH**: Automatic batching.
  + **MANUAL_BATCH**: Manual batching.
  + **NONE**: No batching.

  -> When a scheduled task is associated with a job and the current step does not require a resource instance, this
  value must be **NONE**.
  <br>When the `target_selection` is **ALL**, this value must be **AUTO_BATCH**.
  <br>When the resource instance list is not empty, this value can only be **AUTO_BATCH** or **MANUAL_BATCH**.

* `sub_target_instances` - (Optional, List) Specifies the secondary resource instance information.

  The [sub_target_instances](#target_instances_struct) structure is documented below.

  -> If a step in a job associated with a scheduled task has multiple tasks and the `runbook_instance_mode` is **DIFF**,
  this value is required. This means that each task specifies a separate resource instance for the operation.

<a name="target_instances_target_resource_struct"></a>
The `target_resource` block supports:

* `type` - (Optional, String) Specifies the resource selection method. The default value is **REGION**.
  Values can be as follows:
  + **REGION**: Selects resource instances by region.
  + **APPLICATION**: Selects resource instances by application.

* `id` - (Optional, String) Specifies the region or application ID.

* `app_name` - (Optional, String) Specifies the application name. Separate hierarchical elements with dot(.).

  -> When the resource selection method is **APPLICATION**, this value is required and specifies the name of the
  selected application.

* `region_id` - (Optional, String) Specifies the region ID which the application is associated.

  -> When the resource selection method is **APPLICATION**, this value is required.

* `params` - (Optional, List) Specifies the dynamic query conditions for resource instances.

  The [params](#target_instances_target_resource_params_struct) structure is documented below.

<a name="target_instances_target_resource_params_struct"></a>
The `params` block supports:

* `key` - (Required, String) Specifies the resource attribute name. For example, **ep_id**, **agent_state** and so on.

* `value` - (Required, String) Specifies the parameter value corresponding to the resource attribute.

<a name="ticket_infos_struct"></a>
The `ticket_infos` block supports:

* `ticket_id` - (Optional, String) Specifies the ID of the work order related to the change management.

* `ticket_type` - (Optional, String) Specifies the work order type.
  Values can be as follows:
  + **CHANGE**: Change order.
  + **INCIDENT**: Incident order.
  + **WARROOM**: War room order.

* `target_id` - (Optional, String) Specifies the application ID associated with the work order.

* `scope_id` - (Optional, String) Specifies the region ID.

<a name="message_notification_struct"></a>
The `reviewer_notification` and `message_notification` block supports:

* `notification_endpoint_type` - (Required, String) Specifies the notification endpoint type. The default value is
  **ONCALL**.
  Values can be as follows:
  + **ONCALL**: Scheduled.
  + **USER**: Individual.

* `policy` - (Optional, String) Specifies the notification policy.
  Values can be as follows:
  + **START_EXECUTION**: Starts execution.
  + **EXECUTION_FAILED**: Execution failed.
  + **EXECUTION_SUCCEEDED**: Execution succeeded.

  -> This value is required when `enable_message_notification` is **true**.

* `schedule_scene_id` - (Optional, String) Specifies the schedule scenario ID.

  -> This value is required if the notification recipient type is **ONCALL**.

* `schedule_role_id` - (Optional, String) Specifies the schedule role ID.

  -> This value is required if the notification recipient type is **ONCALL**.

* `recipients` - (Optional, String) Specifies the ID of the notification recipient.

  -> This value is required if the notification recipient type is **USER**.

* `protocol` - (Optional, String) Specifies the notification channel. The default value is **DEFAULT**.
  Values can be as follows:
  + **DEFAULT**: Selects one of your subscribed notification channels.
  + **NONE**: No notification.
  + **SMS**: SMS.
  + **EMAIL**: Email.
  + **DINGDING**: DingTalk.
  + **LARK**: Lark.
  + **CALLNOTIFY**: Voice.
  + **WECHAT**: WeChat Work.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_user_name` - Indicates the nickname of the creator of the scheduled task.

* `approve_status` - Indicates the approval status of the scheduled task.

* `approve_comments` - Indicates the approval opinion for the scheduled task.

* `target_instances` - Indicates the target instance information.

  The [target_instances](#target_instances_out_struct) structure is documented below.

<a name="target_instances_out_struct"></a>
The `target_instances` block supports:

* `id` - Indicates the ID of target instance.

* `schedule_id` - Indicates the schedule task ID.

* `create_time` - Indicates the target instance creation time.

* `update_time` - Indicates the target instance update time.

## Import

The COC scheduled task can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_coc_scheduled_task.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `ticket_infos`, `target_instances.0.batch_strategy` and
`associated_task_enterprise_project_id`.
It is generally recommended running `terraform plan` after importing a scheduled task.
You can then decide if changes should be applied to the scheduled task, or the resource definition should be updated to
align with the scheduled task. Also you can ignore changes as below.

```hcl
resource "huaweicloud_coc_scheduled_task" "test" {
    ...

  lifecycle {
    ignore_changes = [
      ticket_infos, target_instances.0.batch_strategy, associated_task_enterprise_project_id,
    ]
  }
}
```

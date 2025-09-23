---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_schedule_task"
description: |-
  Manages a workspace APP schedule task resource within HuaweiCloud.
---

# huaweicloud_workspace_app_schedule_task

Manages a workspace APP schedule task resource within HuaweiCloud.

## Example Usage

```hcl
variable "schedule_task_name" {}
variable "scheduled_time" {}
variable "target_objects" {
  type = list(object({
    target_type = string
    target_id   = string
  }))
}

resource "huaweicloud_workspace_app_schedule_task" "test" {
  task_name      = var.schedule_task_name
  task_type      = "STOP_SERVER"
  scheduled_type = "WEEK"
  week_list      = "1,3" # Sunday, Tuesday
  scheduled_time = var.scheduled_time
  time_zone      = "Asia/Shanghai"

  dynamic "target_infos" {
    for_each = var.target_objects

    content {
      target_type = target_infos.value.target_type
      target_id   = target_infos.value.target_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the schedule task is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `task_name` - (Required, String) Specifies the name of the schedule task.  
  The name must be `1` to `64` characters, only letters, digits, and underscores (_) are allowed, and the name
  cannot contain spaces.

* `task_type` - (Required, String) Specifies the type of the schedule task.  
  The valid values are as follows:
  + **RESTART_SERVER**: Restart servers.
  + **START_SERVER**: Start servers.
  + **STOP_SERVER**: Stop servers.
  + **REINSTALL_OS**: Reinstall operating system.

* `scheduled_type` - (Required, String) Specifies the execution cycle of the schedule task.  
  The valid values are as follows:
  + **FIXED_TIME**
  + **DAY**
  + **WEEK**
  + **MONTH**

* `scheduled_time` - (Required, String) Specifies the execution time of the schedule task.  
  The format is `HH:mm:ss`.

* `target_infos` - (Required, List) Specifies the target object list of the schedule task.  
  The [target_infos](#app_schedule_task_target_infos) structure is documented below.

* `day_interval` - (Optional, Int) Specifies the execution interval of the scheduled task, in day.
  The valid value ranges from `1` to `31`.  
  This parameter is **required** when `scheduled_type` is set to **DAY**.

* `week_list` - (Optional, String) Specifies the days of week of the schedule task.  
  The valid value ranges from `1` to `7`, separated by commas, e.g. `1,2,7`.  
  `1` means Sunday, `2` means Monday, and so on.  
  This parameter is **required** when `scheduled_type` is set to **WEEK**.

* `month_list` - (Optional, String) Specifies the months of the schedule task.  
  The valid value ranges from `1` to `12`, separated by commas, e.g. `1,3,12`.  
  This parameter is **required** when `scheduled_type` is set to **MONTH**.

* `date_list` - (Optional, String) Specifies the days of month of the schedule task.  
  The valid value ranges from `1` to `31` and `L` (means the last day), separated by commas, e.g. `1,2,28` or `L`.  
  `L` can only be used alone, and cannot be used together with other values.  
  This parameter is **required** when `scheduled_type` is set to **MONTH**.

* `scheduled_date` - (Optional, String) Specifies the fixed date of the schedule task.  
  The format is `YYYY-MM-dd`.  
  This parameter is **required** when `scheduled_type` is set to **FIXED_TIME**.

* `time_zone` - (Optional, String) Specifies the time zone of the schedule task.  
  Defaults to **Asia/Shanghai**.

* `expire_time` - (Optional, String) Specifies the expiration time of the schedule task, in UTC format.

* `description` - (Optional, String) Specifies the description of the schedule task.

* `schedule_task_policy` - (Optional, List) Specifies the policy of the schedule task.
  The [schedule_task_policy](#app_schedule_task_policy) structure is documented below.

* `is_enable` - (Optional, Bool) Specifies whether to enable the schedule task.
  Defaults to **true**.

<a name="app_schedule_task_target_infos"></a>
The `target_infos` block supports:

* `target_id` - (Required, String) Specifies the ID of the target object.

* `target_type` - (Required, String) Specifies the type of the target object.  
  The valid values are as follows:
  + **SERVER**
  + **SERVER_GROUP**

<a name="app_schedule_task_policy"></a>
The `schedule_task_policy` block supports:

* `enforcement_enable` - (Required, Bool) Specifies whether to forcefully execute the task when there are
  active sessions.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the schedule task ID.

## Import

Schedule tasks can be imported using the `id`, e.g.:

```bash
$ terraform import huaweicloud_workspace_app_schedule_task.test <id>
```

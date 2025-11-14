---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_scheduled_tasks"
description: |-
  Use this data source to get the list of Workspace scheduled tasks within HuaweiCloud.
---

# huaweicloud_workspace_scheduled_tasks

Use this data source to get the list of Workspace scheduled tasks within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
data "huaweicloud_workspace_scheduled_tasks" "test" {}
```

### Filter scheduled tasks by task name

```hcl
variable "task_name" {}

data "huaweicloud_workspace_scheduled_tasks" "test" {
  task_name = var.task_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the scheduled tasks are located.  
  If omitted, the provider-level region will be used.

* `task_name` - (Optional, String) Specifies the name of the scheduled task to be queried.  
  Fuzzy match is supported.

* `task_type` - (Optional, String) Specifies the type of the scheduled task.  
  The valid values are as follows:
  + **START**
  + **STOP**
  + **REBOOT**
  + **HIBERNATE**
  + **REBUILD**
  + **EXECUTE_SCRIPT**
  + **CREATE_SNAPSHOT**

* `scheduled_type` - (Optional, String) Specifies the execution cycle type of the scheduled task.  
  The valid values are as follows:
  + **FIXED_TIME**
  + **DAY**
  + **WEEK**
  + **MONTH**
  + **LIFE_CYCLE**

* `last_status` - (Optional, String) Specifies the last execution status of the scheduled task.  
  The valid values are as follows:
  + **SUCCESS**
  + **SKIP**
  + **FAIL**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of scheduled tasks that match the filter parameters.  
  The [tasks](#workspace_scheduled_tasks_attr) structure is documented below.

<a name="workspace_scheduled_tasks_attr"></a>
The `tasks` block supports:

* `id` - The ID of the scheduled task.

* `name` - The name of the scheduled task.

* `type` - The type of the scheduled task.

* `scheduled_type` - The execution cycle type of the scheduled task.

* `life_cycle_type` - The trigger scenario type of the scheduled task.

* `last_status` - The last execution status of the scheduled task.

* `next_execution_time` - The next execution time of the scheduled task, format is **2006-01-02 15:04:05 GMT+08:00**.

* `enable` - Whether the scheduled task is enabled.

* `description` - The description of the scheduled task.

* `time_zone` - The time zone of the scheduled task.

* `priority` - The priority of the scheduled task.

  -> Only trigger task return this attribute.

* `wait_time` - The wait time after the trigger scenario for the scheduled task.

  -> Only trigger task return this attribute.

---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_scheduled_task_records"
description: |-
  Use this data source to get the execution records list of the scheduled task within HuaweiCloud.
---

# huaweicloud_workspace_scheduled_task_records

Use this data source to get the execution records list of the scheduled task within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_workspace_scheduled_task_records" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the scheduled task records are located.  
  If omitted, the provider-level region will be used.

* `task_id` - (Required, String) Specifies the ID of the scheduled task to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of scheduled task execution records.  
  The [records](#scheduled_task_records_attr) structure is documented below.

<a name="scheduled_task_records_attr"></a>
The `records` block supports:

* `id` - The ID of the scheduled task execution record.

* `start_time` - The execution time, in RFC3339 format.

* `task_type` - The type of the scheduled task.
  + **START** - Startup
  + **STOP** - Shutdown
  + **REBOOT** - Restart
  + **HIBERNATE** - Hibernate
  + **REBUILD** - Rebuild system disk

* `scheduled_type` - The execution cycle type.
  + **FIXED_TIME** - Specified time
  + **DAY** - Daily
  + **WEEK** - Weekly
  + **MONTH** - Monthly

* `life_cycle_type` - The trigger scenario type.

* `status` - The execution status of this execution.

* `success_num` - The number of successful desktops.

* `failed_num` - The number of failed desktops.

* `skip_num` - The number of skipped desktops.

* `time_zone` - The time zone information.

* `execute_task_id` - The task ID of executing the scheduled task.

* `execute_object_type` - The object type of executing the scheduled task.

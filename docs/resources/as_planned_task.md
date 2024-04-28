---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_planned_task"
description: ""
---

# huaweicloud_as_planned_task

Manages an AS planned task resource within HuaweiCloud.

## Example Usage

### AS timed planned Task

```hcl
variable "scaling_group_id" {}
variable "name" {}

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = var.scaling_group_id
  name             = var.name

  scheduled_policy {
    launch_time = "2025-11-30T12:00Z"
  }

  instance_number {
    max    = "3"
    min    = "1"
    desire = "2"
  }
}
```

### AS periodic planned Task

```hcl
variable "scaling_group_id" {}
variable "name" {}

resource "huaweicloud_as_planned_task" "test" {
  scaling_group_id = var.scaling_group_id
  name             = var.name

  scheduled_policy {
    launch_time      = "10:00"
    recurrence_type  = "WEEKLY"
    start_time       = "2025-11-30T12:00Z"
    end_time         = "2025-12-30T12:00Z"
    recurrence_value = "6"
  }

  instance_number {
    max    = "3"
    min    = "1"
    desire = "2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `scaling_group_id` - (Required, String, ForceNew) Specifies the AS scaling group ID.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of planned task.

* `scheduled_policy` - (Required, List) Specifies the scheduled or periodic policy for AS planned task.
  The [scheduled_policy](#AS_ScheduledPolicy) structure is documented below.

* `instance_number` - (Required, List) Specifies the numbers of scaling group instance for AS planned task.
  The [instance_number](#AS_InstanceNumber) structure is documented below.

<a name="AS_ScheduledPolicy"></a>
The `scheduled_policy` block supports:

* `launch_time` - (Required, String) Specifies the execution time of the AS planned task.
  + If `recurrence_type` not set or is empty, the time format is **yyyy-MM-ddTHH:mmZ**.
  + If `recurrence_type` is specified, the time format is **HH:mm**.

* `recurrence_type` - (Optional, String) Specifies the triggering type of AS planned task.
  When not set or is empty, the planned task is scheduled execution.
  After setting, the planned task is periodic execution. The valid values are as follows:
  + **DAILY**: by day periodic execution.
  + **WEEKLY**: by week periodic execution.
  + **MONTHLY**: by month periodic execution.

* `recurrence_value` - (Optional, String) Specifies the frequency at which planned task are triggered.
  Required only when `recurrence_type` is **WEEKLY** or **MONTHLY**.
  + When `recurrence_type` is **WEEKLY**, The valid value ranges from `1` to `7`.
  + When `recurrence_type` is **MONTHLY**, The valid value ranges from `1` to `31`.

* `start_time` - (Optional, String) Specifies the effective start time of the planned task.
  Only effective when `recurrence_type` is not empty.
  The time format is **yyyy-MM-ddTHH:mmZ**.
  If not set, the default is the time when the planned task is successfully created.

* `end_time` - (Optional, String) Specifies the effective end time of the planned task.
  Only effective and required when `recurrence_type` is not empty.
  The time format is **yyyy-MM-ddTHH:mmZ**.

<a name="AS_InstanceNumber"></a>
The `instance_number` block supports:

* `max` - (Optional, String) Specifies the maximum number of instances for the scaling group.

* `min` - (Optional, String) Specifies the minimum number of instances for the scaling group.

* `desire` - (Optional, String) Specifies the expected number of instances for the scaling group.

-> At least set one of `max`, `min` or `desire` parameters, at the same time, the `min` can not be
  greater than `desire` or `max`, and `desire` can not be greater than `max`. Parameters that are not set or empty,
  it means that the value of this field remains unchanged compared to the scaling group instance number.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The (UTC) creation time of the event source, in RFC3339 format.

## Import

The AS planned task can be imported using the related `scaling_group_id` and `id`, separated by a slash, e.g.

```shell
$ terraform import huaweicloud_as_planned_task.test <scaling_group_id>/<id>
```

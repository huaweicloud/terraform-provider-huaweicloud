---
subcategory: "Auto Scaling"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_as_planned_tasks"
description: |-
  Use this data source to get a list of planned tasks.
---

# huaweicloud_as_planned_tasks

Use this data source to get a list of planned tasks.

## Example Usage

```hcl
variable "scaling_group_id" {}
variable "task_id" {}

data "huaweicloud_as_planned_tasks" "test" {
  scaling_group_id = var.scaling_group_id
  task_id          = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `scaling_group_id` - (Required, String) Specifies the ID of the AS group to which the planned tasks belong.

* `task_id` - (Optional, String) Specifies the ID of the planned task.

* `name` - (Optional, String) Specifies the name of the planned task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scheduled_tasks` - All planned tasks that match the filter parameters.

  The [scheduled_tasks](#scheduled_tasks_struct) structure is documented below.

<a name="scheduled_tasks_struct"></a>
The `scheduled_tasks` block supports:

* `id` - The ID of the planned task.

* `scaling_group_id` - The ID of the AS group to which the planned task belongs.

* `name` - The name of the planned task.

* `scheduled_policy` - The planned task policy.

  The [scheduled_policy](#scheduled_tasks_scheduled_policy_struct) structure is documented below.

* `instance_number` - The instance number settings of the AS group.

  The [instance_number](#scheduled_tasks_instance_number_struct) structure is documented below.

* `created_at` - The creation time of the planned task, in RFC3339 format.

<a name="scheduled_tasks_scheduled_policy_struct"></a>
The `scheduled_policy` block supports:

* `start_time` - The start time of the valid period of the planned task, in RFC3339 format.

* `end_time` - The end time of the valid period of the planned task, in RFC3339 format.

* `launch_time` - The execute time of the planned task.
  If **recurrence_type** is left empty or null, the time format is RFC3339.
  If **recurrence_type** is specified, the time format is **HH:mm**.

* `recurrence_type` - The triggering type of planned task.

* `recurrence_value` - The frequency at which planned task are triggered.

<a name="scheduled_tasks_instance_number_struct"></a>
The `instance_number` block supports:

* `max` - The maximum number of instances in the AS group.

* `min` - The minimum number of instances in the AS group.

* `desire` - The expected number of instances in the AS group.

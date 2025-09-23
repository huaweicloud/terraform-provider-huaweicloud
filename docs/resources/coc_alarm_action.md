---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_alarm_action"
description: |-
  Manages a COC alarm action resource within HuaweiCloud.
---

# huaweicloud_coc_alarm_action

Manages a COC alarm action resource within HuaweiCloud.

~> Deleting alarm action resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "alarm_id" {}
variable "associated_task_id" {}
variable "associated_task_name" {}
variable "project_id" {}
variable "target_instances" {}

resource "huaweicloud_coc_alarm_action" "test" {
  alarm_id             = var.alarm_id
  task_type            = "SCRIPT"
  associated_task_id   = var.associated_task_id
  associated_task_type = "CUSTOMIZATION"
  associated_task_name = var.associated_task_name
  
  input_param = {
    timeout      = 300
    execute_user = "root"
    success_rate = 100
    project_id   = var.project_id
  }

  target_instances {
    target_selection = "MANUAL"
    order_no         = 0
    batch_strategy   = "NONE"
    target_instances = var.target_instances
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_id` - (Required, String, NonUpdatable) Specifies the alarm ID.

* `task_type` - (Required, String, NonUpdatable) Specifies the task type and execute specific tasks to process the alarm.
  The valid values are as follows:
  + **PLAN**: Contingency plan.
  + **SCRIPT**: Script.
  + **RUNBOOK**: Job.

* `associated_task_id` - (Required, String, NonUpdatable) Specifies the task ID.

* `associated_task_type` - (Required, String, NonUpdatable) Specifies the task type classification.
  The valid values are as follows:
  + **CUSTOMIZATION**: Custom.
  + **COMMUNAL**: Public.

* `associated_task_name` - (Required, String, NonUpdatable) Specifies the task name.

* `associated_task_enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.

* `runbook_instance_mode` - (Optional, String, NonUpdatable) Specifies the job instance mode. The default value is **SAME**.
  The valid values are as follows:
  + **SAME**: The steps are consistent.
  + **DIFF**: The steps are independent.
  + **MIXED**: Customized.

* `input_param` - (Optional, Map, NonUpdatable) Specifies the task execution parameters.

* `target_instances` - (Optional, List, NonUpdatable) Specifies the target instance information.

  The [target_instances](#target_instances_struct) structure is documented below.

* `region_id` - (Required, String, NonUpdatable) Specifies the region ID.

<a name="target_instances_struct"></a>
The `target_instances` block supports:

* `target_selection` - (Required, String, NonUpdatable) Specifies the target selection method.
  The valid values are as follows:
  + **ALL**: All instances.
  + **MANUAL**: Manually select an instance.
  + **NONE**: No instance.

* `order_no` - (Required, Int, NonUpdatable) Specifies the step number.

* `batch_strategy` - (Optional, String, NonUpdatable) Specifies the batching strategy.
  The valid values are as follows:
  + **AUTO_BATCH**: Automatic batching.
  + **MANUAL_BATCH**: Manual batching.
  + **NONE**: No batching.

* `target_instances` - (Optional, String, NonUpdatable) Specifies the instance information.

* `sub_target_instances` - (Optional, List, NonUpdatable) Specifies the incident number.

  The [sub_target_instances](#target_instances_sub_target_instances_struct) structure is documented below.

* `region_id` - (Optional, String, NonUpdatable) Specifies the substep instance target.

<a name="target_instances_sub_target_instances_struct"></a>
The `sub_target_instances` block supports:

* `target_selection` - (Required, String, NonUpdatable) Specifies the target selection method.
  The valid values are as follows:
  + **ALL**: All instances.
  + **MANUAL**: Manually select an instance.
  + **NONE**: No instance.

* `order_no` - (Required, Int, NonUpdatable) Specifies the step number.

* `batch_strategy` - (Optional, String, NonUpdatable) Specifies the batching strategy.
  The valid values are as follows:
  + **AUTO_BATCH**: Automatic batching.
  + **MANUAL_BATCH**: Manual batching.
  + **NONE**: No batching.

* `target_instances` - (Optional, String, NonUpdatable) Specifies the instance information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

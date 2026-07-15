---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_background_task_detail"
description: |-
  Use this data source to query the background task detail of DCS instances within HuaweiCloud.
---

# huaweicloud_dcs_background_task_detail

Use this data source to query the background task detail of DCS instances within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "task_id" {}

data "huaweicloud_dcs_background_task_detail" "test" {
  instance_id = var.instance_id
  task_id     = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the background task detail.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `task_id` - (Required, String) Specifies the background task ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `progress` - The overall progress, in percentage. The value ranges from **0** to **100**.

* `remain_time` - The remaining time, in seconds.

* `step_details` - The list of task step details.
  The [step_details](#dcs_background_task_detail_step_details) structure is documented below.

<a name="dcs_background_task_detail_step_details"></a>
The `step_details` block supports:

* `step_id` - The step ID.

* `step_name` - The step name.

* `step_status` - The step status. The valid values are:
  + **FINISH**: Completed.
  + **FAILED**: Failed.
  + **EXECUTING**: Executing.
  + **WAITING**: Waiting.

* `begin_time` - The step start time, in the format of 2020-06-17T07:38:42.503Z.

* `end_time` - The step end time, in the format of 2020-06-17T07:38:42.503Z.

* `error_code` - The error code.

* `sub_step_details` - The list of sub-step details.
  The [sub_step_details](#dcs_background_task_detail_sub_step_details) structure is documented below.

<a name="dcs_background_task_detail_sub_step_details"></a>
The `sub_step_details` block supports:

* `sub_step_id` - The sub-step ID.

* `sub_step_name` - The sub-step name.

* `sub_step_status` - The sub-step status. The valid values are:
  + **FINISH**: Completed.
  + **SUSPEND**: Suspended.
  + **EXECUTING**: Executing.
  + **WAITING**: Waiting.
  + **CANCELED**: Canceled.
  + **FAILED**: Failed.

* `begin_time` - The sub-step start time, in the format of 2020-06-17T07:38:42.503Z.

* `end_time` - The sub-step end time, in the format of 2020-06-17T07:38:42.503Z.

* `detail` - The additional attribute details of the sub-step.

* `error_code` - The error code.

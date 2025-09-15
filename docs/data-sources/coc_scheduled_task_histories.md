---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_scheduled_task_histories"
description: |-
  Use this data source to get the list of COC scheduled task histories.
---

# huaweicloud_coc_scheduled_task_histories

Use this data source to get the list of COC scheduled task histories.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_coc_scheduled_task_histories" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Specifies the task ID.

* `region` - (Optional, String) Specifies the task region ID.

* `status` - (Optional, String) Specifies the task execution status.
  The valid values are as follows:
  + **READY**: Pending execution.
  + **PROCESSING**: Executing.
  + **ABNORMAL**: Abnormal.
  + **PAUSED**: Paused.
  + **CANCELED**: Cancelled.
  + **FINISHED**: Successful.

* `started_start_time` - (Optional, Int) Specifies the start timestamp of the query interval for the `started_time`
  parameter. Millisecond-level UTC timestamp.

* `started_end_time` - (Optional, Int) Specifies the end timestamp of the query interval for the `started_time`
  parameter. Millisecond-level UTC timestamp.

* `finished_start_time` - (Optional, Int) Specifies the start timestamp of the query interval for the `finished_time`
  parameter. Millisecond-level UTC timestamp.

* `finished_end_time` - (Optional, Int) Specifies the end timestamp of the query interval for the `finished_time`
  parameter. Millisecond-level UTC timestamp.

* `sort_key` - (Optional, String) Specifies the sorting field name. Supported are **started_time** and **finished_time**.

* `sort_dir` - (Optional, String) Specifies the sorting method.
  The value can be **asc** or **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `scheduled_task_history_list` - Indicates the list of scheduled operation and maintenance history records.

  The [scheduled_task_history_list](#scheduled_task_history_list_struct) structure is documented below.

<a name="scheduled_task_history_list_struct"></a>
The `scheduled_task_history_list` block supports:

* `id` - Indicates the history record ID.

* `task_type` - Indicates the reference task type.

* `execution_id` - Indicates the execution ID.

* `associated_task_name` - Indicates the reference task name.

* `associated_task_name_en` - Indicates the referenced task name (in English).

* `region` - Indicates the region.

* `created_by` - Indicates the creator.

* `started_time` - Indicates the start time timestamp.

* `finished_time` - Indicates the end time timestamp.

* `status` - Indicates the status.

* `execution_msg` - Indicates the description of the execution result.

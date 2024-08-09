---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCLoud: huaweicloud_dataarts_quality_tasks"
description: |-
  Use this data source to get the list of quality tasks within HuaweiCloud.
---

# huaweicloud_dataarts_quality_tasks

Use this data source to get the list of quality tasks within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "task_name_to_be_queried" {}

data "huaweicloud_dataarts_quality_tasks" "test" {
  workspace_id = var.workspace_id

  # Filter parameter
  name = var.task_name_to_be_queried
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the quality tasks are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the quality tasks belong.

* `name` - (Optional, String) Specifies the name of the quality task.

* `category_id` - (Optional, String) Specifies the category ID to which the quality tasks belong.

* `schedule_status` - (Optional, String) Specifies the schedule status of the quality task.
  + **UNKNOWN**
  + **NOT_START**
  + **SCHEDULING**
  + **FINISH_SUCCESS**
  + **KILL**
  + **RUNNING_EXCEPTION**

* `start_time` - (Optional, String) Specifies the start time of the query interval for the most recent run time.  
  The valid format is RFC3339 format, e.g. `2024-01-01T10:00:00+08:00`

* `creator` - (Optional, String) Specifies the name of the quality task creator.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - All quality tasks that match the filter parameters.  
  The [tasks](#quality_tasks_elem) structure is documented below.

<a name="quality_tasks_elem"></a>
The `tasks` block supports:

* `id` - The API ID, in UUID format.

* `name` - The name of the quality task.

* `category_id` - The category ID to which the quality task belongs.

* `schedule_status` - The schedule status of the quality task.

* `schedule_period` - The schedule period of the quality task.
  + **MINUTE**
  + **HOUR**
  + **DAY**
  + **WEEK**

* `schedule_interval` - The schedule interval of the quality task.
  + If the `schedule_period` is **MINUTE**,**HOUR** or **DAY**, a numeric string is returned.
  + If the `schedule_period` is **WEEK**, the scheduling week information is returned, e.g. **MONDAY**, **THURSDAY**.

* `created_at` - The creation time of the quality task, in RFC3339 format.

* `last_run_time` - The last run time of the quality task, in RFC3339 format.

* `creator` - The name of the task creator.

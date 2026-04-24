---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_quality_consistency_tasks"
description: |-
  Use this data source to get the list of consistency tasks within HuaweiCloud.
---

# huaweicloud_dataarts_quality_consistency_tasks

Use this data source to get the list of consistency tasks within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "workspace_id" {}
variable "task_name" {}

data "huaweicloud_dataarts_quality_consistency_tasks" "test" {
  workspace_id = var.workspace_id
  name         = var.task_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the consistency tasks are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the consistency tasks belong.

* `category_id` - (Optional, Int) Specifies the category ID to which the consistency tasks belong.

* `name` - (Optional, String) Specifies the name of the consistency task.

* `schedule_status` - (Optional, String) Specifies the schedule status of the consistency task.
  + **UNKNOWN**
  + **NOT_START**
  + **SCHEDULING**
  + **FINISH_SUCCESS**
  + **KILL**
  + **RUNNING_EXCEPTION**

* `start_time` - (Optional, Int) Specifies the start time of the last run time query interval, in `13-digit` millisecond
  timestamp.

* `end_time` - (Optional, Int) Specifies the end time of the last run time query interval, in `13-digit` millisecond
  timestamp.

* `creator` - (Optional, String) Specifies the name of the consistency task creator.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the consistency tasks that matched filter parameters.  
  The [tasks](#dataarts_quality_consistency_tasks) structure is documented below.

<a name="dataarts_quality_consistency_tasks"></a>
The `tasks` block supports:

* `id` - The ID of the consistency task.

* `name` - The name of the consistency task.

* `category_id` - The category ID to which the consistency task belongs.

* `schedule_status` - The schedule status of the consistency task.

* `schedule_period` - The schedule period of the consistency task.
  + **MINUTE**
  + **HOUR**
  + **DAY**
  + **WEEK**

* `schedule_interval` - The schedule interval of the consistency task.

* `create_time` - The creation time of the consistency task, in RFC3339 format.

* `last_run_time` - The last run time of the consistency task, in RFC3339 format.

* `creator` - The name of the task creator.

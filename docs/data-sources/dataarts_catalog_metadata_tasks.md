---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_catalog_metadata_tasks"
description: |-
  Use this data source to query DataArts Catalog metadata tasks within HuaweiCloud.
---

# huaweicloud_dataarts_catalog_metadata_tasks

Use this data source to query DataArts Catalog metadata tasks within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_catalog_metadata_tasks" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the metadata tasks are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the metadata tasks belongs.

* `user_name` - (Optional, String) Specifies the user name which the metadata tasks are created.

* `name` - (Optional, String) Specifies the name of the metadata task.

* `data_source_type` - (Optional, String) Specifies the data source type of the metadata tasks.

* `data_connection_id` - (Optional, String) Specifies the data connection id of the metadata tasks.

* `start_time` - (Optional, String) Specifies the start time of the metadata tasks.

* `end_time` - (Optional, String) Specifies the end time of the metadata tasks.

* `directory_id` - (Optional, String) Specifies the directory ID of the metadata tasks.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metadata_tasks` - The list of metadata tasks that matched filter parameters.  
  The [metadata_tasks](#dataarts_catalog_metadata_tasks) structure is documented below.

<a name="dataarts_catalog_metadata_tasks"></a>
The `metadata_tasks` block supports:

* `id` - The ID of the metadata task, in UUID format.

* `name` - The name of the metadata task.

* `description` - The description of the metadata task.

* `user_id` - The user ID which the metadata task is created.

* `create_time` - The create time of the metadata task.

* `dir_id` - The directory ID of the metadata task.

* `schedule_config` - The dispatch information of the metadata task.  
  The [schedule_config](#dataarts_catalog_metadata_task_schedule_config_attr) structure is documented below.

* `update_time` - The latest update time of the metadata task.

* `user_name` - The user name which the metadata task is created.

* `path` - The directory path of the metadata task.

* `last_run_time` - The last run time of the metadata task.

* `start_time` - The start time of the metadata task.

* `end_time` - The end time of the metadata task.

* `next_run_time` - The next run time of the metadata task.

* `duty_person` - The duty person of the metadata task.

* `data_source_type` - The data source type of the metadata task.

* `task_config` - The config information of the metadata task, in JSON format.

<a name="dataarts_catalog_metadata_task_schedule_config_attr"></a>
The `schedule_config` block supports:

* `cron_expression` - The cron expression of the schedule task.

* `end_time` - The end time of the schedule task.

* `max_time_out` - The max time out of the schedule task.

* `interval` - The interval time of the schedule task.

* `schedule_type` - The schedule type of the schedule task.

* `start_time` - The start time of the schedule task.

* `enabled` - Whether to enable the schedule task.

* `job_id` - The job ID of the schedule task.

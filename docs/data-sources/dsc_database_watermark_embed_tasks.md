---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_database_watermark_embed_tasks"
description: |-
  Use this data source to get the list of DSC database watermark embed tasks within HuaweiCloud.
---

# huaweicloud_dsc_database_watermark_embed_tasks

Use this data source to get the list of DSC database watermark embed tasks within HuaweiCloud.

## Example Usage

### Query all database watermark embed tasks

```hcl
data "huaweicloud_dsc_database_watermark_embed_tasks" "test" {}
```

### Query database watermark embed tasks by task ID

```hcl
variable "task_id" {}

data "huaweicloud_dsc_database_watermark_embed_tasks" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the database watermark embed tasks are located.  
  If omitted, the provider-level region will be used.

* `task_id` - (Optional, String) Specifies the ID of the database watermark embed task.

* `start` - (Optional, String) Specifies the start time of the task running time interval, in RFC3339 format.  

* `end` - (Optional, String) Specifies the end time of the task running time interval, in RFC3339 format.  

* `status` - (Optional, String) Specifies the status of the database watermark embed task.  
  The valid values are as follows:
  + **RUNNING**: The task is executing.
  + **ERROR**: The task execution has failed.
  + **FINISHED**: The task has completed.
  + **WAIT**: The task is waiting to be executed.
  + **CLOSED**: The task is closed.
  
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the database watermark embed tasks.  
  The [tasks](#database_watermark_embed_tasks_tasks) structure is documented below.

<a name="database_watermark_embed_tasks_tasks"></a>
The `tasks` block supports:

* `id` - The ID of the database watermark embed task.

* `task_name` - The name of the database watermark embed task.

* `water_mark` - The watermark content.

* `db_water_param` - The database watermark embedding configuration.  
  The [db_water_param](#database_watermark_embed_tasks_db_water_param) structure is documented below.

* `selected_fields` - The selected field list used for watermark embedding.  
  The [selected_fields](#database_watermark_embed_tasks_selected_fields) structure is documented below.

* `source_db_info` - The source database information.  
  The [source_db_info](#database_watermark_embed_tasks_db_info) structure is documented below.

* `target_db_info` - The target database information.  
  The [target_db_info](#database_watermark_embed_tasks_db_info) structure is documented below.

* `watermark_describe` - The watermark description.

* `watermark_version` - The watermark algorithm version.

* `start_now` - Whether the task starts immediately.

* `start_time` - The scheduled start time of the task, in RFC3339 format.

* `schedule_switch` - Whether task scheduling is enabled.

* `schedule_type` - The schedule type of the task.

* `task_state` - The running state of the watermark embed task.

* `task_create_time` - The creation time of the task, in RFC3339 format.

* `task_end_time` - The end time of the task, in RFC3339 format.

* `water_extract_result` - The watermark extract result.

<a name="database_watermark_embed_tasks_db_water_param"></a>
The `db_water_param` block supports:

* `embed_mode` - The watermark embed mode.

* `row_spacing` - The row spacing used by fake-row watermark.

* `watermark_key` - The watermark key used to embed and extract the watermark.

* `params` - The fake-column embed parameter list.  
  The [params](#database_watermark_embed_tasks_params) structure is documented below.

<a name="database_watermark_embed_tasks_params"></a>
The `params` block supports:

* `new_column_name` - The name of the new fake column.

* `new_column_type` - The data type of the new fake column.

* `fake_strategy` - The strategy used to generate fake data.

* `fake_param` - The configuration of fake data generation.  
  The [fake_param](#database_watermark_embed_tasks_fake_param) structure is documented below.

<a name="database_watermark_embed_tasks_fake_param"></a>
The `fake_param` block supports:

* `address_accuracy` - The accuracy of generated address data.

* `date_begin` - The start date of the generated date range, in RFC3339 format.

* `date_end` - The end date of the generated date range, in RFC3339 format.

* `random_accuracy` - The precision of generated random numbers.

* `random_begin` - The lower bound of the generated random value range.

* `random_distribute` - The distribution mode of generated random values.

* `random_end` - The upper bound of the generated random value range.

* `string_distribute` - The distribution mode of generated string values.

<a name="database_watermark_embed_tasks_selected_fields"></a>
The `selected_fields` block supports:

* `column_name` - The name of the selected column.

* `column_type` - The data type of the selected column.

<a name="database_watermark_embed_tasks_db_info"></a>
The `source_db_info` and `target_db_info` blocks support:

* `db_id` - The ID of the authorized database.

* `db_name` - The name of the database.

* `db_type` - The database type.

* `ins_id` - The ID of the database instance.

* `ins_name` - The name of the database instance.

* `schema_name` - The schema name of the database.

* `table_name` - The name of the database table.

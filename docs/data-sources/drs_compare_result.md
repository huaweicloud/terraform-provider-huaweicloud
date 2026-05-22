---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_compare_result"
description: |-
  Use this data source to get the compare result for specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_compare_result

Use this data source to get the compare result for specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_compare_result" "test" { 
  job_id       = var.job_id
  current_page = 1
  per_page     = 10
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the DRS job ID.

* `current_page` - (Required, Int) Specifies the current page number for pagination.

* `per_page` - (Required, Int) Specifies the number of items per page for pagination.

* `object_level_compare_id` - (Optional, String) Specifies the object-level compare task ID.

* `line_compare_id` - (Optional, String) Specifies the line compare task ID.

* `content_compare_id` - (Optional, String) Specifies the content compare task ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `object_level_compare_results` - The object-level compare results.

  The [object_level_compare_results](#object_level_compare_results_struct) structure is documented below.

* `line_compare_results` - The line compare results.

  The [line_compare_results](#line_compare_results_struct) structure is documented below.

* `content_compare_results` - The content compare results.

  The [content_compare_results](#content_compare_results_struct) structure is documented below.

* `compare_task_list_results` - The compare task list results.

  The [compare_task_list_results](#compare_task_list_results_struct) structure is documented below.

<a name="object_level_compare_results_struct"></a>
The `object_level_compare_results` block supports:

* `compare_task_id` - The object-level compare task ID.

* `object_compare_overview` - The object compare overview.

  The [object_compare_overview](#object_compare_overview_struct) structure is documented below.

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="object_compare_overview_struct"></a>
The `object_compare_overview` block supports:

* `object_type` - The object type.
  The valid values are as follows:
  + **DB**: Database.
  + **TABLE**: Table.
  + **VIEW**: View.
  + **EVENT**: Event.
  + **ROUTINE**: Stored procedure and function.
  + **INDEX**: Index.
  + **TRIGGER**: Trigger.
  + **SYNONYM**: Synonym.
  + **FUNCTION**: Function.
  + **PROCEDURE**: Stored procedure.
  + **TYPE**: Custom type.
  + **RULE**: Rule.
  + **DEFAULT_TYPE**: Default value.
  + **PLAN_GUIDE**: Execution plan.
  + **CONSTRAINT**: Constraint.
  + **FILE_GROUP**: File group.
  + **PARTITION_FUNCTION**: Partition function.
  + **PARTITION_SCHEME**: Partition scheme.
  + **TABLE_COLLATION**: Table collation.
  + **EXTENSIONS**: Plugin.

* `object_compare_result` - The compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

* `target_count` - The number of objects in the target database.

* `source_count` - The number of objects in the source database.

* `diff_count` - The number of differences between source and target databases.

<a name="line_compare_results_struct"></a>
The `line_compare_results` block supports:

* `compare_task_id` - The line compare task ID.

* `line_compare_overview` - The line compare overview.

  The [line_compare_overview](#line_compare_overview_struct) structure is documented below.

* `line_compare_overview_count` - The total number of line compare overview results.

* `line_compare_details` - The line compare details.

  The [line_compare_details](#line_compare_details_struct) structure is documented below.

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="line_compare_overview_struct"></a>
The `line_compare_overview` block supports:

* `source_db_name` - The source database name.

* `target_db_name` - The target database name.

* `line_compare_result` - The compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

<a name="line_compare_details_struct"></a>
The `line_compare_details` block supports:

* `source_db_name` - The source database name.

* `line_compare_detail` - The line compare detail.

  The [line_compare_detail](#line_compare_detail_struct) structure is documented below.

* `line_compare_detail_count` - The total number of line compare detail results.

<a name="line_compare_detail_struct"></a>
The `line_compare_detail` block supports:

* `source_table_name` - The source table name.

* `target_table_name` - The target table name.

* `source_row_num` - The number of rows in the source table.

* `target_row_num` - The number of rows in the target table.

* `diff_row_num` - The number of different rows between source and target tables.

* `line_compare_result` - The compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

* `message` - The additional information.

<a name="content_compare_results_struct"></a>
The `content_compare_results` block supports:

* `compare_task_id` - The content compare task ID.

* `content_compare_overview` - The content compare overview.

  The [content_compare_overview](#content_compare_overview_struct) structure is documented below.

* `content_compare_overview_count` - The total number of content compare overview results.

* `content_compare_details` - The content compare details.

  The [content_compare_details](#content_compare_details_struct) structure is documented below.

* `content_compare_diffs` - The content compare differences.

  The [content_compare_diffs](#content_compare_diffs_struct) structure is documented below.

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="content_compare_overview_struct"></a>
The `content_compare_overview` block supports:

* `source_db_name` - The source database name.

* `target_db_name` - The target database name.

* `content_compare_result` - The compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

<a name="content_compare_details_struct"></a>
The `content_compare_details` block supports:

* `source_db_name` - The source database name.

* `content_compare_detail` - The content compare detail.

  The [content_compare_detail](#content_compare_detail_struct) structure is documented below.

* `content_compare_detail_count` - The total number of content compare detail results.

* `content_uncompare_detail` - The content uncompare detail (tables that cannot be compared).

  The [content_compare_detail](#content_compare_detail_struct) structure is documented below.

* `content_uncompare_detail_count` - The total number of content uncompare detail results.

<a name="content_compare_detail_struct"></a>
The `content_compare_detail` block supports:

* `source_db_name` - The source database name.

* `target_db_name` - The target database name.

* `source_table_name` - The source table name.

* `target_table_name` - The target table name.

* `source_row_num` - The number of rows in the source table.

* `target_row_num` - The number of rows in the target table.

* `diff_row_num` - The number of different rows between source and target tables.

* `line_compare_result` - The line compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

* `content_compare_result` - The content compare result.
  The valid values are as follows:
  + **CONSISTENT**: Consistent.
  + **INCONSISTENT**: Inconsistent.
  + **COMPARING**: Comparing.
  + **WAITING_FOR_COMPARISON**: Waiting for comparison.
  + **FAILED_TO_COMPARE**: Compare failed.
  + **TARGET_DB_NOT_EXIT**: Target database does not exist.
  + **CAN_NOT_COMPARE**: Cannot compare.

* `message` - The additional information.

<a name="content_compare_diffs_struct"></a>
The `content_compare_diffs` block supports:

* `source_db_name` - The source database name.

* `source_table_name` - The source table name.

* `content_compare_diff` - The content compare difference.

  The [content_compare_diff](#content_compare_diff_struct) structure is documented below.

* `content_compare_diff_count` - The total number of content compare differences.

<a name="content_compare_diff_struct"></a>
The `content_compare_diff` block supports:

* `target_select_sql` - The SQL for querying the target database.

* `source_select_sql` - The SQL for querying the source database.

* `source_key_value` - The list of source key values.

* `target_key_value` - The list of target key values.

<a name="compare_task_list_results_struct"></a>
The `compare_task_list_results` block supports:

* `compare_task_list` - The compare task list.

  The [compare_task_list](#compare_task_list_struct) structure is documented below.

* `compare_task_list_count` - The total number of compare tasks.

* `error_msg` - The error message.

* `error_code` - The error code.

<a name="compare_task_list_struct"></a>
The `compare_task_list` block supports:

* `compare_task_id` - The compare task ID.

* `compare_type` - The compare task type.
  The valid values are as follows:
  + **lines**: Line comparison.
  + **contents**: Value comparison.
  + **object_comparison**: Object-level comparison.
  + **account**: Account comparison.
  + **random**: Sampling comparison.
  + **node**: Kernel calculation comparison result.
  + **mgr**: Management calculation comparison result.

* `compare_task_status` - The compare task status.
  The valid values are as follows:
  + **RUNNING**: Running.
  + **WAITING_FOR_RUNNING**: Waiting to start.
  + **SUCCESSFUL**: Completed.
  + **FAILED**: Failed.
  + **CANCELLED**: Cancelled.
  + **TIMEOUT_INTERRUPT**: Timeout interrupted.
  + **FULL_DOING**: Full verification in progress.
  + **INCRE_DOING**: Incremental verification in progress.

* `create_time` - The compare start time.

* `end_time` - The compare end time.

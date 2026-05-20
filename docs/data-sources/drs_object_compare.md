---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_object_compare"
description: |-
  Use this data source to get the object comparison result of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_object_compare

Use this data source to get the object comparison result of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_object_compare" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `create_time` - The creation time of the comparison task, UTC time, for example: **2024-04-09T07:00:57Z**.

* `start_time` - The start time of the comparison task, UTC time, for example: **2024-04-09T07:00:57Z**.

* `status` - The status of the comparison task.
  The valid values are as follows:
  + **RUNNING**: Running.
  + **WAITING_FOR_RUNNING**: Waiting to start.
  + **SUCCESSFUL**: Completed.
  + **FAILED**: Failed.
  + **CANCELLED**: Cancelled.
  + **TIMEOUT_INTERRUPT**: Timeout interrupted.
  + **FULL_DOING**: Full verification in progress.
  + **INCRE_DOING**: Incremental verification in progress.

* `export_status` - The status of generating the comparison result report file.
  The valid values are as follows:
  + **INIT**: Initial state.
  + **EXPORTING**: Comparison result exporting.
  + **EXPORT_COMPLETE**: Comparison result export completed.
  + **EXPORT_COMMON_FAILED**: Comparison result export failed.

* `report_remain_seconds` - The remaining validity time of the comparison result report file, in seconds.
  Returns `-1` when the report is not generated.

* `compare_job_id` - The comparison task ID.

* `error_msg` - The failure reason.

* `compare_result` - The comparison results.

  The [compare_result](#compare_result_struct) structure is documented below.

* `database_info` - The business/disaster recovery database information during the comparison task.

  The [database_info](#database_info_struct) structure is documented below.

<a name="compare_result_struct"></a>
The `compare_result` block supports:

* `type` - The object type.
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

* `source_count` - The number of objects of this type in the source database.

* `target_count` - The number of objects of this type in the target database.

* `status` - The comparison result. `0` for inconsistent, `2` for consistent, `3` for incomplete.

<a name="database_info_struct"></a>
The `database_info` block supports:

* `service_database` - The business database information.

* `disaster_recovery_database` - The disaster recovery database information.

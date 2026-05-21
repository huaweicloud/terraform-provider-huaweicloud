---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_table_compare"
description: |-
  Use this data source to get the table comparison result of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_table_compare

Use this data source to get the table comparison result of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_table_compare" "test" {
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

* `compare_jobs` - The table comparison task information.

  The [compare_jobs](#compare_jobs_struct) structure is documented below.

<a name="compare_jobs_struct"></a>
The `compare_jobs` block supports:

* `id` - The comparison task ID.

* `type` - The comparison type.
  The valid values are as follows:
  + **lines**: Row count comparison.
  + **contents**: Content comparison.
  + **random**: Sample comparison, currently only supports GaussDB distributed to GaussDB distributed, GaussDB
    distributed to PostgreSQL, GaussDB centralized to PostgreSQL synchronization links.

* `options` - The comparison configuration items in key-value format.
  Content comparison supports the following configuration items:
  + Comparison method configuration, key: **contentCompareType**, value: **dynamic** for dynamic comparison, **static**
    for static comparison.
  + LOB field comparison type configuration, key: **lobCompare**, value: **ignore** for ignoring, **length** for length
    comparison. Row count comparison supports the following configuration items:
  + Comparison strategy configuration, applicable for multi-table merging, key: **comparePolicy**, value: **normal** for
    normal comparison, **manyToOne** for many-to-one comparison.

* `start_time` - The start time of the comparison task, UTC time, for example: **2024-04-09T18:50:20Z**.

* `end_time` - The end time of the comparison task, UTC time, for example: **2024-04-09T18:50:20Z**.

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

* `compare_job_tag` - The tags of the comparison task, currently only returned when the comparison strategy is involved.

* `proportion_value` - The sampling proportion, filled in when the comparison type is sample comparison.

* `database_info` - The business/disaster recovery database information during the comparison task.

  The [database_info](#database_info_struct) structure is documented below.

<a name="database_info_struct"></a>
The `database_info` block supports:

* `service_database` - The business database information.

* `disaster_recovery_database` - The disaster recovery database information.

---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_top_sqls"
description: |-
  Use this data source to get the top SQLs of RDS instance.
---

# huaweicloud_rds_top_sqls

Use this data source to get the top SQLs of RDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_top_sqls" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `sort_key` - (Optional, String) Specifies the sort key.
  Value options:
  + **avg_cpu_time**: average CPU time
  + **total_cpu_time**: total cpu time
  + **total_duration_time**: total duration time
  + **avg_duration_time**: average duration time
  + **total_rows**: total rows
  + **avg_rows**: average rows
  + **total_logical_reads**: total logical reads
  + **avg_logical_reads**: average logical reads

* `limit` - (Optional, Int) Specifies the limit of the top SQL. The max value is **15**.

* `statement` - (Optional, String) Specifies the statement.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  Value options: **desc**, **asc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `avg_cpu_time_top_values` - Indicates the top SQL list of average CPU time.

  The [avg_cpu_time_top_values](#avg_cpu_time_top_values_struct) structure is documented below.

* `total_cpu_time_top_values` - Indicates the top SQL list of total CPU time.

  The [total_cpu_time_top_values](#total_cpu_time_top_values_struct) structure is documented below.

* `total_logical_reads_top_values` - Indicates the top SQL list of total logical reads.

  The [total_logical_reads_top_values](#total_logical_reads_top_values_struct) structure is documented below.

* `avg_duration_time_top_values` - Indicates the top SQL list of total logical reads.

  The [avg_duration_time_top_values](#avg_duration_time_top_values_struct) structure is documented below.

* `avg_rows_top_values` - Indicates the top SQL list of average rows.

  The [avg_rows_top_values](#avg_rows_top_values_struct) structure is documented below.

* `avg_logical_top_values` - Indicates the top SQL list of average logical.

  The [avg_logical_top_values](#avg_logical_top_values_struct) structure is documented below.

* `total_duration_time_top_values` - Indicates the top SQL list of total duration time.

  The [total_duration_time_top_values](#total_duration_time_top_values_struct) structure is documented below.

* `total_rows_top_values` - Indicates the top SQL list of total rows.

  The [total_rows_top_values](#total_rows_top_values_struct) structure is documented below.

* `list` - Indicates the top SQL list.

  The [list](#list_struct) structure is documented below.

<a name="avg_cpu_time_top_values_struct"></a>
The `avg_cpu_time_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgWorkerTime**: average CPU time
  + **AvgDuration**: average duration
  + **TotalWorkerTime**: total CPU time
  + **TotalDuration**: total duration

* `value` - Indicates the time consume, unit is **ms**.

<a name="total_cpu_time_top_values_struct"></a>
The `total_cpu_time_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgWorkerTime**: average CPU time
  + **AvgDuration**: average duration
  + **TotalWorkerTime**: total CPU time
  + **TotalDuration**: total duration

* `value` - Indicates the time consume, unit is **ms**.

<a name="total_logical_reads_top_values_struct"></a>
The `total_logical_reads_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgLogicalReads**: average logical reads
  + **TotalLogicalReads**: total logical reads

* `value` - Indicates the logical read consumption.

<a name="avg_duration_time_top_values_struct"></a>
The `avg_duration_time_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgWorkerTime**: average CPU time
  + **AvgDuration**: average duration
  + **TotalWorkerTime**: total CPU time
  + **TotalDuration**: total duration

* `value` - Indicates the time consume, unit is **ms**.

<a name="avg_rows_top_values_struct"></a>
The `avg_rows_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgReturnRows**: average return rows
  + **TotalReturnRows**: total  return rows

* `value` - Indicates the row number.

<a name="avg_logical_top_values_struct"></a>
The `avg_logical_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgLogicalReads**: average logical reads
  + **TotalLogicalReads**: total logical reads

* `value` - Indicates the logical read consumption.

<a name="total_duration_time_top_values_struct"></a>
The `total_duration_time_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgWorkerTime**: average CPU time
  + **AvgDuration**: average duration
  + **TotalWorkerTime**: total CPU time
  + **TotalDuration**: total duration

* `value` - Indicates the time consume, unit is **ms**.

<a name="total_rows_top_values_struct"></a>
The `total_rows_top_values` block supports:

* `id` - Indicates the ID.

* `data_type` - Indicates the data type.
  The value can be:
  + **AvgReturnRows**: average return rows
  + **TotalReturnRows**: total  return rows

* `value` - Indicates the row number.

<a name="list_struct"></a>
The `list` block supports:

* `id` - Indicates the ID.

* `avg_logical_reads` - Indicates the average logical reads.

* `query` - Indicates the SQL full text.

* `execution_count` - Indicates the execution count.

* `avg_duration_time_percent` - Indicates the average duration time percent.

* `avg_rows` - Indicates the average row number.

* `avg_physical_reads` - Indicates the average physical reads.

* `total_duration_time_percent` - Indicates the total duration time percent.

* `total_rows_percent` - Indicates the total rows percent.

* `total_logical_reads_percent` - Indicates the total logical reads percent.

* `avg_logical_reads_percent` - Indicates the average logical reads percent.

* `avg_logical_write_percent` - Indicates the average logical write percent.

* `last_execution_time` - Indicates the last execution time.

* `total_cpu_time` - Indicates the total cpu time, unit is **ms**.

* `avg_rows_percent` - Indicates the average rows percent.

* `total_logical_write` - Indicates the total logical write.

* `total_physical_reads` - Indicates the total physical reads.

* `avg_cpu_time` - Indicates the average cpu time, unit is **ms**.

* `total_physical_reads_percent` - Indicates the total physical reads percent.

* `avg_cpu_time_percent` - Indicates the average cpu time percent.

* `statement` - Indicates the statement of the SQL.

* `avg_duration_time` - Indicates the average duration time.

* `total_rows` - Indicates the total row number.

* `avg_physical_reads_percent` - Indicates the average physical reads percent.

* `total_duration_time` - Indicates the total duration time.

* `avg_logical_write` - Indicates the average logical write.

* `total_logical_write_percent` - Indicates the total logical write percent.

* `db_name` - Indicates the database name.

* `execution_count_percent` - Indicates the execution count percent.

* `total_cpu_time_percent` - Indicates the total cpu time percent.

* `total_logical_reads` - Indicates the total logical reads.

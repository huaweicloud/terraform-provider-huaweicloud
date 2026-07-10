---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_full_sql_statistics"
description: |-
  Use this data source to get the list of full SQL statistics of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_full_sql_statistics

Use this data source to get the list of full SQL statistics of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}

data "huaweicloud_gaussdb_full_sql_statistics" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
}
```

### Query with component ID

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}
variable "component_id" {}

data "huaweicloud_gaussdb_full_sql_statistics" "test" {
  instance_id  = var.instance_id
  begin_time   = var.begin_time
  end_time     = var.end_time
  component_id = var.component_id
}
```

### Query with Multi Queries

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}

data "huaweicloud_gaussdb_full_sql_statistics" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time

  multi_queries {
    name      = "query"
    condition = "OR"
    values    = ["select"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `begin_time` - (Required, String) Specifies the query start time.
  The value must use the format **yyyy-mm-ddThh:mm:ssZ** and conform to the ISO 8601 UTC standard. T is the separator
  between the calendar and the hourly notation of time. Z indicates the time zone offset.
  For example, in the Beijing time zone, the time zone offset is shown as `+0800`, the value is
  **2026-06-01T15:20:00+0800**.

* `end_time` - (Required, String) Specifies the query end time.
  The value must use the format **yyyy-mm-ddThh:mm:ssZ** and conform to the ISO 8601 UTC standard. T is the separator
  between the calendar and the hourly notation of time. Z indicates the time zone offset.
  For example, in the Beijing time zone, the time zone offset is shown as `+0800`, the value is
  **2026-06-01T15:20:00+0800**.

-> The time range from `begin_time` to `end_time` cannot exceed `30` days.

* `node_id` - (Optional, String) Specifies the node ID.

* `component_id` - (Optional, String) Specifies the component ID.

* `query` - (Optional, String) Specifies the SQL text.

* `sql_id` - (Optional, String) Specifies the normalized SQL ID.

* `sql_exec_id` - (Optional, String) Specifies the unique SQL ID.

* `transaction_id` - (Optional, String) Specifies the transaction ID.

* `trace_id` - (Optional, String) Specifies the trace ID.

* `db_name` - (Optional, String) Specifies the database name.

* `schema_name` - (Optional, String) Specifies the schema name.

* `username` - (Optional, String) Specifies the user name.

* `client_addr` - (Optional, String) Specifies the client address.

* `client_port` - (Optional, String) Specifies the client port.

* `application_name` - (Optional, String) Specifies the application name.

* `is_slow_sql` - (Optional, Bool) Specifies whether the SQL statement is slow.
  The valid values are as follows:
  + **true**: The SQL statement is slow.
  + **false**: The SQL statement is not slow.

* `order_by` - (Optional, String) Specifies the sorting field.
  The valid values are as follows:
  + **sql_id**
  + **sql_count** (Defaults)
  + **avg_db_time**
  + **avg_cpu_time**
  + **avg_execution_time**
  + **avg_data_io_time**
  + **start_time_stamp**

* `order` - (Optional, String) Specifies the sorting order.
  The valid values are **DESC** and **ASC**. Defaults to **DESC**.

* `multi_queries` - (Optional, List) Specifies the query conditions for field aggregation.
  The [multi_queries](#multi_queries_struct) structure is documented below.

* `compare_conditions` - (Optional, List) Specifies the list of multi-query conditions.
  The [compare_conditions](#compare_conditions_struct) structure is documented below.

<a name="multi_queries_struct"></a>
The `multi_queries` block supports:

* `name` - (Required, String) Specifies the query field name. The value can be **query**.

* `condition` - (Required, String) Specifies the merge condition. The valid values are **and**, **or**, **AND**, **OR**.

* `values` - (Required, List) Specifies the list of filter values. A list of `1` to `5` strings.

* `is_fuzzy` - (Optional, String) Specifies whether to use fuzzy query.
  The valid values are as follows:
  + **true**: Fuzzy query. (Default)
  + **false**: Exact match.

<a name="compare_conditions_struct"></a>
The `compare_conditions` block supports:

* `name` - (Optional, String) Specifies the query field name.
  The valid values are as follows:
  + **total_sql_time**
  + **sql_count**

* `enable_equal` - (Optional, String) Specifies whether to include the equal to condition.
  The valid values are as follows:
  + **true**: The range boundary values are included.
  + **false**: The range boundary values are not included.

* `min` - (Optional, String) Specifies the value for evaluating the minimum threshold.
  The valid value ranges from `0` to `2^63-1`.

* `max` - (Optional, String) Specifies the value for evaluating the maximum threshold.
  The valid value ranges from `0` to `2^63-1`.

* `value` - (Optional, String) Specifies the value for evaluating the equality threshold.
  The valid value ranges from `0` to `2^63-1`.
  The parameter `value` has the highest priority. If `value` is not left blank, the settings
  for `min` and `max` are ignored. If `value` is left blank, the `min` and `max` filtering conditions are enabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics` - The list of full SQL statistics.
  The [statistics](#full_sql_statistics) structure is documented below.

<a name="full_sql_statistics"></a>
The `statistics` block supports:

* `template` - The SQL template.

* `sql_id` - The normalized SQL ID.

* `sql_count` - The total number of aggregated SQL entries.

* `total_sql_time` - The total SQL elapsed time, in microseconds.

* `avg_sql_time` - The average SQL elapsed time, in microseconds.

* `total_db_time` - The total valid DB elapsed time, in microseconds.

* `avg_db_time` - The average valid DB elapsed time, in microseconds.

* `total_cpu_time` - The total CPU elapsed time, in microseconds.

* `avg_cpu_time` - The average CPU execution elapsed time, in microseconds.

* `avg_execution_time` - The average execution time within the executor, in microseconds.

* `avg_parse_time` - The average parser time, in microseconds.

* `avg_plan_time` - The average execution plan time, in microseconds.

* `total_data_io_time` - The total IO elapsed time, in microseconds.

* `avg_data_io_time` - The average IO elapsed time, in microseconds.

* `avg_n_blocks_hit` - The average number of buffer block hits.

* `avg_n_returned_rows` - The average number of returned rows.

* `avg_n_tuples_fetched` - The average number of scanned rows.

* `start_time_stamp` - The start timestamp, in milliseconds.

* `end_time_stamp` - The end timestamp, in milliseconds.

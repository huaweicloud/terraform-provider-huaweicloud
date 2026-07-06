---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_single_full_sqls"
description: |-
  Use this data source to query full data of a single SQL statement within HuaweiCloud.
---

# huaweicloud_gaussdb_single_full_sqls

Use this data source to query full data of a single SQL statement within HuaweiCloud.

## Example Usage

### Query SQL data

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}

data "huaweicloud_gaussdb_single_full_sqls" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
}
```

### Query SQL data by SQL ID

```hcl
variable "instance_id" {}
variable "begin_time" {}
variable "end_time" {}
variable "sql_id" {}

data "huaweicloud_gaussdb_single_full_sqls" "test" {
  instance_id = var.instance_id
  begin_time  = var.begin_time
  end_time    = var.end_time
  sql_id      = var.sql_id
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

* `query` - (Optional, String) Specifies the SQL text.

* `sql_id` - (Optional, String) Specifies the normalized SQL ID.

* `sql_exec_id` - (Optional, String) Specifies the unique SQL statement ID.

* `transaction_id` - (Optional, String) Specifies the transaction ID.

* `trace_id` - (Optional, String) Specifies the link ID.

* `db_name` - (Optional, String) Specifies the database name.

* `schema_name` - (Optional, String) Specifies the schema name.

* `username` - (Optional, String) Specifies the database user name.

* `client_addr` - (Optional, String) Specifies the client address.

* `client_port` - (Optional, String) Specifies the client port.

* `is_slow_sql` - (Optional, Bool) Specifies whether the SQL statement is slow.
  The valid values are as follows:
  + **true**: The SQL statement is slow.
  + **false**: The SQL statement is not slow.

* `order_by` - (Optional, String) Specifies the sorting field.
  The value can be **begin_time**.

* `order` - (Optional, String) Specifies the sorting mode.
  The valid values are as follows:
  + **DESC**: Descending order.
  + **ASC**: Ascending order.

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
  + **db_time**: The valid DB time, in ms.
  + **cpu_time**: The CPU execution time, in ms.
  + **data_io_time**: The I/O execution time, in ms.
  + **execution_time**: The execution time in the executor, in ms.

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

* `full_sqls` - Indicates the records of a single SQL statement.

  The [full_sqls](#full_sqls_struct) structure is documented below.

<a name="full_sqls_struct"></a>
The `full_sqls` block supports:

* `id` - Indicates the unique key ID of a SQL statement record.

* `instance_id` - Indicates the instance ID.

* `node_id` - Indicates the node ID.

* `component_id` - Indicates the component ID.

* `db_name` - Indicates the database name.

* `schema_name` - Indicates the schema name.

* `username` - Indicates the user name.

* `application_name` - Indicates the name of the application that sends a request.

* `client_addr` - Indicates the IP address of the client that initiated the request.

* `client_port` - Indicates the port number of the client that initiated the request.

* `sql_id` - Indicates the normalized SQL ID.

* `sql_exec_id` - Indicates the unique SQL ID.

* `transaction_id` - Indicates the transaction ID.

* `trace_id` - Indicates the link ID.

* `query` - Indicates the normalized SQL statement.

* `sql` - Indicates the original SQL text after parsing.

* `begin_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `all_time` - Indicates the total execution time, in μs.

* `db_time` - Indicates the valid DB time, in μs.

* `cpu_time` - Indicates the CPU time, in μs.

* `data_io_time` - Indicates the I/O time, in μs.

* `execution_time` - Indicates the execution time in the executor, in μs.

* `scan_lines` - Indicates the scanned rows.

* `insert_rows` - Indicates the number of rows inserted.

* `update_rows` - Indicates the number of rows updated.

* `delete_rows` - Indicates the number of rows deleted.

* `is_slow_sql` - Whether the SQL statement is slow.

* `start_timestamp` - Indicates the start time of SQL statement execution. The value is a 13-digit standard timestamp.

* `finish_timestamp` - Indicates the end time of SQL statement execution. The value is a 13-digit standard timestamp.

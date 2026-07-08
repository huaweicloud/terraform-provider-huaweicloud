---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_single_sql_details"
description: |-
  Use this data source to query the details of a single SQL statement of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_single_sql_details

Use this data source to query the details of a single SQL statement of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "sql_exec_id" {}

data "huaweicloud_gaussdb_single_sql_details" "test" {
  instance_id = var.instance_id
  sql_exec_id = var.sql_exec_id
}
```

### Query with Optional Filters

```hcl
variable "instance_id" {}
variable "sql_exec_id" {}
variable "sql_id" {}

data "huaweicloud_gaussdb_single_sql_details" "test" {
  instance_id = var.instance_id
  sql_exec_id = var.sql_exec_id
  sql_id      = var.sql_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `sql_exec_id` - (Required, String) Specifies the unique SQL ID.

* `key_id` - (Optional, Int) Specifies the unique key ID of the collected full SQL record.

* `sql_id` - (Optional, String) Specifies the normalized SQL ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `trace_statistics` - The trace statistics of the SQL statement.
  The [trace_statistics](#trace_statistics_struct) structure is documented below.

* `components` - The list of SQL execution records on components.
  The [components](#components_struct) structure is documented below.

<a name="trace_statistics_struct"></a>
The `trace_statistics` block supports:

* `hit_rate` - The hit rate.

* `db_time` - The valid DB time, in microseconds.

* `cpu_time` - The CPU time, in microseconds.

* `io_time` - The IO time, in microseconds.

* `execution_time` - The execution time within the executor, in microseconds.

* `scan_rows` - The number of scanned rows.

* `update_rows` - The number of updated rows.

* `insert_rows` - The number of inserted rows.

* `delete_rows` - The number of deleted rows.

<a name="components_struct"></a>
The `components` block supports:

* `component_id` - The component ID.

* `db_name` - The database name.

* `schema_name` - The schema name.

* `origin_node` - The original node.

* `username` - The user name.

* `application_name` - The application name of the user request.

* `client_addr` - The client address of the user request.

* `client_port` - The client port of the user request.

* `parent_sql_id` - The normalized SQL ID of the outer SQL statement.

* `sql_id` - The normalized SQL ID.

* `sql_exec_id` - The unique SQL ID.

* `transaction_id` - The transaction ID.

* `trace_id` - The trace ID.

* `query` - The normalized SQL.

* `sql` - The parsed original SQL text.

* `thread_id` - The thread ID.

* `session_id` - The session ID.

* `start_time` - The start time, in the format of **yyyy-mm-ddThh:mm:ss.SSSSSZ**.

* `finish_time` - The finish time, in the format of **yyyy-mm-ddThh:mm:ss.SSSSSZ**.

* `slow_query_threshold` - The slow SQL threshold.

* `n_soft_parse` - The number of soft parses.

* `n_hard_parse` - The number of hard parses.

* `query_plan` - The execution plan.

* `n_returned_rows` - The number of rows in the result set returned by the SELECT statement.

* `n_tuples_fetched` - The number of randomly scanned rows.

* `n_tuples_returned` - The number of sequentially scanned rows.

* `n_tuples_inserted` - The number of inserted rows.

* `n_tuples_updated` - The number of updated rows.

* `n_tuples_deleted` - The number of deleted rows.

* `n_blocks_fetched` - The number of buffer block accesses.

* `n_blocks_hit` - The number of buffer block hits.

* `db_time` - The valid DB time, in microseconds.

* `cpu_time` - The CPU time, in microseconds.

* `execution_time` - The execution time within the executor, in microseconds.

* `parse_time` - The SQL parsing time, in microseconds.

* `plan_time` - The execution time within the executor, in microseconds.

* `rewrite_time` - The SQL rewriting time, in microseconds.

* `pl_execution_time` - The execution time on PL/pgSQL, in microseconds.

* `pl_compilation_time` - The compilation time on PL/pgSQL, in microseconds.

* `data_io_time` - The IO time, in microseconds.

* `lock_count` - The number of locks.

* `lock_time` - The lock duration.

* `lock_wait_count` - The number of lock waits.

* `lock_wait_time` - The lock wait time.

* `details` - The detailed list.

* `is_slow_sql` - Whether the SQL is a slow SQL.

* `advise` - The risk information that may cause the SQL to be a slow SQL.

* `finish_status` - The statement completion status.

* `net_send_info` - The network status of messages sent through the physical connection.

* `net_recv_info` - The network status of messages received through the physical connection.

* `net_stream_send_info` - The network status of messages sent through the logical connection.

* `net_stream_recv_info` - The network status of messages received through the logical connection.

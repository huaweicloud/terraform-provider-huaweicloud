---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_traces"
description: |-
  Use this data source to query the SQL link information of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_traces

Use this data source to query the SQL link information of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "sql_id" {}

data "huaweicloud_gaussdb_sql_traces" "test" {
  instance_id = var.instance_id
  sql_id      = var.sql_id
}
```

### Query by SQL Execution ID

```hcl
variable "instance_id" {}
variable "sql_exec_id" {}

data "huaweicloud_gaussdb_sql_traces" "test" {
  instance_id = var.instance_id
  sql_exec_id = var.sql_exec_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `sql_id` - (Optional, String) Specifies the normalized SQL ID.

* `sql_exec_id` - (Optional, String) Specifies the unique SQL ID.

* `transaction_id` - (Optional, String) Specifies the transaction ID.

* `trace_id` - (Optional, String) Specifies the trace ID.

-> At least one of `sql_id`, `sql_exec_id`, `transaction_id`, or `trace_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `traces` - The list of SQL link node execution information.
  The [traces](#sql_traces_struct) structure is documented below.

<a name="sql_traces_struct"></a>
The `traces` block supports:

* `component_id` - The component ID.

* `node_id` - The node ID.

* `transaction_id` - The transaction ID.

* `sql_id` - The normalized SQL ID.

* `sql_exec_id` - The unique SQL ID.

* `db_name` - The database name.

* `schema_name` - The schema name.

* `start_time` - The statement start time.

* `finish_time` - The statement finish time.

* `all_time` - The total execution elapsed time, in microseconds.

* `user_name` - The user name.

* `client_addr` - The client address of the user request.

* `client_port` - The client port of the user request.

* `trace_id` - The trace ID passed in by the driver.

* `application_name` - The application name of the user request.

* `session_id` - The user session ID.

* `is_slow_sql` - Whether the SQL is a slow SQL.

* `execution_time_details` - The execution time details.
  The [execution_time_details](#traces_execution_time_details_struct) structure is documented below.

<a name="traces_execution_time_details_struct"></a>
The `execution_time_details` block supports:

* `resource_time` - The resource elapsed time information.
  The [resource_time](#traces_resource_time_struct) structure is documented below.

* `kernel_time` - The kernel module elapsed time information.
  The [kernel_time](#traces_kernel_time_struct) structure is documented below.

* `kernel_execution_time` - The kernel execution module elapsed time information.
  The [kernel_execution_time](#traces_kernel_execution_time_struct) structure is documented below.

* `wait_event_time` - The wait event and statement lock event elapsed time information.
  The [wait_event_time](#traces_wait_event_time_struct) structure is documented below.

<a name="traces_resource_time_struct"></a>
The `resource_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `resource_time_details` - The resource elapsed time details.
  The [resource_time_details](#resource_time_details_struct) structure is documented below.

<a name="resource_time_details_struct"></a>
The `resource_time_details` block supports:

* `cpu_time` - The CPU time, in microseconds.

* `data_io_time` - The IO time, in microseconds.

* `other_time` - The other elapsed time, in microseconds.

<a name="traces_kernel_time_struct"></a>
The `kernel_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `kernel_time_details` - The kernel elapsed time details.
  The [kernel_time_details](#kernel_time_details_struct) structure is documented below.

<a name="kernel_time_details_struct"></a>
The `kernel_time_details` block supports:

* `parse_time` - The SQL parsing time, in microseconds.

* `rewrite_time` - The SQL rewriting time, in microseconds.

* `plan_time` - The SQL plan generation time, in microseconds.

* `execution_time` - The execution time within the executor, in microseconds.

* `other_time` - The other elapsed time, in microseconds.

<a name="traces_kernel_execution_time_struct"></a>
The `kernel_execution_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `kernel_execution_time_details` - The kernel execution elapsed time details.
  The [kernel_execution_time_details](#kernel_execution_time_details_struct) structure is documented below.

<a name="kernel_execution_time_details_struct"></a>
The `kernel_execution_time_details` block supports:

* `execution_time` - The execution time within the executor, in microseconds.

* `other_time` - The other elapsed time, in microseconds.

<a name="traces_wait_event_time_struct"></a>
The `wait_event_time` block supports:

* `code_wait_event_time` - The code wait event elapsed time.
  The [code_wait_event_time](#code_wait_event_time_struct) structure is documented below.

* `resource_wait_event_time` - The resource wait event elapsed time.
  The [resource_wait_event_time](#resource_wait_event_time_struct) structure is documented below.

<a name="code_wait_event_time_struct"></a>
The `code_wait_event_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `code_wait_event_time_details` - The code wait event time details.
  The [code_wait_event_time_details](#code_wait_event_time_details_struct) structure is documented below.

<a name="code_wait_event_time_details_struct"></a>
The `code_wait_event_time_details` block supports:

* `events` - The list of top five event elapsed time information.
  The [events](#code_wait_event_time_details_top_event_struct) structure is documented below.

* `left_time` - The remaining event elapsed time, in microseconds.

* `other_time` - The elapsed time outside events, in microseconds.

<a name="code_wait_event_time_details_top_event_struct"></a>
The `events` block supports:

* `event_name` - The event name.

* `event_time` - The event elapsed time, in microseconds.

<a name="resource_wait_event_time_struct"></a>
The `resource_wait_event_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `other_time` - The elapsed time outside resource wait events, in microseconds.

* `resource_wait_event_time_details` - The resource wait event elapsed time details.
  The [resource_wait_event_time_details](#resource_wait_event_time_details_struct) structure is documented below.

<a name="resource_wait_event_time_details_struct"></a>
The `resource_wait_event_time_details` block supports:

* `data_io_time` - The IO elapsed time information.
  The [data_io_time](#data_io_time_struct) structure is documented below.

* `lock_time` - The lock elapsed time information.
  The [lock_time](#lock_time_struct) structure is documented below.

* `lwlock_time` - The lightweight lock elapsed time information.
  The [lwlock_time](#lwlock_time_struct) structure is documented below.

<a name="data_io_time_struct"></a>
The `data_io_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `data_io_time_details` - The IO time details.
  The [data_io_time_details](#data_io_time_details_struct) structure is documented below.

<a name="data_io_time_details_struct"></a>
The `data_io_time_details` block supports:

* `events` - The list of top five event elapsed time information.
  The [events](#data_io_time_details_top_event_struct) structure is documented below.

* `left_time` - The remaining event elapsed time, in microseconds.

* `other_time` - The elapsed time outside events, in microseconds.

<a name="data_io_time_details_top_event_struct"></a>
The `events` block supports:

* `event_name` - The event name.

* `event_time` - The event elapsed time, in microseconds.

<a name="lock_time_struct"></a>
The `lock_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `lock_time_details` - The lock time details.
  The [lock_time_details](#lock_time_details_struct) structure is documented below.

<a name="lock_time_details_struct"></a>
The `lock_time_details` block supports:

* `events` - The list of top five event elapsed time information.
  The [events](#lock_time_details_top_event_struct) structure is documented below.

* `left_time` - The remaining event elapsed time, in microseconds.

* `other_time` - The elapsed time outside events, in microseconds.

<a name="lock_time_details_top_event_struct"></a>
The `events` block supports:

* `event_name` - The event name.

* `event_time` - The event elapsed time, in microseconds.

<a name="lwlock_time_struct"></a>
The `lwlock_time` block supports:

* `all_time` - The total elapsed time, in microseconds.

* `lwlock_time_details` - The lightweight lock elapsed time details.
  The [lwlock_time_details](#lwlock_time_details_struct) structure is documented below.

<a name="lwlock_time_details_struct"></a>
The `lwlock_time_details` block supports:

* `events` - The list of top five event elapsed time information.
  The [events](#lwlock_time_details_top_event_struct) structure is documented below.

* `left_time` - The remaining event elapsed time, in microseconds.

* `other_time` - The elapsed time outside events, in microseconds.

<a name="lwlock_time_details_top_event_struct"></a>
The `events` block supports:

* `event_name` - The event name.

* `event_time` - The event elapsed time, in microseconds.

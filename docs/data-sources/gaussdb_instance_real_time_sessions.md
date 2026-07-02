---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_real_time_sessions"
description: |-
  Use this data source to query the real-time sessions of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_real_time_sessions

Use this data source to query the real-time sessions of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id
}
```

### Filter Sessions by Query Info

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}
variable "database_name" {}

data "huaweicloud_gaussdb_instance_real_time_sessions" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id

  query_info {
    database_name = var.database_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the real-time sessions.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String) Specifies the node ID of the GaussDB instance.

* `component_id` - (Required, String) Specifies the component ID of the GaussDB instance.

* `query_info` - (Optional, List) Specifies the query filter conditions.
  The [query_info](#gaussdb_instance_real_time_sessions_query_info_arg) structure is documented below.

<a name="gaussdb_instance_real_time_sessions_query_info_arg"></a>
The `query_info` block supports:

* `database_name` - (Optional, String) Specifies the database name for filtering sessions.

* `client_ip` - (Optional, String) Specifies the client IP address for filtering sessions.

* `user_name` - (Optional, String) Specifies the user name for filtering sessions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sessions` - The list of real-time sessions.
  The [sessions](#gaussdb_instance_real_time_sessions_attr) structure is documented below.

<a name="gaussdb_instance_real_time_sessions_attr"></a>
The `sessions` block supports:

* `session_id` - The session ID.

* `pid` - The process ID.

* `unique_sql_id` - The unique SQL ID.

* `database_name` - The database name.

* `client_ip` - The client IP address.

* `user_name` - The user name.

* `wait` - The wait status.

* `block_session` - The blocked session ID.

* `wait_event` - The wait event.

* `state` - The session state.

* `query_runtime` - The query runtime.

* `query` - The query SQL.

* `back_end_start` - The backend start time.

* `transaction_start` - The transaction start time.

* `query_start` - The query start time.

* `application_name` - The application name.

* `exec_time` - The execution time.

* `trans_num` - The transaction number.

* `rollback_num` - The rollback number.

* `sql_num` - The SQL number.

* `client_port` - The client port.

* `query_id` - The query ID.

* `transaction_time_cost` - The transaction time cost.

* `trace_id` - The trace ID.

* `global_session_id` - The global session ID.

* `top_transaction_id` - The top transaction ID.

* `current_transaction_id` - The current transaction ID.

* `xlog_quantity_pretty` - The xlog quantity.

* `wait_status` - The wait status description.

* `lwt_id` - The lightweight transaction ID.

* `thread_name` - The thread name.

* `lock_mode` - The lock mode.

* `parent_session_id` - The parent session ID.

* `smp_id` - The SMP ID.

* `lock_tag` - The lock tag.

* `component_name` - The component name.

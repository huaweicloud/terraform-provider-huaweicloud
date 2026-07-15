---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_real_time_transactions"
description: |-
  Use this data source to query the real-time transaction list of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_real_time_transactions

Use this data source to query the real-time transaction list of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

data "huaweicloud_gaussdb_real_time_transactions" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id
}
```

### Filter by Database Names

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "component_id" {}

data "huaweicloud_gaussdb_real_time_transactions" "test" {
  instance_id  = var.instance_id
  node_id      = var.node_id
  component_id = var.component_id

  transaction_query_info {
    datnames = ["postgres"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the real-time transactions.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String) Specifies the node ID.
  Only nodes with non-log CN or DN components are supported.

* `component_id` - (Required, String) Specifies the component ID.
  Only non-log CN or DN components are supported.

* `transaction_query_info` - (Optional, List) Specifies the query conditions for transactions.
  The [transaction_query_info](#transaction_query_info_attr) structure is documented below.

<a name="transaction_query_info_attr"></a>
The `transaction_query_info` block supports:

* `exec_time` - (Optional, String) Specifies the filter for transaction execution time (seconds).

* `xlog_quantity` - (Optional, String) Specifies the filter for transaction xlog size (bytes).

* `datnames` - (Optional, List) Specifies the list of database names to filter.

* `usenames` - (Optional, List) Specifies the list of user names to filter.

* `client_addrs` - (Optional, List) Specifies the list of client addresses to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rows` - The list of real-time transactions.
  The [rows](#rows_attr) structure is documented below.

<a name="rows_attr"></a>
The `rows` block supports:

* `sessionid` - The transaction ID.

* `pid` - The thread ID.

* `query_id` - The SQL query ID.

* `datname` - The database name.

* `client_addr` - The client address.

* `client_port` - The client port.

* `usename` - The user name.

* `query` - The SQL query statement.

* `backend_start` - The session start time.

* `xact_start` - The transaction start time.

* `application_name` - The application name.

* `state` - The transaction state.

* `state_change` - The transaction state change time.

* `query_start` - The query start time.

* `waiting` - The waiting lock status.

* `unique_sql_id` - The normalized SQL ID.

* `top_xid` - The top-level transaction ID.

* `current_xid` - The current transaction ID.

* `exec_time` - The transaction execution time (seconds).

* `xlog_quantity` - The xlog quantity (bytes).

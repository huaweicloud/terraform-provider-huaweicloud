---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_global_slow_sql_detail"
description: |-
  Use this data source to query the global slow SQL detail of GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_global_slow_sql_detail

Use this data source to query the global slow SQL detail of GaussDB instances within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "sql_id" {}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_global_slow_sql_detail" "test" {
  instance_id    = var.instance_id
  start_time     = "1704067200000"
  end_time       = "1735689600000"
  sql_id         = var.sql_id
  node_ids       = var.node_ids
  component_type = "cn"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the global slow SQL detail.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time in 13-digit UNIX timestamp format (milliseconds), UTC
  timezone.

* `end_time` - (Required, String) Specifies the end time in 13-digit UNIX timestamp format (milliseconds), UTC timezone.

* `sql_id` - (Required, String) Specifies the SQL ID.

* `node_ids` - (Required, List) Specifies the list of node IDs. Cannot be empty.

* `component_type` - (Required, String) Specifies the component type. The valid values are:
  + **cn**: CN component type.
  + **dn**: DN component type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_sql_details` - The list of slow SQL details.
  The [slow_sql_details](#gaussdb_global_slow_sql_detail_slow_sql_details) structure is documented below.

<a name="gaussdb_global_slow_sql_detail_slow_sql_details"></a>
The `slow_sql_details` block supports:

* `db_name` - The database name.

* `schema_name` - The schema name.

* `sql_id` - The SQL ID.

* `user_name` - The user name.

* `client_ip` - The client IP address.

* `client_port` - The client port.

* `node_id` - The node ID.

* `node_name` - The node name.

* `sql_text` - The SQL template.

* `sql` - The complete SQL with variables replaced. When sql_text does not return variable values, sql returns an empty
  string.

* `query_plan` - The query execution plan.

* `start_time` - The start time in the format of yyyy-mm-ddThh:mm:ss+0000.

* `finish_time` - The finish time in the format of yyyy-mm-ddThh:mm:ss+0000.

* `returned_rows` - The number of returned rows.

* `fetched_rows` - The number of fetched rows.

* `fetched_pages` - The number of fetched pages.

* `hit_pages` - The number of hit pages.

* `total_time` - The total time.

* `cpu_time` - The CPU time.

* `plan_time` - The plan time.

* `io_time` - The IO time.

* `lock_count` - The lock count.

* `lock_time` - The lock time.

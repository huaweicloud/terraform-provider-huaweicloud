---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_top_sql_statements"
description: |-
  Use this data source to query the Top SQL list of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_top_sql_statements

Use this data source to query the Top SQL list of a specified GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_top_sql_statements" "test" {
  instance_id = var.instance_id
  node_ids    = var.node_ids
  start_time  = 1750108800000
  end_time    = 1750195200000
}
```

### Query with Filters

```hcl
variable "instance_id" {}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_top_sql_statements" "test" {
  instance_id    = var.instance_id
  node_ids       = var.node_ids
  start_time     = 1750108800000
  end_time       = 1750195200000
  db_name        = "mydb"
  support_system = true
}
```

### Query with Multi Queries

```hcl
variable "instance_id" {}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_top_sql_statements" "test" {
  instance_id = var.instance_id
  node_ids    = var.node_ids
  start_time  = 1750108800000
  end_time    = 1750195200000

  multi_queries {
    name      = "query"
    condition = "and"
    values = ["SELECT", "INSERT"]
    is_fuzzy  = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the Top SQL list. If omitted, the provider-level
  region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_ids` - (Required, List of String) Specifies the list of node IDs to query.

* `start_time` - (Required, Int) Specifies the start time, in 13-digit UNIX timestamp format (milliseconds, UTC).

* `end_time` - (Required, Int) Specifies the end time, in 13-digit UNIX timestamp format (milliseconds, UTC).

* `start_time_utc` - (Optional, String) Specifies the start time in UTC format (e.g. `yyyy-mm-ddThh:mm:ssZ`).

* `end_time_utc` - (Optional, String) Specifies the end time in UTC format (e.g. `yyyy-mm-ddThh:mm:ssZ`).

* `support_system` - (Optional, Boolean) Specifies whether to display system users. Default value is `false`.

* `sql_id` - (Optional, String) Specifies the normalized SQL ID of the Top SQL.

* `db_name` - (Optional, String) Specifies the database name to filter. Only supported for engine version 8.200 and
  above.

* `user_name` - (Optional, String) Specifies the username to filter.

* `sql_text` - (Optional, String) Specifies the SQL text to filter.

* `multi_queries` - (Optional, List) Specifies the list of field aggregation query conditions.
  The [multi_queries](#gaussdb_top_sql_statements_multi_queries) structure is documented below.

<a name="gaussdb_top_sql_statements_multi_queries"></a>
The `multi_queries` block supports:

* `name` - (Required, String) Specifies the query field name. Only `"query"` is supported.

* `condition` - (Required, String) Specifies the merge condition between multiple filter conditions. The valid values
  are `and`, `or`, `AND`, `OR`.

* `values` - (Required, List of String) Specifies the list of filter query values. Contains 1 to 5 strings.

* `is_fuzzy` - (Optional, String) Specifies whether to perform fuzzy query. The valid values are `"true"` and `"false"`.
  Default value is `"true"`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `top_sql_infos` - The list of Top SQL information.
  The [top_sql_infos](#gaussdb_top_sql_infos) structure is documented below.

<a name="gaussdb_top_sql_infos"></a>
The `top_sql_infos` block supports:

* `sql_id` - The normalized SQL ID.

* `user_name` - The username.

* `sql_text` - The SQL text.

* `calls_percent` - The call frequency percentage (0-100).

* `cpu_percent` - The CPU cost percentage (0-100).

* `io_percent` - The IO cost percentage (0-100).

* `calls` - The number of calls.

* `returned_rows` - The number of returned tuples.

* `tuple_read` - The number of read tuples.

* `avg_elapse_time` - The average time cost, in ms.

* `total_elapse_time` - The total time cost, in ms.

* `cpu_time` - The CPU cost, in ms.

* `io_time` - The IO cost, in ms.

* `min_elapse_time` - The minimum execution time, in ms.

* `max_elapse_time` - The maximum execution time, in ms.

* `sql_hit_ratio` - The SQL hit ratio.

* `node_id` - The node ID.

* `node_name` - The node name.

* `db_name` - The database name.

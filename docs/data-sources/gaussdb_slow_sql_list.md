---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_slow_sql_list"
description: |-
  Use this data source to query the slow SQL list of GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_slow_sql_list

Use this data source to query the slow SQL list of GaussDB instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_slow_sql_list" "test" {
  instance_id = var.instance_id
  start_time  = "1704067200000"
  end_time    = "1735689600000"
  threshold   = 5
  node_ids    = var.node_ids

  multi_queries {
    name      = "query"
    condition = "AND"
    is_fuzzy  = true
    values    = ["SELECT"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the slow SQL list.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `start_time` - (Required, String) Specifies the start time in 13-digit UNIX timestamp format (milliseconds), UTC
  timezone.

* `end_time` - (Required, String) Specifies the end time in 13-digit UNIX timestamp format (milliseconds), UTC timezone.

* `threshold` - (Required, Integer) Specifies the slow SQL threshold in seconds. The value must be between **1** and *
  *99**.

* `node_ids` - (Required, List) Specifies the list of node IDs. Cannot be empty.

* `sql_id` - (Optional, String) Specifies the SQL ID for filtering.

* `multi_queries` - (Optional, List) Specifies the list of multi-query conditions.
  The [multi_queries](#gaussdb_slow_sql_list_multi_queries_structure) structure is documented below.

<a name="gaussdb_slow_sql_list_multi_queries_structure"></a>
The `multi_queries` block supports:

* `name` - (Required, String) Specifies the query field name. The valid value is **query**.

* `condition` - (Required, String) Specifies the merge condition. The valid values are **and**, **or**, **AND**, **OR**.

* `values` - (Required, List) Specifies the list of filter values. A list of 1 to 5 strings.

* `is_fuzzy` - (Optional, String) Specifies whether to use fuzzy query. The valid values are:
  + **true**: Fuzzy query. (Default)
  + **false**: Exact match.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `slow_sql_infos` - The list of slow SQL information.
  The [slow_sql_infos](#gaussdb_slow_sql_list_slow_sql_infos) structure is documented below.

<a name="gaussdb_slow_sql_list_slow_sql_infos"></a>
The `slow_sql_infos` block supports:

* `db_name` - The database name.

* `schema_name` - The SCHEMA name.

* `sql_id` - The SQL ID.

* `user_name` - The user name.

* `node_id` - The node ID.

* `node_name` - The node name.

* `sql_text` - The SQL template.

* `sql` - The complete SQL with variables replaced. When sql_text does not return variable values, sql returns an empty
  string.

* `query_plan` - The query execution plan.

* `calls` - The number of executions.

* `avg_exec_time` - The average execution time (in seconds).

* `avg_cpu_time` - The CPU time (in seconds).

* `avg_io_time` - The IO time (in seconds).

* `avg_returned_rows` - The number of returned rows.

* `avg_fetched_rows` - The number of scanned rows.

* `buffer_hit_ratio` - The buffer hit ratio.

* `sql_hit_ratio` - The SQL hit ratio.

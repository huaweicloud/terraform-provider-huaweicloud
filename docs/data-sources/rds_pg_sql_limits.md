---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_sql_limits"
description: |-
  Use this data source to get the list of RDS PostgreSQL SQL limits.
---

# huaweicloud_rds_pg_sql_limits

Use this data source to get the list of RDS PostgreSQL SQL limits.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}

data "huaweicloud_rds_pg_sql_limits" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of RDS PostgreSQL instance.

* `db_name` - (Required, String) Specifies the name of the database.

* `sql_limit_id` - (Optional, String) Specifies the ID of SQL limit.

* `query_id` - (Optional, String) Specifies the query ID.

* `query_string` - (Optional, String) Specifies the text form of SQL statement.

* `max_concurrency` - (Optional, String) Specifies the number of SQL statements executed simultaneously.

* `max_waiting` - (Optional, String) Specifies the max waiting time in seconds.

* `search_path` - (Optional, String) Specifies the query order for names that are not schema qualified.

* `is_effective` - (Optional, String) Specifies whether the SQL limit is effective.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sql_limits` - Indicates the list of SQL limits.

  The [sql_limits](#sql_limits_struct) structure is documented below.

<a name="sql_limits_struct"></a>
The `sql_limits` block supports:

* `id` - Indicates the ID of SQL limit.

* `query_id` - Indicates the query ID.

* `query_string` - Indicates the text form of SQL statement.

* `max_concurrency` - Indicates the number of SQL statements executed simultaneously.

* `max_waiting` - Indicates the max waiting time in seconds.

* `search_path` - Indicates the query order for names that are not schema qualified.

* `is_effective` - Indicates whether the SQL limit is effective.

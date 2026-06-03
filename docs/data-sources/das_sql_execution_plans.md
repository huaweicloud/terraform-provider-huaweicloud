---
subcategory: "Data Admin Service (DAS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_das_sql_execution_plans"
description: |-
  Use this data source to get the list of DAS SQL execution plans.
---

# huaweicloud_das_sql_execution_plans

Use this data source to get the list of DAS SQL execution plans.

-> This data source only supports to query SQL execution plans of **MySQL** instances.

## Example Usage

### Basic usage

```hcl
variable "instance_id" {}
variable "db_user_id" {}
variable "database" {}
variable "sql" {}

data "huaweicloud_das_sql_execution_plans" "test" {
  instance_id = var.instance_id
  db_user_id  = var.db_user_id
  database    = var.database
  sql         = var.sql
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the SQL execution plans are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the database instance.

* `db_user_id` - (Required, String) Specifies the database user ID.

* `database` - (Required, String) Specifies the database name.

* `sql` - (Required, String) Specifies the SQL statement.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plans` - The list of SQL execution plans.  
  The [plans](#sql_execution_plans_plans) structure is documented below.

* `error_message` - The error message if the SQL execution failed.

<a name="sql_execution_plans_plans"></a>
The plans block supports:

* `id` - The ID of the execution plan step.

* `select_type` - The select type of the query.

* `table` - The table name.

* `partitions` - The partitions that the query will match.

* `type` - The access type.

* `possible_keys` - The possible keys that could be used.

* `key` - The actual key used.

* `key_len` - The length of the key used.

* `ref` - The column or constant used with the key to select rows.

* `rows` - The number of rows MySQL estimates it must examine.

* `filtered` - The percentage of rows filtered by the table condition.

* `extra` - Additional information.

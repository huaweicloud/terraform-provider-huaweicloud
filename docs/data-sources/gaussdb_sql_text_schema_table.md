---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_text_schema_table"
description: |-
  Use this data source to identify table information in SQL text of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_text_schema_table

Use this data source to identify table information in SQL text of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_sql_text_schema_table" "test" {
  instance_id = var.instance_id
  sql_text    = "SELECT * FROM public.users u LEFT JOIN sales.orders o ON u.id = o.user_id;"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to identify the SQL text table information.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `sql_text` - (Required, String) Specifies the SQL text to parse for table information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_tables` - The list of table information identified from the SQL text.
  The [database_tables](#database_tables_struct) structure is documented below.

<a name="database_tables_struct"></a>
The `database_tables` block supports:

* `table_name` - The name of the table.

* `schema_name` - The name of the schema.

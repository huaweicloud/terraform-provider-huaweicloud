---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_database_tables"
description: |-
  Use this data source to get the list of database tables in a GaussDB instance.
---

# huaweicloud_gaussdb_instance_database_tables

Use this data source to get the list of database tables in a GaussDB instance.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "schema_name" {}

data "huaweicloud_gaussdb_instance_database_tables" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name
  schema_name = var.schema_name
}
```

### Filter by Table Name Keyword

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "schema_name" {}
variable "table_name_keyword" {}

data "huaweicloud_gaussdb_instance_database_tables" "filter" {
  instance_id        = var.instance_id
  db_name            = var.db_name
  schema_name        = var.schema_name
  table_name_keyword = var.table_name_keyword
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name.
  The value cannot be a template database name.
  The template databases include **postgres**, **template0**, **template1**, **templatea**,
  **template_pdb**, and **templatem**.

* `schema_name` - (Required, String, NonUpdatable) Specifies the schema name.
  The value cannot be **public** or **information_schema**.

* `table_name_keyword` - (Optional, String, NonUpdatable) Specifies the keyword of the table name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - The list of database tables.
  The [tables](#gaussdb_instance_database_tables_tables) structure is documented below.

<a name="gaussdb_instance_database_tables_tables"></a>
The `tables` block supports:

* `table_name` - The table name.

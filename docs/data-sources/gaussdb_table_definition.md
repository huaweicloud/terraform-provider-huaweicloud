---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_table_definition"
description: |-
  Use this data source to query the table definition information of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_table_definition

Use this data source to query the table definition information of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_table_definition" "test" {
  instance_id   = var.instance_id
  database_name = "postgres"
  table_name    = "t1"
}
```

### Query with schema name

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_table_definition" "test" {
  instance_id   = var.instance_id
  database_name = "postgres"
  table_name    = "t1"
  schema_name   = "public"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the table definition.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `database_name` - (Required, String) Specifies the name of the database.

* `table_name` - (Required, String) Specifies the name of the table.

* `schema_name` - (Optional, String) Specifies the name of the schema.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_definitions` - The list of table definition information.
  The [table_definitions](#gaussdb_table_definition_table_definitions) structure is documented below.

<a name="gaussdb_table_definition_table_definitions"></a>
The `table_definitions` block supports:

* `table_definition` - The DDL definition of the table.

* `schema_name` - The name of the schema.

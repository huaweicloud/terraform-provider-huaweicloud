---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_schemas_storage_usage"
description: |-
  Use this data source to query the schema storage usage of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_schemas_storage_usage

Use this data source to query the schema storage usage of a GaussDB instance within HuaweiCloud.

## Example Usage

### Query all schemas in a database

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_schemas_storage_usage" "test" {
  instance_id   = var.instance_id
  database_name = "test_db"
}
```

### Query a specific schema

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_schemas_storage_usage" "test" {
  instance_id   = var.instance_id
  database_name = "test_db"
  schema_name   = "public"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the schema storage usage.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `database_name` - (Required, String) Specifies the name of the database.

* `schema_name` - (Optional, String) Specifies the name of the schema to filter.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schema_volumes` - The list of schema storage usage information.
  The [schema_volumes](#schema_volumes_struct) structure is documented below.

<a name="schema_volumes_struct"></a>
The `schema_volumes` block supports:

* `schema_size` - The size of the schema (e.g., "28 MB").

* `table_count` - The number of tables in the schema.

* `user_name` - The name of the user that owns the schema.

* `schema_name` - The name of the schema.

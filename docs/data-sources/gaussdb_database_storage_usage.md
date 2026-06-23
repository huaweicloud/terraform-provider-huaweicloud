---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_database_storage_usage"
description: |-
  Use this data source to query the database storage usage of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_database_storage_usage

Use this data source to query the database storage usage of a GaussDB instance within HuaweiCloud.

## Example Usage

### Query all databases

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_database_storage_usage" "test" {
  instance_id = var.instance_id
}
```

### Query with filters

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_database_storage_usage" "test" {
  instance_id      = var.instance_id
  database_name    = "postgres"
  table_space_name = "pg_default"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the database storage usage.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `database_name` - (Optional, String) Specifies the name of the database to filter.

* `table_space_name` - (Optional, String) Specifies the name of the default table space to filter.

* `user_name` - (Optional, String) Specifies the name of the user that owns the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `database_volumes` - The list of database storage usage information.
  The [database_volumes](#database_volumes_struct) structure is documented below.

<a name="database_volumes_struct"></a>
The `database_volumes` block supports:

* `database_name` - The name of the database.

* `table_space_name` - The default table space name of the database.

* `user_name` - The name of the user that owns the database.

* `database_size` - The storage usage of the database (e.g., "708 MB").

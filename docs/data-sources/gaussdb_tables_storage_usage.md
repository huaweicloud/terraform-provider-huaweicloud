---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_tables_storage_usage"
description: |-
  Use this data source to query the storage usage of tables in a specified schema of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_tables_storage_usage

Use this data source to query the storage usage of tables in a specified schema of a GaussDB instance within HuaweiCloud.

## Example Usage

### Query tables in specific schemas

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_tables_storage_usage" "test" {
  instance_id   = var.instance_id
  database_name = "test_db"
  schema_names  = ["public", "pg_catalog"]
}
```

### Query with table name filter

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_tables_storage_usage" "test" {
  instance_id   = var.instance_id
  database_name = "test_db"
  schema_names  = ["public"]
  table_name    = "my_table"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the tables storage usage.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `database_name` - (Required, String) Specifies the name of the database.

* `schema_names` - (Required, List) Specifies the list of schema names.

* `table_name` - (Optional, String) Specifies the name of the table to filter.

* `user_name` - (Optional, String) Specifies the name of the table owner to filter.

* `sort_key` - (Optional, String) Specifies the sort field.
  Valid values: **table_size**, **id**, **table_name**, **table_owner**, **database_name**, **schema_name**, **is_part_type**.

* `sort_order` - (Optional, String) Specifies the sort method.
  Valid values: **DESC**, **desc**, **ASC**, **asc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `table_volumes` - The list of table storage usage information.
  The [table_volumes](#table_volumes_struct) structure is documented below.

<a name="table_volumes_struct"></a>
The `table_volumes` block supports:

* `table_size` - The size of the table (e.g., "24 kB").

* `id` - The ID of the table.

* `table_name` - The name of the table.

* `table_owner` - The owner of the table.

* `schema_name` - The name of the schema.

* `database_name` - The name of the database.

* `is_part_type` - Whether the table or index has partition table properties.

* `is_hash_cluster_key` - Whether it contains hash partition column information.

* `tuples` - The number of tuples in the table.

* `create_time` - The creation time of the table.

* `update_time` - The update time of the table.

* `average_size` - The average size of the table.

* `max_ratio` - The maximum ratio.

* `min_ratio` - The minimum ratio.

* `skew_size` - The skew size of the table.

* `skew_ratio` - The skew ratio of the table.

* `skew_stddev` - The skew standard deviation of the table.

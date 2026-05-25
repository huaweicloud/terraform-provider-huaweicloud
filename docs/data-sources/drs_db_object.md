---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_db_object"
description: |-
  Use this data source to get the database object information of specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_db_object

Use this data source to get the database object information of specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_db_object" "test" {
  job_id = var.job_id
  type   = "modified"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the job ID.

* `type` - (Required, String) Specifies the query object information type.  
  The valid values are as follows:
  + **modified**: Query selected (synchronized and not delivered) object information.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `target_root_db` - The target database information for database object migration or synchronization.

  The [target_root_db](#target_root_db_struct) structure is documented below.

* `object_info` - The database object migration or synchronization information list.

  The [object_info](#object_info_struct) structure is documented below.

* `max_table_num` - The threshold of the number of tables under a database.

* `status` - The status of obtaining submitted query object selection information.

* `object_scope` - The database type in real-time synchronization scenarios.

<a name="target_root_db_struct"></a>
The `target_root_db` block supports:

* `db_name` - The database name.

* `db_encoding` - The default encoding format is utf8.

<a name="object_info_struct"></a>
The `object_info` block supports:

* `key` - The key of the object info map.

* `sync_type` - The database type in real-time synchronization scenarios.

* `name` - The database name in the target database (database name mapping).

* `all` - Whether to migrate or synchronize the entire database.

* `schemas` - The schemas to be migrated or synchronized.

  The [schemas](#schemas_struct) structure is documented below.

* `tables` - The tables to be migrated or synchronized.

  The [tables](#tables_struct) structure is documented below.

* `total_table_num` - The number of tables under the database.

* `is_synchronized` - Whether it has been synchronized.

<a name="schemas_struct"></a>
The `schemas` block supports:

* `key` - The key of the schemas map.

* `sync_type` - The schema type in real-time synchronization scenarios.

* `name` - The schema name in the target database (schema name mapping).

* `all` - Whether to migrate or synchronize the entire schema.

* `tables` - The tables to be migrated or synchronized.

  The [tables](#tables_struct) structure is documented below.

<a name="tables_struct"></a>
The `tables` block supports:

* `key` - The key of the tables map.

* `sync_type` - The table type in real-time synchronization scenarios.

* `type` - The object type.

* `name` - The table name in the target database (table name mapping).

* `all` - Whether to migrate or synchronize the entire table.

* `db_alias_name` - The database name mapping at the table level in one-to-many scenarios.

* `schema_alias_name` - The schema name mapping at the table level in one-to-many scenarios.

* `filtered` - Whether data filtering is performed on the table.

* `filter_conditions` - The filter conditions for the table data.

* `config_conditions` - The configuration conditions for advanced data filtering settings of the table.

* `is_synchronized` - Whether it has been synchronized.

* `columns` - The columns to be synchronized, mapped, filtered, or added.

  The [columns](#columns_struct) structure is documented below.

<a name="columns_struct"></a>
The `columns` block supports:

* `key` - The key of the columns map.

* `sync_type` - The column type in real-time synchronization scenarios.

* `primary_key_for_data_filtering` - The primary key column name in advanced data filtering scenarios.

* `index_for_data_filtering` - The index required for query optimization.

* `name` - The column name in the target database (column name mapping).

* `type` - The data type of the column field.

* `primary_key_for_column_filtering` - Whether the column is the primary key in column mapping scenarios.

* `filtered` - Whether the column is filtered.

* `additional` - Whether the column is an additional column.

* `operation_type` - The operation type.

* `value` - The value of the additional column.

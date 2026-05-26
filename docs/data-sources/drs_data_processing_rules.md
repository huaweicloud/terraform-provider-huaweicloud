---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_data_processing_rules"
description: |-
  Use this data source to get the data processing rules for specified DRS job within HuaweiCloud.
---

# huaweicloud_drs_data_processing_rules

Use this data source to get the data processing rules for specified DRS job within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_data_processing_rules" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the DRS job ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_process_info` - The data processing rules information.

  The [data_process_info](#data_process_info_struct) structure is documented below.

<a name="data_process_info_struct"></a>
The `data_process_info` block supports:

* `filter_conditions` - The filter conditions.

  The [filter_conditions](#filter_conditions_struct) structure is documented below.

* `is_batch_process` - Whether it is batch processing.
  + **true**: Database-level or batch table-level processing.
  + **false**: Single table operation.

* `add_columns` - The additional columns.

  The [add_columns](#add_columns_struct) structure is documented below.

* `ddl_operation` - The DDL operations supported.
  The key-value pairs indicate the DDL operation types and their meanings.
  For example: "table": "CREATE TABLE, ALTER TABLE, DROP TABLE, RENAME TABLE".

* `dml_operation` - The DML operations supported.
  The valid values are as follows:
  + **i**: INSERT
  + **u**: UPDATE
  + **d**: DELETE

* `db_object_column_info` - The database object column information.

  The [db_object_column_info](#db_object_column_info_struct) structure is documented below.

* `db_or_table_rename_rule` - The database or table rename rule.

  The [db_or_table_rename_rule](#db_or_table_rename_rule_struct) structure is documented below.

* `db_object` - The database object information.

  The [db_object](#db_object_struct) structure is documented below.

* `is_synchronized` - Whether the rule has been synchronized to the target database.

* `source` - The comparison source.
  The valid values are as follows:
  + **job**: Data synchronization filtering.
  + **compare**: Data comparison filtering.

* `process_rule_level` - The data processing rule level.
  The valid values are as follows:
  + **table**: Data synchronization filtering.
  + **combinations**: Combination set, operations on multiple tables.

<a name="filter_conditions_struct"></a>
The `filter_conditions` block supports:

* `value` - The filter condition.
  When filtering_type is configConditionalFilter, the value defaults to config.
  When filtering_type is contentConditionalFilter, the value defaults to the filter condition.

* `filtering_type` - The filter condition type.
  The valid values are as follows:
  + **contentConditionalFilter**: Simple condition filtering.
  + **configConditionalFilter**: Associated table filtering.

<a name="add_columns_struct"></a>
The `add_columns` block supports:

* `column_type` - The column type.
  The valid values are as follows:
  + **default_value**: Default value.
  + **create_time**: Create time.
  + **update_time**: Update time.
  + **expression**: Expression.
  + **server_database_table**: Server database table.

* `column_name` - The column name.

* `column_value` - The column fill value.

* `data_type` - The data type of the filled column.
  The valid values are as follows:
  + **int**
  + **long**
  + **varchar(256)**
  + **varchar(191)**
  + **datetime**
  + **timestamp**

<a name="db_object_column_info_struct"></a>
The `db_object_column_info` block supports:

* `db_name` - The database name.

* `schema_name` - The database schema name.

* `table_name` - The database table name.

* `column_infos` - The database column information.

  The [column_infos](#column_infos_struct) structure is documented below.

* `total_count` - The total number of database column information.

<a name="column_infos_struct"></a>
The `column_infos` block supports:

* `column_name` - The column name.

* `column_type` - The column type.

* `primary_key_or_unique_index` - The primary key or unique index.

* `column_mapped_name` - The column name after mapping.

* `is_filtered` - Whether the column is filtered.

* `is_partition_key` - Whether the column is a partition key.

<a name="db_or_table_rename_rule_struct"></a>
The `db_or_table_rename_rule` block supports:

* `prefix_name` - The prefix name.
  When type is prefixAndSuffix, fill in prefix_name, the database/table name only adds prefix.
  If suffix_name is also filled, the database/table name adds both prefix and suffix.

* `suffix_name` - The suffix name.
  When type is prefixAndSuffix, fill in suffix_name, the database/table name only adds suffix.
  If prefix_name is also filled, the database/table name adds both prefix and suffix.

* `type` - The database/table mapping type.
  The valid values are as follows:
  + **prefixAndSuffix**: Prefix, suffix, or both prefix and suffix.
  + **manyToOne**: Many-to-one.

<a name="db_object_struct"></a>
The `db_object` block supports:

* `object_scope` - The database object migration or synchronization scope.
  The valid values are as follows:
  + **all**: All migration.
  + **database**: Database-level migration or synchronization.
  + **table**: Table-level migration or synchronization.

* `target_root_db` - The target database for database object migration or synchronization.
  Required for two-to-three layer database synchronization.

  The [target_root_db](#target_root_db_struct) structure is documented below.

* `object_info` - The database object migration or synchronization information in JSON format.
  This field contains detailed information about the objects to be migrated or synchronized,
  including database, schema, table, and column configurations.

<a name="target_root_db_struct"></a>
The `target_root_db` block supports:

* `db_name` - The database name.

* `db_encoding` - The default encoding format is utf8.

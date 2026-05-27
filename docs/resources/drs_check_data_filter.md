---
subcategory: "Data Replication Service (DRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_check_data_filter"
description: |-
  Manages a resource to check DRS data filtering rules within HuaweiCloud.
---

# huaweicloud_drs_check_data_filter

Manages a resource to check DRS data filtering rules within HuaweiCloud.

-> This resource is a one-time action resource used to check DRS data filtering rules. Deleting this resource will not
  delete the check result from the cloud, but will only remove the resource information from the tf state file.

## Example Usage

```hcl
resource "huaweicloud_drs_check_data_filter" "test" {
  job_id = "b10111ca-542a-4d17-8ecb-239ec86jb201"

  data_process_info {
    filter_conditions {
      filtering_type = "contentConditionalFilter"
      value          = "id>1"
    }

    db_object {
      object_scope = "table"

      object_info = <<EOT
      {
        "dyh4" : {
          "name" : "dyh4",
          "all" : false,
          "tables" : {
            "test1_table1" : {
              "name" : "test1_table1",
              "type" : "table",
              "all" : true
            },
            "test1_table10" : {
              "name" : "test1_table10",
              "type" : "table",
              "all" : true
            },
            "test1_table11" : {
              "name" : "test1_table11",
              "type" : "table",
              "all" : true
            }
          }
        }
      }
      EOT
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the job ID.

* `data_process_info` - (Required, List, NonUpdatable) Specifies the data processing information for filtering rules.
  The [data_process_info](#data_process_info_struct) structure is documented below.

<a name="data_process_info_struct"></a>
The `data_process_info` block supports:

* `filter_conditions` - (Optional, List, NonUpdatable) Specifies the filter conditions.
  Required when performing data filtering rule validation.
  The [filter_conditions](#filter_conditions_struct) structure is documented below.

* `is_batch_process` - (Optional, Bool, NonUpdatable) Specifies whether it is a batch process.
  **true**: Database-level or batch table-level processing.
  **false**: Single table operation.

* `add_columns` - (Optional, List, NonUpdatable) Specifies the additional columns.
  Required when selecting additional columns. Used to avoid data conflicts in many-to-one operations.
  The [add_columns](#add_columns_struct) structure is documented below.

* `ddl_operation` - (Optional, Map, NonUpdatable) Specifies the DDL operation support for incremental migration or
  synchronization. If left empty, no DDL operations will be migrated or synchronized.
  Example: **table**: "CREATE TABLE, ALTER TABLE, DROP TABLE, RENAME TABLE".

* `dml_operation` - (Optional, String, NonUpdatable) Specifies the DML operation support.
  If left empty, no DML operations will be incrementally migrated or synchronized.
  The valid values are as follows:
  + **i**: INSERT.
  + **u**: UPDATE.
  + **d**: DELETE.

* `db_object_column_info` - (Optional, List, NonUpdatable) Specifies the database object column information for column
  mapping and filtering. Required when performing column mapping or filtering.
  The [db_object_column_info](#db_object_column_info_struct) structure is documented below.

* `db_or_table_rename_rule` - (Optional, List, NonUpdatable) Specifies the database or table rename rule.
  The [db_or_table_rename_rule](#db_or_table_rename_rule_struct) structure is documented below.

* `db_object` - (Optional, List, NonUpdatable) Specifies the database object information.
  Required when performing mapping or data filtering condition validation.
  The [db_object](#db_object_struct) structure is documented below.

* `is_synchronized` - (Optional, Bool, NonUpdatable) Specifies whether the rule has been synchronized to the target
  database.

* `source` - (Optional, String, NonUpdatable) Specifies the comparison source.
  The valid values are as follows:
  + **job**: Data synchronization filtering.
  + **compare**: Data comparison filtering.

* `process_rule_level` - (Optional, String, NonUpdatable) Specifies the data processing rule level.
  Required when performing data filtering rule validation or updating data processing rules.
  The valid values are as follows:
  + **table**: Data synchronization filtering.
  + **combinations**: Combination set, operations on multiple tables.

<a name="filter_conditions_struct"></a>
The `filter_conditions` block supports:

* `value` - (Optional, String, NonUpdatable) Specifies the filter condition value.
  When `filtering_type` is **configConditionalFilter**, the value defaults to **config**.
  When `filtering_type` is **contentConditionalFilter**, the value should be the filter condition.
  Only one validation rule can be added per table.
  Data filtering supports up to 500 tables at a time.
  The filter expression must use standard SQL syntax and cannot use database-specific packages, functions, variables,
  or constants. Enter only the part after WHERE in the SQL
  statement (excluding WHERE and semicolons, e.g., `sid > 3 and sname like "G %"`).
  Maximum length: `512` characters.

* `filtering_type` - (Optional, String, NonUpdatable) Specifies the filter condition type.
  The valid values are as follows:
  + **contentConditionalFilter**: Simple condition filtering.
  + **configConditionalFilter**: Related table filtering.

<a name="add_columns_struct"></a>
The `add_columns` block supports:

* `column_type` - (Optional, String, NonUpdatable) Specifies the column type.
  The valid values are as follows:
  + **default_value**: Default value.
  + **create_time**: Create time.
  + **update_time**: Update time.
  + **expression**: Expression.
  + **server_database_table**: Server name, database, and table.

* `column_name` - (Optional, String, NonUpdatable) Specifies the column name.

* `column_value` - (Optional, String, NonUpdatable) Specifies the column filling value.

* `data_type` - (Optional, String, NonUpdatable) Specifies the data type of the filled column.
  The valid values are as follows:
  + **int**: Integer.
  + **long**: Long integer.
  + **varchar(256)**: Variable character string with maximum length `256`.
  + **varchar(191)**: Variable character string with maximum length `191`.
  + **datetime**: Date and time.
  + **timestamp**: Timestamp.

<a name="db_object_column_info_struct"></a>
The `db_object_column_info` block supports:

* `db_name` - (Optional, String, NonUpdatable) Specifies the database name.

* `schema_name` - (Optional, String, NonUpdatable) Specifies the database schema name.

* `table_name` - (Optional, String, NonUpdatable) Specifies the database table name.

* `column_infos` - (Optional, List, NonUpdatable) Specifies the column information.
  The [column_infos](#column_infos_struct) structure is documented below.

* `total_count` - (Optional, Int, NonUpdatable) Specifies the total count of column information.
  This is only a return parameter and is unrelated to pagination.

<a name="column_infos_struct"></a>
The `column_infos` block supports:

* `column_name` - (Optional, String, NonUpdatable) Specifies the column name.

* `column_type` - (Optional, String, NonUpdatable) Specifies the column type.

* `primary_key_or_unique_index` - (Optional, String, NonUpdatable) Specifies the primary key or unique index.

* `column_mapped_name` - (Optional, String, NonUpdatable) Specifies the column name after mapping.

* `is_filtered` - (Optional, Bool, NonUpdatable) Specifies whether the column is filtered.

* `is_partition_key` - (Optional, Bool, NonUpdatable) Specifies whether the column is a partition key.

<a name="db_or_table_rename_rule_struct"></a>
The `db_or_table_rename_rule` block supports:

* `prefix_name` - (Optional, String, NonUpdatable) Specifies the prefix name.
  When `type` is **prefixAndSuffix**, fill in `prefix_name` to add a prefix to the database or table name.
  If `suffix_name` is also filled in, both prefix and suffix will be added.

* `suffix_name` - (Optional, String, NonUpdatable) Specifies the suffix name.
  When `type` is **prefixAndSuffix**, fill in `suffix_name` to add a suffix to the database or table name.
  If `prefix_name` is also filled in, both prefix and suffix will be added.

* `type` - (Optional, String, NonUpdatable) Specifies the rename rule type.
  The valid values are as follows:
  + **prefixAndSuffix**: Prefix, suffix, or both prefix and suffix.
  + **manyToOne**: Many-to-one mapping.

<a name="db_object_struct"></a>
The `db_object` block supports:

* `object_scope` - (Required, String, NonUpdatable) Specifies the database object migration or synchronization scope.
  The valid values are as follows:
  + **all**: Full migration.
  + **database**: Database-level migration or synchronization.
  + **table**: Table-level migration or synchronization.

* `target_root_db` - (Optional, List, NonUpdatable) Specifies the target database for database object migration or
  synchronization. Required for two-to-three layer database synchronization.
  The [target_root_db](#target_root_db_struct) structure is documented below.

* `object_info` - (Optional, String, NonUpdatable) Specifies the database object migration or synchronization
  information in JSON format. Required when `object_scope` is **database** or **table**.
  Do not fill when `object_scope` is **all**.

<a name="target_root_db_struct"></a>
The `target_root_db` block supports:

* `db_name` - (Optional, String, NonUpdatable) Specifies the database name.

* `db_encoding` - (Optional, String, NonUpdatable) Specifies the database encoding.
  The default encoding format is **utf8**.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The check status.
  The valid values are as follows:
  + **pending**: Processing.
  + **failed**: Failed.
  + **success**: Successful.

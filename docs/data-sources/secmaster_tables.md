---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_tables"
description: |-
  Use this data source to get the list of SecMaster tables within HuaweiCloud.
---

# huaweicloud_secmaster_tables

Use this data source to get the list of SecMaster tables within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_tables" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `category` - (Optional, String) Specifies the directory type.
  The valid values are as follows:
  + **STREAMING**: Real time streaming.
  + **INDEX**: Index.
  + **APPLICATION**: Application.
  + **TENANT_BUCKET**: Tenant bucket.
  + **LAKE**: Data lake.

* `table_id` - (Optional, String) Specifies the table ID.

* `table_alias` - (Optional, String) Specifies the table alias.

* `table_name` - (Optional, String) Specifies the table name.

* `sort_key` - (Optional, String) Specifies the attribute fields for sorting.

* `sort_dir` - (Optional, String) Specifies the sorting order. Supported values are **ASC** and **DESC**.

* `exists` - (Optional, Bool) Specifies whether the table exists.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The tables list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `table_id` - The table ID.

* `pipe_id` - The pipe ID.

* `table_name` - The table name.

* `table_alias` - The table alias.

* `description` - The table description.

* `directory` - The directory group.

* `category` - The directory type.

* `lock_status` - The table lock status. Valid values are **LOCKED** and **UNLOCKED**.

* `process_status` - The processing status.
  Valid values are **COMPLETED**, **CREATING**, **UPDATING**, **DELETING**, **TRUNCATING**, **UPGRADING**,
  **CREATE_FAILED**, **UPDATE_FAILED**, **DELETE_FAILED**, **TRUNCATE_FAILED**, **UPGRADE_FAILED**.

* `process_error` - The table processing error. Valid values are **NONE**, **MISSING_ASSOCIATIONS**,
  **FAILED_INIT_STORAGE_RESOURCES_WHEN_CREATING**, **FAILED_INIT_FLINK_RESOURCES_WHEN_CREATING**,
  **FAILED_DELETE_FLINK_RESOURCES_WHEN_DELETING**, **FAILED_DELETE_STORAGE_RESOURCES_WHEN_DELETING**,
  **FAILED_DELETE_ALL_RESOURCES_WHEN_DELETING**, **FAILED_UPDATE_STORAGE_SETTING**, **FAILED_UPDATE_FLINK_SCHEMA**,
  **FAILED_UPDATE_STORAGE_SCHEMA**, **FAILED_TO_APPLY_RESOURCE**, **FAILED_TO_UPGRADE_RESOURCE_MODEL**.

* `format` - The table format. Valid values are **JSON**, **DEBEZIUM_JSON**, **CSV**, **PARQUET**, **ORC**.

* `rw_type` - The table read/write type. Valid values are **READ_ONLY**, **READ_WRITE**, **WRITE_ONLY**.

* `owner_type` - The owner type. Valid values are **SYSTEM**, **USER**, **SYSTEM_ALLOWED_DELETE**,
  **USER_ALLOWED_DELETE**.

* `data_layering` - The data layering. Valid values are **ODS**, **DWS**, **ADS**.

* `data_classification` - The data classification. Valid values are **FACTUAL_DATA** and **DIMENSION_DATA**.

* `schema` - The table schema.

  The [schema](#records_schema_struct) structure is documented below.

* `storage_setting` - The table storage setting.

  The [storage_setting](#records_storage_setting_struct) structure is documented below.

* `display_setting` - The table display setting.

  The [display_setting](#records_display_setting_struct) structure is documented below.

* `create_time` - The creation time, millisecond timestamp.

* `update_time` - The update time, millisecond timestamp.

* `delete_time` - The deletion time, millisecond timestamp.

<a name="records_schema_struct"></a>
The `schema` block supports:

* `columns` - The table columns list.

  The [columns](#records_schema_columns_struct) structure is documented below.

* `primary_key` - The table primary key list.

* `partition` - The table partition list.

* `watermark_column` - The table watermark column.

* `watermark_interval` - The table watermark delay interval.

* `time_filter` - The table time filter column.

<a name="records_schema_columns_struct"></a>
The `columns` block supports:

* `column_name` - The table column name.

* `column_type` - The column field type. Valid values are **PHYSICAL**, **METADATA**, **VIRTUAL_METADATA**,
  **COMPUTED**.

* `column_type_setting` - The table column type setting.

* `column_data_type` - The column field data type. Valid values are **ROW**, **MAP_STRING**, **MAP_DECIMAL**,
  **TINYINT**, **SMALLINT**, **INT**, **BIGINT**, **DECIMAL**, **FLOAT**, **DOUBLE**, **CHAR**, **VARCHAR**, **STRING**,
  **KEYWORD**, **BOOLEAN**, **DATE**, **TIME**, **TIMESTAMP**, **TIMESTAMP_LTZ**.

* `column_data_type_setting` - The table column data type setting.

* `nullable` - Whether the column is nullable.

* `array` - Whether the column is an array.

* `depth` - The depth.

* `parent_name` - The parent name.

* `own_name` - The own name.

* `column_display_setting` - The table column display setting.

  The [column_display_setting](#records_schema_columns_column_display_setting_struct) structure is documented below.

* `column_sequence_number` - The column sequence number.

<a name="records_schema_columns_column_display_setting_struct"></a>
The `column_display_setting` block supports:

* `mapping_required` - Whether mapping is required.

* `group_sequence_number` - The group sequence number.

* `intra_group_sequence_number` - The intra-group sequence number.

* `value_type` - The value type.

* `value_qualified` - The qualified value.

* `display_name` - The display name.

* `display_description` - The display description.

* `group_name` - The group name.

<a name="records_storage_setting_struct"></a>
The `storage_setting` block supports:

* `application_index` - The application index.

* `application_topic` - The application topic.

* `application_data_class_id` - The application data class ID.

* `streaming_bandwidth` - The streaming bandwidth (MB/s).

* `streaming_partition` - The streaming partition.

* `streaming_retention_size` - The streaming retention size.

* `streaming_dataspace_id` - The streaming dataspace ID.

* `index_storage_period` - The index storage period.

* `index_storage_size` - The index storage size.

* `index_shards` - The index shards number.

* `index_replicas` - The index replicas number.

* `lake_storage_period` - The data lake storage period.

* `lake_partition_setting` - The data lake partition setting.

* `lake_expiration_status` - The data lake partition status.

<a name="records_display_setting_struct"></a>
The `display_setting` block supports:

* `columns` - The table column display list.

  The [columns](#records_display_setting_columns_struct) structure is documented below.

* `format` - The table display settings. Supported values are **TABLE**, **RAW**.

<a name="records_display_setting_columns_struct"></a>
The `columns` block supports:

* `column_name` - The table column name.

* `column_alias` - The table column alias.

* `display_by_default` - Is it displayed by default.

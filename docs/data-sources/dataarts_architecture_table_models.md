---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_table_models"
description: |-
  Use this data source to get the list of the DataArts Architecture table models within HuaweiCloud.
---
# huaweicloud_dataarts_architecture_table_models

Use this data source to get the list of the DataArts Architecture table models within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "model_id" {}

data "huaweicloud_dataarts_architecture_table_models" "test" {
  workspace_id = var.workspace_id
  model_id     = var.model_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID of DataArts Architecture.

* `model_id` - (Required, String) Specifies the model ID to which the table model belongs.

* `name` - (Optional, String) Specifies the Chinese or English name of the table model.
  Fuzzy search is supported.

* `subject_id` - (Optional, String) Specifies the subject ID to which the table model belongs.

* `status` - (Optional, String) Specifies the status of the table model.  
  The valid values are as follows:
  + **DRAFT**
  + **PUBLISH_DEVELOPING**
  + **PUBLISHED**
  + **OFFLINE**
  + **REJECT**

* `created_by` - (Optional, String) Specifies the creator of the table model.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - All table models that match the filter parameters.
  The [tables](#architectrue_tables) structure is documented below.

<a name="architectrue_tables"></a>
The `tables` block supports:

* `id` - The ID of the table model.

* `model_id` - The model ID corresponding to the table model.

* `table_name` - The Chinese name of the table model.

* `physical_table_name` - The English name of the table model.

* `dw_id` - The ID of the connection corresponding to the table model.

* `dw_type` - The type of the connection corresponding to the table model.

* `catalog_path` - The subject path corresponding to the table model.

* `table_type` - The table type of the table model.

* `description` - The description of the table model.

* `configs` - The advanced configuration information of the table model, in JSON format.

* `parent_table_id` - The parent table ID of table model.

* `parent_table_name` - The parent table name of table model.

* `parent_table_code` - The parent table code of table model.

* `related_logic_table_id` - The logical table ID associated with table model.

* `owner` - The owner of the table model.

* `compression` - The compression type of the table model.

* `db_name` - The database name of the table model.

* `queue_name` - The queue name of the DLI table model.

* `schema` - The schema of the DWS and POSTGRESQL table model.

* `obs_location` - The OBS path of the table model.

* `attributes` - The list of the attributes of the table model.
  The [attributes](#tables_attributes) structure is documented below.

* `distribute` - The attribute distribution mode of the DWS table model.

* `distribute_column` - The HASH column of the attribute distribution.

* `code` - The code of the logical entity.

* `data_format` - The data format of the DLI table model.

* `dlf_task_id` - The DLF task ID of table model.

* `use_recently_partition` - Whether the table model has used latest partition.

* `dirty_out_switch` - The dirty data output switch of the table model.

* `dirty_out_database` - The database where to record the dirty data.

* `dirty_out_prefix` - The prefix of the table recording dirty data.

* `dirty_out_suffix` - The suffix of the table recording dirty data.

* `partition_conf` - The condition expression of the partition.

* `status` - The status of the table model.

* `extend_info` - The extend information of the table model.

* `is_partition` - Whether table is the partition table.

* `tb_guid` - The globally unique ID generated when publishing the table model.

* `logic_tb_guid` - The globally unique ID of the logic table model generated when publishing the table model.

* `logic_tb_id` - The ID of logical entity.

* `has_related_logic_table` - Whether the table model has associated the logical entities.

* `has_related_physical_table` - Whether the logical entity has associated the physical tables.

* `tb_id` - The ID of the data table.

* `physical_table_status` - The physical table status of the table model.

* `dev_physical_table_status` - The dev physical table status of the table model.

* `technical_asset_status` - The technical asset status of the table model.

* `business_asset_status` - The business asset status of the table model.

* `meta_data_link_status` - The meta data link status the table model.

* `data_quality_status` - The data quality status of the table model.

* `summary_status` - The synchronization status of the table model.

* `prod_version` - The production environment version of the table model.

* `dev_version` - The development environment version of the table model.

* `env_type` - The environment type of the table model.
  + **INVALID_TYPE**
  + **DEV_TYPE**
  + **PROD_TYPE**
  + **DEV_PROD_TYPE**

* `created_by` - The creator of the table model.

* `updated_by` - The latest updater of the table model.

* `created_at` - The creation time of the table model, in RFC3339 format.

* `updated_at` - The latest update time of the table model, in RFC3339 format.

<a name="tables_attributes"></a>
The `attributes` block supports:

* `id` - The ID of the attribute.

* `name` - The name of the attribute.

* `name_en` - The English name of the attribute.

* `data_type` - The data type of the attribute.

* `data_type_extend` - The data type extend field of attribute.

* `domain_type` - The domain type of the attribute.
  + **NUMBER**
  + **STRING**
  + **DATETIME**
  + **BLOB**: Large Object (BLOB).
  + **OTHER**

* `description` - The description of the attribute.

* `stand_row_id` - The ID of the data standard associated with attribute.

* `ordinal` - The sequence number of attribute.

* `code` - The code of the logical attribute associated with attribute.

* `extend_field` - The extend field of the attribute.

* `is_foreign_key` - Whether the attribute is foreign key.

* `is_primary_key` - Whether the attribute is primary key.

* `is_partition_key` - Whether the attribute is partition key.

* `not_null` - Whether the attribute is not null.

* `created_at` - The creation time of the attribute, in RFC3339 format.

* `updated_at` - The latest update time of the attribute, in RFC3339 format.

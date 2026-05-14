---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_aggregation_logic_tables"
description: |-
  Use this data source to query DataArts Architecture aggregation logic tables within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_aggregation_logic_tables

Use this data source to query DataArts Architecture aggregation logic tables within HuaweiCloud.

## Example Usage

### Query all aggregation logic tables under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "test" {
  workspace_id = var.workspace_id
}
```

### Filter aggregation logic table by name

```hcl
variable "workspace_id" {}
variable "aggregation_logic_table_name" {}

data "huaweicloud_dataarts_architecture_aggregation_logic_tables" "test" {
  workspace_id = var.workspace_id
  name         = var.aggregation_logic_table_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the aggregation logic tables are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the aggregation logic tables belong.

* `name` - (Optional, String) Specifies the Chinese name or English name of the aggregation logic table
  to be fuzzy queried.

* `name_ch` - (Optional, String) Specifies the Chinese name of the aggregation logic table to be exactly queried.

* `name_en` - (Optional, String) Specifies the English name of the aggregation logic table to be exactly queried.

* `create_by` - (Optional, String) Specifies the creator of the aggregation logic table to be queried.

* `approver` - (Optional, String) Specifies the approver of the aggregation logic table to be queried.

* `owner` - (Optional, String) Specifies the owner of the aggregation logic table to be queried.

* `status` - (Optional, String) Specifies the publishing status of the aggregation logic table to be queried.  
  The valid values are as follows:
  + **DRAFT**
  + **PUBLISH_DEVELOPING**
  + **PUBLISHED**
  + **OFFLINE_DEVELOPING**
  + **OFFLINE**
  + **REJECT**

* `sync_status` - (Optional, String) Specifies the synchronization status of the aggregation logic table
  to be queried.  
  The valid values are as follows:
  + **RUNNING**
  + **NO_NEED**
  + **SUMMARY_SUCCESS**
  + **SUMMARY_FAILED**

* `sync_key` - (Optional, List) Specifies the list of synchronization task types of the aggregation logic table.  
  The valid values are as follows:
  + **BUSINESS_ASSET**
  + **DATA_QUALITY**
  + **TECHNICAL_ASSET**
  + **META_DATA_LINK**
  + **PHYSICAL_TABLE**: Create table in production environment.
  + **DEV_PHYSICAL_TABLE**: Create table in development environment.
  + **DLF_TASK**
  + **MATERIALIZATION**
  + **PUBLISH_TO_DLM**
  + **SUMMARY_STATUS**

* `biz_catalog_id` - (Optional, String) Specifies the business catalog ID to which the aggregation logic table belongs.

* `begin_time` - (Optional, String) Specifies the start time of the modification time for the aggregation
  logic table, in RFC3339 format.  
  Must be `UTC` time, and must be used together with `end_time`.

* `end_time` - (Optional, String) Specifies the end time of the modification time for the aggregation
  logic table, in RFC3339 format.  
  Must be `UTC` time, and must be used together with `begin_time`.

* `auto_generate` - (Optional, Bool) Specifies whether the aggregation logic table is auto-generated.  
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tables` - The list of aggregation logic tables that match the filter parameters.  
  The [tables](#architecture_aggregation_logic_tables) structure is documented below.

<a name="architecture_aggregation_logic_tables"></a>
The `tables` block supports:

* `id` - The ID of the aggregation logic table.

* `tb_name` - The physical table name in English of the aggregation logic table.

* `tb_logic_name` - The display name in Chinese of the aggregation logic table.

* `dw_id` - The ID of the data connection.

* `dw_name` - The name of the data connection.

* `dw_type` - The type of the data connection.

* `db_name` - The name of the database.

* `owner` - The asset owner of the aggregation logic table.

* `description` - The description of the aggregation logic table.

* `alias` - The alias of the aggregation logic table.

* `queue_name` - The queue name of the DLI data connection.

* `schema` - The schema name of the aggregation logic table.

* `table_attributes` - The list of attributes of the aggregation logic table.  
  The [table_attributes](#architecture_aggregation_logic_tables_attributes) structure is documented below.

* `table_type` - The type of the aggregation logic table.

* `distribute` - The distribution mode of the database.

* `distribute_column` - The column used for hash distribution.

* `compression` - The compression level of the DWS table.

* `pre_combine_field` - The column used for table combine or versioning.

* `obs_location` - The OBS storage path of the external table.

* `configs` - The additional configuration of the aggregation logic table, in JSON format.

* `dimension_group` - The ID of the statistical dimension of the derived indicator.

* `group_name` - The name of the dimension group of the derivative metric.

* `group_code` - The dimension group code of the derivative metric.

* `sql` - The SQL statement of the aggregation logic table.

* `partition_conf` - The partition expression of the aggregation logic table.

* `dirty_out_switch` - Whether to enable dirty data output switch.

* `dirty_out_database` - The output database of the dirty data.

* `dirty_out_prefix` - The prefix of the dirty data table.

* `dirty_out_suffix` - The suffix of the dirty data table.

* `secret_type` - The type of the data secrecy level.

* `self_defined_fields` - The list of user-defined extended fields for the aggregation logic table.  
  The [self_defined_fields](#architecture_aggregation_logic_tables_self_defined_fields) structure is documented below.

* `tb_id` - The internal table ID of the aggregation logic table.

* `status` - The publishing status of the aggregation logic table.

* `model_id` - The ID of the model to which the aggregation logic table belongs.

* `create_by` - The creator of the aggregation logic table.

* `create_time` - The creation time of the aggregation logic table, in RFC3339 format.

* `update_time` - The latest update time of the aggregation logic table, in RFC3339 format.

* `env_type` - The publishing environment type of the aggregation logic table.

* `physical_table` - The status of the table creation in the production environment.

* `dev_physical_table` - The status of the table creation in the development environment.

* `technical_asset` - The synchronization status of the technical asset.

* `business_asset` - The synchronization status of the business asset.

* `meta_data_link` - The status of the asset association.

* `tb_guid` - The technical catalog asset GUID after the aggregation logic table is published.

* `tb_logic_guid` - The business catalog asset GUID after the aggregation logic table is published.

* `data_quality` - The status of the data quality job creation.

* `quality_id` - The ID of the data quality.

* `dlf_task` - The status of the DLF task.

* `dlf_task_id` - The ID of the DLF task.

* `publish_to_dlm` - The status of the DLM API generation.

* `api_id` - The ID of the API after publishing.

* `summary_status` - The synchronization status of the summary.

* `reversed` - Whether the aggregation logic table is reversed.

* `table_version` - The version of the aggregation logic table.  
  If the value is `2`, it means that the aggregation logic table is automatically aggregated.

* `dev_version` - The development environment version of the aggregation logic table.

* `prod_version` - The production environment version of the aggregation logic table.

* `dev_version_name` - The development environment version name of the aggregation logic table.

* `prod_version_name` - The production environment version name of the aggregation logic table.

* `approval_info` - The approval information of the aggregation logic table.  
  The [approval_info](#architecture_aggregation_logic_tables_approval_info) structure is documented below.

<a name="architecture_aggregation_logic_tables_attributes"></a>
The `table_attributes` block supports:

* `id` - The ID of the attribute.

* `name_ch` - The Chinese name of the attribute.

* `name_en` - The English name of the attribute.

* `data_type` - The data type of the attribute.

* `attribute_type` - The configuration type of the attribute.

* `is_primary_key` - Whether the attribute is the primary key.

* `is_partition_key` - Whether the attribute is used as a partition key.

* `not_null` - Whether the attribute is not null.

* `description` - The description of the attribute.

* `data_type_extend` - The data type extend field of attribute.

* `domain_type` - The domain type of the attribute.

* `ref_id` - The ID of the object referenced by the attribute.

* `ref_name_ch` - The Chinese name of the object associated with the attribute.

* `ref_name_en` - The English name of the object associated with the attribute.

* `stand_row_id` - The ID of the data standard associated with the attribute.

* `stand_row_name` - The name of the data standard associated with the attribute.

* `ordinal` - The sequence number of the attribute.

* `alias` - The alias of the attribute.

* `secrecy_levels` - The list of secrecy levels associated with the attribute.  
  The [secrecy_levels](#architecture_aggregation_logic_tables_secrecy_levels) structure is documented below.

<a name="architecture_aggregation_logic_tables_self_defined_fields"></a>
The `self_defined_fields` block supports:

* `fd_name_ch` - The Chinese display name of the custom extended field.

* `fd_name_en` - The English name of the custom extended field.

* `not_null` - Whether the custom extended field requires a value.

* `fd_value` - The value of the custom extended field.

<a name="architecture_aggregation_logic_tables_approval_info"></a>
The `approval_info` block supports:

* `id` - The ID of the approval for the aggregation logic table.

* `approver` - The approver of the aggregation logic table.

* `approval_status` - The approval status of the aggregation logic table.

* `msg` - The approval message for the aggregation logic table.

* `approval_time` - The approval time for the aggregation logic table, in RFC3339 format.

<a name="architecture_aggregation_logic_tables_secrecy_levels"></a>
The `secrecy_levels` block supports:

* `id` - The ID of the secrecy level.

* `uuid` - The UUID of the secrecy level.

* `name` - The name of the secrecy level.

* `slevel` - The secrecy level number.

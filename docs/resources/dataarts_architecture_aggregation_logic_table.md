---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_aggregation_logic_table"
description: |-
  Manages a DataArts Architecture aggregation logic table resource within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_aggregation_logic_table

Manages a DataArts Architecture aggregation logic table resource within HuaweiCloud.

## Example Usage

### Create aggregation logic table use DLI data connection

```hcl
variable "workspace_id" {}
variable "tb_name" {}
variable "tb_logic_name" {}
variable "subject_id" {}
variable "connection_id" {}
variable "owner" {}
variable "dli_database_name" {}
variable "dli_queue_name" {}
variable "description" {}
variable "table_attributes" {
  type = list(object({
    name_ch          = string
    name_en          = string
    data_type        = string
    attribute_type   = optional(string)
    is_primary_key   = optional(bool)
    is_partition_key = optional(bool)
    not_null         = optional(bool)
    ref_id           = optional(string)
    stand_row_id     = optional(string)
  }))
}

resource "huaweicloud_dataarts_architecture_aggregation_logic_table" "test" {
  workspace_id  = var.workspace_id
  tb_name       = var.tb_name
  tb_logic_name = var.tb_logic_name
  l3_id         = var.subject_id
  dw_id         = var.connection_id
  db_name       = var.dli_database_name
  dw_type       = "DLI"
  owner         = var.owner
  queue_name    = var.dli_queue_name
  description   = var.description
  table_type    = "MANAGED"

  dynamic "table_attributes" {
    for_each = var.table_attributes

    content {
      name_ch          = table_attributes.value.name_ch
      name_en          = table_attributes.value.name_en
      data_type        = table_attributes.value.data_type
      attribute_type   = table_attributes.value.attribute_type
      is_primary_key   = table_attributes.value.is_primary_key
      is_partition_key = table_attributes.value.is_partition_key
      not_null         = table_attributes.value.not_null
      ref_id           = table_attributes.value.ref_id
      stand_row_id     = table_attributes.value.stand_row_id
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the aggregation logic table is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the aggregation logic table belongs.

* `tb_name` - (Required, String) Specifies the physical table name in English of the aggregation logic table.  
  Only letters, digits, and underscores (_) are allowed, and must start with `dws_`.

* `tb_logic_name` - (Required, String) Specifies the display name in Chinese of the aggregation logic table.
  The characters `\<>%"';` and newline characters are not allowed.

* `l3_id` - (Required, String) Specifies the ID of the business subject to which the aggregation logic table belongs.

* `dw_id` - (Required, String) Specifies the ID of the data connection.

* `dw_type` - (Required, String) Specifies the type of the data connection.  
 The valid values are as follows:
  + **DWS**
  + **MRS_HIVE**
  + **POSTGRESQL**
  + **MRS_SPARK**
  + **CLICKHOUSE**
  + **MYSQL**
  + **ORACLE**
  + **DLI**
  + **DORIS**

* `db_name` - (Required, String) Specifies the name of the database corresponding to the data connection.

* `owner` - (Required, String) Specifies the asset owner of the aggregation logic table.

* `description` - (Required, String) Specifies the description of the aggregation logic table.

* `model_id` - (Optional, String) Specifies the ID of the model to which the aggregation logic table belongs.  
  If omitted, the aggregation logic table will be created in the default model.

* `alias` - (Optional, String) Specifies the alias of the aggregation logic table.

* `queue_name` - (Optional, String) Specifies the queue name corresponding to the DLI data connection.  
  This parameter is required when `dw_type` is **DLI**.

* `schema` - (Optional, String) Specifies the name of the database schema.  
  This parameter is required when `dw_type` is **DWS** or **POSTGRESQL**.

* `table_attributes` - (Optional, List) Specifies the list of attributes of the aggregation logic table.  
  The [table_attributes](#architecture_aggregation_logic_table_attributes) structure is documented below.

* `table_type` - (Optional, String) Specifies the type of the database table.  
  For the DLI data connection, the valid values are as follows:
  + **EXTERNAL**
  + **MANAGED**

  For the DWS data connection, the valid values are as follows:
  + **DWS_COLUMN**
  + **DWS_ROW**
  + **DWS_VIEW**

  For the MRS_HIVE data connection, the valid values are as follows:
  + **HIVE_EXTERNAL_TABLE**
  + **HIVE_TABLE**

  For the MRS_SPARK data connection, the valid values are as follows:
  + **HUDI_COW**
  + **HUDI_MOR**

  For the POSTGRESQL data connection, the valid values are as follows:
  + **POSTGRESQL_TABLE**

  For the CLICKHOUSE data connection, the valid values are as follows:
  + **CLICKHOUSE_TABLE**

  For the MYSQL data connection, the valid values are as follows:
  + **MYSQL_TABLE**

  For the ORACLE data connection, the valid values are as follows:
  + **ORACLE_TABLE**

  For the DORIS data connection, the valid values are as follows:
  + **DORIS_TABLE**

* `distribute` - (Optional, String) Specifies the distribution mode of the database.  
  The valid values are as follows:
  + **HASH**
  + **REPLICATION**

  This parameter is required when `dw_type` is **DWS**.

* `distribute_column` - (Optional, String) Specifies the column used for hash distribution.  
  This parameter is available when `distribute` is **HASH**.

* `compression` - (Optional, String) Specifies the compression level of the DWS table.  
  The valid values are as follows:
  + When `table_type` is **DWS_COLUMN**, the valid values are **NO**, **LOW**, **MIDDLE**, **HIGH**
  + When `table_type` is **DWS_ROW**, the valid values are **NO**, **YES**
  + When `table_type` is **DWS_VIEW**, setting the compression level is not supported.

  This parameter is required when `dw_type` is **DWS**.

* `pre_combine_field` - (Optional, String) Specifies the column used for record combine or versioning.  
  This parameter is available when `dw_type` is **MRS_SPARK**.

* `obs_location` - (Optional, String) Specifies the OBS storage path of the external table.  
  This parameter is available when `table_type` is **EXTERNAL**.

* `configs` - (Optional, String) Specifies the additional configuration of the aggregation logic table, in JSON format.

* `dimension_group` - (Optional, String) Specifies the ID of the statistical dimension of the derived indicator.

* `sql` - (Optional, String) Specifies the SQL statement of the aggregation logic table.

* `partition_conf` - (Optional, String) Specifies the partition expression of the aggregation logic table.

* `dirty_out_switch` - (Optional, Bool) Specifies whether to enable dirty data output switch.  
  Defaults to **false**.

* `dirty_out_database` - (Optional, String) Specifies the output database of the dirty data.

* `dirty_out_prefix` - (Optional, String) Specifies the prefix of the dirty data table.

* `dirty_out_suffix` - (Optional, String) Specifies the suffix of the dirty data table.

* `secret_type` - (Optional, String) Specifies the type of the data secrecy level.  
  The valid values are as follows:
  + **PUBLIC**
  + **SECRET**
  + **CONFIDENTIAL**
  + **SUPER_SECRET**

* `self_defined_fields` - (Optional, List) Specifies the list of user-defined extended fields for the aggregation logic
  table.  
  The [self_defined_fields](#architecture_aggregation_logic_table_self_defined_fields) structure is documented below.

* `del_type` - (Optional, String) Specifies the deletion type when the aggregation logic table is deleted.  
  The valid values are as follows:
  + **PHYSICAL_TABLE**: Whether to delete the database physical table.

<a name="architecture_aggregation_logic_table_attributes"></a>
The `table_attributes` block supports:

* `name_ch` - (Required, String) Specifies the Chinese name of the attribute.  
  The characters `\<>%"';` and newline characters are not allowed.

* `name_en` - (Required, String) Specifies the English name of the attribute.  
  Only letters, digits, and underscores (_) are allowed, and must start with a letter.

* `data_type` - (Required, String) Specifies the data type of the attribute.

* `attribute_type` - (Optional, String) Specifies the configuration type of the attribute.  
  The valid values are as follows:
  + **DERIVATIVE_INDEX**: Derived metric.
  + **COMPOSITE_METRIC**: Composite metric.
  + **BIZ_METRIC**: Business metric.
  + **SUMMARY_TIME**: Summary time.
  + **SUMMARY_DIMENSION_ATTRIBUTE**: Summary dimension attribute.

* `is_primary_key` - (Optional, Bool) Specifies whether the attribute is the primary key.  
  Defaults to **false**.

* `is_partition_key` - (Optional, Bool) Specifies whether the attribute is used as a partition key.  
  Defaults to **false**.

* `not_null` - (Optional, Bool) Specifies whether the attribute is not null.  
  Defaults to **false**.

* `description` - (Optional, String) Specifies the description of the attribute.

* `data_type_extend` - (Optional, String) Specifies the data type extend field of attribute.

* `ref_id` - (Optional, String) Specifies the ID of the object referenced by the attribute.

* `stand_row_id` - (Optional, String) Specifies the ID of the data standard associated with the attribute.

* `alias` - (Optional, String) Specifies the alias of the attribute.

* `secrecy_levels` - (Optional, List) Specifies the list of secrecy levels associated with the attribute.  
  The [secrecy_levels](#architecture_aggregation_logic_table_secrecy_levels) structure is documented below.

<a name="architecture_aggregation_logic_table_self_defined_fields"></a>
The `self_defined_fields` block supports:

* `fd_name_ch` - (Optional, String) Specifies the Chinese display name of the custom extended field.

* `fd_name_en` - (Optional, String) Specifies the English name of the custom extended field.

* `not_null` - (Optional, Bool) Specifies whether the custom extended field requires a value.  
  Defaults to **false**.

* `fd_value` - (Optional, String) Specifies the value of the custom extended field.

<a name="architecture_aggregation_logic_table_secrecy_levels"></a>
The `secrecy_levels` block supports:

* `id` - (Required, String) Specifies the ID of the secrecy level.  
  This parameter must be in UUID format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `dw_name` - The name of the data connection.

* `status` - The publishing status of the aggregation logic table.  
  + **DRAFT**
  + **PUBLISH_DEVELOPING**: Publish pending review.
  + **PUBLISHED**
  + **OFFLINE_DEVELOPING**: Offline pending review.
  + **OFFLINE**
  + **REJECT**

* `group_name` - The name of the dimension group of the derivative metric.

* `group_code` - The dimension group code of the derivative metric.

* `created_by` - The account name of the user who created the aggregation logic table.

* `create_time` - The creation time of the aggregation logic table, in RFC3339 format.

* `update_time` - The latest update time of the aggregation logic table, in RFC3339 format.

* `env_type` - The publishing environment type of the aggregation logic table.  
  + **INVALID_TYPE**: Invalid environment.
  + **DEV_TYPE**: Development environment.
  + **PROD_TYPE**: Production environment.
  + **DEV_PROD_TYPE**: Development and production environment.

* `physical_table` - The status of the table creation in the production environment after the aggregation logic table
  is published.

* `dev_physical_table` - The status of the table creation in the development environment after the aggregation logic
  table is published.

* `technical_asset` - The synchronization status of the technical asset after the aggregation logic table is published.

* `business_asset` - The synchronization status of the business asset after the aggregation logic table is published.

* `meta_data_link` - The status of the asset association after the aggregation logic table is published.

* `tb_guid` - The technical catalog asset GUID after the aggregation logic table is published.

* `tb_logic_guid` - The business catalog asset GUID after the aggregation logic table is published.

* `data_quality` - The status of the data quality job creation after the aggregation logic table is published.

* `quality_id` - The ID of the data quality.

* `dlf_task` - The status of the DLF task after the aggregation logic table is published.

* `dlf_task_id` - The ID of the DLF task after the aggregation logic table is published.

* `publish_to_dlm` - The status of the DLM API generation after the aggregation logic table is published.

* `api_id` - The ID of the API after the aggregation logic table is published.

* `summary_status` - The synchronization status of the summary after the aggregation logic table is published.

* `table_attributes` - The list of attributes of the aggregation logic table.  
  The [table_attributes](#architecture_aggregation_logic_table_attributes_attr) structure is documented below.

<a name="architecture_aggregation_logic_table_attributes_attr"></a>
The `table_attributes` block supports:

* `id` - The ID of the attribute.

* `domain_type` - The domain type of the attribute.

* `ref_name_ch` - The Chinese name of the object associated with the attribute.

* `ref_name_en` - The English name of the object associated with the attribute.

* `stand_row_name` - The name of the data standard associated with the attribute.

* `secrecy_levels` - The list of secrecy levels associated with the attribute.  
  The [secrecy_levels](#architecture_aggregation_logic_table_secrecy_levels_attr) structure is documented below.

* `ordinal` - The sequence number of the attribute.

<a name="architecture_aggregation_logic_table_secrecy_levels_attr"></a>
The `secrecy_levels` block supports:

* `uuid` - The UUID of the secrecy level.

* `name` - The name of the secrecy level.

## Import

The resource can be imported using the `workspace_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_aggregation_logic_table.test <workspace_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `l3_id`, `secret_type` and `del_type`.

It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_dataarts_architecture_aggregation_logic_table" "test" {
  ...

  lifecycle {
    ignore_changes = [
      l3_id, secret_type, del_type,
    ]
  }
}
```

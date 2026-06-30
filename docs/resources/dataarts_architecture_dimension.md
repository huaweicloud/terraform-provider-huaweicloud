---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_dimension"
description: |-
  Manages a DataArts Architecture dimension resource within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_dimension

Manages a DataArts Architecture dimension resource within HuaweiCloud.

## Example Usage

### Create a common dimension

```hcl
variable "workspace_id" {}
variable "name_ch" {}
variable "name_en" {}
variable "subject_id" {}
variable "connection_id" {}
variable "owner" {}
variable "description" {}

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = var.workspace_id
  name_ch        = var.name_ch
  name_en        = var.name_en
  l3_id          = var.subject_id
  dimension_type = "COMMON"
  owner          = var.owner
  description    = var.description

  attributes {
    name_ch        = "attr1_ch"
    name_en        = "dim_attr1_en"
    data_type      = "STRING"
    is_primary_key = true
    ordinal        = 1
  }

  attributes {
    name_ch        = "attr2_ch"
    name_en        = "dim_attr2_en"
    data_type      = "BIGINT"
    is_primary_key = false
    ordinal        = 2
  }

  datasource {
    dw_id   = var.connection_id
    dw_type = "DLI"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dimension is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the dimension belongs.

* `name_ch` - (Required, String) Specifies the Chinese name of the dimension.

* `name_en` - (Required, String) Specifies the English name of the dimension.

  -> The `name_en` must start with **dim_**, such as **dim_test_name**.

* `l3_id` - (Required, String) Specifies the ID of the business object to which the dimension belongs.

* `dimension_type` - (Required, String) Specifies the type of the dimension.  
  The valid values are as follows:
  + **COMMON**: Common dimension.
  + **LOOKUP**: Code table dimension.
  + **HIERARCHIES**: Hierarchy dimension.

* `owner` - (Required, String) Specifies the asset owner of the dimension.

* `description` - (Required, String) Specifies the description of the dimension.

* `attributes` - (Required, List) Specifies the list of dimension attributes.  
  The [attributes](#architecture_dimension_attributes_arg) structure is documented below.

* `datasource` - (Required, List) Specifies the data source information.  
  The [datasource](#architecture_dimension_datasource_arg) structure is documented below.

* `code_table_id` - (Optional, String) Specifies the code table ID of the dimension.

* `is_delete_physical_table` - (Optional, Bool) Specifies whether to delete the physical table corresponding to
  the dimension when deleting the dimension.

<a name="architecture_dimension_attributes_arg"></a>
The `attributes` block supports:

* `name_en` - (Required, String) Specifies the English name of the attribute.

* `name_ch` - (Required, String) Specifies the Chinese name of the attribute.

* `ordinal` - (Required, Int) Specifies the sequence number of the attribute.

* `data_type` - (Required, String) Specifies the data type of the attribute.

* `is_primary_key` - (Required, Bool) Specifies whether the attribute is the primary key.

* `description` - (Optional, String) Specifies the description of the attribute.

* `is_biz_primary` - (Optional, Bool) Specifies whether the attribute is the business primary key.

* `is_partition_key` - (Optional, Bool) Specifies whether the attribute is the partition key.

* `not_null` - (Optional, Bool) Specifies whether the attribute is not null.

<a name="architecture_dimension_datasource_arg"></a>
The `datasource` block supports:

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

* `biz_type` - (Optional, String) Specifies the business type of the data source.

* `db_name` - (Optional, String) Specifies the name of the database corresponding to the data connection.

* `queue_name` - (Optional, String) Specifies the queue name corresponding to the DLI data connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The publishing status of the dimension.
  + **DRAFT**
  + **PUBLISH_DEVELOPING**
  + **PUBLISHED**
  + **OFFLINE_DEVELOPING**
  + **OFFLINE**
  + **REJECT**

* `created_by` - The account name of the user who created the dimension.

* `create_time` - The creation time of the dimension, in RFC3339 format.

* `update_time` - The latest update time of the dimension, in RFC3339 format.

* `l1_id` - The ID of the subject domain group.

* `l2_id` - The ID of the subject domain.

* `l1_name` - The Chinese name of the subject domain group.

* `l2_name` - The Chinese name of the subject domain.

* `l3_name` - The Chinese name of the business object.

* `table_type` - The table type of the dimension.

* `distribute` - The distribute type of the dimension.

* `distribute_column` - The distribute column of the dimension.

* `compression` - The compression type of the dimension.

* `obs_location` - The OBS location of the dimension.

* `pre_combine_field` - The pre-combine field of the dimension.

* `alias` - The alias of the dimension.

* `configs` - The configs of the dimension.

* `env_type` - The environment type of the dimension.

* `model_id` - The model ID of the dimension.

* `model` - The model information of the dimension.  
  The [model](#architecture_dimension_model_attr) structure is documented below.

* `update_by` - The updater of the dimension.

* `code_table_id` - The code table ID of the dimension.

* `dev_version` - The development environment version of the dimension.

* `prod_version` - The production environment version of the dimension.

* `dev_version_name` - The development environment version name of the dimension.

* `prod_version_name` - The production environment version name of the dimension.

* `attributes` - The list of dimension attributes.  
  The [attributes](#architecture_dimension_attributes_attr) structure is documented below.

* `datasource` - The data source information.  
  The [datasource](#architecture_dimension_datasource_attr) structure is documented below.

<a name="architecture_dimension_attributes_attr"></a>
The `attributes` block supports:

* `id` - The ID of the attribute.

* `domain_type` - The domain type of the attribute.

* `code_table_field_id` - The code table field ID of the attribute.

* `create_by` - The creator of the attribute.

* `data_type_extend` - The data type extension of the attribute.

* `stand_row_id` - The associated data standard ID of the attribute.

* `stand_row_name` - The associated data standard name of the attribute.

* `quality_infos` - The quality information of the attribute.

* `secrecy_levels` - The secrecy levels of the attribute.

* `status` - The publishing status of the attribute.

* `create_time` - The creation time of the attribute, in RFC3339 format.

* `update_time` - The latest update time of the attribute, in RFC3339 format.

* `alias` - The alias of the attribute.

* `self_defined_fields` - The self-defined fields of the attribute.

* `description` - The description of the attribute.

* `is_biz_primary` - Whether the attribute is the business primary key.

* `is_partition_key` - Whether the attribute is the partition key.

* `not_null` - Whether the attribute is not null.

<a name="architecture_dimension_datasource_attr"></a>
The `datasource` block supports:

* `id` - The ID of the data source.

* `biz_type` - The business type of the data source.

* `biz_id` - The business ID of the data source.

* `dw_name` - The name of the data connection.

* `db_name` - The name of the database.

* `queue_name` - The queue name for DLI data connection.

* `schema` - The schema name for DWS data connection.

<a name="architecture_dimension_model_attr"></a>
The `model` block supports:

* `id` - The ID of the workspace.

* `name` - The name of the workspace.

* `description` - The description of the workspace.

* `is_physical` - Whether it is a physical table.

* `frequent` - Whether it is frequently used.

* `top` - Whether it is a top-level governance.

* `level` - The data governance level.

* `dw_type` - The data warehouse type.

* `create_time` - The creation time of the workspace, in RFC3339 format.

* `update_time` - The latest update time of the workspace, in RFC3339 format.

* `create_by` - The creator of the workspace.

* `update_by` - The updater of the workspace.

* `type` - The workspace type.

* `biz_catalog_ids` - The associated business catalog IDs.

* `databases` - The database names.

* `table_model_prefix` - The table model prefix.

## Import

The resource can be imported using the `workspace_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_architecture_dimension.test <workspace_id>/<id>
```

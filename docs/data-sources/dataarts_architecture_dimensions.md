---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_dimensions"
description: |-
  Use this data source to query DataArts Architecture dimensions within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_dimensions

Use this data source to query DataArts Architecture dimensions within HuaweiCloud.

## Example Usage

### Query all dimensions under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_dimensions" "test" {
  workspace_id = var.workspace_id
}
```

### Query dimensions by name and type

```hcl
variable "workspace_id" {}
variable "dimension_name" {}

data "huaweicloud_dataarts_architecture_dimensions" "test" {
  workspace_id    = var.workspace_id
  name            = var.dimension_name
  dimension_type  = "COMMON"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dimensions are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the dimensions belong.

* `name` - (Optional, String) Specifies the name or code of the dimension to be fuzzy queried.

* `name_ch` - (Optional, String) Specifies the Chinese name of the dimension to be exactly queried.

* `name_en` - (Optional, String) Specifies the English name of the dimension to be exactly queried.

* `create_by` - (Optional, String) Specifies the creator of the dimension to be queried.

* `approver` - (Optional, String) Specifies the approver of the dimension to be queried.

* `status` - (Optional, String) Specifies the publishing status of the dimension to be queried.

* `l2_id` - (Optional, String) Specifies the subject domain L2 ID to which the dimension belongs.

* `dimension_type` - (Optional, String) Specifies the type of the dimension to be queried.

* `biz_catalog_id` - (Optional, String) Specifies the business catalog ID to which the dimension belongs.

* `fact_logic_id` - (Optional, String) Specifies the fact table ID of the dimension to be queried.

* `begin_time` - (Optional, String) Specifies the start time of the dimension to be queried, in RFC3339 format.

* `end_time` - (Optional, String) Specifies the end time of the dimension to be queried, in RFC3339 format.

* `derivative_ids` - (Optional, List) Specifies the derivative indicator ID list for querying dimensions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dimensions` - The list of dimensions that matched filter parameters.  
  The [dimensions](#dataarts_architecture_dimensions_attr) structure is documented below.

<a name="dataarts_architecture_dimensions_attr"></a>
The `dimensions` block supports:

* `id` - The ID of the dimension.

* `name_ch` - The Chinese name of the dimension.

* `name_en` - The English name of the dimension.

* `dimension_type` - The type of the dimension.

* `description` - The description of the dimension.

* `owner` - The asset owner of the dimension.

* `status` - The publishing status of the dimension.
  + **DRAFT**
  + **PUBLISH_DEVELOPING**
  + **PUBLISHED**
  + **OFFLINE_DEVELOPING**
  + **OFFLINE**
  + **REJECT**

* `code_table_id` - The code table ID of the dimension.

* `l1_id` - The L1 subject domain grouping ID of the dimension.

* `l2_id` - The L2 subject domain ID of the dimension.

* `l3_id` - The L3 business object ID of the dimension.

* `l1_name` - The L1 subject domain grouping name of the dimension.

* `l2_name` - The L2 subject domain name of the dimension.

* `l3_name` - The L3 business object name of the dimension.

* `created_by` - The creator of the dimension.

* `updated_by` - The updater of the dimension.

* `created_at` - The creation time of the dimension, in RFC3339 format.

* `updated_at` - The latest update time of the dimension, in RFC3339 format.

* `table_type` - The table type of the dimension.

* `distribute` - The distribute type of the dimension.

* `distribute_column` - The distribute column of the dimension.

* `obs_location` - The OBS location of the dimension.

* `alias` - The alias of the dimension.

* `configs` - The configs of the dimension.

* `env_type` - The environment type of the dimension.

* `model_id` - The model ID of the dimension.

* `dev_version` - The development environment version of the dimension.

* `prod_version` - The production environment version of the dimension.

* `dev_version_name` - The development environment version name of the dimension.

* `prod_version_name` - The production environment version name of the dimension.

* `datasource` - The data source configuration of the dimension.  
  The [datasource](#dataarts_architecture_dimensions_datasource_attr) structure is documented below.

* `attributes` - The list of attributes of the dimension.  
  The [attributes](#dataarts_architecture_dimensions_attributes_attr) structure is documented below.

* `hierarchies` - The hierarchy attributes of the dimension.  
  The [hierarchies](#dataarts_architecture_dimensions_hierarchies_attr) structure is documented below.

* `code_table` - The code table information of the dimension.  
  The [code_table](#dataarts_architecture_dimensions_code_table_attr) structure is documented below.

* `model` - The model information of the dimension.  
  The [model](#dataarts_architecture_dimensions_model_attr) structure is documented below.

<a name="dataarts_architecture_dimensions_datasource_attr"></a>
The `datasource` block supports:

* `id` - The ID of the data source.

* `biz_type` - The business type of the data source.

* `biz_id` - The business ID of the data source.

* `dw_id` - The ID of the data connection.

* `dw_type` - The type of the data connection.

* `dw_name` - The name of the data connection.

* `db_name` - The name of the database corresponding to the data connection.

* `queue_name` - The queue name corresponding to the DLI data connection.

* `schema` - The name of the database schema.

<a name="dataarts_architecture_dimensions_attributes_attr"></a>
The `attributes` block supports:

* `id` - The ID of the attribute.

* `name_en` - The English name of the attribute.

* `name_ch` - The Chinese name of the attribute.

* `description` - The description of the attribute.

* `data_type` - The data type of the attribute.

* `domain_type` - The domain type of the attribute.

* `data_type_extend` - The data type extension of the attribute.

* `is_primary_key` - Whether the attribute is the primary key.

* `is_biz_primary` - Whether the attribute is the business primary key.

* `is_partition_key` - Whether the attribute is the partition key.

* `ordinal` - The sequence number of the attribute.

* `not_null` - Whether the attribute is not null.

* `code_table_field_id` - The code table field ID of the attribute.

* `create_by` - The creator of the attribute.

* `stand_row_id` - The associated data standard ID of the attribute.

* `stand_row_name` - The associated data standard name of the attribute.

* `quality_infos` - The quality information of the attribute.

* `secrecy_levels` - The secrecy levels of the attribute.

* `status` - The publishing status of the attribute.

* `create_time` - The creation time of the attribute, in RFC3339 format.

* `update_time` - The latest update time of the attribute, in RFC3339 format.

* `alias` - The alias of the attribute.

* `self_defined_fields` - The self-defined fields of the attribute.

<a name="dataarts_architecture_dimensions_hierarchies_attr"></a>
The `hierarchies` block supports:

* `id` - The ID of the hierarchy.

* `name` - The name of the hierarchy.

* `created_by` - The creator of the hierarchy.

* `updated_by` - The updater of the hierarchy.

* `created_at` - The creation time of the hierarchy, in RFC3339 format.

* `updated_at` - The latest update time of the hierarchy, in RFC3339 format.

* `attrs` - The attributes of the hierarchy.  
  The [attrs](#dataarts_architecture_dimensions_hierarchies_attrs_attr) structure is documented below.

<a name="dataarts_architecture_dimensions_hierarchies_attrs_attr"></a>
The `attrs` block supports:

* `id` - The ID of the hierarchy attribute.

* `hierarchies_id` - The hierarchy ID of the attribute.

* `attr_id` - The attribute ID.

* `level` - The level of the hierarchy attribute.

* `attr_name_en` - The English name of the referenced attribute.

* `attr_name_ch` - The Chinese name of the referenced attribute.

<a name="dataarts_architecture_dimensions_code_table_attr"></a>
The `code_table` block supports:

* `id` - The ID of the code table.

* `name_en` - The English name of the code table.

* `name_ch` - The Chinese name of the code table.

* `tb_version` - The version of the code table.

* `directory_id` - The directory ID of the code table.

* `directory_path` - The directory path of the code table.

* `description` - The description of the code table.

* `status` - The publishing status of the code table.

<a name="dataarts_architecture_dimensions_model_attr"></a>
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

---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_dimension"
description: |-
  Use this resource to manage a DataArts Architecture dimension resource within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_dimension

Use this resource to manage a DataArts Architecture dimension resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name_ch" {}
variable "name_en" {}
variable "l3_id" {}
variable "owner" {}
variable "description" {}
variable "connection_id" {}
variable "attributes" {
  type = list(object({
    name_ch        = string,
    name_en        = string,
    data_type      = string,
    is_primary_key = bool,
    ordinal        = string,
  }))
}

resource "huaweicloud_dataarts_architecture_dimension" "test" {
  workspace_id   = var.workspace_id
  name_ch        = var.name_ch
  name_en        = var.name_en
  dimension_type = "DIMENSION"
  l3_id          = var.l3_id
  owner          = var.owner
  description    = var.description

  datasource {
    dw_id   = var.connection_id
    dw_type = "DWS"
  }

  dynamic "attributes" {
    for_each = var.attributes
    content {
      name_ch        = attributes.value.name_ch
      name_en        = attributes.value.name_en
      data_type      = attributes.value.data_type
      is_primary_key = attributes.value.is_primary_key
      ordinal        = attributes.value.ordinal
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the dimension is located.  
  If omitted, the provider-level region will be used.  
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the dimension belongs.  
  Changing this parameter will create a new resource.

* `name_ch` - (Required, String) Specifies the Chinese name of the dimension.

* `name_en` - (Required, String) Specifies the English name of the dimension.

* `dimension_type` - (Required, String) Specifies the type of the dimension.

* `l3_id` - (Required, String) Specifies the business object ID to which the dimension belongs.

* `owner` - (Required, String) Specifies the asset owner of the dimension.

* `attributes` - (Required, List) Specifies the dimension attribute information.  
  The [attributes](#dataarts_architecture_dimension_attributes) structure is documented below.

* `datasource` - (Required, List) Specifies the data source information of the dimension.  
  The [datasource](#dataarts_architecture_dimension_datasource) structure is documented below.

* `alias` - (Optional, String) Specifies the alias of the dimension.

* `code_table_id` - (Optional, String) Specifies the referenced code table ID.

* `configs` - (Optional, String) Specifies the other configuration information.

* `create_by` - (Optional, String) Specifies the creator of the dimension.

* `description` - (Optional, String) Specifies the description of the dimension.

* `distribute` - (Optional, String) Specifies the distribution mode.

* `distribute_column` - (Optional, String) Specifies the DISTRIBUTE BY HASH column.

* `hierarchies` - (Optional, List) Specifies the hierarchy attribute definitions of the dimension.  
  The [hierarchies](#dataarts_architecture_dimension_hierarchies) structure is documented below.

* `id_field` - (Optional, String) Specifies the ID field of the dimension.

* `l1` - (Optional, String) Specifies the Chinese name of the subject domain group.

* `l2` - (Optional, String) Specifies the Chinese name of the subject domain.

* `l2_id` - (Optional, String) Specifies the subject domain ID.

* `l3` - (Optional, String) Specifies the Chinese name of the business object.

* `mappings` - (Optional, List) Specifies the table mapping information of the dimension.  
  The [mappings](#dataarts_architecture_dimension_mappings) structure is documented below.

* `model` - (Optional, List) Specifies the model information of the dimension.  
  The [model](#dataarts_architecture_dimension_model) structure is documented below.

* `model_id` - (Optional, String) Specifies the model ID to which the dimension belongs.

* `obs_location` - (Optional, String) Specifies the OBS external table path.

* `self_defined_fields` - (Optional, List) Specifies the custom extended fields of the dimension.  
  The [self_defined_fields](#dataarts_architecture_dimension_self_defined_fields) structure is documented below.

* `status` - (Optional, String) Specifies the publish status of the dimension.

* `table_type` - (Optional, String) Specifies the table type of the dimension.

* `update_by` - (Optional, String) Specifies the updater of the dimension.

<a name="dataarts_architecture_dimension_attributes"></a>
The `attributes` block supports:

* `name_ch` - (Required, String) Specifies the Chinese name of the attribute.

* `name_en` - (Required, String) Specifies the English name of the attribute.

* `data_type` - (Required, String) Specifies the data type of the attribute.

* `is_primary_key` - (Required, Bool) Specifies whether the attribute is a primary key.

* `ordinal` - (Required, String) Specifies the sequence number of the attribute.

* `alias` - (Optional, String) Specifies the alias of the attribute.

* `code_table_field_id` - (Optional, String) Specifies the code table field ID of the attribute.

* `create_by` - (Optional, String) Specifies the creator of the attribute.

* `data_type_extend` - (Optional, String) Specifies the data type extend field of the attribute.

* `description` - (Optional, String) Specifies the description of the attribute.

* `id` - (Optional, String) Specifies the ID of the attribute.

* `is_biz_primary` - (Optional, Bool) Specifies whether the attribute is a business primary key.

* `is_partition_key` - (Optional, Bool) Specifies whether the attribute is a partition key.

* `not_null` - (Optional, Bool) Specifies whether the attribute is not null.

* `self_defined_fields` - (Optional, List) Specifies the custom extended fields of the attribute.  
  The [self_defined_fields](#dataarts_architecture_dimension_self_defined_fields) structure is documented below.

* `stand_row_id` - (Optional, String) Specifies the ID of the associated data standard.

<a name="dataarts_architecture_dimension_self_defined_fields"></a>
The `self_defined_fields` block supports:

* `fd_name_ch` - (Optional, String) Specifies the Chinese display name of the custom extended field.

* `fd_name_en` - (Optional, String) Specifies the English name of the custom extended field.

* `fd_value` - (Optional, String) Specifies the value of the custom extended field.

* `not_null` - (Optional, Bool) Specifies whether the custom extended field requires a value.

<a name="dataarts_architecture_dimension_datasource"></a>
The `datasource` block supports:

* `dw_id` - (Required, String) Specifies the data connection ID.

* `dw_type` - (Required, String) Specifies the data connection type.

* `biz_id` - (Optional, String) Specifies the business object ID.

* `biz_type` - (Optional, String) Specifies the business object type.

* `db_name` - (Optional, String) Specifies the database name.

* `id` - (Optional, String) Specifies the data source ID.

* `queue_name` - (Optional, String) Specifies the DLI queue name.

* `schema` - (Optional, String) Specifies the DWS schema name.

<a name="dataarts_architecture_dimension_hierarchies"></a>
The `hierarchies` block supports:

* `id` - (Optional, String) Specifies the hierarchy ID.

* `name` - (Optional, String) Specifies the hierarchy name.

<a name="dataarts_architecture_dimension_mappings"></a>
The `mappings` block supports:

* `name` - (Required, String) Specifies the mapping name.

* `source_tables` - (Required, List) Specifies the source table information of the mapping.  
  The [source_tables](#dataarts_architecture_dimension_mappings_source_tables) structure is documented below.

* `details` - (Optional, List) Specifies the mapping details.  
  The [details](#dataarts_architecture_dimension_mappings_details) structure is documented below.

* `description` - (Optional, String) Specifies the mapping description.

* `source_fields` - (Optional, List) Specifies the source field information of the mapping.
  The [source_fields](#dataarts_architecture_dimension_mappings_source_fields) structure is documented below.

* `src_model_id` - (Optional, String) Specifies the source model ID in relational modeling.

* `src_model_name` - (Optional, String) Specifies the source model name in relational modeling.

* `target_table_id` - (Optional, String) Specifies the target table ID.

* `target_table_name` - (Optional, String) Specifies the target table name.

* `view_text` - (Optional, String) Specifies the collected view source.

<a name="dataarts_architecture_dimension_mappings_source_tables"></a>
The `source_tables` block supports:

* `table1_id` - (Required, String) Specifies the table 1 ID.

* `table2_id` - (Optional, String) Specifies the table 2 ID.

* `join_type` - (Optional, String) Specifies the join type.  
  The valid values are as follows:
  + **LEFT**
  + **RIGHT**
  + **INNER**
  + **FULL**

* `table1_name` - (Optional, String) Specifies the table 1 name.

* `table2_name` - (Optional, String) Specifies the table 2 name.

* `join_fields` - (Optional, List) Specifies the ON condition fields.  
  The [join_fields](#dataarts_architecture_dimension_mappings_source_fields_join_fields) structure is documented below.

<a name="dataarts_architecture_dimension_mappings_source_fields_join_fields"></a>
The `join_fields` block supports:

* `field1_id` - (Required, String) Specifies the field 1 ID.

* `field2_id` - (Required, String) Specifies the field 2 ID.

* `field1_name` - (Optional, String) Specifies the field 1 name.

* `field2_name` - (Optional, String) Specifies the field 2 name.

<a name="dataarts_architecture_dimension_mappings_source_fields"></a>
The `source_fields` block supports:

* `target_field_name` - (Required, String) Specifies the target field code.

* `field_ids` - (Optional, String) Specifies the source field IDs, multiple IDs separated by commas.

* `field_names` - (Optional, List) Specifies the source field name list.

* `target_field_id` - (Optional, String) Specifies the target field ID.

* `transform_expression` - (Optional, String) Specifies the transform expression.

<a name="dataarts_architecture_dimension_mappings_details"></a>
The `details` block supports:

* `target_attr_name` - (Required, String) Specifies the target attribute name.

* `id` - (Optional, String) Specifies the detail ID.

* `mapping_id` - (Optional, String) Specifies the mapping ID.

* `remark` - (Optional, String) Specifies the remark of the mapping detail.

* `src_attr_ids` - (Optional, String) Specifies the source attribute IDs.

* `src_table_ids` - (Optional, String) Specifies the source table IDs.

* `target_attr_id` - (Optional, String) Specifies the target attribute ID.

<a name="dataarts_architecture_dimension_model"></a>
The `model` block supports:

* `name` - (Required, String) Specifies the workspace name.

* `type` - (Required, String) Specifies the workspace type.

* `id` - (Optional, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attributes` - The dimension attribute information.  
  The [attributes](#dataarts_architecture_dimension_attributes_attr) structure is documented below.

* `hierarchies` - The hierarchy attribute definitions of the dimension.  
  The [hierarchies](#dataarts_architecture_dimension_hierarchies_attr) structure is documented below.

* `mappings` - The table mapping information of the dimension.  
  The [mappings](#dataarts_architecture_dimension_mappings_attr) structure is documented below.

* `approval_info` - The approval information.  
  The [approval_info](#dataarts_architecture_dimension_approval_info_attr) structure is documented below.

* `code_table` - The referenced code table.  
  The [code_table](#dataarts_architecture_dimension_code_table_attr) structure is documented below.

* `create_time` - The creation time.

* `dev_version` - The development environment version.

* `dev_version_name` - The development environment version name.

* `env_type` - The development and production environment type.

* `l1_id` - The subject domain group ID.

* `new_biz` - The business version management information.  
  The [new_biz](#dataarts_architecture_dimension_new_biz_attr) structure is documented below.

* `prod_version` - The production environment version.

* `prod_version_name` - The production environment version name.

<a name="dataarts_architecture_dimension_attributes_attr"></a>
The `attributes` block supports:

* `create_time` - The creation time of the attribute.

* `dimension_id` - The dimension ID of the attribute.

* `domain_type` - The domain type of the attribute.

* `stand_row_name` - The name of the associated data standard.

* `update_time` - The update time of the attribute.

* `status` - The publish status of the attribute.

* `quality_infos` - The quality information of the attribute.  
  The [quality_infos](#dataarts_architecture_dimension_attributes_quality_info_attr) structure is documented below.

* `secrecy_levels` - The secrecy levels of the attribute.  
  The [secrecy_levels](#dataarts_architecture_dimension_attributes_secrecy_level_attr) structure is documented below.

<a name="dataarts_architecture_dimension_attributes_quality_info_attr"></a>
The `quality_infos` block supports:

* `id` - The ID of the quality info.

* `alert_conf` - The alert configuration.

* `attr_id` - The attribute ID.

* `biz_type` - The business entity type.

* `create_by` - The creator of the quality info.

* `create_time` - The creation time.

* `data_quality_id` - The data quality ID.

* `data_quality_name` - The data quality name.

* `expression` - The regular expression configuration.

* `extend_info` - The extended information.

* `from_standard` - Whether it is from data standard quality configuration.

* `result_description` - The result description.

* `show_control` - Whether to display the regular expression.

* `table_id` - The table ID.

* `update_by` - The updater of the quality info.

* `update_time` - The update time.

<a name="dataarts_architecture_dimension_attributes_secrecy_level_attr"></a>
The `secrecy_levels` block supports:

* `id` - The secrecy level ID.

* `name` - The secrecy level name.

* `slevel` - The secrecy level rank.

* `description` - The description of the secrecy level.

* `uuid` - The data security primary key.

* `create_by` - The creator of the secrecy level.

* `create_time` - The creation time.

* `update_by` - The updater of the secrecy level.

* `update_time` - The update time.

* `new_biz` - The business version management information.  
  The [new_biz](#dataarts_architecture_dimension_secrecy_levels_new_biz_attr) structure is documented below.

<a name="dataarts_architecture_dimension_hierarchies_attr"></a>
The `hierarchies` block supports:

* `attrs` - The attributes contained in the hierarchy.  
  The [attrs](#dataarts_architecture_dimension_hierarchies_attrs_attr) structure is documented below.

* `create_by` - The creator of the hierarchy.

* `created_at` - The creation time.

* `updated_by` - The updater of the hierarchy.

* `updated_at` - The update time.

<a name="dataarts_architecture_dimension_hierarchies_attrs_attr"></a>
The `attrs` block supports:

* `attr` - The referenced attribute details.  
  The [attr](#dataarts_architecture_dimension_hierarchies_attr_item_attr) structure is documented below.

* `detail_attrs` - The detail attributes.  
  The [detail_attrs](#dataarts_architecture_dimension_hierarchies_detail_attrs_attr) structure is documented below.

* `attr_id` - The attribute ID.

* `attr_name_en` - The referenced attribute code.

* `detail_attr_ids` - The detail attribute IDs.

* `detail_attr_name_ens` - The detail attribute English names.

* `hierarchies_id` - The hierarchy ID.

* `id` - The hierarchy attribute ID.

* `level` - The hierarchy level.

<a name="dataarts_architecture_dimension_hierarchies_attr_item_attr"></a>
The `attr` block supports:

* `id` - The attribute ID.

* `name_ch` - The Chinese name of the attribute.

* `name_en` - The English name of the attribute.

* `data_type` - The data type of the attribute.

* `is_primary_key` - Whether it is a primary key.

* `ordinal` - The sequence number of the attribute.

* `create_time` - The creation time.

* `update_time` - The update time.

* `domain_type` - The domain type of the attribute.

* `status` - The publish status of the attribute.

* `alias` - The alias of the attribute.

* `code_table_field_id` - The code table field ID.

* `create_by` - The creator of the attribute.

* `data_type_extend` - The data type extend field of the attribute.

* `description` - The description of the attribute.

* `dimension_id` - The dimension ID.

* `is_biz_primary` - Whether it is a business primary key.

* `is_partition_key` - Whether it is a partition key.

* `not_null` - Whether it is not null.

* `stand_row_id` - The associated data standard ID.

* `stand_row_name` - The associated data standard name.

* `quality_infos` - The quality information of the attribute.  
  The [quality_infos](#dataarts_architecture_dimension_attributes_quality_info_attr) structure is documented below.

* `secrecy_levels` - The secrecy levels of the attribute.  
  The [secrecy_levels](#dataarts_architecture_dimension_attributes_secrecy_level_attr) structure is documented below.

* `self_defined_fields` - The custom extended fields of the attribute.  
  The [self_defined_fields](#dataarts_architecture_dimension_self_defined_fields_attr) structure is documented below.

<a name="dataarts_architecture_dimension_hierarchies_detail_attrs_attr"></a>
The `detail_attrs` block supports:

* `alias` - The alias of the attribute.

* `code_table_field_id` - The code table field ID.

* `create_by` - The creator of the attribute.

* `create_time` - The creation time.

* `data_type` - The data type of the attribute.

* `data_type_extend` - The data type extend field of the attribute.

* `description` - The description of the attribute.

* `dimension_id` - The dimension ID.

* `domain_type` - The domain type of the attribute.

* `id` - The attribute ID.

* `is_biz_primary` - Whether it is a business primary key.

* `is_partition_key` - Whether it is a partition key.

* `is_primary_key` - Whether it is a primary key.

* `name_ch` - The Chinese name of the attribute.

* `name_en` - The English name of the attribute.

* `not_null` - Whether it is not null.

* `ordinal` - The sequence number of the attribute.

* `stand_row_id` - The associated data standard ID.

* `stand_row_name` - The associated data standard name.

* `status` - The publish status of the attribute.

* `update_time` - The update time.

* `quality_infos` - The quality information of the attribute.  
  The [quality_infos](#dataarts_architecture_dimension_attributes_quality_info_attr) structure is documented below.

* `secrecy_levels` - The secrecy levels of the attribute.  
  The [secrecy_levels](#dataarts_architecture_dimension_attributes_secrecy_level_attr) structure is documented below.

* `self_defined_fields` - The custom extended fields of the attribute.  
  The [self_defined_fields](#dataarts_architecture_dimension_self_defined_fields_attr) structure is documented below.

<a name="dataarts_architecture_dimension_self_defined_fields_attr"></a>
The `self_defined_fields` block supports:

* `fd_name_ch` - The Chinese display name of the custom extended field.

* `fd_name_en` - The English name of the custom extended field.

* `fd_value` - The value of the custom extended field.

* `not_null` - Whether the custom extended field requires a value.

<a name="dataarts_architecture_dimension_mappings_attr"></a>
The `mappings` block supports:

* `id` - The mapping ID.

* `details` - The mapping details.  
  The [details](#dataarts_architecture_dimension_mappings_details_attr) structure is documented below.

* `source_fields` - The source field information of the mapping.  
  The [source_fields](#dataarts_architecture_dimension_mappings_source_fields_attr) structure is documented below.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `created_by` - The creator of the mapping.

* `updated_by` - The updater of the mapping.

<a name="dataarts_architecture_dimension_mappings_details_attr"></a>
The `details` block supports:

* `target_attr_name` - The target attribute name.

* `id` - The detail ID.

* `mapping_id` - The mapping ID.

* `remark` - The remark of the mapping detail.

* `src_attr_ids` - The source attribute IDs.

* `src_table_ids` - The source table IDs.

* `target_attr_id` - The target attribute ID.

* `create_by` - The creator of the mapping detail.

* `create_time` - The creation time.

* `update_by` - The updater of the mapping detail.

* `update_time` - The update time.

<a name="dataarts_architecture_dimension_mappings_source_fields_attr"></a>
The `source_fields` block supports:

* `changed` - Whether the field has changed.

* `target_field_id` - The target field ID.

<a name="dataarts_architecture_dimension_approval_info_attr"></a>
The `approval_info` block supports:

* `approval_status` - The approval status.

* `approval_time` - The approval time.

* `approval_type` - The approval type.

* `approver` - The approver.

* `biz_id` - The business ID.

* `biz_info` - The serialized business details.

* `id` - The approval ID.

* `msg` - The approval information.

* `name_ch` - The business Chinese name.

* `name_en` - The business English name.

* `submit_time` - The submit time.

* `tenant_id` - The project ID.

* `directory_path` - The directory tree.

* `email` - The approver email.

* `biz_info_obj` - The business details object.  
  The [biz_info_obj](#dataarts_architecture_dimension_approval_info_biz_info_obj_attr) structure is documented below.

<a name="dataarts_architecture_dimension_approval_info_biz_info_obj_attr"></a>
The `biz_info_obj` block supports:

* `biz_status` - The publish status of the business.

* `biz_type` - The business entity type.

* `biz_version` - The business version.

* `create_by` - The creator of the business.

* `directory_path` - The directory tree.

* `email` - The email address of the creator.

* `id` - The business information ID.

* `l1` - The subject domain group Chinese name.

* `l2` - The subject domain Chinese name.

* `l3` - The business object Chinese name.

* `msg` - The approval or submission message.

* `name_ch` - The Chinese name of the business.

* `name_en` - The English name of the business.

* `submit_time` - The submit time.

* `tenant_id` - The tenant ID.

<a name="dataarts_architecture_dimension_code_table_attr"></a>
The `code_table` block supports:

* `id` - The code table ID.

* `name_ch` - The Chinese name of the code table.

* `name_en` - The English name of the code table.

* `status` - The publish status of the code table.

* `tb_version` - The version of the code table.

* `create_by` - The creator of the code table.

* `create_time` - The creation time.

* `directory_path` - The directory tree.

* `description` - The description of the code table.

* `directory_id` - The directory ID.

* `update_time` - The update time.

* `new_biz` - The business version management information.  
  The [new_biz](#dataarts_architecture_dimension_code_table_new_biz_attr) structure is documented below.

* `approval_info` - The approval information.  
  The [approval_info](#dataarts_architecture_dimension_code_table_approval_info_attr) structure is documented below.

* `code_table_fields` - The code table field information.  
  The [code_table_fields](#dataarts_architecture_dimension_code_table_fields_attr) structure is documented below.

<a name="dataarts_architecture_dimension_code_table_new_biz_attr"></a>
The `new_biz` block supports:

* `id` - The business version ID.

* `biz_type` - The business entity type.

<a name="dataarts_architecture_dimension_code_table_approval_info_attr"></a>
The `approval_info` block supports:

* `approval_status` - The approval status.

* `approval_time` - The approval time.

* `approval_type` - The approval type.

* `approver` - The approver.

* `biz_id` - The business ID.

* `biz_info` - The serialized business details.

* `id` - The approval ID.

* `msg` - The approval information.

* `name_ch` - The business Chinese name.

* `name_en` - The business English name.

* `submit_time` - The submit time.

* `tenant_id` - The project ID.

* `directory_path` - The directory tree.

* `email` - The approver email.

* `biz_info_obj` - The business details object.  
  The [biz_info_obj](#dataarts_architecture_dimension_code_table_approval_info_biz_info_obj_attr) structure is documented below.

<a name="dataarts_architecture_dimension_code_table_approval_info_biz_info_obj_attr"></a>
The `biz_info_obj` block supports:

* `biz_status` - The publish status of the business.

* `biz_type` - The business entity type.

* `biz_version` - The business version.

* `create_by` - The creator of the business.

* `directory_path` - The directory tree.

* `email` - The email address of the creator.

* `id` - The business information ID.

* `l1` - The subject domain group Chinese name.

* `l2` - The subject domain Chinese name.

* `l3` - The business object Chinese name.

* `msg` - The approval or submission message.

* `name_ch` - The Chinese name of the business.

* `name_en` - The English name of the business.

* `submit_time` - The submit time.

* `tenant_id` - The tenant ID.

<a name="dataarts_architecture_dimension_code_table_fields_attr"></a>
The `code_table_fields` block supports:

* `id` - The code table field ID.

* `code_table_id` - The code table ID to which the field belongs.

* `name_ch` - The Chinese name of the field.

* `name_en` - The English name of the field.

* `ordinal` - The sequence number of the field.

* `data_type` - The data type of the field.

* `data_type_extend` - The extended data type information of the field.

* `description` - The description of the field.

* `domain_type` - The domain type of the field.

* `is_unique_key` - Whether the field is unique.

* `count_field_values` - The total count of field values.

* `code_table_field_values` - The code table field values.  
  The [code_table_field_values](#dataarts_architecture_dimension_code_table_field_values_attr) structure is documented below.

<a name="dataarts_architecture_dimension_code_table_field_values_attr"></a>
The `code_table_field_values` block supports:

* `id` - The code table field value ID.

* `fd_id` - The code table attribute ID.

* `fd_value` - The code table attribute value.

* `ordinal` - The sequence number of the field value.

<a name="dataarts_architecture_dimension_new_biz_attr"></a>
The `new_biz` block supports:

* `id` - The business version ID.

* `biz_type` - The business entity type.

<a name="dataarts_architecture_dimension_secrecy_levels_new_biz_attr"></a>
The `new_biz` block supports:

* `id` - The business version ID.

* `biz_type` - The business entity type.

## Import

DataArts Architecture dimension can be imported using `<workspace_id>/<id>`, e.g.

```bash
terraform import huaweicloud_dataarts_architecture_dimension.test <workspace_id>/<id>
```

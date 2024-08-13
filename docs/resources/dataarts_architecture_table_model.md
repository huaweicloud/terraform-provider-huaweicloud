---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_table_model"
description: ""
---

# huaweicloud_dataarts_architecture_table_model

Manages DataArts Architecture table model resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "model_id" {}
variable "model_name" {}
variable "subject_id" {}
variable "table_model_name" {}
variable "table_model_id" {}
variable "attribute_name" {}
variable "attribute_id" {}
variable "relation_name" {}
variable "mapping_name" {}
variable "attribute_names" {}
variable "attribute_ids" {}
variable "transform_expression" {}

resource "huaweicloud_dataarts_architecture_table_model" "test" {
  workspace_id        = var.workspace_id
  model_id            = var.model_id
  subject_id          = var.subject_id
  physical_table_name = var.table_model_name
  table_name          = var.table_model_name
  description         = "demo"
  dw_type             = "UNSPECIFIED"

  attributes {
    name      = "key"
    name_en   = "key_en"
    data_type = "STRING"
    ordinal   = "1"
  }

  relations {
    name              = var.relation_name
    source_type       = "ONE"
    target_table_id   = var.table_model_id
    target_table_name = var.table_model_name
    target_type       = "ONE"
    mappings {
      source_field_name = var.attribute_name
      target_field_id   = var.attribute_id
      target_field_name = var.attribute_name
    }
  }

  mappings {
    name           = var.mapping_name
    src_model_id   = var.model_id
    src_model_name = var.model_name
    source_tables {
      table1_id   = var.table_model_id
      table1_name = var.table_model_name
      table2_id   = var.table_model_id
      table2_name = var.table_model_name
      join_type   = "LEFT"
      join_fields {
        field1_id   = var.attribute_id
        field1_name = var.attribute_name
        field2_id   = var.attribute_id
        field2_name = var.attribute_name
      }
    }
    source_fields {
      target_field_name    = var.attribute_name
      field_ids            = var.attribute_ids
      field_names          = var.attribute_names
      transform_expression = var.transform_expression
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the table model.
  If omitted, the provider-level region will be used. Changing this creates a new table model.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID which the table model in.
  Changing this creates a new table model.

* `model_id` - (Required, String, ForceNew) Specifies the model ID which the table model in.
  Changing this creates a new table model.

* `dw_type` - (Required, String, ForceNew) Specifies the data connection type of table model.
  Changing this creates a new table model.
  + The valid values for physical model are: **DWS**, **DLI**, **MRS_HIVE**, **POSTGRESQL**, **MRS_SPARK**,
  **CLICKHOUSE**, **MYSQL**, **ORACLE**.
  + The valid value for logical model is **UNSPECIFIED**.

* `subject_id` - (Required, String) Specifies the subject ID to which table model belongs.

* `physical_table_name` - (Required, String) Specifies the physical name of table model.

* `table_name` - (Required, String) Specifies the name of table model.

* `description` - (Required, String) Specifies the description of table model.

* `attributes` - (Required, List) Specifies the attributes of table model.
  The [attributes](#block--attributes) structure is documented below.

* `del_type` - (Optional, String) Specifies the delete type of table model. The valid value is **PHYSICAL_TABLE**.

* `configs` - (Optional, String) Specifies the advanced configs of table model.

* `dw_id` - (Optional, String) Specifies the data connection ID of table model.

* `dw_name` - (Optional, String) Specifies the data connection name of table model.

* `db_name` - (Optional, String) Specifies the database name of table model.

* `queue_name` - (Optional, String) Specifies the queue name of table model.

* `schema` - (Optional, String) Specifies the schema of table model.

* `obs_location` - (Optional, String) Specifies the obs location of table model.

* `parent_table_id` - (Optional, String) Specifies the parent table ID of table model.

* `related_logic_table_model_id` - (Optional, String) Specifies the related logic table model ID of table model.

* `related_logic_model_id` - (Optional, String) Specifies the related logic model ID of table model.

* `parent_table_name` - (Optional, String) Specifies the parent table name of table model.

* `related_logic_table_model_name` - (Optional, String) Specifies the related logic table model name of table model.

* `related_logic_model_name` - (Optional, String) Specifies the related logic model name of table model.

* `owner` - (Optional, String) Specifies the owner of table model.

* `dirty_out_switch` - (Optional, Bool) Specifies the dirty out switch of table model.

* `dirty_out_database` - (Optional, String) Specifies the database where to record the dirty data.

* `dirty_out_prefix` - (Optional, String) Specifies the prefix of the table recording dirty data.

* `dirty_out_suffix` - (Optional, String) Specifies the suffix of the table recording dirty data.

* `partition_conf` - (Optional, String) Specifies the **WHERE** condition expression.

* `relations` - (Optional, List) Specifies the relations of table model.
  The [relations](#block--relations) structure is documented below.

* `mappings` - (Optional, List) Specifies the mappings of table model.
  The [mappings](#block--mappings) structure is documented below.

* `table_type` - (Optional, String) Specifies the table type of table model.
  + The valid values for **DWS** are **DWS_COLUMN**, **DWS_ROW**.
  + The valid values for **DLI** are **EXTERNAL**, **MANAGED**.
  + The valid values for **MRS_HIVE** are **HIVE_EXTERNAL_TABLE**, **HIVE_TABLE**.
  + The valid value for **POSTGRESQL** is **POSTGRESQL_TABLE**.
  + The valid values for **MRS_SPARK** are **HUDI_COW**, **HUDI_MOR**.
  + The valid value for **CLICKHOUSE** is **CLICKHOUSE_TABLE**.
  + The valid value for **MYSQL** is **MYSQL_TABLE**.
  + The valid value for **ORACLE** is **ORACLE_TABLE**.
  + The valid value for **UNSPECIFIED** is **LOGIC_TABLE**.

* `compression` - (Optional, String) Specifies the compression level of table model.
  The valid values are **YES**, **NO**. It's **Required** for **DWS** table model.

* `code` - (Optional, String) Specifies the code of table model.

* `distribute` - (Optional, String) Specifies the attribute is distributed by what. The valid values are **HASH**,
  **REPLICATION**.

* `distribute_column` - (Optional, String) Specifies the **HASH** column the attribute distributed by.

* `data_format` - (Optional, String) Specifies the data format of table model. Only for **DLI** table model.
  + The valid values for **MANAGED** are **Parquet**, **Carbon**.
  + The valid values for **EXTERNAL** are **Parquet**, **Carbon**, **CSV**, **ORC**, **JSON**, **Avro**.

* `dlf_task_id` - (Optional, String) Specifies the DLF task ID of table model.

* `use_recently_partition` - (Optional, Bool) Specifies the table model use recently partition or not.

* `reversed` - (Optional, Bool) Specifies whether the table model is reversed.

<a name="block--attributes"></a>
The `attributes` block supports:

* `name` - (Required, String) Specifies the name of attribute.

* `name_en` - (Required, String) Specifies the English name of attribute.

* `data_type` - (Required, String) Specifies the data type of attribute.

* `data_type_extend` - (Optional, String) Specifies the data type extend field of attribute.

* `description` - (Optional, String) Specifies the description of attribute.

* `stand_row_id` - (Optional, String) Specifies the data standard ID of attribute.

* `related_logic_attr_id` - (Optional, String) Specifies the related logic attribute ID of attribute.

* `related_logic_attr_name` - (Optional, String) Specifies the related logic attribute name of attribute.

* `related_logic_attr_name_en` - (Optional, String) Specifies the related logic attribute English name of attribute.

* `stand_row_name` - (Optional, String) Specifies the data standard name of attribute.

* `ordinal` - (Optional, String) Specifies the sequence number of attribute. The input values must start from one, and must
  be continuous numbers.

* `code` - (Optional, String) Specifies the code of attribute.

* `extend_field` - (Optional, Bool) Specifies the extend field of attribute.

* `is_foreign_key` - (Optional, Bool) Specifies the attribute is foreign key or not.

* `is_partition_key` - (Optional, Bool) Specifies the attribute is partition key or not.

* `is_primary_key` - (Optional, Bool) Specifies the attribute is primary key or not.

* `not_null` - (Optional, Bool) Specifies the attribute is not null or null.

<a name="block--relations"></a>
The `relations` block supports:

* `mappings` - (Required, List) Specifies the mappings of the attributes related.
  The [mappings](#block--relations--mappings) structure is documented below.

* `name` - (Required, String) Specifies the name of the relation.

* `source_type` - (Required, String) Specifies the relation type of source to target.
  The valid values are **ONE**, **ZERO_OR_ONE**, **ZERO_OR_N**, **ONE_OR_N**.

* `target_type` - (Required, String) Specifies the relation type of target to source.
  The valid values are **ONE**, **ZERO_OR_ONE**, **ZERO_OR_N**, **ONE_OR_N**.

* `role` - (Optional, String) Specifies the role of the relation.

* `source_table_id` - (Optional, String) Specifies the source table ID. Source table ID and target table ID, one of them
  must be the resource ID, so it have to be empty and the other one is **Required**.

* `source_table_name` - (Optional, String) Specifies the source table name.

* `target_table_id` - (Optional, String) Specifies the target table ID. Source table ID and target table ID, one of them
  must be the resource ID, so it have to be empty and the other one is **Required**.

* `target_table_name` - (Optional, String) Specifies the target table name.

<a name="block--relations--mappings"></a>
The `mappings` block supports:

* `source_field_id` - (Optional, String) Specifies the source attribute ID. Source field ID and target field ID, one of
  them must be the resource attribute ID, so it have to be empty and input its name, the other one ID is **Required**.

* `source_field_name` - (Optional, String) Specifies the source attribute English name. If the source attribute ID is
  empty, it's **Required**.

* `target_field_id` - (Optional, String) Specifies the source attribute ID. Source field ID and target field ID, one of
  them must be the resource attribute ID, so it have to be empty and the other one is **Required**.

* `target_field_name` - (Optional, String) Specifies the target attribute English name. If the target attribute ID is
  empty, it's **Required**.

<a name="block--mappings"></a>
The `mappings` block supports:

* `name` - (Required, String) Specifies the mapping name.

* `source_tables` - (Required, List) Specifies the source table informations of mapping.
  The [source_tables](#block--mappings--source_tables) structure is documented below.

* `src_model_id` - (Optional, String) Specifies the source model ID. It's **Required** for physical table model.

* `src_model_name` - (Optional, String) Specifies the source model name.

* `view_text` - (Optional, String) Specifies the source to capturing the view, using for **DWS** reversed view.

* `source_fields` - (Optional, List) Specifies the source attribute informations of mapping.
  The [source_fields](#block--mappings--source_fields) structure is documented below.

<a name="block--mappings--source_fields"></a>
The `source_fields` block supports:

* `field_ids` - (Optional, String) Specifies the source attribute IDs of mapping. Using **,** to split ID.

* `field_names` - (Optional, List) Specifies the source attribute English name list of mapping.

* `target_field_name` - (Optional, String) Specifies the attribute English name of this resource.

* `transform_expression` - (Optional, String) Specifies the transform expression.

<a name="block--mappings--source_tables"></a>
The `source_tables` block supports:

* `table1_id` - (Required, String) Specifies the table id.

* `join_type` - (Optional, String) Specifies the join type of two table.
  The valid values are **LEFT**, **RIGHT**, **INNER**, **FULL**.

* `table2_id` - (Optional, String) Specifies the table id.

* `join_fields` - (Optional, List) Specifies the attribute informations.
  The [join_fields](#block--mappings--source_tables--join_fields) structure is documented below.

* `table1_name` - (Optional, String) Specifies the table name.

* `table2_name` - (Optional, String) Specifies the table name

<a name="block--mappings--source_tables--join_fields"></a>
The `join_fields` block supports:

* `field1_id` - (Required, String) Specifies the attribute ID.

* `field2_id` - (Required, String) Specifies the attribute ID.

* `field1_name` - (Optional, String) Specifies the attribute English name.

* `field2_name` - (Optional, String) Specifies the attribute English name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attributes` - The attributes of the table model.
  The [attributes](#attrblock--attributes) structure is documented below.

* `parent_table_code` - The parent table code of the table model.

* `extend_info` - The extend information of the table model.

* `tb_guid` - The globally unique ID of the table model, generating when the table model publish.

* `logic_tb_guid` - The globally unique ID of the logic table model, generating when the table model publish.

* `status` - The status of the table model.

* `catalog_path` - The subject path.

* `created_at` - The creating time of the table model.

* `updated_at` - The updating time of the table model.

* `created_by` - The person creating the table model.

* `updated_by` - The person updating the table model.

* `has_related_physical_table` - The table model has related physical table or not.

* `has_related_logic_table` - The table model has related logic table or not.

* `is_partition` - The table model is partition or not.

* `physical_table_status` - The physical table status of the table model.

* `dev_physical_table_status` - The dev physical table status of the table model.

* `technical_asset_status` - The technical asset status of the table model.

* `business_asset_status` - The business asset status of the table model.

* `meta_data_link_status` - The meta data link status the table model.

* `data_quality_status` - The data quality status of the table model.

* `summary_status` - The summary status the table model.

* `env_type` - The env type of the table model.

* `relations` - The relations of table model.
  The [relations](#attrblock--relations) structure is documented below.

* `mappings` - The mappings of table model.
  The [mappings](#attrblock--mappings) structure is documented below.

<a name="attrblock--attributes"></a>
The `attributes` block supports:

* `id` - The ID of the attribute.

* `created_at` - The creating time of the attribute.

* `updated_at` - The updating time of the attribute.

* `domain_type` - The domain type of the attribute.

<a name="attrblock--relations"></a>
The `relations` block supports:

* `id` - The ID of the relation.

* `created_at` - The creating time of the relation.

* `updated_at` - The updating time of the relation.

* `created_by` - The person creating the relation.

* `updated_by` - The person updating the relation.

* `mappings` - The mapping of the attributes related.
  The [mappings](#attrblock--relations--mappings) structure is documented below.

<a name="attrblock--relations--mappings"></a>
The `mappings` block supports:

* `id` - The ID of the mapping.

* `created_at` - The creating time of the mapping

* `updated_at` - The updating time of the mapping.

* `created_by` - The person creating the mapping.

* `updated_by` - The person updating the mapping.

<a name="attrblock--mappings"></a>
The `mappings` block supports:

* `id` - The ID of the mapping.

* `created_at` - The creating time of the mapping

* `updated_at` - The updating time of the mapping.

* `created_by` - The person creating the mapping.

* `updated_by` - The person updating the mapping.

* `source_fields` - The source attributes of mapping.
  The [source_fields](#attrblock--mappings--source_fields) structure is documented below.

<a name="attrblock--mappings--source_fields"></a>
The `source_fields` block supports:

* `changed` - The attributes changed or not.

## Import

DataArts Architecture table model can be imported using `<workspace_id>/<id>`, e.g.

```sh
terraform import huaweicloud_dataarts_architecture_table_model.test <workspace_id>/<id>
```

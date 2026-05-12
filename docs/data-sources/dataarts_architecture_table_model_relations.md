---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_table_model_relations"
description: |-
  Use this data source to get the list of DataArts architecture table model relations.
---

# huaweicloud_dataarts_architecture_table_model_relations

Use this data source to get the list of DataArts architecture table model relations.

## Example Usage

```hcl
variable "workspace_id" {}
variable "model_id" {}

data "huaweicloud_dataarts_architecture_table_model_relations" "test" {
  workspace_id = var.workspace_id
  model_id     = var.model_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the table model relations are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the table model belongs.

* `model_id` - (Required, String) Specifies the model ID to which the table model relations belong.

* `biz_type` - (Optional, String) Specifies the business type.  
  The valid values are as follows:
  + **TABLE_MODEL**
  + **FACT_LOGIC_TABLE**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `relations` - The list of the table model relations matched filter parameters.  
  The [relations](#dataarts_architecture_table_model_relation) structure is documented below.

<a name="dataarts_architecture_table_model_relation"></a>
The `relations` block supports:

* `id` - The ID of the relation.

* `relation_type` - The type of the relation.

* `source_table_id` - The ID of the source table.

* `source_table_name` - The name of the source table.

* `target_table_id` - The ID of the target table.

* `target_table_name` - The name of the target table.

* `created_at` - The creation time of the relation, in RFC3339 format.

* `updated_at` - The update time of the relation, in RFC3339 format.

* `created_by` - The creator of the relation.

* `updated_by` - The updater of the relation.

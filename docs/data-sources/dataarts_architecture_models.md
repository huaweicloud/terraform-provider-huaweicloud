---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_models"
description: |-
  Use this data source to get the list of DataArts architecture models.
---

# huaweicloud_dataarts_architecture_models

Use this data source to get the list of DataArts architecture models.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_models" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the architecture models are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace.

* `workspace_type` - (Optional, String) Specifies the model workspace type.  
  The valid values are as follows:
  + **THIRD_NF**: relational modeling
  + **DIMENSION**: dimensional modeling

* `dw_type` - (Optional, String) Specifies the data connection type.  
  The valid values are as follows:
  + **DWS**
  + **MRS_HIVE**
  + **POSTGRESQL**
  + **MRS_SPARK**
  + **CLICKHOUSE**
  + **MYSQL**
  + **ORACLE**
  + **DORIS**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `models` - The list of the architecture models that matched filter parameters.  
  The [models](#dataarts_architecture_model) structure is documented below.

<a name="dataarts_architecture_model"></a>
The `models` block supports:

* `id` - The ID of the model.

* `name` - The name of the workspace.

* `description` - The description of the model.

* `is_physical` - Whether it is a physical table.

* `frequent` - Whether it is frequently used.

* `top` - Whether it is a top-level governance.

* `level` - The data governance layer.  
  + **SDI**
  + **DWI**
  + **DWR**
  + **DM**

* `dw_type` - The data connection type.

* `create_time` - The creation time of the model, in RFC3339 format.

* `update_time` - The update time of the model, in RFC3339 format.

* `create_by` - The creator of the model.

* `update_by` - The updater of the model.

* `type` - The workspace type.  
  + **THIRD_NF**
  + **DIMENSION**
  + **DM**

* `biz_catalog_ids` - The associated business catalog IDs.

* `databases` - The database name list.

* `table_model_prefix` - The table model validation prefix.

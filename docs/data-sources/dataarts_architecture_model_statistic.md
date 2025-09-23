---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_model_statistic"
description: |-
  Use this data source to get the list of model statistic of the DataArts Architecture within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_model_statistic

Use this data source to get the list of model statistic of the DataArts Architecture within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_model_statistic" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID of DataArts Architecture.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `frequent` - The list of the frequently-used objects.
  The [frequent](#model_statistic_object_attr) structure is documented below.

* `tops` - The list of the first-layer models.
  The [tops](#model_statistic_object_attr) structure is documented below.

* `logics` - The list of the logical models.
  The [logics](#model_statistic_object_attr) structure is documented below.

* `physicals` - The list of the physical models.
  The [physicals](#model_statistic_object_attr) structure is documented below.

* `dwr` - The DWR data reporting layer.
  The [dwr](#model_statistic_object_attr) structure is documented below.

* `dm` - The DM data integration layer.
  The [dm](#model_statistic_object_attr) structure is documented below.

<a name="model_statistic_object_attr"></a>
The `frequent`, `tops`, `logics`, `physicals`, `dwr` and `dm` block supports:

* `biz_type` - The service type.

* `level` - The model level.

* `db` - The total database number.

* `tb` - The total data table number.

* `tb_published` - The published data table number.

* `fd` - The total field number.

* `fd_published` - The published field number.

* `st` - The standard coverage.

* `st_published` - The published standard coverage.

* `model` - The model detail.
  The [model](#model_statistic_object_model_attr) structure is documented below.

<a name="model_statistic_object_model_attr"></a>
The `model` block supports:

* `id` - The ID of the model (workspace).

* `name` - The name of the model (workspace).

* `description` - The description of the model (workspace).

* `is_physical` - Whether a table is a physical table.

* `frequent` - Whether the model (workspace) is frequently used.

* `top` - Whether the model (workspace) is hierarchical governance.

* `level` - The level of the data governance layering.  
  The valid values are as follows:
  + **SDI**
  + **DWI**
  + **DWR**
  + **DM**

* `dw_type` - The type of the data connection.

* `created_at` - The creation time of the model (workspace), in RFC3339 format.

* `updated_at` - The latest update time of the model (workspace), in RFC3339 format.

* `created_by` - The person who creates the model (workspace).

* `updated_by` - The person who updates the model (workspace).

* `type` - The type of the model (workspace).  
  The valid values are as follows:
  + **THIRD_NF**: relational modeling
  + **DIMENSION**: dimensional modeling

* `biz_catalog_ids` - The ID list of associated service catalogs.

* `databases` - The list of database names.

* `table_model_prefix` - The verification prefix of the model (workspace).

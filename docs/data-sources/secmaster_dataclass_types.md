---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_dataclass_types"
description: |-
  Use this data source to get the list of SecMaster dataclass types within HuaweiCloud.
---

# huaweicloud_secmaster_dataclass_types

Use this data source to get the list of SecMaster dataclass types within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "dataclass_id" {}

data "huaweicloud_secmaster_dataclass_types" "test" {
  workspace_id = var.workspace_id
  dataclass_id = var.dataclass_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `dataclass_id` - (Required, String) Specifies the dataclass ID.

* `enabled` - (Optional, String) Specifies whether to filter enabled types.

* `order` - (Optional, String) Specifies the sort order.
  The valid values are as follows:
  + **ASC**: Ascending.
  + **DESC**: Descending.

* `sortby` - (Optional, String) Specifies the sort field.
  The valid values are as follows:
  + **CREATE_TIME**: Sort by creation time.
  + **CATEGORY**: Sort by category.

* `sub_category` - (Optional, String) Specifies the sub-category name used for filtering.

* `name` - (Optional, String) Specifies the type name used for filtering.

* `category_code` - (Optional, String) Specifies the category code used for filtering.

* `is_built_in` - (Optional, String) Specifies whether to filter built-in types.

* `layout_name` - (Optional, String) Specifies the layout name used for filtering.

* `level` - (Optional, Int) Specifies the type level used for filtering.
  The valid values are as follows:
  + **1**: Primary category.
  + **2**: Sub-category.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data` - The dataclass types list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The type ID.

* `dataclass_id` - The dataclass ID.

* `domain_id` - The account ID.

* `project_id` - The project ID.

* `workspace_id` - The workspace ID.

* `region_id` - The region ID.

* `layout_id` - The layout ID.

* `layout_name` - The layout name.

* `category` - The dataclass type category.

* `category_code` - The dataclass type category code.

* `sub_category` - The dataclass sub-category name.

* `sub_category_code` - The dataclass sub-category business code.

* `description` - The dataclass sub-category description.

* `enabled` - Whether the type is enabled.

* `level` - The type level.
  + **1**: Primary category.
  + **2**: Sub-category.

* `is_built_in` - Whether the type is built-in.

* `sla` - The response duration.

* `creator_id` - The creator ID.

* `creator_name` - The creator name.

* `modifier_id` - The modifier ID.

* `modifier_name` - The modifier name.

* `create_time` - The creation time.

* `update_time` - The update time.

* `dataclass_business_code` - The business code of the dataclass.

* `sub_count` - The number of sub-types under the type category.

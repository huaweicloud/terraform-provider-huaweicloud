---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_ai_component_statistics"
description: |-
  Use this data source to get the AI component statistics of HSS within HuaweiCloud.
---

# huaweicloud_hss_ai_component_statistics

Use this data source to get the AI component statistics of HSS within HuaweiCloud.

## Example Usage

```hcl
variable "category" {}
variable "catalogue" {}

data "huaweicloud_hss_ai_component_statistics" "test" {
  category  = var.category
  catalogue = var.catalogue
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `category` - (Required, String) Specifies the asset category.  
  The valid values are as follows:
  + **host**: Host asset.
  + **container**: Container asset.

* `catalogue` - (Required, String) Specifies the AI component category.  
  The valid values are as follows:
  + **app**: Application.
  + **tool**: Tool.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `ai_component_name` - (Optional, String) Specifies the AI component name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `data_list` - The AI component statistics list.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `ai_component_name` - The AI component name.

* `num` - The number of servers where the AI component is located.

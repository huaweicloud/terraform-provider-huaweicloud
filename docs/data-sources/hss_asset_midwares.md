---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_midwares"
description: |-
  Use this data source to get the list of middlewares.
---

# huaweicloud_hss_asset_midwares

Use this data source to get the list of middlewares.

## Example Usage

```hcl
data "huaweicloud_hss_asset_midwares" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Optional, String) Specifies the asset category.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `file_name` - (Optional, String) Specifies the jar package name.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The list of middlewares.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `file_name` - The jar package name.

* `num` - The number of middleware statistics.

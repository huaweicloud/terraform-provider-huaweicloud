---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_overview_status_os"
description: |-
  Use this data source to get the OS statistics of HSS asset overview status within HuaweiCloud.
---

# huaweicloud_hss_asset_overview_status_os

Use this data source to get the OS statistics of HSS asset overview status within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_overview_status_os" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `win_num` - The number of Windows assets.

* `linux_num` - The number of Linux assets.

* `os_list` - The list of OS statistics.

  The [os_list](#os_list_struct) structure is documented below.

<a name="os_list_struct"></a>
The `os_list` block supports:

* `os_name` - The OS name.

* `os_type` - The OS type.

* `number` - The number of assets with this OS.

---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_process_statistics"
description: |-
  Use this data source to get the list of process statistics.
---

# huaweicloud_hss_asset_process_statistics

Use this data source to get the list of process statistics.

## Example Usage

```hcl
data "huaweicloud_hss_asset_process_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `category` - (Optional, String) Specifies the type.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `path` - (Optional, String) Specifies the executable process path.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The process statistics list.
  The [data_list](#process_statistics_struct) structure is documented below.

<a name="process_statistics_struct"></a>
The `data_list` block supports:

* `path` - The executable process path.

* `num` - The number of processes.

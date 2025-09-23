---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_port_statistics"
description: |-
  Use this data source to get the list of HSS asset port statistics within HuaweiCloud.
---

# huaweicloud_hss_asset_port_statistics

Use this data source to get the list of HSS asset port statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_port_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `port` - (Optional, Int) Specifies the port number.

* `port_string` - (Optional, String) Specifies the port string used for fuzzy match.

* `type` - (Optional, String) Specifies the port type.

* `status` - (Optional, String) Specifies the port status.
  The valid values are as follows:
  + **danger**: Dangerous ports.
  + **unknow**: No ports with known risks.

* `sort_key` - (Optional, String) Specifies the sort key, sorting by port number is supported.

* `sort_dir` - (Optional, String) Specifies the sort direction. The default value is **asc**.
  The valid values are as follows:
  + **asc**: Ascending order.
  + **desc**: Descending order.

* `category` - (Optional, String) Specifies the type. The default value is **host**.
  The valid values are as follows:
  + **host**: Host.
  + **container**: Container.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project to which the resource belongs.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The port statistics list.
  The [data_list](#port_statistics_structure) structure is documented below.

<a name="port_statistics_structure"></a>
The `data_list` block supports:

* `port` - The port number.

* `type` - The port type.

* `num` - The number of ports.

* `status` - The port status.

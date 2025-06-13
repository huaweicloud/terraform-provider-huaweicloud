---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_app_statistics"
description: |-
  Use this data source to get the list of HSS app statistics within HuaweiCloud.
---

# huaweicloud_hss_app_statistics

Use this data source to get the list of HSS app statistics within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_app_statistics" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `app_name` - (Optional, String) Specifies the app name.

* `category` - (Optional, String) Specifies the app category. The default value is **host**.
  The valid values are as follows:
  + **host**: Host
  + **container**: Container

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project function is enabled.
  The value **all_granted_eps** indicates all enterprise projects.
  If omitted, the default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_list` - The app statistics list.
  The [data_list](#app_statistics_structure) structure is documented below.

<a name="app_statistics_structure"></a>
The `data_list` block supports:

* `app_name` - The app name.

* `num` - The number of hosts that have this app.

---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_setting_login_common_locations"
description: |-
  Use this data source to get the list of common login locations.
---

# huaweicloud_hss_setting_login_common_locations

Use this data source to get the list of common login locations.

## Example Usage

```hcl
data "huaweicloud_hss_setting_login_common_locations" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `area_code` - (Optional, Int) Specifies the code of countries and cities.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of configured common login locations.

* `data_list` - The list of common login locations information.
  The [data_list](#login_common_locations_data_list) structure is documented below.

<a name="login_common_locations_data_list"></a>
The `data_list` block supports:

* `area_code` - The code of countries and cities.

* `total_num` - The total number of hosts in the common login location.

* `host_id_list` - The list of host IDs.

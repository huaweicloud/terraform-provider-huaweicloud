---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_overview_account_top"
description: |-
  Use this data source to get the list of HSS asset overview account top within HuaweiCloud.
---

# huaweicloud_hss_asset_overview_account_top

Use this data source to get the list of HSS asset overview account top within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_overview_account_top" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number.

* `data_list` - The list of account top data.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `name` - The name.

* `host_num` - The number of hosts.

* `percentage` - The host occupancy percentage.

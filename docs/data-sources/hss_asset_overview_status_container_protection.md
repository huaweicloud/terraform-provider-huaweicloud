---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_overview_status_container_protection"
description: |-
  Use this data source to get the container protection statistics of HSS asset overview status within HuaweiCloud.
---

# huaweicloud_hss_asset_overview_status_container_protection

Use this data source to get the container protection statistics of HSS asset overview status within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_asset_overview_status_container_protection" "test" {}
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

* `no_risk` - The number of containers without risks.

* `risk` - The number of containers with risks.

* `no_protect` - The number of unprotected containers.

* `total_num` - The total number of containers.

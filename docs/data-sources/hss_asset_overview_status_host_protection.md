---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_overview_status_host_protection"
description: |-
  Use this data source to query the number of hosts based on the host protection status.
---

# huaweicloud_hss_asset_overview_status_host_protection

Use this data source to query the number of hosts based on the host protection status.

## Example Usage

```hcl
data "huaweicloud_hss_asset_overview_status_host_protection" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `no_risk` - The number of risk free hosts.

* `risk` - The number of risky hosts.

* `no_protect` - The number of unprotected hosts.

* `total_num` - The total number of hosts.

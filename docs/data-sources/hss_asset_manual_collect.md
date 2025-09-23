---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_asset_manual_collect"
description: |-
  Use this data source to query the asset fingerprint collection status of a host.
---

# huaweicloud_hss_asset_manual_collect

Use this data source to query the asset fingerprint collection status of a host.

## Example Usage

```hcl
variable "type" {}
variable "host_id" {}

data "huaweicloud_hss_asset_manual_collect" "test" {
  type    = var.type
  host_id = var.host_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the asset type.
  The valid values are as follows:
  + **web-app**
  + **web-service**
  + **web-framwork**
  + **web-site**
  + **midware**
  + **database**
  + **kernel-module**

* `host_id` - (Required, String) Specifies the host ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the hosts under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `scan_status` - The manual detection status.

* `scanned_time` - The detection completion time.

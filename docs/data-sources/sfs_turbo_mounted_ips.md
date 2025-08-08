---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_mounted_ips"
description: |-
  Use this data source to get the IP addresses of the clients who have mounted the file system.
---

# huaweicloud_sfs_turbo_mounted_ips

Use this data source to get the IP addresses of the clients who have mounted the file system.

## Example Usage

```hcl
variable "share_id" {}
variable "ips" {}

data "huaweicloud_sfs_turbo_mounted_ips" "test" {
  share_id = var.share_id
  ips      = var.ips
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS turbo file system.

* `ips` - (Optional, String) Specifies the IP addresses of the clients who have mounted the file system.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `ips_attribute` - The IP addresses of the clients who have mounted the file system.

---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_cdn_ips"
description: |-
  Use this data source to get the list of CDN IP addresses information within HuaweiCloud.
---

# huaweicloud_live_cdn_ips

Use this data source to get the list of CDN IP addresses information within HuaweiCloud.

## Example Usage

```hcl
variable "ip_address" {}

data "huaweicloud_live_cdn_ips" "test" {
  ip_addresses = [var.ip_address]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ip_addresses` - (Optional, List) Specifies the list of IP addresses, which can contain up to `20` IP addresses.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `cdn_ips` - The homing information of IP addresses.

  The [cdn_ips](#cdn_ips_struct) structure is documented below.

<a name="cdn_ips_struct"></a>
The `cdn_ips` block supports:

* `platform` - The platform name.

* `ip_address` - The IP address to be queried.

* `belongs` - Whether the IP address is a HuaweiCloud CDN node.

* `region` - The province where IP belongs.

* `isp` - The carrier name.

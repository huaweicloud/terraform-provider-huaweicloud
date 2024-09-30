---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_p2c_gateway_availability_zones"
description: |-
  Use this data source to get the list of VPN P2C gateway availability zones.
---

# huaweicloud_vpn_p2c_gateway_availability_zones

Use this data source to get the list of VPN P2C gateway availability zones.

## Example Usage

```hcl
data "huaweicloud_vpn_p2c_gateway_availability_zones" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `flavor` - (Optional, String) Specifies a flavor. The value can be **Professional1**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `availability_zones` - The list of availability zones.

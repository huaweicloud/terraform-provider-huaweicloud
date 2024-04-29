---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_availability_zones"
description: ""
---

# huaweicloud_vpn_gateway_availability_zones

Use this data source to get the list of VPN gateway availability zones.

## Example Usage

```hcl
variable "flavor" {}

data "huaweicloud_vpn_gateway_availability_zones" "test"{
  flavor = var.flavor
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `flavor` - (Required, String) Specifies the flavor name.
  The value can be **Basic**, **Professional1**, **Professional2** and **GM**.

* `attachment_type` - (Optional, String) Specifies the attachment type.
  The value can be **vpc** and **er**. Defaults to **vpc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `names` - The names of the availability zones.

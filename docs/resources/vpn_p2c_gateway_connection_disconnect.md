---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_p2c_gateway_connection_disconnect"
description: |-
  Manages a VPN P2C gateway connection disconnect resource within HuaweiCloud.
---

# huaweicloud_vpn_p2c_gateway_connection_disconnect

Manages a VPN P2C gateway connection disconnect resource within HuaweiCloud.

## Example Usage

```hcl
variable "p2c_vgw_id" {}
variable "connection_id" {}

resource "huaweicloud_vpn_p2c_gateway_connection_disconnect" "test" {
  p2c_vgw_id    = var.p2c_vgw_id
  connection_id = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `p2c_vgw_id` - (Required, String, NonUpdatable) Specifies the instance ID of a P2C VPN gateway.

* `connection_id` - (Required, String, NonUpdatable) Specifies the connection ID of a P2C VPN gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

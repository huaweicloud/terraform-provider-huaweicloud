---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_upgrade"
description: |-
  Manages a VPN gateway upgrade resource within HuaweiCloud.
---

# huaweicloud_vpn_gateway_upgrade

Manages a VPN gateway upgrade resource within HuaweiCloud.

## Example Usage

```hcl
variable "vgw_id" {}

resource "huaweicloud_vpn_gateway_upgrade" "test" {
  vgw_id = var.vgw_id
  action = "start"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `vgw_id` - (Required, String, NonUpdatable) Specifies the instance ID of a VPN gateway.

* `action` - (Required, String, NonUpdatable) Specifies an upgrade operation.
  Value can be **start**, **finish** or **rollback**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_connection_reset"
description: |-
  Manages a VPN gateway connection reset resource within HuaweiCloud.
---

# huaweicloud_vpn_connection_reset

Manages a VPN gateway connection reset resource within HuaweiCloud.

## Example Usage

```hcl
variable "connection_id" {}

resource "huaweicloud_vpn_connection_reset" "test" {
  connection_id = var.connection_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `connection_id` - (Required, String, NonUpdatable) Specifies the connection ID of a VPN gateway.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

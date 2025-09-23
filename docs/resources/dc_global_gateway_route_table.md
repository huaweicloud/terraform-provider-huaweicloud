---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateway_route_table"
description: |-
  Manages a DC global gateway route table resource within HuaweiCloud.
---

# huaweicloud_dc_global_gateway_route_table

Manages a DC global gateway route table resource within HuaweiCloud.

## Example Usage

```hcl
variable "gdgw_id" {}
variable "nexthop" {}

resource "huaweicloud_dc_global_gateway_route_table" "test" {
  gdgw_id     = var.gdgw_id
  type        = "vif_peer"
  destination = "2.2.2.0/30"
  nexthop     = var.nexthop
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create new resource.

* `gdgw_id` - (Required, String, NonUpdatable) Specifies the global DC gateway ID.

* `type` - (Required, String, NonUpdatable) Specifies the next hop type. Value options:
  + **vif_peer**: virtual interface peer
  + **gdgw**: global DC gateway

* `destination` - (Required, String, NonUpdatable) Specifies the subnet the route is destined for.

* `nexthop` - (Required, String, NonUpdatable) Specifies the next hop ID.

* `description` - (Optional, String) Specifies the description of the route.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `obtain_mode` - Indicates the route type.

* `status` - Indicates the route status.

* `address_family` - Indicates the address family.

## Import

The DC connect gateway resource can be imported using the `gdgw_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dc_global_gateway_route_table.test <gdgw_id>/<id>
```

---
subcategory: "Direct Connect (DC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dc_global_gateway_route_tables"
description: |-
  Use this data source to get a list of DC global gateway route tables.
---

# huaweicloud_dc_global_gateway_route_tables

Use this data source to get a list of DC global gateway route tables.

## Example Usage

```hcl
variable "gdgw_id" {}

data "huaweicloud_dc_global_gateway_route_tables" "test" {
  gdgw_id = var.gdgw_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `gdgw_id` - (Required, String) Specifies the global gateway ID.

* `nexthop` - (Optional, List) Specifies the nexthop IDs to filter the routes.

* `destination` - (Optional, List) Specifies the destination addresses to filter the routes.

* `address_family` - (Optional, List) Specifies the address families to filter the routes.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gdgw_routetables` - Indicates the list of global gateway route tables.
  The [gdgw_routetables](#gdgw_routetables_struct) structure is documented below.

<a name="gdgw_routetables_struct"></a>
The `gdgw_routetables` block supports:

* `id` - Indicates the route ID.

* `gateway_id` - Indicates the gateway ID.

* `nexthop` - Indicates the nexthop ID.

* `obtain_mode` - Indicates the route type. The value can be:
  + **customized**: default route.
  + **specific**: custom route.
  + **bgp**: dynamic route.

* `description` - Indicates the route description.

* `destination` - Indicates the route subnet.

* `status` - Indicates the route status. The value can be:
  + **ACTIVE**: issued normally.
  + **ERROR**: failed to issue.
  + **PENDING_CREATE**: to be issued.

* `address_family` - Indicates the address family.

* `type` - Indicates the nexthop type. The value can be:
  + **vif_peer**: virtual interface peer.
  + **gdgw**: global gateway.

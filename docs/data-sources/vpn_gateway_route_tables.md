---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_gateway_route_tables"
description: |-
  Use this data source to get the list of VPN gateway route tables.
---

# huaweicloud_vpn_gateway_route_tables

Use this data source to get the list of VPN gateway route tables.

## Example Usage

```hcl
variable "vgw_id" {}

data "huaweicloud_vpn_gateway_route_tables" "test" {
  vgw_id = var.vgw_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vgw_id` - (Required, String) Specifies a VPN gateway ID.

* `is_include_nexthop_resource` - (Optional, Bool) Specifies whether to include the nexthop resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `routing_table` - Indicates the route table of a specified VPN gateway.

  The [routing_table](#routing_table_struct) structure is documented below.

<a name="routing_table_struct"></a>
The `routing_table` block supports:

* `destination` - Indicates the destination address of a route.

* `nexthop` - Indicates the next-hop IP address.

* `outbound_interface_ip` - Indicates the IP address of the outbound interface.

* `origin` - Indicates the origin of a BGP route.

* `as_path` - Indicates the AS path of a BGP route.

* `med` - Indicates the MED value of a BGP route.

* `nexthop_resource` - Indicates the next hop resource of a route.

  The [nexthop_resource](#routing_table_nexthop_resource_struct) structure is documented below.

<a name="routing_table_nexthop_resource_struct"></a>
The `nexthop_resource` block supports:

* `id` - Indicates the next-hop resource ID, which is in UUID format.

* `type` - Indicates the next-hop resource type.

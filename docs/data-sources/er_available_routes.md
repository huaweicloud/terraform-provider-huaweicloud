---
subcategory: "Enterprise Router (ER)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_er_available_routes"
description: |-
  Using this data source to query the list of available routes within HuaweiCloud.
---

# huaweicloud_er_available_routes

Using this data source to query the list of available routes within HuaweiCloud.

## Example Usage

```hcl
variable "route_table_id" {}

data "huaweicloud_er_available_routes" "test" {
  route_table_id = var.route_table_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the available routes.
  If omitted, the provider-level region will be used.

* `route_table_id` - (Required, String) The route table ID to which the available routes belong.

* `destination` - (Optional, String) The destination address of the routes to be queried.

* `resource_type` - (Optional, String) The attachment type.
  The valid values are as follows:
  + **vpc**: VPC attachment.
  + **vpn**: VPN gateway attachment.
  + **vgw**: virtual gateway attachment.
  + **peering**: peering connection attachment.
  + **ecn**: ECN attachment.
  + **cfw**: CFW instance attachment.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `routes` - All available routes that match the filter parameters.

  The [routes](#routes_struct) structure is documented below.

<a name="routes_struct"></a>
The `routes` block supports:

* `id` - The route ID.

* `destination` - The destination address of the route.

* `next_hops` - The next hops of the route.

  The [next_hops](#routes_next_hops_struct) structure is documented below.

* `is_blackhole` - Whether the route is a blackhole route.

* `type` - The route type.

<a name="routes_next_hops_struct"></a>
The `next_hops` block supports:

* `resource_id` - The attached resource ID.

* `resource_type` - The attachment type.

* `attachment_id` - The attachment ID.

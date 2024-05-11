---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_routes"
description: |-
  Use this data source to get the list of VPC routes.
---

# huaweicloud_vpc_routes

Use this data source to get the list of VPC routes.

## Example Usage

```hcl
data "huaweicloud_vpc_routes" "test" {
  type = "peering"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the route type.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC to which the route belongs.

* `destination` - (Optional, String) Specifies the route destination.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `routes` - The list of routes.

  The [routes](#routes_struct) structure is documented below.

<a name="routes_struct"></a>
The `routes` block supports:

* `id` - The route ID.

* `type` - The route type.

* `vpc_id` - The ID of the VPC to which the route belongs.

* `destination` - The route destination.

* `nexthop` - The next hop of the route.

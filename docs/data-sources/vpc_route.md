---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_route"
description: ""
---

# huaweicloud_vpc_route

!> **WARNING:** It has been deprecated, use `huaweicloud_vpc_route_table` to get the route details.

Provides details about a specific VPC route.

## Example Usage

```hcl
data "huaweicloud_vpc_route" "vpc_route" {
  vpc_id = var.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available routes in the current tenant. The given
filters must match exactly one route whose data will be exported as attributes.

* `region` - (Optional, String) The region in which to obtain the vpc route. If omitted, the provider-level region will
  be used.

* `id` - (Optional, String) The id of the specific route to retrieve.

* `vpc_id` - (Optional, String) The id of the VPC that the desired route belongs to.

* `destination` - (Optional, String) The route destination address (CIDR).

* `type` - (Optional, String) Route type for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `nexthop` - The next hop of the route. If the route type is peering, it will provide VPC peering connection ID.

---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_route"
sidebar_current: "docs-huaweicloud-datasource-vpc-route"
description: |-
  Provides details about a specific VPC Route.
---

# huaweicloud\_vpc\_route

Provides details about a specific VPC route.
This is an alternative to `huaweicloud_vpc_route_v2`

## Example Usage

```hcl
data "huaweicloud_vpc_route" "vpc_route" {
  vpc_id = var.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
routes in the current tenant. The given filters must match exactly one
route whose data will be exported as attributes.

* `id` (Optional) - The id of the specific route to retrieve.

* `vpc_id` (Optional) - The id of the VPC that the desired route belongs to.

* `destination` (Optional) - The route destination address (CIDR).

* `tenant_id` (Optional) - Only the administrator can specify the tenant ID of other tenants.

* `type` (Optional) - Route type for filtering.

## Attribute Reference

All of the argument attributes are also exported as
result attributes.

* `nexthop` - The next hop of the route. If the route type is peering, it will provide VPC peering connection ID.

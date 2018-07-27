---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_route_v2"
sidebar_current: "docs-huaweicloud-resource-vpc-route-v2"
description: |-
  Provides an VPC route resource.
---

# huaweicloud_vpc_route_v2

Provides a resource to create a route.

## Example Usage

 ```hcl
resource "huaweicloud_vpc_route_v2" "vpc_route" {
  type  = "peering"
  nexthop  = "${var.nexthop}"
  destination = "192.168.0.0/16"
  vpc_id = "${var.vpc_id}"
 }
```

## Argument Reference

The following arguments are supported:

* `destination` (Required) - Specifies the destination IP address or CIDR block. Changing this creates a new Route.

* `nexthop` (Required) - Specifies the next hop. If the route type is peering, enter the VPC peering connection ID. Changing this creates a new Route.

* `type` (Required) - Specifies the route type. Currently, the value can only be **peering**. Changing this creates a new Route.

* `vpc_id` (Required) - Specifies the VPC for which a route is to be added. Changing this creates a new Route.

* `tenant_id` (Optional) - Specifies the tenant ID. Only the administrator can specify the tenant ID of other tenant. Changing this creates a new Route.

## Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - The route ID.

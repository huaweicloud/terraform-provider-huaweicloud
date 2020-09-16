---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_route_ids"
sidebar_current: "docs-huaweicloud-datasource-vpc-route-ids"
description: |-
  Provides a list of route Ids for a VPC
---

# Data Source: huaweicloud_vpc_route_ids

Provides a list of route ids for a vpc_id.
This is an alternative to `huaweicloud_vpc_route_ids_v2`

This resource can be useful for getting back a list of route ids for a vpc.

## Example Usage

 ```hcl
variable "vpc_id" {}

data "huaweicloud_vpc_route_ids" "example" {
  vpc_id = var.vpc_id
}

data "huaweicloud_vpc_route" "vpc_route" {
  count = length(data.huaweicloud_vpc_route_ids.example.ids)
  id    = data.huaweicloud_vpc_route_ids.example.ids[count.index]
}

output "route_nexthop" {
  value = ["${data.huaweicloud_vpc_route.vpc_route.*.nexthop}"]
}
 ```

## Argument Reference

* `vpc_id` (Required) - The VPC ID that you want to filter from.

## Attributes Reference

* `ids` - A list of all the route ids found. This data source will fail if none are found.


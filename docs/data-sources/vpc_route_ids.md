---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_route_ids"
description: ""
---

# huaweicloud_vpc_route_ids

!> **WARNING:** It has been deprecated, use `huaweicloud_vpc_route_table` to get the route details.

Provides a list of route ids for a vpc_id.

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

* `region` - (Optional, String) The region in which to obtain the ids. If omitted, the provider-level region will be
  used.

* `vpc_id` - (Required, String) The VPC ID that you want to filter from.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a data source ID in UUID format.

* `ids` - A list of all the route ids found. This data source will fail if none are found.

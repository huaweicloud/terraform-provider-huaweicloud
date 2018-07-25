---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_peering_connection_v2"
sidebar_current: "docs-huaweicloud-datasource-vpc-peering-v2"
description: |-
  Provides details about a specific VPC peering connection.
---

# Data Source: huaweicloud_vpc_peering_connection_v2

The VPC Peering Connection data source provides details about a specific VPC peering connection.


## Example Usage

 ```hcl

data "huaweicloud_vpc_peering_connection_v2" "peering" {
   vpc_id          = "${huaweicloud_vpc_v1.vpc.id}"
   peer_vpc_id     = "${huaweicloud_vpc_v1.peer_vpc.id}"
 }


resource "huaweicloud_vpc_route_v2" "vpc_route" {
  type       = "peering"
  nexthop    = "${data.huaweicloud_vpc_peering_connection_v2.peering.id}"
  destination = "192.168.0.0/16"
  vpc_id = "${huaweicloud_vpc_v1.vpc.id}"
}
 ```


## Argument Reference

The arguments of this data source act as filters for querying the available VPC peering connection.
The given filters must match exactly one VPC peering connection whose data will be exported as attributes.

* `id` (Optional) - The ID of the specific VPC Peering Connection to retrieve.

* `status` (Optional) - The status of the specific VPC Peering Connection to retrieve.

* `vpc_id` (Optional) - The ID of the requester VPC of the specific VPC Peering Connection to retrieve.

* `peer_vpc_id` (Optional) -  The ID of the accepter/peer VPC of the specific VPC Peering Connection to retrieve.

* `peer_tenant_id` (Optional) - The Tenant ID of the accepter/peer VPC of the specific VPC Peering Connection to retrieve.

* `name` (Optional) - The name of the specific VPC Peering Connection to retrieve.


## Attributes Reference

All of the argument attributes are exported as result attributes.
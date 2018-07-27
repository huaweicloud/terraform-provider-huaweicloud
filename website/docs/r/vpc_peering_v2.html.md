---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_peering_connection_v2"
sidebar_current: "docs-huaweicloud-resource-vpc-peering-v2"
description: |-
  Manage a VPC Peering Connection resource.
---

# huaweicloud_vpc_peering_connection_v2

Provides a resource to manage a VPC Peering Connection resource.

-> **Note:** For cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connections, use the `huaweicloud_vpc_peering_connection_v2` resource to manage the requester's side of the connection and use the `huaweicloud_vpc_peering_connection_accepter_v2` resource to manage the accepter's side of the connection.

## Example Usage

 ```hcl
resource "huaweicloud_vpc_peering_connection_v2" "peering" {
  name = "${var.peer_conn_name}"
  vpc_id = "${var.vpc_id}"
  peer_vpc_id = "${var.accepter_vpc_id}"
}
 ```

## Argument Reference

The following arguments are supported:

* `name` (Required) - Specifies the name of the VPC peering connection. The value can contain 1 to 64 characters.

* `vpc_id` (Required) - Specifies the ID of a VPC involved in a VPC peering connection. Changing this creates a new VPC peering connection.

* `peer_vpc_id` (Required) - Specifies the VPC ID of the accepter tenant. Changing this creates a new VPC peering connection.

* `peer_tenant_id` (Optional) - Specified the Tenant Id of the accepter tenant. Changing this creates a new VPC peering connection.
  
## Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - The VPC peering connection ID.

* `status` - The VPC peering connection status. The value can be PENDING_ACCEPTANCE, REJECTED, EXPIRED, DELETED, or ACTIVE.

## Notes

If you create a VPC peering connection with another VPC of your own, the connection is created without the need for you to accept the connection.

## Import

VPC Peering resources can be imported using the `vpc peering id`, e.g.

> $ terraform import huaweicloud_vpc_peering_connection_v2.test_connection 22b76469-08e3-4937-8c1d-7aad34892be1
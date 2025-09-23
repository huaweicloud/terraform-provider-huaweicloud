---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_peering_connection"
description: ""
---

# huaweicloud_vpc_peering_connection

Provides a resource to manage a VPC Peering Connection resource.

-> **NOTE:** For cross-tenant (requester's tenant differs from the accepter's tenant) VPC Peering Connections,
  use the `huaweicloud_vpc_peering_connection` resource to manage the requester's side of the connection and
  use the `huaweicloud_vpc_peering_connection_accepter` resource to manage the accepter's side of the connection.
  <br/>If you create a VPC peering connection with another VPC of your own, the connection is created without the need
  for you to accept the connection.

## Example Usage

 ```hcl
resource "huaweicloud_vpc_peering_connection" "peering" {
  name        = var.peer_conn_name
  vpc_id      = var.vpc_id
  peer_vpc_id = var.accepter_vpc_id
}
 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the VPC peering connection. If omitted, the
  provider-level region will be used. Changing this creates a new VPC peering connection resource.

* `name` - (Required, String) Specifies the name of the VPC peering connection. The value can contain 1 to 64
  characters.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of a VPC involved in a VPC peering connection. Changing this
  creates a new VPC peering connection.

* `peer_vpc_id` - (Required, String, ForceNew) Specifies the VPC ID of the accepter tenant. Changing this creates a new
  VPC peering connection.

* `peer_tenant_id` - (Optional, String, ForceNew) Specifies the tenant ID of the accepter tenant. Changing this creates
  a new VPC peering connection.

* `description` - (Optional, String) Specifies the description of the VPC peering connection.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The VPC peering connection ID.

* `status` - The VPC peering connection status. The value can be PENDING_ACCEPTANCE, REJECTED, EXPIRED, DELETED, or
  ACTIVE.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

VPC Peering resources can be imported using the `vpc peering id`, e.g.

```bash
$ terraform import huaweicloud_vpc_peering_connection.test_connection 22b76469-08e3-4937-8c1d-7aad34892be1
```

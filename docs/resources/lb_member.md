---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_member"
description: ""
---

# huaweicloud_lb_member

Manages an ELB member resource within HuaweiCloud.

## Example Usage

```hcl
variable "lb_pool_id" {}
variable "ipv4_subnet_id" {}

resource "huaweicloud_lb_member" "member_1" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = var.lb_pool_id
  subnet_id     = var.ipv4_subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB member resource. If omitted, the the
  provider-level region will be used. Changing this creates a new member.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this member will be assigned to.

* `subnet_id` - (Required, String, ForceNew) The **IPv4 subnet ID** of the subnet in which to access the member.

* `name` - (Optional, String) Human-readable name for the member.

* `address` - (Required, String, ForceNew) The IP address of the member to receive traffic from the load balancer.
  Changing this creates a new member.

* `protocol_port` - (Required, Int, ForceNew) The port on which to listen for client traffic. Changing this creates a
  new member.

* `weight` - (Optional, Int)  A positive integer value that indicates the relative portion of traffic that this member
  should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a
  member with a weight of 2.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the member.

* `backend_server_status` - Indicates the administrative status of the backend server.

* `operating_status` - Indicates the health check result of the backend server.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB member can be imported using the pool ID and member ID separated by a slash, e.g.

```bash
$ terraform import huaweicloud_lb_member.member_1 e0bd694a-abbe-450e-b329-0931fd1cc5eb/4086b0c9-b18c-4d1c-b6b8-4c56c3ad2a9e
```

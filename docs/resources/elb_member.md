---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_member"
description: |-
  Manages an ELB member resource within HuaweiCloud.
---

# huaweicloud_elb_member

Manages an ELB member resource within HuaweiCloud.

## Example Usage

```hcl
variable "elb_pool_id" {}
variable "ipv4_subnet_id" {}

resource "huaweicloud_elb_member" "member_1" {
  address       = "192.168.199.23"
  protocol_port = 8080
  pool_id       = var.elb_pool_id
  subnet_id     = var.ipv4_subnet_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB member resource. If omitted, the the
  provider-level region will be used. Changing this creates a new member.

* `pool_id` - (Required, String, ForceNew) The id of the pool that this member will be assigned to.

* `subnet_id` - (Optional, String, ForceNew) The **IPv4 or IPv6 subnet ID** of the subnet in which to access the member.
  + The IPv4 or IPv6 subnet must be in the same VPC as the subnet of the load balancer.
  + This parameter must be specified for gateway load balancers. The subnet of the backend server must be in the same
    VPC as that of the load balancer, and it must be different from the subnet of the load balancer.
  + If this parameter is not specified, **cross-VPC backend** has been enabled for the load balancer.
    In this case, cross-VPC backend servers must use private IPv4 addresses,
    and the protocol of the backend server group must be TCP, HTTP, or HTTPS.

* `name` - (Optional, String) Human-readable name for the member.

* `address` - (Required, String, ForceNew) The IP address of the member to receive traffic from the load balancer.
  Changing this creates a new member.

* `protocol_port` - (Optional, Int) The port on which to listen for client traffic. It must be set to `0` for gateway
  load balancers with IP backend server groups associated. It can be left blank because it does not take effect if
  `any_port_enable` is set to **true** for a backend server group.

* `weight` - (Optional, Int)  A positive integer value that indicates the relative portion of traffic that this member
  should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a
  member with a weight of 2.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the member.

* `instance_id` - The ID of the instance associated with the backend server. If this parameter is left blank, the backend
  server is not a real device. It may be an IP address.

* `ip_version` - The IP version supported by the backend server. The value can be **v4 (IPv4)** or **v6 (IPv6)**, depending
  on the value of address returned by the system.

* `member_type` - The type of the backend server. The value can be:
  + **ip**: IP as backend servers
  + **instance**: ECSs used as backend servers

* `operating_status` - The health status of the backend server if `listener_id` under `status` is not specified. The value
  can be:
  + **ONLINE**: The backend server is running normally.
  + **NO_MONITOR**: No health check is configured for the backend server group to which the backend server belongs.
  + **OFFLINE**: The cloud server used as the backend server is stopped or does not exist.

* `reason` - Why health check fails.
  The [reason](#reason_struct) structure is documented below.

* `status` - The health status of the backend server if `listener_id` under status is specified. If `listener_id` under
  `status` is not specified, `operating_status` of member takes precedence.
  The [status](#status_struct) structure is documented below.

* `created_at` - The time when the backend server was added. The format is **yyyy-MM-dd'T'HH:mm:ss'Z' (UTC time)**.

* `updated_at` - The time when the backend server was updated. The format is **yyyy-MM-dd'T'HH:mm:ss'Z' (UTC time)**.

<a name="reason_struct"></a>
The `reason` block supports:

* `expected_response` - The code of the health check failures.

* `healthcheck_response` - The expected HTTP status code.

* `reason_code` - The returned HTTP status code in the response.

<a name="status_struct"></a>
The `status` block supports:

* `listener_id` - The listener ID.

* `operating_status` - The health status of the backend server.

* `reason` - Why health check fails.
  The [reason](#reason_struct) structure is documented below.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

ELB member can be imported using the `pool_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_elb_member.member_1 <pool_id>/<id>
```

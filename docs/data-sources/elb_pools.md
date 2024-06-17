---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_pools"
description: ""
---

# huaweicloud_elb_pools

Use this data source to get the list of ELB pools.

## Example Usage

```hcl
variable "pool_name" {}

data "huaweicloud_elb_pools" "test" {
  name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the ELB pool.

* `pool_id` - (Optional, String) Specifies the ID of the ELB pool.

* `description` - (Optional, String) Specifies the description of the ELB pool.

* `loadbalancer_id` - (Optional, String) Specifies the loadbalancer ID of the ELB pool.

* `healthmonitor_id` - (Optional, String) Specifies the health monitor ID of the ELB pool.

* `protocol` - (Optional, String) Specifies the protocol of the ELB pool. Value options: **TCP**, **UDP**, **HTTP**,
  **HTTPS**, **QUIC**, **GRPC** or **TLS**.

* `lb_method` - (Optional, String) Specifies the method of the ELB pool. Value options: **ROUND_ROBIN**,
  **LEAST_CONNECTIONS**, **SOURCE_IP** or **QUIC_CID**.

* `listener_id` - (Optional, String) Specifies the listener ID of the ELB pool.

* `type` - (Optional, String) Specifies the type of the backend server group. Value options: **instance**, **ip**.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the backend server group works.

* `protection_status` - (Optional, String) Specifies the protection status for update.
  Value options: **nonProtection**, **consoleProtection**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - Pool list. For details, see data structure of the pool field.
  The [object](#pools_object) structure is documented below.

<a name="pools_object"></a>
The `pools` block supports:

* `id` - The pool ID.

* `name` - The pool name.

* `description` - The description of pool.

* `protocol` - The protocol of pool.

* `lb_method` - The load balancing algorithm to distribute traffic to the pool's members.

* `healthmonitor_id` - The health monitor ID of the LB pool.

* `ip_version` - The IP version of the LB pool.

* `type` - The type of the backend server group.

* `vpc_id` - The ID of the VPC where the backend server group works.

* `protection_status` - The protection status for update.

* `protection_reason` - The reason for update protection.

* `slow_start_enabled` - Whether to enable slow start.

* `slow_start_duration` - The slow start duration, in seconds.

* `connection_drain_enabled` - Whether to enable delayed logout.

* `connection_drain_timeout` - The timeout of the delayed logout in seconds.

* `minimum_healthy_member_count` - The minimum healthy member count.

* `listeners` - The listener list. The [object](#elem_object) structure is documented below.

* `loadbalancers` - The loadbalancer list. The [object](#elem_object) structure is documented below.

* `members` - The member list. The [object](#elem_object) structure is documented below.

* `persistence` - Indicates whether connections in the same session will be processed by the same pool member or not.
  The [object](#persistence_object) structure is documented below.

<a name="elem_object"></a>
The `listeners`,  `loadbalancers` or `members` block supports:

* `id` - The listener, loadbalancer or member ID.

<a name="persistence_object"></a>
The `persistence` block supports:

* `type` - The type of persistence mode.

* `cookie_name` - The name of the cookie if persistence mode is set appropriately.

* `timeout` - The stickiness duration, in minutes.

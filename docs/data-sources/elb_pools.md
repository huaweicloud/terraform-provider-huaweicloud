---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_pools"
description: |-
  Use this data source to get the list of ELB pools.
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

* `any_port_enable` - (Optional, String) Specifies whether forward to same port for a backend server group is enabled.
  Value options:
  + **false**: Disable this option.
  + **true**: Enable this option.

* `connection_drain` - (Optional, String) Specifies whether delayed logout is enabled. Value options:
  + **false**: Disable this option.
  + **true**: Enable this option.

* `enterprise_project_id` - (Optional, String) Specifies the ID of the enterprise project.

* `ip_version` - (Optional, String) Specifies the IP address version supported by the backend server group.

* `member_address` - (Optional, String) Specifies the private IP address bound to the backend server.

* `member_device_id` - (Optional, String) Specifies the ID of the cloud server that serves as a backend server.

* `member_instance_id` - (Optional, String) Specifies the backend server ID.

* `member_deletion_protection_enable` - (Optional, String) Specifies whether deletion protection is enabled. Value options:
  + **false**: Disable this option.
  + **true**: Enable this option.

* `pool_health` - (Optional, String) Specifies whether pool health is enabled. Value options:
  + **minimum_healthy_member_count=0**
  + **minimum_healthy_member_count=1**

* `public_border_group` - (Optional, String) Specifies the public border group.

* `quic_cid_len` - (Optional, Int) Specifies the QUIC connection ID len.

* `quic_cid_offset` - (Optional, Int) Specifies the QUIC connection ID offset.

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

* `any_port_enable` - Whether forward to same port for a backend server group is enabled

* `enterprise_project_id` - The ID of the enterprise project.

* `member_deletion_protection_enable` - Whether deletion protection is enabled

* `public_border_group` - The public border group.

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

* `quic_cid_hash_strategy` - The multi-path forwarding policy based on destination connection IDs.
  The [quic_cid_hash_strategy](#quic_cid_hash_strategy_struct) structure is documented below.

* `created_at` - The time when the backend server group was created

* `updated_at` - The time when the backend server group was updated.

<a name="elem_object"></a>
The `listeners`,  `loadbalancers` or `members` block supports:

* `id` - The listener, loadbalancer or member ID.

<a name="persistence_object"></a>
The `persistence` block supports:

* `type` - The type of persistence mode.

* `cookie_name` - The name of the cookie if persistence mode is set appropriately.

* `timeout` - The stickiness duration, in minutes.

<a name="quic_cid_hash_strategy_struct"></a>
The `quic_cid_hash_strategy` block supports:

* `len` - The length of the hash factor in the connection ID, in byte.

* `offset` - The start position in the connection ID as the hash factor, in byte.

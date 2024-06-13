---
subcategory: Dedicated Load Balance (Dedicated ELB)
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_active_standby_pools"
description: ""
---

# huaweicloud_elb_active_standby_pools

Use this data source to get the list of active standby ELB pools.

## Example Usage

```hcl
variable "pool_name" {}

data "huaweicloud_elb_active_standby_pools" "test" {
  name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the active-standby pool.

* `pool_id` - (Optional, String) Specifies the ID of the active-standby pool.

* `description` - (Optional, String) Specifies supplementary information about the active-standby pool.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the load balancer with which the active-standby pool is
  associated.

* `healthmonitor_id` - (Optional, String) Specifies the ID of the health check configured for the active-standby pool.

* `protocol` - (Optional, String) Specifies the protocol used by the active-standby pool to receive requests from the
  load balancer. Value options: **TCP**, **UDP**, **QUIC** or **TLS**.

* `member_address` - (Optional, String) Specifies the private IP address bound to the member. This parameter is used
  only as a query condition and is not included in the response.

* `member_instance_id` - (Optional, String) Specifies the ID of the ECS used as the member. This parameter is used only
  as a query condition and is not included in the response.

* `listener_id` - (Optional, String) Specifies the ID of the listener to which the forwarding policy is added.

* `type` - (Optional, String) Specifies the type of the active-standby pool.
  The valid values are as follows:
  + **instance**: Any type of backend servers can be added.
  + **ip**: Only IP as backend servers can be added.

  If the `type` is empty, any type of backend servers can be added.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the active-standby pool works.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - The pool list. For details, see data structure of the pool field.
  The [pools](#elb_pools) structure is documented below.

<a name="elb_pools"></a>
The `pools` block supports:

* `id` - The ID of the active-standby pool.

* `name` - The name of the active-standby pool.

* `description` - The description of the active-standby pool.

* `protocol` - The protocol used by the active-standby pool to receive requests.

* `type` - The type of the active-standby pool.

* `any_port_enable` - Whether to enable Forward to same Port for a pool.

* `vpc_id` - The ID of the VPC where the active-standby pool works.

* `connection_drain_enabled` - Whether to enable delayed logout.

* `connection_drain_timeout` - The timeout of the delayed logout in seconds.

* `listeners` - The IDs of the listeners with which the active-standby pool is associated.
  The [listeners](#elb_listeners) structure is documented below.

* `loadbalancers` - The IDs of the load balancers with which the active-standby pool is associated.
  The [loadbalancers](#elb_loadbalancers) structure is documented below.

* `members` - The backend servers in the active-standby pool.
  The [members](#elb_members) structure is documented below.

* `healthmonitor` - The health check configured for the active-standby pool.
  The [healthmonitor](#elb_healthmonitor) structure is documented below.

<a name="elb_listeners"></a>
The `listeners` block supports:

* `id` - The listener ID.

<a name="elb_loadbalancers"></a>
The `loadbalancers` block supports:

* `id` - The loadbalancer ID.

<a name="elb_members"></a>
The `members` block supports:

* `id` - The ID of the member.

* `name` - The name of the member.

* `subnet_id` - The ID of the IPv4 or IPv6 subnet where the member resides.

* `protocol_port` - The port used by the member to receive requests.

* `address` - The private IP address bound to the member.

* `ip_version` - The IP version supported by the member.

* `operating_status` - The health status of the member.

* `member_type` - The type of the member.

* `instance_id` - The ID of the ECS used as the member.

* `role` - The active-standby status of the member.

<a name="elb_healthmonitor"></a>
The `healthmonitor` block supports:

* `id` - The health check ID.

* `name` - The health check name.

* `delay` - The interval between health checks, in seconds.

* `domain_name` - The domain name that HTTP requests are sent to during the health check.

* `expected_codes` - The expected HTTP status code.

* `http_method` - The HTTP method.

* `max_retries_down` - The number of consecutive health checks when the health check result of a backend server changes
  from **ONLINE** to **OFFLINE**.

* `max_retries` - The number of consecutive health checks when the health check result of a backend server changes from
  **OFFLINE** to **ONLINE**.

* `monitor_port` - The port used for the health check.

* `timeout` - The maximum time required for waiting for a response from the health check, in seconds.

* `type` - The health check protocol.

* `url_path` - The HTTP request path for the health check.

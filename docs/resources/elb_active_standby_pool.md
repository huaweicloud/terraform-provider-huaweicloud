---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_active_standby_pool"
description: ""
---

# huaweicloud_elb_active_standby_pool

Manages an ELB active-standby pool resource within HuaweiCloud.

## Example Usage

### Create an active-standby Pool

```hcl
variable "vpc_id" {}

resource "huaweicloud_elb_active_standby_pool" "test" {
  name            = "test_active_standby_pool"
  description     = "test description"
  protocol        = "TCP"
  vpc_id          = var.vpc_id
  type            = "instance"
  any_port_enable = false

  members {
    address       = "192.168.0.1"
    role          = "master"
    protocol_port = 45
  }

  members {
    address       = "192.168.0.2"
    role          = "slave"
    protocol_port = 36
  }

  healthmonitor {
    delay            = 5
    expected_codes   = "200"
    max_retries      = 3
    max_retries_down = 3
    timeout          = 3
    type             = "TCP"
  }
}

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ELB active-standby pool resource.
  If omitted, the provider-level region will be used.

* `protocol` - (Required, String, ForceNew) Specifies the protocol used by the active-standby pool to receive requests.
  Value options: **TCP**, **UDP**, **QUIC** or **TLS**.
  + If the listener's protocol is **UDP**, the value must be **UDP** or **QUIC**.
  + If the listener's protocol is **TCP**, the value must be **TCP**.
  + If the listener's protocol is **TLS**, the value must be **TLS** or **TCP**.

  Changing this parameter will create a new resource.

* `members` - (Required, List, ForceNew) Specifies the members in the active-standby pool.
  The [members](#ELB_members) structure is documented below. Changing this parameter will create a new resource.

* `healthmonitor` - (Required, List, ForceNew) Specifies the health check configured for the active-standby pool.
  The [healthmonitor](#ELB_healthmonitor) structure is documented below. Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the type of the active-standby pool. Value options:
  + **instance**: Any type of backend servers can be added. `vpc_id` must be mandatory.
  + **ip**: Only IP as backend servers can be added. `vpc_id` cannot be specified.

  Changing this parameter will create a new resource.

  -> **NOTE:** If this parameter is not passed, any type of active-standby pool can be added and will return an empty string.
  This parameter can be updated only when it is left blank.

* `ip_version` - (Optional, String, ForceNew) Specifies the IP address version supported by active-standby pool.
  The value can be **dualstack**, **v6**, or **v4**. Changing this parameter will create a new resource.

* `loadbalancer_id` - (Optional, String, ForceNew) Specifies the ID of the load balancer with which the active-standby
  pool is associated. Changing this parameter will create a new resource.

* `listener_id` - (Optional, String, ForceNew) Specifies the ID of the listener with which the active-standby pool is
  associated. Changing this parameter will create a new resource.

  -> **NOTE:** At least one of `loadbalancer_id`, `listener_id`, `type` must be specified.

* `name` - (Optional, String, ForceNew) Specifies the name of the active-standby pool. Changing this parameter will
  create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the active-standby pool. Changing this
  parameter will create a new resource.

* `vpc_id` - (Optional, String, ForceNew) Specifies the ID of the VPC where the active-standby pool works. Changing this
  parameter will create a new resource.

  -> **NOTE:** 1. The active-standby pool must be associated with the VPC.
  <br/> 2. Only backend servers in the VPC or IP as Backend servers can be added.
  <br/> 3. `type` must be set to **instance**.
  <br/> 4. If it is not specified, it is determined by the VPC where the backend server works.
  <br/> 5. This parameter can be updated only when it is left blank.

* `any_port_enable` - (Optional, Bool, ForceNew) Specifies whether to enable forward to same port for active-standby
  pool. If this option is enabled, the listener routes the requests to the backend server over the same port as the
  frontend port. Value options:
  + **false**: Disable forward to same port.
  + **true**: Enable forward to same port.

  Defaults to **false**.

  Changing this parameter will create a new resource.

* `connection_drain_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable delayed logout. This parameter can
  be set to **true** when the `protocol` is set to **TCP**, **UDP** or **QUIC**, and the value of `protocol` of the
  associated listener must be **TCP** or **UDP**. It will be triggered for the following scenes:
  + The pool member is removed from the pool.
  + The health monitor status is abnormal.
  + The pool member weight is changed to 0.

  Changing this parameter will create a new resource.

* `connection_drain_timeout` - (Optional, Int, ForceNew) Specifies the timeout of the delayed logout in seconds. Value
  ranges from `10` to `4,000`.

  Changing this parameter will create a new resource.

<a name="ELB_members"></a>
The `members` block supports:

* `address` - (Required, String, ForceNew) Specifies the private IP address bound to the member.
  + If `subnet_id` is left blank, IP as a Backend is enabled. In this case, the IP address must be an **IPv4** address.
  + If `subnet_id` is not left blank, the IP address can be **IPv4** or **IPv6**. It must be in the subnet specified
    by `subnet_id` and can only be bound to the primary NIC of the backend server.

  Changing this parameter will create a new resource.

* `role` - (Required, String, ForceNew) Specifies the type of the member. Value options:
  + **master**: active backend server.
  + **slave**: standby backend server.

  Changing this parameter will create a new resource.

* `protocol_port` - (Optional, Int, ForceNew) Specifies the port used by the member to receive requests. It is mandatory
  if `any_port_enable` is **false**, and it does not take effect if `any_port_enable` is set to **true**. The value range
  is from `1` to `65,535`. Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the name of the member. It can contain a maximum of `255` characters.
  Changing this parameter will create a new resource.

* `subnet_id` - (Optional, String, ForceNew) Specifies the ID of the IPv4 or IPv6 subnet where the member resides.
  + The IPv4 or IPv6 subnet must be in the same VPC as the subnet of the load balancer.
  + If this parameter is not passed, IP as a Backend has been enabled for the load balancer. In this case, IP as backend
    servers must use private IPv4 addresses, and the protocol of the active-standby pool must be **TCP**, **HTTP**, or
    **HTTPS**.

  Changing this parameter will create a new resource.

<a name="ELB_healthmonitor"></a>
The `healthmonitor` block supports:

* `delay` - (Required, Int, ForceNew) Specifies the interval between health checks, in seconds. The value range is from
  `1` to `50`. Changing this parameter will create a new resource.

* `max_retries` - (Required, Int, ForceNew) Specifies the number of consecutive health checks when the health check
  result of a backend server changes from **OFFLINE** to **ONLINE**. The value range is from `1` to `10`. Changing
  this parameter will create a new resource.

* `timeout` - (Required, Int, ForceNew) Specifies the maximum time required for waiting for a response from the health
  check, in seconds. It is recommended that you set the value less than that of parameter `delay`. The value range is
  from `1` to `50`. Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the health check protocol. Value options: **TCP**, **UDP_CONNECT**,
  **HTTP**, and **HTTPS**.
  + If the protocol of the backend server is **QUIC**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **UDP**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **TCP**, the value can only be **TCP**, **HTTP**, or **HTTPS**.
  + If the protocol of the backend server is **HTTP**, the value can only be **TCP**, **HTTP**, or **HTTPS**.
  + If the protocol of the backend server is **HTTPS**, the value can only be **TCP**, **HTTP**, or **HTTPS**.

  Changing this parameter will create a new resource.

* `domain_name` - (Optional, String, ForceNew) Specifies the domain name that HTTP requests are sent to during the health
  check. The value can contain only digits, letters, hyphens (-), and periods (.) and must start with a digit or letter.
  The value is left blank by default, indicating that the virtual IP address of the load balancer is used as the
  destination address of HTTP requests. This parameter is available only when `type` is set to **HTTP**. The length
  range of value is from `1` to `100`. Changing this parameter will create a new resource.

* `expected_codes` - (Optional, String, ForceNew) Specifies the expected HTTP status code. This parameter will take
  effect only when `type` is set to **HTTP** or **HTTPS**. The default value is 200. Multiple status codes can be
  queried in the format of expected_codes=xxx&expected_codes=xxx. The length range of value is from `1` to `64`.
  Value options:
  + A specific value, for example, **200**
  + A list of values that are separated with commas (,), for example, **200**, **202**
  + A value range, for example, **200-204**

  Changing this parameter will create a new resource.

* `http_method` - (Optional, String, ForceNew) Specifies the HTTP method. The value can be **GET**, **HEAD**, **POST**.
  Default to **GET**. This parameter is available when `type` is set to **HTTP** or **HTTPS**.

  Changing this parameter will create a new resource.

* `max_retries_down` - (Optional, Int, ForceNew) Specifies the number of consecutive health checks when the health check
  result of a backend server changes from ONLINE to OFFLINE. The value range is from `1` to `10`. Defaults to `3`.
  Changing this parameter will create a new resource.

* `monitor_port` - (Optional, Int, ForceNew) Specifies the port used for the health check. If this parameter is left
  blank, a port of the backend server will be used by default. The value range is from `1` to `65,535`. Changing this
  parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the health check name. The length range of value is from `1` to `255`.
  Changing this parameter will create a new resource.

* `url_path` - (Optional, String, ForceNew) Specifies the HTTP request path for the health check. The value must start
  with a slash (/), and the default value is /. The value can contain letters, digits, hyphens (-), slashes (/),
  periods (.), percentage signs (%), question marks (?), pound signs (#), ampersand signs (&), and the extended character
  set **_;~!()*[]@$^:',+**. The length range of value is from `1` to `80`. Changing this parameter will create a new
  resource.

  -> **NOTE:** This parameter is available only when `type` is set to **HTTP** or **HTTPS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the active-standby pool.

* `members` - The backend servers in the active-standby pool.
  The [members](#ELB_membersResp) structure is documented below.

* `healthmonitor` - The health check configured for the active-standby pool.
  The [healthmonitor](#ELB_healthmonitorResp) structure is documented below.

* `created_at` - The create time of the active-standby pool.

* `updated_at` - The update time of the active-standby pool.

<a name="ELB_membersResp"></a>
The `members` block supports:

* `id` - The ID of the member.

* `ip_version` - The IP version supported by the member.

* `instance_id` - The ID of the ECS used as the member.

* `member_type` - The type of the member.

* `operating_status` - The health status of the member.

<a name="ELB_healthmonitorResp"></a>
The `healthmonitor` block supports:

* `id` - The health check ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB active-standby pool can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_active_standby_pool.test <id>
```

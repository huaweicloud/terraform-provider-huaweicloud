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

* `lb_algorithm` - (Required, String, ForceNew) Specifies the load balancing algorithm used by the load balancer to route
  requests to backend servers in the associated backend server group. Value options:
  + **ROUND_ROBIN**: weighted round robin.
  + **LEAST_CONNECTIONS**: weighted least connections.
  + **SOURCE_IP**: source IP hash.
  + **QUIC_CID**: connection ID.

  Defaults to **ROUND_ROBIN**.

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

* `quic_cid_hash_strategy` - The multi-path distribution configuration based on destination connection IDs.
  The [quic_cid_hash_strategy](#ELB_quiccidhashstrategyResp) structure is documented below.

* `created_at` - The create time of the active-standby pool.

* `updated_at` - The update time of the active-standby pool.

<a name="ELB_membersResp"></a>
The `members` block supports:

* `id` - The ID of the member.

* `ip_version` - The IP version supported by the member.

* `instance_id` - The ID of the ECS used as the member.

* `member_type` - The type of the member.

* `operating_status` - The health status of the member.

* `reason` - Why health check fails.
  The [reason](#ELB_member_reasonResp) structure is documented below.

* `status` - The health status of the backend server if `listener_id` under status is specified. If `listener_id` under
  status is not specified, operating_status of member takes precedence.
  The [status](#ELB_member_statusResp) structure is documented below.

<a name="ELB_member_reasonResp"></a>
The `reason` block supports:

* `reason_code` - The code of the health check failures. The value can be:
  + **CONNECT_TIMEOUT**: The connection with the backend server times out during a health check.
  + **CONNECT_REFUSED**: The load balancer rejects connections with the backend server during a health check.
  + **CONNECT_FAILED**: The load balancer fails to establish connections with the backend server during a health check.
  + **CONNECT_INTERRUPT**: The load balancer is disconnected from the backend server during a health check.
  + **SSL_HANDSHAKE_ERROR**: The SSL handshakes with the backend server fail during a health check.
  + **RECV_RESPONSE_FAILED**: The load balancer fails to receive responses from the backend server during a health check.
  + **RECV_RESPONSE_TIMEOUT**: The load balancer does not receive responses from the backend server within the timeout
    duration during a health check.
  + **SEND_REQUEST_FAILED**: The load balancer fails to send a health check request to the backend server during a health
    check.
  + **SEND_REQUEST_TIMEOUT**: The load balancer fails to send a health check request to the backend server within the
    timeout duration.
  + **RESPONSE_FORMAT_ERROR**: The load balancer receives invalid responses from the backend server during a health check.
  + **RESPONSE_MISMATCH**: The response code received from the backend server is different from the preset code.

* `expected_response` - The expected HTTP status code. This parameter will take effect only when `type` is set to **HTTP**,
  **HTTPS** or **GRPC**.
  + A specific status code. If `type` is set to **GRPC**, the status code ranges from **0** to **99**. If `type` is set
    to other values, the status code ranges from **200** to **599**.
  + A list of status codes that are separated with commas (,). A maximum of five status codes are supported.
  + A status code range. Different ranges are separated with commas (,). A maximum of five ranges are supported.

* `healthcheck_response` - The returned HTTP status code in the response. This parameter will take effect only when `type`
  is set to **HTTP**, **HTTPS** or **GRPC**.
  + A specific status code. If type is set to **GRPC**, the status code ranges from **0** to **99**. If `type` is set to
    other values, the status code ranges from **200** to **599**.

<a name="ELB_member_statusResp"></a>
The `status` block supports:

* `listener_id` - The ID of the listener associated with the backend server.

* `operating_status` - The health status of the backend server. The value can be:
  + **ONLINE**: The backend server is running normally.
  + **NO_MONITOR**: No health check is configured for the backend server group to which the backend server belongs.
  + **OFFLINE**: The cloud server used as the backend server is stopped or does not exist.

<a name="ELB_healthmonitorResp"></a>
The `healthmonitor` block supports:

* `id` - The health check ID.

<a name="ELB_quiccidhashstrategyResp"></a>
The `quic_cid_hash_strategy` block supports:

* `len` - The length of the hash factor in the connection ID, in byte. This parameter is valid only when `lb_algorithm`
  is **QUIC_CID**. Value range: **1** to **20**.

* `offset` - The start position in the connection ID as the hash factor, in byte. This parameter is valid only when
  `lb_algorithm` is **QUIC_CID**. Value range: **0** to **19**.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB active-standby pool can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_active_standby_pool.test <id>
```

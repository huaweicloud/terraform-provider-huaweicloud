---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# huaweicloud_elb_active_standby_pool

Manages an ELB active standby pool resource within HuaweiCloud.

## Example Usage

### Create a active standby Pool

```hcl
variable "vpc_id" {}

resource "huaweicloud_elb_active_standby_pool" "test" {
  name            = "%s"
  description     = "test"
  protocol        = "TCP"
  vpc_id          = var.vpc_id
  type            = "instance"
  any_port_enable = false

  members {
    address       = "0.0.0.1"
    protocol_port = 45
    role          = "slave"
  }

  members {
    address       = "0.0.0.0"
    protocol_port = 36
    role          = "master"
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

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ELB pool resource. If omitted, the
  provider-level region will be used.

* `protocol` - (Required, String, ForceNew) Specifies the protocol used by the backend server group to receive requests.
  Value options: **TCP**, **UDP** or **QUIC**. The value range is from **1** to **255**.
  + If the listener's protocol is **UDP**, the value must be **UDP** or **QUIC**.
  + If the listener's protocol is **TCP**, the value must be **TCP**.

* `members` - (Required, List, ForceNew) Specifies the backend servers in the active/standby backend server group.
  The [members](#ELB_members) structure is documented below.

* `healthmonitor` - (Required, List, ForceNew) Specifies the health check configured for the active/standby backend
  server group. The [healthmonitor](#ELB_healthmonitor) structure is documented below.

* `type` - (Optional, String) Specifies the type of the backend server group. Value options:
  + **instance**: Any type of backend servers can be added. `vpc_id` must be mandatory.
  + **ip**: Only IP as Backend servers can be added. `vpc_id` cannot be specified.

  -> **NOTE:** If this parameter is not passed, any type of backend servers can be added and will return an empty string.
  This parameter can be updated only when it is left blank.

* `loadbalancer_id` - (Optional, String, ForceNew) Specifies the ID of the load balancer with which the backend server
  group is associated.

* `listener_id` - (Optional, String, ForceNew) Specifies the ID of the listener with which the backend server group is
  associated.

  -> **NOTE:** At least one of `loadbalancer_id`, `listener_id`, `type` must be specified.

* `name` - (Optional, String, ForceNew) Specifies the backend server group name.

* `description` - (Optional, String, ForceNew) Specifies supplementary information about the active/standby backend
  server group.

* `vpc_id` - (Optional, String, ForceNew) Specifies the ID of the VPC where the backend server group works.

  -> **NOTE:** 1. The backend server group must be associated with the VPC.
  <br/> 2. Only backend servers in the VPC or IP as Backend servers can be added.
  <br/> 3. `type` must be set to **instance**.
  <br/> 4. If it is not specified, it is determined by the VPC where the backend server works.
  <br/> 5. This parameter can be updated only when it is left blank.

* `any_port_enable` - (Optional, Bool, ForceNew) Specifies whether to enable Forward to same port for a backend server
  group. If this option is enabled, the listener routes the requests to the backend server over the same port as the
  frontend port. Value options:
  + **false** (default): Disable Forward to Same Port.
  + **true**: Enable Forward to Same Port.

<a name="ELB_members"></a>
The `members` block supports:

* `address` - (Required, String, ForceNew) Specifies the private IP address bound to the backend server. The value range
  is from **1** to **64**.
  + If `subnet_cidr_id` is left blank, IP as a Backend is enabled. In this case, the IP address must be an **IPv4** address.
  + If `subnet_cidr_id` is not left blank, the IP address can be **IPv4** or **IPv6**. It must be in the subnet specified
    by `subnet_cidr_id` and can only be bound to the primary NIC of the backend server.

* `role` - (Required, String, ForceNew) Specifies the type of the backend server. The value range is from **0** to **36**.
  Value options:
  + **master**: active backend server.
  + **slave**: standby backend server.

* `protocol_port` - (Optional, Int, ForceNew) Specifies the port used by the backend server to receive requests.It is
  mandatory if `any_port_enable` is **true**. The value range is from **1** to **65535**.
  
  -> **NOTE:** This parameter can be left blank because it does not take effect if `any_port_enable` is set to true for a
  backend server group.

* `name` - (Optional, String, ForceNew) Specifies the backend server name. The value range is from **0** to **255**.

* `subnet_cidr_id` - (Optional, String, ForceNew) Specifies the ID of the IPv4 or IPv6 subnet where the backend server
  resides. The value range is from **1** to **36**.

  + The IPv4 or IPv6 subnet must be in the same VPC as the subnet of the load balancer.
  + If this parameter is not passed, IP as a Backend has been enabled for the load balancer. In this case, IP as backend
    servers must use private IPv4 addresses, and the protocol of the backend server group must be **TCP**, **HTTP**, or
    **HTTPS**.

<a name="ELB_healthmonitor"></a>
The `healthmonitor` block supports:

* `delay` - (Required, Int, ForceNew) Specifies the interval between health checks, in seconds. The value range is from
  **1** to **50**.

* `max_retries` - (Required, Int, ForceNew) Specifies the number of consecutive health checks when the health check
  result of a backend server changes from **OFFLINE** to **ONLINE**. The value range is from **1** to **10**.

* `timeout` - (Required, Int, ForceNew) Specifies the maximum time required for waiting for a response from the health
  check, in seconds. It is recommended that you set the value less than that of parameter `delay`. The value range is
  from **1** to **50**.

* `type` - (Required, String, ForceNew) Specifies the health check protocol. Value options: **TCP**, **UDP_CONNECT**,
  **HTTP**, and **HTTPS**.
  + If the protocol of the backend server is **QUIC**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **UDP**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **TCP**, the value can only be **TCP**, **HTTP**, or **HTTPS**.
  + If the protocol of the backend server is **HTTP**, the value can only be **TCP**, **HTTP**, or **HTTPS**.
  + If the protocol of the backend server is **HTTPS**, the value can only be **TCP**, **HTTP**, or **HTTPS**.

* `domain_name` - (Optional, String, ForceNew) Specifies the domain name that HTTP requests are sent to during the health
  check. The value can contain only digits, letters, hyphens (-), and periods (.) and must start with a digit or letter.
  The value is left blank by default, indicating that the virtual IP address of the load balancer is used as the
  destination address of HTTP requests. This parameter is available only when `type` is set to **HTTP**. The value range
  is from **1** to **100**.

* `expected_codes` - (Optional, String, ForceNew) Specifies the expected HTTP status code. This parameter will take
  effect only when `type` is set to **HTTP** or **HTTPS**. The default value is 200. Multiple status codes can be
  queried in the format of expected_codes=xxx&expected_codes=xxx. The value range is from **1** to **64**. Value options:
  + A specific value, for example, **200**
  + A list of values that are separated with commas (,), for example, **200**, **202**
  + A value range, for example, **200**-**204**

* `max_retries_down` - (Optional, Int, ForceNew) Specifies the number of consecutive health checks when the health check
  result of a backend server changes from ONLINE to OFFLINE. The value range is from **1** to **10**, and the default
  value is **3**.

* `monitor_port` - (Optional, Int, ForceNew) Specifies the port used for the health check. If this parameter is left
  blank, a port of the backend server will be used by default. The value range is from **1** to **65535**.

* `name` - (Optional, String, ForceNew) Specifies the health check name. The value range is from **1** to **255**.

* `url_path` - (Optional, String, ForceNew) Specifies the HTTP request path for the health check. The value must start
  with a slash (/), and the default value is /. The value can contain letters, digits, hyphens (-), slashes (/),
  periods (.), percentage signs (%), question marks (?), pound signs (#), ampersand signs (&), and the extended character
  set _;~!()*[]@$^:',+. The value range is from **1** to **80**.

  -> **NOTE:** This parameter is available only when `type` is set to **HTTP** or **HTTPS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the pool.

* `members` - The backend servers in the active/standby backend server group.
  The [members](#ELB_membersResp) structure is documented below.

* `healthmonitor` - The health check configured for the active/standby backend server group.
  The [healthmonitor](#ELB_healthmonitorResp) structure is documented below.

<a name="ELB_membersResp"></a>
The `members` block supports:

* `ip_version` - The IP version supported by the backend server.

* `instance_id` - The ID of the ECS used as the backend server.

* `member_type` - The type of the backend server.

* `operating_status` - The health status of the backend server.

* `id` - The ID of the backend server group.

<a name="ELB_healthmonitorResp"></a>
The `healthmonitor` block supports:

* `id` - The health check ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB active standby pool can be imported using the pool `id`, e.g.

```bash
$ terraform import huaweicloud_elb_active_standby_pool.test <id>
```

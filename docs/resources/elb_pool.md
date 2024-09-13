---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_pool"
description: ""
---

# huaweicloud_elb_pool

Manages an ELB pool resource within HuaweiCloud.

## Example Usage

### Create a Pool and Associate a Load Balancer

```hcl
variable "loadbalancer_id" {}

resource "huaweicloud_elb_pool" "pool_1" {
  name            = "test"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = var.loadbalancer_id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test reason"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

### Create a Pool and Associate a Listener

```hcl
variable "listener_id" {}

resource "huaweicloud_elb_pool" "pool_1" {
  name        = "test"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = var.listener_id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test reason"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

### Create a Pool and Associate later

```hcl
variable "vpc_id" {}

resource "huaweicloud_elb_pool" "pool_1" {
  name      = "test"
  protocol  = "HTTP"
  lb_method = "ROUND_ROBIN"
  type      = "instance"
  vpc_id    = var.vpc_id

  slow_start_enabled  = true
  slow_start_duration = 100

  protection_status = "consoleProtection"
  protection_reason = "test reason"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the ELB pool resource. If omitted, the
  provider-level region will be used. Changing this creates a new pool.

* `protocol` - (Required, String, ForceNew) Specifies the protocol used by the backend server group to receive requests.
  Value options: **TCP**, **UDP**, **HTTP**, **HTTPS**, **QUIC**, **GRPC** or **TLS**.
  + If the listener's protocol is **UDP**, the value must be **UDP** or **QUIC**.
  + If the listener's protocol is **TCP**, the value must be **TCP**.
  + If the listener's protocol is **IP**, the value must be **IP**.
  + If the listener's protocol is **HTTP**, the value must be **HTTP**.
  + If the listener's protocol is **HTTPS**, the value must be **HTTP** or **HTTPS**.
  + If the listener's protocol is **TERMINATED_HTTPS**, the value must be **HTTP**.
  + If the listener's protocol is **QUIC**, the value must be **HTTP**、**HTTPS** or **GRPC**.
  + If the listener's protocol is **TLS**, the value must be **TLS** or **TCP**.
  + If the value is **QUIC**, sticky session must be enabled with `type` set to **SOURCE_IP**.
  + If the value is **GRPC**, the value of `http2_enable` of the associated listener must be **true**.

  Changing this creates a new pool.

* `lb_method` - (Required, String) Specifies the load balancing algorithm used by the load balancer to route requests
  to backend servers in the associated backend server group. Value options:
  + **ROUND_ROBIN**: weighted round-robin.
  + **LEAST_CONNECTIONS**: weighted least connections.
  + **SOURCE_IP**: source IP hash.
  + **QUIC_CID**: connection ID.
  + **2_TUPLE_HASH**: 2-tuple hash that is only available for IP backend server groups.
  + **3_TUPLE_HASH**: 3-tuple hash that is only available for IP backend server groups.
  + **5_TUPLE_HASH**: 5-tuple hash that is only available for IP backend server groups Note.

  -> **NOTE:** 1. If the value is **SOURCE_IP**, the weight parameter will not take effect for backend servers.
  <br/> 2. **QUIC_CID** is supported only when the protocol of the backend server group is **QUIC**.

* `persistence` - (Optional, List) Specifies the sticky session.
  The [object](#persistence) structure is documented below.

* `type` - (Optional, String) Specifies the type of the backend server group. Value options:
  + **instance**: Any type of backend servers can be added. `vpc_id` must be mandatory.
  + **ip**: Only IP as Backend servers can be added. `vpc_id` cannot be specified.

  -> **NOTE:** If this parameter is not passed, any type of backend servers can be added and will return an empty string.
  This parameter can be updated only when it is left blank.

* `loadbalancer_id` - (Optional, String, ForceNew) Specifies the ID of the load balancer with which the backend server
  group is associated. Changing this creates a new pool.

* `listener_id` - (Optional, String, ForceNew) Specifies the ID of the listener with which the backend server group is
  associated. Changing this creates a new pool.

  -> **NOTE:** At least one of `loadbalancer_id`, `listener_id`, `type` must be specified.

* `ip_version` - (Optional, String, ForceNew) Specifies the IP address version supported by the backend server group.
  The value can be **dualstack**, **v6**, or **v4**. If the protocol of the backend server group is HTTP, the value is **v4**.
  Changing this creates a new pool.

* `any_port_enable` - (Optional, Bool, ForceNew) Specifies whether to enable transparent port transmission on the backend.
  If enable, the port of the backend server will be same as the port of the listener.
  Changing this creates a new pool.

* `deletion_protection_enable` - (Optional, Bool) Specifies whether to enable deletion protection.

* `name` - (Optional, String) Specifies the backend server group name.

* `description` - (Optional, String) Specifies the description of the pool.

* `vpc_id` - (Optional, String) Specifies the ID of the VPC where the backend server group works.

  -> **NOTE:** 1. The backend server group must be associated with the VPC.
  <br/> 2. Only backend servers in the VPC or IP as Backend servers can be added.
  <br/> 3. `type` must be set to **instance**.
  <br/> 4. If it is not specified, it is determined by the VPC where the backend server works.
  <br/> 5. This parameter can be updated only when it is left blank.

* `protection_status` - (Optional, String) The protection status for update. Value options:
  + **nonProtection**: No protection.
  + **consoleProtection**: Console modification protection.

  Defaults to **nonProtection**.

* `protection_reason` - (Optional, String) The reason for update protection. Only valid when `protection_status` is
  **consoleProtection**.

* `slow_start_enabled` - (Optional, Bool) Specifies whether to enable slow start. After you enable slow start, new
  backend servers added to the backend server group are warmed up, and the number of requests they can receive
  increases linearly during the configured slow start duration. Defaults to **false**.

  -> **NOTE:** This parameter can be set to **true** when the `protocol` is set to **HTTP** or **HTTPS**, or an error
  will be returned.

* `slow_start_duration` - (Optional, Int) Specifies the slow start duration, in seconds.  
  Value ranges from `30` to `1,200`. Defaults to `30`.

* `connection_drain_enabled` - (Optional, Bool) Specifies whether to enable delayed logout. This parameter can be set to
  **true** when the `protocol` is set to **TCP**, **UDP** or **QUIC**, and the value of `protocol` of the associated
  listener must be **TCP** or **UDP**. It will be triggered for the following scenes:
  + The pool member is removed from the pool.
  + The health monitor status is abnormal.
  + The pool member weight is changed to `0`.

* `connection_drain_timeout` - (Optional, Int) Specifies the timeout of the delayed logout in seconds.  
  Value ranges from `10` to `4000`.

* `minimum_healthy_member_count` - (Optional, Int) Specifies the minimum healthy member count. When the number of online
  members in the health check is less than this number, the status of the pool is determined to be unhealthy. Value options:
  + **0** (default value): Not take effect.
  + **1**: Take effect when all member offline.

<a name="persistence"></a>
The `persistence` block supports:

* `type` - (Required, String) Specifies the sticky session type. Value options: **SOURCE_IP**,
  **HTTP_COOKIE**, and **APP_COOKIE**.

  -> **NOTE:** 1. If the protocol of the backend server group is **TCP** or **UDP**, only **SOURCE_IP** takes effect.
  <br/> 2. If the protocol of the backend server group is **HTTP** or **HTTPS**, the value can only be **HTTP_COOKIE**.
  <br/> 3. If the backend server group protocol is **QUIC**, sticky session must be enabled with type set to
  **SOURCE_IP**.

* `cookie_name` - (Optional, String) Specifies the cookie name. The value can contain only letters, digits,
  hyphens (-), underscores (_), and periods (.). It is required if `type` of `persistence` is set to **APP_COOKIE**.

* `timeout` - (Optional, Int) Specifies the sticky session timeout duration in minutes. This parameter is
  invalid when `type` is set to **APP_COOKIE**. The value range varies depending on the protocol of the backend server
  group:
  + When the protocol of the backend server group is **TCP** or **UDP**, the value ranges from `1` to `60`, and
    defaults to `1`.
  + When the protocol of the backend server group is **HTTP** or **HTTPS**, the value ranges from `1` to `1,440`,
    and defaults to `1,440`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the pool.

* `monitor_id` - The ID of the health check configured for the backend server group.

* `created_at` - The create time of the pool.

* `updated_at` - The update time of the pool.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB pool can be imported using the pool `id`, e.g.

```bash
$ terraform import huaweicloud_elb_pool.pool_1 <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `deletion_protection_enable`. It is
generally recommended running **terraform plan** after importing a pool. You can then decide if changes should be
applied to the pool, or the resource definition should be updated to align with the pool. Also you can ignore changes
as below.

```hcl
resource "huaweicloud_elb_pool" "test" {
    ...

  lifecycle {
    ignore_changes = [
      deletion_protection_enable,
    ]
  }
}
```

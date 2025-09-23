---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_monitor"
description: |-
  Manages an ELB monitor resource within HuaweiCloud.
---

# huaweicloud_lb_monitor

Manages an ELB monitor resource within HuaweiCloud.

## Example Usage

### TCP Health Check

```hcl
resource "huaweicloud_lb_monitor" "monitor_tcp" {
  pool_id     = var.pool_id
  type        = "TCP"
  delay       = 5
  timeout     = 3
  max_retries = 3
}
```

### UDP Health Check

```hcl
resource "huaweicloud_lb_monitor" "monitor_udp" {
  pool_id     = var.pool_id
  type        = "UDP_CONNECT"
  delay       = 5
  timeout     = 3
  max_retries = 3
}
```

### HTTP Health Check

```hcl
resource "huaweicloud_lb_monitor" "monitor_http" {
  pool_id        = var.pool_id
  type           = "HTTP"
  delay          = 5
  timeout        = 3
  max_retries    = 3
  url_path       = "/test"
  expected_codes = "200-202"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB monitor resource. If omitted, the
  provider-level region will be used. Changing this creates a new monitor.

* `pool_id` - (Required, String, ForceNew) Specifies the id of the pool that this monitor will be assigned to. Changing
  this creates a new monitor.

* `type` - (Required, String) Specifies the monitor protocol.
  The value can be **TCP**, **UDP_CONNECT** or **HTTP**.
  If the listener protocol is **UDP**, the monitor protocol must be **UDP_CONNECT**.

* `delay` - (Required, Int) Specifies the maximum time between health checks in the unit of second. The value ranges
  from **1** to **50**.

* `timeout` - (Required, Int) Specifies the health check timeout duration in the unit of second.
  The value ranges from **1** to **50** and must be less than the `delay` value.

* `max_retries` - (Required, Int) Specifies the maximum number of consecutive health checks after which the backend
  servers are declared **healthy**. The value ranges from **1** to **10**.

  -> Backend servers can be declared **unhealthy** after **three** consecutive health checks that detect these backend
  servers are unhealthy, regardless of the value set for `max_retries`. The health check time window is determined
  by [Health Check Time Window](https://support.huaweicloud.com/intl/en-us/usermanual-elb/elb_ug_hc_0001.html#section4).

* `name` - (Optional, String) Specifies the health check name. The value contains a maximum of 255 characters.

* `port` - (Optional, Int) Specifies the health check port. The port number ranges from `1` to `65,535`. If not specified,
  the port of the backend server will be used as the health check port.

* `url_path` - (Optional, String) Specifies the HTTP request path for the health check. Required for HTTP type.
  The value starts with a slash (/) and contains a maximum of 255 characters.

* `http_method` - (Optional, String) Specifies the HTTP request method. Required for HTTP type.
  The default value is **GET**.

* `expected_codes` - (Optional, String) Specifies the expected HTTP status code. Required for HTTP type.
  You can either specify a single status like **200**, or a range like **200-202**.

* `domain_name` - (Optional, String) Specifies the domain name of HTTP requests during the health check. It takes effect
  only when the value of `type` is set to **HTTP**. The value is left blank by default, indicating that the private IP
  address of the load balancer is used as the destination address of HTTP requests. The value can contain only digits,
  letters, hyphens (-), and periods (.) and must start with a digit or letter, the value contains a maximum of 100 characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the monitor.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB monitor can be imported using the monitor ID, e.g.

```bash
$ terraform import huaweicloud_lb_monitor.monitor_1 5c20fdad-7288-11eb-b817-0255ac10158b
```

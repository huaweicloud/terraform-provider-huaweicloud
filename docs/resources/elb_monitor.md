---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_monitor"
description: ""
---

# huaweicloud_elb_monitor

Manages an ELB monitor resource within HuaweiCloud.

## Example Usage

```hcl
variable "pool_id" {}

resource "huaweicloud_elb_monitor" "monitor_1" {
  pool_id     = var.pool_id
  protocol    = "HTTPS"
  interval    = 30
  timeout     = 20
  max_retries = 8
  url_path    = "/bb"
  domain_name = "www.bb.com"
  port        = 8888
  status_code = "200,301,404-500,504"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the ELB monitor resource. If omitted, the
  provider-level region will be used. Changing this creates a new monitor.

* `pool_id` - (Required, String, ForceNew) Specifies the ID of the backend server group for which the health check is
  configured. Changing this creates a new monitor.

* `protocol` - (Required, String) Specifies the health check protocol. Value options: **TCP**, **UDP_CONNECT**,
  **HTTP**, **HTTPS**, **GRPC** or **TLS**.
  + If the protocol of the backend server is **QUIC**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **UDP**, the value can only be **UDP_CONNECT**.
  + If the protocol of the backend server is **TCP**, the value can only be **TCP**, **HTTP** or **HTTPS**.
  + If the protocol of the backend server is **HTTP**, the value can only be **TCP**, **HTTP**, **HTTPS**, **TLS** or **GRPC**.
  + If the protocol of the backend server is **HTTPS**, the value can only be **TCP**, **HTTP**, **HTTPS**, **TLS** or **GRPC**.
  + If the protocol of the backend server is **GRPC**, the value can only be **TCP**, **HTTP**, **HTTPS**, **TLS** or **GRPC**.
  + If the protocol of the backend server is **TLS**, the value can only be **TCP**, **HTTP**, **HTTPS**, **TLS** or **GRPC**.

* `interval` - (Required, Int) Specifies the interval between health checks, in seconds.
  Value ranges from `1` to `50`.

* `timeout` - (Required, Int) Specifies the maximum time required for waiting for a response from the health check,
  in seconds. Value ranges from `1` to `50`. It is recommended that you set the value less than that of
  parameter `interval`.

* `max_retries` - (Required, Int) Specifies the number of consecutive health checks when the health check result of
  a backend server changes from OFFLINE to ONLINE. Value ranges from `1` to `10`.

* `max_retries_down` - (Optional, Int) Specifies the number of consecutive health checks when the health check result of
  a backend server changes from ONLINE to OFFLINE. The value ranges from `1` to `10`, and the default value is `3`.

* `name` - (Optional, String) Specifies the health check name.

* `domain_name` - (Optional, String) Specifies the domain name that HTTP requests are sent to during the health check.
  The domain name consists of 1 to 100 characters, can contain only digits, letters, hyphens (-), and periods (.) and
  must start with a digit or letter. The value is left blank by default, indicating that the virtual IP address of the
  load balancer is used as the destination address of HTTP requests. This parameter is available only when `protocol`
  is set to **HTTP** or **HTTPS**.

* `port` - (Optional, Int) Specifies the port used for the health check. If this parameter is left blank, a port of
  the backend server will be used by default. It is mandatory when the `protocol` of the backend server group is **IP**.
  Value ranges from `1` to `65,535`.

* `url_path` - (Optional, String) Specifies the HTTP request path for the health check. The value must start with a
  slash (/), can contain letters, digits, hyphens (-), slash (/), periods (.), percent signs (%), hashes(#), and(&)
  and the special characters: `~!()*[]@$^:',+`, and the default value is **/**. This parameter is available only when
  `protocol` is set to **HTTP** or **HTTPS**.

* `status_code` - (Optional, String) Specifies the expected HTTP status code. This parameter will take effect only when
  `protocol` is set to **HTTP** or **HTTPS**. Value options are as follows:
  + A specific value, for example: **200**.
  + A list of values that are separated with commas (,), for example: **200,202**.
  + A value range, for example: **200-204**.

  Defaults to **200**.

* `http_method` - (Optional, String) Specifies the HTTP method. Value options: **GET**, **HEAD**, **POST**. Defaults to **GET**.

* `enabled` - (Optional, Bool) Specifies whether the health check is enabled.
  + **true(default)**: Health check is enabled.
  + **false**: Health check is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID of the monitor.

* `created_at` - The creation time of the monitor.

* `updated_at` - The update time of the monitor.

## Import

ELB monitor can be imported using the monitor `id`, e.g.

```bash
$ terraform import huaweicloud_elb_monitor.test <id>
```

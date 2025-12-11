---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_monitors"
description: |-
  Use this data source to get the list of ELB monitors.
---

# huaweicloud_elb_monitors

Use this data source to get the list of ELB monitors.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_elb_monitors" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source. If omitted, the provider-level
  region will be used.

* `pool_id` - (Optional, String) Specifies the ID of backend server groups for which the health check is configured.

* `monitor_id` - (Optional, String) Specifies the health check ID.

* `port` - (Optional, Int) Specifies the port used for the health check.

* `domain_name` - (Optional, String) Specifies the domain name to which HTTP requests are sent during the health check.
  The value can be digits, letters, hyphens (-), or periods (.) and must start with a digit or letter.

* `name` - (Optional, String) Specifies the health check name.

* `interval` - (Optional, Int)  Specifies the interval between health checks, in seconds.  
  The value ranges from `1` to `50`.

* `max_retries` - (Optional, Int) Specifies the number of consecutive health checks when the health check result of a
  backend server changes from **OFFLINE** to **ONLINE**.

* `max_retries_down` - (Optional, Int) Specifies the number of consecutive health checks when the health check result of
  a backend server changes from **ONLINE** to **OFFLINE**. The value ranges from `1` to `10`.

* `timeout` - (Optional, Int) Specifies the maximum time required for waiting for a response from the health check, in
  seconds.

* `protocol` - (Optional, String) Specifies the health check protocol. The value can be **TCP**, **UDP_CONNECT**,
  **HTTP**, **HTTPS**, **GRPC** or **TLS**.

* `status_code` - (Optional, String) Specifies the expected HTTP status code. This parameter will take effect only when
  type is set to **HTTP** or **HTTPS**.Value options:
  + A specific value, for example, **200**
  + A list of values that are separated with commas (,), for example, **200**, **202**
  + A value range, for example, **200**-**204**

* `http_method` - (Optional, String)  Specifies the HTTP method. Value options: **GET**, **HEAD**, **POST**.

* `url_path` - (Optional, String) Specifies the HTTP request path for the health check. The value must start with a
  slash(/), and the default value is **/**. This parameter is available only when type is set to **HTTP**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `monitors` - Lists the monitors.
  The [monitors](#Elb_monitors) structure is documented below.

<a name="Elb_monitors"></a>
The `monitors` block supports:

* `id` - The health check ID.

* `name` - The health check name.

* `domain_name` - The domain name that HTTP requests are sent to during the health check.

* `interval` - The interval between health checks, in seconds.

* `status_code` - The expected HTTP status code.

* `http_method` - The HTTP method

* `max_retries` - The number of consecutive health checks when the health check result of a backend server changes from
  **OFFLINE** to **ONLINE**.

* `max_retries_down` - The number of consecutive health checks when the health check result of a backend server changes
  from **ONLINE** to **OFFLINE**.

* `port` - The port used for the health check.

* `pool_id` - The ID of backend server groups for which the health check is configured.

* `timeout` - The maximum time required for waiting for a response from the health check, in seconds.

* `protocol` - The health check protocol.

* `url_path` - The HTTP request path for the health check.

* `created_at` - The time when the health check was configured.

* `updated_at` - The time when the health check was updated.

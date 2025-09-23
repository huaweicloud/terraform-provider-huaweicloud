---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_listener"
description: |-
  Manages an ELB listener resource within HuaweiCloud.
---

# huaweicloud_lb_listener

Manages an ELB listener resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_lb_listener" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the listener resource. If omitted, the
  provider-level region will be used. Changing this creates a new listener.

* `protocol` - (Required, String, NonUpdatable) The protocol can either be **TCP**, **UDP**, **HTTP** or **TERMINATED_HTTPS**.

* `protocol_port` - (Required, Int, NonUpdatable) The port on which to listen for client traffic.

* `loadbalancer_id` - (Required, String, NonUpdatable) The load balancer on which to provision this listener.

* `name` - (Optional, String) Human-readable name for the listener. Does not have to be unique.

* `default_pool_id` - (Optional, String) The ID of the default pool with which the listener is associated.

* `description` - (Optional, String) Human-readable description for the listener.

* `http2_enable` - (Optional, Bool) Specifies whether to use HTTP/2. The default value is false. This parameter is valid
  only when the `protocol` is set to **TERMINATED_HTTPS**.

* `default_tls_container_ref` - (Optional, String) Specifies the ID of the server certificate used by the listener. This
  parameter is mandatory when `protocol` is set to **TERMINATED_HTTPS**.

* `client_ca_tls_container_ref` - (Optional, String) Specifies the ID of the CA certificate used by the listener. This
  parameter is mandatory when `protocol` is set to **TERMINATED_HTTPS**.

* `sni_container_refs` - (Optional, List) Lists the IDs of SNI certificates (server certificates with a domain name)
  used by the listener. This parameter is valid when `protocol` is set to **TERMINATED_HTTPS**.

* `insert_headers` - (Optional, List) Specifies whether to insert HTTP extension headers and sent them to backend servers.
  All headers are synchronized. If this parameter is not set, default values are used. Information required by backend
  servers can be written into HTTP headers and passed to backend servers. This parameter is mandatory when `protocol` is
  set to **TERMINATED_HTTPS**.
  The [insert_headers](#insert_headers_struct) structure is documented below.

* `tls_ciphers_policy` - (Optional, String) Specifies the security policy used by the listener. This parameter takes effect
  only when the `protocol` used by the listener is set to **TERMINATED_HTTPS**. Value options: **tls-1-0-inherit**,
  **tls-1-0**, **tls-1-1**, **tls-1-2** or **tls-1-2-strict**. Defaults to **tls-1-0**.

* `protection_status` - (Optional, String) Specifies whether modification protection is enabled. Value options:
  + **nonProtection (default)**: Modification protection is not enabled.
  + **consoleProtection**: Modification protection is enabled to avoid that resources are modified by accident on the console.

* `protection_reason` - (Optional, String) The reason to enable modification protection. This parameter is valid only when
  `protection_status` is set to **consoleProtection**.

* `tags` - (Optional, Map) Specifies the reason to enable modification protection.

<a name="insert_headers_struct"></a>
The `insert_headers` block supports:

* `x_forwarded_elb_ip` - (Optional, String) Specifies whether to transparently transmit the load balancer EIP to backend
  servers. After this function is enabled, the load balancer EIP is stored in the HTTP header and passes to backend servers.
  Value options:
  + **true**: This function is enabled.
  + **false (default)**: The function is disabled.

* `x_forwarded_host` - (Optional, String) Specifies whether to rewrite the X-Forwarded-Host header. If this function is
  enabled, **X-Forwarded-Host** is rewritten based on Host in the request and sent to backend servers. Value options:
  + **true (default)**: This function is enabled.
  + **false**: The function is disabled.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the listener.

* `created_at` - The creation time of the listener.

* `updated_at` - The update time of the listener.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB listener can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_lb_listener.test <id>
```

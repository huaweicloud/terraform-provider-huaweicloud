---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
---

# huaweicloud_elb_listener

Manages an ELB listener resource within HuaweiCloud.

## Example Usage

```hcl
variable "loadbalancer_id" {}

resource "huaweicloud_elb_listener" "basic" {
  name            = "basic"
  description     = "basic description"
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = var.loadbalancer_id

  idle_timeout     = 60
  request_timeout  = 60
  response_timeout = 60

  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the listener resource. If omitted, the
  provider-level region will be used. Changing this creates a new listener.

* `protocol` - (Required, String, ForceNew) The protocol can either be TCP, UDP, HTTP or HTTPS. Changing this creates a
  new listener.

* `loadbalancer_id` - (Required, String, ForceNew) The load balancer on which to provision this listener. Changing this
  creates a new listener.

* `protocol_port` - (Optional, Int, ForceNew) Specifies the port used by the listener. The **QUIC** listener port cannot
  be **4789** or the same as the **UDP** listener port. If this parameter is set to **0**, `port_ranges` is required.
  Changing this creates a new listener.

* `port_ranges` - (Optional, List, ForceNew) Specifies the port monitoring range (closed range), specify up to 10 port
  groups, each group range must not overlap. This field can only be passed in when `protocol_port` is **0** or empty.
  Only **TCP**, **UDP**, and **TCPSSL** listener support this field. Changing this creates a new listener.
  The [port_ranges](#ELB_port_ranges) structure is documented below.

* `name` - (Optional, String) Human-readable name for the listener.

* `default_pool_id` - (Optional, String) The ID of the default pool with which the listener is associated. Changing this
  creates a new listener.

* `description` - (Optional, String) Human-readable description for the listener.

* `http2_enable` - (Optional, Bool) Specifies whether to use HTTP/2. The default value is false. This parameter is valid
  only when the protocol is set to *HTTPS*.

* `forward_eip` - (Optional, Bool) Specifies whether transfer the load balancer EIP in the X-Forward-EIP header to
  backend servers. The default value is false. This parameter is valid only when the protocol is set to *HTTP* or
  *HTTPS*.

* `forward_port` - (Optional, Bool) Specifies whether transfer the listening port of the load balancer in the
  X-Forwarded-Port header to backend servers. The default value is false. This parameter is valid only when the
  protocol is set to *HTTP* or *HTTPS*.

* `forward_request_port` - (Optional, Bool) Specifies whether transfer the source port of the client in the
  X-Forwarded-For-Port header to backend servers. The default value is false. This parameter is valid only when the
  protocol is set to *HTTP* or *HTTPS*.

* `forward_host` - (Optional, Bool) Specifies whether to rewrite the X-Forwarded-Host header. If X-Forwarded-Host is
  set to true, X-Forwarded-Host in the request header from the clients can be set to Host in the request header sent
  from the load balancer to backend servers. The default value is true. This parameter is valid only when the protocol
  is set to *HTTP* or *HTTPS*.

* `access_policy` - (Optional, String) Specifies the access policy for the listener. Valid options are *white* and
  *black*.

* `ip_group` - (Optional, String) Specifies the ip group id for the listener.

* `server_certificate` - (Optional, String) Specifies the ID of the server certificate used by the listener. This
  parameter is mandatory when protocol is set to *HTTPS*.

* `sni_certificate` - (Optional, List) Lists the IDs of SNI certificates (server certificates with a domain name) used
  by the listener. This parameter is valid when protocol is set to *HTTPS*.

* `ca_certificate` - (Optional, String) Specifies the ID of the CA certificate used by the listener. This parameter is
  valid when protocol is set to *HTTPS*.

* `tls_ciphers_policy` - (Optional, String) Specifies the TLS cipher policy for the listener. Valid options are:
  tls-1-0-inherit, tls-1-0, tls-1-1, tls-1-2, tls-1-2-strict, tls-1-2-fs, tls-1-0-with-1-3, and tls-1-2-fs-with-1-3.
  This parameter is valid when protocol is set to *HTTPS*.

* `idle_timeout` - (Optional, Int) Specifies the idle timeout for the listener. Value range: 0 to 4000.

* `request_timeout` - (Optional, Int) Specifies the request timeout for the listener. Value range: 1 to 300. This
  parameter is valid when protocol is set to *HTTP* or *HTTPS*.

* `response_timeout` - (Optional, Int) Specifies the response timeout for the listener. Value range: 1 to 300. This
  parameter is valid when protocol is set to *HTTP* or *HTTPS*.

* `advanced_forwarding_enabled` - (Optional, Bool) Specifies whether to enable advanced forwarding.
  If advanced forwarding is enabled, more flexible forwarding policies and rules are supported.

* `protection_status` - (Optional, String) The protection status for update. Value options:
  + **nonProtection**: No protection.
  + **consoleProtection**: Console modification protection.

  Defaults to **nonProtection**.

* `protection_reason` - (Optional, String) The reason for update protection. Only valid when `protection_status` is
  **consoleProtection**.

* `force_delete` - (Optional, Bool) Specifies whether to forcibly delete the listener, remove the listener, l7 policies,
  unbind associated pools. Defaults to **false**.

* `tags` - (Optional, Map) The key/value pairs to associate with the listener.

* `gzip_enable` - (Optional, Bool) Specifies whether to enable gzip compression for a load balancer. The default value
  is **false**. This parameter can be configured only for **HTTP**, **HTTPS**, and **QUIC** listeners.

<a name="ELB_port_ranges"></a>
The `port_ranges` block supports:

* `start_port` - (Required, Int, ForceNew) Specifies the start port. Changing this creates a new listener.

* `end_port` - (Required, Int, ForceNew) Specifies the end port. Changing this creates a new listener.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique ID for the listener.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

ELB listener can be imported using the listener ID, e.g.

```
$ terraform import huaweicloud_elb_listener.listener_1 5c20fdad-7288-11eb-b817-0255ac10158b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `force_delete`. It is generally recommended
running `terraform plan` after importing a listener. You can then decide if changes should be applied to the listener,
or the resource definition should be updated to align with the listener. Also you can ignore changes as below.

```hcl
resource "huaweicloud_elb_listener" "listener_1" {
    ...
  lifecycle {
    ignore_changes = [
      force_delete,
    ]
  }
}
```

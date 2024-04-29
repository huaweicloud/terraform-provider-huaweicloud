---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_listeners"
description: ""
---

# huaweicloud_elb_listeners

Use this data source to get the list of ELB listeners.

## Example Usage

```hcl
variable "listener_name" {}

data "huaweicloud_elb_listeners" "test" {
  name = var.listener_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the ELB listener.

* `description` - (Optional, String) Specifies the description of the ELB listener.

* `listener_id` - (Optional, String) Specifies the ID of the ELB listener.

* `loadbalancer_id` - (Optional, String) Specifies the ID of the load balancer that the listener is added to.

* `protocol_port` - (Optional, Int) Specifies the port used by the listener.

* `protocol` - (Optional, String) Specifies the protocol of the ELB listener. Value options:
  **TCP**, **UDP**, **HTTP**, **HTTPS** or **QUIC**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `listeners` - Lists the listeners.
  The [listeners](#Elb_loadbalancer_listeners) structure is documented below.

<a name="Elb_loadbalancer_listeners"></a>
The `listeners` block supports:

* `id` - The listener ID.

* `name` - The listener name.

* `description` - The description of the listener.

* `protocol` - The protocol used by the listener.

* `protocol_port` - The port used by the listener.

* `default_pool_id` - The ID of the default backend server group.

* `http2_enable` - Whether to use HTTP/2 if you want the clients to use HTTP/2 to communicate with the listener.

* `forward_eip` - Whether to transparently transmit the load balancer EIP to backend servers.

* `forward_port` - Whether to transparently transmit the listening port of the load balancer to backend servers.

* `forward_request_port` - Whether to transparently transmit the source port of the client to backend servers.

* `forward_host` - Whether to rewrite the X-Forwarded-Host header.

* `sni_certificate` - The IDs of SNI certificates (server certificates with domain names) used by the listener.

* `ca_certificate` - The ID of the CA certificate used by the listener.

* `server_certificate` - The ID of the server certificate used by the listener.

* `tls_ciphers_policy` - The security policy used by the listener.

* `idle_timeout` - The idle timeout duration, in seconds.

* `request_timeout` - The timeout duration for waiting for a response from a client, in seconds.

* `response_timeout` - The timeout duration for waiting for a response from a backend server, in seconds.

* `loadbalancer_id` - The ID of the load balancer that the listener is added to.

* `advanced_forwarding_enabled` - Whether to enable advanced forwarding.

* `protection_status` - The protection status for update.

* `protection_reason` - The reason for update protection.

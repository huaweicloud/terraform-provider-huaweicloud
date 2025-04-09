---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_listeners"
description: |-
  Use this data source to get the list of ELB listeners.
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

* `ca_certificate` - (Optional, String) Specifies the ID of the CA certificate used by the listener.

* `request_timeout` - (Optional, Int) Specifies the request timeout for the listener. Value range: **1** to **300**.

* `default_pool_id` - (Optional, String) Specifies the ID of the default pool with which the listener is associated.

* `server_certificate` - (Optional, String) Specifies the ID of the server certificate used by the listener.

* `enable_member_retry` - (Optional, String) Specifies whether the health check retries for backend servers is enabled.
  Value options: **true**, **false**.

* `advanced_forwarding_enabled` - (Optional, String) Specifies whether the advanced forwarding is enabled. Value options:
  **true**, **false**.

* `http2_enable` - (Optional, String) Specifies whether the HTTP/2 is used. Value options: **true**, **false**.

* `idle_timeout` - (Optional, Int) Specifies the idle timeout for the listener.

* `member_address` - (Optional, String) Specifies the private IP address bound to the backend server.

* `member_device_id` - (Optional, String) Specifies the ID of the cloud server that serves as a backend server.

* `member_instance_id` - (Optional, String) Specifies the backend server ID.

* `response_timeout` - (Optional, Int) Specifies the response timeout for the listener.

* `protection_status` - (Optional, String) Specifies the protection status.

* `proxy_protocol_enable` - (Optional, String) Specifies whether the proxy protocol option to pass the source IP addresses
  of the clients to backend servers is enabled. Value options: **true**, **false**.

* `ssl_early_data_enable` - (Optional, String) Specifies whether the 0-RTT capability is enabled. Value options: **true**,
  **false**.

* `tls_ciphers_policy` - (Optional, String) Specifies the TLS cipher policy for the listener.

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

* `forward_elb` - Whether to transfer the load balancer ID to backend servers through the HTTP header of the packet.

* `forward_eip` - Whether to transparently transmit the load balancer EIP to backend servers.

* `forward_port` - Whether to transparently transmit the listening port of the load balancer to backend servers.

* `forward_request_port` - Whether to transparently transmit the source port of the client to backend servers.

* `forward_host` - Whether to rewrite the X-Forwarded-Host header.

* `forward_proto` - Whether to transfer the listener protocol of the load balancer to backend servers through the HTTP
  header of the packet.

* `forward_tls_certificate` - Whether to transfer the certificate ID of the load balancer to backend servers through the
  HTTP header of the packet.

* `forward_tls_cipher` - Whether to transfer the algorithm suite of the load balancer to backend servers through the HTTP
  header of the packet.

* `forward_tls_protocol` - Whether to transfer the algorithm protocol of the load balancer to backend servers through the
  HTTP header of the packet.

* `real_ip` - Whether to transfer the source IP address of the client to backend servers through the HTTP header of the
  packet.

* `ipgroup` - The IP address group associated with the listener.
  The [ipgroup](#ipgroup_stuct) structure is documented below.

* `port_ranges` - The port range, including the start and end port numbers.
  The [port_ranges](#port_ranges_stuct) structure is documented below.

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

* `proxy_protocol_enable` - Whether to enable the proxy protocol option to pass the source IP addresses of the clients
  to backend servers.

* `quic_config` - The QUIC configuration for the current listener.
  The [quic_config](#quic_config_stuct) structure is documented below.

* `security_policy_id` - The ID of the custom security policy.

* `sni_match_algo` - How wildcard domain name matches with the SNI certificates used by the listener.

* `ssl_early_data_enable` - Whether the 0-RTT capability is enabled.

* `max_connection` - The maximum number of concurrent connections that a listener can handle per second.

* `cps` - The maximum number of new connections that a listener can handle per second.

* `enable_member_retry` - Whether the health check retries for backend servers is enabled.

* `enterprise_project_id` - The ID of the enterprise project.

* `gzip_enable` - Whether the gzip compression for a load balancer is enabled.

* `tags` - The key/value pairs to associate with the listener.

* `created_at` - The creation time of the listener.

* `updated_at` - The update time of the listener.

<a name="ipgroup_stuct"></a>
The `ipgroup` block supports:

* `ipgroup_id` - The ID of the IP address group associated with the listener.

* `enable_ipgroup` - Whether access control is enabled.

* `type` - How access to the listener is controlled.

<a name="port_ranges_stuct"></a>
The `port_ranges` block supports:

* `start_port` - The start port.

* `end_port` - The end port.

<a name="quic_config_stuct"></a>
The `quic_config` block supports:

* `quic_listener_id` - The ID of the QUIC listener.

* `enable_quic_upgrade` - Whether to enable QUIC upgrade.

---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_listener_copy"
description: |-
  Manages a listener copy resource within HuaweiCloud.
---

# huaweicloud_elb_listener_copy

Manages a listener copy resource within HuaweiCloud.

-> Before copying listener, ensure that the following conditions are met:
   <br/>1. You can only copy a listener from a load balancer to another in the same VPC.
   <br/>2. Listeners of gateway load balancers cannot be copied, and listeners of other types of load balancers
   cannot be copied to gateway load balancers.
   <br/>3. You can only copy a listener from a load balancer to another of the same type.
   <br/>4. The original listener cannot have more than 1,000 backend servers and 100 forwarding policies.
   <br/>5. The source and destination load balancers cannot be frozen or migrated.
   <br/>6. If quic config of the original listener is configured, quic config of the new listener will be null.
   <br/>7. If a redirection is configured for the original listener, this forwarding policy will not be copied.

## Example Usage

```hcl
variable "listener_id" {}
variable "loadbalancer_id" {}
variable "protocol_port" {}

resource "huaweicloud_elb_listener_copy" "test" {
  listener_id     = var.listener_id
  loadbalancer_id = var.loadbalancer_id
  protocol_port   = var.protocol_port
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new listener.

* `listener_id` - (Required, String, NonUpdatable) Specifies the ID of the listener to be copied.

* `loadbalancer_id` - (Required, String, NonUpdatable) Specifies the ID of the load balancer to which
  the new listener is copied.

  -> 1. Listeners cannot be copied to gateway load balancers.
     <br/>2. The destination load balancer must support the protocol of the original listener. If the original
     listener uses HTTP or HTTPS, the destination load balancer must be an application load balancer. If the
     original listener uses TCP, UDP, or TLS, the destination load balancer must be a network load balancer.
     <br/>3. Listeners of shared load balancers can be copied to shared load balancers, and those of dedicated
     load balancers can be copied to dedicated load balancers.
     <br/>4. If IP as a Backend is enabled for the source load balancer, it must also be enabled for the
     destination load balancer.
     <br/>5. If the original listener uses TLS, the destination load balancer must support TLS listeners.

* `name` - (Optional, String, NonUpdatable) Specifies the name of the new listener.
  The value contains `0` to `255` characters. If not specified, the default value is used.
  The default value is original listener name + **-copy**.

* `protocol_port` - (Optional, Int, NonUpdatable) Specifies the port used by the new listener.

  -> 1. The port cannot be the same as that of any existing listener on the destination load balancer.
     <br/>2. The port of HTTP and TERMINATED_HTTPS listeners added to a shared load balancer cannot be `21`.
     <br/>3. If this parameter is set to `0`, `port_ranges` is required.
     <br/>4. The valid value ranges from `0` to `65,535`.

* `port_ranges` - (Optional, List, NonUpdatable) Specifies the port monitoring range (closed range).
  The [port_ranges](#port_ranges_struct) structure is documented below.

  -> 1. This parameter can be specified only when `protocol_port` is set to `0` or `protocol_port` is not specified.
     <br/>2. This parameter is available for **TCP**, **UDP**, and **TLS** listeners.
     <br/>3. The port cannot conflict with that of any existing listener on the destination load balancer.
     <br/>4. A maximum of `10` port ranges can be specified. Port ranges cannot overlap with each other.

* `reuse_pool` - (Optional, Bool, NonUpdatable) Specifies whether to reuse or copy the backend server groups and
  backend servers of the original listener.
  The valid values are as follows:
  + **true**: Indicates the backend server groups of the original listener will be used.
  + **false**: Indicates new backend server groups with the same settings as those of the original listener will be
  created and associated with the destination load balancer.

* `force_delete` - (Optional, Bool) Specifies whether to forcibly delete the listener, remove the listener,
  l7 policies, unbind associated pools. Defaults to **false**.
  This parameter only used when you delete the listener.

<a name="port_ranges_struct"></a>
The `port_ranges` block supports:

* `start_port` - (Required, Int, NonUpdatable) Specifies the start port number.

* `end_port` - (Required, Int, NonUpdatable) Specifies the end port number.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `client_ca_tls_container_ref` - The ID of the CA certificate used by the listener.

* `connection_limit` - The maximum number of connections that the listener can establish with backend servers.

* `created_at` - The creation time of the listener.

* `default_pool_id` - The ID of the default backend server group of the listener.

* `default_tls_container_ref` - The ID of the server certificate used by the listener.

* `description` - The description of the listener.

* `http2_enable` - Whether to allow clients to use HTTP/2 to communicate with the load balancer for higher
  access performance.

* `insert_headers` - The HTTP headers that can transmit required information to backend servers.
  The [insert_headers](#insert_headers_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID.

* `protocol` - The protocol of the listener.

* `sni_container_refs` - The IDs of SNI certificates used by the listener.

* `sni_match_algo` - How wildcard domain name matches with the SNI certificates used by the listener.

* `tags` - The tags of the listener.

* `updated_at` - The update time of the listener.

* `tls_ciphers_policy` - The security policy used by the listener.

* `security_policy_id` - The ID of the custom security policy.

* `enable_member_retry` - Whether to enable health check retries for backend servers.

* `keepalive_timeout` - The idle timeout duration, in seconds.

* `client_timeout` - The timeout duration for waiting for a response from a client, in seconds.

* `member_timeout` - The timeout duration for waiting for a response from a backend server, in seconds.

* `ipgroup` - The IP address group associated with the listener.
  The [ipgroup](#ipgroup_struct) structure is documented below.

* `transparent_client_ip_enable` - Whether to pass source IP addresses of the clients to backend servers.

* `proxy_protocol_enable` - Whether to enable ProxyProtocol to pass the source IP addresses of the clients to
  backend servers.

* `enhance_l7policy_enable` - Whether to enable advanced forwarding.

* `quic_config` - The QUIC configuration for the current listener.
  The [quic_config](#quic_config_struct) structure is documented below.

* `protection_status` - The protection status.

* `protection_reason` - The reason for setting up protection.

* `gzip_enable` - Whether to enable gzip for a load balancer.

* `ssl_early_data_enable` - Whether to enable zero round trip time resumption (0-RTT) for listeners.

* `cps` - The maximum number of new connections that a listener can handle per second.

* `connection_attr` - The maximum number of concurrent connections that a listener can handle per second.

* `nat64_enable` - Whether to translate between IPv4 and IPv6 addresses.

<a name="insert_headers_struct"></a>
The `insert_headers` block supports:

* `forwarded_elb_ip` - The load balancer EIP can be stored in the X-Forwarded-ELB-IP header and passed to
  backend servers.

* `forwarded_port` - The listening port of the load balancer can be stored in the X-Forwarded-Port header
  and passed to backend servers.

* `forwarded_for_port` - The source port of the client can be stored in the HTTP header and passed to backend servers.

* `forwarded_host` - The client host can be stored in the X-Forwarded-Host header and passed to backend servers.

* `forwarded_proto` - The listener protocol of the load balancer can be stored in the X-Forwarded-Proto header and
  passed to backend servers.

* `real_ip` - The client IP address can be stored in the X-Real-IP header and passed to backend servers.

* `forwarded_elb_id` - The load balancer ID can be stored in the X-Forwarded-ELB-ID header and passed to backend
  servers.

* `forwarded_tls_certificate_id` - The certificate ID of the load balancer can be stored in the
  X-Forwarded-TLS-Certificate-ID header and passed to backend servers.

* `forwarded_tls_protocol` - The algorithm protocol of the load balancer can be stored in the X-Forwarded-TLS-Protocol
  header and passed to backend servers.

* `forwarded_tls_cipher` - The algorithm suite of the load balancer can be stored in the X-Forwarded-TLS-Cipher header
  and passed to backend servers.

* `forwarded_tls_protocol_alias` - The name of the X-Forwarded-TLS-Protocol header.

* `forwarded_tls_cipher_alias` - The name of the X-Forwarded-TLS-Cipher header.

* `forwarded_for_processing_mode` - The X-Forwarded-For header handle mode.

* `forwarded_clientcert_subjectdn_enable` - Whether to transfer the owner information of the client certificate that
  accesses the load balancer through the X-Forwarded-Clientcert-subjectdn header.

* `forwarded_clientcert_subjectdn_alias` - The name of the X-Forwarded-Clientcert-subjectdn header.

* `forwarded_clientcert_issuerdn_enable` - Whether to transfer the issuer information of the client certificate that
  accesses the load balancer through the X-Forwarded-Clientcert-issuerdn header.

* `forwarded_clientcert_issuerdn_alias` - The name of the X-Forwarded-Clientcert-issuerdn header.

* `forwarded_clientcert_fingerprint_enable` - Whether to transfer the fingerprint of the client certificate that
  accesses the load balancer through the X-Forwarded-Clientcert-fingerprint header.

* `forwarded_clientcert_fingerprint_alias` - The name of the X-Forwarded-Clientcert-fingerprint header.

* `forwarded_clientcert_clientverify_enable` - Whether to transfer the verification result of the client certificate
  that accesses the load balancer through the X-Forwarded-Clientcert-clientverify header.

* `forwarded_clientcert_clientverify_alias` - The name of the X-Forwarded-Clientcert-clientverify header.

* `forwarded_clientcert_serialnumber_enable` - Whether to transfer the client certificate serial number through the
  X-Forwarded-Clientcert-serialnumber header.

* `forwarded_clientcert_serialnumber_alias` - The name of the X-Forwarded-Clientcert-serialnumber header.

* `forwarded_clientcert_enable` - Whether to transfer the client certificate content through the X-Forwarded-Clientcert
  header.

* `forwarded_clientcert_alias` - The name of the X-Forwarded-Clientcert header.

* `forwarded_clientcert_ciphers_enable` - Whether to transfer the TLS cipher suite supported by the client through the
  X-Forwarded-Clientcert-ciphers header.

* `forwarded_clientcert_ciphers_alias` - The name of the X-Forwarded-Clientcert-ciphers header.

* `forwarded_clientcert_end_enable` - Whether to transfer the client certificate expiration date through the
  X-Forwarded-Clientcert-end header.

* `forwarded_clientcert_end_alias` - The name of the X-Forwarded-Clientcert-end header.

* `forwarded_tls_alpn_protocol_enable` - Whether to use the X-Forwarded-Tls-Alpn-Protocol header to transfer
  the application protocol that the client and the load balancer negotiate to use during SSL handshakes.

* `forwarded_tls_alpn_protocol_alias` - The name of the X-Forwarded-Tls-Alpn-Protocol header.

* `forwarded_tls_sni_enable` - Whether to transfer the domain name of the SNI certificate accessed by the client
  through the X-Forwarded-Tls-Sni header.

* `forwarded_tls_sni_alias` - The name of the X-Forwarded-Tls-Sni header.

* `forwarded_tls_ja3_enable` - Whether to transfer the Ja3 fingerprint of the client that accesses the load balancer
  through the X-Forwarded-Tls-Ja3 header.

* `forwarded_tls_ja3_alias` - The name of the X-Forwarded-Tls-Ja3 header.

* `forwarded_tls_ja4_enable` - Whether to transfer the Ja4 fingerprint of the client that accesses the load balancer
  through the X-Forwarded-Tls-Ja4 header.

* `forwarded_tls_ja4_alias` - The name of the X-Forwarded-Tls-Ja4 header.

<a name="ipgroup_struct"></a>
The `ipgroup` block supports:

* `ipgroup_id` - The ID of the IP address group associated with the listener.

* `enable_ipgroup` - Whether access control is enabled.

* `type` - The IP address group type.

<a name="quic_config_struct"></a>
The `quic_config` block supports:

* `quic_listener_id` - The ID of the QUIC listener.

* `enable_quic_upgrade` - Whether to enable QUIC upgrade.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_elb_listener_copy.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `listener_id`, `reuse_pool`,
`force_delete`. It is generally recommended running `terraform plan` after importing a listener.
You can then decide if changes should be applied to the listener, or the resource definition should be updated to align
with the listener. Also you can ignore changes as below.

```hcl
resource "huaweicloud_elb_listener" "test" {
    ...
  lifecycle {
    ignore_changes = [
      listener_id,
      reuse_pool,
      force_delete,
    ]
  }
}
```

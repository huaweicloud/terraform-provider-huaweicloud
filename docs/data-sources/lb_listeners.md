---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_listeners"
description: |-
  Use this data source to query the list of ELB listeners.
---

# huaweicloud_lb_listeners

Use this data source to query the list of ELB listeners.

## Example Usage

```hcl
variable "protocol" {}

data "huaweicloud_lb_listeners" "test" {
  protocol  = var.protocol
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The listener name.

* `protocol` - (Optional, String) The listener protocol.  
  The valid values are **TCP**, **UDP**, **HTTP** and **TERMINATED_HTTPS**.

* `protocol_port` - (Optional, String) The front-end listening port of the listener.  
  The valid value is range from `1` to `65535`.

* `client_ca_tls_container_ref` - (Optional, String) The ID of the CA certificate used by the listener.

* `default_pool_id` - (Optional, String) The ID of the default pool with which the listener is associated.

* `default_tls_container_ref` - (Optional, String) The ID of the server certificate used by the listener.

* `description` - (Optional, String) The description for the listener.

* `enterprise_project_id` - (Optional, String) The ID of the enterprise project.

* `http2_enable` - (Optional, String) Whether the ELB listener uses HTTP/2. Value options: **true**, **false**.

* `listener_id` - (Optional, String) ID of the listener.

* `loadbalancer_id` - (Optional, String) The ID of the load balancer that the listener is added to.

* `tls_ciphers_policy` - (Optional, String) The security policy used by the listener.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `listeners` - Listener list.
The [listeners](#listeners_struct) structure is documented below.

<a name="listeners_struct"></a>
The `listeners` block supports:

* `id` - The ELB listener ID.

* `name` - The listener name.

* `protocol` - The listener protocol.

* `protocol_port` - The front-end listening port of the listener.

* `default_pool_id` - The ID of the default pool with which the ELB listener is associated.

* `client_ca_tls_container_ref` - The ID of the CA certificate used by the listener.

* `description` - The description of the ELB listener.

* `connection_limit` - The maximum number of connections allowed for the listener.

* `http2_enable` - Whether the ELB listener uses HTTP/2.

* `insert_headers` - Whether to insert HTTP extension headers and sent them to backend servers.
  The [insert_headers](#insert_headers_struct) structure is documented below.

* `default_tls_container_ref` - The ID of the server certificate used by the listener.

* `sni_container_refs` - List of the SNI certificate (server certificates with a domain name) IDs used by the listener.

* `protection_status` - Whether modification protection is enabled.

* `protection_reason` - The reason to enable modification protection.

* `tls_ciphers_policy` - security policy used by the listener.

* `loadbalancers` - The list of the associated load balancer.
  The [loadbalancers](#loadbalancers_struct) structure is documented below.

* `tags` - The key/value pairs to associate with the listener.

* `created_at` - The time when the listener was created.

* `updated_at` - The time when the listener was updated.

<a name="loadbalancers_struct"></a>
The `loadbalancers` block supports:

* `id` - The ELB loadbalancer ID.

<a name="insert_headers_struct"></a>
The `insert_headers` block supports:

* `x_forwarded_elb_ip` - Whether to transparently transmit the load balancer EIP to backend servers.

* `x_forwarded_host` - Whether to rewrite the X-Forwarded-Host header.

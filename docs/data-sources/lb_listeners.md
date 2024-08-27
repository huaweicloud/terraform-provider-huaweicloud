---
subcategory: "Elastic Load Balance (ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lb_listeners"
description: ""
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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `listeners` - Listener list.
The [object](#listeners_listeners) structure is documented below.

<a name="listeners_listeners"></a>
The `listeners` block supports:

* `id` - The ELB listener ID.

* `name` - The listener name.

* `protocol` - The listener protocol.

* `protocol_port` - The front-end listening port of the listener.

* `default_pool_id` - The ID of the default pool with which the ELB listener is associated.

* `description` - The description of the ELB listener.

* `connection_limit` - The maximum number of connections allowed for the listener.

* `http2_enable` - Whether the ELB listener uses HTTP/2.

* `default_tls_container_ref` - The ID of the server certificate used by the listener.

* `sni_container_refs` - List of the SNI certificate (server certificates with a domain name) IDs used by the listener.

* `loadbalancers` - Listener list.
The [object](#listeners_loadbalancers) structure is documented below.

<a name="listeners_loadbalancers"></a>
The `loadbalancers` block supports:

* `id` - The ELB loadbalancer ID.

---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_es_loadbalancer_config"
description: |-
  Manages CSS ElasticSearch loadbalancer resource within HuaweiCloud.
---

# huaweicloud_css_es_loadbalancer_config

Manages CSS ElasticSearch loadbalancer resource within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "agency" {}
variable "elb_loadbalancer_id" {}
variable "protocol_port" {}
variable "server_cert_id" {}

resource "huaweicloud_css_es_loadbalancer_config" "test" {
  cluster_id      = var.cluster_id
  agency          = var.agency
  loadbalancer_id = var.elb_loadbalancer_id
  protocol_port   = var.protocol_port
  server_cert_id  = var.server_cert_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `cluster_id` - (Required, String, ForceNew) Specifies the ID of the CSS cluster.
  Changing this creates a new resource.

* `agency` - (Required, String, ForceNew) Specifies the IAM agency used to access ELB.
  Changing this creates a new resource.

* `loadbalancer_id` - (Required, String, ForceNew) Specifies the ID of the loadbalancer.
  Changing this creates a new resource.

* `protocol_port` - (Required, Int, ForceNew) Specifies the front-end listening port of the listener.
  Changing this creates a new resource.

* `server_cert_id` - (Optional, String) Specifies the server certificate ID used by the ELB listener.

* `ca_cert_id` - (Optional, String) Specifies the CA certificate ID used by the ELB listener.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `server_cert_name` - The server certificate name.

* `ca_cert_name` - The CA certificate name.

* `elb_enabled` - Whether the loadbalancer is enabled.

* `authentication_type` - The authentication type.

* `loadbalancer` - The ELB loadbalancer information.
  The [loadbalancer](#Css_elb_loadbalancer) structure is documented below.

* `listener` - The listener information.
  The [listener](#Css_elb_listener) structure is documented below.

* `health_monitors` - The health monitors.
  The [health_monitors](#Css_elb_health_monitors) structure is documented below.

<a name="Css_elb_loadbalancer"></a>
The `loadbalancer` block supports:

* `id` - The loadbalancer ID.

* `name` - The loadbalancer name.

* `ip` - The IPv4 virtual IP address of the loadbalancer.

* `public_ip` - The elastic public IP address.

<a name="Css_elb_listener"></a>
The `listener` block supports:

* `id` - The listener ID.

* `name` - The listener name.

* `protocol` - The listening protocol of the listener.

* `protocol_port` - The front-end listening port of the listener.

* `ip_group` - The ipgroup information in the listener object.
  The [ip_group](#Listener_ip_group) structure is documented below.

<a name="Listener_ip_group"></a>
The `ip_group` block supports:

* `id` - The ID of the access control group associated with the listener.

* `enabled` - The status of the access control group.

<a name="Css_elb_health_monitors"></a>
The `health_monitors` block supports:

* `ip` - The IP address corresponding to the backend server.

* `protocol_port` - The backend server business port number.

* `status` - The health status of the backend cloud server.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 3 minutes.

* `delete` - Default is 3 minutes.

## Import

The CSS ElasticSearch loadbalancer config can be imported using `cluster_id`, e.g.

```bash
$ terraform import huaweicloud_css_es_loadbalancer_config.test <cluster_id>
```

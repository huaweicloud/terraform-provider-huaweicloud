---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_service"
description: |-
  Manages a CCI Service resource within HuaweiCloud.
---

# huaweicloud_cciv2_service

Manages a CCI Service resource within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}
variable "name" {}
variable "elb_id" {}

resource "huaweicloud_cciv2_service" "test" {
  namespace = var.namespace
  name      = var.name

  annotations = {
    "kubernetes.io/elb.class" = "elb",
    "kubernetes.io/elb.id"    = var.elb_id,
  }

  ports {
    name         = "service-example-port"
    app_protocol = "TCP"
    protocol     = "TCP"
    port         = 87
    target_port  = 65529
  }

  selector = {
    app = "test2"
  }

  type = "LoadBalancer"

  lifecycle {
    ignore_changes = [
      annotations,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the name of the CCI Service.

* `namespace` - (Required, String, NonUpdatable) Specifies the namespace of the CCI Service.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI Service.

* `labels` - (Optional, Map) Specifies the labels of the CCI Service.

* `ports` - (Optional, List) Specifies the ports of the CCI Service.
  The [ports](#service_ports) structure is documented below.

* `selector` - (Optional, Map) Specifies the selector of the CCI Service.

* `type` - (Optional, String) The type of the CCI Service.

<a name="service_ports"></a>
The `ports` block supports:

* `port` - (Required, Int) The port.

* `app_protocol` - (Optional, String) The app protocol.

* `name` - (Optional, String) The name.

* `protocol` - (Optional, String) The protocol.

* `target_port` - (Optional, Int) The target port.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI Service.

* `cluster_ip` - Specifies the cluster IP of the CCI Service.

* `cluster_ips` - Specifies the cluster IPs of the CCI Service.

* `creation_timestamp` - The creation timestamp of the namespace.

* `external_name` - The external name of the CCI Service.

* `finalizers` - The finalizers of the namespace.

* `ip_families` - The IP families of the CCI Service.

* `ip_family_policy` - The IP family policy of the CCI Service.

* `kind` - The kind of the CCI Service.

* `load_balancer_ip` - The load balancer IP of the CCI Service.

* `publish_not_ready_addresses` - Whether the publish is not ready addresses of the CCI Service.

* `resource_version` - The resource version of the namespace.

* `session_affinity` - The session affinity of the CCI Service.

* `status` - The status of the namespace.
  The [status](#service_status) structure is documented below.

* `uid` - The uid of the namespace.

<a name="service_status"></a>
The `status` block supports:

* `conditions` - Tthe conditions of the CCI Service.
  The [conditions](#service_status_conditions) structure is documented below.

* `loadbalancer` - The loadbalancer of the CCI Service.
  The [loadbalancer](#service_status_loadbalancer) structure is documented below.

<a name="service_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time.

* `message` - The message.

* `observe_generation` - The observe generation.

* `reason` - The reason.

* `status` - Tthe status.

* `type` - The type.

<a name="service_status_loadbalancer"></a>
The `loadbalancer` block supports:

* `ingress` - The ingress of the loadbalancer.
  The [ingress](#service_status_loadbalancer_ingress) structure is documented below.

<a name="service_status_loadbalancer_ingress"></a>
The `ingress` block supports:

* `ip` - The IP of the loadbalancer.

* `ports` - The ports of the loadbalancer.
  The [ports](#service_status_loadbalancer_ingress_ports) structure is documented below.

<a name="service_status_loadbalancer_ingress_ports"></a>
The `ports` block supports:

* `error` - The error.

* `port` - The port.

* `protocol` - The protocol.

## Import

The CCI Service can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_service.test <namespace>/<name>
```

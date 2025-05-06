---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_services"
description: |-
  Use this data source to get the list of CCI Services within HuaweiCloud.
---
# huaweicloud_cciv2_services

Use this data source to get the list of CCI Services within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_services" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the CCI service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `services` - The CCI services list.
  The [services](#services) structure is documented below.

<a name="services"></a>
The `services` block supports:

* `annotations` - The annotations.

* `creation_timestamp` - The creation time.

* `finalizers` - The finalizers.

* `labels` - The labels.

* `name` - The name.

* `namespace` - The namespace.

* `ports` - The ports.
  The [ports](#services_ports) structure is documented below.

* `resource_version` - The resource version.

* `selector` - The selector.

* `session_affinity` - The session affinity.

* `status` - The status.
  The [status](#services_status) structure is documented below.

* `type` - The type.

* `uid` - The uid.

<a name="services_ports"></a>
The `ports` block supports:

* `app_protocol` - The app protocol.

* `name` - The name.

* `port` - The port.

* `protocol` - The protocol.

* `target_port` - The target port.

<a name="services_status"></a>
The `status` block supports:

* `conditions` - The conditions.
  The [conditions](#services_status_conditions) structure is documented below.

* `loadbalancer` - The loadbalancer.
  The [loadbalancer](#services_status_loadbalancer) structure is documented below.

<a name="services_status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time.

* `message` - The message.

* `observe_generation` - The observe generation.

* `reason` - The reason.

* `status` - The status.

* `type` - The type.

<a name="services_status_loadbalancer"></a>
The `loadbalancer` block supports:

* `ingress` - The ingress.
  The [ingress](#services_status_loadbalancer_ingress) structure is documented below.

<a name="services_status_loadbalancer_ingress"></a>
The `ingress` block supports:

* `ip` - The IP.

* `ports` - The ports.
  The [ports](#services_status_loadbalancer_ingress_ports) structure is documented below.

<a name="services_status_loadbalancer_ingress_ports"></a>
The `ports` block supports:

* `error` - The error.

* `port` - The port.

* `protocol` - The protocol.

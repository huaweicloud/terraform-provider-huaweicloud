---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_service"
description: |-
  Manages a CCI Service resource within HuaweiCloud.
---

# huaweicloud_cciv2_service

<!--
please add the description of huaweicloud_cciv2_service
  + For resource: Manages xxx resource within HuaweiCloud.
  + For data source: Use this data source to get the list of xxx.
-->

## Example Usage

<!-- please add the usage of huaweicloud_cciv2_service -->
```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the CCI Service.

* `namespace` - (Required, String) Specifies the namespace.

* `annotations` - (Optional, Map) Specifies the annotations of the CCI Service.

* `enable_force_new` - (Optional, String) <!-- please add the description of the argument -->

* `labels` - (Optional, Map) Specifies the annotations of the CCI Service.

* `selector` - (Optional, Map) Specifies the selector of the CCI Service.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `api_version` - The API version of the CCI Service.

* `client_timeout_seconds` - Specifies the cluster IPs of the CCI Service.

* `cluster_ip` - Specifies the cluster IP of the CCI Service.

* `cluster_ips` - Specifies the cluster IPs of the CCI Service.

* `creation_timestamp` - The creation timestamp of the namespace.

* `external_name` - The external name of the CCI Service.

* `ip_families` - The IP families of the CCI Service.

* `ip_family_policy` - The IP family policy of the CCI Service.

* `kind` - The kind of the CCI Service.

* `load_balancer_ip` - The load balancer IP of the CCI Service.

* `ports` - Specifies the ports of the CCI Service.
  The [ports](#attrblock--ports) structure is documented below.

* `publish_not_ready_addresses` - Whether the publish is not ready addresses of the CCI Service.

* `resource_version` - The resource version of the namespace.

* `session_affinity` - The load balancer IP of the CCI Service.

* `status` - The status of the namespace.
  The [status](#attrblock--status) structure is documented below.

* `type` - The type of the CCI Service.

* `uid` - The uid of the namespace.

<a name="attrblock--ports"></a>
The `ports` block supports:

* `app_protocol` - The app protocol.

* `name` - Tthe name.

* `port` - The port.

* `protocol` - The protocol.

* `target_port` - The target port.

<a name="attrblock--status"></a>
The `status` block supports:

* `conditions` - Tthe conditions of the CCI Service.
  The [conditions](#attrblock--status--conditions) structure is documented below.

* `loadbalancer` - The subnet attributes of the CCI Service.
  The [loadbalancer](#attrblock--status--loadbalancer) structure is documented below.

<a name="attrblock--status--conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time.

* `message` - The message.

* `observe_generation` - The observe generation.

* `reason` - The reason.

* `status` - Tthe status.

* `type` - The type.

<a name="attrblock--status--loadbalancer"></a>
The `loadbalancer` block supports:

* `ingress` - The ID of the CCI Service.
  The [ingress](#attrblock--status--loadbalancer--ingress) structure is documented below.

<a name="attrblock--status--loadbalancer--ingress"></a>
The `ingress` block supports:

* `ip` - The ID of the CCI Service.

* `ports` - The ports.
  The [ports](#attrblock--status--loadbalancer--ingress--ports) structure is documented below.

<a name="attrblock--status--loadbalancer--ingress--ports"></a>
The `ports` block supports:

* `error` - The error.

* `port` - The port.

* `protocol` - The protocol.

## Import

The CCI Service can be imported using `namespace` and `name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_cciv2_service.test <namespace>/<name>
```

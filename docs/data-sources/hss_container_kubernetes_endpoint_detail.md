---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_endpoint_detail"
description: |-
  Use this data source to get HSS container kubernetes endpoint detail within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_endpoint_detail

Use this data source to get HSS container kubernetes endpoint detail within HuaweiCloud.

## Example Usage

```hcl
variable "endpoint_id" {}

data "huaweicloud_hss_container_kubernetes_endpoint_detail" "test" {
  endpoint_id = var.endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS container kubernetes endpoints.
  If omitted, the provider-level region will be used.

* `endpoint_id` - (Required, String) Specifies the endpoint ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID same as `endpoint_id`.

* `name` - The endpoint name.

* `service_name` - The service name.

* `namespace` - The namespace.

* `creation_timestamp` - Create timestamp.

* `cluster_name` - The cluster name.

* `labels` - The labels.

* `association_service` - Is it associated with a service.

* `endpoint_pod_list` - The endpoint associated pod mapping.
  The [endpoint_pod_list](#endpoint_pod_list_struct) structure is documented below.

* `endpoint_port_list` - The endpoint associated port list.
  The [endpoint_port_list](#endpoint_port_list_struct) structure is documented below.

<a name="endpoint_pod_list_struct"></a>
The `endpoint_pod_list` block supports:

* `id` - The ID.

* `endpoint_id` - The associate endpoint ID.

* `pod_ip` - The pod ip.

* `pod_name` - The pod name.

* `available` - Is it available.

<a name="endpoint_port_list_struct"></a>
The `endpoint_port_list` block supports:

* `id` - The ID.

* `endpoint_id` - The associate endpoint ID.

* `name` - The port name.

* `protocol` - The service agreement.

* `port` - The port number.

* `app_protocol` - The application protocol.

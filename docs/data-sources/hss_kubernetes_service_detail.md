---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_service_detail"
description: |-
  Use this data source to get the HSS kubernetes service detail within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_service_detail

Use this data source to get the HSS kubernetes service detail within HuaweiCloud.

## Example Usage

```hcl
variable "service_id" {}

data "huaweicloud_hss_kubernetes_service_detail" "test" {
  service_id = var.service_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `service_id` - (Required, String) Specifies the service ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `name` - The service name.

* `endpoint_name` - The endpoint name.

* `namespace` - The namespace.

* `creation_timestamp` - The creation timestamp.

* `cluster_name` - The cluster name.

* `labels` - The labels.

* `type` - The service type. Valid values are:
  + **ClusterIP**: Services that are only accessible within the cluster.
  + **NodePort**: Services that are exposed through NodePort.
  + **LoadBalancer**: Services that are exposed through LoadBalancer.

* `cluster_ip` - The cluster IP.

* `selector` - The selector.

* `session_affinity` - The session affinity.

* `service_port_list` - The list of service ports.

  The [service_port_list](#service_port_list_struct) structure is documented below.

<a name="service_port_list_struct"></a>
The `service_port_list` block supports:

* `id` - The ID.

* `service_id` - The associated service ID.

* `name` - The port name.

* `protocol` - The service protocol.

* `port` - The port number.

* `target_port` - The backend container port.

* `node_port` - The node port.

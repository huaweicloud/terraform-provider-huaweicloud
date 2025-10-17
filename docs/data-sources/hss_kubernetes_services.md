---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_services"
description: |-
  Use this data source to get the list of HSS kubernetes services within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_services

Use this data source to get the list of HSS kubernetes services within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_kubernetes_services" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `name` - (Optional, String) Specifies the service name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `namespace` - (Optional, String) Specifies the namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `service_info_list` - The list of kubernetes services.

  The [service_info_list](#service_info_list_struct) structure is documented below.

<a name="service_info_list_struct"></a>
The `service_info_list` block supports:

* `id` - The service ID.

* `name` - The service name.

* `endpoint_name` - The endpoint name.

* `namespace` - The namespace.

* `creation_timestamp` - The creation timestamp.

* `type` - The service type. Valid values are:
  + **ClusterIP**: Services that are only accessible within the cluster.
  + **NodePort**: Services that are exposed through NodePort.
  + **LoadBalancer**: Services that are exposed through LoadBalancer.

* `cluster_ip` - The cluster IP.

* `cluster_name` - The cluster name.

* `cluster_type` - The cluster type. Valid values are:
  + **k8s**: Native cluster.
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

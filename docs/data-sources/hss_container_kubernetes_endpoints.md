---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes_endpoints"
description: |-
  Use this data source to get the list of HSS container kubernetes endpoints within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes_endpoints

Use this data source to get the list of HSS container kubernetes endpoints within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes_endpoints" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS container kubernetes endpoints.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the endpoint name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `namespace` - (Optional, String) Specifies the namespace.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of endpoints.

* `last_update_time` - The latest update time.

* `endpoint_info_list` - The list of endpoint information.
  The [endpoint_info_list](#endpoint_info_list_struct) structure is documented below.

<a name="endpoint_info_list_struct"></a>
The `endpoint_info_list` block supports:

* `id` - The ID.

* `name` - The endpoint name.

* `service_name` - The service name.

* `namespace` - The namespace.

* `creation_timestamp` - Create timestamp.

* `cluster_name` - The cluster name.

* `cluster_type` - The cluster type.  
  The valid values are as follows:
  + **k8s**: Native cluster.
  + **cce**: CCE cluster.
  + **ali**: Alibaba cloud cluster.
  + **tencent**: Tencent cloud cluster.
  + **azure**: Microsoft cloud cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: HuaweiCloud self built cluster.
  + **self_built_idc**: IDC self built cluster.

* `association_service` - Is it associated with a service.

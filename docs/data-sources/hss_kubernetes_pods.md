---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_pods"
description: |-
  Use this data source to get the list of HSS kubernetes pods within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_pods

Use this data source to get the list of HSS kubernetes pods within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_kubernetes_pods" "test" {}
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

* `pod_name` - (Optional, String) Specifies the pod name.

* `namespace_name` - (Optional, String) Specifies the namespace name.

* `cluster_name` - (Optional, String) Specifies the cluster name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `total_num` - The total number of pods.

* `data_list` - The list of pods.

  The [data_list](#data_list_struct) structure is documented below.

* `pod_info_list` - The pod basic information list.

  The [pod_info_list](#pod_info_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `pod_name` - The pod name.

* `namespace_name` - The namespace name.

* `cluster_name` - The cluster name.

* `node_name` - The node name.

* `cpu` - The CPU usage.

* `memory` - The memory usage.

* `cpu_limit` - The CPU limit.

* `memory_limit` - The memory limit.

* `node_ip` - The node IP.

* `pod_ip` - The pod IP.

* `status` - The pod status.  
  The valid values are as follows:
  + **Pending**: The pod has been accepted by the Kubernetes system, but one or more container images have not been
    created yet.
  + **Running**: The pod has been bound to a node and all containers have been created.
  + **Succeeded**: All containers in the pod have successfully terminated and will not restart.
  + **Failed**: All containers in the pod have terminated, and at least one container has terminated due to a failure.
  + **Unknown**: Due to some reason, the status of the pod cannot be obtained, usually due to an error in communication
    with the pod's host.

* `create_time` - The creation time.

* `region_id` - The region ID.

* `id` - The ID.

* `cluster_id` - The cluster ID.

* `cluster_type` - The cluster type.  
  The valid values are as follows:
  + **k8s**: Native cluster.
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

<a name="pod_info_list_struct"></a>
The `pod_info_list` block supports:

* `pod_name` - The pod name.

* `namespace_name` - The namespace name.

* `cluster_name` - The cluster name.

* `cpu` - The CPU usage.

* `memory` - The memory usage.

* `cpu_limit` - The CPU limit.

* `memory_limit` - The memory limit.

* `pod_ip` - The pod IP.

* `protect_status` - The protect status.  
  The valid values are as follows:
  + **closed**
  + **opened**
  + **protection_exception**

* `detect_result` - The serverless security detection results.  
  The valid values are as follows:
  + **undetected**: Not detected.
  + **clean**: No risk.
  + **risk**: At risk.
  + **scanning**: Under detection.

* `status` - The pod status.  
  The valid values are as follows:
  + **Pending**: The pod has been accepted by the Kubernetes system, but one or more container images have not been
    created yet.
  + **Running**: The pod has been bound to a node and all containers have been created.
  + **Succeeded**: All containers in the pod have successfully terminated and will not restart.
  + **Failed**: All containers in the pod have terminated, and at least one container has terminated due to a failure.
  + **Unknown**: Due to some reason, the status of the pod cannot be obtained, usually due to an error in communication
    with the pod's host.

* `create_time` - The creation time.

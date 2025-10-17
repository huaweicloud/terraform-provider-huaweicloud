---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_kubernetes_pod_detail"
description: |-
  Use this data source to get the kubernetes pod detail of HSS within HuaweiCloud.
---

# huaweicloud_hss_kubernetes_pod_detail

Use this data source to get the kubernetes pod detail of HSS within HuaweiCloud.

## Example Usage

```hcl
variable "pod_name" {}

data "huaweicloud_hss_kubernetes_pod_detail" "test" {
  pod_name = var.pod_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `pod_name` - (Required, String) Specifies the pod name.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID.
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If it is necessary to operate the asset under all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `pod_name_attr` - The pod name.

* `namespace_name` - The namespace name.

* `cluster_name` - The cluster name.

* `node_name` - The node name.

* `label` - The pod label.

* `cpu` - The CPU usage.

* `memory` - The memory usage.

* `cpu_limit` - The CPU limit.

* `memory_limit` - The memory limit.

* `node_ip` - The node IP address.

* `pod_ip` - The pod IP address.

* `status` - The pod status.
  The valid values are as follows:
  + **Pending**: The pod has been accepted by the Kubernetes system, but one or more container images have not been created.
  + **Running**: The pod has been bound to a node and all of the containers have been created.
  + **Succeeded**: All containers in the pod have voluntarily terminated with a container exit code of 0,
    and will not be restarted.
  + **Failed**: All containers in the pod have terminated, and at least one container has terminated in a failure state.
  + **Unknown**: For some reason the state of the pod could not be obtained.

* `create_time` - The creation timestamp.

* `containers` - The container list.

  The [containers](#containers_struct) structure is documented below.

<a name="containers_struct"></a>
The `containers` block supports:

* `id` - The container ID.

* `region_id` - The region ID.

* `container_id` - The container ID.

* `container_name` - The container name.

* `image_name` - The image name.

* `status` - The container status.
  The valid values are as follows:
  + **Running**: Running.
  + **Terminated**: Terminated.
  + **Waiting**: Waiting.

* `create_time` - The creation timestamp.

* `cpu_limit` - The CPU limit.

* `memory_limit` - The memory limit.

* `restart_count` - The restart count.

* `pod_name` - The pod name.

* `cluster_name` - The cluster name.

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

* `risky` - Whether there is a risk.
  The valid values are as follows:
  + **true**: There is a risk.
  + **false**: No risk.

* `low_risk` - The number of low-risk vulnerabilities.

* `medium_risk` - The number of medium-risk vulnerabilities.

* `high_risk` - The number of high-risk vulnerabilities.

* `fatal_risk` - The number of fatal-risk vulnerabilities.

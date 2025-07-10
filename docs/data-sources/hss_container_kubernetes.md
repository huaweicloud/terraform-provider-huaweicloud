---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_kubernetes"
description: |-
  Use this data source to get the list of HSS container information within HuaweiCloud.
---

# huaweicloud_hss_container_kubernetes

Use this data source to get the list of HSS container information within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_kubernetes" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the HSS container information.
  If omitted, the provider-level region will be used.

* `container_name` - (Optional, String) Specifies the container name.

* `pod_name` - (Optional, String) Specifies the pod name.

* `image_name` - (Optional, String) Specifies the image name.

* `cluster_container` - (Optional, Bool) Specifies whether it is a cluster-managed container.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of containers.

* `last_update_time` - The last update time.

* `data_list` - The list of container information.
  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `id` - The ID.

* `region_id` - The region ID.

* `container_id` - The container ID.

* `container_name` - The container name.

* `image_name` - The image name.

* `status` - The container status.  
  The valid values are as follows:
  + **Running**
  + **Terminated**
  + **Waiting**

* `create_time` - The creation time of the container.

* `cpu_limit` - The CPU limit.

* `memory_limit` - The memory limit.

* `restart_count` - The number of restarts.

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

* `low_risk` - The number of low-risk vulnerabilities.

* `medium_risk` - The number of medium-risk vulnerabilities.

* `high_risk` - The number of high-risk vulnerabilities.

* `fatal_risk` - The number of fatal-risk vulnerabilities.

---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_logs"
description: |-
  Use this data source to get the list of HSS container logs within HuaweiCloud.
---

# huaweicloud_hss_container_logs

Use this data source to get the list of HSS container logs within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_hss_container_logs" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

* `cluster_id` - (Optional, String) Specifies the cluster ID.

* `cluster_name` - (Optional, String) Specifies the cluster name.

* `namespace` - (Optional, String) Specifies the namespace to which the container generating the log belongs.

* `pod_name` - (Optional, String) Specifies the name of the pod to which the container generating the log belongs.

* `pod_id` - (Optional, String) Specifies the ID of the pod to which the container generating the log belongs.

* `pod_ip` - (Optional, String) Specifies the IP address of the pod to which the container generating the log belongs.

* `host_ip` - (Optional, String) Specifies the IP address of the host where the container generating the log is located.

* `container_id` - (Optional, String) Specifies the container ID.

* `container_name` - (Optional, String) Specifies the container name for generating logs.

* `content` - (Optional, String) Specifies the log content.

* `start_time` - (Optional, Int) Specifies the minimum time to query log range.

* `end_time` - (Optional, Int) Specifies the maximum time for querying log range.

* `line_num` - (Optional, String) Specifies the pagination line number that needs to be passed when querying the CCE
  cluster container log.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `total_num` - The total number of container logs.

* `data_list` - The list of container logs.

  The [data_list](#data_list_struct) structure is documented below.

<a name="data_list_struct"></a>
The `data_list` block supports:

* `cluster_id` - The cluster ID.

* `cluster_name` - The cluster name.

* `cluster_type` - The cluster type.  
  The valid values are as follows:
  + **cce**: CCE cluster.
  + **ali**: Alibaba Cloud cluster.
  + **tencent**: Tencent Cloud cluster.
  + **azure**: Microsoft Azure cluster.
  + **aws**: Amazon cluster.
  + **self_built_hw**: Huawei Cloud self-built cluster.
  + **self_built_idc**: IDC self-built cluster.

* `time` - The time when the log was generated.

* `namespace` - The namespace to which the container log belongs.

* `pod_name` - The name of the pod to which the container generating the log belongs.

* `pod_id` - The ID of the pod to which the container generating the log belongs.

* `pod_ip` - The IP address of the pod to which the container generating the log belongs.

* `host_ip` - The IP address of the host where the container generating the log is located.

* `container_name` - The container name for generating logs.

* `container_id` - The container ID for generating logs.

* `content` - The log content.

* `line_num` - The line number of the cce cluster container log.

---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_resource_pool_workloads"
description: |-
  Use this data source to query the workloads under a specified resource pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_resource_pool_workloads

Use this data source to query the workloads under a specified resource pool within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "resource_pool_id" {}

data "huaweicloud_modelartsv2_resource_pool_workloads" "test" {
  pool_id = var.resource_pool_id
}
```

### Filter by sequence

```hcl
variable "resource_pool_id" {}

data "huaweicloud_modelartsv2_resource_pool_workloads" "test" {
  pool_id = var.resource_pool_id
  sort    = "create_time"
  ascend  = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resource pool workloads are located.  
  If omitted, the provider-level region will be used.

* `pool_id` - (Required, String) Specifies the ID of the resource pool.

* `type` - (Optional, String) Specifies the type of the workload to be queried.  
  The valid values are as follows:
  + **train**: Training job.
  + **infer**: Inference service.
  + **notebook**: Notebook job.
  + **x-infer**: New version inference job.

* `status` - (Optional, String) Specifies the status of the workload to be queried.  
  The valid values are as follows:
  + **Queue**: Queued job.
  + **Pending**: Pending job.
  + **Abnormal**: Abnormal job.
  + **Terminating**: Terminating job.
  + **Creating**: Creating job.
  + **Running**: Running job.
  + **Completed**: Completed job.
  + **Terminated**: Terminated job.
  + **Failed**: Failed job.

* `sort` - (Optional, String) Specifies the sort field of the query result. The valid value is **create_time**.

* `ascend` - (Optional, Bool) Whether to sort the query results in ascending order. This parameter needs to be used
  in conjunction with `sort`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workloads` - The list of resource pool workloads that matched filter parameters.  
  The [workloads](#modelartsv2_resource_pool_workloads) structure is documented below.

<a name="modelartsv2_resource_pool_workloads"></a>
The `workloads` block supports:

* `api_version` - The API version of the workload.

* `kind` - The type of the workload resource.

* `type` - The type of the workload.

* `namespace` - The namespace of the workload.

* `name` - The name of the workload.

* `job_name` - The name of the job to which the workload belongs.

* `uid` - The ID of the workload.

* `job_uuid` - The ID of the job to which the workload belongs.

* `flavor` - The resource flavor of the workload.

* `status` - The status of the workload.

* `resource_requirement` - The resource requirement of the workload.  
  The [resource_requirement](#modelartsv2_resource_pool_workloads_resource_requirement) structure is documented below.

* `priority` - The priority of the workload.

* `running_duration` - The running duration of the workload, in seconds.

* `pending_duration` - The pending duration of the workload, in seconds.

* `pending_position` - The pending position of the workload.

* `create_time` - The creation time of the workload.

* `gvk` - The GVK (Group Version Kind) information of the workload.

* `host_ips` - The host IPs of the workload.

* `nodes` - The node information of the workload.  
  The [nodes](#modelartsv2_resource_pool_workloads_nodes) structure is documented below.

<a name="modelartsv2_resource_pool_workloads_resource_requirement"></a>
The `resource_requirement` block supports:

* `cpu` - The CPU resource of the workload.

* `memory` - The memory resource of the workload.

* `nvidia_gpu` - The GPU resource of the workload.

* `huawei_ascend_snt3` - The Ascend Snt3 NPU resource of the workload.

* `huawei_ascend_snt9` - The Ascend Snt9 NPU resource of the workload.

<a name="modelartsv2_resource_pool_workloads_nodes"></a>
The `nodes` block supports:

* `host_ip` - The host IP of the node.

* `npu_topology_placement` - The NPU topology placement information of the node.

* `resource_requirement` - The resource requirement of the node.  
  The [resource_requirement](#modelartsv2_resource_pool_workloads_nodes_resource_requirement) structure is documented below.

<a name="modelartsv2_resource_pool_workloads_nodes_resource_requirement"></a>
The `resource_requirement` block supports:

* `cpu` - The CPU resource of the node.

* `memory` - The memory resource of the node.

* `nvidia_gpu` - The GPU resource of the node.

* `huawei_ascend_snt310` - The Ascend Snt3 NPU resource of the node.

* `huawei_ascend_snt1980` - The Ascend Snt9 NPU resource of the node.

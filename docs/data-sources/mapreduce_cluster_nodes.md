---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_nodes"
description: |-
  Use this data source to query the node list under the specified cluster within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_nodes

Use this data source to query the node list under the specified cluster within HuaweiCloud.

## Example Usage

### Query all nodes under the specified cluster

```hcl
variable "cluster_id" {}

data "huaweicloud_mapreduce_cluster_nodes" "test" {
  cluster_id = var.cluster_id
}
```

### Query Nodes under the specified node group

```hcl
variable "cluster_id" {}
variable "node_group_name" {}

data "huaweicloud_mapreduce_cluster_nodes" "test" {
  cluster_id = var.cluster_id
  node_group = var.node_group_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the nodes are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

* `node_group` - (Optional, String) Specifies the name of the node group to which the node belongs.

* `node_name` - (Optional, String) Specifies the name of the node.  
  Fuzzy search is supported.

* `query_node_detail` - (Optional, Bool) Specifies whether to query node detail.  
  Default to **false**.

* `query_ecs_detail` - (Optional, Bool) Specifies whether to query ECS detail.
  Default to **false**.

* `internal_ip` - (Optional, String) Specifies the internal IP address of the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The list of nodes that match the filter parameters.  
  The [nodes](#cluster_nodes) structure is documented below.

<a name="cluster_nodes"></a>
The `nodes` block supports:

* `node_name` - The name of the node.

* `resource_id` - The resource ID of the node.

* `node_group_name` - The name of the node group to which the node belongs.

* `node_type` - The type of the node.
  + **Master**
  + **Core**
  + **Task**

* `charging_mode` - The billing type of the node.
  + **prePaid**
  + **postPaid**

* `deployment_type` - The deployment type of the node.
  + **SERVER**: Host type.

* `server_info` - The server information of the node.  
  The [server_info](#cluster_nodes_server_info) structure is documented below.

* `tags` - The key/value pairs associated with the node.

* `node_detail` - The monitoring information of the node.  
  The [node_detail](#cluster_nodes_node_detail) structure is documented below.  
  It is not empty only when `query_node_detail` is set to **true**.

* `node_status` - The status of the node.
  + **error**
  + **delete_fail**
  + **shutdown**
  + **started**
  + **attached**
  + **detached**
  + **detached_no_subscription**
  + **detach_fail**
  + **destroyed**
  + **lost**
  + **unknown**
  + **scaling-up**
  + **decommissioned**
  + **decommission_fail**
  + **restore_fail**

* `component_infos` - The list of components deployed on the node.  
  The [component_infos](#cluster_nodes_component_infos) structure is documented below.  
  It is not empty only when `query_node_detail` is set to **true**.

<a name="cluster_nodes_server_info"></a>
The `server_info` block supports:

* `server_id` - The ID of the server.

* `server_name` - The name of the server.

* `server_type` - The type of the server.
  + **ECS**
  + **BMS**

* `data_volumes` - The data disks of the server.  
  The [data_volumes](#cluster_nodes_server_info_volume) structure is documented below.

* `root_volume` - The system disk configuration of the server.  
  The [root_volume](#cluster_nodes_server_info_volume) structure is documented below.

* `cpu_type` - The CPU type of the server.
  + **X86**
  + **ARM**

* `cpu` - The CPU size of the server.  
  It is not empty only when `query_ecs_detail` is set to **true**.

* `mem` - The memory size of the server, in MB.
  It is not empty only when `query_ecs_detail` is set to **true**.

* `internal_ip` - The internal IP address of the server.

<a name="cluster_nodes_server_info_volume"></a>
The `data_volumes` and `root_volume` blocks support:

* `type` - The disk type.
  + **SATA**
  + **SAS**
  + **SSD**
  + **GPSSD**

* `size` - The disk size, in GB.

* `count` - The disk count.

<a name="cluster_nodes_node_detail"></a>
The `node_detail` block supports:

* `running_status` - The running status.
  + **running**
  + **BAD**
  + **UNKNOWN**
  + **ISOLATED**
  + **SUSPENDED**

* `cpu_usage` - The CPU usage.

* `memory_usage` - The memory usage.

* `disk_usage` - The disk usage.

* `total_memory` - The total memory, in MB.

* `available_memory` - The available memory, in MB.

* `total_hard_disk_space` - The total hard disk space, in GB.

* `available_hard_disk_space` - The available hard disk space, in GB.

* `network_read` - The network read speed, in Byte/s.

* `network_write` - The network write speed, in Byte/s.

<a name="cluster_nodes_component_infos"></a>
The `component_infos` block supports:

* `id` - The component ID.

* `name` - The component name.

* `instance_group_name` - The component instance group name.

* `running_status` - The component running status.
  + **running**
  + **BAD**
  + **UNKNOWN**
  + **ISOLATED**
  + **SUSPENDED**
  
* `ha_status` - The HA status.
  + **ACTIVE**
  + **STANDBY**
  + **OBSERVER**
  + **UNKNOWN**

* `config_status` - The config status.
  + **SYNCHRONIZED**  
  + **EXPIRED**
  + **FAILED**
  + **UNKNOWN**

* `role_name` - The role name.

* `role_short_name` - The role short name.

* `role_type` - The role type.

* `service_name` - The service name.

* `pair_name` - The pair name.

* `relation_pairs` - The relation pairs.

* `support_decom` - Whether Decom is supported.

* `support_reinstall` - Whether reinstall is supported.

* `support_collect_stack_info` - Whether stack info collection is supported.

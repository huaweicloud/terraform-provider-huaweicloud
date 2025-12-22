---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node_pool"
description: ""
---

# huaweicloud_cce_node_pool

To get the specified node pool in a cluster.

## Example Usage

```hcl
variable "cluster_id" {}
variable "node_pool_name" {}

data "huaweicloud_cce_node_pool" "node_pool" {
  cluster_id = var.cluster_id
  name       = var.node_pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE node pools.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of container cluster.

* `name` - (Optional, String) Specifies the name of the node pool.

* `node_pool_id` - (Optional, String) Specifies the ID of the node pool.

* `status` - (Optional, String) Specifies the state of the node pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `initial_node_count` - Initial number of nodes in the node pool.

* `current_node_count` - Current number of nodes in the node pool.

* `flavor_id` - The flavor ID.

* `type` - Node Pool type.

* `availability_zone` - The name of the available partition (AZ).

* `os` - Operating System of the node.

* `key_pair` - Key pair name when logging in to select the key pair mode.

* `subnet_id` - The ID of the subnet to which the NIC belongs.

* `max_pods` - The maximum number of instances a node is allowed to create.

* `extend_param` - Extended parameter.

* `extension_scale_groups` - The configurations of extended scaling groups in the node pool.
  The [extension_scale_groups](#extension_scale_groups_struct) structure is documented below.

* `scall_enable` - Whether auto-scaling is enabled.

* `min_node_count` - Minimum number of nodes allowed if auto-scaling is enabled.

* `max_node_count` - Maximum number of nodes allowed if auto-scaling is enabled.

* `scale_down_cooldown_time` - Interval between two scaling operations, in minutes.

* `priority` - Weight of a node pool. A node pool with a higher weight has a higher priority during scaling.

* `labels` - Tags of a Kubernetes node, key/value pair format.

* `tags` - Tags of a VM node, key/value pair format.

* `root_volume` - It corresponds to the system disk related configuration.
  The [root_volume](#volume_struct) structure is documented below.

* `data_volumes` - Represents the data disk to be created.
  The [data_volumes](#volume_struct) structure is documented below.

* `enterprise_project_id` - The enterprise project ID of the node pool.

* `hostname_config` - The hostname config of the kubernetes node.
  The [hostname_config](#hostname_config_struct) structure is documented below.

<a name="extension_scale_groups_struct"></a>
The `extension_scale_groups` block supports:

* `metadata` - The basic information about the extended scaling group.
  The [metadata](#metadata_struct) structure is documented below.

* `spec` - The configurations of the extended scaling group, which carry different configurations from those of the
  default scaling group.
  The [spec](#spec_struct) structure is documented below.

<a name="metadata_struct"></a>
The `metadata` block supports:

* `name` - The name of an extended scaling group.
  Only digits, lowercase letters, and hyphens (-) are allowed.

* `uid` - The extended scaling group UUID.

<a name="spec_struct"></a>
The `spec` block supports:

* `flavor` - The node flavor.

* `az` - The availability zone of a node.

* `capacity_reservation_specification` - The capacity reservation configurations of the extended scaling group.
  The [capacity_reservation_specification](#capacity_reservation_specification_struct) structure is documented below.

* `autoscaling` - The auto-scaling configurations of the extended scaling group.
  The [autoscaling](#autoscaling_struct) structure is documented below.

<a name="capacity_reservation_specification_struct"></a>
The `capacity_reservation_specification` block supports:

* `id` - The private pool ID.

* `preference` - The capacity of a private storage pool.

<a name="autoscaling_struct"></a>
The `autoscaling` block supports:

* `enable` - Whether to enable auto-scaling for the scaling group.

* `extension_priority` - The priority of the scaling group. A higher value indicates a greater priority.

* `min_node_count` - The minimum number of nodes in the scaling group during auto-scaling.

* `max_node_count` - The maximum number of nodes that can be retained in the scaling group during auto-scaling.

<a name="volume_struct"></a>
The `root_volume` and `data_volumes` blocks support:

* `size` - Disk size in GB.

* `volumetype` - Disk type.

* `extend_params` - Disk expansion parameters.

<a name="hostname_config_struct"></a>
The `hostname_config` block supports:

* `type` - The hostname type of the kubernetes node.

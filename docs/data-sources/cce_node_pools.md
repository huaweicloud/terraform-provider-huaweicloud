---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_node_pools"
description: |-
  Use this data source to query CCE node pools within HuaweiCloud.
---

# huaweicloud_cce_node_pools

Use this data source to query CCE node pools within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_cce_node_pools" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the CCE node pools are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the CCE cluster.

* `show_default_node_pool` - (Optional, String) Specifies whether to show the default node pool.  
  The value can be **true** or **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_pools` - The list of CCE node pools that matched filter parameters.  
  The [node_pools](#cce_node_pools) structure is documented below.

<a name="cce_node_pools"></a>
The `node_pools` block supports:

* `id` - The ID of the node pool.

* `name` - The name of the node pool.

* `initial_node_count` - The initial number of nodes in the node pool.

* `current_node_count` - The current number of nodes in the node pool.

* `flavor_id` - The flavor ID of the node pool.

* `type` - The type of the node pool.

* `availability_zone` - The availability zone of the node pool.

* `os` - The operating system of the node pool.

* `key_pair` - The key pair name of the node pool.

* `subnet_id` - The subnet ID of the NIC.

* `subnet_list` - The subnet ID list of the NIC.

* `ecs_group_id` - The ECS group ID of the node pool.

* `max_pods` - The maximum number of pods allowed on a node.

* `extend_param` - The extended parameters of the node pool, in key/value format.

* `extend_params` - The extended parameters of the node pool.  
  The [extend_params](#cce_node_pools_extend_params) structure is documented below.

* `extension_scale_groups` - The extended scaling groups.  
  The [extension_scale_groups](#cce_node_pools_extension_scale_groups) structure is documented below.

* `scall_enable` - Whether auto scaling is enabled.

* `min_node_count` - The minimum number of nodes if auto scaling is enabled.

* `max_node_count` - The maximum number of nodes if auto scaling is enabled.

* `scale_down_cooldown_time` - The interval between two scaling operations, in minutes.

* `priority` - The weight of the node pool during scaling.

* `labels` - The labels of a Kubernetes node.

* `tags` - The tags of a VM node.

* `root_volume` - The system disk configuration of the node pool.  
  The [root_volume](#cce_node_pools_volume) structure is documented below.

* `data_volumes` - The data disk configurations of the node pool.  
  The [data_volumes](#cce_node_pools_volume) structure is documented below.

* `storage` - The disk initialization configuration of the node pool.  
  The [storage](#cce_node_pools_storage) structure is documented below.

* `taints` - The taints configuration of the node pool.  
  The [taints](#cce_node_pools_taints) structure is documented below.

* `security_groups` - The custom security group IDs of the node pool.

* `pod_security_groups` - The pod security group IDs of the node pool.

* `initialized_conditions` - The custom initialization flags of the node pool.

* `label_policy_on_existing_nodes` - The label policy on existing nodes.

* `tag_policy_on_existing_nodes` - The tag policy on existing nodes.

* `taint_policy_on_existing_nodes` - The taint policy on existing nodes.

* `hostname_config` - The hostname configuration of the kubernetes node.  
  The [hostname_config](#cce_node_pools_hostname_config) structure is documented below.

* `partition` - The partition to which the node belongs.

* `enterprise_project_id` - The enterprise project ID of the node pool.

* `runtime` - The runtime of the node pool.

* `billing_mode` - The billing mode of a node.

* `period_unit` - The charging period unit of the node pool.

* `period` - The charging period of the node pool.

* `auto_renew` - Whether auto-renew is enabled.

* `status` - The status of the node pool.

<a name="cce_node_pools_extend_params"></a>
The `extend_params` block supports:

* `max_pods` - The maximum number of pods allowed on a node.

* `docker_base_size` - The available disk space of a single container on a node, in GB.

* `preinstall` - The script to be executed before installation.

* `postinstall` - The script to be executed after installation.

* `node_image_id` - The image ID used to create the node.

* `node_multi_queue` - The number of ENI queues.

* `nic_threshold` - The ENI pre-binding thresholds.

* `agency_name` - The agency name of the node pool.

* `kube_reserved_mem` - The reserved memory for Kubernetes components.

* `system_reserved_mem` - The reserved memory for system components.

* `security_reinforcement_type` - The security reinforcement type.

* `market_type` - The market type of the spot node pool.

* `spot_price` - The highest price per hour for a spot node.

<a name="cce_node_pools_extension_scale_groups"></a>
The `extension_scale_groups` block supports:

* `metadata` - The basic information about the extended scaling group.  
  The [metadata](#cce_node_pools_extension_scale_groups_metadata) structure is documented below.

* `spec` - The configurations of the extended scaling group.  
  The [spec](#cce_node_pools_extension_scale_groups_spec) structure is documented below.

<a name="cce_node_pools_extension_scale_groups_metadata"></a>
The `metadata` block supports:

* `name` - The name of the extended scaling group.

* `uid` - The UUID of the extended scaling group.

<a name="cce_node_pools_extension_scale_groups_spec"></a>
The `spec` block supports:

* `flavor` - The node flavor.

* `az` - The availability zone.

* `capacity_reservation_specification` - The capacity reservation configurations.  
  The [capacity_reservation_specification](#cce_node_pools_capacity_reservation_specification) structure is documented below.

* `autoscaling` - The auto scaling configurations.  
  The [autoscaling](#cce_node_pools_autoscaling) structure is documented below.

<a name="cce_node_pools_capacity_reservation_specification"></a>
The `capacity_reservation_specification` block supports:

* `id` - The private pool ID.

* `preference` - The private pool capacity preference.

<a name="cce_node_pools_autoscaling"></a>
The `autoscaling` block supports:

* `enable` - Whether auto scaling is enabled.

* `extension_priority` - The priority of the scaling group.

* `min_node_count` - The minimum number of nodes.

* `max_node_count` - The maximum number of nodes.

<a name="cce_node_pools_volume"></a>
The `root_volume` and `data_volumes` blocks support:

* `size` - The disk size in GB.

* `volumetype` - The disk type.

* `extend_params` - The disk expansion parameters.

* `kms_key_id` - The KMS key ID of the disk.

* `dss_pool_id` - The DSS pool ID of the disk.

* `iops` - The IOPS of the disk.

* `throughput` - The throughput of the disk.

* `hw_passthrough` - Whether passthrough is enabled.

<a name="cce_node_pools_storage"></a>
The `storage` block supports:

* `selectors` - The disk selection configuration.  
  The [selectors](#cce_node_pools_storage_selectors) structure is documented below.

* `groups` - The storage group configuration.  
  The [groups](#cce_node_pools_storage_groups) structure is documented below.

<a name="cce_node_pools_storage_selectors"></a>
The `selectors` block supports:

* `name` - The selector name.

* `type` - The storage type.

* `match_label_size` - The matched disk size.

* `match_label_volume_type` - The EVS disk type.

* `match_label_metadata_encrypted` - The disk encryption identifier.

* `match_label_metadata_cmkid` - The customer master key ID of an encrypted disk.

* `match_label_count` - The number of disks to be selected.

<a name="cce_node_pools_storage_groups"></a>
The `groups` block supports:

* `name` - The name of a virtual storage group.

* `cce_managed` - Whether the storage space is for kubernetes and runtime components.

* `selector_names` - The list of selector names to match.

* `virtual_spaces` - The detailed space configuration in a group.  
  The [virtual_spaces](#cce_node_pools_storage_virtual_spaces) structure is documented below.

<a name="cce_node_pools_storage_virtual_spaces"></a>
The `virtual_spaces` block supports:

* `name` - The virtual space name.

* `size` - The size of a virtual space.

* `lvm_lv_type` - The LVM write mode.

* `lvm_path` - The absolute path to which the disk is attached.

* `runtime_lv_type` - The LVM write mode of runtime.

<a name="cce_node_pools_taints"></a>
The `taints` block supports:

* `key` - The key of the taint.

* `value` - The value of the taint.

* `effect` - The effect of the taint.

<a name="cce_node_pools_hostname_config"></a>
The `hostname_config` block supports:

* `type` - The hostname type of the kubernetes node.

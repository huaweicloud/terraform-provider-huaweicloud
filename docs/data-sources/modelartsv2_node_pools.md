---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_node_pools"
description: |-
  Use this data source to get the list of node pools under a specified resource pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_node_pools

Use this data source to get the list of node pools under a specified resource pool within HuaweiCloud.

## Example Usage

```hcl
variable "pool_name" {}

data "huaweicloud_modelartsv2_node_pools" "test" {
  pool_name = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the resource nodes are located.  
  If omitted, the provider-level region will be used.

* `pool_name` - (Required, String) Specifies the resource pool name to which the node pool belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `node_pools` - The list of the node pools.  
  The [node_pools](#modelartsv2_node_pools_node_pools) structure is documented below.

<a name="modelartsv2_node_pools_node_pools"></a>
The `node_pools` block supports:

* `metadata` - The metadata information of the node pool.  
  The [metadata](#modelartsv2_node_pools_node_pools_metadata) structure is documented below.

* `spec` - The expectation information of the node pool.  
  The [spec](#modelartsv2_node_pools_node_pools_spec) structure is documented below.

* `status` - The status information of the node pool.  
  The [status](#modelartsv2_node_pools_node_pools_status) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_metadata"></a>
The `metadata` block supports:

* `name` - The name of the node pool.

<a name="modelartsv2_node_pools_node_pools_spec"></a>
The `spec` block supports:

* `resources` - The list of resources in the node pool.  
  The [resources](#modelartsv2_node_pools_node_pools_spec_resource) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_spec_resource"></a>
The `resources` block supports:

* `flavor` - The resource flavor name.

* `count` - The desired usage of the flavor.

* `max_count` - The elastic usage of the resource flavor.

* `azs` - The availability zones information of nodes in the resource pool.  
  The [azs](#modelartsv2_node_pools_node_pools_spec_resource_azs) structure is documented below.

* `node_pool` - The name of the node pool.

* `taints` - The taints information.  
  The [taints](#modelartsv2_node_pools_node_pools_spec_resource_taints) structure is documented below.

* `labels` - The Kubernetes label.

* `tags` - The resource tags.  

* `network` - The network configuration.  
  The [network](#modelartsv2_node_pools_node_pools_spec_resource_network) structure is documented below.

* `extend_params` - The custom configuration.  
  The [extend_params](#modelartsv2_node_pools_node_pools_spec_resource_extend_params) structure is documented below.

* `creating_step` - The information about batch creation.  
  The [creating_step](#modelartsv2_node_pools_node_pools_spec_resource_creating_step) structure is documented below.

* `root_volume` - The custom system disk information.  
  The [root_volume](#modelartsv2_node_pools_node_pools_spec_resource_root_volume) structure is documented below.

* `data_volumes` - The custom data disks information.  
  The [data_volumes](#modelartsv2_node_pools_node_pools_spec_resource_data_volumes) structure is documented below.

* `volume_group_configs` - The advanced disk configuration.  
  The [volume_group_configs](#modelartsv2_node_pools_node_pools_spec_resource_volume_group_configs) structure  
  is documented below.

* `os` - The OS image information.  
  The [os](#modelartsv2_node_pools_node_pools_spec_resource_os) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_spec_resource_azs"></a>
The `azs` block supports:

* `az` - The availability zone name.

* `count` - The number of availability zone resource instances.

<a name="modelartsv2_node_pools_node_pools_spec_resource_taints"></a>
The `taints` block supports:

* `key` - The taint key.

* `value` - The taint value.

* `effect` - The effect of the action.

<a name="modelartsv2_node_pools_node_pools_spec_resource_network"></a>
The `network` block supports:

* `vpc` - The ID of the VPC.

* `subnet` - The ID of the subnet.

* `security_groups` - The IDs of the security groups.

<a name="modelartsv2_node_pools_node_pools_spec_resource_extend_params"></a>
The `extend_params` block supports:

* `docker_base_size` - The container image space size of a node.

* `post_install` - The post-installation script.

* `runtime` - The container runtime.

* `label_policy_on_existing_nodes` - The Kubernetes label update policy of existing nodes.

* `taint_policy_on_existing_nodes` - The Kubernetes taint update policy of existing nodes.

* `tag_policy_on_existing_nodes` - The resource tag update policy of existing nodes.

* `x_parameter_plane_subnet` - The subnet ID used for data transmission on the parameter plane between  
  physical clusters.

* `node_pool_name` - The name of the node pool specified by user.

<a name="modelartsv2_node_pools_node_pools_spec_resource_creating_step"></a>
The `creating_step` block supports:

* `step` - The step of a supernode.

* `type` - The batch creation type.

<a name="modelartsv2_node_pools_node_pools_spec_resource_root_volume"></a>
The `root_volume` block supports:

* `volume_type` - The disk type.

* `size` - The disk size.

<a name="modelartsv2_node_pools_node_pools_spec_resource_data_volumes"></a>
The `data_volumes` block supports:

* `volume_type` - The disk type.

* `size` - The disk size.

* `count` - The number of disks.

* `extend_params` - The custom disk configuration.  
  The [extend_params](#modelartsv2_node_pools_node_pools_spec_resource_data_volumes_extend_params) structure  
  is documented below.

<a name="modelartsv2_node_pools_node_pools_spec_resource_data_volumes_extend_params"></a>
The `extend_params` block supports:

* `volume_group` - The name of the disk group.

<a name="modelartsv2_node_pools_node_pools_spec_resource_volume_group_configs"></a>
The `volume_group_configs` block supports:

* `volume_group` - The disk group name.

* `docker_thin_pool` - The percentage of container disks to data disks on nodes in a resource pool.

* `lvm_config` - The LVM configuration.  
  The [lvm_config](#modelartsv2_node_pools_node_pools_spec_resource_volume_group_config_lvm_config) structure  
  is documented below.

* `types` - The storage type.

<a name="modelartsv2_node_pools_node_pools_spec_resource_volume_group_config_lvm_config"></a>
The `lvm_config` block supports:

* `lv_type` - The LVM write mode.

* `path` - The disk mount path.

<a name="modelartsv2_node_pools_node_pools_spec_resource_os"></a>
The `os` block supports:

* `name` - The OS name and version.

* `image_id` - The OS image ID.

* `image_type` - The OS image type.

* `auto_match` - The automatic OS image matching configuration.

<a name="modelartsv2_node_pools_node_pools_status"></a>
The `status` block supports:

* `resources` - The resources in different states in the node pool.  
  The [resources](#modelartsv2_node_pools_node_pools_status_resources) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_status_resources"></a>
The `resources` block supports:

* `creating` - The number of resources that are being created.  
  The [creating](#modelartsv2_node_pools_node_pools_status_resources_creating) structure is documented below.

* `available` - The number of available resources.  
  The [available](#modelartsv2_node_pools_node_pools_status_resources_available) structure is documented below.

* `abnormal` - The number of abnormal resources.  
  The [abnormal](#modelartsv2_node_pools_node_pools_status_resources_abnormal) structure is documented below.

* `deleting` - The number of resources that are being deleted.  
  The [deleting](#modelartsv2_node_pools_node_pools_status_resources_deleting) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_status_resources_creating"></a>
The `creating` block supports:

* `flavor` - The resource flavor ID.

* `count` - The number of resource specification instances in the resource pool.

* `max_count` - The number of elastic resource specification instances in the resource pool.

* `azs` - The AZ distribution of the resource specification instances to be created in the resource pool.  
  The [azs](#modelartsv2_node_pools_node_pools_status_resources_flavor_azs) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_status_resources_available"></a>
The `available` block supports:

* `flavor` - The resource flavor ID.

* `count` - The number of resource specification instances in the resource pool.

* `max_count` - The number of elastic resource specification instances in the resource pool.

* `azs` - The AZ distribution of the resource specification instances to be created in the resource pool.  
  The [azs](#modelartsv2_node_pools_node_pools_status_resources_flavor_azs) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_status_resources_abnormal"></a>
The `abnormal` block supports:

* `flavor` - The resource flavor ID.

* `count` - The number of resource specification instances in the resource pool.

* `max_count` - The number of elastic resource specification instances in the resource pool.

* `azs` - The AZ distribution of the resource specification instances to be created in the resource pool.  
  The [azs](#modelartsv2_node_pools_node_pools_status_resources_flavor_azs) structure is documented below.

<a name="modelartsv2_node_pools_node_pools_status_resources_deleting"></a>
The `deleting` block supports:

* `flavor` - The resource flavor ID.

* `count` - The number of resource specification instances in the resource pool.

* `max_count` - The number of elastic resource specification instances in the resource pool.

* `azs` - The AZ distribution of the resource specification instances to be created in the resource pool.  
  The [azs](#modelartsv2_node_pools_node_pools_status_resources_flavor_azs) structure is documented below.

* `node_pool` - The node pool ID.

<a name="modelartsv2_node_pools_node_pools_status_resources_flavor_azs"></a>
The `azs` block supports:

* `az` - The availability zone name.

* `count` - The number of availability zone resource instances.

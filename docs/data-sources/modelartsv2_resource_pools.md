---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_resource_pools"
description: |-
  Use this data source to get list of ModelArts resource pools within HuaweiCloud.
---

# huaweicloud_modelartsv2_resource_pools

Use this data source to get list of ModelArts resource pools within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_modelartsv2_resource_pools" "test" {
  status = "created"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Optional, String) Specifies the workspace ID to which the resource pool belongs.

* `status` - (Optional, String) Specifies the status of the resource pool to be queried.  
  The valid values are as follows:
  + **created**
  + **failed**
  + **creating**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_pools` - All resource pools that matched filter parameters.  
  The [resource_pools](#modelarts_v2_resource_pools) structure is documented below.

<a name="modelarts_v2_resource_pools"></a>
The `resource_pools` block supports:

* `metadata` - The metadata configuration of the resource pool.  
  The [metadata](#modelarts_v2_resource_pools_metadata) structure is documented below.

* `spec` - The specification of the resource pool.  
  The [spec](#modelarts_v2_resource_pools_spec) structure is documented below.

* `status` - The status of the resource pool.

<a name="modelarts_v2_resource_pools_metadata"></a>
The `metadata` block supports:

* `name` - The name of the resource pool.

* `labels` - The labels of the resource pool.

* `annotations` - The annotations of the resource pool.

* `created_at` - The creation time of the resource pool, in RFC3339 format.

<a name="modelarts_v2_resource_pools_spec"></a>
The `spec` block supports:

* `resources` - The list of resource specifications in the resource pool.  
  The [resources](#modelarts_v2_resource_pools_spec_resources) structure is documented below.

* `scope` - The list of job types supported by the resource pool.

* `network` - The network of the resource pool.  
  The [network](#modelarts_v2_resource_pools_spec_network) structure is documented below.

* `user_login` - The user login information of the privileged pool.  
  The [user_login](#modelarts_v2_resource_pools_spec_user_login) structure is documented below.

* `clusters` - The cluster information of the privileged pool.  
  The [clusters](#modelarts_v2_resource_pools_spec_clusters) structure is documented below.

<a name="modelarts_v2_resource_pools_spec_resources"></a>
The `resources` block supports:

* `flavor` - The flavor of the resource pool.

* `count` - The count of the resource pool.

* `max_count` - The max number of resources of the corresponding flavors.

* `node_pool` - The name of resource pool nodes.

* `taints` - The taint list of the resource pool.  
  The [taints](#modelarts_v2_resource_pools_spec_resources_taints) structure is documented below.

* `labels` - The key/value pairs labels of resource pool.

* `tags` - The key/value pairs to associate with the resource pool nodes.

* `network` - The network of the privileged pool.  
  The [network](#modelarts_v2_resource_pools_spec_resources_network) structure is documented below.

* `extend_params` - The extend params of the resource pool, in JSON format.

* `creating_step` - The creation step configuration of the resource pool nodes.  
  The [creating_step](#modelarts_v2_resource_pools_spec_resources_creating_step) structure is documented below.

* `root_volume` - The root volume of the resource pool nodes.  
  The [root_volume](#modelarts_v2_resource_pools_spec_resources_root_volume) structure is documented below.

* `data_volumes` - The data volumes of the resource pool nodes.  
  The [data_volumes](#modelarts_v2_resource_pools_spec_resources_data_volumes) structure is documented below.

* `volume_group_configs` - The extend configurations of the volume groups.  
  The [volume_group_configs](#modelarts_v2_resource_pools_spec_resources_volume_group_configs) structure is documented below.
  
* `os` - The image information for the specified OS.  
  The [os](#modelarts_v2_resource_pools_spec_resources_os) structure is documented below.

* `azs` - The AZ list of the resource pool nodes.  
  The [azs](#modelarts_v2_resource_pools_spec_resources_azs) structure is documented below.

<a name="modelarts_v2_resource_pools_spec_network"></a>
The `network` block supports:

* `name` - The name of the network.

* `subnet_id` - The ID of the subnet.

* `vpc_id` - The ID of the VPC.

<a name="modelarts_v2_resource_pools_spec_user_login"></a>
The `user_login` block supports:

* `key_pair_name` - The name of the key pair.

<a name="modelarts_v2_resource_pools_spec_clusters"></a>
The `clusters` block supports:

* `name` - The name of the cluster.

* `provider_id` - The provider ID of the cluster.

<a name="modelarts_v2_resource_pools_spec_resources_taints"></a>
The `taints` block supports:

* `effect` - The effect of the taint.

* `key` - The key of the taint.

* `value` - The value of the taint.

<a name="modelarts_v2_resource_pools_spec_resources_network"></a>
The `network` block supports:

* `security_groups` - The ID list of the security group.

* `subnet` - The ID of the subnet.

* `vpc` - The ID of the VPC.

<a name="modelarts_v2_resource_pools_spec_resources_creating_step"></a>
The `creating_step` block supports:

* `step` - The creation step of the resource pool nodes.

* `type` - The type of the resource pool nodes.

<a name="modelarts_v2_resource_pools_spec_resources_root_volume"></a>
The `root_volume` block supports:

* `size` - The size of the root volume.

* `volume_type` - The type of the root volume.

* `extend_params` - The extend parameters of the root volume, in JSON format.

<a name="modelarts_v2_resource_pools_spec_resources_data_volumes"></a>
The `data_volumes` block supports:

* `count` - The count of the current data volume configuration.

* `extend_params` - The extend parameters of the data volume, in JSON format.

* `size` - The size of the data volume.

* `volume_type` - The type of the data volume.

<a name="modelarts_v2_resource_pools_spec_resources_volume_group_configs"></a>
The `volume_group_configs` block supports:

* `docker_thin_pool` - The percentage of container volumes to data volumes on resource pool nodes.

* `lvm_config` - The configuration of the LVM management.  
  The [lvm_config](#modelarts_v2_resource_pools_spec_resources_volume_group_configs_lvm_config) structure is documented
  below.

* `types` - The storage types of the volume group.

* `volume_group` - The name of the volume group.

<a name="modelarts_v2_resource_pools_spec_resources_os"></a>
The `os` block supports:

* `image_id` - The image ID.

* `image_type` - The image type.

* `name` - The OS name of the image.

<a name="modelarts_v2_resource_pools_spec_resources_azs"></a>
The `azs` block supports:

* `az` - The AZ name

* `count` - The number of nodes in the AZ.

<a name="modelarts_v2_resource_pools_spec_resources_volume_group_configs_lvm_config"></a>
The `lvm_config` block supports:

* `lv_type` - The LVM write mode.

* `path` - The volume mount path.

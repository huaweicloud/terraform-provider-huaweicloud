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

* `workspace_id` - (Optional, String) The workspace ID to which the resource pool belongs.

* `status` - (Optional, String) The status of the resource pool to be queried.  
  The valid values are as follows:
  + **created**
  + **failed**
  + **creating**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `resource_pools` - All resource pools that matched filter parameters.  
  The [resource_pools](#modelarts_resource_pools) structure is documented below.

<a name="modelarts_resource_pools"></a>
The `resource_pools` block supports:

* `metadata` - The metadata configuration of the resource pool.  
  The [metadata](#modelarts_resource_pool_metadata) structure is documented below.

* `name` - The name of the resource pool.

* `scope` - The list of job types supported by the resource pool.

* `resources` - The list of resource specifications in the resource pool.  
  The [resources](#modelarts_resource_pool_resources) structure is documented below.

* `network_id` - The ModelArts network ID of the resource pool.

* `prefix` - The prefix of the user-defined node name of the resource pool.

* `vpc_id` - The ID of the VPC to which the resource pool belongs.

* `subnet_id` - The network ID of the subnet to which the resource pool belongs.

* `clusters` - The list of the CCE clusters.  
  The [clusters](#modelarts_resource_pool_clusters) structure is documented below.

* `user_login` - The user login info of the resource pool.  
  The [user_login](#modelarts_resource_pool_user_login) structure is documented below.

* `workspace_id` - The workspace ID of the resource pool.
  + **0**: Default workspace.

* `description` - The description of the resource pool.

* `charging_mode` - The charging mode of the resource pool.
  + **prePaid**: The yearly/monthly billing mode.

* `status` - The status of the resource pool.

* `resource_pool_id` - The resource ID of the resource pool.

<a name="modelarts_resource_pool_metadata"></a>
The `metadata` block supports:

* `name` - The name of the resource pool.

* `annotations` - The annotations of the resource pool.

* `labels` - The labels of the resource pool.

* `created_at` - The creation time of the resource pool, in RFC3339 format.

<a name="modelarts_resource_pool_resources"></a>
The `resources` block supports:

* `flavor_id` - The resource flavor ID.

* `count` - The number of resources of the corresponding flavors.

* `node_pool` - The name of resource pool nodes.

* `max_count` - The max number of resources of the corresponding flavors.

* `vpc_id` - The ID of the VPC to which the the resource pool nodes belong.

* `subnet_id` - The network ID of a subnet to which the the resource pool nodes belong.

* `security_group_ids` - The security group IDs to which the the resource pool nodes belong.

* `azs` - The availability zones for the resource pool nodes.
  The [azs](#modelarts_resource_pool_resource_azs) structure is documented below.

* `taints` - The taints added to resource pool nodes.  
  The [taints](#modelarts_resource_pool_resource_taints) structure is documented below.

* `labels` - The labels of resource pool nodes.

* `tags` - The key/value pairs to associate with the resource pool nodes.

* `extend_params` - The extend parameters of the resource pool nodes.

* `root_volume` - The root volume of the resource pool nodes.  
  The [root_volume](#modelarts_resource_pool_resource_root_volume) structure is documented below.

* `data_volumes` - The list of data volumes of the resource pool nodes.  
  The [data_volumes](#modelarts_resource_pool_resource_data_volumes) structure is documented below.

* `volume_group_configs` - The extend configurations of the volume groups.  
  The [volume_group_configs](#modelarts_resource_pool_resource_volume_group_configs) structure is documented below.

<a name="modelarts_resource_pool_resource_azs"></a>
The `azs` block supports:

* `az` - The AZ name.

* `count` - The number of nodes for the corresponding AZ

<a name="modelarts_resource_pool_resource_taints"></a>
The `taints` block supports:

* `key` - The key of the taint.

* `value` - The value of the taint.

* `effect` - The effect of the taint.

<a name="modelarts_resource_pool_resource_root_volume"></a>
The `root_volume` block supports:

* `volume_type` - The type of the root volume.

* `size` - The size of the root volume.

<a name="modelarts_resource_pool_resource_data_volumes"></a>
The `root_volume` block supports:

* `volume_type` - The type of the data volume.

* `size` - The size of the data volume.

* `extend_params` - The extend parameters of the data volume.

* `count` - The count of the current data volume configuration.

<a name="modelarts_resource_pool_resource_volume_group_configs"></a>
The `volume_group_configs` block supports:

* `volume_group` - The name of the volume group.

* `docker_thin_pool` - The percentage of container volumes to data volumes on resource pool nodes.

* `lvm_config` - The configuration of the LVM management.  
  The [lvm_config](#modelarts_resource_pool_group_config_lvm_config) structure is documented below.

* `types` - The list of storage types of the volume group.

<a name="modelarts_resource_pool_group_config_lvm_config"></a>
The `lvm_config` block supports:

* `lv_type` - The LVM write mode.

* `path` - The volume mount path.

<a name="modelarts_resource_pool_clusters"></a>
The `clusters` block supports:

* `provider_id` - The ID of the CCE cluster that resource pool used.

* `name` - The name of the CCE cluster that resource pool used.

<a name="modelarts_resource_pool_user_login"></a>
The `user_login` block supports:

* `key_pair_name` - The key pair name of the login user.

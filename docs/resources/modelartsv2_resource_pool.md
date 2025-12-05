---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_resource_pool"
description: |-
  Manages a ModelArts resource pool resource within HuaweiCloud.
---

# huaweicloud_modelartsv2_resource_pool

Manages a ModelArts resource pool resource within HuaweiCloud.

## Example Usage

```hcl

```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `metadata` - (Required, List) Specifies the metadata of the resource pool.  
  The [metadata](#v2_resource_pool_metadata) structure is documented below.  
  For more fileds details, please refer to the [document](https://support.huaweicloud.com/intl/en-us/api-modelarts/CreatePool.html#EN-US_TOPIC_0000001868289874__request_PoolMetadataCreation).

* `spec` - (Required, List) Specifies the specification of the resource pool.  
  The [spec](#v2_resource_pool_spec) structure is documented below.

<a name="v2_resource_pool_metadata"></a>
The `metadata` block supports:

* `labels` - (Required, String) Specifies the labels of the resource pool, in JSON format.

* `annotations` - (Optional, String) Specifies the annotations of the resource pool, in JSON format.

-> When creating a resource pool, the billing-related parameters in this parameter indicate the resource pool,
   and when expanding, it indicates the nodes and will be applied to all nodes that are expanded under
   the `resources` parameter.

<a name="v2_resource_pool_spec"></a>
The `spec` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the type of the resource pool.

* `resources` - (Required, List) Specifies the list of resource specifications in the resource pool.  
  Including resource flavors and the number of resources of the corresponding flavors.  
  The [resources](#v2_resource_pool_spec_resources) structure is documented below.

* `scope` - (Optional, List) Specifies the list of job types supported by the resource pool.  
  The valid values are as follows:
  + **Train**: Training job.
  + **Infer**: Inference job.
  + **Notebook**: Notebook job.

* `network` - (Optional, List, NonUpdatable) Specifies the network of the resource pool.  
  This parameter is only valid and required for physical resource pool.  
  The [network](#v2_resource_pool_spec_network) structure is documented below.

* `user_login` - (Optional, List, NonUpdatable) Specifies the user login information of the privileged pool.  
  The [user_login](#v2_resource_pool_spec_user_login) structure is documented below.
  
* `clusters` - (Optional, List, NonUpdatable) Specifies the cluster information of the privileged pool.  
  The [clusters](#v2_resource_pool_spec_clusters) structure is documented below.

<a name="v2_resource_pool_spec_resources"></a>
The `resources` block supports:

* `flavor` - (Required, String) Specifies the flavor of the resource pool.

* `count` - (Required, Int) Specifies the count of the resource pool.

* `max_count` - (Optional, Int) Specifies the max number of resources of the corresponding flavors.

* `node_pool` - (Optional, String) Specifies the name of resource pool nodes.

* `taints` - (Optional, List) Specifies the taint list of the resource pool.  
  This parameter cannot be specified for non-privileged pools.  
  The [taints](#v2_resource_pool_spec_resources_taints) structure is documented below.

* `labels` - (Optional, Map) Specifies the key/value pairs labels of resource pool.  
  This parameter cannot be specified for non-privileged pools.

* `tags` - (Optional, Map) Specifies the key/value pairs tags of resource pool nodes.  
  The [tags](#v2_resource_pool_spec_resources_tags) structure is documented below.

* `network` - (Optional, List) Specifies the network of the privileged pool.  
  This parameter cannot be specified for non-privileged pools.  
  The [network](#v2_resource_pool_spec_resources_network) structure is documented below.

* `extend_params` - (Optional, String) Specifies the extend params of the resource pool, in JSON format.

* `creating_step` - (Optional, List) Specifies the creation step configuration of the
  resource pool nodes.  
  This parameter cannot be updated.  
  The [creating_step](#v2_resource_pool_spec_resources_creating_step) structure is documented below.

* `root_volume` - (Optional, List) Specifies the root volume of the resource pool nodes.  
  The [root_volume](#v2_resource_pool_spec_resources_root_volume) structure is documented below.

* `data_volumes` - (Optional, List) Specifies the data volumes of the resource pool nodes.  
  The [data_volumes](#v2_resource_pool_spec_resources_data_volumes) structure is documented below.

* `volume_group_configs` - (Optional, List) Specifies the extend configurations of the volume groups.  
  This parameter is required when custom data disks is specified.
  The [volume_group_configs](#v2_resource_pool_spec_resources_volume_group_configs) structure is documented below.

* `os` - (Optional, List) Specifies the image information for the specified OS.  
  The [os](#v2_resource_pool_spec_resources_os) structure is documented below.

* `driver` - (Optional, List) Specifies the driver information.  
  The [driver](#v2_resource_pool_spec_resources_driver) structure is documented below.

<a name="v2_resource_pool_spec_resources_taints"></a>
The `taints` block supports:

* `key` - (Required, String) Specifies the key of the taint.

* `effect` - (Required, String) Specifies the effect of the taint.

* `value` - (Optional, String) Specifies the value of the taint.

<a name="v2_resource_pool_spec_resources_tags"></a>
The `tags` block supports:

* `key` - (Required, String) Specifies the key of the tag.  
  Cannot start with `CCE-` or `__type_baremetal`.

* `value` - (Required, String) Specifies the value of the tag.

<a name="v2_resource_pool_spec_resources_network"></a>
The `network` block supports:

* `vpc` - (Optional, String) Specifies the ID of the VPC.

* `subnet` - (Optional, String) Specifies the ID of the subnet.

* `security_groups` - (Optional, List) Specifies the ID list of the security group.

<a name="v2_resource_pool_spec_resources_creating_step"></a>
The `creating_step` block supports:

* `step` - (Required, Int) Specifies the creation step of the resource pool nodes.

* `type` - (Required, String) Specifies the type of the resource pool nodes.  
  The valid values are as follows:
  + **hyperinstance**

<a name="v2_resource_pool_spec_resources_root_volume"></a>
The `root_volume` block supports:

* `size` - (Required, String) Specifies the size of the root volume.

* `volume_type` - (Required, String) Specifies the type of the root volume, in Gi.  
  The valid values are as follows:
  + **SSD**
  + **GPSSD**
  + **SAS**

<a name="v2_resource_pool_spec_resources_data_volumes"></a>
The `data_volumes` block supports:

* `volume_type` - (Required, String) Specifies the type of the data volume.  
  The valid values are as follows:
  + **SSD**
  + **GPSSD**
  + **SAS**

* `size` - (Required, String) Specifies the size of the data volume, in Gi.

* `extend_params` - (Optional, String) Specifies the extend parameters of the data volume, in JSON format.  
  
* `count` - (Optional, Int) Specifies the count of the current data volume configuration.

<a name="v2_resource_pool_spec_resources_volume_group_configs"></a>
The `volume_group_configs` block supports:

* `volume_group` - (Required, String) Specifies the name of the volume group.

* `docker_thin_pool` - (Optional, Int) Specifies the percentage of container volumes to data volumes
  on resource pool nodes.  
  This is only valid when the disk group name is **vgpass**.

* `lvm_config` - (Optional, List) Specifies the configuration of the LVM management.  
  The [lvm_config](#v2_resource_pool_spec_resources_volume_group_configs_lvm_config) structure is documented below.

* `types` - (Optional, List) Specifies the storage types of the volume group.  
  The valid values are as follows:
  + **volume**
  + **local**: Local disk, this parameter must be specified.

<a name="v2_resource_pool_spec_resources_os"></a>
The `os` block supports:

* `name` - (Optional, String) Specifies the OS name of the image.

* `image_id` - (Optional, String) Specifies the image ID.

* `image_type` - (Optional, String) Specifies the image type.

<a name="v2_resource_pool_spec_resources_driver"></a>
The `driver` block supports:

* `version` - (Optional, String) Specifies the driver version.

<a name="v2_resource_pool_spec_resources_volume_group_configs_lvm_config"></a>
The `lvm_config` block supports:

* `lv_type` - (Required, String) Specifies the LVM write mode.  
  The valid values are as follows:
  + **liner**: Linear mode.
  + **striped**: Stripe mode.

* `path` - (Optional, String) Specifies the volume mount path.  
  Only numbers, letters, dots (.), hyphens (-), and underscores (_) are allowed.

<a name="v2_resource_pool_spec_network"></a>
The `network` block supports:

* `name` - (Optional, String, NonUpdatable) Specifies the name of the network.

* `vpc_id` - (Optional, String, NonUpdatable) Specifies the ID of the VPC.  
  For privileged pool, this parameter is only valid and required.

* `subnet_id` - (Optional, String, NonUpdatable) Specifies the ID of the subnet.  
  For privileged pool, this parameter is only valid and required.

<a name="v2_resource_pool_spec_user_login"></a>
The `user_login` block supports:

* `key_pair_name` - (Optional, String, NonUpdatable) Specifies the name of the key pair.

* `password` - (Optional, String, NonUpdatable) Specifies the password of the resource pool.  
 The value needs to be salted, encrypted and Base64 encoded. Default user is **root**.

<a name="v2_resource_pool_spec_clusters"></a>
The `clusters` block supports:

* `provider_id` - (Optional, String, NonUpdatable) Specifies the provider ID of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the resource pool.

* `clusters` - The cluster information of the privileged pool.  
  The [clusters](#v2_resource_pool_spec_clusters_attr) structure is documented below.

<a name="v2_resource_pool_spec_clusters_attr"></a>
The `clusters` block supports:

* `name` - The name of the cluster.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 90 minutes.
* `update` - Default is 90 minutes.
* `delete` - Default is 30 minutes.

## Import

The resource can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_modelartsv2_resource_pool.test <id>
```

Please add the followings if some attributes are missing when importing the resource.

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `metadata.0.annotations` and `spec.0.user_login`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_modelartsv2_resource_pool" "test" {
  ...

  lifecycle {
    ignore_changes = [
      metadata.0.annotations, spec.0.user_login
    ]
  }
}

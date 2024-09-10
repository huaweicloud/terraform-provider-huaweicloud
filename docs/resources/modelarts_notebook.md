---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook"
description: ""
---

# huaweicloud_modelarts_notebook

Manages ModelArts notebook resource within HuaweiCloud.

## Example Usage

### Create a notebook with the EVS storage type

```hcl
variable "notebook_name" {}
variable "key_pair_name" {}
variable "image_id" {}
variable "allowed_ip_addresses" {
  type = list(string)
}
variable "key_pair_name" {}

resource "huaweicloud_modelarts_notebook" "test" {
  name      = var.notebook_name
  flavor_id = "modelarts.vm.cpu.2u"
  image_id  = var.image_id

  allowed_access_ips = var.allowed_ip_addresses
  key_pair           = var.key_pair_name

  volume {
    type = "EVS"
    size = 5
  }
}
```

### Create a notebook with the EFS storage type

```hcl
variable "notebook_name" {} 
variable "image_id" {}
variable "resource_pool_id" {}
variable "sfs_export_location" {}
variable "sfs_turbo_id" {}

resource "huaweicloud_modelarts_notebook" "test" {
  name      = var.notebook_name
  flavor_id = "modelarts.vm.cpu.2u"
  image_id  = var.image_id
  pool_id   = var.resource_pool_id

  volume {
    type      = "EFS"
    ownership = "DEDICATED"
    uri       = var.sfs_export_location
    id        = var.sfs_turbo_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the notebook. The name consists of 1 to 64 characters,
 starting with a letter. Only letters, digits and underscores (_) are allowed.

* `flavor_id` - (Required, String) Specifies the flavor ID. The options are as follows:
  - **modelarts.vm.cpu.2u**: General-purpose Intel CPU specifications, suitable for data exploration and algorithm
   discovery.
  - **modelarts.vm.cpu.8u**: General computing-plus Intel CPU specifications, suitable for compute-intensive
   applications.
  - **modelarts.bm.gpu.v100NV32**: One NVIDIA V100 GPU with 32GB of memory, suitable for deep learning algorithm
   training and debugging.
  - **modelarts.bm.d910.xlarge.1**: One Ascend 910 NPU with 32GB of memory, suitable for deep learning code running
   and debugging.
  - **modelarts.bm.d910.xlarge.2**: Two Ascend 910 NPU with 32GB of memory, suitable for deep learning code running
   and debugging.
  - **modelarts.bm.d910.xlarge.8**: Eight Ascend 910 NPU with 32GB of memory, suitable for deep learning code running
   and debugging.

* `image_id` - (Required, String) Specifies the image ID of notebook.

* `volume` - (Required, List) Specifies the volume information. Structure is documented below.

* `description` - (Optional, String) Specifies the description of notebook. It contains a maximum of `512` characters and
 cannot contain special characters `&<>"'/`.

* `key_pair` - (Optional, String, ForceNew) Specifies the key pair name for remote SSH access.
 Changing this parameter will create a new resource.

* `allowed_access_ips` - (Optional, List) Specifies public IP addresses that are allowed for remote SSH access.
 If the parameter is not specified, all IP addresses will be allowed for remote SSH access.

* `pool_id` - (Optional, String, ForceNew) Specifies the ID of Dedicated resource pool which the notebook used.
 Changing this parameter will create a new resource.

* `workspace_id` - (Optional, String, ForceNew) Specifies the workspace ID which the notebook belongs to.
 The default value is `0`. Changing this parameter will create a new resource.

The `volume` block supports:

* `type` - (Required, String, ForceNew) Specifies the volume type. The options are as follows:
  - *EFS*: use Scalable File Service, default 50GB is **free**.
  - *EVS*: use Elastic Volume Service, default size is 5 GB.
  
 Changing this parameter will create a new resource.

* `size` - (Optional, Int) Specifies the volume size. Its value range is from `5` GB to `4,096` GB.

* `ownership` - (Optional, String, ForceNew) Specifies the volume ownership. The options are as follows:
  - *MANAGED*: shared storage disk of the ModelArts service.
  - *DEDICATED*: dedicated storage disk, only supported when the category is `EFS`.

 Changing this parameter will create a new resource.

* `uri` - (Optional, String, ForceNew) Specifies the URL of dedicated storage disk, which is mandatory when the `type`
 is `EFS` and the `ownership` is `DEDICATED`. Example: `192.168.0.1:/user-9sfdsdgdfgh5ea4d56871e75d6966aa274/mount/`.
 Changing this parameter will create a new resource.

* `id` - (Optional, String, ForceNew) Specifies the ID of dedicated storage disk, which is mandatory when the `type`
 is `EFS` and the `ownership` is `DEDICATED`.
 Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `auto_stop_enabled` - Whether enabled the notebook instance to automatically stop.
* `status` - Notebook status. Valid values include: `INIT`, `CREATING`, `STARTING`, `STOPPING`, `DELETING`, `RUNNING`,
 `STOPPED`, `SNAPSHOTTING`, `CREATE_FAILED`, `START_FAILED`, `DELETE_FAILED`, `ERROR`, `DELETED`, `FROZEN`.
* `image_name` - The image name.
* `image_swr_path` - The image path in swr.
* `image_type` - The image type. Valid values include: `BUILD_IN`, `DEDICATED`.
* `created_at` - The notebook creation time.
* `updated_at` - The notebook update time.
* `pool_name` - The name of Dedicated resource pool which the notebook used.
* `url` - The web url of the notebook.
* `ssh_uri` - The uri for remote SSH access.
* `volume/mount_path` - The local mount path of volume.
* `mount_storages` - An array of storages which mount to the notebook. Structure is documented below.

The `mount_storages` block contains:

* `id` - The mount ID.
* `type` - The type of storage which be mounted.
* `path` - The path of storage which be mounted.
* `mount_path` - The local mount path.
* `status` - The status of mount.

## Import

The notebook can be imported by `id`.

```bash
terraform import huaweicloud_modelarts_notebook.test b11b407c-e604-4e8d-8bc4-92398320b847
```

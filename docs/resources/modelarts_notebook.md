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

  tags = {
    foo   = "bar"
    owner = "terraform"
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

* `region` - (Optional, String, ForceNew) Specifies the region where the notebook is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the notebook.  
  The valid length is limited from `1` to `64`, only letters, digits and underscores (_) are allowed.
  The name must starting with a letter.

* `flavor_id` - (Required, String) Specifies the flavor ID of the notebook.

* `image_id` - (Required, String) Specifies the image ID of the notebook.

* `volume` - (Required, List) Specifies the volume configuration of the notebook.  
  The [volume](#modelarts_notebook_volume) structure is documented below.

* `description` - (Optional, String) Specifies the description of the notebook.  
  It contains a maximum of `512` characters and cannot contain special characters `&<>"'/`.

* `key_pair` - (Optional, String, NonUpdatable) Specifies the key pair name for remote SSH access.  
  Required if the parameter `allowed_access_ips` is set.

* `allowed_access_ips` - (Optional, List) Specifies the public IP addresses that are allowed for remote SSH access.  
  If the parameter is not specified, all IP addresses will be allowed for remote SSH access.

* `pool_id` - (Optional, String, NonUpdatable) Specifies the ID of the dedicated resource pool which the notebook used.

* `workspace_id` - (Optional, String, NonUpdatable) Specifies the workspace ID to which the notebook belongs.  
  The default value is **0** (default workspace).

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the notebook.

<a name="modelarts_notebook_volume"></a>
The `volume` block supports:

* `type` - (Required, String, NonUpdatable) Specifies the storage type.  
  The valid values are as follows:
  + **EFS**
  + **EVS**

* `ownership` - (Optional, String, NonUpdatable) Specifies the storage ownership.  
  The valid values are as follows:
  + **MANAGED**: The shared storage disk of the ModelArts service.
  + **DEDICATED**: The dedicated storage disk, only supported when the value of parameter `volume.type` is **EFS**.

  Defaults to **MANAGED**.

* `size` - (Optional, Int) Specifies the storage size, in GB.  
  The valid value range is from `5` to `4,096`.  
  Required if the value of parameter `volume.type` is **EVS** and the value of parameter `volume.ownership` is
  **MANAGED**.

* `uri` - (Optional, String, NonUpdatable) Specifies the storage URL of the dedicated disk.  
  E.g. **9048456e-4e9a-11f1-ba0c-fa16520f410a.c67b108f-4e99-11f1-81fd-fa1652d6f62f.sfsturbo.internal:/**.  
  Required if the value of parameter `volume.type` is **EFS** and the value of parameter `volume.ownership` is
  **DEDICATED**.

* `id` - (Optional, String, NonUpdatable) Specifies the storage ID of the dedicated disk.  
  Required if the value of parameter `volume.type` is **EFS** and the value of parameter `volume.ownership` is
  **DEDICATED**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.

* `auto_stop_enabled` - Whether the notebook auto stop is enabled.

* `status` - The status of the notebook.
  + **RUNNING**
  + **STOPPED**
  + **FROZEN**

* `image_type` - The type of the image which the notebook used.
  + **BUILD_IN**
  + **DEDICATED**

* `image_name` - The name of the image which the notebook used.

* `image_swr_path` - The SWR repository path of the image which the notebook used.

* `created_at` - The creation time of the notebook, in UTC format.

* `updated_at` - The latest update time of the notebook, in UTC format.

* `pool_name` - The name of the dedicated resource pool which the notebook used.

* `url` - The web URL of the notebook.

* `ssh_uri` - The URL for remote SSH access of the notebook.

* `mount_storages` - The storages which are mounted to the notebook.  
  The [mount_storages](#modelarts_notebook_mount_storages_attr) structure is documented below.

<a name="modelarts_notebook_mount_storages_attr"></a>
The `mount_storages` block supports:

* `id` - The ID of the storage which is mounted to the notebook.

* `type` - The type of the storage which is mounted to the notebook.

* `mount_path` - The local mount path of the storage which is mounted to the notebook.

* `path` - The source path of the storage which is mounted to the notebook.

* `status` - The status of the storage which is mounted to the notebook.

## Import

The notebook can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_modelarts_notebook.test <id>
```

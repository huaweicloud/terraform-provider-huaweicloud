---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook_mount_storage"
description: ""
---

# huaweicloud_modelarts_notebook_mount_storage

Manage storages mounted to the notebook resource within HuaweiCloud. A maximum of 10 storages can be mounted.

## Example Usage

```hcl
variable "notebook_id" {}
variable "uri_parallel_obs" {}

resource "huaweicloud_modelarts_notebook_mount_storage" "test" {
  notebook_id           = var.notebook_id
  storage_path          = var.uri_parallel_obs
  local_mount_directory = "/data/test/"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `notebook_id` - (Required, String, ForceNew) Specifies ID of notebook which the storage be mounted to.
 Changing this parameter will create a new resource.

* `storage_path` - (Required, String, ForceNew) Specifies the path of Parallel File System (PFS) or its folders in OBS.
 The format is : `obs://obs-bucket/folder/`. Changing this parameter will create a new resource.

* `local_mount_directory` - (Required, String, ForceNew) Specifies the local mount directory.
  Only the subdirectory of `/data/` can be mounted. The format is : `/data/dir1/`.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **notebook_id/mount_id**. It is composed of the ID of notebook and mount ID,
 separated by a slash.

* `mount_id` - The mount ID.

* `status` - mount status. Valid values include: `MOUNTING`, `MOUNT_FAILED`, `MOUNTED`, `UNMOUNTING`,
 `UNMOUNT_FAILED`, `UNMOUNTED`.

* `type` - The type of storage system.  The value is `OBSFS`.

## Import

The mount storage can be imported by `id`, It is composed of the ID of notebook and mount ID, separated by a slash.
 For example,

```bash
terraform import huaweicloud_modelarts_notebook_mount_storage.test b11b407c-e604-4e8d-8bc4-92398320b847/4e206d3c-6831-4267-b93d-e236105cda38
```

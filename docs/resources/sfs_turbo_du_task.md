---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_du_task"
description: |-
  Manages a DU task resource under the SFS Turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_du_task

Manages a DU task resource under the SFS Turbo within HuaweiCloud.

-> Only SFS Turbo File systems created after August 1, 2023 can use this resource. Creating new task is not supported
  when there are `10` or more executing tasks under the system. In order not to affect the performance of the file
  system, it is recommended that the number of tasks submitted simultaneously under one system should not exceed `4`.

-> This resource is only available for the following SFS Turbo file system types:
  **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,**250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC缓存型**.

## Example Usage

```hcl
variable "share_id" {}
variable "path" {}

resource "huaweicloud_sfs_turbo_du_task" "test" {
  share_id = var.share_id
  path     = var.path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, ForceNew) Specifies the ID of the SFS Turbo file system to which the DU task belongs.
  Changing this creates a new resource.

* `path` - (Required, String, ForceNew) Specifies the full path to a legal directory in the file system.
  The length of a single level directory is not allowed to exceed `255`, and the length of the full path is not allowed
  to exceed `4,096`.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the DU task. The valid values are as following:
  + **SUCCESS**
  + **DOING**
  + **FAIL**

* `dir_usage` - The direcrory resource usage (including subdirectories).
  The [dir_usage](#du_task_dirUsage) structure is documented below.

* `begin_time` - The start time of the DU task, in RFC3339 format.

* `end_time` - The end time of the DU task, in RFC3339 format.

<a name="du_task_dirUsage"></a>
The `dir_usage` block supports:

* `path` - The full path to a legal directory in the file system.

* `used_capacity` - The used capacity, in byte.

* `file_count` - The total number of files in the directory.
  The [file_count](#dir_usage_fileCount) structure is documented below.

* `message` - The error messages.

<a name="dir_usage_fileCount"></a>
The `file_count` block supports:

* `dir` - The number of directories.

* `regular` - The number of common files.

* `pipe` - The number of pipe files.

* `char` - The number of character devices.

* `block` - The number of block devices.

* `socket` - The number of sockets.

* `symlink` - The number of symbolic links.

## Import

The SFS Turbo DU task can be imported using the related `share_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_sfs_turbo_du_task.test <share_id>/<id>
```

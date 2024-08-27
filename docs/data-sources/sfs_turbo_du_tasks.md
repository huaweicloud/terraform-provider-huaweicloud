---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_du_tasks"
description: |-
  Use this datasource to get list of the DU tasks.
---

# huaweicloud_sfs_turbo_du_tasks

Use this datasource to get list of the DU tasks.

-> Only SFS Turbo File systems created after August 1, 2023 can use this datasource.

-> This datasource is only available for the following SFS Turbo file system types:
  **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,**250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC缓存型**.

## Example Usage

```hcl
variable "share_id" {}

data "huaweicloud_sfs_turbo_du_tasks" "test" {
  share_id = var.share_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS Turbo file system to which the DU tasks belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of DU tasks.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - The ID of the DU task.

* `status` - The status of the DU task.
  The value can be **SUCCESS**, **DOING**, or **FAIL**.

* `dir_usage` - The resource usages of a directory (including subdirectories).

  The [dir_usage](#tasks_dir_usage_struct) structure is documented below.

* `begin_time` - The start time of the DU task, in RFC3339 format.

* `end_time` - The end time of the DU task, in RFC3339 format.

<a name="tasks_dir_usage_struct"></a>
The `dir_usage` block supports:

* `path` - The full path to a legal directory in the file system.

* `used_capacity` - The used capacity, in byte.

* `message` - The error message.

* `file_count` - The total number of files in the directory.

  The [file_count](#dir_usage_file_count_struct) structure is documented below.

<a name="dir_usage_file_count_struct"></a>
The `file_count` block supports:

* `dir` - The number of directories.

* `regular` - The number of common files.

* `char` - The number of character devices.

* `block` - The number of block devices.

* `pipe` - The number of pipe files.

* `socket` - The number of sockets.

* `symlink` - The number of symbolic links.

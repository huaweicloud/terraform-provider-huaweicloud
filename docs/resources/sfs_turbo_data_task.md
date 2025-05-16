---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_data_task"
description: |-
  Manages a data import or export task resource under the SFS Turbo within HuaweiCloud.
---

# huaweicloud_sfs_turbo_data_task

Manages a data import or export task resource under the SFS Turbo within HuaweiCloud.

-> Please pay attention to the results of task execution through `status`.

-> This resource is only available for the following SFS Turbo file system types:
  **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,**250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC缓存型**.

## Example Usage

```hcl
variable "share_id" {}
variable "type" {}
variable "src_target" {}

resource "huaweicloud_sfs_turbo_data_task" "test" {
  share_id    = var.share_id
  type        = var.type
  src_target  = var.src_target
  dest_target = var.src_target
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `share_id` - (Required, String, ForceNew) Specifies the ID of the SFS Turbo file system to which the data task
  belongs. Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the data task.
  Changing this creates a new resource.
  The valid values are as following:
  + **import**: Additional metadata import. This import method will import the metadata (name, size, and latest update
    time) of the OBS object and the additional metadata (such as uid, gid and mode) exported from SFS Turbo.
  + **import_metadata**: Quick import. In this mode, only the metadata (name, size, and latest update time) of the OBS
    object will be imported, other additional metadata (such as uid, gid and mode) is not imported, SFS Turbo will
    generate the default additional metadata (uid: `0`, gid: `0`, directory permission: `755`, file permission: `644`).
  + **preload**: Data preheat. Both metadata and data content will be imported. Metadata is imported in quick mode,
    other additional metadata (such as uid, gid and mode) is not imported.
  + **export**: Data export. This method will export the files which created in the linkage directory or which have
    been modified after imported from OBS to the OBS bucket.

* `src_target` - (Required, String, ForceNew) Specifies the linkage directory name.
  Changing this creates a new resource.

* `dest_target` - (Required, String, ForceNew) Specifies target end information of the data task.
  Changing this creates a new resource.
  Currently, the value only support keep the same as `src_target`.

* `src_prefix` - (Optional, String, ForceNew) Specifies source path prefix of the data task.
  Changing this creates a new resource.
  + For an import task, do not include the OBS bucket name.
  + For an export task, do not include the name of the linkage directory.
  + For a preload task, the source path prefix must be a directory or specific object ended with a slash (/).
  + If not specified, all objects in the bound OBS bucket will be imported when importing, all files in the linkage
    directory will be exported when exporting.

* `dest_prefix` - (Optional, String, ForceNew) Specifies destination path prefix of the data task.
  Changing this creates a new resource.
  Currently, the value only support keep the same as `src_prefix`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the data task. The valid values are as following:
  + **SUCCESS**
  + **FAIL**

* `message` - The execution result information of the data task.

* `start_time` - The start time of the data task, in RFC3339 format.

* `end_time` - The end time of the data task, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

* `delete` - Default is 5 minutes.

## Import

The SFS Turbo data task can be imported using the related `share_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_sfs_turbo_data_task.test <share_id>/<id>
```

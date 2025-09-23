---
subcategory: "SFS Turbo"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_turbo_data_tasks"
description: |-
  Use this datasource to get list of the data import or export tasks.
---

# huaweicloud_sfs_turbo_data_tasks

Use this datasource to get list of the data import or export tasks.

-> This datasource is only available for the following SFS Turbo file system types:
  **20MB/s/TiB**, **40MB/s/TiB**, **125MB/s/TiB**,**250MB/s/TiB**, **500MB/s/TiB**, **1,000MB/s/TiB**, **HPC缓存型**.

## Example Usage

```hcl
variable "share_id" {}

data "huaweicloud_sfs_turbo_data_tasks" "test" {
  share_id = var.share_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `share_id` - (Required, String) Specifies the ID of the SFS Turbo file system to which the data tasks belong.

* `type` - (Optional, String) Specifies the type of the data task.
  The valid values are as following:
  + **import**: Additional metadata import.
  + **import_metadata**: Quick import.
  + **preload**: Data preheat.
  + **export**: Data export.

* `status` - (Optional, String) Specifies the status of the data task.
  The value can be **SUCCESS** or **FAIL**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the data task.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - The ID of the data task.

* `type` - The type of the data task.

* `status` - The status of the data task.

* `src_target` - The linkage directory name.

* `dest_target` - The target end information of the data task.

* `src_prefix` - The source path prefix of the data task.

* `dest_prefix` - The destination path prefix of the data task.

* `message` - The data task execution result information.

* `start_time` - The start time of the data task, in RFC3339 format.

* `end_time` - The end time of the data task, in RFC3339 format.

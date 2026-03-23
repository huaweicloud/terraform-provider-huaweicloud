---
subcategory: "MapReduce Service (MRS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_mapreduce_cluster_files"
description: |-
  Use this data source to query the file list under the specified directory within HuaweiCloud.
---

# huaweicloud_mapreduce_cluster_files

Use this data source to query the file list under the specified directory within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}
variable "file_path" {}

data "huaweicloud_mapreduce_cluster_files" "test" {
  cluster_id = var.cluster_id
  path       = var.file_path
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the cluster files are located.  
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the ID of the cluster.

* `path` - (Required, String) Specifies the directory path of the file.  
  The maximum length is `255` characters.  
  This parameter cannot start and end with a dot (.), and cannot include special characters (``*?"<>|;&,'`!{}[]$%+``).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `files` - The list of files under the specified path.  
  The [files](#cluster_files) structure is documented below.

<a name="cluster_files"></a>
The `files` block supports:

* `path_suffix` - The path suffix of the file under the queried directory.

* `owner` - The owner of the file.

* `group` - The group to which the file belongs.

* `permission` - The permission of the file.

* `replication` - The replication factor of the file.

* `block_size` - The block size of the file.

* `length` - The length of the file.

* `type` - The type of the file.
  - **FILE**
  - **DIRECTORY**

* `children_num` - The number of entries under the queried directory.

* `access_time` - The time when the file was accessed, in RFC3339 format.

* `modification_time` - The time when the file was modified, in RFC3339 format.

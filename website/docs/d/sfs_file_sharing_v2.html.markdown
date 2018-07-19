---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_file_sharing_v2"
sidebar_current: "docs-huaweicloud-sfs-file-sharing-v2"
description: |-
  Get information on an Huawei Cloud shared file system.
---

# Data Source: huaweicloud_sfs_file_sharing_v2

Provides information about an Shared File System (SFS).

## Example Usage

```hcl
    variable "share_name" { }

    variable "share_id" { }

    data "huaweicloud_sfs_file_sharing_v2" "shared_file" 
    {
           name = "${var.share_name}"
           id = "${var.share_id}"
    }

```

## Argument Reference
The following arguments are supported:

* `name` - (Optional) The name of the shared file system.

* `id` - (Optional) The UUID of the shared file system.

* `status` - (Optional) The status of the shared file system.


## Attributes Reference

The following attributes are exported:

* `availability_zone` - The availability zone name.

* `size` - 	The size (GB) of the shared file system.

* `share_type` - The storage service type for the shared file system, such as high-performance storage (composed of SSDs) or large-capacity storage (composed of SATA disks).

* `status` - The status of the shared file system.

* `host` - The host name of the shared file system.

* `is_public` - The level of visibility for the shared file system.

* `share_proto` - The protocol for sharing file systems.

* `volume_type` - The volume type.

* `metadata` - Metadata key and value pairs as a dictionary of strings.

* `export_location` - The path for accessing the shared file system.

* `access_level` - The level of the access rule.

* `access_rules_status` - The status of the share access rule.

* `access_type` - The type of the share access rule.

* `access_to` - The access that the back end grants or denies.

* `share_access_id` - The UUID of the share access rule.

* `mount_id` - The UUID of the mount location of the shared file system.

* `share_instance_id` - The access that the back end grants or denies.

* `preferred` - Identifies which mount locations are most efficient and are used preferentially when multiple mount locations exist.



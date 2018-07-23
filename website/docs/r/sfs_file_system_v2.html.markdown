---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_file_system_v2"
sidebar_current: "docs-huaweicloud-resource-sfs-file-system-v2"
description: |-
   Provides an Scalable File Resource (SFS) resource.
---

# huaweicloud_sfs_file_system_v2

Provides an Shared File System (SFS) resource.

## Example Usage

 ```hcl
    variable "share_name" { }

    variable "share_description" { }

    variable "vpc_id" { }

    resource "huaweicloud_sfs_file_system_v2" "sfs1" {
            size = 50
            name = "${var.share_name}"
            access_to = "${var.vpc_id}"
            access_level = "rw"
            description = "${var.share_description}"
            metadata = {
                "type"="nfs"
            }
    }
 ```

## Argument Reference
The following arguments are supported:

* `size` - (Required) The size (GB) of the shared file system.

* `share_proto` - (Optional) The protocol for sharing file systems. The default value is NFS.

* `name` - (Optional) The name of the shared file system.

* `description` - (Optional) Describes the shared file system.

* `is_public` - (Optional) The level of visibility for the shared file system.

* `metadata` - (Optional) Metadata key and value pairs as a dictionary of strings.Changing this will create a new resource.

* `availability_zone` - (Optional) The availability zone name.Changing this parameter will create a new resource.

* `access_level` - (Required) The access level of the shared file system. Changing this will create a new access rule.

* `access_type` - (Optional) The type of the share access rule. Changing this will create a new access rule.

* `access_to` - (Required) The access that the back end grants or denies. Changing this will create a new access rule

## Attributes Reference
In addition to all arguments above, the following attributes are exported:

* `id` - The UUID of the shared file system.

* `status` - The status of the shared file system.

* `share_type` - The storage service type assigned for the shared file system, such as high-performance storage (composed of SSDs) and large-capacity storage (composed of SATA disks).

* `volume_type` - The volume type.

* `export_location` - The address for accessing the shared file system.

* `host` - The host name of the shared file system.

* `share_access_id` - The UUID of the share access rule.

* `access_rules_status` - The status of the share access rule.

## Import

SFS can be imported using the `id`, e.g.

```
> $ terraform import huaweicloud_sfs_file_system_v2 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```









   
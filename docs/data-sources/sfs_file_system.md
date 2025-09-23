---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sfs_file_system"
description: ""
---

# huaweicloud_sfs_file_system

!> **WARNING:** It has been deprecated.

Provides information about an Shared File System (SFS) within HuaweiCloud.

## Example Usage

```hcl
variable "share_name" {}

data "huaweicloud_sfs_file_system" "shared_file" {
  name = var.share_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the shared file system.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) The name of the shared file system.

* `id` - (Optional, String) The UUID of the shared file system.

* `status` - (Optional, String) The status of the shared file system.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The availability zone name.

* `description` - The description of the shared file system.

* `state` - The state of the shared file system.

* `size` - The size (GB) of the shared file system.

* `is_public` - Whether a file system can be publicly seen.

* `share_proto` - The protocol for sharing file systems.

* `metadata` - The key and value pairs information of the shared file system.

* `export_location` - The path for accessing the shared file system.

* `share_access_id` - The UUID of the share access rule.

* `access_rules_status` - The status of the share access rule.

* `access_level` - The level of the access rule.

* `access_type` - The type of the share access rule.

* `access_to` - The access that the back end grants or denies.

* `mount_id` - The UUID of the mount location of the shared file system.

* `share_instance_id` - The access that the back end grants or denies.

* `preferred` - Identifies which mount locations are most efficient and are used preferentially when multiple mount
  locations exist.

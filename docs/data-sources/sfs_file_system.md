---
subcategory: "Scalable File Service (SFS)"
---

# huaweicloud_sfs_file_system

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

* `name` - (Optional, String) The name of the shared file system.

* `id` - (Optional, String) The UUID of the shared file system.

* `status` - (Optional, String) The status of the shared file system.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The availability zone name.

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
